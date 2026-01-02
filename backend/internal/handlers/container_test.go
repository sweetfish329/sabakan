package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/sweetfish329/sabakan/backend/internal/container"
	"github.com/sweetfish329/sabakan/backend/internal/models"
)

// mockServer creates a test server that mocks Podman API responses.
func mockServer(t *testing.T, handlers map[string]http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for path, handler := range handlers {
			if r.URL.Path == path {
				handler(w, r)
				return
			}
		}
		http.NotFound(w, r)
	}))
}

func TestContainerHandler_List(t *testing.T) {
	// Create mock Podman server
	mockPodman := mockServer(t, map[string]http.HandlerFunc{
		"/v5.0.0/libpod/containers/json": func(w http.ResponseWriter, _ *http.Request) {
			containers := []map[string]interface{}{
				{
					"Id":      "abc123",
					"Names":   []string{"test-container"},
					"Image":   "nginx:latest",
					"State":   "running",
					"Status":  "Up 10 minutes",
					"Created": 1704067200,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(containers)
		},
	})
	defer mockPodman.Close()

	// Create service pointing to mock server
	svc := container.NewService(mockPodman.URL)
	handler := NewContainerHandler(svc)

	// Create Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/containers", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call handler
	err := handler.List(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify response
	var containers []models.Container
	err = json.Unmarshal(rec.Body.Bytes(), &containers)
	assert.NoError(t, err)
	assert.Len(t, containers, 1)
	assert.Equal(t, "abc123", containers[0].ID)
	assert.Equal(t, "test-container", containers[0].Name)
	assert.Equal(t, models.StateRunning, containers[0].State)
}

func TestContainerHandler_Get(t *testing.T) {
	mockPodman := mockServer(t, map[string]http.HandlerFunc{
		"/v5.0.0/libpod/containers/test-container/json": func(w http.ResponseWriter, _ *http.Request) {
			data := map[string]interface{}{
				"Id":      "abc123",
				"Name":    "test-container",
				"Created": "2024-01-01T00:00:00Z",
				"State": map[string]string{
					"Status": "running",
				},
				"Config": map[string]interface{}{
					"Image":  "nginx:latest",
					"Labels": map[string]string{"app": "test"},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(data)
		},
	})
	defer mockPodman.Close()

	svc := container.NewService(mockPodman.URL)
	handler := NewContainerHandler(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/containers/test-container", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("test-container")

	err := handler.Get(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var container models.Container
	err = json.Unmarshal(rec.Body.Bytes(), &container)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", container.ID)
	assert.Equal(t, "test-container", container.Name)
}

func TestContainerHandler_Start(t *testing.T) {
	mockPodman := mockServer(t, map[string]http.HandlerFunc{
		"/v5.0.0/libpod/containers/test-container/start": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer mockPodman.Close()

	svc := container.NewService(mockPodman.URL)
	handler := NewContainerHandler(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/containers/test-container/start", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("test-container")

	err := handler.Start(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestContainerHandler_Stop(t *testing.T) {
	mockPodman := mockServer(t, map[string]http.HandlerFunc{
		"/v5.0.0/libpod/containers/test-container/stop": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer mockPodman.Close()

	svc := container.NewService(mockPodman.URL)
	handler := NewContainerHandler(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/containers/test-container/stop", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("test-container")

	err := handler.Stop(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestContainerHandler_Get_NotFound(t *testing.T) {
	mockPodman := mockServer(t, map[string]http.HandlerFunc{
		"/v5.0.0/libpod/containers/nonexistent/json": func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"cause":"no such container","message":"no such container","response":404}`))
		},
	})
	defer mockPodman.Close()

	svc := container.NewService(mockPodman.URL)
	handler := NewContainerHandler(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/containers/nonexistent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nonexistent")

	err := handler.Get(c)
	assert.Error(t, err)

	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, he.Code)
}
