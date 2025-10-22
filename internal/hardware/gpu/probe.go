package gpu

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	VendorNvidia Vendor = "nvidia"
	VendorIntel  Vendor = "intel"
	VendorAMD    Vendor = "amd"
)

var driverVersionPaths = map[Vendor]string{
	VendorNvidia: "/proc/driver/nvidia/version",
	VendorIntel:  "/sys/module/i915/version",
	VendorAMD:    "/sys/module/amdgpu/version",
}

var cardRegex = regexp.MustCompile(`^card\d+$`)

// NewGPUReader crée un lecteur GPU pour la carte spécifiée.
func NewGPUReader(cardName string) *GPUReader {
	return &GPUReader{
		BasePath: "/sys/class/drm",
		CardName: cardName,
	}
}

// readDriverVersion retourne la version du driver pour le vendor spécifié.
// Supporte nvidia, intel et amd (insensible à la casse).
func readDriverVersion(vendor string) (string, error) {
	v := Vendor(strings.ToLower(vendor))

	path, ok := driverVersionPaths[v]
	if !ok {
		return "", fmt.Errorf("vendor non supporté: %q (attendus: nvidia, intel, amd)", vendor)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("lecture version driver %s depuis %q: %w", v, path, err)
	}

	text := strings.TrimSpace(string(data))

	const (
		nvidiaVersionKeyword = "for"
		nvidiaVersionOffset  = 2
	)

	if v == VendorNvidia {
		fields := strings.Fields(text)
		for i, f := range fields {
			if f == nvidiaVersionKeyword && len(fields) > i+nvidiaVersionOffset {
				return fields[i+nvidiaVersionOffset], nil
			}
		}
		return "", fmt.Errorf("format inattendu du fichier version nvidia")
	}

	return text, nil
}

// readUevent lit et parse le fichier uevent d'une carte GPU.
// Extrait le driver, l'ID vendor PCI et l'adresse PCI slot.
func readUevent(cardPath string) (UeventInfo, error) {
	ueventPath := filepath.Join(cardPath, "device", "uevent")
	data, err := os.ReadFile(ueventPath)
	if err != nil {
		return UeventInfo{}, fmt.Errorf("lecture uevent %q: %w", ueventPath, err)
	}

	var info UeventInfo

	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "DRIVER="):
			info.Driver = strings.TrimPrefix(line, "DRIVER=")

		case strings.HasPrefix(line, "PCI_ID="):
			pciID := strings.TrimPrefix(line, "PCI_ID=")
			if !strings.Contains(pciID, ":") {
				return UeventInfo{}, fmt.Errorf(
					"format PCI_ID invalide dans %q: %q",
					ueventPath,
					pciID,
				)
			}

			parts := strings.SplitN(pciID, ":", 2)
			if len(parts) == 2 {
				if len(parts[0]) != 4 || !isHexString(parts[0]) {
					return UeventInfo{}, fmt.Errorf("vendor ID invalide: %q", parts[0])
				}
				info.VendorID = "0x" + strings.ToLower(parts[0])
			}

		case strings.HasPrefix(line, "PCI_SLOT_NAME="):
			pciSlot := strings.TrimPrefix(line, "PCI_SLOT_NAME=")
			if !isValidPCISlot(pciSlot) {
				return UeventInfo{}, fmt.Errorf(
					"PCI_SLOT_NAME invalide dans %q: %q",
					ueventPath,
					pciSlot,
				)
			}
			info.PCISlot = pciSlot
		}
	}

	return info, nil
}

// readLspci interroge lspci pour le périphérique au slot PCI spécifié.
// Retourne une erreur si le slot est invalide ou non trouvé.
func readLspci(pciSlot string) (PCIDevice, error) {
	if !isValidPCISlot(pciSlot) {
		return PCIDevice{}, fmt.Errorf(
			"format PCI slot invalide: %q (attendu: XXXX:XX:XX.X)",
			pciSlot,
		)
	}

	out, err := exec.Command("lspci", "-mm", "-nn", "-D", "-s", pciSlot).Output()
	if err != nil {
		return PCIDevice{}, fmt.Errorf("échec lspci pour slot %q: %w", pciSlot, err)
	}

	output := strings.TrimSpace(string(out))
	if output == "" {
		return PCIDevice{}, fmt.Errorf("aucun périphérique trouvé pour slot %q", pciSlot)
	}

	parts := strings.Split(output, "\"")
	if len(parts) < 6 {
		return PCIDevice{}, fmt.Errorf("format lspci inattendu pour slot %q: %q", pciSlot, output)
	}

	return PCIDevice{
		Vendor: parts[3],
		Model:  parts[5],
	}, nil
}

// listCards retourne les chemins de toutes les cartes DRM du système.
func listCards() ([]string, error) {
	const drmPath = "/sys/class/drm/"
	entries, err := os.ReadDir(drmPath)
	if err != nil {
		return nil, fmt.Errorf("lecture répertoire DRM %q: %w", drmPath, err)
	}

	cards := make([]string, 0, len(entries)/2)
	for _, entry := range entries {
		if cardRegex.MatchString(entry.Name()) {
			cards = append(cards, filepath.Join(drmPath, entry.Name()))
		}
	}

	return cards, nil
}

// readConnectorStatus lit le statut d'un connecteur DRM.
func readConnectorStatus(path string) (string, error) {
	fullPath := filepath.Join(path, "status")
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("lecture du statut du connecteur %q: %w", path, err)
	}
	return strings.TrimSpace(string(data)), nil
}

// listConnectors retourne les connecteurs connectés d'une carte GPU.
func listConnectors(cardPath string) ([]string, error) {
	const drmRoot = "/sys/class/drm"

	cleanPath := filepath.Clean(cardPath)
	if !strings.HasPrefix(cleanPath, drmRoot) {
		return nil, fmt.Errorf("cardPath invalide (hors de %s): %q", drmRoot, cardPath)
	}

	cardName := filepath.Base(cleanPath)

	entries, err := os.ReadDir(drmRoot)
	if err != nil {
		return nil, fmt.Errorf("lecture drm root: %w", err)
	}

	connectedOutputs := make([]string, 0, len(entries))

	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), cardName+"-") {
			continue
		}

		connectorPath := filepath.Join(drmRoot, entry.Name())

		status, err := readConnectorStatus(connectorPath)
		if err != nil {
			return nil, fmt.Errorf("vérification connecteur %s: %w", entry.Name(), err)
		}

		if status == "connected" {
			name := strings.TrimPrefix(entry.Name(), cardName+"-")
			connectedOutputs = append(connectedOutputs, name)
		}
	}

	return connectedOutputs, nil
}
