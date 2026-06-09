package utils

import "testing"

func TestFormatPopulation(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"billions", 1_500_000_000, "1.5B"},
		{"millions", 2_400_000, "2.4M"},
		{"thousands", 47_400, "47.4K"},
		{"under thousand", 500, "500"},
		{"exact million", 1_000_000, "1.0M"},
		{"exact billion", 1_000_000_000, "1.0B"},
		{"exact thousand", 1_000, "1.0K"},
		{"zero", 0, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatPopulation(tt.input)
			if got != tt.expected {
				t.Errorf("FormatPopulation(%d) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestToSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "Albania", "albania"},
		{"with space", "United States", "united-states"},
		{"with apostrophe", "Côte d'Ivoire", "cte-divoire"},
		{"multiple spaces", "United  Arab  Emirates", "united-arab-emirates"},
		{"already lowercase", "france", "france"},
		{"special chars", "São Tomé & Príncipe", "so-tom-prncipe"},
		{"empty string", "", ""},
		{"hyphens", "Guinea-Bissau", "guinea-bissau"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToSlug(tt.input)
			if got != tt.expected {
				t.Errorf("ToSlug(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatCurrency(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		curr     string
		expected string
	}{
		{"full", "ALL", "Albanian lek", "ALL (Albanian lek)"},
		{"code only", "USD", "", "USD"},
		{"empty code", "", "Dollar", ""},
		{"both empty", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatCurrency(tt.code, tt.curr)
			if got != tt.expected {
				t.Errorf("FormatCurrency(%q,%q) = %q, want %q", tt.code, tt.curr, got, tt.expected)
			}
		})
	}
}