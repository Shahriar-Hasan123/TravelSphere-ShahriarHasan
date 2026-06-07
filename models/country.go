// Country DTOs - shaped from REST Countries API responses.

package models

// Country represents a single country for listing and detail views.
type Country struct {
	Name              string
	OfficialName      string
	Slug              string
	Flag              string
	Capital           string
	Population        int64
	FormattedPop      string   
	Region            string
	Subregion         string
	Currency          string
	Languages         []string
	Latlng            []float64 // [latitude, longitude] for OpenTripMap
}

type FeaturedCountry struct {
	Name    string
	Slug    string
	Capital string
	Region  string
	Flag    string
}