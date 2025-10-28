package ports

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SerialPort représente un port série
type SerialPort struct {
	Name   string // Ex: "ttyS0", "ttyUSB0"
	Driver string // Ex: "serial8250", "usbserial"
	Device string // Chemin device (/dev/ttyS0)
}

// ListSerialPorts liste tous les ports série physiques
func ListSerialPorts() ([]SerialPort, error) {
	const ttyRoot = "/sys/class/tty/"

	entries, err := os.ReadDir(ttyRoot)
	if err != nil {
		return nil, fmt.Errorf("lecture %s: %w", ttyRoot, err)
	}

	ports := make([]SerialPort, 0)

	for _, entry := range entries {
		name := entry.Name()

		// Filtrer console et ttys système (non-série)
		if name == "console" ||
			name == "tty" ||
			strings.HasPrefix(name, "tty") && len(name) == 4 && name[3] >= '0' && name[3] <= '9' {
			continue // tty0-tty9 virtuels
		}

		devicePath := filepath.Join(ttyRoot, name)

		// Vérifier présence de /device/ (indique port physique ou USB)
		deviceLinkPath := filepath.Join(devicePath, "device")
		if _, err := os.Stat(deviceLinkPath); os.IsNotExist(err) {
			continue
		}

		port := SerialPort{
			Name:   name,
			Device: "/dev/" + name,
		}

		// Lire nom du driver
		driverPath := filepath.Join(devicePath, "device", "driver")
		if target, err := os.Readlink(driverPath); err == nil {
			port.Driver = filepath.Base(target)
		}

		ports = append(ports, port)
	}

	return ports, nil
}

func validatePortName(name string) error {
	if name == "" {
		return fmt.Errorf("nom de port vide")
	}
	if strings.Contains(name, "..") || strings.Contains(name, "/") {
		return fmt.Errorf("nom de port invalide : %s", name)
	}
	return nil
}
