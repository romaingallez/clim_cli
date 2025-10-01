package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/romaingallez/clim_cli/internals/search"
)

const (
	StorageFileName = "devices.json"
)

// GetStoragePath returns the path to the storage file
func GetStoragePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %v", err)
	}

	climDir := filepath.Join(configDir, "clim_cli")
	if err := os.MkdirAll(climDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create clim_cli config directory: %v", err)
	}

	return filepath.Join(climDir, StorageFileName), nil
}

// LoadDeviceStorage loads the device storage from file
func LoadDeviceStorage() (*DeviceStorage, error) {
	path, err := GetStoragePath()
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return empty storage if file doesn't exist
		return &DeviceStorage{
			Devices:     make(map[string]*DeviceHistory),
			LastUpdated: time.Now(),
		}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage file: %v", err)
	}

	var storage DeviceStorage
	if err := json.Unmarshal(data, &storage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal storage data: %v", err)
	}

	// Initialize devices map if nil
	if storage.Devices == nil {
		storage.Devices = make(map[string]*DeviceHistory)
	}

	return &storage, nil
}

// SaveDeviceStorage saves the device storage to file
func SaveDeviceStorage(storage *DeviceStorage) error {
	path, err := GetStoragePath()
	if err != nil {
		return err
	}

	storage.LastUpdated = time.Now()

	data, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal storage data: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write storage file: %v", err)
	}

	log.Printf("Saved device storage to %s", path)
	return nil
}

// SaveDevices saves a list of discovered devices to storage
func SaveDevices(devices []search.Device) error {
	storage, err := LoadDeviceStorage()
	if err != nil {
		return fmt.Errorf("failed to load device storage: %v", err)
	}

	now := time.Now()

	for _, device := range devices {
		snapshot := DeviceSnapshot{
			IP:           device.IP,
			MAC:          device.MAC,
			Name:         device.Name,
			Model:        device.Model,
			Status:       device.Status,
			BasicInfo:    device.BasicInfo,
			ControlInfo:  device.ControlInfo,
			DiscoveredAt: now,
			LastSeenAt:   now,
		}

		// Check if device already exists
		if history, exists := storage.Devices[device.MAC]; exists {
			// Update existing device
			oldSnapshot := history.Device

			// Add current snapshot to history
			history.Snapshots = append(history.Snapshots, oldSnapshot)
			history.Device = snapshot

			// Detect changes
			changes := detectChanges(oldSnapshot, snapshot)
			history.Changes = append(history.Changes, changes...)

			log.Printf("Updated existing device %s (%s)", device.Name, device.MAC)
		} else {
			// Add new device
			history := &DeviceHistory{
				MAC:       device.MAC,
				Device:    snapshot,
				Snapshots: []DeviceSnapshot{snapshot},
				Changes:   []DeviceChange{},
			}
			storage.Devices[device.MAC] = history
			log.Printf("Added new device %s (%s)", device.Name, device.MAC)
		}
	}

	return SaveDeviceStorage(storage)
}

// GetDeviceHistories returns all device histories sorted by device name
func GetDeviceHistories() ([]*DeviceHistory, error) {
	storage, err := LoadDeviceStorage()
	if err != nil {
		return nil, err
	}

	histories := make([]*DeviceHistory, 0, len(storage.Devices))
	for _, history := range storage.Devices {
		histories = append(histories, history)
	}

	// Sort by device name
	sort.Slice(histories, func(i, j int) bool {
		return strings.ToLower(histories[i].Device.Name) < strings.ToLower(histories[j].Device.Name)
	})

	return histories, nil
}

// GetDeviceHistory returns the history for a specific device by MAC address
func GetDeviceHistory(mac string) (*DeviceHistory, error) {
	storage, err := LoadDeviceStorage()
	if err != nil {
		return nil, err
	}

	history, exists := storage.Devices[mac]
	if !exists {
		return nil, fmt.Errorf("device with MAC %s not found", mac)
	}

	return history, nil
}

// detectChanges compares two device snapshots and returns a list of changes
func detectChanges(old, new DeviceSnapshot) []DeviceChange {
	var changes []DeviceChange
	now := time.Now()

	// Compare basic fields
	if old.IP != new.IP {
		changes = append(changes, DeviceChange{
			Field:     "ip",
			OldValue:  old.IP,
			NewValue:  new.IP,
			ChangedAt: now,
		})
	}

	if old.Name != new.Name {
		changes = append(changes, DeviceChange{
			Field:     "name",
			OldValue:  old.Name,
			NewValue:  new.Name,
			ChangedAt: now,
		})
	}

	if old.Model != new.Model {
		changes = append(changes, DeviceChange{
			Field:     "model",
			OldValue:  old.Model,
			NewValue:  new.Model,
			ChangedAt: now,
		})
	}

	if old.Status != new.Status {
		changes = append(changes, DeviceChange{
			Field:     "status",
			OldValue:  old.Status,
			NewValue:  new.Status,
			ChangedAt: now,
		})
	}

	// Compare maps
	changes = append(changes, compareMaps("basic_info", old.BasicInfo, new.BasicInfo, now)...)
	changes = append(changes, compareMaps("control_info", old.ControlInfo, new.ControlInfo, now)...)

	return changes
}

// compareMaps compares two string maps and returns changes
func compareMaps(prefix string, old, new map[string]string, timestamp time.Time) []DeviceChange {
	var changes []DeviceChange

	// Check for changed/added values
	for key, newValue := range new {
		field := fmt.Sprintf("%s.%s", prefix, key)
		if oldValue, exists := old[key]; !exists {
			// New field added
			changes = append(changes, DeviceChange{
				Field:     field,
				OldValue:  "",
				NewValue:  newValue,
				ChangedAt: timestamp,
			})
		} else if oldValue != newValue {
			// Field value changed
			changes = append(changes, DeviceChange{
				Field:     field,
				OldValue:  oldValue,
				NewValue:  newValue,
				ChangedAt: timestamp,
			})
		}
	}

	// Check for removed values
	for key, oldValue := range old {
		field := fmt.Sprintf("%s.%s", prefix, key)
		if _, exists := new[key]; !exists {
			// Field removed
			changes = append(changes, DeviceChange{
				Field:     field,
				OldValue:  oldValue,
				NewValue:  "",
				ChangedAt: timestamp,
			})
		}
	}

	return changes
}

// GetRecentChanges returns devices that have changed within the specified duration
func GetRecentChanges(since time.Duration) ([]*DeviceHistory, error) {
	storage, err := LoadDeviceStorage()
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().Add(-since)
	var changed []*DeviceHistory

	for _, history := range storage.Devices {
		for _, change := range history.Changes {
			if change.ChangedAt.After(cutoff) {
				changed = append(changed, history)
				break
			}
		}
	}

	// Sort by device name
	sort.Slice(changed, func(i, j int) bool {
		return strings.ToLower(changed[i].Device.Name) < strings.ToLower(changed[j].Device.Name)
	})

	return changed, nil
}
