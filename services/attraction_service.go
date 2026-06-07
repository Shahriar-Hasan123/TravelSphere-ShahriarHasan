// AttractionService retrieves and transforms tourist attraction data
// from OpenTripMap for a given country's coordinates.

package services

import (
	"TravelSphere/models"
	"TravelSphere/utils/clients"
	"strings"
)

const (
	attractionRadiusM = 100_000 // 100 km radius around country centre
	attractionLimit   = 15      // Maximum attractions to display
	minNameLength     = 3       // Skip unnamed or near-empty place names
)

// AttractionService provides tourist attraction data for a destination.
type AttractionService struct {
	client *clients.OpenTripMapClient
}

// NewAttractionService creates an AttractionService with its API client.
func NewAttractionService() *AttractionService {
	return &AttractionService{
		client: clients.NewOpenTripMapClient(),
	}
}

// GetAttractionsByCoords fetches attractions near the given coordinates
// and returns clean Attraction DTOs ready for template rendering.
func (s *AttractionService) GetAttractionsByCoords(
	lat, lon float64,
) ([]models.Attraction, error) {
	raw, err := s.client.FetchAttractionsByCoords(
		lat, lon, attractionRadiusM, attractionLimit,
	)
	if err != nil {
		return nil, err
	}

	attractions := make([]models.Attraction, 0, len(raw))
	for _, place := range raw {
		// Skip places without a meaningful name.
		name := strings.TrimSpace(place.Name)
		if len(name) < minNameLength {
			continue
		}

		kinds := parseKinds(place.Kinds)
		attractions = append(attractions, models.Attraction{
			Name:  name,
			Kinds: kinds,
		})
	}
	return attractions, nil
}

// parseKinds splits the comma-separated kinds string from OpenTripMap
// into a clean slice, removing empty entries and internal underscores.
func parseKinds(raw string) []string {
	parts := strings.Split(raw, ",")
	kinds := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		// Replace underscores with spaces for readable display.
		p = strings.ReplaceAll(p, "_", " ")
		kinds = append(kinds, p)
	}
	return kinds
}