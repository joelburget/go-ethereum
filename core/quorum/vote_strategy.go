package quorum

import (
	"math/rand"
	"sync"
	"time"

	"encoding/json"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

type BlockMakerStatus int

const (
	_                       = iota
	Active BlockMakerStatus = iota
	Paused                  = iota
)

type BlockMakerStrategy interface {
	// Start generating blocks
	Start() error
	// (temporary) stop generating blocks
	Pause() error
	// Resume after a pause
	Resume() error
	// Status returns indication if this implementation
	// is generation CreateBlock events.
	Status() BlockMakerStatus
}

// randomDeadlineStrategy asks the block voter to generate blocks
// after a deadline is passed without importing a new head. This
// deadline is chosen random between 2 limits.
type randomDeadlineStrategy struct {
	mux           *event.TypeMux
	min, max      int // min and max deadline
	activeMu      sync.Mutex
	active        bool
	deadlineTimer *time.Timer
}

// NewRandomDeadlineStrategy returns a block maker strategy that
// generated blocks randomly between the given min and max seconds.
func NewRandomDeadlineStrategy(mux *event.TypeMux, min, max uint) *randomDeadlineStrategy {
	if min > max {
		min, max = max, min
	}
	if min == 0 {
		log.Info("Set minimum block deadline interval to 1 second")
		min += 1
	}
	if max < min+1 {
		max = min + 1
		log.Info("Set maximum block deadline interval to ", max)
	}
	s := &randomDeadlineStrategy{
		mux:    mux,
		min:    int(min),
		max:    int(max),
		active: true,
	}
	return s
}

func resetBlockMakerTimer(t *time.Timer, min, max int) {
	t.Stop()
	select {
	case <-t.C:
	default:
	}
	t.Reset(time.Duration(min+rand.Intn(max-min)) * time.Second)
}

// Start generating block create request events.
func (s *randomDeadlineStrategy) Start() error {
	log.Info("Random deadline strategy configured with", "min", s.min, "max", s.max)
	s.deadlineTimer = time.NewTimer(time.Duration(s.min+rand.Intn(s.max-s.min)) * time.Second)
	go func() {
		sub := s.mux.Subscribe(core.ChainHeadEvent{})
		for {
			select {
			case <-s.deadlineTimer.C:
				s.activeMu.Lock()
				if s.active {
					s.mux.Post(CreateBlock{})
				}
				s.activeMu.Unlock()
				resetBlockMakerTimer(s.deadlineTimer, s.min, s.max)
			case <-sub.Chan():
				resetBlockMakerTimer(s.deadlineTimer, s.min, s.max)
			}
		}
	}()

	return nil
}

// Pause stops generating block create requests.
// Can be resumed with Resume.
func (s *randomDeadlineStrategy) Pause() error {
	s.activeMu.Lock()
	s.active = false
	s.activeMu.Unlock()
	return nil
}

// Resume if paused.
func (s *randomDeadlineStrategy) Resume() error {
	s.activeMu.Lock()
	s.active = true
	s.activeMu.Unlock()
	return nil
}

// Status returns an indication if this strategy is currently
// generating block create request.
func (s *randomDeadlineStrategy) Status() BlockMakerStatus {
	s.activeMu.Lock()
	defer s.activeMu.Unlock()

	if s.active {
		return Active
	}
	return Paused
}

func (s *randomDeadlineStrategy) MarshalJSON() ([]byte, error) {
	s.activeMu.Lock()
	defer s.activeMu.Unlock()

	status := "active"
	if !s.active {
		status = "paused"
	}

	return json.Marshal(map[string]interface{}{
		"type":         "deadline",
		"minblocktime": s.min,
		"maxblocktime": s.max,
		"status":       status,
	})
}
