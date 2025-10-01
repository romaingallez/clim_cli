package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/romaingallez/clim_cli/internals/storage"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("46")).
			Background(lipgloss.Color("235"))

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Italic(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))
)

// Model represents the TUI model for device selection
type Model struct {
	devices     []*storage.DeviceHistory
	cursor      int
	selected    map[int]bool
	filter      string
	showDetail  bool
	detailIndex int
	quitting    bool
	err         error
}

// NewModel creates a new TUI model for device selection
func NewModel() Model {
	devices, err := storage.GetDeviceHistories()
	if err != nil {
		return Model{err: err}
	}

	return Model{
		devices:  devices,
		selected: make(map[int]bool),
	}
}

// Init initializes the TUI model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles TUI updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.devices)-1 {
				m.cursor++
			}

		case "enter", " ":
			if m.showDetail {
				m.showDetail = false
			} else {
				m.selected[m.cursor] = !m.selected[m.cursor]
			}

		case "d":
			if !m.showDetail && m.cursor < len(m.devices) {
				m.showDetail = true
				m.detailIndex = m.cursor
			}

		case "esc":
			if m.showDetail {
				m.showDetail = false
			}

		case "r":
			// Refresh devices
			devices, err := storage.GetDeviceHistories()
			if err != nil {
				m.err = err
			} else {
				m.devices = devices
				m.err = nil
				if m.cursor >= len(m.devices) {
					if len(m.devices) > 0 {
						m.cursor = len(m.devices) - 1
					} else {
						m.cursor = 0
					}
				}
			}
		}
	}

	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n\nPress 'q' to quit."
	}

	if m.quitting {
		selected := m.getSelectedDevices()
		if len(selected) > 0 {
			return "Selected devices:\n" + formatSelectedDevices(selected) + "\n\nPress 'q' to quit."
		}
		return "No devices selected.\n\nPress 'q' to quit."
	}

	if m.showDetail {
		return m.renderDetailView()
	}

	return m.renderListView()
}

// renderListView renders the main device list view
func (m Model) renderListView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Climate Device Selector"))
	b.WriteString("\n")

	if len(m.devices) == 0 {
		b.WriteString(infoStyle.Render("No devices found. Run 'clim_cli search' first."))
		b.WriteString("\n\n")
	} else {
		b.WriteString(fmt.Sprintf("Found %d device(s). Use ↑/↓ to navigate, Space/Enter to select, 'd' for details, 'r' to refresh.\n\n", len(m.devices)))

		for i, device := range m.devices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			checked := "☐"
			if m.selected[i] {
				checked = "☑"
			}

			line := fmt.Sprintf("%s %s %s (%s) - %s",
				cursor,
				checked,
				device.Device.Name,
				device.Device.IP,
				formatLastSeen(device.Device.LastSeenAt),
			)

			if i == m.cursor {
				b.WriteString(selectedStyle.Render(line))
			} else {
				b.WriteString(normalStyle.Render(line))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(infoStyle.Render("Controls: ↑/↓ navigate • Space select • d details • r refresh • q quit"))
	b.WriteString("\n")

	return b.String()
}

// renderDetailView renders the detailed view of a device
func (m Model) renderDetailView() string {
	if m.detailIndex >= len(m.devices) {
		return "Invalid device index"
	}

	device := m.devices[m.detailIndex]
	var b strings.Builder

	b.WriteString(titleStyle.Render(fmt.Sprintf("Device Details: %s", device.Device.Name)))
	b.WriteString("\n\n")

	// Basic information
	b.WriteString(lipgloss.NewStyle().Bold(true).Render("Basic Information:"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  IP Address: %s\n", device.Device.IP))
	b.WriteString(fmt.Sprintf("  MAC Address: %s\n", device.Device.MAC))
	b.WriteString(fmt.Sprintf("  Name: %s\n", device.Device.Name))
	b.WriteString(fmt.Sprintf("  Model: %s\n", device.Device.Model))
	b.WriteString(fmt.Sprintf("  Status: %s\n", device.Device.Status))
	b.WriteString(fmt.Sprintf("  First Seen: %s\n", device.Device.DiscoveredAt.Format("2006-01-02 15:04:05")))
	b.WriteString(fmt.Sprintf("  Last Seen: %s\n", device.Device.LastSeenAt.Format("2006-01-02 15:04:05")))
	b.WriteString("\n")

	// Basic info from device
	if len(device.Device.BasicInfo) > 0 {
		b.WriteString(lipgloss.NewStyle().Bold(true).Render("Device Basic Info:"))
		b.WriteString("\n")
		for key, value := range device.Device.BasicInfo {
			b.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
		}
		b.WriteString("\n")
	}

	// Control info from device
	if len(device.Device.ControlInfo) > 0 {
		b.WriteString(lipgloss.NewStyle().Bold(true).Render("Device Control Info:"))
		b.WriteString("\n")
		for key, value := range device.Device.ControlInfo {
			b.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
		}
		b.WriteString("\n")
	}

	// Recent changes
	if len(device.Changes) > 0 {
		b.WriteString(lipgloss.NewStyle().Bold(true).Render("Recent Changes:"))
		b.WriteString("\n")
		// Show last 5 changes
		start := len(device.Changes) - 5
		if start < 0 {
			start = 0
		}
		for i := start; i < len(device.Changes); i++ {
			change := device.Changes[i]
			b.WriteString(fmt.Sprintf("  %s: '%s' → '%s' (%s)\n",
				change.Field,
				change.OldValue,
				change.NewValue,
				change.ChangedAt.Format("2006-01-02 15:04:05")))
		}
	}

	b.WriteString("\n")
	b.WriteString(infoStyle.Render("Press Esc to go back, q to quit"))
	b.WriteString("\n")

	return b.String()
}

// getSelectedDevices returns the selected devices
func (m Model) getSelectedDevices() []*storage.DeviceHistory {
	var selected []*storage.DeviceHistory
	for i, device := range m.devices {
		if m.selected[i] {
			selected = append(selected, device)
		}
	}
	return selected
}

// formatSelectedDevices formats selected devices for display
func formatSelectedDevices(devices []*storage.DeviceHistory) string {
	var lines []string
	for _, device := range devices {
		lines = append(lines, fmt.Sprintf("• %s (%s)", device.Device.Name, device.Device.IP))
	}
	return strings.Join(lines, "\n")
}

// formatLastSeen formats the last seen time
func formatLastSeen(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	} else {
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	}
}
