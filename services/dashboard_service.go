// DashboardService aggregates wishlist data for the dashboard view.
package services

import "TravelSphere/models"

type DashboardService struct{}

func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// Summary returns total, planned, and visited counts for a user.
func (s *DashboardService) Summary(username string) (total, planned, visited int) {
	return GetWishlistService().Summary(username)
}

// GetItems returns the full wishlist item list for the dashboard destination list.
func (s *DashboardService) GetItems(username string) []models.WishlistItem {
	return GetWishlistService().GetAll(username)
}
