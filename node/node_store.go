package proto

import "sync"

type StoreKey string
type Confidence uint

type NodeKVMap struct {
	lock sync.Mutex
	kvMap map[StoreKey]Confidence
}

func CreateNodeKVMap() NodeKVMap {
	m := NodeKVMap{}
	m.kvMap = make(map[StoreKey]Confidence)
	return m
}

/* Insert a new key, initialize confidence to 0 */
func (m *NodeKVMap) AddNewKey(k StoreKey) {
	m.lock.Lock()
	m.kvMap[k] = 0
	m.lock.Unlock()
}

func (m *NodeKVMap) RemoveKey(k StoreKey) {
	m.lock.Lock()
	delete(m.kvMap, k)
	m.lock.Unlock()
}

func (m *NodeKVMap) IncreaseConfidence(k StoreKey) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.kvMap[k]; ok { // key exists
		m.kvMap[k] += 1
	} else {
		m.kvMap[k] = 1
	}
}

func (m *NodeKVMap) GetKeyWithMostConfidence() (StoreKey, Confidence) {
	var maxConfidence Confidence
	var bestKey StoreKey

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