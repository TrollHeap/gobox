package battery

import "fmt"

func convertWattToAmpere(path string, energyFull float64) (float64, error) {
	voltDesign, err := readFloat(path, "voltage_min_design")
	if err != nil {
		return 0, fmt.Errorf("reading voltage: %w", err)
	}

	mAH := (energyFull / 1_000_000) * 1000 / (voltDesign / 1_000_000)
	return mAH, nil
}
