package main

import (
	"gobox/internal/hardware/network"
)

func main() {
	// print.CPUInfo()
	// ports.PrintSerialPorts()
	// ports.PrintUSBInfo()
	network.PrintNetworkInterfaces()
	// ports.PrintUSBCPorts()
	// print.DiskInfo()
}
