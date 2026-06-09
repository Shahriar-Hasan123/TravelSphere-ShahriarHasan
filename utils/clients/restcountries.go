// This layer is responsible only for making HTTP requests and decoding

package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const restCountriesDefaultURL = "https://restcountries.com/v3.1"

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
}

// NewRestCountriesClientWithURL creates a client with a custom base URL — used in tests.
func NewRestCountriesClientWithURL(baseURL string) *RestCountriesClient {
	return &RestCountriesClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
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
	}
}

// FetchAll retrieves all countries with only the fields we need.
func (c *RestCountriesClient) FetchAll() ([]RawCountry, error) {
	fields := "name,flags,capital,population,region,subregion,currencies,languages,latlng"
	url := fmt.Sprintf("%s/all?fields=%s", c.baseURL, fields)
	return c.fetch(url)
}

// FetchByName retrieves countries whose common name matches the query.
func (c *RestCountriesClient) FetchByName(name string) ([]RawCountry, error) {
	fields := "name,flags,capital,population,region,subregion,currencies,languages,latlng"
	url := fmt.Sprintf("%s/name/%s?fields=%s", c.baseURL, name, fields)
	return c.fetch(url)
}

// fetch is the shared internal method that executes the HTTP GET and
func (c *RestCountriesClient) fetch(url string) ([]RawCountry, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("restcountries request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil 
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("restcountries returned status %d", resp.StatusCode)
	}

	var countries []RawCountry
	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		return nil, fmt.Errorf("restcountries decode failed: %w", err)
	}
	return countries, nil
}
