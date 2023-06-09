// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc_stdio.proto

package plugin

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import empty "github.com/golang/protobuf/ptypes/empty"

import (
	context "context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type StdioData_Channel int32

const (
	StdioData_INVALID StdioData_Channel = 0
	StdioData_STDOUT  StdioData_Channel = 1
	StdioData_STDERR  StdioData_Channel = 2
)

var StdioData_Channel_name = map[int32]string{
	0: "INVALID",
	1: "STDOUT",
	2: "STDERR",
}
var StdioData_Channel_value = map[string]int32{
	"INVALID": 0,
	"STDOUT":  1,
	"STDERR":  2,
}

func (x StdioData_Channel) String() string {
	return proto.EnumName(StdioData_Channel_name, int32(x))
}
func (StdioData_Channel) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_grpc_stdio_db2934322ca63bd5, []int{0, 0}
}

// StdioData is a single chunk of stdout or stderr data that is streamed
// from GRPCStdio.
type StdioData struct {
	Channel              StdioData_Channel `protobuf:"varint,1,opt,name=channel,proto3,enum=plugin.StdioData_Channel" json:"channel,omitempty"`
	Data                 []byte            `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *StdioData) Reset()         { *m = StdioData{} }
func (m *StdioData) String() string { return proto.CompactTextString(m) }
func (*StdioData) ProtoMessage()    {}
func (*StdioData) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_stdio_db2934322ca63bd5, []int{0}
}
func (m *StdioData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StdioData.Unmarshal(m, b)
}
func (m *StdioData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StdioData.Marshal(b, m, deterministic)
}
func (dst *StdioData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StdioData.Merge(dst, src)
}
func (m *StdioData) XXX_Size() int {
	return xxx_messageInfo_StdioData.Size(m)
}
func (m *StdioData) XXX_DiscardUnknown() {
	xxx_messageInfo_StdioData.DiscardUnknown(m)
}

var xxx_messageInfo_StdioData proto.InternalMessageInfo

func (m *StdioData) GetChannel() StdioData_Channel {
	if m != nil {
		return m.Channel
	}
	return StdioData_INVALID
}

func (m *StdioData) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*StdioData)(nil), "plugin.StdioData")
	proto.RegisterEnum("plugin.StdioData_Channel", StdioData_Channel_name, StdioData_Channel_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GRPCStdioClient is the client API for GRPCStdio service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GRPCStdioClient interface {
	// StreamStdio returns a stream that contains all the stdout/stderr.
	// This RPC endpoint must only be called ONCE. Once stdio data is consumed
	// it is not sent again.
	//
	// Callers should connect early to prevent blocking on the plugin process.
	StreamStdio(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (GRPCStdio_StreamStdioClient, error)
}

type gRPCStdioClient struct {
	cc *grpc.ClientConn
}

func NewGRPCStdioClient(cc *grpc.ClientConn) GRPCStdioClient {
	return &gRPCStdioClient{cc}
}

func (c *gRPCStdioClient) StreamStdio(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (GRPCStdio_StreamStdioClient, error) {
	stream, err := c.cc.NewStream(ctx, &_GRPCStdio_serviceDesc.Streams[0], "/plugin.GRPCStdio/StreamStdio", opts...)
	if err != nil {
		return nil, err
	}
	x := &gRPCStdioStreamStdioClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GRPCStdio_StreamStdioClient interface {
	Recv() (*StdioData, error)
	grpc.ClientStream
}

type gRPCStdioStreamStdioClient struct {
	grpc.ClientStream
}

func (x *gRPCStdioStreamStdioClient) Recv() (*StdioData, error) {
	m := new(StdioData)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GRPCStdioServer is the server API for GRPCStdio service.
type GRPCStdioServer interface {
	// StreamStdio returns a stream that contains all the stdout/stderr.
	// This RPC endpoint must only be called ONCE. Once stdio data is consumed
	// it is not sent again.
	//
	// Callers should connect early to prevent blocking on the plugin process.
	StreamStdio(*empty.Empty, GRPCStdio_StreamStdioServer) error
}

func RegisterGRPCStdioServer(s *grpc.Server, srv GRPCStdioServer) {
	s.RegisterService(&_GRPCStdio_serviceDesc, srv)
}

func _GRPCStdio_StreamStdio_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GRPCStdioServer).StreamStdio(m, &gRPCStdioStreamStdioServer{stream})
}

type GRPCStdio_StreamStdioServer interface {
	Send(*StdioData) error
	grpc.ServerStream
}

type gRPCStdioStreamStdioServer struct {
	grpc.ServerStream
}

func (x *gRPCStdioStreamStdioServer) Send(m *StdioData) error {
	return x.ServerStream.SendMsg(m)
}

var _GRPCStdio_serviceDesc = grpc.ServiceDesc{
	ServiceName: "plugin.GRPCStdio",
	HandlerType: (*GRPCStdioServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamStdio",
			Handler:       _GRPCStdio_StreamStdio_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "grpc_stdio.proto",
}

func init() { proto.RegisterFile("grpc_stdio.proto", fileDescriptor_grpc_stdio_db2934322ca63bd5) }

var fileDescriptor_grpc_stdio_db2934322ca63bd5 = []byte{
	// 221 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x48, 0x2f, 0x2a, 0x48,
	0x8e, 0x2f, 0x2e, 0x49, 0xc9, 0xcc, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2b, 0xc8,
	0x29, 0x4d, 0xcf, 0xcc, 0x93, 0x92, 0x4e, 0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5, 0x07, 0x8b, 0x26,
	0x95, 0xa6, 0xe9, 0xa7, 0xe6, 0x16, 0x94, 0x54, 0x42, 0x14, 0x29, 0xb5, 0x30, 0x72, 0x71, 0x06,
	0x83, 0x34, 0xb9, 0x24, 0x96, 0x24, 0x0a, 0x19, 0x73, 0xb1, 0x27, 0x67, 0x24, 0xe6, 0xe5, 0xa5,
	0xe6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0xf0, 0x19, 0x49, 0xea, 0x41, 0x0c, 0xd1, 0x83, 0xab, 0xd1,
	0x73, 0x86, 0x28, 0x08, 0x82, 0xa9, 0x14, 0x12, 0xe2, 0x62, 0x49, 0x49, 0x2c, 0x49, 0x94, 0x60,
	0x52, 0x60, 0xd4, 0xe0, 0x09, 0x02, 0xb3, 0x95, 0xf4, 0xb8, 0xd8, 0xa1, 0xea, 0x84, 0xb8, 0xb9,
	0xd8, 0x3d, 0xfd, 0xc2, 0x1c, 0x7d, 0x3c, 0x5d, 0x04, 0x18, 0x84, 0xb8, 0xb8, 0xd8, 0x82, 0x43,
	0x5c, 0xfc, 0x43, 0x43, 0x04, 0x18, 0xa1, 0x6c, 0xd7, 0xa0, 0x20, 0x01, 0x26, 0x23, 0x77, 0x2e,
	0x4e, 0xf7, 0xa0, 0x00, 0x67, 0xb0, 0x2d, 0x42, 0x56, 0x5c, 0xdc, 0xc1, 0x25, 0x45, 0xa9, 0x89,
	0xb9, 0x10, 0xae, 0x98, 0x1e, 0xc4, 0x03, 0x7a, 0x30, 0x0f, 0xe8, 0xb9, 0x82, 0x3c, 0x20, 0x25,
	0x88, 0xe1, 0x36, 0x03, 0x46, 0x27, 0x8e, 0x28, 0xa8, 0xb7, 0x93, 0xd8, 0xc0, 0xca, 0x8d, 0x01,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x5d, 0xbb, 0xe0, 0x69, 0x19, 0x01, 0x00, 0x00,
}
