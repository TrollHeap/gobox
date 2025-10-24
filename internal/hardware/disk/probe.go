package disk

// DiskInfo contient les informations d'un disque
type DiskInfo struct {
	Name      string // Nom du périphérique (ex: "sda", "nvme0n1")
	Vendor    string // Fabricant (ex: "Samsung", peut être vide)
	Model     string // Modèle (ex: "860 EVO", peut être vide)
	Type      string // Type de disque : "SSD" ou "HDD"
	SizeBytes int64  // Taille totale en octets
}

// cat /sys/block/nvme0n1/queue/rotational
// =====> 0
// 3. Type (rotationnel)
// 0  ← Non rotationnel = SSD
// Règle :
//     0 → SSD
//     1 → HDD
//
// Model + Vendor
// =================
// cat  /sys/block/nvme0n1/device/model
// =====> SKHynix_HFS001TEJ9X115N
//
// SizeBytes
// =======================
// cat /sys/block/nvme0n1/size
// =====> 2000409264
