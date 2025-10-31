package report

import (
	"time"

	"gobox/internal/hardware/battery"
	"gobox/internal/hardware/cpu"
	"gobox/internal/hardware/disk"
	"gobox/internal/hardware/gpu"
	"gobox/internal/hardware/network"
	"gobox/internal/hardware/ram"
)

type ReconditioningReport struct {
	Metadata        ReportMetadata    `json:"metadata"`        // Métadonnées (qui, quand, où)
	Hardware        HardwareInventory `json:"hardware"`        // Informations matériel (détection)
	Tests           TestResults       `json:"tests"`           // Résultats des tests (fonctionnels)
	Grade           Grade             `json:"grade"`           // Évaluation globale
	Status          string            `json:"status"`          // "in_progress", "completed", "failed"
	Recommendations []string          `json:"recommendations"` // Recommandations d'action
}

type ReportMetadata struct {
	DeviceID     string    // ID unique de la machine
	Timestamp    time.Time // Quand le test a été fait
	GoBoxVersion string    // Version de l'outil
	Operator     string    // Qui a lancé le test (optionnel)
	Location     string    // Où (atelier A, B, etc.)
}

type HardwareInventory struct {
	CPU     cpu.CPUInfo
	RAM     ram.MemoryInfo
	Disks   []disk.DiskInfo // Plusieurs disques possibles
	GPU     *gpu.GPUInfo    // Optionnel (serveurs n'ont pas de GPU)
	Network []network.NetworkInterface
	Battery *battery.BatteryInfo // Optionnel (desktops n'ont pas de batterie)
}

type TestResults struct {
	// DiskHealth    *DiskHealthTest
	// MemoryStress  *MemoryStressTest
	// BatteryHealth *BatteryHealthTest
	// NetworkPing   *NetworkPingTest
}

type Grade struct {
	Overall   string    // "A", "B", "C", "D", "F"
	CPU       string    // Note CPU
	RAM       string    // Note RAM
	Disk      string    // Note Disque
	Battery   string    // Note Batterie
	Breakdown Breakdown // Détail du calcul
}

type Breakdown struct {
	CPUScore     int // 0-100
	RAMScore     int // 0-100
	DiskScore    int // 0-100
	BatteryScore int // 0-100
}
