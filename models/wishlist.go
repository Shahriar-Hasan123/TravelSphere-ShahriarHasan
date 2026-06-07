package models

import "time"

type WishlistStatus string

const (
	StatusPlanned WishlistStatus = "Planned"
	StatusVisited WishlistStatus = "Visited"
)

// WishlistItem represents one saved destination in a user's travel wishlist.
type WishlistItem struct {
	ID          string         `json:"id"`
	CountryName string         `json:"country_name"`
	Note        string         `json:"note"`
	Status      WishlistStatus `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
}
