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
)

// Structure of each gossiper
type Gossiper struct {
	ip string // IP of self
	peers []string // known peers
	peersLock sync.Mutex

	hashes map[string]bool
	requests []string
}

func (g *Gossiper) Poke(ctx context.Context, reqId *pb.ReqId) (*pb.Bool, error) {
	_, exists := g.hashes[string(reqId.Hash)]
	return &pb.Bool{Status: exists}, nil
}

func (g *Gossiper) Push(ctx context.Context, reqBody *pb.ReqBody) (*pb.Void, error) {
	g.requests = append(g.requests, string(reqBody.Body))
	reqHash := hashBytes(reqBody.Body)
	g.hashes[string(reqHash)] = true

	//share request with one random neighbor
	g.peersLock.Lock()
	randomNeighbor := g.peers[rand.Intn(len(g.peers))]
	g.peersLock.Unlock()

	conn, err := grpc.Dial(randomNeighbor)
	if err != nil {
		log.WithFields(log.Fields{
			"ip": g.ip,
			"peer": randomNeighbor,
		}).Error("Cannot dial peer")
	}
	defer conn.Close()

	client := pb.NewGossipClient(conn)
	exists, err := client.Poke(context.Background(), &pb.ReqId{Hash: reqHash})
	if !exists.Status { // if peer doesn't have this hash
		client.Push(context.Background(), reqBody)
	}
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
	}).Info("Added peer")
}

func CreateGossiper(ip string) *Gossiper {
	newGossiper := new(Gossiper)
	newGossiper.ip = ip
	newGossiper.peers = make([]string, 5)
	newGossiper.hashes = make(map[string]bool)
	newGossiper.requests = make([]string, 5)

	return newGossiper
}

func GossiperUp(g *Gossiper) {
	lis, err := net.Listen("tcp", tcpString(g.ip))
	if err != nil {
		log.WithField("ip", g.ip).Error("cannot listen on tcp")
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

//func (g *Gossiper) logInfo(msg string, fields log.Fields) {
//	log.WithField("ip", g.ip).WithFields(fields).Info(msg)
//}
//
//func (g *Gossiper) logWarn(msg string, fields log.Fields) {
//	log.WithField("ip", g.ip).WithFields(fields).Warn(msg)
//}
//
//func (g *Gossiper) logError(msg string, fields log.Fields) {
//	log.WithField("ip", g.ip).WithFields(fields).Error(msg)
//}
