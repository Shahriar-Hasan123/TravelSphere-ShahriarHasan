package services

import (
	"TravelSphere/models"
	"sync"
	"testing"
)

// newTestWishlistService creates a fresh isolated instance for each test.
func newTestWishlistService() *WishlistService {
	return &WishlistService{
		store: make(map[string][]models.WishlistItem),
	}
}

func TestWishlistCreate(t *testing.T) {
	svc := newTestWishlistService()
	item := svc.Create("beta", "France", "Visit Eiffel Tower", "Planned")

	if item.ID == "" {
		t.Error("expected non-empty ID")
	}
	if item.CountryName != "France" {
		t.Errorf("expected 'France', got %q", item.CountryName)
	}
	if item.Note != "Visit Eiffel Tower" {
		t.Errorf("expected note, got %q", item.Note)
	}
	if item.Status != models.StatusPlanned {
		t.Errorf("expected Planned, got %q", item.Status)
	}
	if item.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestWishlistGetAll_Empty(t *testing.T) {
	svc := newTestWishlistService()
	items := svc.GetAll("beta")

	if len(items) != 0 {
		t.Errorf("expected empty slice, got %d items", len(items))
	}
}

func TestWishlistGetAll_ReturnsCopy(t *testing.T) {
	svc := newTestWishlistService()
	svc.Create("beta", "France", "", "Planned")

	items := svc.GetAll("beta")
	items[0].CountryName = "MODIFIED"

	original := svc.GetAll("beta")
	if original[0].CountryName == "MODIFIED" {
		t.Error("GetAll should return a copy, not a reference to internal state")
	}
}

func TestWishlistUpdate_Success(t *testing.T) {
	svc := newTestWishlistService()
	item := svc.Create("beta", "Japan", "", "Planned")

	updated, ok := svc.Update("beta", item.ID, "Cherry blossom tour", "Visited")
	if !ok {
		t.Fatal("expected update to succeed")
	}
	if updated.Note != "Cherry blossom tour" {
		t.Errorf("expected updated note, got %q", updated.Note)
	}
	if updated.Status != models.StatusVisited {
		t.Errorf("expected Visited, got %q", updated.Status)
	}
}

func TestWishlistUpdate_NotFound(t *testing.T) {
	svc := newTestWishlistService()

	_, ok := svc.Update("beta", "non-existent-id", "note", "Planned")
	if ok {
		t.Error("expected update to fail for unknown ID")
	}
}

func TestWishlistDelete_Success(t *testing.T) {
	svc := newTestWishlistService()
	item := svc.Create("beta", "Germany", "", "Planned")

	ok := svc.Delete("beta", item.ID)
	if !ok {
		t.Error("expected delete to succeed")
	}

	items := svc.GetAll("beta")
	if len(items) != 0 {
		t.Errorf("expected 0 items after delete, got %d", len(items))
	}
}

func TestWishlistDelete_NotFound(t *testing.T) {
	svc := newTestWishlistService()

	ok := svc.Delete("beta", "non-existent-id")
	if ok {
		t.Error("expected delete to return false for unknown ID")
	}
}

func TestWishlistExists_True(t *testing.T) {
	svc := newTestWishlistService()
	svc.Create("beta", "France", "", "Planned")

	if !svc.Exists("beta", "France") {
		t.Error("expected Exists to return true")
	}
}

func TestWishlistExists_False(t *testing.T) {
	svc := newTestWishlistService()

	if svc.Exists("beta", "France") {
		t.Error("expected Exists to return false for empty store")
	}
}

func TestWishlistExists_DifferentUser(t *testing.T) {
	svc := newTestWishlistService()
	svc.Create("alice", "France", "", "Planned")

	if svc.Exists("beta", "France") {
		t.Error("expected Exists false for different user")
	}
}

func TestWishlistSummary(t *testing.T) {
	svc := newTestWishlistService()
	svc.Create("beta", "France", "", "Planned")
	svc.Create("beta", "Japan", "", "Planned")
	svc.Create("beta", "Bangladesh", "", "Visited")

	total, planned, visited := svc.Summary("beta")

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

func TestWishlistSummary_Empty(t *testing.T) {
	svc := newTestWishlistService()

	total, planned, visited := svc.Summary("beta")
	if total != 0 || planned != 0 || visited != 0 {
		t.Errorf("expected all zeros, got %d/%d/%d", total, planned, visited)
	}
}

func TestWishlistPerUserIsolation(t *testing.T) {
	svc := newTestWishlistService()
	svc.Create("alice", "France", "", "Planned")
	svc.Create("beta", "Japan", "", "Planned")

	aliceItems := svc.GetAll("alice")
	betaItems := svc.GetAll("beta")

	if len(aliceItems) != 1 || aliceItems[0].CountryName != "France" {
		t.Error("alice should only see France")
	}
	if len(betaItems) != 1 || betaItems[0].CountryName != "Japan" {
		t.Error("beta should only see Japan")
	}
}

func TestWishlistConcurrentAccess(t *testing.T) {
	svc := newTestWishlistService()
	var wg sync.WaitGroup

	// Concurrent writes from multiple goroutines.
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			svc.Create("beta", string(rune('A'+n)), "", "Planned")
		}(i)
	}
	wg.Wait()

	items := svc.GetAll("beta")
	if len(items) != 20 {
		t.Errorf("expected 20 items after concurrent writes, got %d", len(items))
	}
}

func TestGetWishlistService_Singleton(t *testing.T) {
	s1 := GetWishlistService()
	s2 := GetWishlistService()

	if s1 != s2 {
		t.Error("GetWishlistService must return the same instance")
	}
}
