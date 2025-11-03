package disktest

import "time"

// DiskHealthTest résultat du test disque
type DiskHealthTest struct {
	DiskName      string   // ex: "nvme0n1"
	Type          string   // "SSD" ou "HDD"
	Grade         string   // "A", "B", "C", "F"
	HealthPercent float64  // Santé globale (ex: 98.5%)
	UsagePercent  float64  // % d'espace utilisé
	SmartStatus   string   // "PASSED", "FAILED", "UNKNOWN"
	Issues        []string // Problèmes détectés
	Timestamp     time.Time
}

// DiskGradingCriteria critères de notation
type DiskGradingCriteria struct {
	// Santé
	MinHealthForA float64 // ex: 95.0
	MinHealthForB float64 // ex: 85.0
	MinHealthForC float64 // ex: 70.0

	// Utilisation espace
	MaxUsageForA float64 // ex: 70.0 (%)
	MaxUsageForB float64 // ex: 85.0
	MaxUsageForC float64 // ex: 95.0
}
