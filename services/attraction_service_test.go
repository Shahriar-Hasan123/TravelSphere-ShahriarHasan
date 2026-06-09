package services

import (
	"TravelSphere/utils/clients"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockAttractionServer(places []clients.RawPlace) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Features []struct {
				Properties clients.RawPlace `json:"properties"`
			} `json:"features"`
		}{}
		for _, p := range places {
			resp.Features = append(resp.Features, struct {
				Properties clients.RawPlace `json:"properties"`
			}{Properties: p})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestGetAttractionsByCoords_Success(t *testing.T) {
	places := []clients.RawPlace{
		{Name: "Eiffel Tower", Kinds: "architecture,historic"},
		{Name: "Louvre",       Kinds: "museums"},
		{Name: "X",            Kinds: "misc"}, // too short — should be filtered
	}

	server := mockAttractionServer(places)
	defer server.Close()

	svc := &AttractionService{
		client: clients.NewOpenTripMapClientWithURL(server.URL, "test-key"),
	}

	attractions, err := svc.GetAttractionsByCoords(48.85, 2.29)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "X" should be filtered out — name too short
	if len(attractions) != 2 {
		t.Errorf("expected 2 attractions, got %d", len(attractions))
	}
	if attractions[0].Name != "Eiffel Tower" {
		t.Errorf("expected 'Eiffel Tower', got %q", attractions[0].Name)
	}
}

func TestGetAttractionsByCoords_EmptyResponse(t *testing.T) {
	server := mockAttractionServer([]clients.RawPlace{})
	defer server.Close()

	svc := &AttractionService{
		client: clients.NewOpenTripMapClientWithURL(server.URL, "test-key"),
	}

	attractions, err := svc.GetAttractionsByCoords(48.85, 2.29)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(attractions) != 0 {
		t.Errorf("expected empty slice, got %d", len(attractions))
	}
}

func TestParseKinds(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"architecture,historic", []string{"architecture", "historic"}},
		{"interesting_places",    []string{"interesting places"}},
		{"",                      []string{}},
		{"museums,,historic",     []string{"museums", "historic"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseKinds(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("parseKinds(%q): expected %v, got %v", tt.input, tt.expected, got)
				return
			}
			for i := range tt.expected {
				if got[i] != tt.expected[i] {
					t.Errorf("parseKinds(%q)[%d]: expected %q, got %q", tt.input, i, tt.expected[i], got[i])
				}
			}
		})
	}
}

func TestGetAttractionsByCoords_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer server.Close()

	svc := &AttractionService{
		client: clients.NewOpenTripMapClientWithURL(server.URL, "test-key"),
	}

	_, err := svc.GetAttractionsByCoords(48.85, 2.29)
	if err == nil {
		t.Error("expected error from API failure")
	}
}

func TestNewAttractionService_NotNil(t *testing.T) {
	svc := NewAttractionService()
	if svc == nil {
		t.Error("expected non-nil AttractionService")
	}
}