// models/wishlist.go
package models

import "time"

type WishlistStatus string

const (
	StatusPlanned WishlistStatus = "Planned"
	StatusVisited WishlistStatus = "Visited"
)

type WishlistItem struct {
	ID          string         `json:"id"`
	CountryName string         `json:"CountryName"`
	Note        string         `json:"Note"`
	Status      WishlistStatus `json:"Status"`
	CreatedAt   time.Time      `json:"created_at"`
}
