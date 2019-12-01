package proto

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	pb "github.com/spockqin/leaderless-bft/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"fmt"
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
	ReqMsg        *pb.PbftReq
	PrepareMsgs   map[string]*pb.PrepareMsg
	CommitMsgs    map[string]*pb.CommitMsg
}

type State struct {
	ViewID		   int64
	MsgLogs        *MsgLogs
	LastSequenceID int64
	CurrentStage   Stage
}

type Pbfter struct {
	Gossiper
	NodeID string
	ViewID int64
	CurrentState *State
	CommittedMsgs []*pb.PbftReq
}

const f = 1;

func (p *Pbfter) GetReq(ctx context.Context, req *pb.PbftReq) (*pb.Void, error) {
	LogMsg(req)

	// begin new consensus
	if p.CurrentState != nil {
		return &pb.Void{}, errors.New("another consensus is ongoing")
	}

	err := p.createStateForNewConsensus()
	if err != nil {
		return &pb.Void{}, nil
	}

	prePrepareMsg, err := p.CurrentState.StartConsensus(req)
	if err != nil {
		log.Error("Error happens when starting consensus [SendReq]")
	}

	LogStage(fmt.Sprintf("Consensus Process (ViewID:%d)", p.CurrentState.ViewID), false)

	// broadcast prePrepareMsg
	err = p.SendPrePrepare(prePrepareMsg)
	if err != nil {
		log.Error("SendPrePrepare Error [pbfter]")
		panic(err)
	}

	LogStage("Pre-prepare", true)

	return &pb.Void{}, nil
}

func (p *Pbfter) SendPrePrepare(prePrepareMsg *pb.PrePrepareMsg) (error){
	if prePrepareMsg != nil {
		for _, ip := range p.Gossiper.peers {
			conn, err := grpc.Dial(ip, grpc.WithInsecure())
			if err != nil {
				log.WithFields(log.Fields{
					"ip": p.ip,
					"node": ip,
					"error": err,
				}).Error("Cannot dial node\n")
				continue
			}
			client := pb.NewPbftClient(conn)
			_, err = client.GetPrePrepare(context.Background(), prePrepareMsg)
			if err != nil {
				log.WithFields(log.Fields{
					"ip": p.ip,
					"node": ip,
					"error": err,
				}).Error("Cannot send PrePrepareMsg\n")
				continue
			}
			conn.Close()
		}

	} else {
		log.Error("[SendPrePrepare] prePrepareMsg is empty")
		return errors.New("empty prePrepareMsg")
	}
	return nil
}

// node sends prePrepare msg
func (p *Pbfter) GetPrePrepare(ctx context.Context, msg *pb.PrePrepareMsg) (*pb.Void, error) {
	LogMsg(msg)
	err := p.createStateForNewConsensus()
	if err != nil {
		log.Error("[GetPrePrepare] createStateForNewConsensus error")
		panic(err)
	}

	prePareMsg, err := p.CurrentState.PrePrepare(msg)
	if err != nil {
		log.Error("[GetPrePrepare] PrePrepare error")
		panic(err)
	}

	if prePareMsg != nil {
		prePareMsg.NodeID = p.NodeID
		LogStage("Pre-prepare", true)
		err = p.SendPrepare(prePareMsg)
		if err != nil {
			log.Error("[GetPrePrepare] sendPrepare error")
			panic(err)
		}
		fmt.Printf("[pre-prepare] after broadcast")
		LogStage("Prepare", false)
	}

	return &pb.Void{}, nil
}

func (state *State) PrePrepare(prePrepareMsg *pb.PrePrepareMsg) (*pb.PrepareMsg, error) {
	// Get ReqMsgs and save it to its logs like the primary.
	state.MsgLogs.ReqMsg = prePrepareMsg.GetReq()

	// Verify if v, n(a.k.a. sequenceID), d are correct.
	if !state.verifyMsg(prePrepareMsg.ViewID, prePrepareMsg.SequenceID, prePrepareMsg.Digest) {
		return nil, errors.New("pre-prepare message is corrupted")
	}

	// Change the stage to pre-prepared.
	state.CurrentStage = PrePrepared

	return &pb.PrepareMsg{
		ViewID: state.ViewID,
		SequenceID: prePrepareMsg.SequenceID,
		Digest: prePrepareMsg.Digest,
	}, nil
}

