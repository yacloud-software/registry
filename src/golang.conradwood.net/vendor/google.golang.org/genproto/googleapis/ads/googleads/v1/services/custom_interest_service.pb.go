// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v1/services/custom_interest_service.proto

package services

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/wrappers"
	resources "google.golang.org/genproto/googleapis/ads/googleads/v1/resources"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// Request message for
// [CustomInterestService.GetCustomInterest][google.ads.googleads.v1.services.CustomInterestService.GetCustomInterest].
type GetCustomInterestRequest struct {
	// The resource name of the custom interest to fetch.
	ResourceName         string   `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetCustomInterestRequest) Reset()         { *m = GetCustomInterestRequest{} }
func (m *GetCustomInterestRequest) String() string { return proto.CompactTextString(m) }
func (*GetCustomInterestRequest) ProtoMessage()    {}
func (*GetCustomInterestRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddfad66fd19c3e99, []int{0}
}

func (m *GetCustomInterestRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetCustomInterestRequest.Unmarshal(m, b)
}
func (m *GetCustomInterestRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetCustomInterestRequest.Marshal(b, m, deterministic)
}
func (m *GetCustomInterestRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetCustomInterestRequest.Merge(m, src)
}
func (m *GetCustomInterestRequest) XXX_Size() int {
	return xxx_messageInfo_GetCustomInterestRequest.Size(m)
}
func (m *GetCustomInterestRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetCustomInterestRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetCustomInterestRequest proto.InternalMessageInfo

func (m *GetCustomInterestRequest) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

// Request message for
// [CustomInterestService.MutateCustomInterests][google.ads.googleads.v1.services.CustomInterestService.MutateCustomInterests].
type MutateCustomInterestsRequest struct {
	// The ID of the customer whose custom interests are being modified.
	CustomerId string `protobuf:"bytes,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	// The list of operations to perform on individual custom interests.
	Operations []*CustomInterestOperation `protobuf:"bytes,2,rep,name=operations,proto3" json:"operations,omitempty"`
	// If true, the request is validated but not executed. Only errors are
	// returned, not results.
	ValidateOnly         bool     `protobuf:"varint,4,opt,name=validate_only,json=validateOnly,proto3" json:"validate_only,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MutateCustomInterestsRequest) Reset()         { *m = MutateCustomInterestsRequest{} }
func (m *MutateCustomInterestsRequest) String() string { return proto.CompactTextString(m) }
func (*MutateCustomInterestsRequest) ProtoMessage()    {}
func (*MutateCustomInterestsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddfad66fd19c3e99, []int{1}
}

func (m *MutateCustomInterestsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MutateCustomInterestsRequest.Unmarshal(m, b)
}
func (m *MutateCustomInterestsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MutateCustomInterestsRequest.Marshal(b, m, deterministic)
}
func (m *MutateCustomInterestsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MutateCustomInterestsRequest.Merge(m, src)
}
func (m *MutateCustomInterestsRequest) XXX_Size() int {
	return xxx_messageInfo_MutateCustomInterestsRequest.Size(m)
}
func (m *MutateCustomInterestsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MutateCustomInterestsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MutateCustomInterestsRequest proto.InternalMessageInfo

func (m *MutateCustomInterestsRequest) GetCustomerId() string {
	if m != nil {
		return m.CustomerId
	}
	return ""
}

func (m *MutateCustomInterestsRequest) GetOperations() []*CustomInterestOperation {
	if m != nil {
		return m.Operations
	}
	return nil
}

func (m *MutateCustomInterestsRequest) GetValidateOnly() bool {
	if m != nil {
		return m.ValidateOnly
	}
	return false
}

// A single operation (create, update) on a custom interest.
type CustomInterestOperation struct {
	// FieldMask that determines which resource fields are modified in an update.
	UpdateMask *field_mask.FieldMask `protobuf:"bytes,4,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	// The mutate operation.
	//
	// Types that are valid to be assigned to Operation:
	//	*CustomInterestOperation_Create
	//	*CustomInterestOperation_Update
	Operation            isCustomInterestOperation_Operation `protobuf_oneof:"operation"`
	XXX_NoUnkeyedLiteral struct{}                            `json:"-"`
	XXX_unrecognized     []byte                              `json:"-"`
	XXX_sizecache        int32                               `json:"-"`
}

func (m *CustomInterestOperation) Reset()         { *m = CustomInterestOperation{} }
func (m *CustomInterestOperation) String() string { return proto.CompactTextString(m) }
func (*CustomInterestOperation) ProtoMessage()    {}
func (*CustomInterestOperation) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddfad66fd19c3e99, []int{2}
}

