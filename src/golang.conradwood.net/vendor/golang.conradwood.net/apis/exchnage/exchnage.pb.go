// Code generated by protoc-gen-go.
// source: golang.conradwood.net/apis/exchnage/exchnage.proto
// DO NOT EDIT!

/*
Package exchnage is a generated protocol buffer package.

It is generated from these files:
	golang.conradwood.net/apis/exchnage/exchnage.proto

It has these top-level messages:
	Calendar
	CalendarEntry
	CalendarRequest
	CalendarResponse
	CalendarConfig
*/
package exchnage

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
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

// a users calendar
type Calendar struct {
	Entries []*CalendarEntry `protobuf:"bytes,1,rep,name=Entries" json:"Entries,omitempty"`
	ID      string           `protobuf:"bytes,2,opt,name=ID" json:"ID,omitempty"`
}

func (m *Calendar) Reset()                    { *m = Calendar{} }
func (m *Calendar) String() string            { return proto.CompactTextString(m) }
func (*Calendar) ProtoMessage()               {}
func (*Calendar) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Calendar) GetEntries() []*CalendarEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *Calendar) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

// a single entry in a calendar. expect this to grow more complicated
type CalendarEntry struct {
	ID          string `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	Summary     string `protobuf:"bytes,2,opt,name=Summary" json:"Summary,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=Description" json:"Description,omitempty"`
	Start       uint32 `protobuf:"varint,4,opt,name=Start" json:"Start,omitempty"`
	End         uint32 `protobuf:"varint,5,opt,name=End" json:"End,omitempty"`
}

func (m *CalendarEntry) Reset()                    { *m = CalendarEntry{} }
func (m *CalendarEntry) String() string            { return proto.CompactTextString(m) }
func (*CalendarEntry) ProtoMessage()               {}
func (*CalendarEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CalendarEntry) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *CalendarEntry) GetSummary() string {
	if m != nil {
		return m.Summary
	}
	return ""
}

func (m *CalendarEntry) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *CalendarEntry) GetStart() uint32 {
	if m != nil {
		return m.Start
	}
	return 0
}

func (m *CalendarEntry) GetEnd() uint32 {
	if m != nil {
		return m.End
	}
	return 0
}

type CalendarRequest struct {
	CalendarID string `protobuf:"bytes,1,opt,name=CalendarID" json:"CalendarID,omitempty"`
}

func (m *CalendarRequest) Reset()                    { *m = CalendarRequest{} }
func (m *CalendarRequest) String() string            { return proto.CompactTextString(m) }
func (*CalendarRequest) ProtoMessage()               {}
func (*CalendarRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *CalendarRequest) GetCalendarID() string {
	if m != nil {
		return m.CalendarID
	}
	return ""
}

type CalendarResponse struct {
	Calendar    *Calendar `protobuf:"bytes,1,opt,name=Calendar" json:"Calendar,omitempty"`
	RetrievedAt uint32    `protobuf:"varint,2,opt,name=RetrievedAt" json:"RetrievedAt,omitempty"`
}

func (m *CalendarResponse) Reset()                    { *m = CalendarResponse{} }
func (m *CalendarResponse) String() string            { return proto.CompactTextString(m) }
func (*CalendarResponse) ProtoMessage()               {}
func (*CalendarResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *CalendarResponse) GetCalendar() *Calendar {
	if m != nil {
		return m.Calendar
	}
	return nil
}

func (m *CalendarResponse) GetRetrievedAt() uint32 {
	if m != nil {
		return m.RetrievedAt
	}
	return 0
}

type CalendarConfig struct {
	ID     uint64 `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	UserID string `protobuf:"bytes,2,opt,name=UserID" json:"UserID,omitempty"`
}

