package ports

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const usbRoot = "/sys/bus/usb/devices/"

// USBPort représente un port USB physique
type USBPort struct {
	Name    string // Ex: "1-4" (bus 1, port 4)
	Speed   string // Ex: "480" (Mbps)
	Product string // Ex: "Mass Storage Device"
	Vendor  string // Ex: "SanDisk"
	BusNum  string // Numéro du bus
	DevNum  string // Numéro du device
}

// ListUSBPorts liste tous les ports USB avec devices connectés
func ListUSBPorts() ([]USBPort, error) {
	entries, err := os.ReadDir(usbRoot)
	if err != nil {
		return nil, fmt.Errorf("lecture %s: %w", usbRoot, err)
	}

	ports := make([]USBPort, 0, len(entries))

	for _, entry := range entries {
		name := entry.Name()

		// Filtrer hubs racines et interfaces
		if strings.HasPrefix(name, "usb") || strings.Contains(name, ":") {
			continue
		}

		devicePath := filepath.Join(usbRoot, name)

		// Vérifier que c'est un device réel (pas un lien mort)
		if _, err := os.Stat(devicePath); os.IsNotExist(err) {
			continue
		}

		port := USBPort{Name: name}

		// Lecture infos (avec gestion erreurs silencieuse)
		if data, err := os.ReadFile(filepath.Join(devicePath, "speed")); err == nil {
			port.Speed = strings.TrimSpace(string(data))
		}

		if data, err := os.ReadFile(filepath.Join(devicePath, "product")); err == nil {
			port.Product = strings.TrimSpace(string(data))
		}

		if data, err := os.ReadFile(filepath.Join(devicePath, "manufacturer")); err == nil {
			port.Vendor = strings.TrimSpace(string(data))
		}

		if data, err := os.ReadFile(filepath.Join(devicePath, "busnum")); err == nil {
			port.BusNum = strings.TrimSpace(string(data))
		}

		if data, err := os.ReadFile(filepath.Join(devicePath, "devnum")); err == nil {
			port.DevNum = strings.TrimSpace(string(data))
		}

		ports = append(ports, port)
	}

	return ports, nil
}
