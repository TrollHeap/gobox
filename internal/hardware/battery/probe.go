package battery

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// readFile lit et nettoie un fichier sysfs
func readFile(basePath, filename string) (string, error) {
	data, err := os.ReadFile(filepath.Join(basePath, filename))
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", filename, err)
	}
	return strings.TrimSpace(string(data)), nil
}

func readInt(batteryPath, filename string) (int, error) {
	s, err := readFile(batteryPath, filename)
	if err != nil {
		return 0, err // ✅ déjà enrobée par readFile
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parsing %s as int: %w", filename, err)
	}
	return val, nil
}

func readFloat(batteryPath, filename string) (float64, error) {
	s, err := readFile(batteryPath, filename)
	if err != nil {
		return 0, err
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing %s as float: %w", filename, err)
	}
	return val, nil
}

func readString(batteryPath string, filename string) (string, error) {
	return readFile(batteryPath, filename)
}

func findBatteryPath() (string, error) {
	entries, err := os.ReadDir("/sys/class/power_supply")
	if err != nil {
		return "", fmt.Errorf("listing power_supply: %w", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, "BAT") {
			return filepath.Join("/sys/class/power_supply", name), nil
		}
	}
	return "", errors.New("no battery found")
}
