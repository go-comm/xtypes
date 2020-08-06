package objectid

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/go-comm/xtypes/internal/machine"
)

const (
	sequenceMask = 0x00FFFFFF
)

var (
	// Starting from January 1, 2020
	defaultStartEpoch       = time.Unix(1577808000, 0).UnixNano()
	defaultFactory          = NewFactory()
	hashIncrement     int32 = 0x61c88647
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
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fac.sequence = r.Int31() & sequenceMask
	fac.initSeq = fac.sequence
	return fac
}

type factory struct {
	mutex         sync.Mutex
	startEpoch    int64
	lastTimestamp int64
	initSeq       int32
	sequence      int32
	nodeID        uint
}

func (f *factory) New() (id ID, err error) {
	v := ((f.nodeID & 0x0FFF) << 12) | uint(machine.PID()&0x0FFF)
	id[6] = byte(v >> 16)
	id[7] = byte(v >> 8)
	id[8] = byte(v)

	var seq int32
LOOP:
	// divide by 1048576, giving pseudo-milliseconds
	ts := (time.Now().UnixNano() - f.startEpoch) >> 20

	f.mutex.Lock()
	if ts < f.lastTimestamp {
		f.mutex.Unlock()
		return nilID, ErrTimeBackwards
	}

	seq = (f.sequence + hashIncrement) & sequenceMask
	if f.lastTimestamp == ts {
		if seq == f.initSeq {
			f.mutex.Unlock()
			goto LOOP
		}
		f.sequence = seq
	} else {
		f.sequence = seq
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
