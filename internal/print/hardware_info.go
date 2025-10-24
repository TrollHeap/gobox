package print

import (
	"fmt"
	"log"
	"os"

	"gobox/internal/hardware/battery"
	"gobox/internal/hardware/cpu"
	"gobox/internal/hardware/gpu"
	"gobox/internal/hardware/ram"
)

func BatteryInfo() {
	info, err := battery.GetBatteryInfo()
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	if info.Capacity < 0 {
		fmt.Println("Batterie        : Pas de batterie détectée")
		return
	}

	fmt.Printf("Batterie         : %d%% (%s)\n", info.Capacity, info.Status)
	fmt.Printf("Fabricant        : %s\n", info.Manufacturer)
	fmt.Printf("Modèle           : %s\n", info.Model)
	fmt.Printf("Numéro série     : %s\n", info.Serial)
	fmt.Printf("Technologie      : %s\n", info.Technology)
	fmt.Printf("Cycles           : %d\n", info.Cycle)
	fmt.Printf("Tension actuelle : %.2f V\n", info.VoltageNow/1_000_000)
	fmt.Printf("Capacité totale  : %.0f mAh\n", info.EnergyAH)
}

func RamInfo() {
	info, err := ram.GetMemoryInfo()
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Println("─── Rapport Mémoire ───────────────────────────────────────────────")

	for _, s := range info.Slots {
		fmt.Printf("Slot: %-30s | Type: %-6s | %6d MB @ %-5d MHz\n",
			s.Slot, s.Type, s.SizeMB, s.Speed)

		if s.ConfiguredSpeed > 0 && s.ConfiguredSpeed != s.Speed {
			fmt.Printf("  ↳ Vitesse configurée : %d MHz\n", s.ConfiguredSpeed)
		}

		fmt.Printf("  Fabricant     : %-15s  |  Forme : %-8s  |  Rank : %d\n",
			s.Manufacturer, s.FormFactor, s.Rank)

		if s.ConfiguredVoltage > 0 {
			fmt.Printf("  Tension       : %.1f V\n", s.ConfiguredVoltage)
		}

		if s.PartNumber != "" {
			fmt.Printf("  Numéro pièce  : %s\n", s.PartNumber)
		}

		if s.SerialNumber != "" {
			fmt.Printf("  Numéro série  : %s\n", s.SerialNumber)
		}

		if s.AssetTag != "" {
			fmt.Printf("  Tag inventaire: %s\n", s.AssetTag)
		}

		fmt.Printf("  Technologie   : %-10s | Mode : %s\n",
			s.Technology, s.OperatingMode)

		fmt.Println("───────────────────────────────────────────────────────────────────")
	}

	fmt.Printf("Mémoire totale détectée : %d MB (%.1f GB)\n",
		info.TotalMB, float64(info.TotalMB)/1024)
}

func GPUInfo() {
	gpus, err := gpu.DetectGPUs()
	if err != nil {
		log.Fatalf("Erreur de détection GPU : %v", err)
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

func CPUInfo() {
	info, err := cpu.GetCPUInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Vendor ID:       %s\n", info.VendorID)
	fmt.Printf("Model Name:      %s\n", info.ModelName)
	fmt.Printf("Architecture:    %s\n", info.Architect)
	fmt.Printf("CPU Cores:       %d\n", info.NumberCore)
	fmt.Printf("Total Cache:     %d octets (%.2f MB)\n",
		info.CacheSize, float64(info.CacheSize)/(1024*1024))

	if info.FreqMaxMHz > 0 {
		fmt.Printf("Freq Min:        %.0f MHz\n", info.FreqMinMHz)
		fmt.Printf("Freq Max:        %.0f MHz\n", info.FreqMaxMHz)
	}
}
