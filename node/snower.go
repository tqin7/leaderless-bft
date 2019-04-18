package proto

import (
	"math/rand"
	"context"
	pb "github.com/spockqin/leaderless-bft/proto"
	"errors"
	"github.com/spockqin/leaderless-bft/util"
	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"
	"github.com/spockqin/leaderless-bft/types"
	. "github.com/spockqin/leaderless-bft/types"
)

type Snower struct {
	Gossiper
	allIps []string //ip of every node in the network
	confidences ConfidenceMap //storing num of queries that yield each output
	ordered []string //an ordered lists of requests, i.e. ones with consensus
	proposal string //what's believed to be the next request
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
	if s.proposal == NO_PROPOSAL {
		s.setProposal()
		go s.performQueries()
	}

	proposalHash := util.HashBytes([]byte(s.proposal))
	return &pb.ReqId{Hash:proposalHash}, nil
}

func (s *Snower) SendReq(ctx context.Context, req *pb.ReqBody) (*pb.Void, error) {
	go s.Push(ctx, req)

	go s.performQueries()

	return &pb.Void{}, nil
}

func (s *Snower) setProposal() {
	s.requestsLock.Lock()
	defer s.requestsLock.Unlock()

	s.proposal = s.requests[rand.Intn(len(s.requests))] //random request for now
}

func (s *Snower) performQueries() {
	for i := 0; i < types.SNOWBALL_SAMPLE_ROUNDS; i++ {
		s.getMajorityVote()
	}
	s.ordered = append(s.ordered, s.proposal)
	s.confidences.RemoveKey(s.proposal)
	s.proposal = NO_PROPOSAL
}

//queries random subset, get majority, update confidence and proposal
func (s *Snower) getMajorityVote() string {
	sampleSize := 3
	networkSubset := util.UniqueRandomSample(s.allIps, sampleSize)
	votes := CreateConfidenceMap()
	for _, ip := range networkSubset {
		conn, err := grpc.Dial(tcpString(ip), grpc.WithInsecure())
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
		votes.IncreaseConfidence(string(vote.Hash))
	}

	winner, _ := votes.GetKeyWithMostConfidence()

	s.confidences.IncreaseConfidence(winner)
	if s.confidences.Get(winner) > s.confidences.Get(s.proposal) {
		s.proposal = winner
	}

	log.WithFields(log.Fields{
		"ip": s.ip,
		"majority": winner,
		"result": votes.kvMap,
	}).Info("Completed querying one network sample")
	return winner
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
	newSnower.proposal = NO_PROPOSAL

	return newSnower
}

