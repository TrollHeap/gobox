package battery

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"gobox/internal/sysfs"
	"gobox/internal/utils"
)

type BatteryInfo struct {
	Status           string
	Manufacturer     *string
	Model            *string
	Serial           *string
	Technology       *string
	Capacity         int
	Cycle            int
	EnergyAH         float64
	VoltageNow       float64
	EnergyFull       float64 // ← Capacité actuelle (réelle)
	EnergyFullDesign float64 // 🆕 Capacité théorique (neuve)
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

	// Required fields
	capacity, err := sysfs.ReadInt(filepath.Join(path, "capacity"))
	if err != nil {
		return BatteryInfo{}, fmt.Errorf("reading capacity: %w", err)
	}

	status, err := sysfs.ReadFile(filepath.Join(path, "status"))
	if err != nil {
		return BatteryInfo{}, fmt.Errorf("reading status: %w", err)
	}

	// Optional fields (no error handling, use empty string if missing)
	manufacturer, _ := sysfs.ReadFileOptional(filepath.Join(path, "manufacturer"))
	model, _ := sysfs.ReadFileOptional(filepath.Join(path, "model_name"))
	serialNumber, _ := sysfs.ReadFileOptional(filepath.Join(path, "serial_number"))
	technology, _ := sysfs.ReadFileOptional(filepath.Join(path, "technology"))

	cycle := 0
	if c, err := sysfs.ReadInt(filepath.Join(path, "cycle_count")); err == nil {
		cycle = c
	}

	voltNow := 0.0
	if v, err := sysfs.ReadFloat(filepath.Join(path, "voltage_now")); err == nil {
		voltNow = v
	}

	energyFull := 0.0
	if e, err := sysfs.ReadFloat(filepath.Join(path, "energy_full")); err == nil {
		energyFull = e
	}

	// Capacité théorique (neuve)
	energyFullDesign := 0.0
	if e, err := sysfs.ReadFloat(filepath.Join(path, "energy_full_design")); err == nil {
		energyFullDesign = e
	}

	// Convert Wh to Ah (silent error handling)
	energyAH, _ := utils.ConvertWattToAmpere(path, energyFull)

	return BatteryInfo{
		Status:           status,
		Manufacturer:     strPtrIfNotEmpty(manufacturer),
		Model:            strPtrIfNotEmpty(model),
		Serial:           strPtrIfNotEmpty(serialNumber),
		Technology:       strPtrIfNotEmpty(technology),
		Capacity:         capacity,
		Cycle:            cycle,
		EnergyAH:         energyAH,
		VoltageNow:       voltNow,
		EnergyFull:       energyFull,
		EnergyFullDesign: energyFullDesign,
	}, nil
}
