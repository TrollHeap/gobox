package battery

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
