package batterytest

import (
	"fmt"

	"gobox/internal/hardware/battery"
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

	if info.Manufacturer != nil {
		fmt.Printf("Fabricant        : %s\n", *info.Manufacturer)
	}
	if info.Model != nil {
		fmt.Printf("Modèle           : %s\n", *info.Model)
	}
	if info.Serial != nil {
		fmt.Printf("Numéro série     : %s\n", *info.Serial)
	}
	if info.Technology != nil {
		fmt.Printf("Technologie      : %s\n", *info.Technology)
	}

	fmt.Printf("Cycles           : %d\n", info.Cycle)
	fmt.Printf("Tension actuelle : %.2f V\n", info.VoltageNow/1_000_000)
	fmt.Printf("Capacité totale  : %.0f mAh\n", info.EnergyAH)
	fmt.Printf("Capacité actuelle : %.0f \n", info.DesignCapacity)
	fmt.Printf("Capacité Neuve : %.0f \n", info.DesignCapacity)
}
