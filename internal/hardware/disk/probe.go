package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	pathRoot        = "/sys/block/"
	sectorSizeBytes = 512
)

func listPartitions(diskName string) ([]string, error) {
	diskPath := filepath.Join(pathRoot, diskName)

	entries, err := os.ReadDir(diskPath)
	if err != nil {
		return nil, fmt.Errorf("lecture %s: %w", diskPath, err)
	}

	partitions := make([]string, 0, 4)

	for _, entry := range entries {
		// Garder seulement répertoires (partitions sont des sous-répertoires)
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Vérifier si c'est une partition du disque
		if suffix, found := strings.CutPrefix(name, diskName); found {
			if suffix != "" && (strings.HasPrefix(suffix, "p") ||
				(len(suffix) > 0 && suffix[0] >= '0' && suffix[0] <= '9')) {
				partitions = append(partitions, name)
			}
		}
	}

	return partitions, nil
}

func validateDiskName(diskName string) error {
	if strings.Contains(diskName, "..") || strings.Contains(diskName, "/") {
		return fmt.Errorf("nom de disque invalide : %s", diskName)
	}
	if diskName == "" {
		return fmt.Errorf("nom de disque vide")
	}
	return nil
}

func readDiskSize(diskName string) (int64, error) {
	path := filepath.Join(pathRoot, diskName, "size")
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("lecture %s: %w", path, err)
	}

	dataStr := strings.TrimSpace(string(data))
	dataInt, err := strconv.ParseInt(dataStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("conversion taille disque %s: %w", diskName, err)
	}

	sizeBytes := dataInt * sectorSizeBytes
	return sizeBytes, nil
}

func readDiskType(diskName string) (string, error) {
	path := filepath.Join(pathRoot, diskName, "queue", "rotational")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("lecteur %s : %w", path, err)
	}
	dataStr := strings.TrimSpace(string(data))
	typeDIsk := "HDD"
	if dataStr == "0" {
		typeDIsk = "SSD"
	}
	return typeDIsk, nil
}

func readModelAndVendor(diskName string) (string, string, error) {
	// Chemins potentiels
	vendorPath := filepath.Join(pathRoot, diskName, "device", "vendor")
	modelPath := filepath.Join(pathRoot, diskName, "device", "model")

	// Lecture vendor (optionnel)
	vendor := ""
	if data, err := os.ReadFile(vendorPath); err == nil {
		vendor = strings.TrimSpace(string(data))
	} else if !os.IsNotExist(err) {
		// Si erreur autre que fichier absent, renvoyer erreur
		return "", "", fmt.Errorf("lecture vendor %s: %w", vendorPath, err)
	}

	// Lecture modèle (obligatoire)
	modelData, err := os.ReadFile(modelPath)
	if err != nil {
		return "", "", fmt.Errorf("lecture model %s: %w", modelPath, err)
	}
	modelStr := strings.TrimSpace(string(modelData))

	// Si vendor vide, tenter d'extraire vendor + modèle à partir de modelStr
	if vendor == "" {
		parts := strings.SplitN(modelStr, "_", 2)
		if len(parts) == 2 {
			vendor, modelStr = parts[0], parts[1]
		} else {
			// Pas de séparation possible, vendor vide, modèle complet
			vendor = ""
		}
	}

	return vendor, modelStr, nil
}