func (p *Pbfter) SendPrepare(prePareMsg *pb.PrepareMsg) (error) {
	if prePareMsg != nil {
		for _, ip := range p.Gossiper.peers {
			conn, err := grpc.Dial(ip, grpc.WithInsecure())
			if err != nil {
				log.WithFields(log.Fields{
					"ip": p.ip,
					"node": ip,
					"error": err,
				}).Error("Cannot dial node\n")
				continue
			}
			client := pb.NewPbftClient(conn)
			_, err = client.GetPrepare(context.Background(), prePareMsg)
			if err != nil {
				log.WithFields(log.Fields{
					"ip": p.ip,
					"node": ip,
					"error": err,
				}).Error("Cannot send PrepareMsg\n")
				continue
			}
			conn.Close()
		}

	} else {
		log.Error("[SendPrepare] prepareMsg is empty")
		return errors.New("empty prepareMsg")
	}
	return nil
}

// node sends prepare msg
func (p *Pbfter) GetPrepare(ctx context.Context, msg *pb.PrepareMsg) (*pb.Void, error) {
	LogMsg(msg)

	commitMsg, err := p.CurrentState.Prepare(msg)
	if err != nil {
		log.Error("[GetPrepare] Prepare Error")
		panic(err)
	}

	if commitMsg != nil {
		commitMsg.NodeID = p.NodeID

		LogStage("Prepare", true)
		err = p.SendCommit(commitMsg)
		if err != nil {
			log.Error("[GetPrepare] sendCommit error")
			panic(err)
		}
		fmt.Printf("[GetPrepare] after broadcast")
		LogStage("Commit", false)
	}

	return &pb.Void{}, nil
}

func (state *State) Prepare(prepareMsg *pb.PrepareMsg) (*pb.CommitMsg, error)  {
	if !state.verifyMsg(prepareMsg.ViewID, prepareMsg.SequenceID, prepareMsg.Digest) {
		return nil, errors.New("prepare message is corrupted")
	}

	// Append msg to its logs
	state.MsgLogs.PrepareMsgs[prepareMsg.NodeID] = prepareMsg

	// Print current voting status
	fmt.Printf("[Prepare-Vote]: %d\n", len(state.MsgLogs.PrepareMsgs))

	if state.prepared() {
		// Change the stage to prepared.
		state.CurrentStage = Prepared

		return &pb.CommitMsg{
			ViewID:               state.ViewID,
			SequenceID:           prepareMsg.SequenceID,
			Digest:               prepareMsg.Digest,
		}, nil
	}

	return nil, errors.New("[Prepare] haven't received 2f+1 prepare msg")
}

func (p *Pbfter) SendCommit(commitMsg *pb.CommitMsg) (error) {
	if commitMsg != nil {
		for _, ip := range p.Gossiper.peers {
			conn, err := grpc.Dial(ip, grpc.WithInsecure())
			if err != nil {
				log.WithFields(log.Fields{
					"ip": p.ip,
					"node": ip,
					"error": err,
				}).Error("Cannot dial node\n")
				continue
			}
			client := pb.NewPbftClient(conn)
			_, err = client.GetCommit(context.Background(), commitMsg)
			if err != nil {
				log.WithFields(log.Fields{
					"ip": p.ip,
					"node": ip,
					"error": err,
				}).Error("Cannot send CommitMsg\n")
				continue
			}
			conn.Close()
		}

	} else {
		log.Error("[SendCommit] commitMsg is empty")
		return errors.New("empty commitMsg")
	}
	return nil
}

// node sends commit msg
func (p *Pbfter) GetCommit(ctx context.Context, msg *pb.CommitMsg) (*pb.Void, error) {
	LogMsg(msg)

	replyMsg, committedReq, err := p.CurrentState.Commit(msg)
	if err != nil {
		log.Error("[GetCommit] Commit Error")
		panic(err)
	}

	if replyMsg != nil {
		if committedReq == nil {
			return &pb.Void{}, errors.New("committed message is nil, even though the reply message is not nil")
		}

		replyMsg.NodeID = p.NodeID
		p.CommittedMsgs = append(p.CommittedMsgs, committedReq)

		LogStage("Commit", true)
		//_, err := p.SendReply(context.Background(), replyMsg)
		//if err != nil {
		//	log.Error("[GetCommit] SendReply Error")
		//	panic(err)
		//}
		LogStage("Reply", true)
	}

	return &pb.Void{}, nil
}

func (state *State) Commit(commitMsg *pb.CommitMsg) (*pb.ReplyMsg, *pb.PbftReq, error) {
	if !state.verifyMsg(commitMsg.ViewID, commitMsg.SequenceID, commitMsg.Digest) {
		return nil, nil, errors.New("commit message is corrupted")
	}

	// Append msg to its logs
	state.MsgLogs.CommitMsgs[commitMsg.NodeID] = commitMsg

	// Print current voting status
	fmt.Printf("[Commit-Vote]: %d\n", len(state.MsgLogs.CommitMsgs))

	if state.committed() {
		// This node executes the requested operation locally and gets the result.
		result := "Executed"

		// Change the stage to prepared.
		state.CurrentStage = Committed

		return &pb.ReplyMsg{
			ViewID:               state.ViewID,
			Timestamp:            state.MsgLogs.ReqMsg.Timestamp,
			ClientID:             state.MsgLogs.ReqMsg.ClientID,
			Result:               result,
		}, state.MsgLogs.ReqMsg, nil
	}

	return nil, nil, errors.New("[Commit] haven't received 2f+1 commit msg")
}

