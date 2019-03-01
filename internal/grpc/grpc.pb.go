// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc.proto

package grpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import empty "github.com/golang/protobuf/ptypes/empty"

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

// JobInfoRequest defines required input for retrieving job-related information.
type JobInfoRequest struct {
	Jid                  string   `protobuf:"bytes,1,opt,name=jid,proto3" json:"jid,omitempty"`
	Xml                  bool     `protobuf:"varint,2,opt,name=xml,proto3" json:"xml,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JobInfoRequest) Reset()         { *m = JobInfoRequest{} }
func (m *JobInfoRequest) String() string { return proto.CompactTextString(m) }
func (*JobInfoRequest) ProtoMessage()    {}
func (*JobInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_91cf05b1f35c2068, []int{0}
}
func (m *JobInfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JobInfoRequest.Unmarshal(m, b)
}
func (m *JobInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JobInfoRequest.Marshal(b, m, deterministic)
}
func (dst *JobInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JobInfoRequest.Merge(dst, src)
}
func (m *JobInfoRequest) XXX_Size() int {
	return xxx_messageInfo_JobInfoRequest.Size(m)
}
func (m *JobInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_JobInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_JobInfoRequest proto.InternalMessageInfo

func (m *JobInfoRequest) GetJid() string {
	if m != nil {
		return m.Jid
	}
	return ""
}

func (m *JobInfoRequest) GetXml() bool {
	if m != nil {
		return m.Xml
	}
	return false
}

// UserInfoRequest defines required input for retrieving user-related information.
type UserInfoRequest struct {
	Uid                  string   `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Xml                  bool     `protobuf:"varint,2,opt,name=xml,proto3" json:"xml,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserInfoRequest) Reset()         { *m = UserInfoRequest{} }
func (m *UserInfoRequest) String() string { return proto.CompactTextString(m) }
func (*UserInfoRequest) ProtoMessage()    {}
func (*UserInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_91cf05b1f35c2068, []int{1}
}
func (m *UserInfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserInfoRequest.Unmarshal(m, b)
}
func (m *UserInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserInfoRequest.Marshal(b, m, deterministic)
}
func (dst *UserInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserInfoRequest.Merge(dst, src)
}
func (m *UserInfoRequest) XXX_Size() int {
	return xxx_messageInfo_UserInfoRequest.Size(m)
}
func (m *UserInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UserInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UserInfoRequest proto.InternalMessageInfo

func (m *UserInfoRequest) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *UserInfoRequest) GetXml() bool {
	if m != nil {
		return m.Xml
	}
	return false
}

// GeneralResponse is a very simple and naive output message.
type GeneralResponse struct {
	ResponseData         string   `protobuf:"bytes,1,opt,name=responseData,proto3" json:"responseData,omitempty"`
	ExitCode             int32    `protobuf:"varint,2,opt,name=exitCode,proto3" json:"exitCode,omitempty"`
	ErrorMessage         string   `protobuf:"bytes,3,opt,name=errorMessage,proto3" json:"errorMessage,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GeneralResponse) Reset()         { *m = GeneralResponse{} }
func (m *GeneralResponse) String() string { return proto.CompactTextString(m) }
func (*GeneralResponse) ProtoMessage()    {}
func (*GeneralResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_91cf05b1f35c2068, []int{2}
}
func (m *GeneralResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeneralResponse.Unmarshal(m, b)
}
func (m *GeneralResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeneralResponse.Marshal(b, m, deterministic)
}
func (dst *GeneralResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeneralResponse.Merge(dst, src)
}
func (m *GeneralResponse) XXX_Size() int {
	return xxx_messageInfo_GeneralResponse.Size(m)
}
func (m *GeneralResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GeneralResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GeneralResponse proto.InternalMessageInfo

func (m *GeneralResponse) GetResponseData() string {
	if m != nil {
		return m.ResponseData
	}
	return ""
}

func (m *GeneralResponse) GetExitCode() int32 {
	if m != nil {
		return m.ExitCode
	}
	return 0
}

func (m *GeneralResponse) GetErrorMessage() string {
	if m != nil {
		return m.ErrorMessage
	}
	return ""
}

// ServerListResponse returns a list of server instances.  Each server instance has identifier and owner.
type ServerListResponse struct {
	ExitCode             int32                        `protobuf:"varint,1,opt,name=exitCode,proto3" json:"exitCode,omitempty"`
	ErrorMessage         string                       `protobuf:"bytes,2,opt,name=errorMessage,proto3" json:"errorMessage,omitempty"`
	Servers              []*ServerListResponse_Server `protobuf:"bytes,3,rep,name=servers,proto3" json:"servers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *ServerListResponse) Reset()         { *m = ServerListResponse{} }
func (m *ServerListResponse) String() string { return proto.CompactTextString(m) }
func (*ServerListResponse) ProtoMessage()    {}
func (*ServerListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_91cf05b1f35c2068, []int{3}
}
func (m *ServerListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerListResponse.Unmarshal(m, b)
}
func (m *ServerListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerListResponse.Marshal(b, m, deterministic)
}
func (dst *ServerListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerListResponse.Merge(dst, src)
}
func (m *ServerListResponse) XXX_Size() int {
	return xxx_messageInfo_ServerListResponse.Size(m)
}
func (m *ServerListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ServerListResponse proto.InternalMessageInfo

func (m *ServerListResponse) GetExitCode() int32 {
	if m != nil {
		return m.ExitCode
	}
	return 0
}

func (m *ServerListResponse) GetErrorMessage() string {
	if m != nil {
		return m.ErrorMessage
	}
	return ""
}

func (m *ServerListResponse) GetServers() []*ServerListResponse_Server {
	if m != nil {
		return m.Servers
	}
	return nil
}

type ServerListResponse_Server struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Owner                string   `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServerListResponse_Server) Reset()         { *m = ServerListResponse_Server{} }
func (m *ServerListResponse_Server) String() string { return proto.CompactTextString(m) }
func (*ServerListResponse_Server) ProtoMessage()    {}
func (*ServerListResponse_Server) Descriptor() ([]byte, []int) {
	return fileDescriptor_grpc_91cf05b1f35c2068, []int{3, 0}
}
func (m *ServerListResponse_Server) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerListResponse_Server.Unmarshal(m, b)
}
func (m *ServerListResponse_Server) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerListResponse_Server.Marshal(b, m, deterministic)
}
func (dst *ServerListResponse_Server) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerListResponse_Server.Merge(dst, src)
}
func (m *ServerListResponse_Server) XXX_Size() int {
	return xxx_messageInfo_ServerListResponse_Server.Size(m)
}
func (m *ServerListResponse_Server) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerListResponse_Server.DiscardUnknown(m)
}

