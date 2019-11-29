package proto

import (
	"encoding/json"
	"errors"
	node "github.com/spockqin/leaderless-bft/proto"
	"github.com/spockqin/leaderless-bft/types"
	"google.golang.org/grpc"
	"net"
	"time"
	log "github.com/sirupsen/logrus"
	"fmt"
)

type Pbfter struct {
	Gossiper
	NodeID			string
	NodeTable		map[string]string // map[nodeID] -> ip
	ViewID			int64
	Leader			string // leader nodeID
	CurrentState	*node.State
	CommittedMsgs	[]*node.RequestMsg
	MsgBuffer		*MsgBuffer
	MsgEntrance		chan interface{}
	MsgDelivery		chan interface{}
	Alarm			chan bool
}

type MsgBuffer struct {
	ReqMsgs        []*node.RequestMsg
	PrePrepareMsgs []*node.PrePrepareMsg
	PrepareMsgs    []*node.VoteMsg
	CommitMsgs     []*node.VoteMsg
}

const ResolvingTimeDuration = time.Millisecond * 2000 // 2 second.

func (pbfter *Pbfter) PbfterUp() {
	lis, err := net.Listen("tcp", pbfter.ip)
	if err != nil {
		log.WithField("ip", s.ip).Error("Cannot listen on tcp [pbfter]")
	}

	grpcServer := grpc.NewServer()

	node.RegisterGossipServer(grpcServer, pbfter)

	grpcServer.Serve(lis)
}

func CreatePbfter(nodeID string, ip string, nodeTable map[string]string) *Pbfter {
	const viewID = 100000

	pbfter := &Pbfter{
		NodeID:        nodeID,
		NodeTable:	   nodeTable,
		ViewID:        viewID,
		Leader:		   "0",
		CurrentState:  nil,
		CommittedMsgs: make([]*node.RequestMsg, 0),
		MsgBuffer:     &MsgBuffer{
			ReqMsgs:        make([]*node.RequestMsg, 0),
			PrePrepareMsgs: make([]*node.PrePrepareMsg, 0),
			PrepareMsgs:    make([]*node.VoteMsg, 0),
			CommitMsgs:     make([]*node.VoteMsg, 0),
		},
		MsgEntrance:   make(chan interface{}),
		MsgDelivery:   make(chan interface{}),
		Alarm:         make(chan bool),
	}

	pbfter.ip = ip
	pbfter.peers = make([]string, 0)

	// start deliver messages
	go pbfter.dispatchMsg()

	// start timeout alarm trigger
	go pbfter.timeoutAlarm()

	// start message resolver
	go pbfter.resolveMsg()

	return pbfter
}

// deliver message
func (pbfter *Pbfter) dispatchMsg() {
	for {
		select {
		case msg := <-pbfter.MsgEntrance:
			err := pbfter.routeMsg(msg)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("dispatchMsg routeMsg error\n")
			}
		case <-pbfter.Alarm:
			err := pbfter.routeMsgWhenAlarmed()
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("dispatchMsg routeMsgWhenAlarmed error\n")
			}
		}
	}
}

