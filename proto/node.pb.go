// Code generated by protoc-gen-go. DO NOT EDIT.
// source: node.proto

package node

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Requests struct {
	Requests             []string `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Requests) Reset()         { *m = Requests{} }
func (m *Requests) String() string { return proto.CompactTextString(m) }
func (*Requests) ProtoMessage()    {}
func (*Requests) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{0}
}

func (m *Requests) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Requests.Unmarshal(m, b)
}
func (m *Requests) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Requests.Marshal(b, m, deterministic)
}
func (m *Requests) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Requests.Merge(m, src)
}
func (m *Requests) XXX_Size() int {
	return xxx_messageInfo_Requests.Size(m)
}
func (m *Requests) XXX_DiscardUnknown() {
	xxx_messageInfo_Requests.DiscardUnknown(m)
}

var xxx_messageInfo_Requests proto.InternalMessageInfo

func (m *Requests) GetRequests() []string {
	if m != nil {
		return m.Requests
	}
	return nil
}

type ReqId struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqId) Reset()         { *m = ReqId{} }
func (m *ReqId) String() string { return proto.CompactTextString(m) }
func (*ReqId) ProtoMessage()    {}
func (*ReqId) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{1}
}

func (m *ReqId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqId.Unmarshal(m, b)
}
func (m *ReqId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqId.Marshal(b, m, deterministic)
}
func (m *ReqId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqId.Merge(m, src)
}
func (m *ReqId) XXX_Size() int {
	return xxx_messageInfo_ReqId.Size(m)
}
func (m *ReqId) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqId.DiscardUnknown(m)
}

var xxx_messageInfo_ReqId proto.InternalMessageInfo

func (m *ReqId) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

type ReqBody struct {
	Body                 []byte   `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqBody) Reset()         { *m = ReqBody{} }
func (m *ReqBody) String() string { return proto.CompactTextString(m) }
func (*ReqBody) ProtoMessage()    {}
func (*ReqBody) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{2}
}

func (m *ReqBody) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqBody.Unmarshal(m, b)
}
func (m *ReqBody) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqBody.Marshal(b, m, deterministic)
}
func (m *ReqBody) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqBody.Merge(m, src)
}
func (m *ReqBody) XXX_Size() int {
	return xxx_messageInfo_ReqBody.Size(m)
}
func (m *ReqBody) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqBody.DiscardUnknown(m)
}

var xxx_messageInfo_ReqBody proto.InternalMessageInfo

