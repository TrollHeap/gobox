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

// readFile lit un fichier et retourne son contenu nettoyé
func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// readLevel lit le niveau de cache
func readLevel(index string) (int, error) {
	content, err := readFile(filepath.Join(index, "level"))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(content)
}

// readType lit le type de cache
func readType(index string) (string, error) {
	return readFile(filepath.Join(index, "type"))
}

// readSharedCPU lit la liste des CPUs partageant ce cache
func readSharedCPU(index string) (string, error) {
	return readFile(filepath.Join(index, "shared_cpu_list"))
}

// readSize lit et parse la taille du cache
func readSize(index string) (int64, error) {
	sizeStr, err := readFile(filepath.Join(index, "size"))
	if err != nil {
		return 0, err
	}
	return parseSize(sizeStr)
}

// parseSize convertit "64K" → 65536 octets
func parseSize(sizeStr string) (int64, error) {
	if sizeStr == "" {
		return 0, fmt.Errorf("taille vide")
	}

	lastChar := sizeStr[len(sizeStr)-1]

	// Cas sans unité
	if lastChar >= '0' && lastChar <= '9' {
		return strconv.ParseInt(sizeStr, 10, 64)
	}

	value, err := strconv.ParseInt(sizeStr[:len(sizeStr)-1], 10, 64)
	if err != nil {
		return 0, err
	}

	switch lastChar {
	case 'K', 'k':
		return value << 10, nil // *1024
	case 'M', 'm':
		return value << 20, nil // *1024*1024
	case 'G', 'g':
		return value << 30, nil // *1024*1024*1024
	default:
		return 0, fmt.Errorf("unité inconnue: %c", lastChar)
	}
}

// getTotalCacheSize calcule la taille totale des caches avec déduplication
func getTotalCacheSize() (int64, error) {
	var totalCache int64
	seenCaches := make(map[CacheID]struct{}) // struct{} au lieu de bool

	cpuPaths, err := filepath.Glob(filepath.Join(pathSystemCpu, "cpu[0-9]*"))
	if err != nil {
		return 0, fmt.Errorf("glob CPUs: %w", err)
	}

	for _, cpu := range cpuPaths {
		indexPaths, err := filepath.Glob(filepath.Join(cpu, "cache", "index*"))
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

			sharedCPU, err := readSharedCPU(index)
			if err != nil {
				continue
			}

			cacheID := CacheID{
				Level:      level,
				Type:       cacheType,
				SharedCPUs: sharedCPU,
			}

			if _, seen := seenCaches[cacheID]; seen {
				continue
			}

			size, err := readSize(index)
			if err != nil {
				continue
			}

			totalCache += size
			seenCaches[cacheID] = struct{}{}
		}
	}

	return totalCache, nil
}
