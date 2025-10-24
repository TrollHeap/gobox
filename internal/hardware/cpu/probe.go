package cpu

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// CPUInfo contient les informations du CPU
type CPUInfo struct {
	Architect  string
	VendorID   string
	ModelName  string
	NumberCore int
	CacheSize  int64
	FreqMaxMHz float64
	FreqMinMHz float64
}

const (
	vendorIDKey   = "vendor_id"
	modelNameKey  = "model name"
	numCoresKey   = "cpu cores"
	pathCpuInfo   = "/proc/cpuinfo"
	pathSystemCpu = "/sys/devices/system/cpu/"
)

// parseCPUInfo parse les lignes de /proc/cpuinfo
func parseCPUInfo(lines []string) CPUInfo {
	info := CPUInfo{
		Architect: runtime.GOARCH, // Pas besoin de fonction séparée
	}

	for _, line := range lines {
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch key {
		case vendorIDKey:
			if info.VendorID == "" {
				info.VendorID = value
			}
		case modelNameKey:
			if info.ModelName == "" {
				info.ModelName = value
			}
		case numCoresKey:
			if cores, err := strconv.Atoi(value); err == nil {
				info.NumberCore = cores
			}
		}
	}

	return info
}

// readCPUInfo lit et filtre /proc/cpuinfo avec lecture bufférisée
func readCPUInfo() ([]string, error) {
	file, err := os.Open(pathCpuInfo)
	if err != nil {
		return nil, fmt.Errorf("ouverture %s: %w", pathCpuInfo, err)
	}
	defer file.Close()

	result := make([]string, 0, 32)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		// Optimisation : vérifier le préfixe avant de faire strings.HasPrefix complet
		if len(line) > 0 && (line[0] == 'm' || line[0] == 'c' || line[0] == 'v') {
			if strings.HasPrefix(line, modelNameKey) ||
				strings.HasPrefix(line, numCoresKey) ||
				strings.HasPrefix(line, vendorIDKey) {
				result = append(result, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("lecture %s: %w", pathCpuInfo, err)
	}

	return result, nil
}

// readFreqKHz lit une fréquence en KHz depuis sysfs
func readFreqKHz(path string) (int64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	freqStr := strings.TrimSpace(string(data))
	return strconv.ParseInt(freqStr, 10, 64)
}

// getCPUFrequencies retourne les fréquences min et max du CPU
func getCPUFrequencies() (minMHz, maxMHz float64, err error) {
	basePath := filepath.Join(pathSystemCpu, "cpu0", "cpufreq")

	maxKHz, err := readFreqKHz(filepath.Join(basePath, "cpuinfo_max_freq"))
	if err != nil {
		return 0, 0, fmt.Errorf("freq max: %w", err)
	}

	minKHz, err := readFreqKHz(filepath.Join(basePath, "cpuinfo_min_freq"))
	if err != nil {
		return 0, 0, fmt.Errorf("freq min: %w", err)
	}

	// Optimisation : division unique au lieu de deux constantes .0
	const khzToMhz = 1000.0
	return float64(minKHz) / khzToMhz, float64(maxKHz) / khzToMhz, nil
}

// GetCPUInfo retourne les informations complètes du CPU
func GetCPUInfo() (CPUInfo, error) {
	lines, err := readCPUInfo()
	if err != nil {
		return CPUInfo{}, err
	}

	info := parseCPUInfo(lines)

	cacheSize, err := getTotalCacheSize()
	if err != nil {
		return CPUInfo{}, fmt.Errorf("calcul cache: %w", err)
	}
	info.CacheSize = cacheSize

	// Fréquences optionnelles - ne pas échouer si absentes
	minMHz, maxMHz, _ := getCPUFrequencies()
	info.FreqMinMHz = minMHz
	info.FreqMaxMHz = maxMHz

	return info, nil
}
