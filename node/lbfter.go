package proto

import (
	"context"
	"encoding/json"
	"errors"
	// "fmt"
	log "github.com/sirupsen/logrus"
	pb "github.com/spockqin/leaderless-bft/proto"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/reflection"
	"net"
	"sync"
	tp "github.com/spockqin/leaderless-bft/types"
	util "github.com/spockqin/leaderless-bft/util"
)

type Lbfter struct {
	Gossiper
	Snower
	Pbfter
	// SeqIdDecided  map[string]chan bool
}

func (l *Lbfter) LSendReq(ctx context.Context, request *pb.ReqBody) (*pb.Void, error) {
	var req tp.PbftReq
	err := json.Unmarshal(request.GetBody(), &req)
	if err != nil {
		log.Error("Error when unmarshal request [LSendReq]")
	}

	LogMsg(req)

	err = l.createStateForNewConsensus()
	if err != nil {
		log.Error("Error when create consensus state")
		return &pb.Void{}, err
	}

	prePrepareMsg, err := l.startLbftConsensus(&req)
	if err != nil {
		log.Error("Error when starting consensus [LSendReq]")
	}

	// broadcast prePrepareMsg
	prePrepareMsgBytes, err := json.Marshal(*prePrepareMsg)
	if err != nil {
		return &pb.Void{}, errors.New("[LSendReq] prePrepareMsg marshal error!")
	} else {
		_, pushErr := l.Pbfter.Push(ctx, &pb.ReqBody{Body:prePrepareMsgBytes})
		if pushErr != nil {
			panic(errors.New("[LSendReq] push prePrepareMsgBytes error!"))
		}
	}

	LogStage("Pre-prepare", true, l.NodeID)

	return &pb.Void{}, nil
}

func (l *Lbfter) runSeqIdConsensus(operation string) {
	l.SendReq(context.Background(), &pb.ReqBody{Body: []byte(operation)})
	// l.SeqIdDecided[reqHash] <- true
}

func (l *Lbfter) startLbftConsensus(req *tp.PbftReq) (*tp.PrePrepareMsg, error) {
	reqHash := string(util.HashBytes([]byte(req.Operation)))
	// l.SeqIdDecided[reqHash] = make(chan bool)
	// go l.RunSeqIdConsensus(reqHash)
	// <-l.SeqIdDecided[reqHash]
	l.runSeqIdConsensus(req.Operation)
	req.SequenceID = l.finalSeqNums[string(reqHash)]

	state := l.CurrentState

	state.MsgLogs.ReqMsg = req

	digest, err := util.Digest(req)
	if err != nil {
		log.Error("Error happens when getting digest of request [StartLbftConsensus]")
		panic(err)
	}

	state.CurrentStage = PrePrepared

	return &tp.PrePrepareMsg{
		ViewID:               state.ViewID,
		SequenceID:           l.finalSeqNums[string(reqHash)],
		Digest:               digest,
		Req:                  req,
		MsgType:			  "PrePrepareMsg",
	}, nil
}

func CreateLbfter(nodeID string, viewID int64, ip string, allIps []string) *Lbfter {
	newLbfter := &Lbfter{
		Gossiper:      Gossiper{
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

		Snower:		   Snower{
			allIps:			  allIps,
			confidences:	  CreateConfidenceMap(),
			seqNum: 		  -1,
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
			CurrentState:  nil,
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
		// SeqIdDecided:  make(map[string]chan bool)
	}

	return newLbfter
}

func (l *Lbfter) LbfterUp() {
	lis, err := net.Listen("tcp", l.ip)
	if err != nil {
		log.WithField("ip", l.ip).Error("Cannot listen on tcp [lbfter]")
		panic(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterGossipServer(grpcServer, l)
	pb.RegisterSnowballServer(grpcServer, l)
	pb.RegisterPbftServer(grpcServer, l)
	pb.RegisterLbftServer(grpcServer, l)
	// reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.WithField("ip", l.ip).Error("Cannot serve [lbfter]")
		panic(err)
	}
}
