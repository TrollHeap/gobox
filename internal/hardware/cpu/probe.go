package cpu

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// CPUInfo contient les informations du CPU
type CPUInfo struct {
	Architect  string  // Architecture (ex: "amd64", "arm64")
	VendorID   string  // Fabricant (ex: "GenuineIntel")
	ModelName  string  // Nom complet du CPU
	NumberCore int     // Nombre de cores physiques
	CacheSize  int64   // Taille totale cache en octets
	FreqMaxMHz float64 // Fréquence max en MHz
	FreqMinMHz float64 // Fréquence min en MHz
}

// Constantes de chemins et clés
const (
	VendorIDKey   = "vendor_id"
	ModelNameKey  = "model name"
	NumCoresKey   = "cpu cores"
	pathCpuInfo   = "/proc/cpuinfo"
	pathSystemCpu = "/sys/devices/system/cpu/"
)

// ParseCPUINFO parse les lignes de /proc/cpuinfo
func ParseCPUINFO(lines []string) (*CPUInfo, error) {
	info := &CPUInfo{}

	architect, err := getArchitect()
	if err != nil {
		return nil, fmt.Errorf("lecture architecture: %w", err)
	}
	info.Architect = architect

	for _, line := range lines {
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch key {
		case VendorIDKey:
			if info.VendorID == "" {
				info.VendorID = value
			}
		case ModelNameKey:
			if info.ModelName == "" {
				info.ModelName = value
			}
		case NumCoresKey:
			cores, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("parsing cpu cores: %w", err)
			}
			info.NumberCore = cores
		}
	}

	return info, nil
}

// readPathProcCPUInfo lit et filtre /proc/cpuinfo
func readPathProcCPUInfo() ([]string, error) {
	data, err := os.ReadFile(pathCpuInfo)
	if err != nil {
		return nil, fmt.Errorf("lecture %s: %w", pathCpuInfo, err)
	}

	trimData := strings.TrimSpace(string(data))
	result := make([]string, 0, 50)

	for line := range strings.SplitSeq(trimData, "\n") {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, ModelNameKey) ||
			strings.HasPrefix(line, NumCoresKey) ||
			strings.HasPrefix(line, VendorIDKey) {
			result = append(result, line)
		}
	}

	return result, nil
}

// getArchitect retourne l'architecture du CPU
func getArchitect() (string, error) {
	return runtime.GOARCH, nil
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

	// Fréquence max
	maxFreqPath := filepath.Join(basePath, "cpuinfo_max_freq")
	maxKHz, err := readFreqKHz(maxFreqPath)
	if err != nil {
		return 0, 0, fmt.Errorf("lecture freq max: %w", err)
	}
	maxMHz = float64(maxKHz) / 1000.0

	// Fréquence min
	minFreqPath := filepath.Join(basePath, "cpuinfo_min_freq")
	minKHz, err := readFreqKHz(minFreqPath)
	if err != nil {
		return 0, 0, fmt.Errorf("lecture freq min: %w", err)
	}
	minMHz = float64(minKHz) / 1000.0

	return minMHz, maxMHz, nil
}

// GetCPUInfo retourne les informations complètes du CPU
func GetCPUInfo() (*CPUInfo, error) {
	// Lire /proc/cpuinfo
	lines, err := readPathProcCPUInfo()
	if err != nil {
		return nil, err
	}

	info, err := ParseCPUINFO(lines)
	if err != nil {
		return nil, err
	}

	// Ajouter taille cache
	cacheSize, err := getTotalCacheSize()
	if err != nil {
		return nil, fmt.Errorf("calcul cache: %w", err)
	}
	info.CacheSize = cacheSize

	// TODO: Faire la moyenne avg fréquence
	// Ajouter fréquences
	minMHz, maxMHz, err := getCPUFrequencies()
	if err != nil {
		// Non-critique si cpufreq absent
		fmt.Fprintf(os.Stderr, "Attention: impossible de lire les fréquences: %v\n", err)
	} else {
		info.FreqMinMHz = minMHz
		info.FreqMaxMHz = maxMHz
	}

	return info, nil
}

// PrintMain affiche les informations CPU
func PrintMain() {
	info, err := GetCPUInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Vendor ID:       %s\n", info.VendorID)
	fmt.Printf("Model Name:      %s\n", info.ModelName)
	fmt.Printf("Architecture:    %s\n", info.Architect)
	fmt.Printf("CPU Cores:       %d\n", info.NumberCore)
	fmt.Printf("Total Cache:     %d octets (%.2f MB)\n",
		info.CacheSize, float64(info.CacheSize)/(1024*1024))

	if info.FreqMaxMHz > 0 {
		fmt.Printf("Freq Min:        %.0f MHz\n", info.FreqMinMHz)
		fmt.Printf("Freq Max:        %.0f MHz\n", info.FreqMaxMHz)
	}
}
