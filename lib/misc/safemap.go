package misc

import "sync"

type SafeMap struct {
	sync.RWMutex
	Map map[int]int
}

func newSafeMap(size int) *SafeMap {
	sm := new(SafeMap)
	sm.Map = make(map[int]int)
	return sm

}

func (sm *SafeMap) readMap(key int) int {
	sm.RLock()
	value := sm.Map[key]
	sm.RUnlock()
	return value
}

func (sm *SafeMap) writeMap(key int, value int) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}
