package battery

func convertWattToAmpere(path string, energyFull float64) float64 {
	voltDesign, _ := readFloat(path, "voltage_min_design")

	mAH := (energyFull / 1_000_000) * 1000 / (voltDesign / 1_000_000)
	return mAH
}
