package proto

import (
	"sync"
	"context"
	pb "github.com/spockqin/leaderless-bft/proto"
	. "github.com/spockqin/leaderless-bft/types"
	log "github.com/sirupsen/logrus"
	. "github.com/spockqin/leaderless-bft/util"
	"crypto/sha256"
	"net"
	"google.golang.org/grpc"
	"math/rand"
	"fmt"
	"time"
)

// Structure of each gossiper
type Gossiper struct {
	ip string // IP of self

	peers []string // known peers
	peersLock sync.Mutex

	hashes map[string]bool

	requests []string
	requestsLock sync.Mutex
}

func (g *Gossiper) Poke(ctx context.Context, reqId *pb.ReqId) (*pb.Bool, error) {
	fmt.Println("poke request received")
	_, exists := g.hashes[string(reqId.Hash)]

	log.WithFields(log.Fields{
		"ip": g.ip,
		"exists?": exists,
	}).Info("got poked")

	return &pb.Bool{Status: exists}, nil
}

func (g *Gossiper) Push(ctx context.Context, reqBody *pb.ReqBody) (*pb.Void, error) {
	g.requestsLock.Lock()
	g.requests = append(g.requests, string(reqBody.Body))
	g.requestsLock.Unlock()

	reqHash := hashBytes(reqBody.Body)
	g.hashes[string(reqHash)] = true

	log.WithFields(log.Fields{
		"ip": g.ip,
		"request": string(reqBody.Body),
	}).Info("stored new request")

	return &pb.Void{}, nil
}

func (g *Gossiper) StartGossip() {
	log.WithField("ip", g.ip).Info("Start to gossip")
	gossipTicker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-gossipTicker.C:
			g.RandomGossip()
		}
	}
}

func (g *Gossiper) RandomGossip() {
	//share request with one random neighbor
	randomNeighbor := g.selectRandomPeer()
	randomRequest := g.selectRandomRequest()
	if randomNeighbor == "" || randomRequest == nil {
		return
	}

	conn, err := grpc.Dial(tcpString(randomNeighbor), grpc.WithInsecure()) // this should be unreliable
	if err != nil {
		log.WithFields(log.Fields{
			"ip": g.ip,
			"peer": randomNeighbor,
		}).Error("Cannot dial peer\n")
		return
	}
	defer conn.Close()

	client := pb.NewGossipClient(conn)


	randomReqHash := hashBytes(randomRequest)
	exists, err := client.Poke(context.Background(), &pb.ReqId{Hash: randomReqHash})

	if err != nil {
		log.WithFields(log.Fields{
			"ip": g.ip,
			"peer": randomNeighbor,
			"error": err,
		}).Info("failed to poke peer\n")
	} else {
		log.WithFields(log.Fields{
			"ip": g.ip,
			"peer": randomNeighbor,
			"exists": exists.Status,
		}).Info("poked peer\n")
		if !exists.Status { // if peer doesn't have this hash
			client.Push(context.Background(), &pb.ReqBody{Body: randomRequest})
			log.WithFields(log.Fields{
				"ip": g.ip,
				"peer": randomNeighbor,
				"request": randomRequest,
			}).Info("pushed request to peer\n")
		}
	}
}

func (g *Gossiper) selectRandomPeer() string {
	g.peersLock.Lock()
	defer g.peersLock.Unlock()
	if len(g.peers) == 0 {
		return ""
	}
	return g.peers[rand.Intn(len(g.peers))]
}

func (g *Gossiper) selectRandomRequest() []byte {
	g.requestsLock.Lock()
	defer g.requestsLock.Unlock()
	if len(g.requests) == 0 {
		return nil
	}
	return []byte(g.requests[rand.Intn(len(g.requests))])
}

func (g *Gossiper) AddPeer(peerIp string) {
	g.peersLock.Lock()
	defer g.peersLock.Unlock()

	if StringInArray(peerIp, g.peers) {
		log.WithFields(log.Fields{
			"ip": g.ip,
			"peer": peerIp,
		}).Warn("Peer already exists")
		return
	}

	g.peers = append(g.peers, peerIp)
	log.WithFields(log.Fields{
		"ip": g.ip,
		"peer": peerIp,
	}).Info("Added peer\n")
}

func CreateGossiper(ip string) *Gossiper {
	newGossiper := new(Gossiper)
	newGossiper.ip = ip
	newGossiper.peers = make([]string, 0)
	newGossiper.hashes = make(map[string]bool)
	newGossiper.requests = make([]string, 0)

	return newGossiper
}

func GossiperUp(g *Gossiper) {
	lis, err := net.Listen("udp", tcpString(g.ip))
	if err != nil {
		log.WithField("ip", g.ip).Error("Cannot listen on tcp")
	}

	grpcServer := grpc.NewServer()

	pb.RegisterGossipServer(grpcServer, g)

	grpcServer.Serve(lis)
}

func tcpString(ip string) string {
	return ip + ":" + CONN_TCP_PORT
}

func hashBytes(bytes []byte) []byte {
	h := sha256.New()
	h.Write(bytes)
	return h.Sum(nil)
}
