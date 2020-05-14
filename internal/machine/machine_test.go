package machine

import "testing"

func Test_Hostname(t *testing.T) {
	t.Log(Hostname())
}

func Test_HardwareAddr(t *testing.T) {
	t.Log(HardwareAddr())
}
