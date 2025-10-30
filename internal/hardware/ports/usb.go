package ports

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	usbRoot          = "/sys/bus/usb/devices/"
	typecRoot        = "/sys/class/typec/"
	pciRoot          = "/sys/bus/pci/devices/"
	maxSysfsFileSize = 4096 // 4 KB max pour fichiers sysfs
)

// USBInfo contient toutes les informations USB du système
type USBInfo struct {
	Controllers []USBController // Contrôleurs USB détectés
	Devices     []USBDevice     // Devices USB connectés
	USBCPorts   []USBCPort      // Ports USB-C physiques
}

// USBDevice représente un device USB connecté
type USBDevice struct {
	Name       string // Ex: "1-4" (bus 1, port 4)
	Speed      string // Ex: "480" (Mbps)
	SpeedClass string // Ex: "USB 2.0 High-Speed"
	Product    string // Ex: "Mass Storage Device"
	Vendor     string // Ex: "SanDisk"
	BusNum     string // Numéro du bus
	DevNum     string // Numéro du device
}

// USBCPort représente un port USB-C physique
type USBCPort struct {
	Name        string // Ex: "port0"
	PowerRole   string // "source", "sink", "dual"
	DataRole    string // "host", "device", "dual"
	PowerOpMode string // "default", "1.5A", "3.0A"
}

// USBController représente un contrôleur USB sur le bus PCI
type USBController struct {
	Type     string // "USB 2.0 (EHCI)", "USB 3.0 (xHCI)", etc.
	MaxPorts int    // Nombre maximum de ports gérés
	PCIAddr  string // Adresse PCI (ex: "00:14.0")
}

func validateSysfsName(name string) error {
	if name == "" {
		return fmt.Errorf("nom vide")
	}
	if strings.Contains(name, "..") || strings.Contains(name, "/") {
		return fmt.Errorf("traversée de chemin détectée : %s", name)
	}
	// Vérifier caractères alphanumériques + tirets/underscores/points/deux-points
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == ':' || c == '.') {
			return fmt.Errorf("caractère invalide dans nom : %s", name)
		}
	}
	return nil
}

// readSysfsFile lit un fichier sysfs avec limite de taille (sécurité)
func readSysfsFile(path string, buf []byte) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Lire max 4KB (fichiers sysfs sont toujours petits)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	return strings.TrimSpace(string(buf[:n])), nil
}

// ListUSBControllers détecte les contrôleurs USB via PCI
func ListUSBControllers() ([]USBController, error) {
	entries, err := os.ReadDir(pciRoot)
	if err != nil {
		return nil, fmt.Errorf("lecture %s: %w", pciRoot, err)
	}

	// 🆕 Optimisation #3 : Pré-allocation (rarement plus de 4 contrôleurs)
	controllers := make([]USBController, 0, 4)
	buf := make([]byte, 4096)

	for _, entry := range entries {
		pciAddr := entry.Name()

		// 🆕 Optimisation #2 : Validation sécurité
		if err := validateSysfsName(pciAddr); err != nil {
			continue
		}

		devicePath := filepath.Join(pciRoot, pciAddr)

		// 🆕 Optimisation #4 : Lecture avec limite de taille
		class, err := readSysfsFile(filepath.Join(devicePath, "class"), buf)
		if err != nil {
			continue
		}

		// 0x0c03 = USB controller class
		if !strings.HasPrefix(class, "0x0c03") {
			continue
		}

		controller := USBController{PCIAddr: pciAddr}

		// Déterminer le type de contrôleur (UHCI/OHCI/EHCI/xHCI)
		classCode := class

		// Derniers 2 chiffres = programming interface
		switch {
		case strings.HasSuffix(classCode, "00"):
			controller.Type = "USB 1.1 (UHCI)"
		case strings.HasSuffix(classCode, "10"):
			controller.Type = "USB 1.1 (OHCI)"
		case strings.HasSuffix(classCode, "20"):
			controller.Type = "USB 2.0 (EHCI)"
		case strings.HasSuffix(classCode, "30"):
			controller.Type = "USB 3.0 (xHCI)"
		case strings.HasSuffix(classCode, "40"):
			controller.Type = "USB 4.0"
		default:
			controller.Type = "USB (type inconnu)"
		}

		// Tenter de lire le nombre de ports (pas toujours disponible)
		if maxChild, err := readSysfsFile(filepath.Join(devicePath, "max_child_bus_number"), buf); err == nil {
			fmt.Sscanf(maxChild, "%d", &controller.MaxPorts)
		}

		controllers = append(controllers, controller)
	}

	return controllers, nil
}

