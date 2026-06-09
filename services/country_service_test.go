package services

import (
	"TravelSphere/utils/clients"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockCountryServer(countries []clients.RawCountry) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(countries)
	}))
}

func buildRawCountry(name, official, capital, region, subregion string, pop int64) clients.RawCountry {
	c := clients.RawCountry{}
	c.Name.Common = name
	c.Name.Official = official
	c.Capital = []string{capital}
	c.Region = region
	c.Subregion = subregion
	c.Population = pop
	c.Flags.PNG = "https://flag.example.com/" + name + ".png"
	c.Latlng = []float64{41.0, 20.0}
	c.Currencies = map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	}{"ALL": {Name: "Albanian lek", Symbol: "L"}}
	c.Languages = map[string]string{"sqi": "Albanian"}
	return c
}

func newTestCountryService(server *httptest.Server) *CountryService {
	return &CountryService{
		client: &clients.RestCountriesClient{},
	}
}

func TestGetAllCountries_NoFilter(t *testing.T) {
	raw := []clients.RawCountry{
		buildRawCountry("Albania", "Republic of Albania", "Tirana", "Europe", "Southeast Europe", 2837743),
		buildRawCountry("France", "French Republic", "Paris", "Europe", "Western Europe", 67000000),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(raw)
	}))
	defer server.Close()

	svc := &CountryService{
		client: &clients.RestCountriesClient{},
	}
	// Directly inject the mock server URL into the client.
	svc.client = clients.NewRestCountriesClientWithURL(server.URL)

	countries, err := svc.GetAllCountries("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(countries) != 2 {
		t.Errorf("expected 2 countries, got %d", len(countries))
	}
}

func TestGetAllCountries_SearchFilter(t *testing.T) {
	raw := []clients.RawCountry{
		buildRawCountry("Albania", "Republic of Albania", "Tirana", "Europe", "Southeast Europe", 2837743),
		buildRawCountry("France", "French Republic", "Paris", "Europe", "Western Europe", 67000000),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(raw)
	}))
	defer server.Close()

	svc := &CountryService{client: clients.NewRestCountriesClientWithURL(server.URL)}

	countries, err := svc.GetAllCountries("alb", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(countries) != 1 || countries[0].Name != "Albania" {
		t.Errorf("expected only Albania, got %v", countries)
	}
}

func TestGetAllCountries_RegionFilter(t *testing.T) {
	raw := []clients.RawCountry{
		buildRawCountry("Albania", "Republic of Albania", "Tirana", "Europe", "Southeast Europe", 2837743),
		buildRawCountry("Japan", "Japan", "Tokyo", "Asia", "Eastern Asia", 125000000),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(raw)
	}))
	defer server.Close()

	svc := &CountryService{client: clients.NewRestCountriesClientWithURL(server.URL)}

	countries, err := svc.GetAllCountries("", "Asia")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(countries) != 1 || countries[0].Name != "Japan" {
		t.Errorf("expected only Japan, got %v", countries)
	}
}

func TestGetCountryBySlug_Found(t *testing.T) {
	raw := []clients.RawCountry{
		buildRawCountry("Albania", "Republic of Albania", "Tirana", "Europe", "Southeast Europe", 2837743),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(raw)
	}))
	defer server.Close()

	svc := &CountryService{client: clients.NewRestCountriesClientWithURL(server.URL)}

	country, err := svc.GetCountryBySlug("albania")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if country == nil {
		t.Fatal("expected country, got nil")
	}
	if country.Name != "Albania" {
		t.Errorf("expected Albania, got %q", country.Name)
	}
}

func TestGetCountryBySlug_NotFound(t *testing.T) {
	raw := []clients.RawCountry{
		buildRawCountry("Albania", "Republic of Albania", "Tirana", "Europe", "Southeast Europe", 2837743),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(raw)
	}))
	defer server.Close()

	svc := &CountryService{client: clients.NewRestCountriesClientWithURL(server.URL)}

	country, err := svc.GetCountryBySlug("zzz-invalid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if country != nil {
		t.Errorf("expected nil for unknown slug, got %v", country)
	}
}

