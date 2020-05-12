package hardware

import (
	"net"
)

func GetHardwareAddr() []byte {
	var ifaceAddr []byte
	ifaces, _ := net.Interfaces()
	if len(ifaces) > 0 {
		for _, iface := range ifaces {
			if len(iface.HardwareAddr) >= 6 {
				ifaceAddr = iface.HardwareAddr
				break
			}
		}
	}
	return ifaceAddr
}
