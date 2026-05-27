package util

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	epoch            = int64(1704067200000)
	timestampBits    = uint(41)
	workerIDBits     = uint(10)
	sequenceBits     = uint(12)
	timestampShift   = workerIDBits + sequenceBits
	workerIDShift    = sequenceBits
	sequenceMask     = int64(-1) ^ (int64(-1) << sequenceBits)
	maxWorkerID      = int64(-1) ^ (int64(-1) << workerIDBits)
)

type Snowflake struct {
	sync.Mutex
	timestamp int64
	workerID   int64
	sequence   int64
}

var instance *Snowflake
var once sync.Once

func NewSnowflake(workerID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, errors.New("worker ID must be between 0 and 1023")
	}
	return &Snowflake{
		timestamp: 0,
		workerID:  workerID,
		sequence:  0,
	}, nil
}

func GetSnowflakeInstance() *Snowflake {
	once.Do(func() {
		instance, _ = NewSnowflake(1)
	})
	return instance
}

func (s *Snowflake) Generate() int64 {
	s.Lock()
	defer s.Unlock()

	now := time.Now().UnixNano() / 1000000

	if now == s.timestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = now

	return ((now - epoch) << timestampShift) |
		(s.workerID << workerIDShift) |
		s.sequence
}

func GenerateOrderNo() string {
	sf := GetSnowflakeInstance()
	id := sf.Generate()
	return fmt.Sprintf("V%d", id)
}
