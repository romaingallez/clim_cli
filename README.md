# CLIM CLI

A CLI application to control climate devices at work.

It allows setting and getting parameters of climate devices, discovering devices on the network, and tracking device changes over time.

## Installation

Download the latest release from the release page or build it yourself:

```bash
# With a specific version
go install github.com/romaingallez/clim_cli@v0.x.x
# Latest (beta)
go install github.com/romaingallez/clim_cli@latest
```

- Requires Go 1.23+

## Prerequisites

- Network discovery uses `arp-scan`.
  - Debian/Ubuntu: `sudo apt-get install arp-scan`
  - macOS (Homebrew): `brew install arp-scan`

If `arp-scan` is not installed, the `search` command returns an actionable error with install hints.

## Usage

```bash
clim-cli --help
```

### Recommended Workflow

- Search → Browse → Get/Set
  - `clim-cli search --tui` to discover and interactively pick your device
  - This saves the selection to the config so `get`/`set` can use it
  - Or pass `--ip` directly to `get`/`set`

### Default IP Behavior

- The default `ip` is empty. If you run `get` or `set` without `--ip` and no device is selected in the config, you’ll be guided to run `search --tui` or `browse`.

## Device Discovery and Management

### Search for Devices

Discover climate devices on your network (requires `arp-scan`):

```bash
# Basic search
clim-cli search

# Search with custom interface and timeout
clim-cli search -I eth0 --timeout 10

# Limit concurrency for HTTP info fetches
clim-cli search -w 5

# Interactive search with TUI for device selection
clim-cli search --tui
```

Notes:
- `--workers` controls parallel HTTP fetches of per-device details; `arp-scan` itself runs once.
- If your default interface is a VPN/virtual adapter, pass a physical one with `-I`.

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