func (m *CalendarConfig) Reset()                    { *m = CalendarConfig{} }
func (m *CalendarConfig) String() string            { return proto.CompactTextString(m) }
func (*CalendarConfig) ProtoMessage()               {}
func (*CalendarConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *CalendarConfig) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *CalendarConfig) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func init() {
	proto.RegisterType((*Calendar)(nil), "exchnage.Calendar")
	proto.RegisterType((*CalendarEntry)(nil), "exchnage.CalendarEntry")
	proto.RegisterType((*CalendarRequest)(nil), "exchnage.CalendarRequest")
	proto.RegisterType((*CalendarResponse)(nil), "exchnage.CalendarResponse")
	proto.RegisterType((*CalendarConfig)(nil), "exchnage.CalendarConfig")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ExchnageService service

type ExchnageServiceClient interface {
	// get a users calendar from cache
	GetCalendar(ctx context.Context, in *CalendarRequest, opts ...grpc.CallOption) (*CalendarResponse, error)
}

type exchnageServiceClient struct {
	cc *grpc.ClientConn
}

func NewExchnageServiceClient(cc *grpc.ClientConn) ExchnageServiceClient {
	return &exchnageServiceClient{cc}
}

func (c *exchnageServiceClient) GetCalendar(ctx context.Context, in *CalendarRequest, opts ...grpc.CallOption) (*CalendarResponse, error) {
	out := new(CalendarResponse)
	err := grpc.Invoke(ctx, "/exchnage.ExchnageService/GetCalendar", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ExchnageService service

type ExchnageServiceServer interface {
	// get a users calendar from cache
	GetCalendar(context.Context, *CalendarRequest) (*CalendarResponse, error)
}

func RegisterExchnageServiceServer(s *grpc.Server, srv ExchnageServiceServer) {
	s.RegisterService(&_ExchnageService_serviceDesc, srv)
}

func _ExchnageService_GetCalendar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CalendarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchnageServiceServer).GetCalendar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/exchnage.ExchnageService/GetCalendar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchnageServiceServer).GetCalendar(ctx, req.(*CalendarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ExchnageService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "exchnage.ExchnageService",
	HandlerType: (*ExchnageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCalendar",
			Handler:    _ExchnageService_GetCalendar_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "golang.conradwood.net/apis/exchnage/exchnage.proto",
}

func init() { proto.RegisterFile("golang.conradwood.net/apis/exchnage/exchnage.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 339 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x52, 0x4d, 0x4b, 0xc3, 0x40,
	0x14, 0x24, 0xfd, 0xf6, 0x85, 0x7e, 0xb0, 0x88, 0xc6, 0x1e, 0x24, 0x44, 0x84, 0x9e, 0x52, 0x1a,
	0x2f, 0x5e, 0xb5, 0x29, 0xd2, 0x83, 0x97, 0x2d, 0xe2, 0x79, 0x4d, 0x9e, 0x31, 0xd0, 0xee, 0xc6,
	0xdd, 0x6d, 0xb5, 0x67, 0xff, 0xb8, 0x24, 0xcd, 0xa6, 0x11, 0x7b, 0xdb, 0x99, 0x9d, 0xe1, 0xbd,
	0x37, 0x0c, 0x04, 0x89, 0x58, 0x33, 0x9e, 0xf8, 0x91, 0xe0, 0x92, 0xc5, 0x5f, 0x42, 0xc4, 0x3e,
	0x47, 0x3d, 0x65, 0x59, 0xaa, 0xa6, 0xf8, 0x1d, 0x7d, 0x70, 0x96, 0x60, 0xf5, 0xf0, 0x33, 0x29,
	0xb4, 0x20, 0x3d, 0x83, 0xbd, 0x67, 0xe8, 0xcd, 0xd9, 0x1a, 0x79, 0xcc, 0x24, 0x99, 0x41, 0x77,
	0xc1, 0xb5, 0x4c, 0x51, 0x39, 0x96, 0xdb, 0x9c, 0xd8, 0xc1, 0xa5, 0x5f, 0xf9, 0x8c, 0x28, 0x17,
	0xec, 0xa9, 0xd1, 0x91, 0x01, 0x34, 0x96, 0xa1, 0xd3, 0x70, 0xad, 0xc9, 0x19, 0x6d, 0x2c, 0x43,
	0xef, 0xc7, 0x82, 0xfe, 0x1f, 0x69, 0xa9, 0xb0, 0x8c, 0x82, 0x38, 0xd0, 0x5d, 0x6d, 0x37, 0x1b,
	0x26, 0xf7, 0xa5, 0xcd, 0x40, 0xe2, 0x82, 0x1d, 0xa2, 0x8a, 0x64, 0x9a, 0xe9, 0x54, 0x70, 0xa7,
	0x59, 0xfc, 0xd6, 0x29, 0x72, 0x0e, 0xed, 0x95, 0x66, 0x52, 0x3b, 0x2d, 0xd7, 0x9a, 0xf4, 0xe9,
	0x01, 0x90, 0x11, 0x34, 0x17, 0x3c, 0x76, 0xda, 0x05, 0x97, 0x3f, 0xbd, 0x19, 0x0c, 0xcd, 0x12,
	0x14, 0x3f, 0xb7, 0xa8, 0x34, 0xb9, 0x06, 0x30, 0x54, 0xb5, 0x4e, 0x8d, 0xf1, 0x62, 0x18, 0x1d,
	0x2d, 0x2a, 0x13, 0x5c, 0x21, 0xf1, 0x8f, 0xd9, 0x14, 0x0e, 0x3b, 0x20, 0xff, 0x03, 0xa1, 0xc7,
	0xfc, 0x5c, 0xb0, 0x29, 0xe6, 0xb9, 0xec, 0x30, 0x7e, 0xd0, 0xc5, 0x79, 0x7d, 0x5a, 0xa7, 0xbc,
	0x7b, 0x18, 0x18, 0xf5, 0x5c, 0xf0, 0xf7, 0x34, 0xa9, 0xc5, 0xd3, 0x2a, 0xe2, 0xb9, 0x80, 0xce,
	0x8b, 0x42, 0x59, 0x85, 0x5a, 0xa2, 0xe0, 0x15, 0x86, 0x8b, 0x72, 0xf4, 0x0a, 0xe5, 0x2e, 0x8d,
	0x90, 0x84, 0x60, 0x3f, 0xa1, 0xae, 0xa6, 0x5f, 0x9d, 0xd8, 0xed, 0x70, 0xfc, 0x78, 0x7c, 0xea,
	0xeb, 0x70, 0xe4, 0xe3, 0x2d, 0xdc, 0x70, 0xd4, 0xf5, 0xf6, 0x94, 0x7d, 0xca, 0x0b, 0x54, 0xf9,
	0xde, 0x3a, 0x45, 0x71, 0xee, 0x7e, 0x03, 0x00, 0x00, 0xff, 0xff, 0xc6, 0x86, 0x1a, 0x7f, 0x6e,
	0x02, 0x00, 0x00,
}
