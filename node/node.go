package proto
//
//import (
//	"sync"
//	. "github.com/spockqin/leaderless-bft/util"
//	log "github.com/sirupsen/logrus"
//	. "github.com/spockqin/leaderless-bft/types"
//	"math/rand"
//	"net"
//	"bytes"
//	"encoding/json"
//	"io"
//	"fmt"
//	"bufio"
//)
//
//// Structure of each gossiper/node
//type Node struct {
//	ip string // IP of self
//	peers []string // known peers
//	peersLock sync.Mutex
//
//	store NodeKVMap // maps each proposal to confidence value
//	value StoreKey
//
//	hashes map[string]bool
//	requests []string
//}
//
//func CreateNode(ip string) *Node {
//	newNode := new(Node)
//	newNode.ip = ip
//	newNode.peers = make([]string, 5)
//	newNode.store = CreateNodeKVMap()
//
//	newNode.logInfo("new node created", log.Fields{})
//
//	return newNode
//}
//
//func (n *Node) NodeUp() {
//	go n.NodeTcpUp()
//	n.NodeUdpUp()
//}
//
///* starts listening for TCP connections on its IP addr (TCP for client requests) */
//func (n *Node) NodeTcpUp() {
//	l, err := net.Listen(CONN_TYPE, tcpString(n.ip))
//	if err != nil {
//		n.logError("Error listening: " + err.Error(), log.Fields{})
//		return
//	}
//
//	defer l.Close()
//	n.logInfo("Listening for TCP connections on " + tcpString(n.ip), log.Fields{})
//
//	for {
//		conn, err := l.Accept()
//		if err != nil {
//			n.logError("Error accepting TCP connection: " + err.Error(), log.Fields{})
//			continue
//		}
//		n.logInfo("accepted TCP connection", log.Fields{
//			"remote": conn.RemoteAddr(),
//		})
//		go n.handleClientConnection(conn)
//	}
//}
//
///* starts listening for UDP packets on its IP addr (UDP for gossip exchange) */
//func (n *Node) NodeUdpUp() {
//	pkc, err := net.ListenPacket("udp", udpString(n.ip))
//	if err != nil {
//		n.logError("Cannot listen for UDP packets", log.Fields{
//			"error": err,
//		})
//	}
//
//	defer pkc.Close()
//	n.logInfo("Listening for UDP packets on " + udpString(n.ip), log.Fields{})
//
//	buffer := make([]byte, 1024)
//	for {
//		bytesRead, remoteAddr, err := pkc.ReadFrom(buffer)
//		if err != nil {
//			n.logError("Error reading UDP packet: " + err.Error(), log.Fields{})
//			continue
//		}
//
//		var msg Message
//		err = json.Unmarshal(buffer[:bytesRead], &msg)
//		if err != nil {
//			n.logError("failed to unmarshal message", log.Fields{
//				"remote": remoteAddr,
//			})
//			continue
//		}
//		n.logInfo("received UDP packet", log.Fields{
//			"remote": remoteAddr,
//			"message": msg,
//		})
//
//		if msg.Purpose == GOSSIP_REQUEST {
//			go n.handleGossipReq(msg.Value, remoteAddr, pkc) //send node's value
//		}
//
//	}
//}
//
////TODO: modify so that a connection, rather than a request, is handled - a for loop?
//func (n *Node) handleClientConnection(conn net.Conn) {
//	n.logInfo("Received client request", log.Fields{})
//
//	clientMsg, _ := bufio.NewReader(conn).ReadString('\n')
//
//	if n.value == "" {
//		n.value = StoreKey(clientMsg)
//		n.logInfo("Initialized value to " + clientMsg, log.Fields{})
//	}
//
//
//
//}
//
//func (n *Node) handleGossipReq(v StoreKey, remoteAddr net.Addr, pkc net.PacketConn) {
//	if n.value == "" {
//		n.value = v
//		n.logInfo("initialized value upon receiving query", log.Fields{
//			"value": v,
//		})
//	}
//	pkc.WriteTo([]byte(n.value), remoteAddr)
//	n.logInfo("sent gossip response", log.Fields{
//		"remote": remoteAddr,
//	})
//}
//
//func (n *Node) samplePeers() {
//	//fmt.Fprintf(conn, "node's value is " + string(n.value))
//
//	peerSubset := n.choosePeerSubset()
//	subsetInfo := CreateNodeKVMap()
//	n.logInfo("start peer sampling process", log.Fields{
//		"sample": peerSubset,
//	})
//	for _, peerIp := range peerSubset {
//		res := n.fetchPeerInfo(peerIp)  //TODO: use Goroutines here and use channel to gather info
//		subsetInfo.IncreaseConfidence(res)
//	}
//	peerMajority, confidence := subsetInfo.GetKeyWithMostConfidence()
//	n.logInfo("finished sampling peer subset", log.Fields{
//		//"majority": peerMajority,
//		"votes": confidence,
//	})
//
//	n.store.IncreaseConfidence(peerMajority)
//	n.updateValue()
//}
//
///* Obtain a peer's value using UDP protocol */
//func (n *Node) fetchPeerInfo(peerIp string) StoreKey {
//	conn, err := net.Dial("udp", udpString(peerIp))
//	if err != nil {
//		n.logError("failed in UDP dialing: " + err.Error(), log.Fields{
//			"peer": peerIp,
//		})
//	}
//	defer conn.Close()
//
//	m := Message{Purpose: GOSSIP_REQUEST, Value: n.value}
//	mBytes, err := json.Marshal(m)
//	if err != nil {
//		n.logError("failed to marshal Message", log.Fields{
//			"message": m,
//		})
//		return "" //if failed, is empty string correct to return?
//	}
//
//	conn.Write(mBytes)
//
//	buffer := make([]byte, 1024)
//	bytesRead, err := conn.Read(buffer)
//
//	return StoreKey(buffer[:bytesRead])
//}
//
//func (n *Node) sendData(msg interface{}, conn net.Conn) error {
//	buf, err := json.Marshal(msg)
//	if err != nil {
//		n.logWarn("Could not serialize message", log.Fields{
//			"error": err,
//			"message": msg,
//		})
//		return err
//	}
//
//	for len(buf) > 0 {
//		l, err := conn.Write(buf)
//		if err != nil {
//			n.logWarn("Could not send message", log.Fields{
//				"error": err,
//			})
//			return err
//		}
//		buf = buf[l:]
//	}
//
//	return nil
//}
//
//func (n *Node) rcvData(msg interface{}, conn net.Conn) error {
//	var buf bytes.Buffer
//	io.Copy(&buf, conn)
//	fmt.Println("received bf: ", buf)
//	err := json.Unmarshal(buf.Bytes(), msg)
//	if err != nil {
//		n.logWarn("Could not un-serialize message", log.Fields{
//			"error": err,
//		})
//		return err
//	}
//	return nil
//}
//
//func (n *Node) AddPeer(peerIp string) {
//	n.peersLock.Lock()
//	defer n.peersLock.Unlock()
//
//	if StringInArray(peerIp, n.peers) {
//		msg := "Peer already exists"
//		n.logWarn(msg, log.Fields{
//			"peer": peerIp,
//		})
//		return
//	}
//
//	n.peers = append(n.peers, peerIp)
//	n.logInfo("Added peer", log.Fields{
//		"peer": peerIp,
//	})
//}
//
///* Randomly selects a subset of peers (to gossip with) */
//func (n *Node) choosePeerSubset() []string {
//	n.peersLock.Lock()
//	defer n.peersLock.Unlock()
//
//	numPeers := len(n.peers)
//	if numPeers <= GOSSIP_SUBSET_SIZE {
//		n.logInfo("not enough peers to select from, returning all peers", log.Fields{
//			"subset size": GOSSIP_SUBSET_SIZE,
//			"num of peers": numPeers,
//		})
//		return n.peers
//	}
//
//	selected := make([]string, GOSSIP_SUBSET_SIZE)
//	for i := 0; i < GOSSIP_SUBSET_SIZE; i++ {
//		selected[i] = n.peers[rand.Intn(numPeers)]
//	}
//
//	return selected
//}
//
///* set a node's value to what it has the most confidence in */
//func (n *Node) updateValue() {
//	newValue, confidence := n.store.GetKeyWithMostConfidence()
//	n.value = newValue
//	fmt.Println(newValue)
//	fmt.Println("red")
//	n.logInfo("updated value", log.Fields{
//		"value": newValue,
//		"confidence": confidence,
//	})
//}
//
//func udpString(ip string) string {
//	return ip + ":" + CONN_UDP_PORT
//}
//
//
