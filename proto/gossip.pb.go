// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gossip.proto

package gossip

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
	return fileDescriptor_878fa4887b90140c, []int{0}
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
	return fileDescriptor_878fa4887b90140c, []int{1}
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
	return fileDescriptor_878fa4887b90140c, []int{2}
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
	return fileDescriptor_878fa4887b90140c, []int{3}
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

func init() {
	proto.RegisterType((*ReqId)(nil), "ReqId")
	proto.RegisterType((*ReqBody)(nil), "ReqBody")
	proto.RegisterType((*Bool)(nil), "Bool")
	proto.RegisterType((*Void)(nil), "Void")
}

func init() { proto.RegisterFile("gossip.proto", fileDescriptor_878fa4887b90140c) }

var fileDescriptor_878fa4887b90140c = []byte{
	// 168 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x8e, 0x3d, 0x0f, 0x82, 0x40,
	0x0c, 0x40, 0x21, 0x39, 0x4e, 0xd2, 0x30, 0x75, 0xf0, 0x03, 0xa3, 0x31, 0x9d, 0x9c, 0x6e, 0xd0,
	0xd5, 0x89, 0xc5, 0xb8, 0x91, 0x1b, 0xdc, 0x21, 0x47, 0x84, 0x68, 0x52, 0xb0, 0x30, 0xf0, 0xef,
	0xcd, 0x5d, 0x70, 0x7b, 0xcd, 0x4b, 0xfb, 0x0a, 0xd9, 0x8b, 0x45, 0xba, 0xde, 0xf4, 0x5f, 0x1e,
	0x99, 0xf6, 0x90, 0xd8, 0x66, 0x78, 0x38, 0x44, 0x50, 0x6d, 0x25, 0xed, 0x36, 0x3e, 0xc5, 0xe7,
	0xcc, 0x06, 0xa6, 0x03, 0xac, 0x6c, 0x33, 0x14, 0xec, 0x66, 0xaf, 0x6b, 0x76, 0xf3, 0x5f, 0x7b,
	0xa6, 0x23, 0xa8, 0x82, 0xf9, 0x83, 0x6b, 0xd0, 0x32, 0x56, 0xe3, 0x24, 0xc1, 0xa6, 0x76, 0x99,
	0x48, 0x83, 0x7a, 0x72, 0xe7, 0x2e, 0x37, 0xd0, 0xf7, 0xd0, 0xc4, 0x0d, 0xa8, 0x92, 0xdf, 0x0d,
	0x6a, 0x13, 0xa2, 0x79, 0x62, 0xfc, 0x01, 0x8a, 0x70, 0x07, 0xaa, 0x9c, 0xa4, 0xc5, 0xd4, 0x2c,
	0xc1, 0x3c, 0x31, 0x7e, 0x97, 0xa2, 0x5a, 0x87, 0x47, 0xaf, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x17, 0xc7, 0x36, 0x71, 0xb8, 0x00, 0x00, 0x00,
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

// GossipServer is the server API for Gossip service.
type GossipServer interface {
	//checks whether hash already exists
	Poke(context.Context, *ReqId) (*Bool, error)
	//sends new info
	Push(context.Context, *ReqBody) (*Void, error)
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
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gossip.proto",
}