func (m *CustomInterestOperation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CustomInterestOperation.Unmarshal(m, b)
}
func (m *CustomInterestOperation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CustomInterestOperation.Marshal(b, m, deterministic)
}
func (m *CustomInterestOperation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CustomInterestOperation.Merge(m, src)
}
func (m *CustomInterestOperation) XXX_Size() int {
	return xxx_messageInfo_CustomInterestOperation.Size(m)
}
func (m *CustomInterestOperation) XXX_DiscardUnknown() {
	xxx_messageInfo_CustomInterestOperation.DiscardUnknown(m)
}

var xxx_messageInfo_CustomInterestOperation proto.InternalMessageInfo

func (m *CustomInterestOperation) GetUpdateMask() *field_mask.FieldMask {
	if m != nil {
		return m.UpdateMask
	}
	return nil
}

type isCustomInterestOperation_Operation interface {
	isCustomInterestOperation_Operation()
}

type CustomInterestOperation_Create struct {
	Create *resources.CustomInterest `protobuf:"bytes,1,opt,name=create,proto3,oneof"`
}

type CustomInterestOperation_Update struct {
	Update *resources.CustomInterest `protobuf:"bytes,2,opt,name=update,proto3,oneof"`
}

func (*CustomInterestOperation_Create) isCustomInterestOperation_Operation() {}

func (*CustomInterestOperation_Update) isCustomInterestOperation_Operation() {}

func (m *CustomInterestOperation) GetOperation() isCustomInterestOperation_Operation {
	if m != nil {
		return m.Operation
	}
	return nil
}

func (m *CustomInterestOperation) GetCreate() *resources.CustomInterest {
	if x, ok := m.GetOperation().(*CustomInterestOperation_Create); ok {
		return x.Create
	}
	return nil
}

func (m *CustomInterestOperation) GetUpdate() *resources.CustomInterest {
	if x, ok := m.GetOperation().(*CustomInterestOperation_Update); ok {
		return x.Update
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*CustomInterestOperation) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*CustomInterestOperation_Create)(nil),
		(*CustomInterestOperation_Update)(nil),
	}
}

// Response message for custom interest mutate.
type MutateCustomInterestsResponse struct {
	// All results for the mutate.
	Results              []*MutateCustomInterestResult `protobuf:"bytes,2,rep,name=results,proto3" json:"results,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *MutateCustomInterestsResponse) Reset()         { *m = MutateCustomInterestsResponse{} }
func (m *MutateCustomInterestsResponse) String() string { return proto.CompactTextString(m) }
func (*MutateCustomInterestsResponse) ProtoMessage()    {}
func (*MutateCustomInterestsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddfad66fd19c3e99, []int{3}
}

func (m *MutateCustomInterestsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MutateCustomInterestsResponse.Unmarshal(m, b)
}
func (m *MutateCustomInterestsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MutateCustomInterestsResponse.Marshal(b, m, deterministic)
}
func (m *MutateCustomInterestsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MutateCustomInterestsResponse.Merge(m, src)
}
func (m *MutateCustomInterestsResponse) XXX_Size() int {
	return xxx_messageInfo_MutateCustomInterestsResponse.Size(m)
}
func (m *MutateCustomInterestsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MutateCustomInterestsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MutateCustomInterestsResponse proto.InternalMessageInfo

func (m *MutateCustomInterestsResponse) GetResults() []*MutateCustomInterestResult {
	if m != nil {
		return m.Results
	}
	return nil
}

// The result for the custom interest mutate.
type MutateCustomInterestResult struct {
	// Returned for successful operations.
	ResourceName         string   `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MutateCustomInterestResult) Reset()         { *m = MutateCustomInterestResult{} }
func (m *MutateCustomInterestResult) String() string { return proto.CompactTextString(m) }
func (*MutateCustomInterestResult) ProtoMessage()    {}
func (*MutateCustomInterestResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_ddfad66fd19c3e99, []int{4}
}

func (m *MutateCustomInterestResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MutateCustomInterestResult.Unmarshal(m, b)
}
func (m *MutateCustomInterestResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MutateCustomInterestResult.Marshal(b, m, deterministic)
}
func (m *MutateCustomInterestResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MutateCustomInterestResult.Merge(m, src)
}
func (m *MutateCustomInterestResult) XXX_Size() int {
	return xxx_messageInfo_MutateCustomInterestResult.Size(m)
}
func (m *MutateCustomInterestResult) XXX_DiscardUnknown() {
	xxx_messageInfo_MutateCustomInterestResult.DiscardUnknown(m)
}

