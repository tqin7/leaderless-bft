package proto

import "sync"

type ConfidenceMap struct {
	lock sync.Mutex
	kvMap map[int64]uint
}

func CreateConfidenceMap() ConfidenceMap {
	m := ConfidenceMap{}
	m.kvMap = make(map[int64]uint)
	return m
}

func (m *ConfidenceMap) Size() int {
	return len(m.kvMap)
}

func (m *ConfidenceMap) Get(k int64) uint {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.kvMap[k]
}

func (m *ConfidenceMap) RemoveKey(k int64) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.kvMap, k)
}

func (m *ConfidenceMap) IncreaseConfidence(k int64) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.kvMap[k] += 1
}

func (m *ConfidenceMap) GetKeyWithMostConfidence() (int64, uint) {
	var maxConfidence uint
	var bestKey int64

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