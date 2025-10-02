/*
Copyright © 2023 GALLEZ Romain
*/
package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

// These variables are set during build time via ldflags
var (
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
	GoVersion = runtime.Version()
	Platform  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

// init attempts to detect version information from build info
func init() {
	// If version is still "dev", try to get it from build info
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok {
			// Look for vcs info in build settings
			for _, setting := range info.Settings {
				switch setting.Key {
				case "vcs.revision":
					if GitCommit == "none" {
						GitCommit = setting.Value
					}
				case "vcs.modified":
					// If the working directory was modified, append "-dirty"
					if setting.Value == "true" && !strings.Contains(Version, "dirty") {
						if Version == "dev" {
							Version = fmt.Sprintf("dev-%s", GitCommit[:8])
						}
						if len(GitCommit) > 8 {
							GitCommit += "-dirty"
						}
					}
				}
			}

			// Try to extract version from module path if it contains a version
			if strings.Contains(info.Main.Path, "@") {
				parts := strings.Split(info.Main.Path, "@")
				if len(parts) > 1 && parts[1] != "latest" && parts[1] != "main" {
					Version = parts[1]
				}
			}

			// If we still have "dev" and have a commit, use commit-based version
			if Version == "dev" && GitCommit != "none" && len(GitCommit) > 7 {
				Version = fmt.Sprintf("dev-%s", GitCommit[:8])
			}
		}
	}

	// Set build date if unknown
	if BuildDate == "unknown" {
		BuildDate = time.Now().UTC().Format("2006-01-02_15:04:05")
	}
}

// GetVersionFromGit attempts to get version information from git
func GetVersionFromGit() (string, string, error) {
	// This would require importing os/exec, but we'll keep it simple
	// for now and rely on build info
	return "", "", fmt.Errorf("not implemented")
}

// Info holds version information
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

// Get returns version information
func Get() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: GoVersion,
		Platform:  Platform,
	}
}

// String returns a string representation of version info
func (v Info) String() string {
	return fmt.Sprintf("clim_cli %s\nGit Commit: %s\nBuild Date: %s\nGo Version: %s\nPlatform: %s",
		v.Version, v.GitCommit, v.BuildDate, v.GoVersion, v.Platform)
}

// Short returns a short version string
func (v Info) Short() string {
	return fmt.Sprintf("clim_cli %s", v.Version)
}

// IsDev returns true if this is a development build
func (v Info) IsDev() bool {
	return Version == "dev" || strings.Contains(Version, "dev")
}

// IsPreRelease returns true if this is a pre-release version
func (v Info) IsPreRelease() bool {
	return strings.Contains(Version, "alpha") || strings.Contains(Version, "beta") || strings.Contains(Version, "rc")
}

// GetBuildType returns the type of build (release, dev, pre-release)
func (v Info) GetBuildType() string {
	if v.IsDev() {
		return "development"
	}
	if v.IsPreRelease() {
		return "pre-release"
	}
	return "release"
}