// handle arrived messages
func (pbfter *Pbfter) routeMsg(msg interface{}) []error {
	switch msg.(type) {
	case *node.RequestMsg:
		if pbfter.CurrentState == nil {
			// Copy buffered messages first
			msgs := make([]*node.RequestMsg, len(pbfter.MsgBuffer.ReqMsgs))
			copy(msgs, pbfter.MsgBuffer.ReqMsgs)

			// append a newly arrived message
			msgs = append(msgs, msg.(*node.RequestMsg))

			// empty the buffer
			pbfter.MsgBuffer.ReqMsgs = make([]*node.RequestMsg, 0)

			// send messages
			pbfter.MsgDelivery <- msgs

		} else {
			pbfter.MsgBuffer.ReqMsgs = append(pbfter.MsgBuffer.ReqMsgs, msg.(*node.RequestMsg))
		}
	case *node.PrePrepareMsg:
		if pbfter.CurrentState == nil {
			// Copy buffered message first
			msgs := make([]*node.PrePrepareMsg, len(pbfter.MsgBuffer.PrePrepareMsgs))
			copy(msgs, pbfter.MsgBuffer.PrePrepareMsgs)

			// append a newly arrived message
			msgs = append(msgs, msg.(*node.PrePrepareMsg))

			// empty the buffer
			pbfter.MsgBuffer.PrePrepareMsgs = make([]*node.PrePrepareMsg, 0)

			// send message
			pbfter.MsgDelivery <- msgs
		} else {
			pbfter.MsgBuffer.PrePrepareMsgs = append(pbfter.MsgBuffer.PrePrepareMsgs, msg.(*node.PrePrepareMsg))
		}
	case *node.VoteMsg:
		if msg.(*node.VoteMsg).MsgType == node.PrepareMsg {
			if pbfter.CurrentState == nil || pbfter.CurrentState.CurrentStage != node.PrePrepared {
				pbfter.MsgBuffer.PrepareMsgs = append(pbfter.MsgBuffer.PrepareMsgs, msg.(*node.VoteMsg))
			} else {
				// copy buffered message first
				msgs := make([]*node.VoteMsg, len(pbfter.MsgBuffer.PrepareMsgs))
				copy(msgs, pbfter.MsgBuffer.PrepareMsgs)

				// append a newly arrived message
				msgs = append(msgs, msg.(*node.VoteMsg))

				// empty the buffer
				pbfter.MsgBuffer.PrepareMsgs = make([]*node.VoteMsg, 0)

				// send messages
				pbfter.MsgDelivery <- msgs
			}
		} else if msg.(*node.VoteMsg).MsgType == node.CommitMsg{
			if pbfter.CurrentState == nil || pbfter.CurrentState.CurrentStage != node.Prepared {
				pbfter.MsgBuffer.CommitMsgs = append(pbfter.MsgBuffer.CommitMsgs, msg.(*node.VoteMsg))
			} else {
				// copy buffered message first
				msgs := make([]*node.VoteMsg, len(pbfter.MsgBuffer.CommitMsgs))
				copy(msgs, pbfter.MsgBuffer.CommitMsgs)

				// append a newly arrived message
				msgs = append(msgs, msg.(*node.VoteMsg))

				// empty the buffer
				pbfter.MsgBuffer.PrepareMsgs = make([]*node.VoteMsg, 0)

				// send messages
				pbfter.MsgDelivery <- msgs
			}
		}
	}

	return nil
}

// handle arrived messages when timeout (just send out without caching)
func (pbfter *Pbfter) routeMsgWhenAlarmed() []error {
	if pbfter.CurrentState == nil {
		// check ReqMsgs, send them
		if len(pbfter.MsgBuffer.ReqMsgs) != 0 {
			msgs := make([]*node.RequestMsg, len(pbfter.MsgBuffer.ReqMsgs))
			copy(msgs, pbfter.MsgBuffer.ReqMsgs)

			pbfter.MsgDelivery <- msgs
		}

		// check PrePrepareMsgs, send them
		if len(pbfter.MsgBuffer.PrePrepareMsgs) != 0 {
			msgs := make([]*node.PrePrepareMsg, len(pbfter.MsgBuffer.PrePrepareMsgs))
			copy(msgs, pbfter.MsgBuffer.PrePrepareMsgs)

			pbfter.MsgDelivery <- msgs
		}
	} else {
		switch pbfter.CurrentState.CurrentStage {
		case node.PrePrepared:
			// Check Prepare Msgs, send them
			if len(pbfter.MsgBuffer.PrepareMsgs) != 0 {
				msgs := make([]*node.VoteMsg, len(pbfter.MsgBuffer.PrepareMsgs))
				copy(msgs, pbfter.MsgBuffer.PrepareMsgs)

				pbfter.MsgDelivery <- msgs
			}
		case node.Prepared:
			// Check CommitMsgs, send them
			if len(pbfter.MsgBuffer.CommitMsgs) != 0 {
				msgs := make([]*node.VoteMsg, len(pbfter.MsgBuffer.CommitMsgs))
				copy(msgs, pbfter.MsgBuffer.CommitMsgs)

				pbfter.MsgDelivery <- msgs
			}
		}
	}

	return nil
}

func (pbfter *Pbfter) timeoutAlarm() {
	for {
		time.Sleep(ResolvingTimeDuration)
		pbfter.Alarm <- true
	}
}

