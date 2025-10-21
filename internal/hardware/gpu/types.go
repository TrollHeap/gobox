package gpu

type UeventInfo struct {
	Driver   string
	VendorID string
	DeviceID string
	PCISlot  string
}

type GPUInfo struct {
	Model       string
	Vendor      string
	Driver      string
	Version     string
	VendorID    string
	DeviceID    string
	Outputs     []string
	DisplayInfo string
	Message     string
}
