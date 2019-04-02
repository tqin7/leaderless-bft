package proto

import (
	"sync"
	log "github.com/sirupsen/logrus"
	"context"
	pb "github.com/spockqin/leaderless-bft/proto"
	"crypto/sha256"
)

// Structure of each gossiper
type Gossiper struct {
	ip string // IP of self
	peers []string // known peers
	peersLock sync.Mutex

	hashes map[*pb.ReqId]*pb.Bool
	requests []*pb.ReqBody
}

func CreateGossiper(ip string) *Gossiper {
	newGossiper := new(Gossiper)
	newGossiper.ip = ip
	newGossiper.peers = make([]string, 5)
	newGossiper.hashes = make(map[*pb.ReqId]*pb.Bool)
	newGossiper.requests = make([]*pb.ReqBody, 5)

	return newGossiper
}

func (g *Gossiper) Poke(ctx context.Context, reqId *pb.ReqId) (*pb.Bool, error) {
	if val, exists := g.hashes[reqId]; exists {
		return val, nil
	} else {
		return &pb.Bool{Status: false}, nil
	}
}

func (g *Gossiper) Push(ctx context.Context, reqInfo *pb.ReqInfo) (*pb.Void, error) {
	g.requests = append(g.requests, reqInfo.Body)
	g.hashes[reqInfo.Id] = &pb.Bool{Status: true}
	//hash := hashBytes(reqBody.Body)
	//g.hashes[&pb.ReqId{Hash: hash}] = &pb.Bool{Status: true}
	//TODO: push req to one neighbor
}

func hashBytes(bytes []byte) []byte {
	h := sha256.New()
	h.Write(bytes)
	return h.Sum(nil)
}
