package batterytest

import (
	"fmt"
)

const (
	StatusCharging    = "Charging"
	StatusDischarging = "Discharging"
	StatusFull        = "Full"
	StatusNotCharging = "Not charging"
)

func DefaultBatteryGradingCriteria() BatteryGradingCriteria {
	return BatteryGradingCriteria{
		MinHealthForA: 85.0, // 85% de santé = A
		MinHealthForB: 70.0, // 70% de santé = B
		MinHealthForC: 50.0, // 50% de santé = C
		MaxCyclesForA: 300,  // <= 300 cycles = A
		MaxCyclesForB: 500,  // <= 500 cycles = B
		MaxCyclesForC: 800,  // <= 800 cycles = C
	}
}

func CalculateHealthPercent(designCapacity, currentCapacity float64) float64 {
	if designCapacity <= 0 {
		return 0.0
	}
	healthPercentage := (currentCapacity / designCapacity) * 100
	return healthPercentage
}

func ComputeGrade(criteria BatteryGradingCriteria, healthPercent float64, cycleCount int) string {
	healthGrade := gradeFromHealth(criteria, healthPercent)
	cycleGrade := gradeFromCycle(criteria, cycleCount)

	if gradeToScore(healthGrade) < gradeToScore(cycleGrade) {
		return healthGrade
	}
	return cycleGrade
}

func DetectIssues(
	status string,
	healthPercentage float64,
	cycleCount int,
	criteria BatteryGradingCriteria,
) []string {
	issues := []string{}

	validStatuses := map[string]bool{
		StatusCharging:    true,
		StatusDischarging: true,
		StatusFull:        true,
		StatusNotCharging: true,
	}

	if !validStatuses[status] {
		issues = append(issues, fmt.Sprintf("Statut invalide: %s", status))
		return issues
	}

	if status == "Discharging" && healthPercentage < 20 {
		issues = append(issues, "Batterie en décharge avec une santé très basse. (<20%)")
	}
	if healthPercentage < criteria.MinHealthForC {
		issues = append(issues, "La santé de la batterie est trop basse. (<50%)")
	}
	if healthPercentage > 100 {
		issues = append(issues, "Anormal santé de la batterie (>100%), vérifier les capteurs.")
	}

	if cycleCount > criteria.MaxCyclesForC {
		issues = append(
			issues,
			"Haut cycle de charge de la batterie (>800), considérer un remplacement.",
		)
	}

	return issues
}

func gradeFromHealth(criteria BatteryGradingCriteria, healthPercent float64) string {
	switch {
	case healthPercent >= criteria.MinHealthForA:
		return "A"
	case healthPercent >= criteria.MinHealthForB:
		return "B"
	case healthPercent >= criteria.MinHealthForC:
		return "C"
	default:
		return "F"
	}
}

func gradeFromCycle(criteria BatteryGradingCriteria, cycleCount int) string {
	switch {
	case cycleCount <= criteria.MaxCyclesForA:
		return "A"
	case cycleCount <= criteria.MaxCyclesForB:
		return "B"
	case cycleCount <= criteria.MaxCyclesForC:
		return "C"
	default:
		return "F"
	}
}

func gradeToScore(grade string) int {
	switch grade {
	case "A":
		return 4
	case "B":
		return 3
	case "C":
		return 2
	case "F":
		return 1
	default:
		return 0
	}
}
