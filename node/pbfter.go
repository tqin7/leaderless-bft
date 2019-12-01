package proto

import (
	"context"
	pb "github.com/spockqin/leaderless-bft/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Stage int
const (
	Idle        Stage = iota // Node is created successfully, but the consensus process is not started yet.
	PrePrepared              // The ReqMsgs is processed successfully. The node is ready to head to the Prepare stage.
	Prepared                 // Same with `prepared` stage explained in the original paper.
	Committed                // Same with `committed-local` stage explained in the original paper.
)

type MsgLogs struct {
	ReqMsg        *ReqMsg
	PrepareMsgs   map[string]*PrePareMsg
	CommitMsgs    map[string]*CommitMsg
}

type State struct {
	MsgLogs        *MsgLogs
	LastSequenceID int64
	CurrentStage   Stage
}

type Pbfter struct {
	Gossiper
	NodeID string
	ViewID int64
	CurrentState *State
	CommittedMsgs []ReqMsg
}

func (p *Pbfter) SendReq(ctx context.Context, req *pb.PbftReq) (*pb.Void, error) {

	p.Push(ctx, &pb.ReqBody{Body:req})

	s.Push(ctx, req)

	s.setProposal()
	go s.performQueries()

	return &pb.Void{}, nil
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