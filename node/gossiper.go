package proto

import (
	"context"
	log "github.com/sirupsen/logrus"
	pb "github.com/spockqin/leaderless-bft/proto"
	"github.com/spockqin/leaderless-bft/types"
	"github.com/spockqin/leaderless-bft/util"
	"google.golang.org/grpc"
	"net"
	"sync"
)

// Structure of each gossiper
type Gossiper struct {
	ip string // IP of self

	peers []string // known peers
	peersLock sync.Mutex
	clients map[string]pb.GossipClient

	hashes map[string]bool // hashes of requests known
	hashesLock sync.Mutex
	poked map[string]bool // hashes of requests poked for
	pokedLock sync.Mutex
	requests []string // known requests
	requestsLock sync.Mutex
}

func (g *Gossiper) Poke(ctx context.Context, reqId *pb.ReqId) (*pb.Bool, error) {
	hashString := string(reqId.Hash)

	g.hashesLock.Lock()
	exists := g.hashes[hashString]
	g.hashesLock.Unlock()

	// check whether g is already poked for this request to
	// avoid double pushing caused by concurrency
	if !exists {
		g.pokedLock.Lock()
		exists = g.poked[hashString]
		g.poked[hashString] = true
		g.pokedLock.Unlock()
	}

	//log.WithFields(log.Fields{
	//	"ip": g.ip,
	//	"exists": exists,
	//}).Info("got poked")

	// log.Info("Timestamp: ",
	// 	time.Now().Format("2006-01-01 15:04:05 .000"))
	//fmt.Println("Testing Timestamp:", time.Now().Unix())

	return &pb.Bool{Status: exists}, nil
}

func (g *Gossiper) Push(ctx context.Context, reqBody *pb.ReqBody) (*pb.Void, error) {
	g.hashesLock.Lock()
	reqHash := util.HashBytes(reqBody.Body)
	g.hashes[string(reqHash)] = true
	g.hashesLock.Unlock()

	g.requestsLock.Lock()
	g.requests = append(g.requests, string(reqBody.Body))
	g.requestsLock.Unlock()

	g.pokedLock.Lock()
	delete(g.poked, string(reqHash))
	g.pokedLock.Unlock()

	// log.WithFields(log.Fields{
	// 	"ip": g.ip,
	// 	"request": string(reqBody.Body),
	// }).Info("stored new request")

	maxSoc := make(chan bool, types.MAX_SOCKETS)
	for _, peerIp := range g.peers {
		maxSoc <- true // blocks if maxSoc is full
		go g.sendGossip(peerIp, reqBody.Body, maxSoc)
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

func (g *Gossiper) sendGossip(neighborIp string, request []byte, c chan bool) {
	reqHash := util.HashBytes(request)

	// conn, err := grpc.Dial(neighborIp, grpc.WithInsecure())
	defer func(c chan bool) { <-c }(c)
	// if err != nil {
	// 	log.WithFields(log.Fields{
	// 		"ip": g.ip,
	// 		"peer": neighborIp,
	// 		"error": err,
	// 	}).Error("Cannot dial peer\n")
	// 	return
	// }
	// defer conn.Close()

	// client := pb.NewGossipClient(g.connections[neighborIndex])
	client := g.clients[neighborIp]

	// ctx, cancel := context.WithTimeout(context.Background(), types.GRPC_TIMEOUT)
	exists, err := client.Poke(context.Background(), &pb.ReqId{Hash: reqHash}, grpc.WaitForReady(true))
	// defer cancel()

	if err != nil {
		log.WithFields(log.Fields{
			"ip": g.ip,
			"peer": neighborIp,
			"error": err,
		}).Info("failed to poke peer\n")
	} else {
		//log.WithFields(log.Fields{
		//	"ip": g.ip,
		//	"peer": neighborIp,
		//	"exists": exists.Status,
		//}).Info("poked peer\n")
		if !exists.Status { // if peer doesn't have this hash
			client.Push(context.Background(), &pb.ReqBody{Body: request}, grpc.WaitForReady(true))
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

	conn, err := grpc.Dial(peerIp, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{
			"ip": g.ip,
			"peer": peerIp,
			"error": err,
		}).Error("Cannot dial peer\n")
		return
	} else {
		client := pb.NewGossipClient(conn)
		g.clients[peerIp] = client
	}

	log.WithFields(log.Fields{
		"ip": g.ip,
		"peer": peerIp,
	}).Info("Added and connected to peer\n")
}

func CreateGossiper(ip string) *Gossiper {
	newGossiper := new(Gossiper)
	newGossiper.ip = ip
	newGossiper.peers = make([]string, 0)
	newGossiper.clients = make(map[string]pb.GossipClient)
	newGossiper.hashes = make(map[string]bool)
	newGossiper.poked = make(map[string]bool)
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
