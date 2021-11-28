package models

import (
	"sync/atomic"
)

type ServiceInformation struct {
	BackEnds []*Backend
	Routes   []Route
	current  uint64
}

func (s *ServiceInformation) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.BackEnds)))
}

// GetNextPeer todo: fix here to be in a round robin way not base on alive stuff
func (s *ServiceInformation) GetNextPeer() *Backend {
	// loop entire backends to find out an Alive backend
	next := s.NextIndex()
	l := len(s.BackEnds) + next // start from next and move a full cycle
	for i := next; i < l; i++ {
		idx := i % len(s.BackEnds) // take an index by modding with length
		// if we have an alive backend, use it and store if its not the original one
		if s.BackEnds[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx)) // mark the current one
			}
			return s.BackEnds[idx]
		}
	}
	return nil
}
