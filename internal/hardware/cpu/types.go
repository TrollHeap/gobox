package cpu

type CPUInfo struct {
	Architect string // Architecture du CPU (ex: "amd64", "arm64")
	VendorID  string // Fabricant
	ModelName string // Nom complet du CPU
	// CacheSize  int64  // Cache size de tous les processeurs
	NumberCore int // Nombre de coeurs
	// FrequenceAVG float64 // Fréquence AVG
}

type CacheID struct {
	Level      int    // Niveau du cache (1, 2, 3, etc.)
	SharedCPUs string // Ex: "0-23", "0-1", etc.
}
