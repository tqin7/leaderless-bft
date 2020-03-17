package proto

import (
	// "math/rand"
	"context"
	pb "github.com/spockqin/leaderless-bft/proto"
	// "errors"
	"github.com/spockqin/leaderless-bft/util"
	grpc "google.golang.org/grpc"
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
	seqNum int64 // current sequence number, non-decreasing
	finalSeqNums map[string]int64 // map request hashes (stringified) to their final sequence numbers
	finalSeqNumsLock sync.Mutex
}

/* return this node's opinion on sequence number */
func (s *Snower) GetVote(ctx context.Context, msg *pb.SeqNumMsg) (*pb.SeqNumMsg, error) {
	// return existing seq num if it's determined already
	s.finalSeqNumsLock.Lock()
	seqNum, exists := s.finalSeqNums[string(msg.ReqHash)]
	s.finalSeqNumsLock.Unlock()
	if exists {
		return &pb.SeqNumMsg{SeqNum: seqNum, ReqHash: msg.ReqHash}, nil
	} else {
		s.seqNum = msg.SeqNum
		s.performQueries(msg)
		return &pb.SeqNumMsg{SeqNum: s.seqNum, ReqHash: msg.ReqHash}, nil
	}
}

func (s *Snower) SendReq(ctx context.Context, req *pb.ReqBody) (*pb.Void, error) {
	// check if request is already known
	reqHash := util.HashBytes(req.Body)

	if !s.hashes[string(reqHash)] {
		s.Push(ctx, req)

		s.seqNum += 1
		seqNumProposal := pb.SeqNumMsg{SeqNum: s.seqNum, ReqHash: reqHash}
		s.performQueries(&seqNumProposal)
	}

	return &pb.Void{}, nil
}

func (s *Snower) performQueries(msg *pb.SeqNumMsg) {
	var wg sync.WaitGroup

	for i := 0; i < types.SNOWBALL_SAMPLE_ROUNDS; i++ {
		wg.Add(1)
		go s.getMajorityVote(msg, &wg)
	}

	wg.Wait()

	if s.confidences.Size() == 0 {
		return
	}

	finalSeqNum, _ := s.confidences.GetKeyWithMostConfidence()

	log.WithFields(log.Fields{
		"ip": s.ip,
		"majority": finalSeqNum,
	}).Info("Completed entire query")

	s.finalSeqNumsLock.Lock()
	s.finalSeqNums[string(msg.ReqHash)] = finalSeqNum
	s.finalSeqNumsLock.Unlock()

	s.confidences = CreateConfidenceMap()
	s.seqNum = finalSeqNum
}

//queries random subset, get majority, update confidence and proposal
func (s *Snower) getMajorityVote(msg *pb.SeqNumMsg, wg *sync.WaitGroup) {
	defer wg.Done()

	sampleSize := 4		//TODO: determine sample size as sqrt?
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
		//CallOption for above line:    , grpc.WaitForReady(true)
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
		fmt.Println("vote is: ", vote.SeqNum)
		votes.IncreaseConfidence(vote.SeqNum)
	}

	if votes.Size() != 0 {
		fmt.Println("votes map: ", votes.kvMap)
		majorityValue, _ := votes.GetKeyWithMostConfidence()
		s.confidences.IncreaseConfidence(majorityValue)
	}
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
	newSnower.seqNum = -1
	newSnower.finalSeqNums = make(map[string]int64)

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