// ListUSBDevices liste tous les devices USB connectés
func ListUSBDevices() ([]USBDevice, error) {
	buf := make([]byte, 4096)
	entries, err := os.ReadDir(usbRoot)
	if err != nil {
		return nil, fmt.Errorf("lecture %s: %w", usbRoot, err)
	}

	devices := make([]USBDevice, 0, len(entries))

	for _, entry := range entries {
		name := entry.Name()

		// 🆕 Optimisation #2 : Validation sécurité
		if err := validateSysfsName(name); err != nil {
			continue
		}

		// Filtrer hubs racines et interfaces
		if strings.HasPrefix(name, "usb") || strings.Contains(name, ":") {
			continue
		}

		devicePath := filepath.Join(usbRoot, name)

		// Vérifier que c'est un device réel (pas un lien mort)
		if _, err := os.Stat(devicePath); os.IsNotExist(err) {
			continue
		}

		device := USBDevice{Name: name}

		if speed, err := readSysfsFile(filepath.Join(devicePath, "speed"), buf); err == nil {
			device.Speed = speed

			// Classifier immédiatement après lecture
			switch speed {
			case "1.5", "12":
				device.SpeedClass = "USB 2.0"
			case "480":
				device.SpeedClass = "USB 2.0 High-Speed"
			case "5000":
				device.SpeedClass = "USB 3.0"
			case "10000":
				device.SpeedClass = "USB 3.1"
			case "20000":
				device.SpeedClass = "USB 3.2"
			default:
				device.SpeedClass = "USB " + speed + " Mbps"
			}
		}

		if product, err := readSysfsFile(filepath.Join(devicePath, "product"), buf); err == nil {
			device.Product = product
		}

		if vendor, err := readSysfsFile(filepath.Join(devicePath, "manufacturer"), buf); err == nil {
			device.Vendor = vendor
		}

		if busNum, err := readSysfsFile(filepath.Join(devicePath, "busnum"), buf); err == nil {
			device.BusNum = busNum
		}

		if devNum, err := readSysfsFile(filepath.Join(devicePath, "devnum"), buf); err == nil {
			device.DevNum = devNum
		}

		devices = append(devices, device)
	}

	return devices, nil
}

// ListUSBCPorts liste les ports USB-C physiques
// Retourne une liste vide si non disponible (kernel ancien ou pas de USB-C)
func ListUSBCPorts() ([]USBCPort, error) {
	// Vérifier si le système expose les infos USB-C
	buf := make([]byte, 4096)
	if _, err := os.Stat(typecRoot); os.IsNotExist(err) {
		return []USBCPort{}, nil // Pas d'erreur, juste vide
	}

	entries, err := os.ReadDir(typecRoot)
	if err != nil {
		return nil, fmt.Errorf("lecture %s: %w", typecRoot, err)
	}

	ports := make([]USBCPort, 0, len(entries))

	for _, entry := range entries {
		name := entry.Name()

		if err := validateSysfsName(name); err != nil {
			continue
		}

		// Garder seulement les ports (pas les câbles ou partenaires)
		if !strings.HasPrefix(name, "port") {
			continue
		}

		portPath := filepath.Join(typecRoot, name)
		port := USBCPort{Name: name}

		if powerRole, err := readSysfsFile(filepath.Join(portPath, "power_role"), buf); err == nil {
			port.PowerRole = powerRole
		}

		if dataRole, err := readSysfsFile(filepath.Join(portPath, "data_role"), buf); err == nil {
			port.DataRole = dataRole
		}

		if powerOpMode, err := readSysfsFile(filepath.Join(portPath, "power_operation_mode"), buf); err == nil {
			port.PowerOpMode = powerOpMode
		}

		ports = append(ports, port)
	}

	return ports, nil
}
