package ram

import (
	"os/exec"
)

type MemoryInfo struct {
	Slots   []MemorySlot
	TotalMB int
}

type MemorySlot struct {
	Slot              string  // ex: Controller0-ChannelA-DIMM0
	BankLocator       string  // ex: BANK 0
	FormFactor        string  // ex: SODIMM, DIMM
	Type              string  // ex: DDR5
	TypeDetail        string  // ex: Synchronous
	Manufacturer      string  // ex: Samsung
	SerialNumber      string  // ex: 4735802A
	PartNumber        string  // ex: M425R2GA3BB0-CWMOD
	AssetTag          string  // ex: 9876543210
	SizeMB            int     // en MB
	Speed             int     // ex: 5600
	ConfiguredSpeed   int     // ex: 4800
	ConfiguredVoltage float64 // ex: 1.1
	Rank              int     // ex: 1
	Technology        string  // ex: DRAM
	OperatingMode     string  // ex: Volatile memory
}

func GetMemoryInfo() (MemoryInfo, error) {
	cmd := exec.Command("dmidecode", "-t", "memory")
	out, err := cmd.Output()
	if err != nil {
		return MemoryInfo{}, err
	}
	slots := parseDMIDecodeOutput(out)

	total := 0
	for _, s := range slots {
		total += s.SizeMB
	}

	return MemoryInfo{
		Slots:   slots,
		TotalMB: total,
	}, nil
}