func (pbfter *Pbfter) resolveMsg() {
	for {
		// get buffered message from the dispatcher
		msgs := <- pbfter.MsgDelivery
		switch msgs.(type) {
		case []*node.RequestMsg:
			errs := pbfter.resolveRequestMsg(msgs.([]*node.RequestMsg))
			if len(errs) != 0 {
				log.Error("PBFTER resolveMsg resolveRequestMsg error", errs)
			}
		case []*node.PrePrepareMsg:
			errs := pbfter.resolvePrePrepareMsg(msgs.([]*node.PrePrepareMsg))
			if len(errs) != 0 {
				log.Error("PBFTER resolveMsg resolvePrePrepareMsg error", errs)
			}
		case []*node.VoteMsg:
			voteMsgs := msgs.([]*node.VoteMsg)
			if len(voteMsgs) == 0 {
				break
			}
			if voteMsgs[0].MsgType == node.PrepareMsg {
				errs := pbfter.resolvePrepareMsg(voteMsgs)
				if len(errs) != 0 {
					log.Error("PBFTER resolveMsg resolvePrepareMsg error", errs)
				}
			} else if voteMsgs[0].MsgType == node.CommitMsg {
				errs := pbfter.resolveCommitMsg(voteMsgs)
				if len(errs) != 0 {
					log.Error("PBFTER resolveMsg resolveCommitMsg error", errs)
				}
			}
		}
	}
}

