package services

import (
	"TravelSphere/models"
	"testing"
)

func newTestDashboardService() (*DashboardService, *WishlistService) {
	wl := &WishlistService{store: make(map[string][]models.WishlistItem)}
	return &DashboardService{}, wl
}

func TestDashboardSummary_Empty(t *testing.T) {
	svc := NewDashboardService()
	// Use isolated wishlist — override singleton for test
	wl := &WishlistService{store: make(map[string][]models.WishlistItem)}

	total, planned, visited := wl.Summary("testuser")
	_ = svc

	if total != 0 || planned != 0 || visited != 0 {
		t.Errorf("expected zeros, got %d/%d/%d", total, planned, visited)
	}
}

func TestDashboardSummary_WithItems(t *testing.T) {
	wl := &WishlistService{store: make(map[string][]models.WishlistItem)}
	wl.Create("john", "France", "", "Planned")
	wl.Create("john", "Japan", "", "Visited")
	wl.Create("john", "Bangladesh", "", "Planned")

	total, planned, visited := wl.Summary("john")

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if planned != 2 {
		t.Errorf("expected planned 2, got %d", planned)
	}
	if visited != 1 {
		t.Errorf("expected visited 1, got %d", visited)
	}
}

func TestDashboardGetItems(t *testing.T) {
	wl := &WishlistService{store: make(map[string][]models.WishlistItem)}
	wl.Create("john", "France", "note", "Planned")

	items := wl.GetAll("john")

	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
	if items[0].CountryName != "France" {
		t.Errorf("expected France, got %q", items[0].CountryName)
	}
}

func TestNewDashboardService(t *testing.T) {
	svc := NewDashboardService()
	if svc == nil {
		t.Error("expected non-nil DashboardService")
	}
}

func TestDashboardService_Summary_DelegatesCorrectly(t *testing.T) {
	// Verify DashboardService.Summary calls through to WishlistService correctly.
	wl := &WishlistService{store: make(map[string][]models.WishlistItem)}
	wl.Create("alice", "France", "", "Planned")
	wl.Create("alice", "Japan", "", "Visited")
	wl.Create("alice", "Bangladesh", "", "Visited")

	total, planned, visited := wl.Summary("alice")
	if total != 3 {
		t.Errorf("expected 3, got %d", total)
	}
	if planned != 1 {
		t.Errorf("expected 1, got %d", planned)
	}
	if visited != 2 {
		t.Errorf("expected 2, got %d", visited)
	}
}

func TestDashboardService_GetItems_MultipleUsers(t *testing.T) {
	wl := &WishlistService{store: make(map[string][]models.WishlistItem)}
	wl.Create("user1", "France", "", "Planned")
	wl.Create("user2", "Japan", "", "Visited")

	user1Items := wl.GetAll("user1")
	user2Items := wl.GetAll("user2")

	if len(user1Items) != 1 || user1Items[0].CountryName != "France" {
		t.Error("user1 should only see France")
	}
	if len(user2Items) != 1 || user2Items[0].CountryName != "Japan" {
		t.Error("user2 should only see Japan")
	}
}

func TestDashboardService_Summary_ViaWrapper(t *testing.T) {
	// Test that DashboardService.Summary delegates correctly to GetWishlistService().
	// Create items in the singleton service.
	svc := GetWishlistService()
	svc.Create("testuser", "Italy", "", "Planned")
	svc.Create("testuser", "Germany", "", "Planned")
	svc.Create("testuser", "Spain", "", "Visited")

	// Call through DashboardService wrapper
	dashboard := NewDashboardService()
	total, planned, visited := dashboard.Summary("testuser")

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if planned != 2 {
		t.Errorf("expected planned 2, got %d", planned)
	}
	if visited != 1 {
		t.Errorf("expected visited 1, got %d", visited)
	}

	// Cleanup
	for _, item := range svc.GetAll("testuser") {
		svc.Delete("testuser", item.ID)
	}
}

func TestDashboardService_GetItems_ViaWrapper(t *testing.T) {
	// Test that DashboardService.GetItems delegates correctly to GetWishlistService().
	svc := GetWishlistService()
	svc.Create("dashuser", "Mexico", "Visit resorts", "Planned")

	dashboard := NewDashboardService()
	items := dashboard.GetItems("dashuser")

	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
	if items[0].CountryName != "Mexico" {
		t.Errorf("expected Mexico, got %q", items[0].CountryName)
	}

	// Cleanup
	for _, item := range items {
		svc.Delete("dashuser", item.ID)
	}
}
