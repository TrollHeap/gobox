package ui

import (
	"fmt"
	"strings"

	"gobox/internal/hardware/network"
)

// PrintNetworkInterfaces displays all detected network interfaces
func PrintNetworkInterfaces() {
	interfaces, err := network.ListNetworkInterfaces()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(interfaces) == 0 {
		fmt.Println("\nNo network interfaces detected.")
		return
	}

	fmt.Printf("\n" + strings.Repeat("=", 70) + "\n")
	fmt.Printf("  NETWORK INTERFACES DETECTED: %d\n", len(interfaces))
	fmt.Printf(strings.Repeat("=", 70) + "\n\n")

	for i, iface := range interfaces {
		fmt.Printf("Interface #%d\n", i+1)
		fmt.Println(strings.Repeat("-", 70))
		fmt.Printf("  Name            : %s\n", iface.Name)

		if iface.MACAddress != "" {
			fmt.Printf("  MAC Address     : %s\n", iface.MACAddress)
		}

		if iface.Type != "" {
			fmt.Printf("  Type            : %s\n", iface.Type)
		}

		fmt.Printf("  State           : ")
		if iface.IsUp {
			fmt.Printf("UP ✅\n")
		} else {
			fmt.Printf("DOWN ❌\n")
		}

		if iface.Speed != "" {
			fmt.Printf("  Speed           : %s\n", iface.Speed)
		}

		fmt.Printf("  Carrier         : ")
		if iface.Carrier {
			fmt.Printf("Connected 🔌\n")
		} else {
			fmt.Printf("Disconnected\n")
		}

		if iface.IPAddress != "" {
			fmt.Printf("  IP              : %s\n", iface.IPAddress)
		}

		fmt.Println(strings.Repeat("-", 70))

		if i < len(interfaces)-1 {
			fmt.Println()
		}
	}
	fmt.Println()
}
