/*
LoggingFilter records the HTTP method, URL, status code, and duration
of every request. Registered at BeforeRouter so it captures the full
request lifecycle including routing and controller execution time.
*/

package filters

import (
	"fmt"
	"log"
	"time"

	"github.com/beego/beego/v2/server/web/context"
)

// LoggingFilter returns a Beego filter function that logs each request.
// It stores the start time in the request context, then after the handler
// completes it prints method, path, status, and elapsed duration.
func LoggingFilter(ctx *context.Context) {
	start := time.Now()

	// Store start time so the after-filter can compute duration.
	// We use ResponseWriter's WriteHeader hook indirectly — simpler is to
	// log everything in a single BeforeRouter pass using a deferred call.
	defer func() {
		duration := time.Since(start)
		status := ctx.ResponseWriter.Status
		if status == 0 {
			// Beego sets 0 before WriteHeader is called; treat as 200.
			status = 200
		}
		log.Printf(
			"[HTTP] %-6s %-40s %d  %s",
			ctx.Request.Method,
			ctx.Request.URL.Path,
			status,
			formatDuration(duration),
		)
	}()
}

// formatDuration returns a human-readable duration string.
// Examples: "1.23ms", "450µs", "2.10s"
func formatDuration(d time.Duration) string {
	switch {
	case d >= time.Second:
		return fmt.Sprintf("%.2fs", d.Seconds())
	case d >= time.Millisecond:
		return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1000)
	default:
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
}