func (pbfter *Pbfter) resolveRequestMsg(msgs []*node.RequestMsg) []error {
	errs := make([]error, 0)

	for _, resMsg := range msgs {
		err := pbfter.GetReq(resMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

func (pbfter *Pbfter) resolvePrePrepareMsg(msgs []*node.PrePrepareMsg) []error {
	errs := make([]error, 0)

	for _, prePrepareMsg := range msgs {
		err := pbfter.GetPrePrepare(prePrepareMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

func (pbfter *Pbfter) resolvePrepareMsg(msgs []*node.VoteMsg) [] error {
	errs := make([]error, 0)

	for _, prePareMsg := range msgs {
		err := pbfter.GetPrepare(prePareMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

func (pbfter *Pbfter) resolveCommitMsg(msgs []*node.VoteMsg) []error {
	errs := make([]error, 0)

	for _, commitMsg := range msgs {
		err := pbfter.GetCommit(commitMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

func (pbfter *Pbfter) GetReq(reqMsg *node.RequestMsg) error {
	LogMsg(reqMsg)

	// create a new state for the new consensus
	err := pbfter.createStateForNewConsensus()
	if err != nil {
		log.Error("PBFTER: GetReq createStateForNewConsensus error")
		return err
	}

	// start the consensus process
	prePrepareMsg, err := pbfter.CurrentState.StartConsensus(reqMsg)
	if err != nil {
		log.Error("PBFTER: GetReq StartConsensus error")
		return err
	}

	LogStage(fmt.Sprintf("Consensus Process (ViewID:%d)", pbfter.CurrentState.ViewID), false)

	// send getPrePrepare message
	if prePrepareMsg != nil {
		pbfter.Broadcast(prePrepareMsg)
		LogStage("Pre-prepare", true)
	}

	return nil
}

func (pbfter *Pbfter) GetPrePrepare(prePrepareMsg *node.PrePrepareMsg) error {
	LogMsg(prePrepareMsg)

	// create a new state
	err := pbfter.createStateForNewConsensus()
	if err != nil {
		log.Error("PBFTER: GetPrePrepare createStateForNewConsensus error")
		return err
	}

	prePareMsg, err := pbfter.CurrentState.PrePrepare(prePrepareMsg)
	if err != nil {
		log.Error("PBFTER: GetPrePrepare PrePrepare error")
		return err
	}

	if prePareMsg != nil {
		prePareMsg.NodeID = pbfter.NodeID
		LogStage("Pre-prepare", true)
		pbfter.Broadcast(prePareMsg)
		LogStage("Prepare", false)
	}

	return nil
}

func (pbfter *Pbfter) GetPrepare(prePareMsg *node.VoteMsg) error {
	LogMsg(prePareMsg)
	commitMsg, err := pbfter.CurrentState.Prepare(prePareMsg)
	if err != nil {
		log.Error("PBFTER: GetPrepare Prepare error")
		return err
	}

	if commitMsg != nil {
		commitMsg.NodeID = pbfter.NodeID
		LogStage("Prepare", true)
		pbfter.Broadcast(commitMsg)
		LogStage("Commit", false)
	}

	return nil
}

func (pbfter *Pbfter) GetCommit(commitMsg *node.VoteMsg) error {
	LogMsg(commitMsg)

	replyMsg, committedMsg, err := pbfter.CurrentState.Commit(commitMsg)
	if err != nil {
		log.Error("PBFTER: GetCommit Commit error")
		return err
	}

	if replyMsg != nil {
		if committedMsg == nil {
			log.Error("PBFTER: committed message is nil, even though the reply message is not nil")
			return errors.New("committed message is nil, even though the reply message is not nil")
		}

		replyMsg.NodeID = pbfter.NodeID

		// save the last version of committed messages to ndoe
		pbfter.CommittedMsgs = append(pbfter.CommittedMsgs, committedMsg)
		LogStage("Commit", true)
		err:= pbfter.Reply(replyMsg)
		if err != nil {
			log.Error("PBFTER: GetCommit Reply error")
			return err
		}
		LogStage("Reply", true)
	}

	return nil
}

func (pbfter *Pbfter) Reply(msg *node.ReplyMsg) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Error("PBFTER Reply json marshal error")
		return err
	}

	// QUESTION
	maxSoc := make(chan bool, types.MAX_SOCKETS)
	maxSoc <- true
	pbfter.sendGossip(pbfter.NodeTable[pbfter.Leader], jsonMsg, maxSoc)

	return nil
}

func (pbfter *Pbfter) createStateForNewConsensus() error {
	// check if there is an ongoing consensus process
	if pbfter.CurrentState != nil {
		return errors.New("another consensus is ongoing")
	}

	// get the last sequence ID
	var lastSequenceID int64
	if len(pbfter.CommittedMsgs) == 0{
		lastSequenceID = -1
	} else {
		lastSequenceID = pbfter.CommittedMsgs[len(pbfter.CommittedMsgs) - 1].SequenceID
	}

	// calculate bad nodes tolerance
	f := (len(pbfter.NodeTable) - 1) / 3

	// create a new state
	pbfter.CurrentState = node.CreateState(pbfter.ViewID, lastSequenceID, f)

	LogStage("Create the replica status", true)

	return nil
}

func (pbfter *Pbfter) Broadcast (msg interface{}) map[string]error {
	errorMap := make(map[string]error)
	maxSoc := make(chan bool, types.MAX_SOCKETS)
	for nodeID, ip := range pbfter.NodeTable {
		if nodeID == pbfter.NodeID {
			continue
		}

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			errorMap[nodeID] = err
			continue
		}

		maxSoc <- true // blocks if maxSoc is full
		go pbfter.sendGossip(ip, jsonMsg, maxSoc)
	}
	if len(errorMap) == 0 {
		return nil
	} else {
		return errorMap
	}
}

func LogMsg(msg interface{}) {
	switch msg.(type) {
	case *node.RequestMsg:
		reqMsg := msg.(*node.RequestMsg)
		log.WithFields(log.Fields{
			"msgType": "REQUEST",
			"ClientID": reqMsg.ClientID,
			"TimeStamp": reqMsg.Timestamp,
			"Operation": reqMsg.Operation,
		}).Info("Log REQUEST message")
	case *node.PrePrepareMsg:
		prePrepareMsg := msg.(*node.PrePrepareMsg)
		log.WithFields(log.Fields{
			"msgType": "PREPREPARE",
			"ClientID": prePrepareMsg.RequestMsg.ClientID,
			"TimeStamp": prePrepareMsg.RequestMsg.Timestamp,
			"Operation": prePrepareMsg.RequestMsg.Operation,
			"SequenceID": prePrepareMsg.SequenceID,
		}).Info("Log PREPREPARE message")
	case *node.VoteMsg:
		voteMsg := msg.(*node.VoteMsg)
		if voteMsg.MsgType == node.PrepareMsg {
			log.WithFields(log.Fields{
				"msgType": "PREPARE",
				"ViewID": voteMsg.ViewID,
				"NodeID": voteMsg.NodeID,
				"SequenceID": voteMsg.SequenceID,
			}).Info("Log PREPARE message")
		} else if voteMsg.MsgType == node.CommitMsg {
			log.WithFields(log.Fields{
				"msgType": "COMMIT",
				"ViewID": voteMsg.ViewID,
				"NodeID": voteMsg.NodeID,
				"SequenceID": voteMsg.SequenceID,
			}).Info("Log COMMIT message")
		}
	}
}

func LogStage(stage string, isDone bool) {
	if isDone {
		log.Info("[STAGE-DONE] %s\n", stage)
	} else {
		log.Info("[STAGE-BEGIN] %s\n", stage)
	}
}