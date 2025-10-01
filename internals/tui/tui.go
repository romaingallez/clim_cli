package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/romaingallez/clim_cli/internals/storage"
)

// RunDeviceSelector runs the interactive device selector TUI
func RunDeviceSelector() ([]*storage.DeviceHistory, error) {
	model := NewModel()
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("error running TUI: %v", err)
	}

	m, ok := finalModel.(Model)
	if !ok {
		return nil, fmt.Errorf("unexpected model type")
	}

	if m.err != nil {
		return nil, m.err
	}

	return m.getSelectedDevices(), nil
}

// RunDeviceBrowser runs the device browser TUI (non-interactive, just for browsing)
func RunDeviceBrowser() error {
	model := NewModel()
	p := tea.NewProgram(model)

	_, err := p.Run()
	return err
}

// GetSelectedDeviceIPs returns just the IP addresses of selected devices
func GetSelectedDeviceIPs() ([]string, error) {
	devices, err := RunDeviceSelector()
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, device := range devices {
		ips = append(ips, device.Device.IP)
	}

	return ips, nil
}

// PrintDeviceSummary prints a summary of stored devices to stdout
func PrintDeviceSummary() error {
	devices, err := storage.GetDeviceHistories()
	if err != nil {
		return fmt.Errorf("failed to load devices: %v", err)
	}

	if len(devices) == 0 {
		fmt.Println("No devices found in storage.")
		fmt.Println("Run 'clim_cli search' to discover devices.")
		return nil
	}

	fmt.Printf("Found %d device(s) in storage:\n\n", len(devices))

	for i, device := range devices {
		fmt.Printf("%d. %s\n", i+1, device.Device.Name)
		fmt.Printf("   IP: %s\n", device.Device.IP)
		fmt.Printf("   MAC: %s\n", device.Device.MAC)
		fmt.Printf("   Status: %s\n", device.Device.Status)
		fmt.Printf("   Last Seen: %s\n", device.Device.LastSeenAt.Format("2006-01-02 15:04:05"))

		if len(device.Changes) > 0 {
			fmt.Printf("   Changes: %d detected\n", len(device.Changes))
		}

		if len(device.Device.BasicInfo) > 0 {
			fmt.Printf("   Basic Info: %d fields\n", len(device.Device.BasicInfo))
		}

		if len(device.Device.ControlInfo) > 0 {
			fmt.Printf("   Control Info: %d fields\n", len(device.Device.ControlInfo))
		}

		fmt.Println()
	}

	return nil
}
