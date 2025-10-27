package usb

import "time"

// USBPort représente un port USB détecté
type USBPort struct {
	BusNum     int    // Numéro de bus USB (ex: 1)
	DevNum     int    // Numéro de device (ex: 3)
	VendorID   string // ID vendeur (ex: "1d6b")
	ProductID  string // ID produit (ex: "0002")
	Speed      string // Vitesse (ex: "480" = USB 2.0)
	MaxPower   int    // Consommation max en mA
	DevicePath string // Chemin sysfs complet
}

// TestResult représente le résultat d'un test fonctionnel
type TestResult struct {
	Port     USBPort
	Passed   bool
	Latency  time.Duration // Temps de réponse loopback
	ErrorMsg string
}
