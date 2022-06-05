package main

import "sync"

type RequestCount struct {
	count uint64
	lock  sync.RWMutex
}

func (rc *RequestCount) Inc() {
	rc.lock.Lock()
	defer rc.lock.Unlock()
	rc.count++
}

func (rc *RequestCount) Value() uint64 {
	rc.lock.RLock()
	defer rc.lock.RUnlock()
	return rc.count
}
