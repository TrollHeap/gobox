package ram

import (
	"os/exec"
)

func GetMemoryInfo() (MemoryInfo, error) {
	cmd := exec.Command("dmidecode", "-t", "memory")
	out, err := cmd.Output()
	if err != nil {
		return MemoryInfo{}, err
	}
	slots := parseDMIDecodeOutput(out)

	total := 0
	for _, s := range slots {
		total += s.SizeMB
	}

	return MemoryInfo{
		Slots:   slots,
		TotalMB: total,
	}, nil
}
