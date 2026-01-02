// Package models provides data structures for the Sabakan application.
package models

import "time"

// ContainerState represents the current state of a container.
type ContainerState string

const (
	// StateRunning indicates a running container.
	StateRunning ContainerState = "running"
	// StateStopped indicates a stopped container.
	StateStopped ContainerState = "stopped"
	// StateCreated indicates a created but not started container.
	StateCreated ContainerState = "created"
	// StatePaused indicates a paused container.
	StatePaused ContainerState = "paused"
	// StateRestarting indicates a restarting container.
	StateRestarting ContainerState = "restarting"
	// StateExited indicates an exited container.
	StateExited ContainerState = "exited"
	// StateUnknown indicates an unknown state.
	StateUnknown ContainerState = "unknown"
)

// PortMapping represents a container port mapping.
type PortMapping struct {
	// HostIP is the host IP address.
	HostIP string `json:"hostIp,omitempty"`
	// HostPort is the port on the host.
	HostPort uint16 `json:"hostPort"`
	// ContainerPort is the port inside the container.
	ContainerPort uint16 `json:"containerPort"`
	// Protocol is the protocol (tcp, udp).
	Protocol string `json:"protocol"`
}

// Container represents a game server container.
type Container struct {
	// ID is the unique identifier of the container.
	ID string `json:"id"`
	// Name is the human-readable name of the container.
	Name string `json:"name"`
	// Image is the container image used.
	Image string `json:"image"`
	// State is the current state of the container.
	State ContainerState `json:"state"`
	// Status is a human-readable status string.
	Status string `json:"status"`
	// Created is the creation timestamp.
	Created time.Time `json:"created"`
	// Ports is a list of port mappings.
	Ports []PortMapping `json:"ports"`
	// Labels are container labels.
	Labels map[string]string `json:"labels"`
}

// ContainerLogEntry represents a single log entry.
type ContainerLogEntry struct {
	// Timestamp is when the log was generated.
	Timestamp time.Time `json:"timestamp,omitempty"`
	// Stream is the output stream (stdout or stderr).
	Stream string `json:"stream"`
	// Message is the log message content.
	Message string `json:"message"`
}
