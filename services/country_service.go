// CountryService orchestrates country data retrieval and transformation.
// All business logic (filtering, sorting, DTO mapping) lives here.

package services

import (
	"TravelSphere/models"
	"TravelSphere/utils"
	"TravelSphere/utils/clients"
	"sort"
	"strings"
)

// CountryService provides country-related business operations.
type CountryService struct {
	client *clients.RestCountriesClient
}

// NewCountryService creates a CountryService with its required API client.
func NewCountryService() *CountryService {
	return &CountryService{
		client: clients.NewRestCountriesClient(),
	}
}

// GetAllCountries fetches all countries, applies optional search and region
// filters, sorts them alphabetically, and returns clean Country DTOs.
func (s *CountryService) GetAllCountries(search, region string) ([]models.Country, error) {
	raw, err := s.client.FetchAll()
	if err != nil {
		return nil, err
	}

	countries := make([]models.Country, 0, len(raw))
	search = strings.ToLower(strings.TrimSpace(search))
	region = strings.TrimSpace(region)

	for _, r := range raw {
		country := s.toDTO(r)

		// Apply region filter.
		if region != "" && !strings.EqualFold(country.Region, region) {
			continue
		}

		// Apply search filter — matches country name or capital.
		if search != "" {
			nameMatch := strings.Contains(strings.ToLower(country.Name), search)
			capitalMatch := strings.Contains(strings.ToLower(country.Capital), search)
			if !nameMatch && !capitalMatch {
				continue
			}
		}

		countries = append(countries, country)
	}

	// Sort alphabetically by country name.
	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Name < countries[j].Name
	})

	return countries, nil
}

// GetCountryBySlug finds a single country whose slug matches the given value.
// Returns nil, nil when not found (caller handles the 404 case).
func (s *CountryService) GetCountryBySlug(slug string) (*models.Country, error) {
	raw, err := s.client.FetchAll()
	if err != nil {
		return nil, err
	}

	slug = strings.ToLower(strings.TrimSpace(slug))

	for _, r := range raw {
		country := s.toDTO(r)
		if country.Slug == slug {
			return &country, nil
		}
	}
	return nil, nil // Not found.
}

// GetFeaturedCountries returns a curated short list for the home page.
// The featured slugs are hardcoded to match the UI screenshot exactly.
func (s *CountryService) GetFeaturedCountries() ([]models.FeaturedCountry, error) {
	featuredSlugs := []string{
		"united-states", "france", "japan", "australia", "brazil", "bangladesh",
	}

	raw, err := s.client.FetchAll()
	if err != nil {
		return nil, err
	}

	// Build a slug → DTO map for quick lookup.
	bySlug := make(map[string]models.Country, len(raw))
	for _, r := range raw {
		c := s.toDTO(r)
		bySlug[c.Slug] = c
	}

	featured := make([]models.FeaturedCountry, 0, len(featuredSlugs))
	for _, slug := range featuredSlugs {
		if c, ok := bySlug[slug]; ok {
			featured = append(featured, models.FeaturedCountry{
				Name:    c.Name,
				Slug:    c.Slug,
				Capital: c.Capital,
				Region:  c.Region,
				Flag:    c.Flag,
			})
		}
	}
	return featured, nil
}

// SearchSuggestions returns lightweight name+capital pairs for the home
// page search autocomplete dropdown.
func (s *CountryService) SearchSuggestions(query string) ([]models.FeaturedCountry, error) {
	if strings.TrimSpace(query) == "" {
		return nil, nil
	}
	countries, err := s.GetAllCountries(query, "")
	if err != nil {
		return nil, err
	}

	// Cap suggestions at 10 results.
	if len(countries) > 10 {
		countries = countries[:10]
	}

	suggestions := make([]models.FeaturedCountry, 0, len(countries))
	for _, c := range countries {
		suggestions = append(suggestions, models.FeaturedCountry{
			Name:    c.Name,
			Slug:    c.Slug,
			Capital: c.Capital,
			Region:  c.Region,
			Flag:    c.Flag,
		})
	}
	return suggestions, nil
}

// toDTO transforms a raw RawCountry from the API client into a clean
// models.Country DTO that the rest of the application uses.
func (s *CountryService) toDTO(r clients.RawCountry) models.Country {
	// Extract the first capital city if available.
	capital := ""
	if len(r.Capital) > 0 {
		capital = r.Capital[0]
	}

	// Build a single currency display string from the currencies map.
	currency := ""
	for code, cur := range r.Currencies {
		currency = utils.FormatCurrency(code, cur.Name)
		break // REST Countries may have multiple; we display the first.
	}

	// Collect all language names into a slice.
	languages := make([]string, 0, len(r.Languages))
	for _, lang := range r.Languages {
		languages = append(languages, lang)
	}
	sort.Strings(languages)

	// Use PNG flag
	flag := r.Flags.PNG
	if flag == "" {
		flag = r.Flags.SVG
	}

	name := r.Name.Common

	return models.Country{
		Name:         name,
		OfficialName: r.Name.Official,
		Slug:         utils.ToSlug(name),
		Flag:         flag,
		Capital:      capital,
		Population:   r.Population,
		FormattedPop: utils.FormatPopulation(r.Population),
		Region:       r.Region,
		Subregion:    r.Subregion,
		Currency:     currency,
		Languages:    languages,
		Latlng:       r.Latlng,
	}
}
