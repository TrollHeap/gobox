package batterytest

import "time"

type BatteryHealthTest struct {
	Status           string
	Grade            string
	HealthPercentage float64
	CycleCount       int
	Issues           []string
	Timestamp        time.Time
}

type BatteryGradingCriteria struct {
	MinHealthForA float64
	MinHealthForB float64
	MinHealthForC float64
	MaxCyclesForA int
	MaxCyclesForB int
	MaxCyclesForC int
}
