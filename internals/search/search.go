package search

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/romaingallez/clim_cli/internals/api"
	"github.com/romaingallez/clim_cli/internals/config"
)

// Device represents a found climate device
type Device struct {
	IP          string
	Status      string
	Name        string
	Model       string
	MAC         string
	BasicInfo   map[string]string
	ControlInfo map[string]string
}

// SearchDevices searches for climate devices in the specified IP range using arp-scan
func SearchDevices(ifaceString string, timeout int, workers int) ([]Device, error) {
	log.Printf("Searching on interface: %s, timeout: %d, workers: %d", ifaceString, timeout, workers)

	// Validate interface exists
	iface, err := net.InterfaceByName(ifaceString)
	if err != nil {
		return nil, fmt.Errorf("interface %s not found: %v", ifaceString, err)
	}

	// Quick preflight for arp-scan
	if _, err := exec.LookPath("arp-scan"); err != nil {
		return nil, fmt.Errorf("arp-scan is not installed or not in PATH: %w", err)
	}

	// Get the network address for the interface
	networkAddr, err := getNetworkAddress(iface)
	if err != nil {
		return nil, fmt.Errorf("failed to get network address for interface %s: %v", ifaceString, err)
	}

	// Execute arp-scan command
	devices, err := executeArpScan(ifaceString, networkAddr, timeout)
	if err != nil {
		return nil, fmt.Errorf("arp-scan failed: %v", err)
	}

	log.Printf("Found %d devices", len(devices))
	return devices, nil
}

// getNetworkAddress returns the network address for the given interface
func getNetworkAddress(iface *net.Interface) (string, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if ip4 := ipNet.IP.To4(); ip4 != nil {
				// Skip localhost
				if ip4[0] == 127 {
					continue
				}
				// Return the network in CIDR notation
				return ipNet.String(), nil
			}
		}
	}

	return "", errors.New("no valid IPv4 address found on interface")
}

// executeArpScan runs the arp-scan command and parses its output
func executeArpScan(iface string, networkAddr string, timeout int) ([]Device, error) {
	// Build the arp-scan command
	cmd := exec.Command("arp-scan")

	switch {
	case iface != "" && networkAddr != "":
		cmd.Args = append(cmd.Args, "-I", iface)
	case iface != "":
		cmd.Args = append(cmd.Args, "-I", iface)
	case networkAddr != "":
		// Could derive interface from CIDR if needed
	}

	cmd.Args = append(cmd.Args, "--timeout", fmt.Sprintf("%d", timeout))
	cmd.Args = append(cmd.Args, "--localnet", "-x", "--format=${ip};${mac};${vendor}")

	log.Printf("Executing command: %s", strings.Join(cmd.Args, " "))

	// Execute the command
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute arp-scan: %v", err)
	}

	// Parse the output
	devices, err := parseArpScanOutput(string(output))
	if err != nil {
		return nil, fmt.Errorf("failed to parse arp-scan output: %v", err)
	}

	return devices, nil
}

// parseArpScanOutput parses the output from arp-scan command
// Expected format: IP;MAC;Vendor
func parseArpScanOutput(output string) ([]Device, error) {
	var devices []Device

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and header lines
		if line == "" || strings.Contains(line, "Starting arp-scan") ||
			strings.Contains(line, "Interface:") || strings.Contains(line, "Starting:") ||
			strings.Contains(line, "Ending arp-scan") || strings.Contains(line, "packets") {
			continue
		}

		// Parse the line: IP;MAC;Vendor
		parts := strings.Split(line, ";")
		if len(parts) < 2 {
			continue // Skip malformed lines
		}

		ip := strings.TrimSpace(parts[0])
		mac := strings.TrimSpace(parts[1])
		vendor := ""
		if len(parts) > 2 {
			vendor = strings.TrimSpace(parts[2])
		}

		// Validate IP address
		if net.ParseIP(ip) == nil {
			continue
		}

		device := Device{
			IP:     ip,
			MAC:    mac,
			Name:   vendor,
			Status: "online",
		}

		devices = append(devices, device)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

// FuzzySearchDevices searches for devices matching a fuzzy pattern and saves AC manufacturer MACs to config
func FuzzySearchDevices(ifaceString string, timeout int, workers int, pattern string) ([]Device, error) {
	// First, get all devices
	devices, err := SearchDevices(ifaceString, timeout, workers)
	if err != nil {
		return nil, err
	}

	// Filter devices using fuzzy search
	filteredDevices := fuzzyFilter(devices, pattern)

	// Parallel fetch basic/control info per device with worker cap
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup
	for i := range filteredDevices {
		i := i
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			log.Printf("Getting clim info for device %s", filteredDevices[i].IP)
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
			defer cancel()
			basicInfo, berr := api.FetchBasicInfo(ctx, filteredDevices[i].IP)
			if berr == nil {
				filteredDevices[i].BasicInfo = basicInfo
				if name, ok := basicInfo["name"]; ok && name != "" {
					filteredDevices[i].Name = name
				}
			}
			controlInfo, cerr := api.FetchControlInfo(ctx, filteredDevices[i].IP)
			if cerr == nil {
				filteredDevices[i].ControlInfo = controlInfo
			}
		}()
	}
	wg.Wait()

	// Save AC manufacturer MACs to config via central config package
	var acMACs []string
	for _, d := range filteredDevices {
		if strings.Contains(strings.ToLower(d.Name), "murata") {
			acMACs = append(acMACs, d.MAC)
		}
	}
	if err := config.SaveACManufacturerMACs(acMACs); err != nil {
		log.Printf("Warning: Failed to save AC manufacturer MACs to config: %v", err)
	}

	return filteredDevices, nil
}

// fuzzyFilter filters devices based on fuzzy pattern matching
func fuzzyFilter(devices []Device, pattern string) []Device {
	var filtered []Device

	// Convert pattern to lowercase for case-insensitive matching
	pattern = strings.ToLower(pattern)

	for _, device := range devices {
		// Check if device name contains the pattern (case-insensitive)
		if strings.Contains(strings.ToLower(device.Name), pattern) {
			filtered = append(filtered, device)
		}
	}

	return filtered
}
