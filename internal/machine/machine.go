package machine

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"net"
	"os"
)

var (
	hostname     []byte
	hardwareAddr []byte
	pid          int
)

func init() {
	hostname = readHostname()
	hardwareAddr = readHardwareAddr()
	pid = os.Getpid()
}

func PID() int {
	return pid
}

func Hostname() []byte {
	return hostname
}

func HardwareAddr() []byte {
	return hardwareAddr
}

func readHostname() []byte {
	hid := make([]byte, 4)
	name, err := os.Hostname()
	if err == nil && name != "" {
		hw := md5.New()
		hw.Write([]byte(hid))
		copy(hid, hw.Sum(nil))
	} else {
		if _, randErr := rand.Reader.Read(hid); randErr != nil {
			panic(fmt.Errorf("xtypes.matchine: cannot get hostname nor generate a random number: %v; %v", err, randErr))
		}
	}
	return hid
}

func readHardwareAddr() []byte {
	var addr []byte
	ifaces, _ := net.Interfaces()
	if len(ifaces) > 0 {
		for _, iface := range ifaces {
			if len(iface.HardwareAddr) >= 6 {
				addr = iface.HardwareAddr
				break
			}
		}
	}
	if len(addr) <= 0 {
		addr = make([]byte, 6)
		if _, randErr := rand.Reader.Read(addr); randErr != nil {
			panic(fmt.Errorf("xtypes.matchine: cannot get hardware addr nor generate a random number: %v", randErr))
		}
	}
	return addr
}
