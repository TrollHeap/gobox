package cpu

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CacheID identifie un cache physique unique
type CacheID struct {
	Level      int
	Type       string
	SharedCPUs string
}

// readLevel lit le niveau de cache depuis index*/level
func readLevel(index string) (int, error) {
	levelPath := filepath.Join(index, "level")
	levelBytes, err := os.ReadFile(levelPath)
	if err != nil {
		return 0, err
	}

	levelStr := strings.TrimSpace(string(levelBytes))
	return strconv.Atoi(levelStr)
}

// readType lit le type de cache depuis index*/type
func readType(index string) (string, error) {
	typePath := filepath.Join(index, "type")
	typeBytes, err := os.ReadFile(typePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(typeBytes)), nil
}

// readSharedCPU lit la liste des CPUs partageant ce cache
func readSharedCPU(index string) (string, error) {
	sharedPath := filepath.Join(index, "shared_cpu_list")
	sharedBytes, err := os.ReadFile(sharedPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(sharedBytes)), nil
}

// readSize lit et parse la taille du cache
func readSize(index string) (int64, error) {
	sizePath := filepath.Join(index, "size")
	sizeBytes, err := os.ReadFile(sizePath)
	if err != nil {
		return 0, err
	}

	sizeStr := strings.TrimSpace(string(sizeBytes))
	return parseSize(sizeStr)
}

// parseSize convertit "64K" → 65536 octets
func parseSize(sizeStr string) (int64, error) {
	sizeStr = strings.TrimSpace(sizeStr)
	if len(sizeStr) == 0 {
		return 0, fmt.Errorf("taille vide")
	}

	unit := sizeStr[len(sizeStr)-1]

	// Cas sans unité (tous chiffres)
	if unit >= '0' && unit <= '9' {
		return strconv.ParseInt(sizeStr, 10, 64)
	}

	// Parser la valeur numérique
	valueStr := sizeStr[:len(sizeStr)-1]
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, err
	}

	// Appliquer le multiplicateur
	switch unit {
	case 'K', 'k':
		return value * 1024, nil
	case 'M', 'm':
		return value * 1024 * 1024, nil
	case 'G', 'g':
		return value * 1024 * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("unité inconnue: %c", unit)
	}
}

// getTotalCacheSize calcule la taille totale des caches avec déduplication
func getTotalCacheSize() (int64, error) {
	var totalCache int64
	seenCaches := make(map[CacheID]bool)

	pattern := filepath.Join(pathSystemCpu, "cpu[0-9]*")
	cpuPaths, err := filepath.Glob(pattern)
	if err != nil {
		return 0, fmt.Errorf("glob CPUs: %w", err)
	}

	for _, cpu := range cpuPaths {
		cachePath := filepath.Join(cpu, "cache")
		indexPaths, err := filepath.Glob(filepath.Join(cachePath, "index*"))
		if err != nil {
			continue
		}

		for _, index := range indexPaths {
			level, err := readLevel(index)
			if err != nil {
				continue
			}

			cacheType, err := readType(index)
			if err != nil {
				continue
			}

			sharedCPUList, err := readSharedCPU(index)
			if err != nil {
				continue
			}

			cacheID := CacheID{
				Level:      level,
				Type:       cacheType,
				SharedCPUs: sharedCPUList,
			}

			if seenCaches[cacheID] {
				continue
			}

			size, err := readSize(index)
			if err != nil {
				continue
			}

			totalCache += size
			seenCaches[cacheID] = true
		}
	}

	return totalCache, nil
}
