package proto

import (
	// "math/rand"
	"context"
	pb "github.com/spockqin/leaderless-bft/proto"
	// "errors"
	"github.com/spockqin/leaderless-bft/util"
	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"
	"net"
	// "bytes"
	"github.com/spockqin/leaderless-bft/types"
	"fmt"
	"sync"
)

type Snower struct {
	Gossiper
	allIps []string // ip of every node in the network
	confidences ConfidenceMap // storing num of queries that yield each output
	confidencesLock sync.Mutex // lock on above confidences map
	seqNum uint64 // current sequence number, non-decreasing
	finalSeqNums map[string]uint64 // map request hashes (stringified) to their final sequence numbers
}

/* return this node's opinion on sequence number */
func (s *Snower) GetVote(ctx context.Context, msg *pb.SeqNumMsg) (*pb.SeqNumMsg, error) {
	// return existing seq num if it's determined already
	if seqNum, exists := s.finalSeqNums[string(msg.ReqHash)]; exists {
		return &pb.SeqNumMsg{SeqNum: seqNum, ReqHash: msg.ReqHash}, nil
	} else {
		s.seqNum = msg.SeqNum
		go s.performQueries(msg)
	}

	return &pb.SeqNumMsg{SeqNum: s.seqNum, ReqHash: msg.ReqHash}, nil
}

func (s *Snower) SendReq(ctx context.Context, req *pb.ReqBody) (*pb.Void, error) {
	// check if request is already known
	reqHash := util.HashBytes(req.Body)

	if !s.hashes[string(reqHash)] {
		s.Push(ctx, req)

		s.seqNum += 1
		seqNumProposal := pb.SeqNumMsg{SeqNum: s.seqNum, ReqHash: reqHash}
		go s.performQueries(&seqNumProposal)
	}

	return &pb.Void{}, nil
}

func (s *Snower) performQueries(msg *pb.SeqNumMsg) {
	for i := 0; i < types.SNOWBALL_SAMPLE_ROUNDS; i++ {
		go s.getMajorityVote(msg)
	}

	fmt.Println("I'm most confident in: ", s.seqNum, "with votes ", s.confidences.Get(s.seqNum))
	log.WithFields(log.Fields{
		"ip": s.ip,
		"majority": s.seqNum,
	}).Info("Completed one entire query")

	s.finalSeqNums[string(msg.ReqHash)] = s.seqNum
	// s.confidences = CreateConfidenceMap()
}

//queries random subset, get majority, update confidence and proposal
func (s *Snower) getMajorityVote(msg *pb.SeqNumMsg) uint64 {
	sampleSize := 2		//TODO: determine sample size as sqrt?
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
		vote, err := client.GetVote(context.Background(), msg)
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
		votes.IncreaseConfidence(vote.SeqNum)
	}

	majorityValue, _ := votes.GetKeyWithMostConfidence()

	s.confidencesLock.Lock()
	s.confidences.IncreaseConfidence(majorityValue)
	if s.confidences.Get(majorityValue) > s.confidences.Get(s.seqNum) {
		s.seqNum = majorityValue
	}
	s.confidencesLock.Unlock()

	return majorityValue
}

func (s *Snower) completeRequest(reqId *pb.ReqId) {
	delete(s.finalSeqNums, string(reqId.Hash))
	delete(s.hashes, string(reqId.Hash))
	//TODO: remove request from s.requests
}

func CreateSnower(ip string, allIps []string) *Snower {
	newSnower := new(Snower)
	newSnower.ip = ip
	newSnower.peers = make([]string, 0)
	newSnower.hashes = make(map[string]bool)
	newSnower.poked = make(map[string]bool)
	newSnower.requests = make([]string, 0)

	newSnower.allIps = allIps
	newSnower.confidences = CreateConfidenceMap()
	newSnower.seqNum = 0
	newSnower.finalSeqNums = make(map[string]uint64)

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

