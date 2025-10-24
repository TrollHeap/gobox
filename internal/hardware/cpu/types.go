package cpu

// CPUInfo contient les informations du CPU
type CPUInfo struct {
	Architect  string  // Architecture (ex: "amd64", "arm64")
	VendorID   string  // Fabricant (ex: "GenuineIntel")
	ModelName  string  // Nom complet du CPU
	NumberCore int     // Nombre de cores physiques
	CacheSize  int64   // Taille totale cache en octets
	FreqMaxMHz float64 // Fréquence max en MHz
	FreqMinMHz float64 // Fréquence min en MHz
}

// CacheID identifie un cache physique unique
type CacheID struct {
	Level      int
	Type       string
	SharedCPUs string
}
