// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.18.0
// source: coupons.proto

package coupons_pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SendCouponsToUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mobile    string `protobuf:"bytes,1,opt,name=mobile,proto3" json:"mobile,omitempty"`
	CouponsId uint32 `protobuf:"varint,2,opt,name=coupons_id,json=couponsId,proto3" json:"coupons_id,omitempty"`
	Num       uint32 `protobuf:"varint,3,opt,name=num,proto3" json:"num,omitempty"`
}

func (x *SendCouponsToUserRequest) Reset() {
	*x = SendCouponsToUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coupons_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendCouponsToUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendCouponsToUserRequest) ProtoMessage() {}

func (x *SendCouponsToUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_coupons_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendCouponsToUserRequest.ProtoReflect.Descriptor instead.
func (*SendCouponsToUserRequest) Descriptor() ([]byte, []int) {
	return file_coupons_proto_rawDescGZIP(), []int{0}
}

func (x *SendCouponsToUserRequest) GetMobile() string {
	if x != nil {
		return x.Mobile
	}
	return ""
}

func (x *SendCouponsToUserRequest) GetCouponsId() uint32 {
	if x != nil {
		return x.CouponsId
	}
	return 0
}

func (x *SendCouponsToUserRequest) GetNum() uint32 {
	if x != nil {
		return x.Num
	}
	return 0
}

var File_coupons_proto protoreflect.FileDescriptor

var file_coupons_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x63, 0x6f, 0x75, 0x70, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x63, 0x0a, 0x18,
	0x53, 0x65, 0x6e, 0x64, 0x43, 0x6f, 0x75, 0x70, 0x6f, 0x6e, 0x73, 0x54, 0x6f, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x6f, 0x62, 0x69,
	0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65,
	0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6f, 0x75, 0x70, 0x6f, 0x6e, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x63, 0x6f, 0x75, 0x70, 0x6f, 0x6e, 0x73, 0x49, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x6e, 0x75, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x6e, 0x75,
	0x6d, 0x32, 0x51, 0x0a, 0x07, 0x43, 0x6f, 0x75, 0x70, 0x6f, 0x6e, 0x73, 0x12, 0x46, 0x0a, 0x11,
	0x53, 0x65, 0x6e, 0x64, 0x43, 0x6f, 0x75, 0x70, 0x6f, 0x6e, 0x73, 0x54, 0x6f, 0x55, 0x73, 0x65,
	0x72, 0x12, 0x19, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x6f, 0x75, 0x70, 0x6f, 0x6e, 0x73, 0x54,
	0x6f, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x42, 0x1b, 0x5a, 0x19, 0x67, 0x65, 0x6e, 0x2f, 0x63, 0x6f, 0x75, 0x70,
	0x6f, 0x6e, 0x73, 0x5f, 0x70, 0x62, 0x3b, 0x63, 0x6f, 0x75, 0x70, 0x6f, 0x6e, 0x73, 0x5f, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_coupons_proto_rawDescOnce sync.Once
	file_coupons_proto_rawDescData = file_coupons_proto_rawDesc
)

func file_coupons_proto_rawDescGZIP() []byte {
	file_coupons_proto_rawDescOnce.Do(func() {
		file_coupons_proto_rawDescData = protoimpl.X.CompressGZIP(file_coupons_proto_rawDescData)
	})
	return file_coupons_proto_rawDescData
}

var file_coupons_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_coupons_proto_goTypes = []interface{}{
	(*SendCouponsToUserRequest)(nil), // 0: SendCouponsToUserRequest
	(*emptypb.Empty)(nil),            // 1: google.protobuf.Empty
}
var file_coupons_proto_depIdxs = []int32{
	0, // 0: Coupons.SendCouponsToUser:input_type -> SendCouponsToUserRequest
	1, // 1: Coupons.SendCouponsToUser:output_type -> google.protobuf.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_coupons_proto_init() }
func file_coupons_proto_init() {
	if File_coupons_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_coupons_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendCouponsToUserRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_coupons_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_coupons_proto_goTypes,
		DependencyIndexes: file_coupons_proto_depIdxs,
		MessageInfos:      file_coupons_proto_msgTypes,
	}.Build()
	File_coupons_proto = out.File
	file_coupons_proto_rawDesc = nil
	file_coupons_proto_goTypes = nil
	file_coupons_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CouponsClient is the client API for Coupons service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CouponsClient interface {
	SendCouponsToUser(ctx context.Context, in *SendCouponsToUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type couponsClient struct {
	cc grpc.ClientConnInterface
}

func NewCouponsClient(cc grpc.ClientConnInterface) CouponsClient {
	return &couponsClient{cc}
}

func (c *couponsClient) SendCouponsToUser(ctx context.Context, in *SendCouponsToUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/Coupons/SendCouponsToUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CouponsServer is the server API for Coupons service.
type CouponsServer interface {
	SendCouponsToUser(context.Context, *SendCouponsToUserRequest) (*emptypb.Empty, error)
}

// UnimplementedCouponsServer can be embedded to have forward compatible implementations.
type UnimplementedCouponsServer struct {
}

func (*UnimplementedCouponsServer) SendCouponsToUser(context.Context, *SendCouponsToUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCouponsToUser not implemented")
}

func RegisterCouponsServer(s *grpc.Server, srv CouponsServer) {
	s.RegisterService(&_Coupons_serviceDesc, srv)
}

func _Coupons_SendCouponsToUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendCouponsToUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CouponsServer).SendCouponsToUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Coupons/SendCouponsToUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CouponsServer).SendCouponsToUser(ctx, req.(*SendCouponsToUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Coupons_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Coupons",
	HandlerType: (*CouponsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendCouponsToUser",
			Handler:    _Coupons_SendCouponsToUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "coupons.proto",
}
