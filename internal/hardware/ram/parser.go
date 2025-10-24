package ram

import (
	"os/exec"
	"strings"
)

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
