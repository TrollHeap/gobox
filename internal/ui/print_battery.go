package ui

import (
	"fmt"

	"gobox/internal/hardware/battery"
)

func PrintBattery() {
	info, err := battery.GetBatteryInfo()
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	if info.Capacity < 0 {
		fmt.Println("🔋 Batterie        : Pas de batterie détectée")
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
