package messageid

import (
	"errors"
	"sync"
	"time"
)

var (
	nodeIDBits     = uint64(10)
	sequenceBits   = uint64(12)
	nodeIDShift    = sequenceBits
	timestampShift = sequenceBits + nodeIDBits
	sequenceMask   = int64(-1) ^ (int64(-1) << sequenceBits)

	// Starting from January 1, 2020
	defaultStartEpoch = time.Unix(1577808000, 0).UnixNano()

	defaultFactory = NewFactory()
)

var (
	ErrTimeBackwards = errors.New("messageid: time has gone backwards")
)

type Option func(s *factory)

func WithStartEpoch(startEpoch int64) Option {
	return func(s *factory) {
		s.startEpoch = startEpoch
	}
}

func WithNodeID(nodeID int64) Option {
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

type factory struct {
	mutex sync.Mutex

	startEpoch    int64
	lastTimestamp int64
	sequence      int64
	nodeID        int64
	lastID        ID
}

func NewFactory(options ...Option) Factory {
	fac := &factory{}
	fac.startEpoch = defaultStartEpoch
	for _, opt := range options {
		opt(fac)
	}
	return fac
}

func (f *factory) New() (ID, error) {
	var seq int64
LOOP:
	// divide by 1048576, giving pseudo-milliseconds
	ts := time.Now().UnixNano() >> 20

	f.mutex.Lock()
	if ts < f.lastTimestamp {
		f.mutex.Unlock()
		return 0, ErrTimeBackwards
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

	id := ID(((ts - f.startEpoch) << timestampShift) |
		(f.nodeID << nodeIDShift) |
		seq)

	return id, nil
}
