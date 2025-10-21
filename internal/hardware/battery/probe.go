package battery

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func readInt(batteryPath string, filename string) (int, error) {
	data, err := os.ReadFile(filepath.Join(batteryPath, filename))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

func readFloat(path, filename string) (float64, error) {
	data, err := os.ReadFile(filepath.Join(path, filename))
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
}

func readString(batteryPath string, filename string) (string, error) {
	data, err := os.ReadFile(filepath.Join(batteryPath, filename))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

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
