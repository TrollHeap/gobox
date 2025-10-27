package disk

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// DiskInfo contient les informations d'un disque
type DiskInfo struct {
	Name       string   // Nom du périphérique (ex: "sda", "nvme0n1")
	Vendor     string   // Fabricant (ex: "Samsung", peut être vide)
	Model      string   // Modèle (ex: "860 EVO", peut être vide)
	Type       string   // Type de disque : "SSD" ou "HDD"
	SizeBytes  int64    // Taille totale en octets
	Partitions []string // Liste des partitions (ex: ["nvme0n1p1", "nvme0n1p2"])
}

func GetDiskInfo(diskName string) (*DiskInfo, error) {
	if err := validateDiskName(diskName); err != nil {
		return nil, err
	}
	info := &DiskInfo{
		Name: diskName,
	}
	size, err := readDiskSize(diskName)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture taille disque %s: %w", diskName, err)
	}

	info.SizeBytes = size // 2. Lire taille

	diskType, err := readDiskType(diskName)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture type disque %s: %w", diskName, err)
	}

	info.Type = diskType // 3. Lire type
	vendor, model, err := readModelAndVendor(diskName)
	if err != nil {
		log.Printf("Warning: impossible de lire vendor/model du disque %s : %v", diskName, err)
		vendor, model = "", ""
	}

	info.Vendor = vendor // 4. Lire vendor
	info.Model = model   // 5. Lire modèle

	// Ajouter les partitions
	partitions, err := listPartitions(diskName)
	if err != nil {
		log.Printf("Warning: impossible de lire partitions du disque %s : %v", diskName, err)
		partitions = []string{} // Liste vide si erreur
	}
	info.Partitions = partitions

	return info, nil
}

func ListDisks() ([]string, error) {
	entries, err := os.ReadDir(pathRoot)
	if err != nil {
		return nil, fmt.Errorf("lecture %s: %w", pathRoot, err)
	}

	diskNames := make([]string, 0, len(entries))

	for _, entry := range entries {
		name := entry.Name()

		// Filtrer périphériques virtuels en premier (plus rapide)
		if strings.HasPrefix(name, "loop") ||
			strings.HasPrefix(name, "zram") ||
			strings.HasPrefix(name, "ram") {
			continue
		}

		// Vérifier si c'est un symlink pointant vers un répertoire
		if entry.Type()&os.ModeSymlink != 0 {
			fullPath := filepath.Join(pathRoot, name)
			info, err := os.Stat(fullPath)
			if err != nil || !info.IsDir() {
				continue
			}
			diskNames = append(diskNames, name)
		} else if entry.IsDir() {
			// Cas où ce serait un répertoire direct (rare dans /sys/block/)
			diskNames = append(diskNames, name)
		}
	}

	return diskNames, nil
}
