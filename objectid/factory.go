package objectid

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/go-comm/xtypes"
	"github.com/go-comm/xtypes/internal/hardware"
)

const (
	// Size of id
	Size = 12
)

var (
	// Starting from January 1, 2020
	defaultStartEpoch int64 = 1577808000 * int64(time.Second)
	defaultFactory          = NewFactory(defaultStartEpoch)
)

func DefaultFactory() xtypes.Factory {
	return defaultFactory
}

func NewFactory(startEpoch int64) xtypes.Factory {
	f := &factory{
		startEpoch: startEpoch,
		rand:       rand.Reader,
	}
	return f
}

type factory struct {
	clockSeqOnce     sync.Once
	hardwareAddrOnce sync.Once

	rand io.Reader

	seq          uint32
	startEpoch   int64
	hardwareAddr []byte
	mutex        sync.RWMutex
	lastTime     uint32
}

func (f *factory) New() xtypes.Object {
	id := new(ID)
	f.Reset(id)
	return id
}

func (f *factory) Reset(o xtypes.Object) error {
	id, ok := o.(*ID)
	if !ok {
		return fmt.Errorf("objectid: expect type %T", o)
	}
	now, next := f.getClockSeq()

	binary.BigEndian.PutUint16(id[0:], uint16(now&0xFFFF))
	binary.BigEndian.PutUint16(id[2:], uint16((now>>16)&0xFFFF))

	copy(id[4:], f.getHardwareAddr()[:3])
	binary.BigEndian.PutUint16(id[7:], uint16(os.Getpid()&0xFFFF))

	id[9] = uint8((next >> 16) & 0xFF)
	binary.BigEndian.PutUint16(id[10:], uint16(next&0xFFFF))
	return nil
}

func (f *factory) getClockSeq() (now uint32, next uint32) {
	f.clockSeqOnce.Do(func() {
		buf := make([]byte, 4)
		io.ReadFull(f.rand, buf)
		f.seq = binary.BigEndian.Uint32(buf)
	})
	now = uint32((time.Now().UnixNano() - f.startEpoch) / int64(time.Second))

	f.mutex.Lock()
	if now <= f.lastTime {
		f.seq++
	}
	next = f.seq
	f.lastTime = now
	f.mutex.Unlock()
	return
}

func (f *factory) getHardwareAddr() []byte {
	f.hardwareAddrOnce.Do(func() {
		ifaceAddr := hardware.GetHardwareAddr()
		if len(ifaceAddr) >= 3 {
			ifaceAddr = ifaceAddr[len(ifaceAddr)-3:]
		} else {
			ifaceAddr = make([]byte, 3)
			io.ReadFull(f.rand, ifaceAddr)
		}
		f.hardwareAddr = ifaceAddr
	})
	return f.hardwareAddr
}

func New() xtypes.Object {
	return defaultFactory.New()
}

func Reset(o xtypes.Object) error {
	return defaultFactory.Reset(o)
}
