package gpu

type UeventInfo struct {
	Driver   string
	VendorID string
	DeviceID string
	PCISlot  string
}

// GPUInfo représente les informations complètes d'une carte GPU détectée.
type GPUInfo struct {
	Model       string   // Modèle du GPU (ex: "AD106M [GeForce RTX 4070 Max-Q]")
	Vendor      string   // Fabricant (ex: "NVIDIA Corporation")
	Driver      string   // Driver noyau (ex: "nvidia", "i915")
	Version     string   // Version du driver
	VendorID    string   // ID hexadécimal du vendor (ex: "0x10de")
	Outputs     []string // Connecteurs connectés (ex: ["HDMI-A-1", "eDP-1"])
	DisplayInfo string   // Modes d'affichage concaténésj
}

// PCIDevice représente un périphérique PCI extrait de lspci.
type PCIDevice struct {
	Vendor string
	Model  string
}

// GPUReader lit les informations d'affichage d'une carte GPU.
type GPUReader struct {
	BasePath string
	CardName string
}

// Vendor représente un fabricant de GPU supporté.
type Vendor string
