// Package container provides container management services using Podman REST API.
package container

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sweetfish329/sabakan/backend/internal/models"
)

// Service provides container management operations using Podman REST API.
type Service struct {
	socketPath string
	baseURL    string
	client     *http.Client
}

// NewService creates a new container service with the specified socket path.
// socketPath should be in the format "unix:///path/to/socket", "tcp://host:port", or "http://host:port".
func NewService(socketPath string) *Service {
	// Parse the socket path to determine the transport type
	u, err := url.Parse(socketPath)
	if err != nil {
		// Default to unix socket
		u = &url.URL{Scheme: "unix", Path: socketPath}
	}

	var transport *http.Transport
	var baseURL string

	switch u.Scheme {
	case "unix":
		// Unix socket transport - use dummy host for URL
		path := u.Path
		if path == "" {
			// Handle "unix:///path" format
			path = u.Host + u.Path
		}
		transport = &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", path)
			},
		}
		baseURL = "http://d" // Dummy host for unix socket
	case "http", "https", "tcp":
		// TCP/HTTP transport - use actual URL
		transport = &http.Transport{}
		if u.Scheme == "tcp" {
			baseURL = "http://" + u.Host
		} else {
			baseURL = strings.TrimSuffix(socketPath, "/")
		}
	default:
		// Assume http URL
		transport = &http.Transport{}
		baseURL = strings.TrimSuffix(socketPath, "/")
	}

	return &Service{
		socketPath: socketPath,
		baseURL:    baseURL,
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}
}

// apiURL constructs the full URL for an API endpoint.
func (s *Service) apiURL(path string) string {
	return s.baseURL + path
}

// List returns all containers.
func (s *Service) List(ctx context.Context) ([]models.Container, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.apiURL("/v5.0.0/libpod/containers/json?all=true"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var podmanContainers []podmanListContainer
	if err := json.NewDecoder(resp.Body).Decode(&podmanContainers); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	result := make([]models.Container, 0, len(podmanContainers))
	for _, c := range podmanContainers {
		result = append(result, c.toModel())
	}

	return result, nil
}

// Get returns a specific container by ID or name.
func (s *Service) Get(ctx context.Context, id string) (*models.Container, error) {
	endpoint := fmt.Sprintf("/v5.0.0/libpod/containers/%s/json", url.PathEscape(id))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.apiURL(endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get container %s: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("container %s not found", id)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var data podmanInspectData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	container := data.toModel()
	return &container, nil
}

// Start starts a container by ID or name.
func (s *Service) Start(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/v5.0.0/libpod/containers/%s/start", url.PathEscape(id))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.apiURL(endpoint), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to start container %s: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to start container %s: status %d: %s", id, resp.StatusCode, string(body))
	}

	return nil
}

// Stop stops a container by ID or name with optional timeout in seconds.
func (s *Service) Stop(ctx context.Context, id string, timeout uint) error {
	endpoint := fmt.Sprintf("/v5.0.0/libpod/containers/%s/stop", url.PathEscape(id))
	if timeout > 0 {
		endpoint = fmt.Sprintf("%s?timeout=%d", endpoint, timeout)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.apiURL(endpoint), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to stop container %s: %w", id, err)
	}
	defer resp.Body.Close()

	// 304 means already stopped, which is acceptable
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to stop container %s: status %d: %s", id, resp.StatusCode, string(body))
	}

	return nil
}

// Logs returns the last N lines of container logs.
func (s *Service) Logs(ctx context.Context, id string, lines int) ([]models.ContainerLogEntry, error) {
	endpoint := fmt.Sprintf("/v5.0.0/libpod/containers/%s/logs?stdout=true&stderr=true&tail=%d", url.PathEscape(id), lines)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.apiURL(endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for container %s: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	// Read log output - each line is prefixed with stream type
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read logs: %w", err)
	}

	entries := parseLogOutput(body)
	return entries, nil
}

// podmanListContainer represents a container in Podman list response.
type podmanListContainer struct {
	ID      string            `json:"Id"`
	Names   []string          `json:"Names"`
	Image   string            `json:"Image"`
	State   string            `json:"State"`
	Status  string            `json:"Status"`
	Created int64             `json:"Created"`
	Ports   []podmanPort      `json:"Ports"`
	Labels  map[string]string `json:"Labels"`
}

type podmanPort struct {
	HostIP        string `json:"host_ip"`
	HostPort      uint16 `json:"host_port"`
	ContainerPort uint16 `json:"container_port"`
	Protocol      string `json:"protocol"`
}

func (c *podmanListContainer) toModel() models.Container {
	name := ""
	if len(c.Names) > 0 {
		name = c.Names[0]
	}

	ports := make([]models.PortMapping, 0, len(c.Ports))
	for _, p := range c.Ports {
		ports = append(ports, models.PortMapping{
			HostIP:        p.HostIP,
			HostPort:      p.HostPort,
			ContainerPort: p.ContainerPort,
			Protocol:      p.Protocol,
		})
	}

	return models.Container{
		ID:      c.ID,
		Name:    name,
		Image:   c.Image,
		State:   mapState(c.State),
		Status:  c.Status,
		Created: time.Unix(c.Created, 0),
		Ports:   ports,
		Labels:  c.Labels,
	}
}

// podmanInspectData represents container inspect response.
type podmanInspectData struct {
	ID      string `json:"Id"`
	Name    string `json:"Name"`
	Created string `json:"Created"`
	State   struct {
		Status string `json:"Status"`
	} `json:"State"`
	Config struct {
		Image  string            `json:"Image"`
		Labels map[string]string `json:"Labels"`
	} `json:"Config"`
}

func (d *podmanInspectData) toModel() models.Container {
	created, _ := time.Parse(time.RFC3339Nano, d.Created)

	return models.Container{
		ID:      d.ID,
		Name:    d.Name,
		Image:   d.Config.Image,
		State:   mapState(d.State.Status),
		Status:  d.State.Status,
		Created: created,
		Ports:   []models.PortMapping{},
		Labels:  d.Config.Labels,
	}
}

// mapState converts a Podman state string to our ContainerState.
func mapState(state string) models.ContainerState {
	switch state {
	case "running":
		return models.StateRunning
	case "stopped":
		return models.StateStopped
	case "created":
		return models.StateCreated
	case "paused":
		return models.StatePaused
	case "restarting":
		return models.StateRestarting
	case "exited":
		return models.StateExited
	default:
		return models.StateUnknown
	}
}

// parseLogOutput parses the Docker/Podman log output format.
func parseLogOutput(data []byte) []models.ContainerLogEntry {
	var entries []models.ContainerLogEntry

	// Simple line-by-line parsing
	lines := splitLines(data)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		entries = append(entries, models.ContainerLogEntry{
			Stream:  "stdout",
			Message: string(line),
		})
	}

	return entries
}

// splitLines splits byte data into lines.
func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
