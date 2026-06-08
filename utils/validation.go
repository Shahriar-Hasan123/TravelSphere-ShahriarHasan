// Request payload validators for wishlist operations.

package utils

import (
	"TravelSphere/models"
	"strings"
)

// ValidateWishlistCreate checks required fields for a create request.
func ValidateWishlistCreate(countryName string) string {
	if strings.TrimSpace(countryName) == "" {
		return "country_name is required"
	}
	return ""
}

// ValidateWishlistStatus ensures the status value is an allowed enum.
func ValidateWishlistStatus(status string) string {
	s := models.WishlistStatus(status)
	if s != models.StatusPlanned && s != models.StatusVisited {
		return "status must be 'Planned' or 'Visited'"
	}
	return ""
}