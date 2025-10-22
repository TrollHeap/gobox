package battery

import (
	"fmt"
	"log"
)

func GetBatteryInfo() (BatteryInfo, error) {
	path, err := findBatteryPath()
	if err != nil {
		return BatteryInfo{}, err
	}

	// Champs critiques (obligatoires selon la spec Linux)
	capacity, err := readInt(path, "capacity")
	if err != nil {
		return BatteryInfo{}, fmt.Errorf("reading capacity: %w", err)
	}

	status, err := readString(path, "status")
	if err != nil {
		return BatteryInfo{}, fmt.Errorf("reading status: %w", err)
	}

	// Champs optionnels (peuvent être absents selon le matériel/driver)
	manufacturer, err := readString(path, "manufacturer")
	if err != nil {
		log.Printf("manufacturer unavailable: %v", err)
		manufacturer = "N/A"
	}

	model, err := readString(path, "model_name")
	if err != nil {
		log.Printf("model_name unavailable: %v", err)
		model = "N/A"
	}

	serialNumber, err := readString(path, "serial_number")
	if err != nil {
		log.Printf("serial_number unavailable: %v", err)
		serialNumber = "N/A"
	}

	technology, err := readString(path, "technology")
	if err != nil {
		log.Printf("technology unavailable: %v", err)
		technology = "N/A"
	}

	cycle, err := readInt(path, "cycle_count")
	if err != nil {
		log.Printf("cycle_count unavailable: %v", err)
		cycle = 0
	}

	voltNow, err := readFloat(path, "voltage_now")
	if err != nil {
		log.Printf("voltage_now unavailable: %v", err)
		voltNow = 0.0
	}

	// energy_full_design peut être absent (certains systèmes utilisent charge_full_design)
	energyFull, err := readFloat(path, "energy_full_design")
	if err != nil {
		log.Printf("energy_full_design unavailable: %v", err)
		energyFull = 0.0
	}

	// energyAH dépend de energyFull et voltNow, donc peut aussi échouer
	energyAH, err := convertWattToAmpere(path, energyFull)
	if err != nil {
		log.Printf("energy amperes calculation failed: %v", err)
		energyAH = 0.0
	}

	return BatteryInfo{
		Status:       status,
		Manufacturer: manufacturer,
		Model:        model,
		Serial:       serialNumber,
		Technology:   technology,
		Capacity:     capacity,
		Cycle:        cycle,
		EnergyAH:     energyAH,
		VoltageNow:   voltNow,
		EnergyFull:   energyFull,
	}, nil
}
