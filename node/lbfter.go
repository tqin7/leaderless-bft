package proto

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	pb "github.com/spockqin/leaderless-bft/proto"
	"google.golang.org/grpc"
	//"time"
	// "google.golang.org/grpc/reflection"
	"net"
	"sync"
	tp "github.com/spockqin/leaderless-bft/types"
	util "github.com/spockqin/leaderless-bft/util"
)

type Lbfter struct {
	// Gossiper
	Snower
	Pbfter
}

func (l *Lbfter) LSendReq(ctx context.Context, request *pb.ReqBody) (*pb.Void, error) {
	var req tp.PbftReq
	err := json.Unmarshal(request.GetBody(), &req)
	if err != nil {
		log.Error("Error when unmarshal request [LSendReq]")
	}

	LogMsg(req)

	err, seqID := l.createLbftStateForNewConsensus(&req,-1)
	if err != nil {
		log.Error("Error when create consensus state")
		return &pb.Void{}, err
	}

	prePrepareMsg, err := l.Pbfter.CurrentState[seqID].StartConsensus(&req, seqID)
	if err != nil {
		log.Error("Error when starting consensus [LSendReq]")
	}

	// broadcast prePrepareMsg
	prePrepareMsgBytes, err := json.Marshal(*prePrepareMsg)
	if err != nil {
		return &pb.Void{}, errors.New("[LSendReq] prePrepareMsg marshal error!")
	} else {
		fmt.Println("l Pbfter Push")
		_, pushErr := l.Pbfter.Push(ctx, &pb.ReqBody{Body:prePrepareMsgBytes})
		fmt.Println("Pushed")
		if pushErr != nil {
			panic(errors.New("[LSendReq] push prePrepareMsgBytes error!"))
		}
	}

	LogStage("Pre-prepare", true, l.NodeID)

	return &pb.Void{}, nil
}

func (l *Lbfter) runSeqIdConsensus(operation string) {
	l.Snower.SendReq(context.Background(), &pb.ReqBody{Body: []byte(operation)}) // snowball
}

func (l *Lbfter) createLbftStateForNewConsensus(req *tp.PbftReq, msgSeqID int64) (error, int64) {
	var lastSeqID int64
	if len(l.Pbfter.CommittedMsgs) == 0 {
		lastSeqID = -1
	} else {
		lastSeqID = l.Pbfter.CommittedMsgs[len(l.Pbfter.CommittedMsgs) - 1].SequenceID
	}

	var seqID int64
	if msgSeqID == -1 {
		reqHash := string(util.HashBytes([]byte(req.Operation)))
		l.runSeqIdConsensus(req.Operation)
		seqID = l.Snower.finalSeqNums[string(reqHash)]
	} else {
		seqID = msgSeqID
	}

	l.Pbfter.CurrentState[seqID] = createState(l.Pbfter.ViewID, lastSeqID)

	return nil, seqID
}

//func (l *Lbfter) startLbftConsensus(req *tp.PbftReq, seqID int64) (*tp.PrePrepareMsg, error) {
//	reqHash := string(util.HashBytes([]byte(req.Operation)))
//	l.runSeqIdConsensus(req.Operation)
//	req.SequenceID = l.Snower.finalSeqNums[string(reqHash)]
//
//	state := l.Pbfter.CurrentState
//
//	state.MsgLogs.ReqMsg = req
//
//	digest, err := util.Digest(req)
//	if err != nil {
//		log.Error("Error happens when getting digest of request [StartLbftConsensus]")
//		panic(err)
//	}
//
//	state.CurrentStage = PrePrepared
//
//	return &tp.PrePrepareMsg{
//		ViewID:               state.ViewID,
//		SequenceID:           l.Snower.finalSeqNums[string(reqHash)],
//		Digest:               digest,
//		Req:                  req,
//		MsgType:			  "PrePrepareMsg",
//	}, nil
//}

func CreateLbfter(nodeID string, viewID int64, ip string, allIps []string) *Lbfter {
	newLbfter := &Lbfter{
		// Gossiper:      Gossiper{
		// 	ip:           ip,
		// 	peers:        make([]string, 0),
		// 	peersLock:    sync.Mutex{},
		// 	hashes:       make(map[string]bool),
		// 	hashesLock:   sync.Mutex{},
		// 	poked:        make(map[string]bool),
		// 	pokedLock:    sync.Mutex{},
		// 	requests:     make([]string, 0),
		// 	requestsLock: sync.Mutex{},
		// },

		Snower:		   Snower{
			allIps:			  allIps,
			confidences:	  make(map[string]ConfidenceMap),
			seqNum: 		  1,
			finalSeqNums:     make(map[string]int64),
			finalSeqNumsLock:  sync.Mutex{},
			Gossiper: Gossiper{
				ip:           ip,
				peers:        make([]string, 0),
				peersLock:    sync.Mutex{},
				hashes:       make(map[string]bool),
				hashesLock:   sync.Mutex{},
				poked:        make(map[string]bool),
				pokedLock:    sync.Mutex{},
				requests:     make([]string, 0),
				requestsLock: sync.Mutex{},
			},
		},


		Pbfter:        Pbfter{
			NodeID:        nodeID,
			ViewID:        viewID,
			CurrentState:  make(map[int64]*State),
			CommittedMsgs: make([]*tp.PbftReq, 0),
			MsgBuffer:     &MsgBuffer{
				ReqMsgs:        make([]*tp.PbftReq, 0),
				PrePrepareMsgs: make([]*tp.PrePrepareMsg, 0),
				PrepareMsgs:    make([]*tp.PrepareMsg, 0),
				CommitMsgs:     make([]*tp.CommitMsg, 0),
			},
			MsgDelivery:   make(chan interface{}),
			Gossiper: Gossiper{
				ip:           ip,
				peers:        make([]string, 0),
				peersLock:    sync.Mutex{},
				hashes:       make(map[string]bool),
				hashesLock:   sync.Mutex{},
				poked:        make(map[string]bool),
				pokedLock:    sync.Mutex{},
				requests:     make([]string, 0),
				requestsLock: sync.Mutex{},
			},
		},
	}

	return newLbfter
}

func (l *Lbfter) LbfterUp() {
	nodeIp := l.Snower.ip
	lis, err := net.Listen("tcp", nodeIp)
	if err != nil {
		log.WithField("ip", nodeIp).Error("Cannot listen on tcp [lbfter]")
		panic(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterGossipServer(grpcServer, l)
	pb.RegisterSnowballServer(grpcServer, l)
	pb.RegisterPbftServer(grpcServer, l)
	pb.RegisterLbftServer(grpcServer, l)
	// reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.WithField("ip", nodeIp).Error("Cannot serve [lbfter]")
		panic(err)
	}
}


/* ******************************************
    Override other gRPC service methods to 
    direct method calls to correct embedding
   ****************************************** */

/* Gossip */
func (l *Lbfter) Poke(ctx context.Context, reqId *pb.ReqId) (*pb.Bool, error) {
	return l.Pbfter.Poke(ctx, reqId)
}
func (l *Lbfter) Push(ctx context.Context, reqBody *pb.ReqBody) (*pb.Void, error) {
	return l.Pbfter.Push(ctx, reqBody)
}

func (l *Lbfter) GetAllRequests(ctx context.Context, void *pb.Void) (*pb.Requests, error) {
	return l.Pbfter.GetAllRequests(ctx, void)
}

/* Snowball */
func (l *Lbfter) GetVote(ctx context.Context, msg *pb.SeqNumMsg) (*pb.SeqNumMsg, error) {
	return l.Snower.GetVote(ctx, msg)
}

/* Pbfter */
