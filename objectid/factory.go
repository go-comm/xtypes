package objectid

import (
	"encoding/binary"
	"sync/atomic"
	"time"

	"github.com/go-comm/xtypes/internal/couter"
	"github.com/go-comm/xtypes/internal/machine"
)

var (
	// Starting from January 1, 2020
	defaultStartEpoch = time.Unix(1577808000, 0)

	defaultFactory = NewFactory()
)

type Factory interface {
	New() (id ID)
	Update(id *ID)
}

func New() ID {
	return defaultFactory.New()
}

func Update(id *ID) {
	defaultFactory.Update(id)
}

func NewFactory() *factory {
	return NewFactoryWithStartEpoch(defaultStartEpoch)
}

func NewFactoryWithStartEpoch(startEpoch time.Time) *factory {
	return &factory{startEpoch: startEpoch, idCounter: couter.Couter()}
}

type factory struct {
	startEpoch time.Time
	idCounter  uint32
}

func (f *factory) New() (id ID) {
	f.Update(&id)
	return id
}

func (f *factory) Update(id *ID) {
	ts := uint32(time.Now().Unix() - f.startEpoch.Unix())
	binary.BigEndian.PutUint32((*id)[:], ts)

	machineID := machine.HardwareAddr()
	(*id)[4] = machineID[0]
	(*id)[5] = machineID[1]
	(*id)[6] = machineID[2]

	pid := machine.PID()
	(*id)[7] = byte(pid >> 8)
	(*id)[8] = byte(pid)

	i := atomic.AddUint32(&f.idCounter, 1)
	(*id)[9] = byte(i >> 16)
	(*id)[10] = byte(i >> 8)
	(*id)[11] = byte(i)
}
