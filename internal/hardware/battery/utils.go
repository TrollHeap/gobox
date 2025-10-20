package battery

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func findBatteryPath() (string, error) {
	entries, err := os.ReadDir("/sys/class/power_supply")
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, "BAT") {
			return filepath.Join("/sys/class/power_supply", name), nil
		}
	}
	return "", fmt.Errorf("no battery found")
}

func convertWattToAmpere(path string, energyFull float64) float64 {
	voltDesign, _ := readFloat(path, "voltage_min_design")

	mAH := (energyFull / 1_000_000) * 1000 / (voltDesign / 1_000_000)
	return mAH
}
