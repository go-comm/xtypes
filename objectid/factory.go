package objectid

import (
	"errors"
	"sync"
	"time"
)

const (
	sequenceMask = 0x00FFFFFF
)

var (
	// Starting from January 1, 2020
	defaultStartEpoch = time.Unix(1577808000, 0).UnixNano()
	defaultFactory    = NewFactory()
)

var (
	ErrTimeBackwards = errors.New("objectid: time has gone backwards")
)

type Option func(s *factory)

func WithStartEpoch(startEpoch int64) Option {
	return func(s *factory) {
		s.startEpoch = startEpoch
	}
}

func WithNodeID(nodeID uint) Option {
	return func(s *factory) {
		s.nodeID = nodeID
	}
}

func New() (ID, error) {
	return defaultFactory.New()
}

type Factory interface {
	New() (ID, error)
}

func NewFactory(options ...Option) Factory {
	fac := &factory{}
	fac.startEpoch = defaultStartEpoch
	for _, opt := range options {
		opt(fac)
	}
	return fac
}

type factory struct {
	mutex         sync.Mutex
	startEpoch    int64
	lastTimestamp int64
	sequence      uint32
	nodeID        uint
}

func (f *factory) New() (id ID, err error) {
	nodeID := f.nodeID
	id[6] = byte(nodeID >> 16)
	id[7] = byte((nodeID >> 8))
	id[8] = byte((nodeID))

	var seq uint32
LOOP:
	// divide by 1048576, giving pseudo-milliseconds
	ts := (time.Now().UnixNano() - f.startEpoch) >> 20

	f.mutex.Lock()
	if ts < f.lastTimestamp {
		f.mutex.Unlock()
		return nilID, ErrTimeBackwards
	}

	if f.lastTimestamp == ts {
		seq = (f.sequence + 1) & sequenceMask
		if seq == 0 {
			f.mutex.Unlock()
			goto LOOP
		}
		f.sequence = seq
	} else {
		f.sequence = 0
	}
	f.lastTimestamp = ts
	f.mutex.Unlock()

	id[0] = byte(ts >> 40)
	id[1] = byte(ts >> 32)
	id[2] = byte(ts >> 24)
	id[3] = byte(ts >> 16)
	id[4] = byte(ts >> 8)
	id[5] = byte(ts)

	id[9] = byte(seq >> 16)
	id[10] = byte(seq >> 8)
	id[11] = byte(seq)
	return
}
