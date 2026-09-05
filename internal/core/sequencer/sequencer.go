package sequencer

import (
	"log"
	"sync"
	"time"

	"github.com/arsrivastawa/shawttyy/internal/error"
)

const (
	epoch     = 1700000000000 // custom epoch in ms
	nodeBits  = 10
	seqBits   = 12
	maxNode   = -1 ^ (-1 << nodeBits)
	maxSeq    = -1 ^ (-1 << seqBits)
	nodeShift = seqBits
	timeShift = seqBits + nodeBits
)

type Sequencer struct {
	mu       sync.Mutex
	lastTime int64
	seq      int64
	nodeID   int64
}

func New(nodeID int64) (*Sequencer, error) {
	if nodeID < 0 || nodeID > maxNode {
		return nil, error.ErrNodeOutOfRange
	}
	return &Sequencer{nodeID: nodeID}, nil
}

// Next returns the next snowflake-style 64-bit ID.
func (s *Sequencer) Next() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	if now < s.lastTime {
		log.Printf("sequencer: clock moved backwards, waiting")
		for now < s.lastTime {
			now = time.Now().UnixMilli()
		}
	}

	if now == s.lastTime {
		s.seq = (s.seq + 1) & maxSeq
		if s.seq == 0 {
			for now <= s.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.seq = 0
	}

	s.lastTime = now

	return (now-epoch)<<timeShift | s.nodeID<<nodeShift | s.seq
}
