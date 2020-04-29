package proto

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	pb "github.com/spockqin/leaderless-bft/proto"
	"google.golang.org/grpc"
	"sort"
	tp "github.com/spockqin/leaderless-bft/types"
	"github.com/spockqin/leaderless-bft/util"
	// "google.golang.org/grpc/reflection"
	"net"
	"sync"
	"time"
)

type Stage int
const (
	Idle        Stage = iota // Node is created successfully, but the consensus process is not started yet.
	PrePrepared              // The ReqMsgs is processed successfully. The node is ready to head to the Prepare stage.
	Prepared                 // Same with `prepared` stage explained in the original paper.
	Committed                // Same with `committed-local` stage explained in the original paper.
)

type MsgLogs struct {
	ReqMsg        *tp.PbftReq
	PrepareMsgs   map[string]*tp.PrepareMsg
	CommitMsgs    map[string]*tp.CommitMsg
}

type State struct {
	ViewID		   int64
	MsgLogs        *MsgLogs
	LastSequenceID int64
	CurrentStage   Stage
}

type MsgBuffer struct {
	ReqMsgs        []*tp.PbftReq
	PrePrepareMsgs []*tp.PrePrepareMsg
	PrepareMsgs    []*tp.PrepareMsg
	CommitMsgs     []*tp.CommitMsg
}

type Pbfter struct {
	Gossiper
	NodeID        string
	ViewID        int64
	CurrentState  map[int64]*State
	CurrentStateLock sync.Mutex
	// CurrentState sync.Map
	CommittedMsgs []*tp.PbftReq
	MsgBuffer     *MsgBuffer
	MsgDelivery   chan interface{}
}

const f = 1; //TODO: set f to (R-1)/3

func (p *Pbfter) GetReq(ctx context.Context, request *pb.ReqBody) (*pb.Void, error) {
	var req tp.PbftReq
	err := json.Unmarshal(request.GetBody(), &req)
	if err != nil {
		log.Error("Error happens when unmarshal request [GetReq]")
	}

	// LogMsg(req)

	err, seqID := p.createStateForNewConsensus(-1)
	if err != nil {
		return &pb.Void{}, err
	}

	p.CurrentStateLock.Lock()
	state := p.CurrentState[seqID]
	p.CurrentStateLock.Unlock()
	prePrepareMsg, err := state.StartConsensus(&req, seqID)
	// prePrepareMsg, err := p.CurrentState.Load(seqID).StartConsensus(&req, seqID)
	if err != nil {
		log.Error("Error happens when starting consensus [SendReq]")
	}

	// broadcast prePrepareMsg
	prePrepareMsgBytes, err := json.Marshal(*prePrepareMsg)
	if err != nil {
		return &pb.Void{}, errors.New("[GetReq] prePrepareMsg marshal error!")
	} else {
		_, pushErr := p.Push(ctx, &pb.ReqBody{Body:prePrepareMsgBytes})
		if pushErr != nil {
			panic(errors.New("[GetReq] push prePrepareMsgBytes error!"))
		}
	}

	// LogStage("Pre-prepare", true, p.NodeID)

	return &pb.Void{}, nil
}

