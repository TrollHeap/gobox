package disk

import (
	"fmt"
	"os"
)

// DiskInfo contient les informations d'un disque
type DiskInfo struct {
	Name      string // Nom du périphérique (ex: "sda", "nvme0n1")
	Vendor    string // Fabricant (ex: "Samsung", peut être vide)
	Model     string // Modèle (ex: "860 EVO", peut être vide)
	Type      string // Type de disque : "SSD" ou "HDD"
	SizeBytes int64  // Taille totale en octets
}

type Disks struct {
	Disk DiskInfo
}

const (
	pathRoot = "/sys/block/"
)

func ListDisks() {
	path, _ := os.ReadFile(pathRoot)
	// if err != nil {
	// 	return nil, err
	// }
	for _, line := range path {
		println(line)
	}
}

func PrintDisk() {
	fmt.Println("Disk Info:")
	ListDisks()
}

// TODO:
/*
1. Lister /sys/block/* (tous les périphériques)
2. Filtrer :
   - Exclure loop* (loop devices)
   - Exclure zram* (swap compressé)
   - Exclure ram* (ramdisk)
   - Garder seulement les disques physiques
3. Pour chaque disque trouvé :
   - Appeler GetDiskInfo(diskName)
4. Retourner slice de DiskInfo
**Exemple d'appel** :
 disks, err := ListDisks()
 disks = [nvme0n1, sda, sdb, ...]
 for _, disk := range disks {
     fmt.Printf("%s: %d bytes\n", disk.Name, disk.SizeBytes)
 }
*/

/*
Nom : readDisks
path : /sys/block/nvme0n1/device/model
Entrée : nom du disque (string)
Sortie : taille en octets (int64), erreur (error)
=======================
cat /sys/block/nvme0n1/size
=====> 2000409264
Étapes :
  1. Construire chemin vers /sys/block/{disk}/size
  2. Lire fichier
  3. TrimSpace
  4. Convertir string → int64 (secteurs)
  5. Multiplier par 512 → octets
  6. Return
  Exemple : SKHynix_HFS001TEJ9X115N

*/

/*
Nom : readModelAndVendor
DiskInfo.Vendor
DiskInfo.Model
path : /sys/block/nvme0n1/device/model
Entrée : nom du disque (string)
Sortie : taille en octets (int64), erreur (error)
=======================
cat /sys/block/nvme0n1/size
=====> 2000409264
Étapes :
  1. Construire chemin vers /sys/block/{disk}/size
  2. Lire fichier
  3. TrimSpace
  4. Convertir string → int64 (secteurs)
  5. Multiplier par 512 → octets
  6. Return
  Exemple : SKHynix_HFS001TEJ9X115N

*/

/*
Nom : readDiskSize
DiskInfo.SizeBytes
path : /sys/block/nvme0n1/size
Entrée : nom du disque (string)
Sortie : taille en octets (int64), erreur (error)
=======================
cat /sys/block/nvme0n1/size
=====> 2000409264
Étapes :
  1. Construire chemin vers /sys/block/{disk}/size
  2. Lire fichier
  3. TrimSpace
  4. Convertir string → int64 (secteurs)
  5. Multiplier par 512 → octets
  6. Return

*/

/*
Nom : readDiskType
DiskInfo.Type
path: /sys/block/nvme0n1/queue/rotational
Entrée : nom du disque (string)
Sortie : type (string "SSD" ou "HDD"), erreur (error)
Étapes :
  1. Construire chemin vers /sys/block/{disk}/queue/rotational
  2. Lire fichier
  3. TrimSpace
  4. Condition : si "0" → "SSD", si "1" → "HDD"
  5. Return
Règle :
    0 → SSD
    1 → HDD
*/

/*
func GetDiskInfo(diskName string) (*DiskInfo, error) {
		// 1. Créer struct vide

		// 2. Remplir champ Name (facile, c'est l'argument)

		// 3. Lire Vendor (appeler fonction helper)
		// Condition : si erreur, mettre "" (optionnel)

		// 4. Lire Model (appeler fonction helper)
		// Condition : si erreur, mettre ""

		// 5. Lire Type (appeler fonction helper)
		// Condition : si erreur, ??? (obligatoire ou optionnel ?)

		//6. Lire Size (appeler fonction helper)
		//Condition : si erreur, ??? (obligatoire !)

		//7. Return struct complétée
}
info, err := GetDiskInfo("nvme0n1")
info contient Name, Vendor, Model, Type, Size

*/

/*
#### Étape 1 : Glob `/sys/block/*`

Fonction : filepath.Glob
Pattern : "/sys/block/*"
Résultat : []string{
    "/sys/block/loop0",
    "/sys/block/nvme0n1",
    "/sys/block/sda",
    "/sys/block/zram0",
}

#### Étape 2 : Filtrer les noms

**Condition pour garder**  :[3][2]

```
Garder si :
  - Nom ne commence PAS par "loop"
  - Nom ne commence PAS par "zram"
  - Nom ne commence PAS par "ram"
  - (Optionnel) Nom ne contient PAS un chiffre AVANT une lettre
    (exclut partitions comme "nvme0n1p1")
```

**Exemple de filtrage** :
```
/sys/block/loop0     → SKIP (loop device)
/sys/block/nvme0n1   → KEEP (disque physique)
/sys/block/nvme0n1p1 → N'existe PAS dans /sys/block (partitions ailleurs)
/sys/block/sda       → KEEP
/sys/block/zram0     → SKIP (swap)
```

**Structure sysfs**  :[3][2]

```
/sys/block/
├── nvme0n1/           ← DISQUE (tu le vois)
│   ├── nvme0n1p1/     ← PARTITION (subdirectory)
│   ├── nvme0n1p2/
│   ├── nvme0n1p3/
│   └── nvme0n1p4/
├── sda/               ← DISQUE
│   ├── sda1/          ← PARTITION
│   └── sda2/
└── zram0/
```
*/
