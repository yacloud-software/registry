// Code generated by protoc-gen-go.
// source: golang.conradwood.net/apis/geoip/geoip.proto
// DO NOT EDIT!

/*
Package geoip is a generated protocol buffer package.

It is generated from these files:
	golang.conradwood.net/apis/geoip/geoip.proto

It has these top-level messages:
	LookupRequest
	LookupResponse
	IPRecord
	LookupBatchRequest
	LookupBatchResponse
*/
package geoip

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

type LookupRequest struct {
	IP string `protobuf:"bytes,1,opt,name=IP" json:"IP,omitempty"`
}

func (m *LookupRequest) Reset()                    { *m = LookupRequest{} }
func (m *LookupRequest) String() string            { return proto.CompactTextString(m) }
func (*LookupRequest) ProtoMessage()               {}
func (*LookupRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *LookupRequest) GetIP() string {
	if m != nil {
		return m.IP
	}
	return ""
}

type LookupResponse struct {
	AS           string  `protobuf:"bytes,1,opt,name=AS" json:"AS,omitempty"`
	City         string  `protobuf:"bytes,2,opt,name=City" json:"City,omitempty"`
	Country      string  `protobuf:"bytes,3,opt,name=Country" json:"Country,omitempty"`
	CountryCode  string  `protobuf:"bytes,4,opt,name=CountryCode" json:"CountryCode,omitempty"`
	Latitude     float64 `protobuf:"fixed64,5,opt,name=Latitude" json:"Latitude,omitempty"`
	Longitude    float64 `protobuf:"fixed64,6,opt,name=Longitude" json:"Longitude,omitempty"`
	Organisation string  `protobuf:"bytes,7,opt,name=Organisation" json:"Organisation,omitempty"`
	Region       string  `protobuf:"bytes,8,opt,name=Region" json:"Region,omitempty"`
	RegionName   string  `protobuf:"bytes,9,opt,name=RegionName" json:"RegionName,omitempty"`
	Timezone     string  `protobuf:"bytes,10,opt,name=Timezone" json:"Timezone,omitempty"`
	IP           string  `protobuf:"bytes,11,opt,name=IP" json:"IP,omitempty"`
}

func (m *LookupResponse) Reset()                    { *m = LookupResponse{} }
func (m *LookupResponse) String() string            { return proto.CompactTextString(m) }
func (*LookupResponse) ProtoMessage()               {}
func (*LookupResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *LookupResponse) GetAS() string {
	if m != nil {
		return m.AS
	}
	return ""
}

func (m *LookupResponse) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *LookupResponse) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *LookupResponse) GetCountryCode() string {
	if m != nil {
		return m.CountryCode
	}
	return ""
}

func (m *LookupResponse) GetLatitude() float64 {
	if m != nil {
		return m.Latitude
	}
	return 0
}

func (m *LookupResponse) GetLongitude() float64 {
	if m != nil {
		return m.Longitude
	}
	return 0
}

func (m *LookupResponse) GetOrganisation() string {
	if m != nil {
		return m.Organisation
	}
	return ""
}

func (m *LookupResponse) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *LookupResponse) GetRegionName() string {
	if m != nil {
		return m.RegionName
	}
	return ""
}

func (m *LookupResponse) GetTimezone() string {
	if m != nil {
		return m.Timezone
	}
	return ""
}

func (m *LookupResponse) GetIP() string {
	if m != nil {
		return m.IP
	}
	return ""
}

