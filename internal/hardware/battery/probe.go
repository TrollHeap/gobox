package battery

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	pathPowerSupply = "/sys/class/power_supply"
)

func findBatteryPath() (string, error) {
	entries, err := os.ReadDir(pathPowerSupply)
	if err != nil {
		return "", fmt.Errorf("listing power_supply: %w", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, "BAT") {
			return filepath.Join(pathPowerSupply, name), nil
		}
	}
	return "", errors.New("no battery found")
}
