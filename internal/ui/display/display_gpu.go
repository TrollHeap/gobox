package ui

import (
	"fmt"

	"gobox/internal/hardware/gpu"
)

func DisplayGPUInfo() {
	gpus, err := gpu.DetectGPUs()
	if err != nil {
		fmt.Printf("Erreur de détection GPU : %v\n", err)
		return
	}

	if len(gpus) == 0 {
		fmt.Println("Aucune carte graphique détectée.")
		return
	}

	for _, g := range gpus {
		fmt.Println("─────────────────────────────────────────────")
		fmt.Printf("Modèle          : %s\n", g.Model)
		fmt.Printf("Marque          : %s (%s)\n", g.Vendor, g.VendorID)
		fmt.Printf("Driver          : %s\n", g.Driver)
		fmt.Printf("Version         : %s\n", g.Version)
		fmt.Printf("Sorties actives : %v\n", g.Outputs)
	}
	fmt.Println("─────────────────────────────────────────────")
}
