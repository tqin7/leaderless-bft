package types

type PbftReq struct {
	ClientID   string;
	Operation  string;
	Timestamp  int64;
	SequenceID int64;
	MsgType    string;
}

type PrePrepareMsg struct {
	ViewID     int64;
	SequenceID int64;
	Digest     string;
	Req        *PbftReq;
	MsgType    string;
}

type PrepareMsg struct {
	ViewID     int64;
	SequenceID int64;
	Digest     string;
	NodeID     string;
	MsgType    string;
}

type CommitMsg struct {
	ViewID     int64;
	SequenceID int64;
	Digest     string;
	NodeID     string;
	MsgType    string;
}

type ReplyMsg struct {
	ViewID    int64;
	Timestamp int64;
	ClientID  string;
	NodeID    string;
	Result    string;
	MsgType    string;
}