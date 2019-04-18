package proto

import "sync"

type ConfidenceMap struct {
	lock sync.Mutex
	kvMap map[string]uint
}

func CreateConfidenceMap() ConfidenceMap {
	m := ConfidenceMap{}
	m.kvMap = make(map[string]uint)
	return m
}

func (m *ConfidenceMap) Get(k string) uint {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.kvMap[k]
}

func (m *ConfidenceMap) RemoveKey(k string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.kvMap, k)
}

func (m *ConfidenceMap) IncreaseConfidence(k string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.kvMap[k] += 1
}

func (m *ConfidenceMap) GetKeyWithMostConfidence() (string, uint) {
	var maxConfidence uint
	var bestKey string

	m.lock.Lock()
	defer m.lock.Unlock()

	for k, v := range m.kvMap {
		if v >= maxConfidence {
			maxConfidence = v
			bestKey = k
		}
	}

	return bestKey, maxConfidence
}