type IPRecord struct {
	ID           uint64 `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	ASRecord     string `protobuf:"bytes,2,opt,name=ASRecord" json:"ASRecord,omitempty"`
	City         string `protobuf:"bytes,3,opt,name=City" json:"City,omitempty"`
	Country      string `protobuf:"bytes,4,opt,name=Country" json:"Country,omitempty"`
	CountryCode  string `protobuf:"bytes,5,opt,name=CountryCode" json:"CountryCode,omitempty"`
	Organisation string `protobuf:"bytes,6,opt,name=Organisation" json:"Organisation,omitempty"`
	Region       string `protobuf:"bytes,7,opt,name=Region" json:"Region,omitempty"`
	RegionName   string `protobuf:"bytes,8,opt,name=RegionName" json:"RegionName,omitempty"`
	Timezone     string `protobuf:"bytes,9,opt,name=Timezone" json:"Timezone,omitempty"`
	IP           string `protobuf:"bytes,10,opt,name=IP" json:"IP,omitempty"`
}

func (m *IPRecord) Reset()                    { *m = IPRecord{} }
func (m *IPRecord) String() string            { return proto.CompactTextString(m) }
func (*IPRecord) ProtoMessage()               {}
func (*IPRecord) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *IPRecord) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *IPRecord) GetASRecord() string {
	if m != nil {
		return m.ASRecord
	}
	return ""
}

func (m *IPRecord) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *IPRecord) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *IPRecord) GetCountryCode() string {
	if m != nil {
		return m.CountryCode
	}
	return ""
}

func (m *IPRecord) GetOrganisation() string {
	if m != nil {
		return m.Organisation
	}
	return ""
}

func (m *IPRecord) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *IPRecord) GetRegionName() string {
	if m != nil {
		return m.RegionName
	}
	return ""
}

func (m *IPRecord) GetTimezone() string {
	if m != nil {
		return m.Timezone
	}
	return ""
}

func (m *IPRecord) GetIP() string {
	if m != nil {
		return m.IP
	}
	return ""
}

type LookupBatchRequest struct {
	IPs []string `protobuf:"bytes,1,rep,name=IPs" json:"IPs,omitempty"`
}

func (m *LookupBatchRequest) Reset()                    { *m = LookupBatchRequest{} }
func (m *LookupBatchRequest) String() string            { return proto.CompactTextString(m) }
func (*LookupBatchRequest) ProtoMessage()               {}
func (*LookupBatchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *LookupBatchRequest) GetIPs() []string {
	if m != nil {
		return m.IPs
	}
	return nil
}

type LookupBatchResponse struct {
	Responses []*LookupResponse `protobuf:"bytes,1,rep,name=Responses" json:"Responses,omitempty"`
}

func (m *LookupBatchResponse) Reset()                    { *m = LookupBatchResponse{} }
func (m *LookupBatchResponse) String() string            { return proto.CompactTextString(m) }
func (*LookupBatchResponse) ProtoMessage()               {}
func (*LookupBatchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *LookupBatchResponse) GetResponses() []*LookupResponse {
	if m != nil {
		return m.Responses
	}
	return nil
}

func init() {
	proto.RegisterType((*LookupRequest)(nil), "geoip.LookupRequest")
	proto.RegisterType((*LookupResponse)(nil), "geoip.LookupResponse")
	proto.RegisterType((*IPRecord)(nil), "geoip.IPRecord")
	proto.RegisterType((*LookupBatchRequest)(nil), "geoip.LookupBatchRequest")
	proto.RegisterType((*LookupBatchResponse)(nil), "geoip.LookupBatchResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for GeoIPService service

type GeoIPServiceClient interface {
	Lookup(ctx context.Context, in *LookupRequest, opts ...grpc.CallOption) (*LookupResponse, error)
	LookupBatch(ctx context.Context, in *LookupBatchRequest, opts ...grpc.CallOption) (*LookupBatchResponse, error)
}

type geoIPServiceClient struct {
	cc *grpc.ClientConn
}

func NewGeoIPServiceClient(cc *grpc.ClientConn) GeoIPServiceClient {
	return &geoIPServiceClient{cc}
}

func (c *geoIPServiceClient) Lookup(ctx context.Context, in *LookupRequest, opts ...grpc.CallOption) (*LookupResponse, error) {
	out := new(LookupResponse)
	err := grpc.Invoke(ctx, "/geoip.GeoIPService/Lookup", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *geoIPServiceClient) LookupBatch(ctx context.Context, in *LookupBatchRequest, opts ...grpc.CallOption) (*LookupBatchResponse, error) {
	out := new(LookupBatchResponse)
	err := grpc.Invoke(ctx, "/geoip.GeoIPService/LookupBatch", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GeoIPService service

type GeoIPServiceServer interface {
	Lookup(context.Context, *LookupRequest) (*LookupResponse, error)
	LookupBatch(context.Context, *LookupBatchRequest) (*LookupBatchResponse, error)
}

func RegisterGeoIPServiceServer(s *grpc.Server, srv GeoIPServiceServer) {
	s.RegisterService(&_GeoIPService_serviceDesc, srv)
}

func _GeoIPService_Lookup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LookupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeoIPServiceServer).Lookup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/geoip.GeoIPService/Lookup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeoIPServiceServer).Lookup(ctx, req.(*LookupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GeoIPService_LookupBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LookupBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeoIPServiceServer).LookupBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/geoip.GeoIPService/LookupBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeoIPServiceServer).LookupBatch(ctx, req.(*LookupBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _GeoIPService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "geoip.GeoIPService",
	HandlerType: (*GeoIPServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Lookup",
			Handler:    _GeoIPService_Lookup_Handler,
		},
		{
			MethodName: "LookupBatch",
			Handler:    _GeoIPService_LookupBatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "golang.conradwood.net/apis/geoip/geoip.proto",
}

func init() { proto.RegisterFile("golang.conradwood.net/apis/geoip/geoip.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 440 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x93, 0xdf, 0x8a, 0xd3, 0x50,
	0x10, 0xc6, 0x49, 0xda, 0xa6, 0xcd, 0x74, 0x5d, 0x64, 0xfc, 0xc3, 0xb1, 0x88, 0x86, 0x5c, 0x48,
	0x2f, 0xa4, 0x0b, 0xbb, 0xf8, 0x00, 0xdd, 0x16, 0x24, 0x52, 0x34, 0xa4, 0xbe, 0x40, 0x6c, 0x86,
	0x78, 0xd0, 0x3d, 0x13, 0x93, 0x53, 0x65, 0x7d, 0x05, 0xc1, 0x87, 0xf1, 0x09, 0x25, 0xe7, 0x24,
	0xd9, 0x74, 0xb7, 0xf4, 0x26, 0xcc, 0xcc, 0xf7, 0x9d, 0x61, 0xf2, 0x63, 0x06, 0xde, 0xe6, 0xfc,
	0x3d, 0x55, 0xf9, 0x62, 0xc7, 0xaa, 0x4c, 0xb3, 0x5f, 0xcc, 0xd9, 0x42, 0x91, 0xbe, 0x48, 0x0b,
	0x59, 0x5d, 0xe4, 0xc4, 0xb2, 0xb0, 0xdf, 0x45, 0x51, 0xb2, 0x66, 0x1c, 0x99, 0x24, 0x7c, 0x0d,
	0x8f, 0x36, 0xcc, 0xdf, 0xf6, 0x45, 0x42, 0x3f, 0xf6, 0x54, 0x69, 0x3c, 0x07, 0x37, 0x8a, 0x85,
	0x13, 0x38, 0x73, 0x3f, 0x71, 0xa3, 0x38, 0xfc, 0xe7, 0xc2, 0x79, 0xeb, 0xa8, 0x0a, 0x56, 0x15,
	0xd5, 0x96, 0xe5, 0xb6, 0xb5, 0x2c, 0xb7, 0x88, 0x30, 0x5c, 0x49, 0x7d, 0x2b, 0x5c, 0x53, 0x31,
	0x31, 0x0a, 0x18, 0xaf, 0x78, 0xaf, 0x74, 0x79, 0x2b, 0x06, 0xa6, 0xdc, 0xa6, 0x18, 0xc0, 0xb4,
	0x09, 0x57, 0x9c, 0x91, 0x18, 0x1a, 0xb5, 0x5f, 0xc2, 0x19, 0x4c, 0x36, 0xa9, 0x96, 0x7a, 0x9f,
	0x91, 0x18, 0x05, 0xce, 0xdc, 0x49, 0xba, 0x1c, 0x5f, 0x82, 0xbf, 0x61, 0x95, 0x5b, 0xd1, 0x33,
	0xe2, 0x5d, 0x01, 0x43, 0x38, 0xfb, 0x54, 0xe6, 0xa9, 0x92, 0x55, 0xaa, 0x25, 0x2b, 0x31, 0x36,
	0xcd, 0x0f, 0x6a, 0xf8, 0x1c, 0xbc, 0x84, 0xf2, 0x5a, 0x9d, 0x18, 0xb5, 0xc9, 0xf0, 0x15, 0x80,
	0x8d, 0x3e, 0xa6, 0x37, 0x24, 0x7c, 0xa3, 0xf5, 0x2a, 0xf5, 0x54, 0x9f, 0xe5, 0x0d, 0xfd, 0x66,
	0x45, 0x02, 0x8c, 0xda, 0xe5, 0x0d, 0xb4, 0x69, 0x07, 0xed, 0xaf, 0x0b, 0x93, 0x28, 0x4e, 0x68,
	0xc7, 0x65, 0x66, 0xc4, 0xb5, 0xc1, 0x35, 0x4c, 0xdc, 0x68, 0x5d, 0x37, 0x5a, 0x6e, 0xad, 0xd6,
	0x20, 0xeb, 0xf2, 0x0e, 0xe5, 0xe0, 0x38, 0xca, 0xe1, 0x49, 0x94, 0xa3, 0x87, 0x28, 0xef, 0x03,
	0xf1, 0x4e, 0x02, 0x19, 0x9f, 0x00, 0x32, 0x39, 0x09, 0xc4, 0x3f, 0x0a, 0x04, 0x3a, 0x20, 0x6f,
	0x00, 0xed, 0x12, 0x5d, 0xa7, 0x7a, 0xf7, 0xb5, 0xdd, 0xb5, 0xc7, 0x30, 0x88, 0xe2, 0x4a, 0x38,
	0xc1, 0x60, 0xee, 0x27, 0x75, 0x18, 0x7e, 0x80, 0x27, 0x07, 0xbe, 0x66, 0xe3, 0xae, 0xc0, 0x6f,
	0x63, 0x6b, 0x9f, 0x5e, 0x3e, 0x5b, 0xd8, 0x6d, 0x3e, 0xdc, 0xcd, 0xe4, 0xce, 0x77, 0xf9, 0xc7,
	0x81, 0xb3, 0xf7, 0xc4, 0x51, 0xbc, 0xa5, 0xf2, 0xa7, 0xdc, 0x11, 0xbe, 0x03, 0xcf, 0xba, 0xf1,
	0xe9, 0xbd, 0xc7, 0x66, 0x9c, 0xd9, 0xf1, 0x96, 0xb8, 0x86, 0x69, 0x6f, 0x26, 0x7c, 0x71, 0xe0,
	0xea, 0xff, 0xcf, 0x6c, 0x76, 0x4c, 0xb2, 0x5d, 0xae, 0x43, 0x08, 0x14, 0xe9, 0xfe, 0x71, 0x36,
	0xe7, 0x5a, 0xdf, 0xa7, 0x7d, 0xf7, 0xc5, 0x33, 0xa7, 0x79, 0xf5, 0x3f, 0x00, 0x00, 0xff, 0xff,
	0xd6, 0x2d, 0x91, 0x9c, 0xca, 0x03, 0x00, 0x00,
}
