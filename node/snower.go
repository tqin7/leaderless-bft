package proto

import (
	"math/rand"
	"context"
	pb "github.com/spockqin/leaderless-bft/proto"
	"errors"
	"github.com/spockqin/leaderless-bft/util"
	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"
	"net"
	"bytes"
	"github.com/spockqin/leaderless-bft/types"
	"fmt"
	"sync"
)

type Snower struct {
	Gossiper
	allIps []string //ip of every node in the network
	confidences ConfidenceMap //storing num of queries that yield each output
	confidencesLock sync.Mutex
	ordered []string //an ordered lists of requests, i.e. ones with consensus
	proposal []byte //the hash of what's believed to be the next request
}

//return what this node thinks is the next request
func (s *Snower) GetVote(ctx context.Context, req *pb.Void) (*pb.ReqId, error) {
	s.requestsLock.Lock()
	reqLen := len(s.requests)
	s.requestsLock.Unlock()

	if reqLen == 0 {
		return &pb.ReqId{Hash:nil}, errors.New("no known request")
	}

	//if s doesn't have a proposal yet, set proposal and initiate query
	//otherwise, simply return proposal
	if s.proposal == nil {
		s.setProposal()
		go s.performQueries()
	}

	return &pb.ReqId{Hash:s.proposal}, nil
}

func (s *Snower) SendReq(ctx context.Context, req *pb.ReqBody) (*pb.Void, error) {
	s.Push(ctx, req)

	s.setProposal()
	go s.performQueries()

	return &pb.Void{}, nil
}

func (s *Snower) GetOrderedReqs(ctx context.Context, req *pb.Void) (*pb.Requests, error) {
	return &pb.Requests{Requests: s.ordered}, nil
}

func (s *Snower) setProposal() {
	s.requestsLock.Lock()
	randomRequest := s.requests[rand.Intn(len(s.requests))]
	s.requestsLock.Unlock()

	s.proposal = util.HashBytes([]byte(randomRequest)) //random request for now
}

func (s *Snower) performQueries() {
	for i := 0; i < types.SNOWBALL_SAMPLE_ROUNDS; i++ {
		s.getMajorityVote()
	}
	//most confident request and its index
	bestReq, bestIndex := s.findReqFromHash(s.proposal)
	fmt.Println("I'm most confident in: ", bestReq, s.confidences.Get(bestReq))
	log.WithFields(log.Fields{
		"ip": s.ip,
		"majority": bestReq,
		//"result": s.confidences,
	}).Info("Completed one entire query")
	s.ordered = append(s.ordered, bestReq)
	s.confidencesLock.Lock()
	fmt.Println("removing ", bestReq, s.ip)
	s.confidences.RemoveKey(bestReq)
	s.confidencesLock.Unlock()
	s.requests = util.RemoveOneFromArr(s.requests, bestIndex)
	s.proposal = nil
}

func (s *Snower) findReqFromHash(hash []byte) (string, int) {
	for i, request := range s.requests {
		if bytes.Compare(util.HashBytes([]byte(request)), hash) == 0 {
			return request, i
		}
	}
	return "", -1
}

//queries random subset, get majority, update confidence and proposal
func (s *Snower) getMajorityVote() string {
	sampleSize := 3
	networkSubset := util.UniqueRandomSample(s.allIps, sampleSize)
	votes := CreateConfidenceMap()
	for _, ip := range networkSubset {
		conn, err := grpc.Dial(ip, grpc.WithInsecure())
		if err != nil {
			log.WithFields(log.Fields{
				"ip": s.ip,
				"node": ip,
				"error": err,
			}).Error("Cannot dial node\n")
			continue
		}
		client := pb.NewSnowballClient(conn)
		vote, err := client.GetVote(context.Background(), &pb.Void{})
		conn.Close()
		if err != nil {
			log.WithFields(log.Fields{
				"ip": s.ip,
				"node": ip,
				"error": err,
			}).Error("Cannot get vote\n")
			continue
		}
		log.WithFields(log.Fields{
			"ip": s.ip,
			"from": ip,
		}).Info("Got vote")
		votes.IncreaseConfidence(string(vote.Hash))
	}

	majorityHashStr, _ := votes.GetKeyWithMostConfidence()
	majorityHash := []byte(majorityHashStr)

	s.confidencesLock.Lock()
	majorityReq, _ := s.findReqFromHash(majorityHash)
	s.confidences.IncreaseConfidence(majorityReq)
	if s.confidences.Get(majorityReq) > s.confidences.Get(string(s.proposal)) {
		s.proposal = majorityHash
	}
	s.confidencesLock.Unlock()

	return majorityHashStr
}

func CreateSnower(ip string, allIps []string) *Snower {
	newSnower := new(Snower)
	newSnower.ip = ip
	newSnower.peers = make([]string, 0)
	newSnower.hashes = make(map[string]bool)
	newSnower.requests = make([]string, 0)

	newSnower.allIps = allIps
	newSnower.confidences = CreateConfidenceMap()
	newSnower.ordered = make([]string, 0)
	newSnower.proposal = nil

	return newSnower
}

func (s *Snower) SnowerUp() {
	lis, err := net.Listen("tcp", s.ip)
	if err != nil {
		log.WithField("ip", s.ip).Error("Cannot listen on tcp [snower]")
	}

	grpcServer := grpc.NewServer()

	pb.RegisterGossipServer(grpcServer, s)
	pb.RegisterSnowballServer(grpcServer, s)

	grpcServer.Serve(lis)
}

