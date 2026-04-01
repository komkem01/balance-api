package config

import (
	"os"
	"os/exec"
	"strings"
	"sync"
)

var version = "0.0.0"

var resolvedVersionCache string
var resolvedVersionOnce sync.Once

func resolvedVersion() string {
	resolvedVersionOnce.Do(func() {
		envVersion := strings.TrimSpace(os.Getenv("APP_VERSION"))
		if envVersion != "" {
			resolvedVersionCache = envVersion
			return
		}

		buildVersion := strings.TrimSpace(version)
		if buildVersion != "" && buildVersion != "0.0.0" {
			resolvedVersionCache = buildVersion
			return
		}

		cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
		cmd.Dir = "."
		output, err := cmd.Output()
		if err == nil {
			tagVersion := strings.TrimSpace(string(output))
			if tagVersion != "" {
				resolvedVersionCache = tagVersion
				return
			}
		}

		if buildVersion == "" {
			buildVersion = "0.0.0"
		}
		resolvedVersionCache = buildVersion
	})

	return resolvedVersionCache
}