func TestSearchSuggestions_ReturnsMax10(t *testing.T) {
	raw := make([]clients.RawCountry, 15)
	for i := range raw {
		raw[i] = buildRawCountry(
			"Country"+string(rune('A'+i)), "", "Capital", "Europe", "", 1000000,
		)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(raw)
	}))
	defer server.Close()

	svc := &CountryService{client: clients.NewRestCountriesClientWithURL(server.URL)}

	suggestions, err := svc.SearchSuggestions("country")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(suggestions) > 10 {
		t.Errorf("expected max 10 suggestions, got %d", len(suggestions))
	}
}

func TestSearchSuggestions_EmptyQuery(t *testing.T) {
	svc := &CountryService{client: clients.NewRestCountriesClientWithURL("http://unused")}

	suggestions, err := svc.SearchSuggestions("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if suggestions != nil {
		t.Error("expected nil for empty query")
	}
}

func TestGetFeaturedCountries_Success(t *testing.T) {
	raw := []clients.RawCountry{
		buildRawCountry("United States", "United States of America", "Washington D.C.", "Americas", "North America", 331000000),
		buildRawCountry("France", "French Republic", "Paris", "Europe", "Western Europe", 67000000),
		buildRawCountry("Japan", "Japan", "Tokyo", "Asia", "Eastern Asia", 125000000),
		buildRawCountry("Australia", "Australia", "Canberra", "Oceania", "Australia and New Zealand", 25000000),
		buildRawCountry("Brazil", "Federative Republic of Brazil", "Brasília", "Americas", "South America", 213000000),
		buildRawCountry("Bangladesh", "People's Republic of Bangladesh", "Dhaka", "Asia", "Southern Asia", 165000000),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(raw)
	}))
	defer server.Close()

	svc := &CountryService{client: clients.NewRestCountriesClientWithURL(server.URL)}

	featured, err := svc.GetFeaturedCountries()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(featured) == 0 {
		t.Error("expected featured countries, got empty slice")
	}
}

func TestGetFeaturedCountries_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer server.Close()

	svc := &CountryService{client: clients.NewRestCountriesClientWithURL(server.URL)}

	_, err := svc.GetFeaturedCountries()
	if err == nil {
		t.Error("expected error from API failure")
	}
}

func TestGetAllCountries_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer server.Close()

	svc := &CountryService{client: clients.NewRestCountriesClientWithURL(server.URL)}

	_, err := svc.GetAllCountries("", "")
	if err == nil {
		t.Error("expected error from API failure")
	}
}

func TestGetCountryBySlug_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer server.Close()

	svc := &CountryService{client: clients.NewRestCountriesClientWithURL(server.URL)}

	_, err := svc.GetCountryBySlug("france")
	if err == nil {
		t.Error("expected error from API failure")
	}
}

func TestGetAllCountries_SortedAlphabetically(t *testing.T) {
	raw := []clients.RawCountry{
		buildRawCountry("Zimbabwe", "Republic of Zimbabwe", "Harare", "Africa", "Southern Africa", 15000000),
		buildRawCountry("Albania", "Republic of Albania", "Tirana", "Europe", "Southeast Europe", 2837743),
		buildRawCountry("France", "French Republic", "Paris", "Europe", "Western Europe", 67000000),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(raw)
	}))
	defer server.Close()

	svc := &CountryService{client: clients.NewRestCountriesClientWithURL(server.URL)}

	countries, err := svc.GetAllCountries("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if countries[0].Name != "Albania" {
		t.Errorf("expected Albania first, got %q", countries[0].Name)
	}
	if countries[2].Name != "Zimbabwe" {
		t.Errorf("expected Zimbabwe last, got %q", countries[2].Name)
	}
}

