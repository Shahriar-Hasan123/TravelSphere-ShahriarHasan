// This layer is responsible only for making HTTP requests and decoding

package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const restCountriesDefaultURL = "https://api.restcountries.com/countries/v5"

// RawCountry holds the raw fields we extract from the REST Countries API.
type RawCountry struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	Flags struct {
		PNG string `json:"png"`
		SVG string `json:"svg"`
	} `json:"flags"`
	Capital    []string `json:"capital"`
	Population int64    `json:"population"`
	Region     string   `json:"region"`
	Subregion  string   `json:"subregion"`
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Languages map[string]string `json:"languages"`
	Latlng    []float64         `json:"latlng"`
}

// RestCountriesClient makes requests to the REST Countries API.
type RestCountriesClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

// NewRestCountriesClientWithURL creates a client with a custom base URL — used in tests.
func NewRestCountriesClientWithURL(baseURL string) *RestCountriesClient {
	return &RestCountriesClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		apiKey:     os.Getenv("RESTCOUNTRIES_API_KEY"),
	}
}

// NewRestCountriesClient creates a client with a sensible timeout.
func NewRestCountriesClient() *RestCountriesClient {
	baseURL := os.Getenv("RESTCOUNTRIES_BASE_URL")
	if baseURL == "" {
		baseURL = restCountriesDefaultURL
	}
	return &RestCountriesClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     os.Getenv("RESTCOUNTRIES_API_KEY"),
	}
}

// FetchAll retrieves all countries with only the fields we need.
func (c *RestCountriesClient) FetchAll() ([]RawCountry, error) {
	fields := "name,flags,capital,population,region,subregion,currencies,languages,latlng"

	// v5 requires pagination
	if strings.Contains(c.baseURL, "/v5") {
		var allCountries []RawCountry
		for offset := 0; offset < 300; offset += 25 {
			url := fmt.Sprintf("%s?fields=%s&offset=%d", c.baseURL, fields, offset)
			page, err := c.fetch(url)
			if err != nil {
				return nil, err
			}
			if len(page) == 0 {
				break // No more results
			}
			allCountries = append(allCountries, page...)
		}
		return allCountries, nil
	}

	// v3 format
	url := fmt.Sprintf("%s/all?fields=%s", c.baseURL, fields)
	return c.fetch(url)
}

// FetchByName retrieves countries whose common name matches the query.
func (c *RestCountriesClient) FetchByName(name string) ([]RawCountry, error) {
	fields := "name,flags,capital,population,region,subregion,currencies,languages,latlng"
	var url string
	if strings.Contains(c.baseURL, "/v5") {
		// v5 variant: use query param 'name' on the base endpoint.
		url = fmt.Sprintf("%s?name=%s&fields=%s", c.baseURL, name, fields)
	} else {
		url = fmt.Sprintf("%s/name/%s?fields=%s", c.baseURL, name, fields)
	}
	return c.fetch(url)
}

// fetch is the shared internal method that executes the HTTP GET and
func (c *RestCountriesClient) fetch(url string) ([]RawCountry, error) {
	log.Printf("[RestCountries] GET %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("restcountries request failed: %w", err)
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("restcountries request failed: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("[RestCountries] Status: %d", resp.StatusCode)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("restcountries returned status %d", resp.StatusCode)
	}

	// Read full body so we can attempt multiple decode strategies (v3 array vs v5 wrapper).
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("restcountries read failed: %w", err)
	}

	// Try v3 format: direct array
	var countries []RawCountry
	if err := json.Unmarshal(body, &countries); err == nil {
		log.Printf("[RestCountries] Decoded v3 format: %d countries", len(countries))
		return countries, nil
	}

	// Try v5 format: { "data": { "objects": [...] } }
	var wrapperV5 struct {
		Data struct {
			Objects []struct {
				Names struct {
					Common   string `json:"common"`
					Official string `json:"official"`
				} `json:"names"`
				Flag struct {
					URLPNG string `json:"url_png"`
					URLSVG string `json:"url_svg"`
				} `json:"flag"`
				Capitals []struct {
					Name        string `json:"name"`
					Coordinates struct {
						Lat float64 `json:"lat"`
						Lng float64 `json:"lng"`
					} `json:"coordinates"`
				} `json:"capitals"`
				Population int64  `json:"population"`
				Region     string `json:"region"`
				Subregion  string `json:"subregion"`
				Currencies []struct {
					Code   string `json:"code"`
					Name   string `json:"name"`
					Symbol string `json:"symbol"`
				} `json:"currencies"`
				Languages []struct {
					Name string `json:"name"`
				} `json:"languages"`
				Coordinates struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"coordinates"`
				Latlng []float64 `json:"latlng"`
			} `json:"objects"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapperV5); err == nil {
		v5 := wrapperV5.Data.Objects
		log.Printf("[RestCountries] Decoded v5 format: %d countries", len(v5))
		mapped := make([]RawCountry, 0, len(v5))
		for _, v := range v5 {
			var rc RawCountry
			// names
			rc.Name.Common = v.Names.Common
			rc.Name.Official = v.Names.Official
			// flags
			rc.Flags.PNG = v.Flag.URLPNG
			rc.Flags.SVG = v.Flag.URLSVG
			// capitals
			if len(v.Capitals) > 0 {
				rc.Capital = []string{v.Capitals[0].Name}
			}
			// population/region/subregion
			rc.Population = v.Population
			rc.Region = v.Region
			rc.Subregion = v.Subregion
			// currencies -> map
			rc.Currencies = make(map[string]struct {
				Name   string `json:"name"`
				Symbol string `json:"symbol"`
			})
			for _, cur := range v.Currencies {
				key := cur.Code
				if key == "" {
					key = cur.Name
				}
				rc.Currencies[key] = struct {
					Name   string `json:"name"`
					Symbol string `json:"symbol"`
				}{Name: cur.Name, Symbol: cur.Symbol}
			}
			// languages -> map
			rc.Languages = make(map[string]string)
			for i, lang := range v.Languages {
				rc.Languages[fmt.Sprintf("l%d", i)] = lang.Name
			}
			// latlng
			if len(v.Latlng) == 2 {
				rc.Latlng = v.Latlng
			} else if v.Coordinates.Lat != 0 || v.Coordinates.Lng != 0 {
				rc.Latlng = []float64{v.Coordinates.Lat, v.Coordinates.Lng}
			}
			mapped = append(mapped, rc)
		}
		return mapped, nil
	}

	return nil, fmt.Errorf("restcountries decode failed: %w", err)
}
