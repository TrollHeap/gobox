package gpu

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func readDisplayModes(outputs []string) []string {
	modes := []string{}
	for _, out := range outputs {
		path := filepath.Join("/sys/class/drm", "card1-"+out, "modes")
		data, err := os.ReadFile(path)
		if err == nil {
			modes = append(modes, strings.Split(strings.TrimSpace(string(data)), "\n")...)
		}
	}
	return modes
}

func readDriverVersion(vendor string) string {
	var path string

	switch strings.ToLower(vendor) {
	case "nvidia":
		path = "/proc/driver/nvidia/version"
	case "intel":
		path = "/sys/module/i915/version"
	case "amd":
		path = "/sys/module/amdgpu/version"
	default:
		return "Unknown"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "Unknown"
	}

	text := strings.TrimSpace(string(data))

	// Extraction selon constructeur :
	if vendor == "nvidia" {
		fields := strings.Fields(text)
		for i, f := range fields {
			if f == "for" && len(fields) > i+2 {
				return fields[i+2] // saute "x86_64"
			}
		}
	}
	return text
}

func readUevent(cardPath string) (UeventInfo, error) {
	ueventPath := filepath.Join(cardPath, "device", "uevent")
	data, err := os.ReadFile(ueventPath)
	if err != nil {
		return UeventInfo{}, fmt.Errorf("lecture du uevent: %w", err)
	}

	var driver, vendorID, deviceID, pciSlot string
	for line := range strings.SplitSeq(string(data), "\n") {
		switch {
		case strings.HasPrefix(line, "DRIVER="):
			driver = strings.TrimPrefix(line, "DRIVER=")
		case strings.HasPrefix(line, "PCI_ID="):
			parts := strings.Split(strings.TrimPrefix(line, "PCI_ID="), ":")
			if len(parts) == 2 {
				vendorID = "0x" + strings.ToLower(parts[0])
				deviceID = "0x" + strings.ToLower(parts[1])
			}
		case strings.HasPrefix(line, "PCI_SLOT_NAME="):
			pciSlot = strings.TrimPrefix(line, "PCI_SLOT_NAME=")
		}
	}

	return UeventInfo{
		Driver:   driver,
		VendorID: vendorID,
		DeviceID: deviceID,
		PCISlot:  pciSlot,
	}, nil
}

func readConnectorStatus(path string) string {
	data, _ := os.ReadFile(filepath.Join(path, "status"))
	return strings.TrimSpace(string(data))
}

func readLspci(pciSlot string) (string, string) {
	out, err := exec.Command("lspci", "-mm", "-nn", "-D").Output()
	if err != nil {
		return "Unknown", "Unknown"
	}
	for line := range strings.SplitSeq(string(out), "\n") {
		if strings.Contains(line, pciSlot) {
			parts := strings.Split(line, "\"")
			if len(parts) >= 6 {
				vendor := parts[3] // ex: NVIDIA Corporation
				model := parts[5]  // ex: AD106M [GeForce RTX 4070 Max-Q / Mobile]
				return vendor, model
			}
		}
	}
	return "Unknown", "Unknown"
}

func listCards() ([]string, error) {
	drmPath := "/sys/class/drm/"
	cardRegex := regexp.MustCompile("^card[0-9]+$")

	entries, err := os.ReadDir(drmPath)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire %s: %w", drmPath, err)
	}

	cards := []string{}
	for _, entry := range entries {
		if cardRegex.MatchString(entry.Name()) {
			path := filepath.Join(drmPath, entry.Name())
			cards = append(cards, path)
		}
	}

	return cards, err
}

func listConnectors(cardPath string) ([]string, error) {
	drmRoot := "/sys/class/drm"
	cardName := filepath.Base(cardPath) // ex: "card1"

	entries, err := os.ReadDir(drmRoot)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire %s: %w", drmRoot, err)
	}

	connectedOutputs := []string{}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), cardName+"-") {
			connectorPath := filepath.Join(drmRoot, entry.Name())
			if readConnectorStatus(connectorPath) == "connected" {
				name := strings.TrimPrefix(entry.Name(), cardName+"-")
				connectedOutputs = append(connectedOutputs, name)
			}
		}
	}

	return connectedOutputs, nil
}