func TestGetAllCountries_CapitalSearchFilter(t *testing.T) {
	raw := []clients.RawCountry{
		buildRawCountry("France", "French Republic", "Paris", "Europe", "Western Europe", 67000000),
		buildRawCountry("Japan", "Japan", "Tokyo", "Asia", "Eastern Asia", 125000000),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(raw)
	}))
	defer server.Close()

	svc := &CountryService{client: clients.NewRestCountriesClientWithURL(server.URL)}

	// Search by capital city name.
	countries, err := svc.GetAllCountries("paris", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(countries) != 1 || countries[0].Name != "France" {
		t.Errorf("expected France by capital search, got %v", countries)
	}
}

func TestNewCountryService_NotNil(t *testing.T) {
	svc := NewCountryService()
	if svc == nil {
		t.Error("expected non-nil CountryService")
	}
}

func TestToDTO_WithAllFields(t *testing.T) {
	// Test toDTO with complete data
	raw := buildRawCountry("France", "French Republic", "Paris", "Europe", "Western Europe", 67000000)
	raw.Flags.SVG = "" // Only PNG flag

	svc := NewCountryService()
	dto := svc.toDTO(raw)

	if dto.Name != "France" {
		t.Errorf("expected France, got %q", dto.Name)
	}
	if dto.Capital != "Paris" {
		t.Errorf("expected Paris, got %q", dto.Capital)
	}
	if dto.Region != "Europe" {
		t.Errorf("expected Europe, got %q", dto.Region)
	}
	if dto.Population != 67000000 {
		t.Errorf("expected population 67000000, got %d", dto.Population)
	}
	if dto.Flag != "https://flag.example.com/France.png" {
		t.Errorf("expected PNG flag, got %q", dto.Flag)
	}
}

func TestToDTO_NoCapitals(t *testing.T) {
	// Test toDTO when capital list is empty
	raw := buildRawCountry("TestCountry", "Test Country Official", "Capital", "Region", "Subregion", 1000000)
	raw.Capital = []string{} // Empty capital list

	svc := NewCountryService()
	dto := svc.toDTO(raw)

	if dto.Capital != "" {
		t.Errorf("expected empty capital, got %q", dto.Capital)
	}
}

func TestToDTO_NoFlagPNG_UsesSVG(t *testing.T) {
	// Test toDTO fallback to SVG when PNG is empty
	raw := buildRawCountry("TestCountry", "Official", "Capital", "Region", "Sub", 1000000)
	raw.Flags.PNG = ""
	raw.Flags.SVG = "https://flag.example.com/test.svg"

	svc := NewCountryService()
	dto := svc.toDTO(raw)

	if dto.Flag != "https://flag.example.com/test.svg" {
		t.Errorf("expected SVG flag fallback, got %q", dto.Flag)
	}
}

func TestToDTO_NoCurrencies(t *testing.T) {
	// Test toDTO with empty currencies map
	raw := buildRawCountry("TestCountry", "Official", "Capital", "Region", "Sub", 1000000)
	raw.Currencies = make(map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	})

	svc := NewCountryService()
	dto := svc.toDTO(raw)

	if dto.Currency != "" {
		t.Errorf("expected empty currency, got %q", dto.Currency)
	}
}

func TestToDTO_MultipleLanguages(t *testing.T) {
	// Test toDTO with multiple languages (should be sorted)
	raw := buildRawCountry("TestCountry", "Official", "Capital", "Region", "Sub", 1000000)
	raw.Languages = map[string]string{
		"eng": "English",
		"fra": "French",
		"ger": "German",
	}

	svc := NewCountryService()
	dto := svc.toDTO(raw)

	if len(dto.Languages) != 3 {
		t.Errorf("expected 3 languages, got %d", len(dto.Languages))
	}
	// Languages should be sorted
	if dto.Languages[0] != "English" {
		t.Errorf("expected English first (sorted), got %q", dto.Languages[0])
	}
}

func TestSearchSuggestions_WithError(t *testing.T) {
	// Test SearchSuggestions error propagation
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer server.Close()

	svc := &CountryService{client: clients.NewRestCountriesClientWithURL(server.URL)}

	suggestions, err := svc.SearchSuggestions("test")
	if err == nil {
		t.Error("expected error from API failure")
	}
	if suggestions != nil {
		t.Error("expected nil suggestions on error")
	}
}
