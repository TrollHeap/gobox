package gpu

import (
	"fmt"
	"strings"
)

// DetectGPUs retourne les informations de toutes les cartes GPU du système.
// Ignore les cartes individuelles en erreur et retourne les résultats partiels.
func DetectGPUs() ([]GPUInfo, error) {
	cards, err := listCards()
	if err != nil {
		return nil, fmt.Errorf("impossible de lister les cartes GPU: %w", err)
	}

	if len(cards) == 0 {
		return []GPUInfo{}, nil
	}

	infos := make([]GPUInfo, 0, len(cards))
	var detectionErrors []string

	for _, card := range cards {
		info, err := detectSingleGPU(card)
		if err != nil {
			detectionErrors = append(detectionErrors,
				fmt.Sprintf("carte %s: %v", card, err))
			continue
		}
		infos = append(infos, info)
	}

	if len(infos) == 0 && len(detectionErrors) > 0 {
		return nil, fmt.Errorf("échec détection de toutes les cartes:\n%s",
			strings.Join(detectionErrors, "\n"))
	}

	return infos, nil
}

// detectSingleGPU collecte les informations d'une carte GPU.
// Utilise des valeurs par défaut pour les données optionnelles en erreur.
func detectSingleGPU(cardPath string) (GPUInfo, error) {
	uevent, err := readUevent(cardPath)
	if err != nil {
		return GPUInfo{}, fmt.Errorf("lecture uevent: %w", err)
	}

	pciDevice, err := readLspci(uevent.PCISlot)
	if err != nil {
		return GPUInfo{}, fmt.Errorf("lecture lspci pour slot %q: %w", uevent.PCISlot, err)
	}

	driverVersion, err := readDriverVersion(uevent.Driver)
	if err != nil {
		driverVersion = "inconnu"
	}

	outputs, err := listConnectors(cardPath)
	if err != nil {
		outputs = []string{}
	}

	var displayModes []string
	if len(outputs) > 0 {
		reader := NewGPUReader(cardPath)
		displayModes, err = reader.ReadDisplayModes(outputs)
		if err != nil {
			displayModes = []string{}
		}
	}

	return GPUInfo{
		Model:       pciDevice.Model,
		Vendor:      pciDevice.Vendor,
		Driver:      uevent.Driver,
		Version:     driverVersion,
		VendorID:    uevent.VendorID,
		Outputs:     outputs,
		DisplayInfo: strings.Join(displayModes, ", "),
	}, nil
}
