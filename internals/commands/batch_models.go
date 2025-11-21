package commands

// ClimParams represents climate control parameters
// Empty string values mean "keep current value"
type ClimParams struct {
	Power   string `json:"power,omitempty"`   // "0" or "1"
	Mode    string `json:"mode,omitempty"`    // "0"=AUTO, "1"=HEAT, "2"=DRY, "3"=FAN, "4"=COOL
	Temp    string `json:"temp,omitempty"`    // Temperature (e.g., "22.0")
	FanRate string `json:"fan-rate,omitempty"` // Fan rate (e.g., "A", "3"-"7")
	FanDir  string `json:"fan-dir,omitempty"`  // Fan direction: "0"=all wings stopped, "1"=vertical, "2"=horizontal, "3"=both
}

// DeviceOverride represents per-device parameter overrides
// Device is matched by name or IP (at least one must be specified)
type DeviceOverride struct {
	Name   string     `json:"name,omitempty"`   // Device name to match
	IP     string     `json:"ip,omitempty"`     // Device IP to match
	Params ClimParams `json:"params"`           // Parameters to apply to this device
}

// GroupConfig represents a group configuration
// All devices with matching grp_name will receive the default params,
// unless overridden by a DeviceOverride
type GroupConfig struct {
	GroupName string           `json:"group_name"`           // grp_name from basic_info to match
	Params    ClimParams       `json:"params"`               // Default parameters for all devices in group
	Overrides []DeviceOverride `json:"overrides,omitempty"`  // Per-device overrides
}

// BatchScript represents the root JSON structure for batch operations
type BatchScript struct {
	Groups []GroupConfig `json:"groups"` // List of group configurations to apply
}

