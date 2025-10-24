package utils

import (
	"fmt"

	"gobox/internal/sysfs"
)

func ConvertWattToAmpere(path string, energyFull float64) (float64, error) {
	voltDesign, err := sysfs.ReadFloat(path, "voltage_min_design")
	if err != nil {
		return 0, fmt.Errorf("reading voltage: %w", err)
	}

	mAH := (energyFull / 1_000_000) * 1000 / (voltDesign / 1_000_000)
	return mAH, nil
}
