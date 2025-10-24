package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DiskInfo contient les informations d'un disque
type DiskInfo struct {
	Name      string // Nom du périphérique (ex: "sda", "nvme0n1")
	Vendor    string // Fabricant (ex: "Samsung", peut être vide)
	Model     string // Modèle (ex: "860 EVO", peut être vide)
	Type      string // Type de disque : "SSD" ou "HDD"
	SizeBytes int64  // Taille totale en octets
}

const (
	pathRoot = "/sys/block/"
)

func ListDisks() ([]string, error) {
	entries, err := os.ReadDir(pathRoot)
	if err != nil {
		return nil, fmt.Errorf("lecture %s: %w", pathRoot, err)
	}

	var diskNames []string

	for _, entry := range entries {
		// Ignorer fichiers (garder seulement répertoires)
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Filtrer périphériques virtuels
		if strings.HasPrefix(name, "loop") ||
			strings.HasPrefix(name, "zram") ||
			strings.HasPrefix(name, "ram") {
			continue
		}

		diskNames = append(diskNames, name)
	}

	return diskNames, nil
}

func ListPartitions(diskName string) ([]string, error) {
	diskPath := filepath.Join(pathRoot, diskName)

	entries, err := os.ReadDir(diskPath)
	if err != nil {
		return nil, fmt.Errorf("lecture %s: %w", diskPath, err)
	}

	var partitions []string

	for _, entry := range entries {
		// Garder seulement répertoires (partitions sont des sous-répertoires)
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Vérifier si c'est une partition du disque
		if suffix, found := strings.CutPrefix(name, diskName); found {
			if suffix != "" && (strings.HasPrefix(suffix, "p") ||
				(len(suffix) > 0 && suffix[0] >= '0' && suffix[0] <= '9')) {
				partitions = append(partitions, name)
			}
		}
	}

	return partitions, nil
}

func PrintDisk() {
	fmt.Println("Disk Info:")
	fmt.Println(ListDisks())
	fmt.Println("Partitions :")
	DiskPartitions, _ := ListPartitions("nvme0n1")
	fmt.Println(DiskPartitions)
}

// TODO:

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
