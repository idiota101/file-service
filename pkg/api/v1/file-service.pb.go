// Code generated by protoc-gen-go. DO NOT EDIT.
// source: file-service.proto

package v1

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

type UserDetails struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Email                string   `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Password             string   `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserDetails) Reset()         { *m = UserDetails{} }
func (m *UserDetails) String() string { return proto.CompactTextString(m) }
func (*UserDetails) ProtoMessage()    {}
func (*UserDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_b50cba0e4ffba287, []int{0}
}

func (m *UserDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserDetails.Unmarshal(m, b)
}
func (m *UserDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserDetails.Marshal(b, m, deterministic)
}
func (m *UserDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserDetails.Merge(m, src)
}
func (m *UserDetails) XXX_Size() int {
	return xxx_messageInfo_UserDetails.Size(m)
}
func (m *UserDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_UserDetails.DiscardUnknown(m)
}

var xxx_messageInfo_UserDetails proto.InternalMessageInfo

func (m *UserDetails) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *UserDetails) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UserDetails) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type RegisterUserRequest struct {
	Api                  string       `protobuf:"bytes,1,opt,name=api,proto3" json:"api,omitempty"`
	UserDetails          *UserDetails `protobuf:"bytes,2,opt,name=userDetails,proto3" json:"userDetails,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *RegisterUserRequest) Reset()         { *m = RegisterUserRequest{} }
func (m *RegisterUserRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterUserRequest) ProtoMessage()    {}
func (*RegisterUserRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b50cba0e4ffba287, []int{1}
}

func (m *RegisterUserRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterUserRequest.Unmarshal(m, b)
}
func (m *RegisterUserRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterUserRequest.Marshal(b, m, deterministic)
}
func (m *RegisterUserRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterUserRequest.Merge(m, src)
}
func (m *RegisterUserRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterUserRequest.Size(m)
}
func (m *RegisterUserRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterUserRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterUserRequest proto.InternalMessageInfo

func (m *RegisterUserRequest) GetApi() string {
	if m != nil {
		return m.Api
	}
	return ""
}

func (m *RegisterUserRequest) GetUserDetails() *UserDetails {
	if m != nil {
		return m.UserDetails
	}
	return nil
}

type RegisterUserResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterUserResponse) Reset()         { *m = RegisterUserResponse{} }
func (m *RegisterUserResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterUserResponse) ProtoMessage()    {}
func (*RegisterUserResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b50cba0e4ffba287, []int{2}
}

func (m *RegisterUserResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterUserResponse.Unmarshal(m, b)
}
func (m *RegisterUserResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterUserResponse.Marshal(b, m, deterministic)
}
func (m *RegisterUserResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterUserResponse.Merge(m, src)
}
func (m *RegisterUserResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterUserResponse.Size(m)
}
func (m *RegisterUserResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterUserResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterUserResponse proto.InternalMessageInfo

func (m *RegisterUserResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*UserDetails)(nil), "v1.UserDetails")
	proto.RegisterType((*RegisterUserRequest)(nil), "v1.RegisterUserRequest")
	proto.RegisterType((*RegisterUserResponse)(nil), "v1.RegisterUserResponse")
}

func init() { proto.RegisterFile("file-service.proto", fileDescriptor_b50cba0e4ffba287) }

var fileDescriptor_b50cba0e4ffba287 = []byte{
	// 223 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xcf, 0x4e, 0xc4, 0x20,
	0x10, 0xc6, 0xdd, 0x5d, 0xff, 0x0e, 0x26, 0x9a, 0x71, 0x13, 0xc9, 0x9e, 0x0c, 0x27, 0x2f, 0x36,
	0x76, 0x7d, 0x04, 0x8d, 0x0f, 0x40, 0xe3, 0xc5, 0x1b, 0xea, 0xd8, 0x90, 0xd0, 0x82, 0x0c, 0xad,
	0xaf, 0x6f, 0x4a, 0xad, 0xd6, 0xb8, 0xb7, 0xf9, 0x18, 0xf2, 0xfb, 0x7d, 0x00, 0xf8, 0x6e, 0x1d,
	0xdd, 0x30, 0xc5, 0xde, 0xbe, 0x52, 0x11, 0xa2, 0x4f, 0x1e, 0x97, 0x7d, 0xa9, 0x2a, 0x10, 0x4f,
	0x4c, 0xf1, 0x81, 0x92, 0xb1, 0x8e, 0x11, 0x61, 0xbf, 0x35, 0x0d, 0xc9, 0xc5, 0xd5, 0xe2, 0xfa,
	0x44, 0xe7, 0x19, 0xd7, 0x70, 0x40, 0x8d, 0xb1, 0x4e, 0x2e, 0xf3, 0xe1, 0x18, 0x70, 0x03, 0xc7,
	0xc1, 0x30, 0x7f, 0xfa, 0xf8, 0x26, 0x57, 0x79, 0xf1, 0x93, 0xd5, 0x33, 0x5c, 0x68, 0xaa, 0x2d,
	0x27, 0x8a, 0x03, 0x5c, 0xd3, 0x47, 0x47, 0x9c, 0xf0, 0x1c, 0x56, 0x26, 0xd8, 0x6f, 0xf6, 0x30,
	0x62, 0x09, 0xa2, 0xfb, 0xb5, 0x67, 0x81, 0xd8, 0x9e, 0x15, 0x7d, 0x59, 0xcc, 0x4a, 0xe9, 0xf9,
	0x1d, 0x75, 0x0b, 0xeb, 0xbf, 0x6c, 0x0e, 0xbe, 0x65, 0x42, 0x09, 0x47, 0x0d, 0x31, 0x9b, 0x7a,
	0x2a, 0x3f, 0xc5, 0xad, 0x06, 0xf1, 0x68, 0x1d, 0x55, 0xe3, 0xdb, 0xf1, 0x1e, 0x4e, 0xe7, 0x00,
	0xbc, 0x1c, 0x74, 0x3b, 0xea, 0x6e, 0xe4, 0xff, 0xc5, 0xe8, 0x52, 0x7b, 0x2f, 0x87, 0xf9, 0x07,
	0xef, 0xbe, 0x02, 0x00, 0x00, 0xff, 0xff, 0xf3, 0x12, 0x57, 0xe0, 0x57, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// FileServiceClient is the client API for FileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FileServiceClient interface {
	//simple RPC
	RegisterUser(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*RegisterUserResponse, error)
}

type fileServiceClient struct {
	cc *grpc.ClientConn
}

func NewFileServiceClient(cc *grpc.ClientConn) FileServiceClient {
	return &fileServiceClient{cc}
}

func (c *fileServiceClient) RegisterUser(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*RegisterUserResponse, error) {
	out := new(RegisterUserResponse)
	err := c.cc.Invoke(ctx, "/v1.FileService/RegisterUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileServiceServer is the server API for FileService service.
type FileServiceServer interface {
	//simple RPC
	RegisterUser(context.Context, *RegisterUserRequest) (*RegisterUserResponse, error)
}

// UnimplementedFileServiceServer can be embedded to have forward compatible implementations.
type UnimplementedFileServiceServer struct {
}

func (*UnimplementedFileServiceServer) RegisterUser(ctx context.Context, req *RegisterUserRequest) (*RegisterUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterUser not implemented")
}

func RegisterFileServiceServer(s *grpc.Server, srv FileServiceServer) {
	s.RegisterService(&_FileService_serviceDesc, srv)
}

func _FileService_RegisterUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).RegisterUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.FileService/RegisterUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).RegisterUser(ctx, req.(*RegisterUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _FileService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.FileService",
	HandlerType: (*FileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterUser",
			Handler:    _FileService_RegisterUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "file-service.proto",
}
