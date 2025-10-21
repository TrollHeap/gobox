package gpu

import (
	"strings"
)

func DetectGPUs() ([]GPUInfo, error) {
	cards, err := listCards()
	if err != nil {
		return nil, err
	}

	infos := []GPUInfo{}
	for _, card := range cards {
		u, _ := readUevent(card)
		vendorName, model := readLspci(u.PCISlot)
		driverVersion := readDriverVersion(u.Driver)
		outputs, _ := listConnectors(card)
		displayModes := readDisplayModes(outputs)

		info := GPUInfo{
			Model:       model,
			Vendor:      vendorName,
			Driver:      u.Driver,
			Version:     driverVersion,
			VendorID:    u.VendorID,
			DeviceID:    u.DeviceID,
			Outputs:     outputs,
			DisplayInfo: strings.Join(displayModes, ", "),
			Message:     "GL data unavailable for root.",
		}
		infos = append(infos, info)
	}

	return infos, nil
}
