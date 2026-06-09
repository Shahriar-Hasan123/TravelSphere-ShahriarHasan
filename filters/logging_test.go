package filters

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/beego/beego/v2/server/web/context"
)

func newTestContext(method, path string) *context.Context {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	ctx := context.NewContext()
	ctx.Reset(w, req)
	return ctx
}

func TestLoggingFilter_DoesNotPanic(t *testing.T) {
	ctx := newTestContext(http.MethodGet, "/countries")
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("LoggingFilter panicked: %v", r)
		}
	}()
	LoggingFilter(ctx)
}

func TestLoggingFilter_AllMethods(t *testing.T) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
	}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			ctx := newTestContext(method, "/api/wishlist")
			LoggingFilter(ctx)
		})
	}
}

func TestLoggingFilter_ZeroStatus(t *testing.T) {
	// Covers the branch where ResponseWriter.Status == 0 (treated as 200).
	ctx := newTestContext(http.MethodGet, "/dashboard")
	LoggingFilter(ctx)
}

func TestFormatDuration_Microseconds(t *testing.T) {
	result := formatDuration(500 * time.Nanosecond)
	if result == "" {
		t.Error("expected non-empty string")
	}
}

func TestFormatDuration_Milliseconds(t *testing.T) {
	result := formatDuration(2 * time.Millisecond)
	if result == "" {
		t.Error("expected non-empty string")
	}
}

func TestFormatDuration_Seconds(t *testing.T) {
	result := formatDuration(3 * time.Second)
	if result == "" {
		t.Error("expected non-empty string")
	}
}

func TestFormatDuration_ExactMillisecond(t *testing.T) {
	result := formatDuration(time.Millisecond)
	if result == "" {
		t.Error("expected non-empty string")
	}
}

func TestFormatDuration_ExactSecond(t *testing.T) {
	result := formatDuration(time.Second)
	if result == "" {
		t.Error("expected non-empty string")
	}
}

func TestNanosecondDuration(t *testing.T) {
	// Test the nanosecondDuration helper function
	result := nanosecondDuration(1000000000) // 1 second in nanoseconds
	if result != time.Second {
		t.Errorf("expected 1 second, got %v", result)
	}

	result = nanosecondDuration(1000000) // 1 millisecond in nanoseconds
	if result != time.Millisecond {
		t.Errorf("expected 1 millisecond, got %v", result)
	}

	result = nanosecondDuration(1000) // 1 microsecond in nanoseconds
	if result != time.Microsecond {
		t.Errorf("expected 1 microsecond, got %v", result)
	}
}
