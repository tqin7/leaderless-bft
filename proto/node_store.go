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