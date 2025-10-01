package storage

import (
	"time"
)

// DeviceSnapshot represents a snapshot of device information at a specific point in time
type DeviceSnapshot struct {
	IP           string            `json:"ip"`
	MAC          string            `json:"mac"`
	Name         string            `json:"name"`
	Model        string            `json:"model"`
	Status       string            `json:"status"`
	BasicInfo    map[string]string `json:"basic_info"`
	ControlInfo  map[string]string `json:"control_info"`
	DiscoveredAt time.Time         `json:"discovered_at"`
	LastSeenAt   time.Time         `json:"last_seen_at"`
}

// DeviceHistory represents the complete history of a device including all snapshots and changes
type DeviceHistory struct {
	MAC       string           `json:"mac"`       // Primary key
	Device    DeviceSnapshot   `json:"device"`    // Latest snapshot
	Snapshots []DeviceSnapshot `json:"snapshots"` // All historical snapshots
	Changes   []DeviceChange   `json:"changes"`   // Detected changes over time
}

// DeviceChange represents a detected change in device information
type DeviceChange struct {
	Field     string    `json:"field"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	ChangedAt time.Time `json:"changed_at"`
}

// DeviceStorage represents the complete storage structure
type DeviceStorage struct {
	Devices     map[string]*DeviceHistory `json:"devices"` // Keyed by MAC address
	LastUpdated time.Time                 `json:"last_updated"`
}
