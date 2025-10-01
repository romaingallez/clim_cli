# CLIM CLI

A CLI application to control climate devices at work.

It allows setting and getting parameters of climate devices, discovering devices on the network, and tracking device changes over time.

## Installation

Download the latest release from the [release page](https://github.com/romaingallez/clim_cli/releases)
or build it yourself with the following command:


To get a version using golang install
```bash
## With a specific release
go install github.com/romaingallez/clim_cli@v0.x.x
## the latest commit (beta)
go install github.com/romaingallez/clim_cli@latest
```

## Usage

```bash
clim-cli --help
```

## Device Discovery and Management

### Search for Devices

Discover climate devices on your network:

```bash
# Basic search
clim-cli search

# Search with custom interface and timeout
clim-cli search -I eth0 --timeout 10

# Interactive search with TUI for device selection
clim-cli search --tui
```

### Browse Stored Devices

Browse previously discovered devices:

```bash
# Interactive browser
clim-cli browse

# List all stored devices
clim-cli list
```

### Device Storage

Devices are automatically stored in `~/.config/clim_cli/devices.json` with:
- Historical snapshots with timestamps
- Change tracking (IP changes, name changes, etc.)
- Device information and status

## Commands

- `search` - Discover climate devices on the network
- `browse` - Interactive device browser
- `list` - List all stored devices
- `get` - Get current climate device settings
- `set` - Set climate device parameters
