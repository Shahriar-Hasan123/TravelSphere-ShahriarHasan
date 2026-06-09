// mockSession implements beego's SessionStore interface for unit tests.
package filters

import (
	"context"
	"net/http"
)

type mockSession struct {
	data map[interface{}]interface{}
}

func (m *mockSession) Set(ctx context.Context, key, val interface{}) error {
	m.data[key] = val
	return nil
}

func (m *mockSession) Get(ctx context.Context, key interface{}) interface{} {
	return m.data[key]
}

func (m *mockSession) Delete(ctx context.Context, key interface{}) error {
	delete(m.data, key)
	return nil
}

func (m *mockSession) SessionID(ctx context.Context) string { return "mock-session-id" }
func (m *mockSession) SessionRelease(ctx context.Context, w http.ResponseWriter) {
	// no-op
}
func (m *mockSession) SessionReleaseIfPresent(ctx context.Context, w http.ResponseWriter) {
	// no-op
}
func (m *mockSession) Flush(ctx context.Context) error {
	m.data = make(map[interface{}]interface{})
	return nil
}