// node sends reply msg to the client
//func (p *Pbfter) SendReply(ctx context.Context, msg *pb.ReplyMsg) (*pb.Void, error) {
//
//}

func (p *Pbfter) createStateForNewConsensus() error {
	if p.CurrentState != nil {
		return errors.New("another consensus is ongoing")
	}

	var lastSeqID int64
	if len(p.CommittedMsgs) == 0 {
		lastSeqID = -1
	} else {
		lastSeqID = p.CommittedMsgs[len(p.CommittedMsgs) - 1].SequenceID
	}

	p.CurrentState = createState(p.ViewID, lastSeqID)

	LogStage("Create new state for new consensus done", true)

	return nil
}

func createState (viewID int64, seqID int64) *State {
	return &State{
		ViewID:		    viewID,
		MsgLogs:        &MsgLogs{
			ReqMsg:      nil,
			PrepareMsgs: make(map[string]*pb.PrepareMsg),
			CommitMsgs:  make(map[string]*pb.CommitMsg),
		},
		LastSequenceID: seqID,
		CurrentStage:   Idle,
	}
}

func (state *State) StartConsensus(req *pb.PbftReq) (*pb.PrePrepareMsg, error) {
	seqID := time.Now().UnixNano()

	if state.LastSequenceID != -1 {
		seqID = state.LastSequenceID + 1
	}

	req.SequenceID = seqID

	state.MsgLogs.ReqMsg = req

	digest, err := digest(req)
	if err != nil {
		log.Error("Error happens when getting digest of request [StartConsensus]")
		panic(err)
	}

	state.CurrentStage = PrePrepared

	return &pb.PrePrepareMsg{
		ViewID:               state.ViewID,
		SequenceID:           seqID,
		Digest:               digest,
		Req:                  req,
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

	digest, err := digest(state.MsgLogs.ReqMsg)
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
	newPbfter := new(Pbfter)
	newPbfter.ip = ip
	newPbfter.peers = make([]string, 0)
	newPbfter.hashes = make(map[string]bool)
	newPbfter.requests = make([]string, 0)

	newPbfter.NodeID = nodeID
	newPbfter.ViewID = viewID
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
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.WithField("ip", p.ip).Error("Cannot serve [pbfter]")
		panic(err)
	}
}

func LogMsg(msg interface{}) {
	switch msg.(type) {
	case *pb.PbftReq:
		reqMsg := msg.(*pb.PbftReq)
		fmt.Printf("[REQUEST] ClientID: %s, Timestamp: %d, Operation: %s\n", reqMsg.ClientID, reqMsg.Timestamp, reqMsg.Operation)
	case *pb.PrePrepareMsg:
		prePrepareMsg := msg.(*pb.PrePrepareMsg)
		fmt.Printf("[PREPREPARE] ClientID: %s, Operation: %s, SequenceID: %d\n", prePrepareMsg.Req.ClientID, prePrepareMsg.Req.Operation, prePrepareMsg.SequenceID)
	case *pb.PrepareMsg:
		prePareMsg := msg.(*pb.PrepareMsg)
		fmt.Printf("[PREPARE] NodeID: %s\n", prePareMsg.NodeID)
	case *pb.CommitMsg:
		commitMsg := msg.(*pb.CommitMsg)
		fmt.Printf("[COMMIT] NodeID: %s\n", commitMsg.NodeID)
	}
}

func LogStage(stage string, isDone bool) {
	if isDone {
		fmt.Printf("[STAGE-DONE] %s\n", stage)
	} else {
		fmt.Printf("[STAGE-BEGIN] %s\n", stage)
	}
}

func digest(object interface{}) (string, error) {
	msg, err := json.Marshal(object)

	if err != nil {
		return "", err
	}

	return Hash(msg), nil
}

func Hash(content []byte) string {
	h := sha256.New()
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
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

func (p *Pbfter) SendMsg(msg interface{}) {
	var buf []byte
	buf, err := json.Marshal(msg)
	if err != nil {
		log.Error("[SendMsg] Cannot Marshal")
		panic(err)
	}

	p.Push(context.Background(), &pb.ReqBody{Body:buf})
}