func (m *ReqBody) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type Bool struct {
	Status               bool     `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Bool) Reset()         { *m = Bool{} }
func (m *Bool) String() string { return proto.CompactTextString(m) }
func (*Bool) ProtoMessage()    {}
func (*Bool) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{3}
}

func (m *Bool) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Bool.Unmarshal(m, b)
}
func (m *Bool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Bool.Marshal(b, m, deterministic)
}
func (m *Bool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bool.Merge(m, src)
}
func (m *Bool) XXX_Size() int {
	return xxx_messageInfo_Bool.Size(m)
}
func (m *Bool) XXX_DiscardUnknown() {
	xxx_messageInfo_Bool.DiscardUnknown(m)
}

var xxx_messageInfo_Bool proto.InternalMessageInfo

func (m *Bool) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

type Void struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Void) Reset()         { *m = Void{} }
func (m *Void) String() string { return proto.CompactTextString(m) }
func (*Void) ProtoMessage()    {}
func (*Void) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{4}
}

func (m *Void) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Void.Unmarshal(m, b)
}
func (m *Void) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Void.Marshal(b, m, deterministic)
}
func (m *Void) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Void.Merge(m, src)
}
func (m *Void) XXX_Size() int {
	return xxx_messageInfo_Void.Size(m)
}
func (m *Void) XXX_DiscardUnknown() {
	xxx_messageInfo_Void.DiscardUnknown(m)
}

var xxx_messageInfo_Void proto.InternalMessageInfo

type SeqNumMsg struct {
	SeqNum               int64    `protobuf:"varint,1,opt,name=seq_num,json=seqNum,proto3" json:"seq_num,omitempty"`
	ReqHash              []byte   `protobuf:"bytes,2,opt,name=req_hash,json=reqHash,proto3" json:"req_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SeqNumMsg) Reset()         { *m = SeqNumMsg{} }
func (m *SeqNumMsg) String() string { return proto.CompactTextString(m) }
func (*SeqNumMsg) ProtoMessage()    {}
func (*SeqNumMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{5}
}

func (m *SeqNumMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SeqNumMsg.Unmarshal(m, b)
}
func (m *SeqNumMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SeqNumMsg.Marshal(b, m, deterministic)
}
func (m *SeqNumMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SeqNumMsg.Merge(m, src)
}
func (m *SeqNumMsg) XXX_Size() int {
	return xxx_messageInfo_SeqNumMsg.Size(m)
}
func (m *SeqNumMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_SeqNumMsg.DiscardUnknown(m)
}

var xxx_messageInfo_SeqNumMsg proto.InternalMessageInfo

func (m *SeqNumMsg) GetSeqNum() int64 {
	if m != nil {
		return m.SeqNum
	}
	return 0
}

func (m *SeqNumMsg) GetReqHash() []byte {
	if m != nil {
		return m.ReqHash
	}
	return nil
}

func init() {
	proto.RegisterType((*Requests)(nil), "Requests")
	proto.RegisterType((*ReqId)(nil), "ReqId")
	proto.RegisterType((*ReqBody)(nil), "ReqBody")
	proto.RegisterType((*Bool)(nil), "Bool")
	proto.RegisterType((*Void)(nil), "Void")
	proto.RegisterType((*SeqNumMsg)(nil), "SeqNumMsg")
}

func init() { proto.RegisterFile("node.proto", fileDescriptor_0c843d59d2d938e7) }

var fileDescriptor_0c843d59d2d938e7 = []byte{
	// 316 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xc1, 0x4f, 0xfa, 0x30,
	0x14, 0xc7, 0xc7, 0x8f, 0xd1, 0x8d, 0x97, 0x5f, 0x3c, 0xbc, 0x83, 0xc0, 0x10, 0x43, 0x8a, 0x1a,
	0x4e, 0x3d, 0xe0, 0x1f, 0x60, 0xe4, 0x32, 0x4d, 0xc0, 0x90, 0x92, 0x70, 0x25, 0x5b, 0xf6, 0x74,
	0xc4, 0xb1, 0x52, 0xda, 0xc5, 0xf0, 0xdf, 0x9b, 0x55, 0x20, 0x1e, 0xe4, 0xf6, 0x7d, 0xfd, 0xbc,
	0xf4, 0xfb, 0x49, 0x0b, 0x50, 0xaa, 0x8c, 0xc4, 0x6e, 0xaf, 0xac, 0xe2, 0x0f, 0x10, 0x4a, 0xd2,
	0x15, 0x19, 0x6b, 0x30, 0x82, 0x70, 0x7f, 0xcc, 0xdd, 0xc6, 0xb0, 0x39, 0x6e, 0xcb, 0xf3, 0xcc,
	0xfb, 0xd0, 0x92, 0xa4, 0x5f, 0x33, 0x44, 0xf0, 0xf3, 0xc4, 0xe4, 0xdd, 0xc6, 0xb0, 0x31, 0xfe,
	0x2f, 0x5d, 0xe6, 0x03, 0x08, 0x24, 0xe9, 0xa9, 0xca, 0x0e, 0x35, 0x4e, 0x55, 0x76, 0x38, 0xe1,
	0x3a, 0xf3, 0x5b, 0xf0, 0xa7, 0x4a, 0x15, 0x78, 0x0d, 0xcc, 0xd8, 0xc4, 0x56, 0xc6, 0xd1, 0x50,
	0x1e, 0x27, 0xce, 0xc0, 0x5f, 0xa9, 0x4d, 0xc6, 0x9f, 0xa0, 0xbd, 0x24, 0xfd, 0x56, 0x6d, 0xe7,
	0xe6, 0x03, 0x3b, 0x10, 0x18, 0xd2, 0xeb, 0xb2, 0xda, 0xba, 0xed, 0xa6, 0x64, 0xc6, 0x31, 0xec,
	0x39, 0xcb, 0xb5, 0x93, 0xf8, 0xe7, 0x5a, 0x82, 0x3d, 0xe9, 0x97, 0xc4, 0xe4, 0x93, 0x14, 0x58,
	0xac, 0x8c, 0xd9, 0xec, 0xb0, 0x03, 0xfe, 0x42, 0x7d, 0x12, 0x32, 0xe1, 0xac, 0xa3, 0x96, 0xa8,
	0x0d, 0xb8, 0x87, 0x3d, 0xf0, 0x17, 0x95, 0xc9, 0x31, 0x14, 0x47, 0xe3, 0xa8, 0x25, 0x5c, 0xb9,
	0x87, 0x77, 0x70, 0x15, 0x93, 0x7d, 0x2e, 0x8a, 0xf3, 0x83, 0xfc, 0xa0, 0xa8, 0x2d, 0x4e, 0x27,
	0xdc, 0x9b, 0xcc, 0x21, 0x5c, 0x96, 0xea, 0x2b, 0x4d, 0x8a, 0x02, 0x47, 0x10, 0xc4, 0x64, 0x57,
	0xca, 0x12, 0x82, 0x38, 0xab, 0x47, 0xbf, 0x32, 0xf7, 0xf0, 0x06, 0x82, 0x25, 0x95, 0x99, 0x24,
	0xfd, 0x47, 0xe9, 0x64, 0x04, 0xfe, 0x22, 0x7d, 0xb7, 0xd8, 0x07, 0x16, 0x93, 0xbd, 0xb0, 0x74,
	0x0f, 0xfe, 0xac, 0x5e, 0x1a, 0x40, 0x38, 0xbb, 0x7c, 0x57, 0xca, 0xdc, 0x97, 0x3e, 0x7e, 0x07,
	0x00, 0x00, 0xff, 0xff, 0xf6, 0xe8, 0x76, 0x0e, 0xe0, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GossipClient is the client API for Gossip service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GossipClient interface {
	//checks whether hash already exists
	Poke(ctx context.Context, in *ReqId, opts ...grpc.CallOption) (*Bool, error)
	//sends new info
	Push(ctx context.Context, in *ReqBody, opts ...grpc.CallOption) (*Void, error)
	//gets requests
	GetAllRequests(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Requests, error)
}

type gossipClient struct {
	cc *grpc.ClientConn
}

func NewGossipClient(cc *grpc.ClientConn) GossipClient {
	return &gossipClient{cc}
}

func (c *gossipClient) Poke(ctx context.Context, in *ReqId, opts ...grpc.CallOption) (*Bool, error) {
	out := new(Bool)
	err := c.cc.Invoke(ctx, "/Gossip/Poke", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gossipClient) Push(ctx context.Context, in *ReqBody, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/Gossip/Push", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gossipClient) GetAllRequests(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Requests, error) {
	out := new(Requests)
	err := c.cc.Invoke(ctx, "/Gossip/GetAllRequests", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GossipServer is the server API for Gossip service.
type GossipServer interface {
	//checks whether hash already exists
	Poke(context.Context, *ReqId) (*Bool, error)
	//sends new info
	Push(context.Context, *ReqBody) (*Void, error)
	//gets requests
	GetAllRequests(context.Context, *Void) (*Requests, error)
}

// UnimplementedGossipServer can be embedded to have forward compatible implementations.
type UnimplementedGossipServer struct {
}

func (*UnimplementedGossipServer) Poke(ctx context.Context, req *ReqId) (*Bool, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Poke not implemented")
}
func (*UnimplementedGossipServer) Push(ctx context.Context, req *ReqBody) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Push not implemented")
}
func (*UnimplementedGossipServer) GetAllRequests(ctx context.Context, req *Void) (*Requests, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllRequests not implemented")
}

func RegisterGossipServer(s *grpc.Server, srv GossipServer) {
	s.RegisterService(&_Gossip_serviceDesc, srv)
}

func _Gossip_Poke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GossipServer).Poke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Gossip/Poke",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GossipServer).Poke(ctx, req.(*ReqId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gossip_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqBody)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GossipServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Gossip/Push",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GossipServer).Push(ctx, req.(*ReqBody))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gossip_GetAllRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GossipServer).GetAllRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Gossip/GetAllRequests",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GossipServer).GetAllRequests(ctx, req.(*Void))
	}
	return interceptor(ctx, in, info, handler)
}

