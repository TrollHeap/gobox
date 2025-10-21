package battery

import (
	"errors"
	"os"
)

func GetBatteryInfo() (BatteryInfo, error) {
	path, err := findBatteryPath()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return BatteryInfo{Capacity: -1, Status: "N/A"}, nil
		}
		return BatteryInfo{}, err
	}

	capacity, _ := readInt(path, "capacity")
	cycle, _ := readInt(path, "cycle_count")
	status, _ := readString(path, "status")
	manufacturer, _ := readString(path, "manufacturer")
	model, _ := readString(path, "model_name")
	serial, _ := readString(path, "serial_number")
	technology, _ := readString(path, "technology")
	volt, _ := readFloat(path, "voltage_now")
	energyFull, _ := readFloat(path, "energy_full_design")
	energyAH := convertWattToAmpere(path, energyFull)

	return BatteryInfo{
		Status:       status,
		Manufacturer: manufacturer,
		Model:        model,
		Serial:       serial,
		Technology:   technology,
		Capacity:     capacity,
		Cycle:        cycle,
		EnergyAH:     energyAH,
		VoltageNow:   volt,
		EnergyFull:   energyFull,
	}, nil
}
