package battery

import (
	"fmt"
	"log"

	"gobox/internal/sysfs"
	"gobox/internal/utils"
)

type BatteryInfo struct {
	Status       string
	Manufacturer string
	Model        string
	Serial       string
	Technology   string
	Capacity     int
	Cycle        int
	EnergyAH     float64
	VoltageNow   float64
	EnergyFull   float64
}

func GetBatteryInfo() (BatteryInfo, error) {
	path, err := findBatteryPath()
	if err != nil {
		return BatteryInfo{}, err
	}

	// Champs critiques (obligatoires selon la spec Linux)
	capacity, err := sysfs.ReadInt(path, "capacity")
	if err != nil {
		return BatteryInfo{}, fmt.Errorf("reading capacity: %w", err)
	}

	status, err := sysfs.ReadString(path, "status")
	if err != nil {
		return BatteryInfo{}, fmt.Errorf("reading status: %w", err)
	}

	// Champs optionnels (peuvent être absents selon le matériel/driver)
	manufacturer, err := sysfs.ReadString(path, "manufacturer")
	if err != nil {
		log.Printf("manufacturer unavailable: %v", err)
		manufacturer = "N/A"
	}

	model, err := sysfs.ReadString(path, "model_name")
	if err != nil {
		log.Printf("model_name unavailable: %v", err)
		model = "N/A"
	}

	serialNumber, err := sysfs.ReadString(path, "serial_number")
	if err != nil {
		log.Printf("serial_number unavailable: %v", err)
		serialNumber = "N/A"
	}

	technology, err := sysfs.ReadString(path, "technology")
	if err != nil {
		log.Printf("technology unavailable: %v", err)
		technology = "N/A"
	}

	cycle, err := sysfs.ReadInt(path, "cycle_count")
	if err != nil {
		log.Printf("cycle_count unavailable: %v", err)
		cycle = 0
	}

	voltNow, err := sysfs.ReadFloat(path, "voltage_now")
	if err != nil {
		log.Printf("voltage_now unavailable: %v", err)
		voltNow = 0.0
	}

	// energy_full_design peut être absent (certains systèmes utilisent charge_full_design)
	energyFull, err := sysfs.ReadFloat(path, "energy_full_design")
	if err != nil {
		log.Printf("energy_full_design unavailable: %v", err)
		energyFull = 0.0
	}

	// energyAH dépend de energyFull et voltNow, donc peut aussi échouer
	energyAH, err := utils.ConvertWattToAmpere(path, energyFull)
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
