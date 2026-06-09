package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

const openTripMapDefaultURL = "https://api.opentripmap.com/0.1/en"

// Kinds is the set of OpenTripMap category filters we request.
const attractionKinds = "interesting_places,historic,cultural,architecture,museums,natural"

// RawPlace holds the fields we extract from a single OpenTripMap place result.
type RawPlace struct {
	Name  string  `json:"name"`
	Kinds string  `json:"kinds"` // Comma-separated category tags
	Dist  float64 `json:"dist"`  // Distance from centre in metres
	Xid   string  `json:"xid"`   // Unique place identifier
}

// rawPlacesResponse wraps the /radius endpoint response.
type rawPlacesResponse struct {
	Features []struct {
		Properties RawPlace `json:"properties"`
	} `json:"features"`
}

// OpenTripMapClient makes requests to the OpenTripMap places API.
type OpenTripMapClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewOpenTripMapClient creates a client reading credentials from the environment.
func NewOpenTripMapClient() *OpenTripMapClient {
	baseURL := os.Getenv("OPENTRIPMAP_BASE_URL")
	if baseURL == "" {
		baseURL = openTripMapDefaultURL
	}
	return &OpenTripMapClient{
		baseURL:    baseURL,
		apiKey:     os.Getenv("OPENTRIPMAP_API_KEY"),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// NewOpenTripMapClientWithURL creates a client with a custom base URL — used in tests.
func NewOpenTripMapClientWithURL(baseURL, apiKey string) *OpenTripMapClient {
	return &OpenTripMapClient{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// FetchAttractionsByCoords fetches tourist places within `radiusM` metres of the given coordinates. Returns up to `limit` results.
func (c *OpenTripMapClient) FetchAttractionsByCoords(
	lat, lon float64,
	radiusM, limit int,
) ([]RawPlace, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("OPENTRIPMAP_API_KEY is not set")
	}

	endpoint := fmt.Sprintf(
		"%s/places/radius?apikey=%s&radius=%d&lon=%.6f&lat=%.6f&kinds=%s&limit=%d&format=geojson",
		c.baseURL,
		url.QueryEscape(c.apiKey),
		radiusM,
		lon, lat,
		attractionKinds,
		limit,
	)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("opentripmap request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("opentripmap returned status %d", resp.StatusCode)
	}

	var result rawPlacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("opentripmap decode failed: %w", err)
	}

	places := make([]RawPlace, 0, len(result.Features))
	for _, f := range result.Features {
		places = append(places, f.Properties)
	}
	return places, nil
}
