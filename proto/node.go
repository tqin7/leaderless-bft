package proto

import (
	"sync"
	. "github.com/spockqin/leaderless-bft/util"
	log "github.com/sirupsen/logrus"
)

// Structure of each gossiper/node
type Node struct {
	ip string // IP of self
	peers []string // known peers
	peersLock sync.Mutex

	store NodeKVMap // maps each proposal to confidence value
	value StoreKey
}

func CreateNode(ip string) *Node {
	newNode := new(Node)
	newNode.ip = ip
	newNode.peers = make([]string, 5)
	newNode.store = CreateNodeKVMap()
}

func (n *Node) AddPeer(peerIp string) {
	n.peersLock.Lock()
	defer n.peersLock.Unlock()

	if StringInArray(peerIp, n.peers) {
		msg := "Peer already exists"
		n.logWarn(msg, log.Fields{
			"peer": peerIp,
		})
		return
	}

	n.peers = append(n.peers, peerIp)
	n.logInfo("Added peer", log.Fields{
		"peer": peerIp,
	})
}

/* set a node's value to what it has the most confidence in */
func (n *Node) updateValue() {
	newValue, confidence := n.store.GetKeyWithMostConfidence()
	n.value = newValue
	n.logInfo("Updated value", log.Fields{
		"value": newValue,
		"confidence": confidence,
	})
}

func (n *Node) logInfo(msg string, fields log.Fields) {
	log.WithField("node", n.ip).WithFields(fields).Info(msg)
}

func (n *Node) logWarn(msg string, fields log.Fields) {
	log.WithField("node", n.ip).WithFields(fields).Warn(msg)
}

func (n *Node) logError(msg string, fields log.Fields) {
	log.WithField("node", n.ip).WithFields(fields).Error(msg)
}