func (p *Pbfter) CheckCurrentPbfterMatchesThisMessenge(msgs interface{}) bool{
	p.CurrentStateLock.Lock()
	stateLen := len(p.CurrentState)
	p.CurrentStateLock.Unlock()
	if stateLen == 0 {
		return true
	}
	var seqIDs []int64
	p.CurrentStateLock.Lock()
	for k := range p.CurrentState {
		seqIDs = append(seqIDs, k)
	}
	p.CurrentStateLock.Unlock()
	
	sort.Slice(seqIDs, func(i, j int) bool {
		return seqIDs[i] < seqIDs[j]
	})
	switch msgs.(type) {
	case []*tp.PrePrepareMsg:
		return true
	case []*tp.PrepareMsg:
		for _, msg := range msgs.([]*tp.PrepareMsg) {
			if seqIDs[len(seqIDs)-1] == msg.SequenceID {
				return true
			} else {
				return false
			}
		}
	case []*tp.CommitMsg:
		for _, msg := range msgs.([]*tp.CommitMsg) {
			if seqIDs[len(seqIDs)-1] == msg.SequenceID {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func (p *Pbfter) CheckCommittedMsg(seqID int64) bool{
	flag := false
	if p.CommittedMsgs != nil {
		for _, commitedMsg := range p.CommittedMsgs {
			if seqID == commitedMsg.SequenceID {
				flag = true
				break
			}
		}
	}
	return flag
}

func (p *Pbfter) ResolveMsg() {
	for {
		msgs := <-p.MsgDelivery
		//if !p.CheckCurrentPbfterMatchesThisMessenge(msgs) {
		//	p.MsgDelivery <- msgs
		//	return
		//}
		switch msgs.(type) {
		case []*tp.PrePrepareMsg:
			for _, msg := range msgs.([]*tp.PrePrepareMsg) {
				err := p.GetPrePrepare(msg)
				if err != nil {
					log.Error("[ResolveMsg] resolve PrePrepare Error")
					panic(-1)
				}
			}
		case []*tp.PrepareMsg:
			for _, msg := range msgs.([]*tp.PrepareMsg) {
				if p.CheckCommittedMsg(msg.SequenceID) {
					continue
				}

				err := p.GetPrepare(msg)
				if err != nil {
					log.Error("[ResolveMsg] resolve Prepare Error")
					panic(-1)
				}
			}
		case []*tp.CommitMsg:
			for _, msg := range msgs.([]*tp.CommitMsg) {
				if p.CheckCommittedMsg(msg.SequenceID) {
					continue
				}

				err := p.GetCommit(msg)
				if err != nil {
					log.Error("[ResolveMsg] resolve Commit Error")
					panic(-1)
				}
			}
		}

	}
}

// node sends prePrepare msg
func (p *Pbfter) GetPrePrepare(msg *tp.PrePrepareMsg) (error) {
	//LogMsg(msg)
	err, _ := p.createStateForNewConsensus(msg.SequenceID)
	if err != nil {
		log.Error("[GetPrePrepare] createStateForNewConsensus error")
		panic(err)
	}
	p.CurrentStateLock.Lock()
	state := p.CurrentState[msg.SequenceID]
	p.CurrentStateLock.Unlock()
	prePareMsg, err := state.PrePrepare(msg)
	if err != nil {
		log.Error("[GetPrePrepare] PrePrepare error")
		panic(err)
	}
	if prePareMsg != nil {
		prePareMsg.NodeID = p.NodeID
		// LogStage("Pre-prepare", true, p.NodeID)
		// add msg to its own msglogs
		state.MsgLogs.PrepareMsgs[p.NodeID] = prePareMsg
		prepareMsgBytes, err := json.Marshal(*prePareMsg)
		if err != nil {
			return errors.New("[GetPrePrepare] prepareMsg marshal error!")
		} else {
			//ctx, cancel := context.WithTimeout(context.Background(), tp.GRPC_TIMEOUT)
			_, pushErr := p.Push(context.Background(), &pb.ReqBody{Body:prepareMsgBytes})
			if pushErr != nil {
				panic(errors.New("[GetPrePrepare] push prepareMsgBytes error!"))
			}
			//defer cancel()
		}
		//LogStage("Prepare", false, p.NodeID)
	} else {
		panic(errors.New("[GetPrePrepare] get empty prepareMsg"))
	}
	return nil
}

func (state *State) PrePrepare(prePrepareMsg *tp.PrePrepareMsg) (*tp.PrepareMsg, error) {
	// Get ReqMsgs and save it to its logs like the primary.
	state.MsgLogs.ReqMsg = prePrepareMsg.Req

	// Verify if v, n(a.k.a. sequenceID), d are correct.
	if !state.verifyMsg(prePrepareMsg.ViewID, prePrepareMsg.SequenceID, prePrepareMsg.Digest) {
		return nil, errors.New("pre-prepare message is corrupted")
	}

	// Change the stage to pre-prepared.
	state.CurrentStage = PrePrepared

	return &tp.PrepareMsg{
		ViewID: state.ViewID,
		SequenceID: prePrepareMsg.SequenceID,
		Digest: prePrepareMsg.Digest,
		MsgType: "PrepareMsg",
	}, nil
}

// node handles prepare msg
func (p *Pbfter) GetPrepare(msg *tp.PrepareMsg) (error) {
	//LogMsg(msg)

	p.CurrentStateLock.Lock()
	state := p.CurrentState[msg.SequenceID]
	p.CurrentStateLock.Unlock()
	commitMsg, err := state.Prepare(msg)
	if err != nil {
		log.Info(err)
	}

	if commitMsg != nil {
		commitMsg.NodeID = p.NodeID
		// LogStage("Prepare", true, p.NodeID)
		// add msg to its own msglogs
		state.MsgLogs.CommitMsgs[p.NodeID] = commitMsg
		commitMsgBytes, err := json.Marshal(*commitMsg)
		if err != nil {
			return errors.New("[GetPrepare] commitMsg marshal error!")
		} else {
			//ctx, cancel := context.WithTimeout(context.Background(), tp.GRPC_TIMEOUT)
			_, pushErr := p.Push(context.Background(), &pb.ReqBody{Body:commitMsgBytes})
			if pushErr != nil {
				panic(errors.New("[GetPrepare] push commitMsg error!"))
			}
			//defer cancel()
		}
		//LogStage("Commit", false, p.NodeID)
	}

	return nil
}

func (state *State) Prepare(prepareMsg *tp.PrepareMsg) (*tp.CommitMsg, error)  {
	if !state.verifyMsg(prepareMsg.ViewID, prepareMsg.SequenceID, prepareMsg.Digest) {
		log.Error("prepare message is corrupted")
		panic("prepare message is corrupted")
	}

	// Append msg to its logs
	state.MsgLogs.PrepareMsgs[prepareMsg.NodeID] = prepareMsg

	// Print current voting status
	//fmt.Printf("[Prepare-Vote]: %d\n", len(state.MsgLogs.PrepareMsgs))

	if state.prepared() {
		// Change the stage to prepared.
		state.CurrentStage = Prepared

		return &tp.CommitMsg{
			ViewID:               state.ViewID,
			SequenceID:           prepareMsg.SequenceID,
			Digest:               prepareMsg.Digest,
			MsgType:              "CommitMsg",
		}, nil
	}

	return nil, errors.New("[Prepare] haven't received 2f+1 prepare msg yet")
}

// node handles commit msg
func (p *Pbfter) GetCommit(msg *tp.CommitMsg) (error) {
	//LogMsg(msg)

	p.CurrentStateLock.Lock()
	state := p.CurrentState[msg.SequenceID]
	p.CurrentStateLock.Unlock()
	replyMsg, committedReq, err := state.Commit(msg)
	if err != nil {
		log.Info(err)
	} 

	if replyMsg != nil {
		if committedReq == nil {
			return errors.New("committed message is nil, even though the reply message is not nil")
		}

		replyMsg.NodeID = p.NodeID
		p.CommittedMsgs = append(p.CommittedMsgs, committedReq)

		p.CurrentStateLock.Lock()
		p.CurrentState[msg.SequenceID] = nil
		p.CurrentStateLock.Unlock()

		// LogStage("Commit", true, p.NodeID)
		log.Info("one commit")
		fmt.Println("Testing Timestamp:", time.Now().Unix())
		// log.Info("Committed message: ", msg)
		// LogStage("Reply", true, p.NodeID)
	}

	return nil
}

func (state *State) Commit(commitMsg *tp.CommitMsg) (*tp.ReplyMsg, *tp.PbftReq, error) {
	if !state.verifyMsg(commitMsg.ViewID, commitMsg.SequenceID, commitMsg.Digest) {
		log.Error("commit message is corrupted")
		panic("commit message is corrupted")
	}

	// Append msg to its logs
	state.MsgLogs.CommitMsgs[commitMsg.NodeID] = commitMsg

	// Print current voting status
	// fmt.Printf("[Commit-Vote]: %d\n", len(state.MsgLogs.CommitMsgs))

	if state.committed() {
		// This node executes the requested operation locally and gets the result.
		result := "Executed"

		// Change the stage to prepared.
		state.CurrentStage = Committed

		return &tp.ReplyMsg{
			ViewID:               state.ViewID,
			Timestamp:            state.MsgLogs.ReqMsg.Timestamp,
			ClientID:             state.MsgLogs.ReqMsg.ClientID,
			Result:               result,
			MsgType:              "ReplyMsg",
		}, state.MsgLogs.ReqMsg, nil
	}

	return nil, nil, errors.New("[Commit] haven't received 2f+1 commit msg yet")
}

func (p *Pbfter) createStateForNewConsensus(msgSeqID int64) (error, int64) {
	var lastSeqID int64
	if len(p.CommittedMsgs) == 0 {
		lastSeqID = -1
	} else {
		lastSeqID = p.CommittedMsgs[len(p.CommittedMsgs) - 1].SequenceID
	}

	var seqID int64
	if msgSeqID == -1 {
		seqID = time.Now().UnixNano()
		if lastSeqID != -1 {
			for lastSeqID >= seqID {
				seqID += 1
			}
		}
	} else {
		seqID = msgSeqID
	}

	p.CurrentStateLock.Lock()
	p.CurrentState[seqID] = createState(p.ViewID, lastSeqID)
	p.CurrentStateLock.Unlock()

	return nil, seqID
}

func createState (viewID int64, seqID int64) *State {
	return &State{
		ViewID:		    viewID,
		MsgLogs:        &MsgLogs{
			ReqMsg:      nil,
			PrepareMsgs: make(map[string]*tp.PrepareMsg),
			CommitMsgs:  make(map[string]*tp.CommitMsg),
		},
		LastSequenceID: seqID,
		CurrentStage:   Idle,
	}
}

func (state *State) StartConsensus(req *tp.PbftReq, seqID int64) (*tp.PrePrepareMsg, error) {
	req.SequenceID = seqID

	state.MsgLogs.ReqMsg = req

	digest, err := util.Digest(req)
	if err != nil {
		log.Error("Error happens when getting digest of request [StartConsensus]")
		panic(err)
	}

	state.CurrentStage = PrePrepared

	return &tp.PrePrepareMsg{
		ViewID:               state.ViewID,
		SequenceID:           seqID,
		Digest:               digest,
		Req:                  req,
		MsgType:			  "PrePrepareMsg",
	}, nil
}

func (state *State) verifyMsg(viewID int64, sequenceID int64, digestGot string) bool {
	// Wrong view. That is, wrong configurations of peers to start the consensus.
	if state.ViewID != viewID {
		return false
	}

	// Check if the Primary sent fault sequence number. => Faulty primary.
	// TODO: adopt upper/lower bound check.
	if state.LastSequenceID != -1 {
		if state.LastSequenceID >= sequenceID {
			return false
		}
	}

	digest, err := util.Digest(state.MsgLogs.ReqMsg)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Check digest.
	if digestGot != digest {
		return false
	}

	return true
}

func CreatePbfter(nodeID string, viewID int64, ip string, allIps []string) *Pbfter {
	newPbfter := &Pbfter{
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
	}

	return newPbfter
}

func (p *Pbfter) PbfterUp() {
	lis, err := net.Listen("tcp", p.ip)
	if err != nil {
		log.WithField("ip", p.ip).Error("Cannot listen on tcp [pbfter]")
		panic(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterGossipServer(grpcServer, p)
	pb.RegisterPbftServer(grpcServer, p)
	// reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.WithField("ip", p.ip).Error("Cannot serve [pbfter]")
		panic(err)
	}
}

func LogMsg(msg interface{}) {
	switch msg.(type) {
	case *tp.PbftReq:
		reqMsg := msg.(*tp.PbftReq)
		fmt.Printf("[REQUEST] ClientID: %s, Timestamp: %d, Operation: %s\n", reqMsg.ClientID, reqMsg.Timestamp, reqMsg.Operation)
	case *tp.PrePrepareMsg:
		prePrepareMsg := msg.(*tp.PrePrepareMsg)
		fmt.Printf("[PREPREPARE] ClientID: %s, Operation: %s, SequenceID: %d\n", prePrepareMsg.Req.ClientID, prePrepareMsg.Req.Operation, prePrepareMsg.SequenceID)
	case *tp.PrepareMsg:
		prePareMsg := msg.(*tp.PrepareMsg)
		fmt.Printf("[PREPARE] NodeID: %s\n", prePareMsg.NodeID)
	case *tp.CommitMsg:
		commitMsg := msg.(*tp.CommitMsg)
		fmt.Printf("[COMMIT] NodeID: %s\n", commitMsg.NodeID)
	}
}

func LogStage(stage string, isDone bool, nodeID string) {
	if isDone {
		log.Info(fmt.Sprintf("[STAGE-DONE] %s nodeID: %v\n", stage, nodeID))
		// fmt.Printf("[STAGE-DONE] %s nodeID: %v\n", stage, nodeID)
	} else {
		log.Info(fmt.Sprintf("[STAGE-BEGIN] %s nodeID: %v\n", stage, nodeID))
		// fmt.Printf("[STAGE-BEGIN] %s nodeID: %v\n", stage, nodeID)
	}
}

func (state *State) prepared() bool {
	if state.MsgLogs.ReqMsg == nil {
		return false
	}

	if len(state.MsgLogs.PrepareMsgs) < (2*f + 1) { // 2f + 1 (itself)
		return false
	}

	return true
}

func (state *State) committed() bool {
	if !state.prepared() {
		return false
	}

	if len(state.MsgLogs.CommitMsgs) < (2*f + 1) { // 2f + 1 (itself)
		return false
	}

	return true
}