var xxx_messageInfo_ServerListResponse_Server proto.InternalMessageInfo

func (m *ServerListResponse_Server) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ServerListResponse_Server) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func init() {
	proto.RegisterType((*JobInfoRequest)(nil), "grpc.JobInfoRequest")
	proto.RegisterType((*UserInfoRequest)(nil), "grpc.UserInfoRequest")
	proto.RegisterType((*GeneralResponse)(nil), "grpc.GeneralResponse")
	proto.RegisterType((*ServerListResponse)(nil), "grpc.ServerListResponse")
	proto.RegisterType((*ServerListResponse_Server)(nil), "grpc.ServerListResponse.Server")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TorqueHelperSrvServiceClient is the client API for TorqueHelperSrvService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TorqueHelperSrvServiceClient interface {
	Ping(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GeneralResponse, error)
	TraceJob(ctx context.Context, in *JobInfoRequest, opts ...grpc.CallOption) (*GeneralResponse, error)
	TorqueConfig(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GeneralResponse, error)
	MoabConfig(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GeneralResponse, error)
	GetJobBlockReason(ctx context.Context, in *JobInfoRequest, opts ...grpc.CallOption) (*GeneralResponse, error)
	GetBlockedJobsOfUser(ctx context.Context, in *UserInfoRequest, opts ...grpc.CallOption) (*GeneralResponse, error)
	Qstat(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GeneralResponse, error)
	Qstatx(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GeneralResponse, error)
}

type torqueHelperSrvServiceClient struct {
	cc *grpc.ClientConn
}

func NewTorqueHelperSrvServiceClient(cc *grpc.ClientConn) TorqueHelperSrvServiceClient {
	return &torqueHelperSrvServiceClient{cc}
}

func (c *torqueHelperSrvServiceClient) Ping(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GeneralResponse, error) {
	out := new(GeneralResponse)
	err := c.cc.Invoke(ctx, "/grpc.TorqueHelperSrvService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *torqueHelperSrvServiceClient) TraceJob(ctx context.Context, in *JobInfoRequest, opts ...grpc.CallOption) (*GeneralResponse, error) {
	out := new(GeneralResponse)
	err := c.cc.Invoke(ctx, "/grpc.TorqueHelperSrvService/TraceJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *torqueHelperSrvServiceClient) TorqueConfig(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GeneralResponse, error) {
	out := new(GeneralResponse)
	err := c.cc.Invoke(ctx, "/grpc.TorqueHelperSrvService/TorqueConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *torqueHelperSrvServiceClient) MoabConfig(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GeneralResponse, error) {
	out := new(GeneralResponse)
	err := c.cc.Invoke(ctx, "/grpc.TorqueHelperSrvService/MoabConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *torqueHelperSrvServiceClient) GetJobBlockReason(ctx context.Context, in *JobInfoRequest, opts ...grpc.CallOption) (*GeneralResponse, error) {
	out := new(GeneralResponse)
	err := c.cc.Invoke(ctx, "/grpc.TorqueHelperSrvService/GetJobBlockReason", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *torqueHelperSrvServiceClient) GetBlockedJobsOfUser(ctx context.Context, in *UserInfoRequest, opts ...grpc.CallOption) (*GeneralResponse, error) {
	out := new(GeneralResponse)
	err := c.cc.Invoke(ctx, "/grpc.TorqueHelperSrvService/GetBlockedJobsOfUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *torqueHelperSrvServiceClient) Qstat(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GeneralResponse, error) {
	out := new(GeneralResponse)
	err := c.cc.Invoke(ctx, "/grpc.TorqueHelperSrvService/Qstat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *torqueHelperSrvServiceClient) Qstatx(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GeneralResponse, error) {
	out := new(GeneralResponse)
	err := c.cc.Invoke(ctx, "/grpc.TorqueHelperSrvService/Qstatx", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TorqueHelperSrvServiceServer is the server API for TorqueHelperSrvService service.
type TorqueHelperSrvServiceServer interface {
	Ping(context.Context, *empty.Empty) (*GeneralResponse, error)
	TraceJob(context.Context, *JobInfoRequest) (*GeneralResponse, error)
	TorqueConfig(context.Context, *empty.Empty) (*GeneralResponse, error)
	MoabConfig(context.Context, *empty.Empty) (*GeneralResponse, error)
	GetJobBlockReason(context.Context, *JobInfoRequest) (*GeneralResponse, error)
	GetBlockedJobsOfUser(context.Context, *UserInfoRequest) (*GeneralResponse, error)
	Qstat(context.Context, *empty.Empty) (*GeneralResponse, error)
	Qstatx(context.Context, *empty.Empty) (*GeneralResponse, error)
}

func RegisterTorqueHelperSrvServiceServer(s *grpc.Server, srv TorqueHelperSrvServiceServer) {
	s.RegisterService(&_TorqueHelperSrvService_serviceDesc, srv)
}

func _TorqueHelperSrvService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TorqueHelperSrvServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.TorqueHelperSrvService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TorqueHelperSrvServiceServer).Ping(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TorqueHelperSrvService_TraceJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TorqueHelperSrvServiceServer).TraceJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.TorqueHelperSrvService/TraceJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TorqueHelperSrvServiceServer).TraceJob(ctx, req.(*JobInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TorqueHelperSrvService_TorqueConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TorqueHelperSrvServiceServer).TorqueConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.TorqueHelperSrvService/TorqueConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TorqueHelperSrvServiceServer).TorqueConfig(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TorqueHelperSrvService_MoabConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TorqueHelperSrvServiceServer).MoabConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.TorqueHelperSrvService/MoabConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TorqueHelperSrvServiceServer).MoabConfig(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TorqueHelperSrvService_GetJobBlockReason_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TorqueHelperSrvServiceServer).GetJobBlockReason(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.TorqueHelperSrvService/GetJobBlockReason",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TorqueHelperSrvServiceServer).GetJobBlockReason(ctx, req.(*JobInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TorqueHelperSrvService_GetBlockedJobsOfUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TorqueHelperSrvServiceServer).GetBlockedJobsOfUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.TorqueHelperSrvService/GetBlockedJobsOfUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TorqueHelperSrvServiceServer).GetBlockedJobsOfUser(ctx, req.(*UserInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TorqueHelperSrvService_Qstat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TorqueHelperSrvServiceServer).Qstat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.TorqueHelperSrvService/Qstat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TorqueHelperSrvServiceServer).Qstat(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TorqueHelperSrvService_Qstatx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TorqueHelperSrvServiceServer).Qstatx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.TorqueHelperSrvService/Qstatx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TorqueHelperSrvServiceServer).Qstatx(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _TorqueHelperSrvService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.TorqueHelperSrvService",
	HandlerType: (*TorqueHelperSrvServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _TorqueHelperSrvService_Ping_Handler,
		},
		{
			MethodName: "TraceJob",
			Handler:    _TorqueHelperSrvService_TraceJob_Handler,
		},
		{
			MethodName: "TorqueConfig",
			Handler:    _TorqueHelperSrvService_TorqueConfig_Handler,
		},
		{
			MethodName: "MoabConfig",
			Handler:    _TorqueHelperSrvService_MoabConfig_Handler,
		},
		{
			MethodName: "GetJobBlockReason",
			Handler:    _TorqueHelperSrvService_GetJobBlockReason_Handler,
		},
		{
			MethodName: "GetBlockedJobsOfUser",
			Handler:    _TorqueHelperSrvService_GetBlockedJobsOfUser_Handler,
		},
		{
			MethodName: "Qstat",
			Handler:    _TorqueHelperSrvService_Qstat_Handler,
		},
		{
			MethodName: "Qstatx",
			Handler:    _TorqueHelperSrvService_Qstatx_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc.proto",
}

// TorqueHelperMomServiceClient is the client API for TorqueHelperMomService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TorqueHelperMomServiceClient interface {
	JobMemInfo(ctx context.Context, in *JobInfoRequest, opts ...grpc.CallOption) (*GeneralResponse, error)
}

type torqueHelperMomServiceClient struct {
	cc *grpc.ClientConn
}

func NewTorqueHelperMomServiceClient(cc *grpc.ClientConn) TorqueHelperMomServiceClient {
	return &torqueHelperMomServiceClient{cc}
}

func (c *torqueHelperMomServiceClient) JobMemInfo(ctx context.Context, in *JobInfoRequest, opts ...grpc.CallOption) (*GeneralResponse, error) {
	out := new(GeneralResponse)
	err := c.cc.Invoke(ctx, "/grpc.TorqueHelperMomService/JobMemInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TorqueHelperMomServiceServer is the server API for TorqueHelperMomService service.
type TorqueHelperMomServiceServer interface {
	JobMemInfo(context.Context, *JobInfoRequest) (*GeneralResponse, error)
}

func RegisterTorqueHelperMomServiceServer(s *grpc.Server, srv TorqueHelperMomServiceServer) {
	s.RegisterService(&_TorqueHelperMomService_serviceDesc, srv)
}

func _TorqueHelperMomService_JobMemInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TorqueHelperMomServiceServer).JobMemInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.TorqueHelperMomService/JobMemInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TorqueHelperMomServiceServer).JobMemInfo(ctx, req.(*JobInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TorqueHelperMomService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.TorqueHelperMomService",
	HandlerType: (*TorqueHelperMomServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "JobMemInfo",
			Handler:    _TorqueHelperMomService_JobMemInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc.proto",
}

// TorqueHelperAccServiceClient is the client API for TorqueHelperAccService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TorqueHelperAccServiceClient interface {
	GetVNCServers(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ServerListResponse, error)
}

type torqueHelperAccServiceClient struct {
	cc *grpc.ClientConn
}

func NewTorqueHelperAccServiceClient(cc *grpc.ClientConn) TorqueHelperAccServiceClient {
	return &torqueHelperAccServiceClient{cc}
}

func (c *torqueHelperAccServiceClient) GetVNCServers(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ServerListResponse, error) {
	out := new(ServerListResponse)
	err := c.cc.Invoke(ctx, "/grpc.TorqueHelperAccService/GetVNCServers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TorqueHelperAccServiceServer is the server API for TorqueHelperAccService service.
type TorqueHelperAccServiceServer interface {
	GetVNCServers(context.Context, *empty.Empty) (*ServerListResponse, error)
}

func RegisterTorqueHelperAccServiceServer(s *grpc.Server, srv TorqueHelperAccServiceServer) {
	s.RegisterService(&_TorqueHelperAccService_serviceDesc, srv)
}

func _TorqueHelperAccService_GetVNCServers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TorqueHelperAccServiceServer).GetVNCServers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.TorqueHelperAccService/GetVNCServers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TorqueHelperAccServiceServer).GetVNCServers(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _TorqueHelperAccService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.TorqueHelperAccService",
	HandlerType: (*TorqueHelperAccServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVNCServers",
			Handler:    _TorqueHelperAccService_GetVNCServers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc.proto",
}

func init() { proto.RegisterFile("grpc.proto", fileDescriptor_grpc_91cf05b1f35c2068) }

var fileDescriptor_grpc_91cf05b1f35c2068 = []byte{
	// 461 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xe5, 0xb8, 0x09, 0x61, 0x28, 0x2d, 0xac, 0x42, 0x65, 0x85, 0x03, 0x51, 0x4e, 0x39,
	0xb9, 0x52, 0xa0, 0xa0, 0x0a, 0x21, 0x51, 0x52, 0x14, 0xb0, 0x08, 0x7f, 0x9c, 0xc2, 0x85, 0x93,
	0xd7, 0x99, 0x58, 0x06, 0xc7, 0xe3, 0xee, 0xae, 0x4b, 0x78, 0x3a, 0x9e, 0x0c, 0x09, 0xad, 0xd7,
	0x09, 0x4d, 0x82, 0x85, 0xdc, 0xdb, 0xcc, 0xa7, 0xf9, 0x7d, 0x33, 0x9a, 0x19, 0x80, 0x48, 0x64,
	0xa1, 0x9b, 0x09, 0x52, 0xc4, 0xf6, 0x74, 0xdc, 0x7d, 0x18, 0x11, 0x45, 0x09, 0x1e, 0x17, 0x1a,
	0xcf, 0xe7, 0xc7, 0xb8, 0xc8, 0xd4, 0x4f, 0x53, 0xd2, 0x7f, 0x02, 0x07, 0x1e, 0xf1, 0xb7, 0xe9,
	0x9c, 0x7c, 0xbc, 0xcc, 0x51, 0x2a, 0x76, 0x0f, 0xec, 0x6f, 0xf1, 0xcc, 0xb1, 0x7a, 0xd6, 0xe0,
	0xb6, 0xaf, 0x43, 0xad, 0x2c, 0x17, 0x89, 0xd3, 0xe8, 0x59, 0x83, 0xb6, 0xaf, 0xc3, 0xfe, 0x09,
	0x1c, 0x7e, 0x96, 0x28, 0xb6, 0xb0, 0xfc, 0x2f, 0x96, 0xff, 0x13, 0xcb, 0xe1, 0x70, 0x8c, 0x29,
	0x8a, 0x20, 0xf1, 0x51, 0x66, 0x94, 0x4a, 0x64, 0x7d, 0xd8, 0x17, 0x65, 0x7c, 0x1e, 0xa8, 0xa0,
	0xe4, 0x37, 0x34, 0xd6, 0x85, 0x36, 0x2e, 0x63, 0x35, 0xa2, 0x19, 0x16, 0x6e, 0x4d, 0x7f, 0x9d,
	0x6b, 0x1e, 0x85, 0x20, 0x31, 0x41, 0x29, 0x83, 0x08, 0x1d, 0xdb, 0xf0, 0xd7, 0xb5, 0xfe, 0x2f,
	0x0b, 0xd8, 0x14, 0xc5, 0x15, 0x8a, 0x77, 0xb1, 0x54, 0xeb, 0xd6, 0xd7, 0x6d, 0xad, 0xff, 0xd8,
	0x36, 0x76, 0x6d, 0xd9, 0x29, 0xdc, 0x92, 0x85, 0xab, 0x74, 0xec, 0x9e, 0x3d, 0xb8, 0x33, 0x7c,
	0xe4, 0x16, 0xbb, 0xdf, 0x6d, 0x55, 0x4a, 0xfe, 0xaa, 0xbe, 0xeb, 0x42, 0xcb, 0x48, 0xec, 0x00,
	0x1a, 0xeb, 0xad, 0x35, 0xe2, 0x19, 0xeb, 0x40, 0x93, 0x7e, 0xa4, 0x28, 0xca, 0x8e, 0x26, 0x19,
	0xfe, 0xb6, 0xe1, 0xe8, 0x82, 0xc4, 0x65, 0x8e, 0x6f, 0x30, 0xc9, 0x50, 0x4c, 0xc5, 0x95, 0xe6,
	0xe3, 0x10, 0xd9, 0x09, 0xec, 0x7d, 0x8c, 0xd3, 0x88, 0x1d, 0xb9, 0xe6, 0xcc, 0xee, 0xea, 0xcc,
	0xee, 0x6b, 0x7d, 0xe6, 0xee, 0x03, 0x33, 0xd4, 0xf6, 0xde, 0x9f, 0x41, 0xfb, 0x42, 0x04, 0x21,
	0x7a, 0xc4, 0x59, 0xc7, 0x94, 0x6c, 0xfe, 0x41, 0x15, 0xf8, 0x02, 0xf6, 0xcd, 0x24, 0x23, 0x4a,
	0xe7, 0x71, 0xed, 0xbe, 0xcf, 0x01, 0x26, 0x14, 0xf0, 0x9b, 0xc1, 0x2f, 0xe1, 0xfe, 0x18, 0x95,
	0x47, 0xfc, 0x55, 0x42, 0xe1, 0x77, 0x1f, 0x03, 0x49, 0x69, 0xbd, 0xe9, 0xcf, 0xa1, 0x33, 0x46,
	0x55, 0xe0, 0x38, 0xf3, 0x88, 0xcb, 0x0f, 0x73, 0xfd, 0xc8, 0xac, 0x2c, 0xdf, 0x7a, 0xea, 0x2a,
	0x97, 0xa7, 0xd0, 0xfc, 0x24, 0x55, 0xa0, 0xea, 0x2f, 0xbd, 0x55, 0x70, 0xcb, 0x9a, 0xe0, 0x70,
	0xba, 0x79, 0xfe, 0x09, 0x2d, 0x56, 0xe7, 0x3f, 0x05, 0xf0, 0x88, 0x4f, 0x70, 0xa1, 0xc7, 0xae,
	0xb5, 0x8b, 0xe1, 0xd7, 0x4d, 0xd3, 0xb3, 0x30, 0x5c, 0x99, 0x9e, 0xc1, 0xdd, 0x31, 0xaa, 0x2f,
	0xef, 0x47, 0xe6, 0x49, 0x65, 0xe5, 0xb8, 0x4e, 0xd5, 0xc7, 0xf3, 0x56, 0x51, 0xf9, 0xf8, 0x4f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x42, 0x58, 0xc0, 0x78, 0x8f, 0x04, 0x00, 0x00,
}
