package utils

import (
	"fmt"
	"os"
	"strings"

	"gobox/internal/sysfs"
)

func ReadFileOptional(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// Ne PAS retourner l'erreur si fichier absent
		if os.IsNotExist(err) {
			return "", nil // Pas d'erreur, juste vide
		}
		return "", err // Autre erreur (permissions, etc.)
	}
	return strings.TrimSpace(string(data)), nil
}

func ConvertWattToAmpere(path string, energyFull float64) (float64, error) {
	voltDesign, err := sysfs.ReadFloat(path, "voltage_min_design")
	if err != nil {
		return 0, fmt.Errorf("reading voltage: %w", err)
	}

	mAH := (energyFull / 1_000_000) * 1000 / (voltDesign / 1_000_000)
	return mAH, nil
}
