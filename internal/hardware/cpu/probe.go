package cpu

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	VendorIDKey  = "vendor_id"
	ModelNameKey = "model name"
	NumCoresKey  = "cpu cores"
	// FreqMHzKey   = "cpu MHz" // NOTE: Voir si il y a toujours MHz sur les processeurs
)

func ParseCPUINFO(lines []string) (*CPUInfo, error) {
	info := &CPUInfo{}

	// ✅ Appels AVANT la boucle (1 seule fois)
	architect, err := getArchitect()
	if err != nil {
		return &CPUInfo{}, fmt.Errorf("Erreur architect: %w", err)
	}
	info.Architect = architect

	// cacheSize, err := getTotalCacheSize()
	// if err != nil {
	// 	return &CPUInfo{}, fmt.Errorf("Erreur cache: %w", err)
	// }
	// info.CacheSize = cacheSize

	// Boucle de parsing
	for _, line := range lines {
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch key {
		case VendorIDKey:
			if info.VendorID == "" { // ← Garde première occurrence
				info.VendorID = value
			}
		case ModelNameKey:
			if info.ModelName == "" {
				info.ModelName = value
			}
		case NumCoresKey:
			cores, err := strconv.Atoi(value)
			if err != nil {
				return &CPUInfo{}, fmt.Errorf("parsing cpu cores: %w", err)
			}
			info.NumberCore = cores
		}
	}

	return info, nil
}

func getArchitect() (string, error) {
	return runtime.GOARCH, nil
}

func ReadPathProcCPUInfo() ([]string, error) {
	path := "/proc/cpuinfo"

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Problème de lecture du fichier %s: %w", path, err)
	}
	trimData := strings.TrimSpace(string(data))

	result := make([]string, 0, 50)

	for line := range strings.SplitSeq(string(trimData), "\n") {
		if strings.HasPrefix(line, ModelNameKey) ||
			// strings.HasPrefix(line, FreqMHzKey) ||
			strings.HasPrefix(line, NumCoresKey) ||
			strings.HasPrefix(line, VendorIDKey) {
			result = append(result, line)
		}
	}

	return result, nil
}

func parseSize(sizeStr string) (int64, error) {
	// sizeStr exemple : "32K", "256K", "8M"
	sizeStr = strings.TrimSpace(sizeStr)
	if len(sizeStr) == 0 {
		return 0, fmt.Errorf("taille vide")
	}
	unit := sizeStr[len(sizeStr)-1]
	valueStr := sizeStr[:len(sizeStr)-1]
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, err
	}
	switch unit {
	case 'K', 'k':
		return value * 1024, nil
	case 'M', 'm':
		return value * 1024 * 1024, nil
	case 'G', 'g':
		return value * 1024 * 1024 * 1024, nil
	default:
		// pas d’unité, on retourne la valeur brute
		return strconv.ParseInt(sizeStr, 10, 64)
	}
}

// TODO: Implementer getTotalCacheSize
// func getTotalCacheSize() (int64, error) {
// 	var totalCache int64
// 	seenCaches := make(map[CacheID]bool)
//
// 	cpuPaths, err := "/sys/devices/system/cpu/"
// }
