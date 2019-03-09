package proto

import (
	"sync"
	. "github.com/spockqin/leaderless-bft/util"
	log "github.com/sirupsen/logrus"
	. "github.com/spockqin/leaderless-bft/types"
	"math/rand"
	"net"
	"bytes"
	"encoding/json"
	"io"
	"fmt"
	"bufio"
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

	newNode.logInfo("new node created", log.Fields{})

	return newNode
}

/* starts listening on its IP addr */
func (n *Node) NodeUp() {
	l, err := net.Listen(CONN_TYPE, n.ip + ":" + CONN_PORT)
	if err != nil {
		n.logError("Error listening: " + err.Error(), log.Fields{})
		return
	}

	defer l.Close()
	n.logInfo("Listening on " + n.ip + ":" + CONN_PORT, log.Fields{})
	for {
		conn, err := l.Accept()
		if err != nil {
			n.logError("Error accepting connection: " + err.Error(), log.Fields{})
			return
		}

		go n.handleClientConnection(conn)
	}
}

func (n *Node) handleClientConnection(conn net.Conn) {
	n.logInfo("Received client request", log.Fields{})

	//var req Message
	//n.rcvData(&req, conn)
	//n.logInfo("Parsed client request message", log.Fields{
	//	"message": req.value,
	//})

	clientMsg, _ := bufio.NewReader(conn).ReadString('\n')

	if n.value == "" {
		n.value = StoreKey(clientMsg)
		n.logInfo("Initialized value to " + clientMsg, log.Fields{})
	}

	fmt.Fprintf(conn, "node's value is " + string(n.value))

	peerSubset := n.choosePeerSubset()
	for _, peerIp := range peerSubset {
		//TODO: use UDP to ask each peer for their value and update confidence
	}


	//n.store.IncreaseConfidence()
}

func (n *Node) sendData(msg interface{}, conn net.Conn) error {
	buf, err := json.Marshal(msg)
	if err != nil {
		n.logWarn("Could not serialize message", log.Fields{
			"error": err,
			"message": msg,
		})
		return err
	}

	for len(buf) > 0 {
		l, err := conn.Write(buf)
		if err != nil {
			n.logWarn("Could not send message", log.Fields{
				"error": err,
			})
			return err
		}
		buf = buf[l:]
	}

	return nil
}

func (n *Node) rcvData(msg interface{}, conn net.Conn) error {
	var buf bytes.Buffer
	io.Copy(&buf, conn)
	fmt.Println("received bf: ", buf)
	err := json.Unmarshal(buf.Bytes(), msg)
	if err != nil {
		n.logWarn("Could not un-serialize message", log.Fields{
			"error": err,
		})
		return err
	}
	return nil
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

/* Randomly selects a subset of peers (to gossip with) */
func (n *Node) choosePeerSubset() []string {
	n.peersLock.Lock()
	defer n.peersLock.Unlock()

	numPeers := len(n.peers)
	if numPeers <= GOSSIP_SUBSET_SIZE {
		n.logInfo("not enough peers to select from, returning all peers", log.Fields{
			"subset size": GOSSIP_SUBSET_SIZE,
			"num of peers": numPeers,
		})
		return n.peers
	}

	selected := make([]string, GOSSIP_SUBSET_SIZE)
	for i := 0; i < GOSSIP_SUBSET_SIZE; i++ {
		selected[i] = n.peers[rand.Intn(numPeers)]
	}

	return selected
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
