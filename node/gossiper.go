package proto

import (
	"sync"
	"context"
	pb "github.com/spockqin/leaderless-bft/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spockqin/leaderless-bft/util"
	"net"
	"google.golang.org/grpc"
)

// Structure of each gossiper
type Gossiper struct {
	ip string // IP of self

	peers []string // known peers
	peersLock sync.Mutex

	hashes map[string]bool // hash of requests known
	requests []string // known requests
	requestsLock sync.Mutex
}

func (g *Gossiper) Poke(ctx context.Context, reqId *pb.ReqId) (*pb.Bool, error) {
	_, exists := g.hashes[string(reqId.Hash)]

	log.WithFields(log.Fields{
		"ip": g.ip,
		"exists": exists,
	}).Info("got poked")

	return &pb.Bool{Status: exists}, nil
}

func (g *Gossiper) Push(ctx context.Context, reqBody *pb.ReqBody) (*pb.Void, error) {
	g.requestsLock.Lock()
	g.requests = append(g.requests, string(reqBody.Body))
	g.requestsLock.Unlock()

	reqHash := util.HashBytes(reqBody.Body)
	g.hashes[string(reqHash)] = true

	log.WithFields(log.Fields{
		"ip": g.ip,
		"request": string(reqBody.Body),
	}).Info("stored new request")

	for _, peerIp := range g.peers {
		go g.sendGossip(peerIp, reqBody.Body)
	}

	return &pb.Void{}, nil
}

func (g *Gossiper) GetAllRequests(ctx context.Context, void *pb.Void) (*pb.Requests, error) {
	g.requestsLock.Lock()
	defer g.requestsLock.Unlock()

	log.WithFields(log.Fields{
		"ip": g.ip,
		"requests": g.requests,
	}).Info("print all known requests")
	return &pb.Requests{Requests: g.requests}, nil
}

func (g *Gossiper) sendGossip(neighborIp string, request []byte) {
	conn, err := grpc.Dial(neighborIp, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{
			"ip": g.ip,
			"peer": neighborIp,
			"error": err,
		}).Error("Cannot dial peer\n")
		return
	}
	defer conn.Close()

	client := pb.NewGossipClient(conn)

	reqHash := util.HashBytes(request)
	exists, err := client.Poke(context.Background(), &pb.ReqId{Hash: reqHash})

	if err != nil {
		log.WithFields(log.Fields{
			"ip": g.ip,
			"peer": neighborIp,
			"error": err,
		}).Info("failed to poke peer\n")
	} else {
		log.WithFields(log.Fields{
			"ip": g.ip,
			"peer": neighborIp,
			"exists": exists.Status,
		}).Info("poked peer\n")
		if !exists.Status { // if peer doesn't have this hash
			log.WithFields(log.Fields{
				"ip": g.ip,
				"peer": neighborIp,
				"request": string(request),
			}).Info("push request to peer\n")
			client.Push(context.Background(), &pb.ReqBody{Body: request})
		}
	}
}

func (g *Gossiper) AddPeer(peerIp string) {
	g.peersLock.Lock()
	defer g.peersLock.Unlock()

	if util.StringInArray(peerIp, g.peers) {
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

func (g *Gossiper) GossiperUp() {
	lis, err := net.Listen("tcp", g.ip)
	if err != nil {
		log.WithField("ip", g.ip).Error("Cannot listen on tcp")
	}

	grpcServer := grpc.NewServer()

	pb.RegisterGossipServer(grpcServer, g)

	grpcServer.Serve(lis)
}

//func tcpString(ip string) string {
//	return ip + ":" + CONN_TCP_PORT
//}
