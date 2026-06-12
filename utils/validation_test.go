package utils

import "testing"

func TestValidateWishlistCreate(t *testing.T) {
	tests := []struct {
		name        string
		countryName string
		wantErr     bool
	}{
		{"valid name", "France", false},
		{"empty string", "", true},
		{"whitespace only", "   ", true},
		{"single char", "A", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := ValidateWishlistCreate(tt.countryName)
			if tt.wantErr && msg == "" {
				t.Error("expected error message, got empty string")
			}
			if !tt.wantErr && msg != "" {
				t.Errorf("expected no error, got %q", msg)
			}
		})
	}
}

func TestValidateWishlistStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		wantErr bool
	}{
		{"planned", "Planned", false},
		{"visited", "Visited", false},
		{"lowercase planned", "planned", true},
		{"empty", "", true},
		{"invalid", "Maybe", true},
		{"random", "Done", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := ValidateWishlistStatus(tt.status)
			if tt.wantErr && msg == "" {
				t.Error("expected error message, got empty string")
			}
			if !tt.wantErr && msg != "" {
				t.Errorf("expected no error, got %q", msg)
			}
		})
	}
}
