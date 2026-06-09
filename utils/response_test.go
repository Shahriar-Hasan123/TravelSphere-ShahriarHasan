package utils

import "testing"

func TestOKResponse(t *testing.T) {
	data := map[string]string{"key": "value"}
	resp := OKResponse(data)

	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", resp.Status)
	}
	if resp.Data == nil {
		t.Error("expected data, got nil")
	}
}

func TestCreatedResponse(t *testing.T) {
	data := map[string]string{"id": "123"}
	resp := CreatedResponse(data)

	if resp.Status != "created" {
		t.Errorf("expected status 'created', got %q", resp.Status)
	}
	if resp.Data == nil {
		t.Error("expected data, got nil")
	}
}

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name    string
		message string
		code    int
	}{
		{"not found", "Country not found", 404},
		{"bad request", "invalid JSON body", 400},
		{"server error", "internal error", 500},
		{"unauthorized", "Authentication required", 401},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ErrorResponse(tt.message, tt.code)
			if resp.Status != "error" {
				t.Errorf("expected status 'error', got %q", resp.Status)
			}
			if resp.Message != tt.message {
				t.Errorf("expected message %q, got %q", tt.message, resp.Message)
			}
			if resp.Code != tt.code {
				t.Errorf("expected code %d, got %d", tt.code, resp.Code)
			}
		})
	}
}

func TestOKResponseNilData(t *testing.T) {
	resp := OKResponse(nil)
	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", resp.Status)
	}
}