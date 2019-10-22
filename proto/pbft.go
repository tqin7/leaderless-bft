package node

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
	log "github.com/sirupsen/logrus"
)

type State struct {
	ViewID         int64
	MsgLogs        *MsgLogs
	CurrentStage   Stage
	LastSequenceID int64
	f 			   int
}

type MsgLogs struct {
	ReqMsg        *RequestMsg
	PrepareMsgs   map[string]*VoteMsg
	CommitMsgs    map[string]*VoteMsg
}

type Stage int
const (
	Idle        Stage = iota // Node is created successfully, but the consensus process is not started yet.
	PrePrepared              // The ReqMsgs is processed successfully. The node is ready to head to the Prepare stage.
	Prepared                 // Same with `prepared` stage explained in the original paper.
	Committed                // Same with `committed-local` stage explained in the original paper.
)

type RequestMsg struct {
	Timestamp  int64  `json:"timestamp"`
	ClientID   string `json:"clientID"`
	Operation  string `json:"operation"`
	SequenceID int64  `json:"sequenceID"`
}

type ReplyMsg struct {
	ViewID    int64  `json:"viewID"`
	Timestamp int64  `json:"timestamp"`
	ClientID  string `json:"clientID"`
	NodeID    string `json:"nodeID"`
	Result    string `json:"result"`
}

type PrePrepareMsg struct {
	ViewID     int64       `json:"viewID"`
	SequenceID int64       `json:"sequenceID"`
	Digest     string      `json:"digest"`
	RequestMsg *RequestMsg `json:"requestMsg"`
}

type VoteMsg struct {
	ViewID     int64  `json:"viewID"`
	SequenceID int64  `json:"sequenceID"`
	Digest     string `json:"digest"`
	NodeID     string `json:"nodeID"`
	MsgType           `json:"msgType"`
}

type MsgType int
const (
	PrepareMsg MsgType = iota
	CommitMsg
)

func CreateState(viewID int64, lastSequenceID int64, f int) *State{
	state := &State{
		ViewID:  viewID,
		MsgLogs: &MsgLogs{
			ReqMsg:      nil,
			PrepareMsgs: make(map[string]*VoteMsg),
			CommitMsgs:  make(map[string]*VoteMsg),
		},
		LastSequenceID: lastSequenceID,
		CurrentStage: Idle,
		f: f,
	}
	return state
}

func (state * State) StartConsensus(request *RequestMsg) (*PrePrepareMsg, error){
	// sequenceID is the index of the message
	sequenceID := time.Now().UnixNano()

	// get unique and largest sequence ID
	if state.LastSequenceID != -1 {
		if state.LastSequenceID >= sequenceID {
			sequenceID = state.LastSequenceID + 1
		}
	}

	request.SequenceID = sequenceID

	state.MsgLogs.ReqMsg = request

	// get the hash
	digest, err := digest(request)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Hash error in startConsensus\n")
		return nil, err
	}

	// change the stage to pre-prepared
	state.CurrentStage = PrePrepared

	return &PrePrepareMsg{
		ViewID:     state.ViewID,
		SequenceID: sequenceID,
		Digest:     digest,
		RequestMsg: request,
	}, nil
}

func (state *State) PrePrepare(prePrepareMsg * PrePrepareMsg) (*VoteMsg, error) {
	// save the msg to logs
	state.MsgLogs.ReqMsg = prePrepareMsg.RequestMsg

	// verify view, sequenceID, digest are correct
	if !state.verifyMsg(prePrepareMsg.ViewID, prePrepareMsg.SequenceID, prePrepareMsg.Digest) {
		return nil, errors.New("pre-prepare message is corrupted")
	}

	// change the stage
	state.CurrentStage = PrePrepared

	return &VoteMsg{
		ViewID:     state.ViewID,
		SequenceID: prePrepareMsg.SequenceID,
		Digest:     prePrepareMsg.Digest,
		MsgType:    PrepareMsg,
	}, nil
}

func (state *State) Prepare(prePareMsg *VoteMsg) (*VoteMsg, error) {
	// verify corruption
	if !state.verifyMsg(prePareMsg.ViewID, prePareMsg.SequenceID, prePareMsg.Digest) {
		return nil, errors.New("prepare message is corrupted")
	}

	// append msg to logs
	state.MsgLogs.PrepareMsgs[prePareMsg.NodeID] = prePareMsg

	// receive no less than 2*f prepare messages
	if state.prepared() {
		state.CurrentStage = Prepared

		return &VoteMsg{
			ViewID:     state.ViewID,
			SequenceID: prePareMsg.SequenceID,
			Digest:     prePareMsg.Digest,
			MsgType:    CommitMsg,
		}, nil
	}

	return nil, nil
}

func (state *State) Commit(commitMsg *VoteMsg) (*ReplyMsg, *RequestMsg, error) {
	// verify corruption
	if !state.verifyMsg(commitMsg.ViewID, commitMsg.SequenceID, commitMsg.Digest) {
		return nil, nil, errors.New("prepare message is corrupted")
	}

	// append msg to logs
	state.MsgLogs.CommitMsgs[commitMsg.NodeID] = commitMsg

	if state.committed() {
		// execute the request locally and gets the result
		result := "Successfully Executed"

		state.CurrentStage = Committed

		return &ReplyMsg{
			ViewID:    state.ViewID,
			Timestamp: state.MsgLogs.ReqMsg.Timestamp,
			ClientID:  state.MsgLogs.ReqMsg.ClientID,
			Result:    result,
		}, state.MsgLogs.ReqMsg, nil
	}

	return nil, nil, nil
}

func (state *State) prepared() bool {
	if state.MsgLogs.ReqMsg == nil {
		return false
	}

	// check whether receive no less than 2*f prepare messages
	if len(state.MsgLogs.PrepareMsgs) < 2 * state.f {
		return false
	}

	return true
}

func (state *State) committed() bool {
	if !state.prepared() {
		return false
	}

	// check whether receive no less than 2*f commit messages
	if len(state.MsgLogs.CommitMsgs) < 2 * state.f {
		return false
	}

	return true
}

// verify view, sequenceID, digest are consistent between state and input
func (state *State) verifyMsg(viewID int64, sequenceID int64, digestGot string) bool {

	if state.ViewID != viewID {
		return false
	}

	if state.LastSequenceID != -1 {
		if state.LastSequenceID >= sequenceID {
			return false
		}
	}

	digest, err := digest(state.MsgLogs.ReqMsg)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Hash error in verifyMsg\n")
		return false
	}

	if digestGot != digest {
		return false
	}

	return true
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