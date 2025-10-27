package network

// NetworkInterface représente une carte réseau physique
type NetworkInterface struct {
	Name       string // Nom (ex: "eth0", "wlan0")
	MACAddress string // Adresse MAC matérielle
	Type       string // Type ("ethernet", "wifi")
	Driver     string // Pilote utilisé (ex: "e1000e", "iwlwifi")
	IsUp       bool   // État matériel (pas l'état réseau)
	IPAddress  string // IP actuelle (info bonus)
}
