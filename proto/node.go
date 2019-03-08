package proto

import (
	"sync"
	"github.com/spockqin/leaderless-bft/util"
)

// Structure of each gossiper/node
type Node struct {
	ip string // IP of self
	peers []string // known peers
	peersLock sync.Mutex

	store NodeKVMap
}

func CreateNode(ip string) *Node {
	newNode := new(Node)
	newNode.ip = ip
	newNode.peers = make([]string, 5)
	newNode.store = CreateNodeKVMap()
}

func (n *Node) AddPeer(peerIp string) error {
	if util.StringInArray(peerIp, n.peers) {
		return util.LogAndGetError("Peer already exists")
	}

	n.peers = append(n.peers, peerIp)
	return nil
}