var _Gossip_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Gossip",
	HandlerType: (*GossipServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Poke",
			Handler:    _Gossip_Poke_Handler,
		},
		{
			MethodName: "Push",
			Handler:    _Gossip_Push_Handler,
		},
		{
			MethodName: "GetAllRequests",
			Handler:    _Gossip_GetAllRequests_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "node.proto",
}

// SnowballClient is the client API for Snowball service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SnowballClient interface {
	//get vote from this node
	GetVote(ctx context.Context, in *SeqNumMsg, opts ...grpc.CallOption) (*SeqNumMsg, error)
	//a client sends request using this method
	SendReq(ctx context.Context, in *ReqBody, opts ...grpc.CallOption) (*Void, error)
}

type snowballClient struct {
	cc *grpc.ClientConn
}

func NewSnowballClient(cc *grpc.ClientConn) SnowballClient {
	return &snowballClient{cc}
}

func (c *snowballClient) GetVote(ctx context.Context, in *SeqNumMsg, opts ...grpc.CallOption) (*SeqNumMsg, error) {
	out := new(SeqNumMsg)
	err := c.cc.Invoke(ctx, "/Snowball/GetVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *snowballClient) SendReq(ctx context.Context, in *ReqBody, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/Snowball/SendReq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SnowballServer is the server API for Snowball service.
type SnowballServer interface {
	//get vote from this node
	GetVote(context.Context, *SeqNumMsg) (*SeqNumMsg, error)
	//a client sends request using this method
	SendReq(context.Context, *ReqBody) (*Void, error)
}

// UnimplementedSnowballServer can be embedded to have forward compatible implementations.
type UnimplementedSnowballServer struct {
}

func (*UnimplementedSnowballServer) GetVote(ctx context.Context, req *SeqNumMsg) (*SeqNumMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVote not implemented")
}
func (*UnimplementedSnowballServer) SendReq(ctx context.Context, req *ReqBody) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendReq not implemented")
}

func RegisterSnowballServer(s *grpc.Server, srv SnowballServer) {
	s.RegisterService(&_Snowball_serviceDesc, srv)
}

func _Snowball_GetVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SeqNumMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SnowballServer).GetVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Snowball/GetVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SnowballServer).GetVote(ctx, req.(*SeqNumMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _Snowball_SendReq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqBody)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SnowballServer).SendReq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Snowball/SendReq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SnowballServer).SendReq(ctx, req.(*ReqBody))
	}
	return interceptor(ctx, in, info, handler)
}

