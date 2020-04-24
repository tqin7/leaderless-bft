package proto

import (
	"context"
	pb "github.com/spockqin/leaderless-bft/proto"
	"github.com/spockqin/leaderless-bft/util"
	grpc "google.golang.org/grpc"
	log "github.com/sirupsen/logrus"
	"net"
	"github.com/spockqin/leaderless-bft/types"
	"fmt"
	"sync"
	"time"
	"math"
)

type Snower struct {
	Gossiper
	allIps []string // ip of every node in the network
	//confidences ConfidenceMap // storing num of queries that yield each output
	confidences map[string]ConfidenceMap
	seqNum int64 // current sequence number, non-decreasing
	finalSeqNums map[string]int64 // map request hashes (stringified) to their final sequence numbers
	finalSeqNumsLock sync.Mutex
}

/* return this node's opinion on sequence number */
func (s *Snower) GetVote(ctx context.Context, msg *pb.SeqNumMsg) (*pb.SeqNumMsg, error) {
	reqHashStr := string(msg.ReqHash)

	// return existing seq num if it's determined already
	s.finalSeqNumsLock.Lock()
	seqNum, exists := s.finalSeqNums[reqHashStr]
	s.finalSeqNumsLock.Unlock()
	if exists {
		return &pb.SeqNumMsg{SeqNum: seqNum, ReqHash: msg.ReqHash}, nil
	} else {
		//TODO: check whether msg.SeqNum is valid
		// if valid, adopt the seq num
		s.seqNum = msg.SeqNum
		s.finalSeqNumsLock.Lock()
		s.finalSeqNums[reqHashStr] = msg.SeqNum
		s.finalSeqNumsLock.Unlock()

		s.performQueries(reqHashStr, msg)
		return &pb.SeqNumMsg{SeqNum: msg.SeqNum, ReqHash: msg.ReqHash}, nil
	}
}

func (s *Snower) SendReq(ctx context.Context, req *pb.ReqBody) (*pb.Void, error) {
	// check if request is already known
	reqHash := util.HashBytes(req.Body)
	reqHashStr := string(reqHash)

	if !s.hashes[reqHashStr] {
		// s.Push(ctx, req)

		s.seqNum += 1
		seqNumProposal := pb.SeqNumMsg{SeqNum: s.seqNum, ReqHash: reqHash}
		s.performQueries(reqHashStr, &seqNumProposal)
	}

	return &pb.Void{}, nil
}

func (s *Snower) performQueries(reqHashStr string, msg *pb.SeqNumMsg) {
	s.confidences[reqHashStr] = CreateConfidenceMap()

	var wg sync.WaitGroup

	for i := 0; i < types.SNOWBALL_SAMPLE_ROUNDS; i++ {
		wg.Add(1)
		go s.getMajorityVote(reqHashStr, msg, &wg)
	}

	wg.Wait()

	reqConfidenceMap := s.confidences[reqHashStr]
	if reqConfidenceMap.Size() == 0 {
		log.WithFields(log.Fields{
			"ip": s.ip,
		}).Info("Not able to obtain any snowball vote")
		return
	}

	finalSeqNum, _ := reqConfidenceMap.GetKeyWithMostConfidence()

	log.WithFields(log.Fields{
		"ip": s.ip,
		"majority": finalSeqNum,
	}).Info("Completed entire snowball query")

	s.finalSeqNumsLock.Lock()
	s.finalSeqNums[reqHashStr] = finalSeqNum
	s.finalSeqNumsLock.Unlock()

	fmt.Println("Testing Timestamp:", time.Now().Unix())
}

//queries random subset, get majority, update confidence and proposal
func (s *Snower) getMajorityVote(reqHashStr string, msg *pb.SeqNumMsg, wg *sync.WaitGroup) {
	defer wg.Done()

	sampleSize := int(math.Floor(math.Sqrt(float64(len(s.allIps)))))
	alpha := uint(math.Floor(float64(sampleSize / 2)))

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
		//CallOption "grpc.WaitForReady(true)" elminates "TransientFailure error" but gets stuck
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
		// fmt.Println("vote is: ", vote.SeqNum)
		votes.IncreaseConfidence(vote.SeqNum)
	}

	if votes.Size() != 0 {
		// fmt.Println("votes map: ", votes.kvMap)
		majorityValue, majorityNumVotes := votes.GetKeyWithMostConfidence()
		if majorityNumVotes >= alpha {
			reqConfidenceMap := s.confidences[reqHashStr]
			reqConfidenceMap.IncreaseConfidence(majorityValue)
		}
	}
}

func (s *Snower) completeRequest(reqId *pb.ReqId) {
	delete(s.finalSeqNums, string(reqId.Hash))
	delete(s.hashes, string(reqId.Hash))
	delete(s.confidences, string(reqId.Hash))
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
	newSnower.confidences = make(map[string]ConfidenceMap)
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

