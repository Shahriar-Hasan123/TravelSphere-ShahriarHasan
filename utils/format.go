// utils/format.go
// Formatting helpers for population, currency, languages, and slugs.
package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// FormatPopulation converts a raw int64 population into a human-readable
// string. Examples: 2400000 → "2.4M", 47400000 → "47.4M", 1300 → "1.3K".
func FormatPopulation(pop int64) string {
	switch {
	case pop >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", float64(pop)/1_000_000_000)
	case pop >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(pop)/1_000_000)
	case pop >= 1_000:
		return fmt.Sprintf("%.1fK", float64(pop)/1_000)
	default:
		return fmt.Sprintf("%d", pop)
	}
}

// ToSlug converts a country name into a URL-safe lowercase slug.
// Examples: "United States" → "united-states", "Côte d'Ivoire" → "cote-divoire".
func ToSlug(name string) string {
	// Lowercase everything.
	s := strings.ToLower(name)

	// Replace spaces and underscores with hyphens.
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// Remove all characters that are not alphanumeric or hyphens.
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	s = reg.ReplaceAllString(s, "")

	// Collapse multiple consecutive hyphens into one.
	reg2 := regexp.MustCompile(`-+`)
	s = reg2.ReplaceAllString(s, "-")

	return strings.Trim(s, "-")
}

// FormatCurrency builds a display string from a currency code and name.
// Example: code="ALL", name="Albanian lek" → "ALL (Albanian lek)".
func FormatCurrency(code, name string) string {
	if code == "" {
		return ""
	}
	if name == "" {
		return code
	}
	return fmt.Sprintf("%s (%s)", code, name)
}
