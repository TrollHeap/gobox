package probe

import (
	"os/exec"
	"strconv"
	"strings"
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

func cutPrefixString(line, prefix string, field *string) {
	if after, ok := strings.CutPrefix(line, prefix); ok {
		*field = strings.TrimSpace(after)
	}
}

func parseIntField(line, prefix string, field *int, mult int, suffixes ...string) bool {
	if after, ok := strings.CutPrefix(line, prefix); ok {
		val := strings.TrimSpace(after)
		for _, s := range suffixes {
			val = strings.TrimSuffix(val, s)
		}
		val = strings.TrimSpace(val)

		if val == "" || val == "Unknown" || val == "0" || strings.Contains(val, "No Module") {
			return false
		}

		num, err := strconv.Atoi(val)
		if err == nil {
			*field = num * mult
			return true
		}
	}
	return false
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

func parseDMIDecodeOutput(data []byte) []MemorySlot {
	seen := map[string]bool{}
	var slots []MemorySlot
	var slot MemorySlot

	for line := range strings.SplitSeq(string(data), "\n") {

		line = strings.TrimSpace(line)

		if line == "Memory Device" {
			slot = MemorySlot{}
			continue
		}

		cutPrefixString(line, "Locator:", &slot.Slot)
		cutPrefixString(line, "Type:", &slot.Type)
		cutPrefixString(line, "Manufacturer:", &slot.Manufacturer)
		cutPrefixString(line, "Serial Number:", &slot.SerialNumber)
		cutPrefixString(line, "Bank Locator:", &slot.BankLocator)
		cutPrefixString(line, "Form Factor:", &slot.FormFactor)
		cutPrefixString(line, "Type Detail:", &slot.TypeDetail)
		cutPrefixString(line, "Serial Number:", &slot.SerialNumber)
		cutPrefixString(line, "Part Number:", &slot.PartNumber)
		cutPrefixString(line, "Asset Tag:", &slot.AssetTag)
		cutPrefixString(line, "Memory Technology:", &slot.Technology)
		cutPrefixString(line, "Memory Operating Mode Capability:", &slot.OperatingMode)

		parseIntField(line, "Size:", &slot.SizeMB, 1024, "GB", "MB")
		parseIntField(line, "Speed:", &slot.Speed, 1, "MT/s", "MHz")
		parseIntField(line, "Configured Memory Speed:", &slot.ConfiguredSpeed, 1, "MT/s", "MHz")

		if slot.SizeMB > 0 && line == "" && !seen[slot.Slot] {
			slots = append(slots, slot)
			seen[slot.Slot] = true
		}
	}

	return slots
}
