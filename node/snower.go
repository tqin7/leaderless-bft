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
)

type Snower struct {
	Gossiper
	allIps []string
	confidences ConfidenceMap
	ordered []string
}

//return what this node thinks is the next request
func (s *Snower) GetVote(ctx context.Context, req *pb.Void) (*pb.ReqId, error) {
	s.requestsLock.Lock()
	reqLen := len(s.requests)
	s.requestsLock.Unlock()

	if reqLen == 0 {
		return &pb.ReqId{Hash:nil}, errors.New("no known request")
	}

	//select what this node believes is the next req. random for now
	s.requestsLock.Lock()
	randomReq := s.requests[rand.Intn(reqLen)]
	s.requestsLock.Unlock()

	reqHash := util.HashBytes([]byte(randomReq))
	return &pb.ReqId{Hash:reqHash}, nil
}

func (s *Snower) SendReq(ctx context.Context, req *pb.ReqBody) (*pb.Void, error) {
	s.Push(ctx, req)

	for i := 0; i < types.SNOWBALL_SAMPLE_ROUNDS; i++ {
		s.getMajorityVote()
	}
}

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

	return newSnower
}

