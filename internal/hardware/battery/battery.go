package battery

import (
	"errors"
	"fmt"
	"sync"

	"gobox/internal/sysfs"
	"gobox/internal/utils"
)

type BatteryInfo struct {
	Status       string
	Manufacturer *string
	Model        *string
	Serial       *string
	Technology   *string
	Capacity     int
	Cycle        int
	EnergyAH     float64
	VoltageNow   float64
	EnergyFull   float64
}

var (
	batteryPath string
	once        sync.Once
)

func findBatteryPathCached() (string, error) {
	var err error
	once.Do(func() {
		batteryPath, err = findBatteryPath()
	})
	if batteryPath == "" || err != nil {
		return "", errors.New("battery path not found")
	}
	return batteryPath, nil
}

func strPtrIfNotEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func GetBatteryInfo() (BatteryInfo, error) {
	path, err := findBatteryPathCached()
	if err != nil {
		return BatteryInfo{}, err
	}

	capacity, err := sysfs.ReadInt(path, "capacity")
	if err != nil {
		return BatteryInfo{}, fmt.Errorf("reading capacity: %w", err)
	}

	status, err := sysfs.ReadString(path, "status")
	if err != nil {
		return BatteryInfo{}, fmt.Errorf("reading status: %w", err)
	}

	// Lecture des champs optionnels sans logger, uniquement gestion d'absence avec pointeur nil
	manufacturer, _ := sysfs.ReadString(path, "manufacturer")
	model, _ := sysfs.ReadString(path, "model_name")
	serialNumber, _ := sysfs.ReadString(path, "serial_number")
	technology, _ := sysfs.ReadString(path, "technology")

	cycle, _ := sysfs.ReadInt(path, "cycle_count")

	voltNow, _ := sysfs.ReadFloat(path, "voltage_now")

	energyFull, _ := sysfs.ReadFloat(path, "energy_full_design")

	// Conversion avec gestion d'erreur silencieuse
	energyAH, _ := utils.ConvertWattToAmpere(path, energyFull)

	return BatteryInfo{
		Status:       status,
		Manufacturer: strPtrIfNotEmpty(manufacturer),
		Model:        strPtrIfNotEmpty(model),
		Serial:       strPtrIfNotEmpty(serialNumber),
		Technology:   strPtrIfNotEmpty(technology),
		Capacity:     capacity,
		Cycle:        cycle,
		EnergyAH:     energyAH,
		VoltageNow:   voltNow,
		EnergyFull:   energyFull,
	}, nil
}
