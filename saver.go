package event

import "sync"

type Saver struct {
	Type      string
	Listeners []*Listener
	sync.RWMutex
}
