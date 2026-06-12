// WishlistService manages per-user wishlist data in an in-memory store.
// Thread-safe via sync.RWMutex. Single instance via sync.Once.

package services

import (
	"TravelSphere/models"
	"fmt"
	"sync"
	"time"
)

// WishlistService is the single in-memory wishlist store for all users.
type WishlistService struct {
	mu    sync.RWMutex
	store map[string][]models.WishlistItem // keyed by username
}

var (
	wishlistOnce     sync.Once
	wishlistInstance *WishlistService
)

// GetWishlistService returns the singleton WishlistService instance. Safe to call concurrently — initialization runs exactly once.
func GetWishlistService() *WishlistService {
	wishlistOnce.Do(func() {
		wishlistInstance = &WishlistService{
			store: make(map[string][]models.WishlistItem),
		}
	})
	return wishlistInstance
}

// GetAll returns all wishlist items for the given user.
func (s *WishlistService) GetAll(username string) []models.WishlistItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := s.store[username]
	if items == nil {
		return []models.WishlistItem{}
	}
	result := make([]models.WishlistItem, len(items))
	copy(result, items)
	return result
}

// Create adds a new wishlist entry for the user and returns it.
func (s *WishlistService) Create(username, countryName, note, status string) models.WishlistItem {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := models.WishlistItem{
		ID:          generateID(username, countryName),
		CountryName: countryName,
		Note:        note,
		Status:      models.WishlistStatus(status),
		CreatedAt:   time.Now(),
	}

	s.store[username] = append(s.store[username], item)
	return item
}

// Update modifies the note and status of an existing wishlist entry.
func (s *WishlistService) Update(username, id, note, status string) (models.WishlistItem, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := s.store[username]
	for i, item := range items {
		if item.ID == id {
			items[i].Note = note
			items[i].Status = models.WishlistStatus(status)
			s.store[username] = items
			return items[i], true
		}
	}
	return models.WishlistItem{}, false
}

// Delete removes a wishlist entry by ID.
func (s *WishlistService) Delete(username, id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := s.store[username]
	for i, item := range items {
		if item.ID == id {
			s.store[username] = append(items[:i], items[i+1:]...)
			return true
		}
	}
	return false
}

// Summary returns total, planned, and visited counts for a user.
func (s *WishlistService) Summary(username string) (total, planned, visited int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.store[username] {
		total++
		switch item.Status {
		case models.StatusPlanned:
			planned++
		case models.StatusVisited:
			visited++
		}
	}
	return
}

// generateID creates a unique ID from username, country, and current nanosecond.
func generateID(username, country string) string {
	return fmt.Sprintf("%s-%s-%d", username, country, time.Now().UnixNano())
}

// Exists returns true if the user already has this country in their wishlist.
func (s *WishlistService) Exists(username, countryName string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.store[username] {
		if item.CountryName == countryName {
			return true
		}
	}
	return false
}
