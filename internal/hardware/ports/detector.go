package ports

import (
	"fmt"
	"strings"
)

// GetUSBInfo collecte toutes les informations USB du système
func GetUSBInfo() (*USBInfo, error) {
	info := &USBInfo{}

	// Collecter contrôleurs
	controllers, err := ListUSBControllers()
	if err != nil {
		// Non-bloquant : continue avec liste vide
		controllers = []USBController{}
	}
	info.Controllers = controllers

	// Collecter devices
	devices, err := ListUSBDevices()
	if err != nil {
		return nil, fmt.Errorf("erreur devices USB: %w", err)
	}
	info.Devices = devices

	// Collecter ports USB-C
	usbcPorts, err := ListUSBCPorts()
	if err != nil {
		return nil, fmt.Errorf("erreur ports USB-C: %w", err)
	}
	info.USBCPorts = usbcPorts

	return info, nil
}

// PrintUSBInfo affiche toutes les informations USB
func PrintUSBInfo() {
	info, err := GetUSBInfo()
	if err != nil {
		fmt.Printf("Erreur: %v\n", err)
		return
	}

	fmt.Printf("\n" + strings.Repeat("=", 70) + "\n")
	fmt.Printf("  CONNECTIVITÉ USB\n")
	fmt.Printf(strings.Repeat("=", 70) + "\n")

	// Section 1 : Contrôleurs
	fmt.Println("\n  Contrôleurs USB détectés")
	fmt.Println("  " + strings.Repeat("-", 66))

	if len(info.Controllers) == 0 {
		fmt.Println("    Aucun contrôleur détecté (erreur de lecture PCI)")
	} else {
		for _, ctrl := range info.Controllers {
			fmt.Printf("    • %s", ctrl.Type)
			if ctrl.MaxPorts > 0 {
				fmt.Printf(" (%d ports max)", ctrl.MaxPorts)
			}
			fmt.Printf(" [%s]\n", ctrl.PCIAddr)
		}
	}

	// Section 2 : Ports USB-C
	fmt.Println("\n  Ports USB-C physiques")
	fmt.Println("  " + strings.Repeat("-", 66))

	if len(info.USBCPorts) == 0 {
		fmt.Println("    Aucun port USB-C détecté ou non supporté")
	} else {
		for _, port := range info.USBCPorts {
			fmt.Printf("    • %s : %s, %s",
				port.Name,
				port.PowerRole,
				port.DataRole)
			if port.PowerOpMode != "" {
				fmt.Printf(" (%s)", port.PowerOpMode)
			}
			fmt.Println()
		}
	}

	// Section 3 : Devices connectés
	fmt.Println("\n  Devices USB connectés")
	fmt.Println("  " + strings.Repeat("-", 66))

	if len(info.Devices) == 0 {
		fmt.Println("    Aucun device USB détecté")
	} else {
		// Grouper par classe de vitesse
		usb2Count := 0
		usb3Count := 0

		for _, dev := range info.Devices {
			if strings.Contains(dev.SpeedClass, "USB 2") {
				usb2Count++
			} else if strings.Contains(dev.SpeedClass, "USB 3") {
				usb3Count++
			}
		}

		fmt.Printf("    Total : %d devices\n", len(info.Devices))
		if usb2Count > 0 {
			fmt.Printf("      • USB 2.0 : %d devices\n", usb2Count)
		}
		if usb3Count > 0 {
			fmt.Printf("      • USB 3.0+ : %d devices\n", usb3Count)
		}

		fmt.Println("\n    Détails :")
		for i, dev := range info.Devices {
			vendor := dev.Vendor
			if vendor == "" {
				vendor = "(fabricant inconnu)"
			}
			product := dev.Product
			if product == "" {
				product = "(produit inconnu)"
			}

			fmt.Printf("      %d. %s - %s\n", i+1, vendor, product)
			fmt.Printf("         %s (%s Mbps)\n", dev.SpeedClass, dev.Speed)
		}
	}

	fmt.Println("\n  " + strings.Repeat("-", 66))
	fmt.Println("  Note : Nombre exact de ports USB-A nécessite inspection visuelle")
	fmt.Println(strings.Repeat("=", 70) + "\n")
}

// PrintSerialPorts affiche tous les ports série détectés
func PrintSerialPorts() {
	ports, err := ListSerialPorts()
	if err != nil {
		fmt.Printf("Erreur: %v\n", err)
		return
	}

	if len(ports) == 0 {
		fmt.Println("\nAucun port série détecté.")
		return
	}

	fmt.Printf("\n" + strings.Repeat("=", 70) + "\n")
	fmt.Printf("  Ports série détectés : %d\n", len(ports))
	fmt.Printf(strings.Repeat("=", 70) + "\n\n")

	for i, port := range ports {
		fmt.Printf("Port série #%d\n", i+1)
		fmt.Println(strings.Repeat("-", 70))
		fmt.Printf("  Nom             : %s\n", port.Name)
		fmt.Printf("  Device          : %s\n", port.Device)

		if port.Driver != "" {
			fmt.Printf("  Driver          : %s\n", port.Driver)
		} else {
			fmt.Printf("  Driver          : (non disponible)\n")
		}

		fmt.Println(strings.Repeat("-", 70))

		if i < len(ports)-1 {
			fmt.Println()
		}
	}
	fmt.Println()
}
