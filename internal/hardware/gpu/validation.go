package gpu

import "regexp"

// isHexString vérifie qu'une chaîne ne contient que des caractères hexadécimaux.
func isHexString(s string) bool {
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}
	return true
}

// pciSlotRegex valide le format PCI slot standard (XXXX:XX:XX.X).
var pciSlotRegex = regexp.MustCompile(`^[0-9a-fA-F]{4}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}\.[0-7]$`)

// isValidPCISlot vérifie le format d'un slot PCI.
func isValidPCISlot(slot string) bool {
	return pciSlotRegex.MatchString(slot)
}
