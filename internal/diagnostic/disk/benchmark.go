package disk

import "time"

// DiskHealthTest résultat du test disque (diagnostic)
type DiskHealthTest struct {
	DiskName      string   // ex: "nvme0n1"
	Vendor        string   // ex: "Samsung"
	Model         string   // ex: "980 PRO"
	Type          string   // "SSD" ou "HDD"
	Grade         string   // "A", "B", "C", "F"
	SizeGB        float64  // Taille en GB
	HasPartitions bool     // true si partitions détectées
	PartitionList []string // Liste des partitions
	Issues        []string // Problèmes détectés
	Timestamp     time.Time
}

// DiskGradingCriteria critères de notation
type DiskGradingCriteria struct {
	// Pas de health% pour l'instant (pas de SMART encore)
	// On va grader sur : type, présence partitions, taille

	// Taille minimale acceptable (en GB)
	MinSizeForA float64 // ex: 200 GB
	MinSizeForB float64 // ex: 120 GB
	MinSizeForC float64 // ex: 60 GB

	// Bonus : SSD = meilleure note que HDD
	SSDBonus bool // true = SSD automatiquement >= B
}
