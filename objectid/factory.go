package objectid

import (
	"encoding/binary"
	"sync/atomic"
	"time"

	"github.com/go-comm/xtypes/internal/machine"
)

var (
	// Starting from January 1, 2020
	defaultStartEpoch = time.Unix(1577808000, 0)

	defaultFactory = NewFactory()
)

type Option func(s *factory)

func WithStartEpoch(startEpoch time.Time) Option {
	return func(s *factory) {
		s.startEpoch = startEpoch
	}
}

func WithMachineID(machineID uint) Option {
	return func(s *factory) {
		s.machineID = machineID
	}
}

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

func NewFactory(options ...Option) *factory {
	fac := &factory{}
	fac.startEpoch = defaultStartEpoch
	for _, opt := range options {
		opt(fac)
	}
	return fac
}

type factory struct {
	startEpoch time.Time
	idCounter  uint32
	machineID  uint
}

func (f *factory) New() (id ID) {
	f.Update(&id)
	return id
}

func (f *factory) Update(id *ID) {
	ts := uint32(time.Now().Unix() - f.startEpoch.Unix())
	binary.BigEndian.PutUint32((*id)[:], ts)

	machineID := f.machineID
	(*id)[4] = byte(machineID & 0xFF)
	(*id)[5] = byte((machineID >> 8) & 0xFF)
	(*id)[6] = byte((machineID >> 16) & 0xFF)

	pid := machine.PID()
	(*id)[7] = byte(pid >> 8)
	(*id)[8] = byte(pid)

	i := atomic.AddUint32(&f.idCounter, 1)
	(*id)[9] = byte(i >> 16)
	(*id)[10] = byte(i >> 8)
	(*id)[11] = byte(i)
}