var xxx_messageInfo_MutateCustomInterestResult proto.InternalMessageInfo

func (m *MutateCustomInterestResult) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func init() {
	proto.RegisterType((*GetCustomInterestRequest)(nil), "google.ads.googleads.v1.services.GetCustomInterestRequest")
	proto.RegisterType((*MutateCustomInterestsRequest)(nil), "google.ads.googleads.v1.services.MutateCustomInterestsRequest")
	proto.RegisterType((*CustomInterestOperation)(nil), "google.ads.googleads.v1.services.CustomInterestOperation")
	proto.RegisterType((*MutateCustomInterestsResponse)(nil), "google.ads.googleads.v1.services.MutateCustomInterestsResponse")
	proto.RegisterType((*MutateCustomInterestResult)(nil), "google.ads.googleads.v1.services.MutateCustomInterestResult")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v1/services/custom_interest_service.proto", fileDescriptor_ddfad66fd19c3e99)
}

var fileDescriptor_ddfad66fd19c3e99 = []byte{
	// 633 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0x41, 0x6b, 0xd4, 0x40,
	0x14, 0xc7, 0x4d, 0x2a, 0xd5, 0x4e, 0xea, 0xc1, 0x40, 0x71, 0x09, 0x55, 0x97, 0xd8, 0xc3, 0xb2,
	0x87, 0x84, 0x6c, 0x05, 0xd9, 0xd4, 0x76, 0xc9, 0x0a, 0x6e, 0x8b, 0xd4, 0x96, 0x08, 0x0b, 0xca,
	0x42, 0x98, 0x6e, 0xa6, 0x21, 0x34, 0xc9, 0xc4, 0x99, 0xc9, 0x96, 0x52, 0x7a, 0xd0, 0xaf, 0xe0,
	0x37, 0xf0, 0xe8, 0xd5, 0xb3, 0x78, 0xf7, 0xea, 0xc9, 0xbb, 0x27, 0xbf, 0x83, 0x20, 0xc9, 0x64,
	0x62, 0x77, 0xdd, 0xb0, 0x5a, 0x6f, 0x2f, 0x93, 0xf7, 0x7e, 0xef, 0xfd, 0xdf, 0x9b, 0x37, 0x60,
	0x27, 0xc0, 0x38, 0x88, 0x90, 0x09, 0x7d, 0x6a, 0x72, 0x33, 0xb7, 0x26, 0x96, 0x49, 0x11, 0x99,
	0x84, 0x63, 0x44, 0xcd, 0x71, 0x46, 0x19, 0x8e, 0xbd, 0x30, 0x61, 0x88, 0x20, 0xca, 0xbc, 0xf2,
	0x87, 0x91, 0x12, 0xcc, 0xb0, 0xda, 0xe4, 0x41, 0x06, 0xf4, 0xa9, 0x51, 0xc5, 0x1b, 0x13, 0xcb,
	0x10, 0xf1, 0xda, 0xa3, 0xba, 0x0c, 0x04, 0x51, 0x9c, 0x91, 0x39, 0x29, 0x38, 0x5a, 0x5b, 0x17,
	0x81, 0x69, 0x68, 0xc2, 0x24, 0xc1, 0x0c, 0xb2, 0x10, 0x27, 0xb4, 0xfc, 0x5b, 0x26, 0x36, 0x8b,
	0xaf, 0xa3, 0xec, 0xd8, 0x3c, 0x0e, 0x51, 0xe4, 0x7b, 0x31, 0xa4, 0x27, 0xa5, 0xc7, 0xbd, 0x59,
	0x8f, 0x53, 0x02, 0xd3, 0x14, 0x91, 0x92, 0xa0, 0xf7, 0x40, 0x63, 0x80, 0xd8, 0x93, 0x22, 0xf7,
	0x5e, 0x99, 0xda, 0x45, 0xaf, 0x33, 0x44, 0x99, 0xfa, 0x00, 0xdc, 0x12, 0xe5, 0x79, 0x09, 0x8c,
	0x51, 0x43, 0x6a, 0x4a, 0xad, 0x15, 0x77, 0x55, 0x1c, 0x3e, 0x87, 0x31, 0xd2, 0x3f, 0x4b, 0x60,
	0x7d, 0x3f, 0x63, 0x90, 0xa1, 0x69, 0x08, 0x15, 0x94, 0xfb, 0x40, 0xe1, 0xd2, 0x10, 0xf1, 0x42,
	0xbf, 0x64, 0x00, 0x71, 0xb4, 0xe7, 0xab, 0x2f, 0x01, 0xc0, 0x29, 0x22, 0x5c, 0x58, 0x43, 0x6e,
	0x2e, 0xb5, 0x94, 0x4e, 0xd7, 0x58, 0xd4, 0x52, 0x63, 0x3a, 0xdd, 0x81, 0x20, 0xb8, 0x97, 0x60,
	0xb9, 0x82, 0x09, 0x8c, 0x42, 0x1f, 0x32, 0xe4, 0xe1, 0x24, 0x3a, 0x6b, 0x5c, 0x6f, 0x4a, 0xad,
	0x9b, 0xee, 0xaa, 0x38, 0x3c, 0x48, 0xa2, 0x33, 0xfd, 0xa7, 0x04, 0xee, 0xd4, 0xc0, 0xd4, 0x2d,
	0xa0, 0x64, 0x69, 0x11, 0x9e, 0xf7, 0xb4, 0x08, 0x57, 0x3a, 0x9a, 0x28, 0x4e, 0x34, 0xd5, 0x78,
	0x9a, 0xb7, 0x7d, 0x1f, 0xd2, 0x13, 0x17, 0x70, 0xf7, 0xdc, 0x56, 0x9f, 0x81, 0xe5, 0x31, 0x41,
	0x90, 0xf1, 0xc6, 0x29, 0x1d, 0xab, 0x56, 0x54, 0x75, 0x0b, 0x66, 0x54, 0xed, 0x5e, 0x73, 0x4b,
	0x44, 0x0e, 0xe3, 0xe8, 0x86, 0xfc, 0x1f, 0x30, 0x8e, 0xe8, 0x2b, 0x60, 0xa5, 0xea, 0x92, 0x7e,
	0x0a, 0xee, 0xd6, 0x0c, 0x90, 0xa6, 0x38, 0xa1, 0x48, 0x1d, 0x82, 0x1b, 0x04, 0xd1, 0x2c, 0x62,
	0x62, 0x3a, 0x8f, 0x17, 0x4f, 0x67, 0x1e, 0xd1, 0x2d, 0x20, 0xae, 0x80, 0xe9, 0x0e, 0xd0, 0xea,
	0xdd, 0xfe, 0xea, 0xf6, 0x75, 0x3e, 0x2e, 0x81, 0xb5, 0xe9, 0xe8, 0x17, 0xbc, 0x02, 0xf5, 0x93,
	0x04, 0x6e, 0xff, 0x71, 0xb3, 0x55, 0x7b, 0x71, 0xe5, 0x75, 0xeb, 0xa0, 0xfd, 0x7b, 0xc7, 0xf5,
	0xee, 0xdb, 0xaf, 0xdf, 0xdf, 0xc9, 0x9b, 0xaa, 0x95, 0xaf, 0xfa, 0xf9, 0x94, 0x9c, 0x6d, 0xb1,
	0x01, 0xd4, 0x6c, 0x97, 0xbb, 0x5f, 0x75, 0xde, 0x6c, 0x5f, 0xa8, 0xdf, 0x24, 0xb0, 0x36, 0x77,
	0x2c, 0xea, 0xce, 0xd5, 0xba, 0x2f, 0x16, 0x52, 0xeb, 0x5d, 0x39, 0x9e, 0xdf, 0x07, 0xbd, 0x57,
	0xa8, 0xea, 0xea, 0x0f, 0x73, 0x55, 0xbf, 0x65, 0x9c, 0x5f, 0x5a, 0xf3, 0xed, 0xf6, 0xc5, 0xac,
	0x28, 0x3b, 0x2e, 0xa0, 0xb6, 0xd4, 0xee, 0xbf, 0x91, 0xc1, 0xc6, 0x18, 0xc7, 0x0b, 0xeb, 0xe8,
	0x6b, 0x73, 0x67, 0x7b, 0x98, 0xaf, 0xdd, 0xa1, 0xf4, 0x6a, 0xb7, 0x8c, 0x0f, 0x70, 0x04, 0x93,
	0xc0, 0xc0, 0x24, 0x30, 0x03, 0x94, 0x14, 0x4b, 0x29, 0x1e, 0xd9, 0x34, 0xa4, 0xf5, 0xaf, 0xfa,
	0x96, 0x30, 0xde, 0xcb, 0x4b, 0x03, 0xc7, 0xf9, 0x20, 0x37, 0x07, 0x1c, 0xe8, 0xf8, 0xd4, 0xe0,
	0x66, 0x6e, 0x0d, 0x2d, 0xa3, 0x4c, 0x4c, 0xbf, 0x08, 0x97, 0x91, 0xe3, 0xd3, 0x51, 0xe5, 0x32,
	0x1a, 0x5a, 0x23, 0xe1, 0xf2, 0x43, 0xde, 0xe0, 0xe7, 0xb6, 0xed, 0xf8, 0xd4, 0xb6, 0x2b, 0x27,
	0xdb, 0x1e, 0x5a, 0xb6, 0x2d, 0xdc, 0x8e, 0x96, 0x8b, 0x3a, 0x37, 0x7f, 0x05, 0x00, 0x00, 0xff,
	0xff, 0x62, 0xea, 0xe7, 0xce, 0x7c, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CustomInterestServiceClient is the client API for CustomInterestService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CustomInterestServiceClient interface {
	// Returns the requested custom interest in full detail.
	GetCustomInterest(ctx context.Context, in *GetCustomInterestRequest, opts ...grpc.CallOption) (*resources.CustomInterest, error)
	// Creates or updates custom interests. Operation statuses are returned.
	MutateCustomInterests(ctx context.Context, in *MutateCustomInterestsRequest, opts ...grpc.CallOption) (*MutateCustomInterestsResponse, error)
}

type customInterestServiceClient struct {
	cc *grpc.ClientConn
}

func NewCustomInterestServiceClient(cc *grpc.ClientConn) CustomInterestServiceClient {
	return &customInterestServiceClient{cc}
}

func (c *customInterestServiceClient) GetCustomInterest(ctx context.Context, in *GetCustomInterestRequest, opts ...grpc.CallOption) (*resources.CustomInterest, error) {
	out := new(resources.CustomInterest)
	err := c.cc.Invoke(ctx, "/google.ads.googleads.v1.services.CustomInterestService/GetCustomInterest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customInterestServiceClient) MutateCustomInterests(ctx context.Context, in *MutateCustomInterestsRequest, opts ...grpc.CallOption) (*MutateCustomInterestsResponse, error) {
	out := new(MutateCustomInterestsResponse)
	err := c.cc.Invoke(ctx, "/google.ads.googleads.v1.services.CustomInterestService/MutateCustomInterests", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CustomInterestServiceServer is the server API for CustomInterestService service.
type CustomInterestServiceServer interface {
	// Returns the requested custom interest in full detail.
	GetCustomInterest(context.Context, *GetCustomInterestRequest) (*resources.CustomInterest, error)
	// Creates or updates custom interests. Operation statuses are returned.
	MutateCustomInterests(context.Context, *MutateCustomInterestsRequest) (*MutateCustomInterestsResponse, error)
}

// UnimplementedCustomInterestServiceServer can be embedded to have forward compatible implementations.
type UnimplementedCustomInterestServiceServer struct {
}

func (*UnimplementedCustomInterestServiceServer) GetCustomInterest(ctx context.Context, req *GetCustomInterestRequest) (*resources.CustomInterest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCustomInterest not implemented")
}
func (*UnimplementedCustomInterestServiceServer) MutateCustomInterests(ctx context.Context, req *MutateCustomInterestsRequest) (*MutateCustomInterestsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MutateCustomInterests not implemented")
}

func RegisterCustomInterestServiceServer(s *grpc.Server, srv CustomInterestServiceServer) {
	s.RegisterService(&_CustomInterestService_serviceDesc, srv)
}

func _CustomInterestService_GetCustomInterest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCustomInterestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomInterestServiceServer).GetCustomInterest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ads.googleads.v1.services.CustomInterestService/GetCustomInterest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomInterestServiceServer).GetCustomInterest(ctx, req.(*GetCustomInterestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CustomInterestService_MutateCustomInterests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutateCustomInterestsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomInterestServiceServer).MutateCustomInterests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ads.googleads.v1.services.CustomInterestService/MutateCustomInterests",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomInterestServiceServer).MutateCustomInterests(ctx, req.(*MutateCustomInterestsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CustomInterestService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.ads.googleads.v1.services.CustomInterestService",
	HandlerType: (*CustomInterestServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCustomInterest",
			Handler:    _CustomInterestService_GetCustomInterest_Handler,
		},
		{
			MethodName: "MutateCustomInterests",
			Handler:    _CustomInterestService_MutateCustomInterests_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/ads/googleads/v1/services/custom_interest_service.proto",
}