var _Snowball_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Snowball",
	HandlerType: (*SnowballServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVote",
			Handler:    _Snowball_GetVote_Handler,
		},
		{
			MethodName: "SendReq",
			Handler:    _Snowball_SendReq_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "node.proto",
}

// PbftClient is the client API for Pbft service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PbftClient interface {
	// client sends request
	GetReq(ctx context.Context, in *ReqBody, opts ...grpc.CallOption) (*Void, error)
}

type pbftClient struct {
	cc *grpc.ClientConn
}

func NewPbftClient(cc *grpc.ClientConn) PbftClient {
	return &pbftClient{cc}
}

func (c *pbftClient) GetReq(ctx context.Context, in *ReqBody, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/Pbft/GetReq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PbftServer is the server API for Pbft service.
type PbftServer interface {
	// client sends request
	GetReq(context.Context, *ReqBody) (*Void, error)
}

// UnimplementedPbftServer can be embedded to have forward compatible implementations.
type UnimplementedPbftServer struct {
}

func (*UnimplementedPbftServer) GetReq(ctx context.Context, req *ReqBody) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReq not implemented")
}

func RegisterPbftServer(s *grpc.Server, srv PbftServer) {
	s.RegisterService(&_Pbft_serviceDesc, srv)
}

func _Pbft_GetReq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqBody)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PbftServer).GetReq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Pbft/GetReq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PbftServer).GetReq(ctx, req.(*ReqBody))
	}
	return interceptor(ctx, in, info, handler)
}

var _Pbft_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Pbft",
	HandlerType: (*PbftServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetReq",
			Handler:    _Pbft_GetReq_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "node.proto",
}

// LbftClient is the client API for Lbft service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LbftClient interface {
	LSendReq(ctx context.Context, in *ReqBody, opts ...grpc.CallOption) (*Void, error)
}

type lbftClient struct {
	cc *grpc.ClientConn
}

func NewLbftClient(cc *grpc.ClientConn) LbftClient {
	return &lbftClient{cc}
}

func (c *lbftClient) LSendReq(ctx context.Context, in *ReqBody, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/Lbft/LSendReq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LbftServer is the server API for Lbft service.
type LbftServer interface {
	LSendReq(context.Context, *ReqBody) (*Void, error)
}

// UnimplementedLbftServer can be embedded to have forward compatible implementations.
type UnimplementedLbftServer struct {
}

func (*UnimplementedLbftServer) LSendReq(ctx context.Context, req *ReqBody) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LSendReq not implemented")
}

func RegisterLbftServer(s *grpc.Server, srv LbftServer) {
	s.RegisterService(&_Lbft_serviceDesc, srv)
}

func _Lbft_LSendReq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqBody)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LbftServer).LSendReq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Lbft/LSendReq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LbftServer).LSendReq(ctx, req.(*ReqBody))
	}
	return interceptor(ctx, in, info, handler)
}

var _Lbft_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Lbft",
	HandlerType: (*LbftServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LSendReq",
			Handler:    _Lbft_LSendReq_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "node.proto",
}
