package proto

import (
	"encoding/json"
	"errors"
	"strings"
	tp "github.com/spockqin/leaderless-bft/types"
	"time"
	log "github.com/sirupsen/logrus"
)

func (p *Pbfter) GetMsgFromGossip(){
	for {
		p.requestsLock.Lock()

		for _, req := range p.requests {
			if strings.Contains(req, "PrePrepareMsg") {
				log.WithFields(log.Fields{
					"ip": p.ip,
				}).Info("got PrePrepareMsg")
				var prePrepareMsg tp.PrePrepareMsg
				err := json.Unmarshal([]byte(req), &prePrepareMsg)
				if err != nil {	
					panic(errors.New("[GetMsgFromGossip] unmarshal PrePrepareMsg error"))
				}
				p.RouteMsg(&prePrepareMsg)
			} else if strings.Contains(req, "PrepareMsg") {
				log.WithFields(log.Fields{
					"ip": p.ip,
				}).Info("got PrepareMsg")
				var prePareMsg tp.PrepareMsg
				err := json.Unmarshal([]byte(req), &prePareMsg)
				if err != nil {
					panic(errors.New("[GetMsgFromGossip] unmarshal PrepareMsg error"))
				}
				p.RouteMsg(&prePareMsg)
			} else if strings.Contains(req, "CommitMsg") {
				log.WithFields(log.Fields{
					"ip": p.ip,
				}).Info("got CommitMsg")
				var commitMsg tp.CommitMsg
				err := json.Unmarshal([]byte(req), &commitMsg)
				if err != nil {
					panic(errors.New("[GetMsgFromGossip] unmarshal CommitMsg error"))
				}
				p.RouteMsg(&commitMsg)
			} else if strings.Contains(req, "ReplyMsg") {
				log.WithFields(log.Fields{
					"ip": p.ip,
				}).Info("got ReplyMsg")
				var replyMsg tp.ReplyMsg
				err := json.Unmarshal([]byte(req), &replyMsg)
				if err != nil {
					panic(errors.New("[GetMsgFromGossip] unmarshal ReplyMsg error"))
				}
				p.RouteMsg(&replyMsg)
			}
		}

		p.requests = make([]string, 0)
		p.requestsLock.Unlock()
		time.Sleep(3 * time.Second)
	}
}

func (p *Pbfter) RouteMsg(msg interface{}) []error{
	switch msg.(type) {
	case *tp.PrePrepareMsg:
		if p.CurrentState == nil {
			msgs := make([]*tp.PrePrepareMsg, len(p.MsgBuffer.PrePrepareMsgs))
			copy(msgs, p.MsgBuffer.PrePrepareMsgs)
			msgs = append(msgs, msg.(*tp.PrePrepareMsg))
			p.MsgBuffer.PrePrepareMsgs = make([]*tp.PrePrepareMsg, 0)
			p.MsgDelivery <- msgs
		} else {
			p.MsgBuffer.PrePrepareMsgs = append(p.MsgBuffer.PrePrepareMsgs, msg.(*tp.PrePrepareMsg))
		}
	case *tp.PrepareMsg:
		if p.CurrentState == nil || p.CurrentState.CurrentStage != PrePrepared {
			p.MsgBuffer.PrepareMsgs = append(p.MsgBuffer.PrepareMsgs, msg.(*tp.PrepareMsg))
		} else {
			msgs := make([]*tp.PrepareMsg, len(p.MsgBuffer.PrepareMsgs))
			copy(msgs, p.MsgBuffer.PrepareMsgs)
			msgs = append(msgs, msg.(*tp.PrepareMsg))
			p.MsgBuffer.PrepareMsgs = make([]*tp.PrepareMsg, 0)
			p.MsgDelivery <- msgs
		}
	case *tp.CommitMsg:
		if p.CurrentState == nil || p.CurrentState.CurrentStage != Prepared {
			p.MsgBuffer.CommitMsgs = append(p.MsgBuffer.CommitMsgs, msg.(*tp.CommitMsg))
		} else {
			msgs := make([]*tp.CommitMsg, len(p.MsgBuffer.CommitMsgs))
			copy(msgs, p.MsgBuffer.CommitMsgs)
			msgs = append(msgs, msg.(*tp.CommitMsg))
			p.MsgBuffer.CommitMsgs = make([]*tp.CommitMsg, 0)
			p.MsgDelivery <- msgs
		}
	}
	return nil
}