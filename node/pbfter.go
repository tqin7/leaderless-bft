package proto

import (
	"context"
	"encoding/json"
	"errors"
	pb "github.com/spockqin/leaderless-bft/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"fmt"
	"time"
	"github.com/spockqin/leaderless-bft/util"
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

func (p *Pbfter) SendReq(ctx context.Context, req *pb.PbftReq) (*pb.Void, error) {
	LogMsg(req)

	// begin new consensus
	if p.CurrentState != nil {
		return &pb.Void{}, errors.New("another consensus is ongoing")
	}

	var lastSeqID int64
	if len(p.CommittedMsgs) == 0 {
		lastSeqID = -1
	} else {
		lastSeqID = p.CommittedMsgs[len(p.CommittedMsgs) - 1].SequenceID
	}

	p.CurrentState = createState(p.ViewID, lastSeqID)

	LogStage("Prepare for Consensus", true)

	prePrepareMsg, err := p.CurrentState

	//var buf []byte
	//buf, err := json.Marshal(req)
	//if err != nil {
	//	log.Error("SendReq Cannot Marshal [pbfter]")
	//	panic(err)
	//}
	//
	//p.Push(ctx, &pb.ReqBody{Body:buf})
	//
	//return &pb.Void{}, nil
}
// node sends prePrepare msg
func (p *Pbfter) SendPrePrepare(ctx context.Context, msg *pb.PrePrepareMsg) (*pb.Void, error) {

}
// node sends prepare msg
func (p *Pbfter) SendPrepare(ctx context.Context, msg *pb.PrepareMsg) (*pb.Void, error) {

}
// node sends commit msg
func (p *Pbfter) SendCommit(ctx context.Context, msg *pb.CommitMsg) (*pb.Void, error) {

}
// node sends reply msg to the client
func (p *Pbfter) SendReply(ctx context.Context, msg *pb.ReplyMsg) (*pb.Void, error) {

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

	digest, err := util.digest

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