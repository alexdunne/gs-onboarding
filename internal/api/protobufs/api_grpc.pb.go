// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protobufs

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// APIClient is the client API for API service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type APIClient interface {
	ListAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_ListAllClient, error)
	ListStories(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_ListStoriesClient, error)
	ListJobs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_ListJobsClient, error)
}

type aPIClient struct {
	cc grpc.ClientConnInterface
}

func NewAPIClient(cc grpc.ClientConnInterface) APIClient {
	return &aPIClient{cc}
}

func (c *aPIClient) ListAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_ListAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[0], "/api.API/ListAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPIListAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_ListAllClient interface {
	Recv() (*Item, error)
	grpc.ClientStream
}

type aPIListAllClient struct {
	grpc.ClientStream
}

func (x *aPIListAllClient) Recv() (*Item, error) {
	m := new(Item)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) ListStories(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_ListStoriesClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[1], "/api.API/ListStories", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPIListStoriesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_ListStoriesClient interface {
	Recv() (*Item, error)
	grpc.ClientStream
}

type aPIListStoriesClient struct {
	grpc.ClientStream
}

func (x *aPIListStoriesClient) Recv() (*Item, error) {
	m := new(Item)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) ListJobs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_ListJobsClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[2], "/api.API/ListJobs", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPIListJobsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_ListJobsClient interface {
	Recv() (*Item, error)
	grpc.ClientStream
}

type aPIListJobsClient struct {
	grpc.ClientStream
}

func (x *aPIListJobsClient) Recv() (*Item, error) {
	m := new(Item)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// APIServer is the server API for API service.
// All implementations must embed UnimplementedAPIServer
// for forward compatibility
type APIServer interface {
	ListAll(*emptypb.Empty, API_ListAllServer) error
	ListStories(*emptypb.Empty, API_ListStoriesServer) error
	ListJobs(*emptypb.Empty, API_ListJobsServer) error
	mustEmbedUnimplementedAPIServer()
}

// UnimplementedAPIServer must be embedded to have forward compatible implementations.
type UnimplementedAPIServer struct {
}

func (UnimplementedAPIServer) ListAll(*emptypb.Empty, API_ListAllServer) error {
	return status.Errorf(codes.Unimplemented, "method ListAll not implemented")
}
func (UnimplementedAPIServer) ListStories(*emptypb.Empty, API_ListStoriesServer) error {
	return status.Errorf(codes.Unimplemented, "method ListStories not implemented")
}
func (UnimplementedAPIServer) ListJobs(*emptypb.Empty, API_ListJobsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListJobs not implemented")
}
func (UnimplementedAPIServer) mustEmbedUnimplementedAPIServer() {}

// UnsafeAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to APIServer will
// result in compilation errors.
type UnsafeAPIServer interface {
	mustEmbedUnimplementedAPIServer()
}

func RegisterAPIServer(s grpc.ServiceRegistrar, srv APIServer) {
	s.RegisterService(&API_ServiceDesc, srv)
}

func _API_ListAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).ListAll(m, &aPIListAllServer{stream})
}

type API_ListAllServer interface {
	Send(*Item) error
	grpc.ServerStream
}

type aPIListAllServer struct {
	grpc.ServerStream
}

func (x *aPIListAllServer) Send(m *Item) error {
	return x.ServerStream.SendMsg(m)
}

func _API_ListStories_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).ListStories(m, &aPIListStoriesServer{stream})
}

type API_ListStoriesServer interface {
	Send(*Item) error
	grpc.ServerStream
}

type aPIListStoriesServer struct {
	grpc.ServerStream
}

func (x *aPIListStoriesServer) Send(m *Item) error {
	return x.ServerStream.SendMsg(m)
}

func _API_ListJobs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).ListJobs(m, &aPIListJobsServer{stream})
}

type API_ListJobsServer interface {
	Send(*Item) error
	grpc.ServerStream
}

type aPIListJobsServer struct {
	grpc.ServerStream
}

func (x *aPIListJobsServer) Send(m *Item) error {
	return x.ServerStream.SendMsg(m)
}

// API_ServiceDesc is the grpc.ServiceDesc for API service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var API_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.API",
	HandlerType: (*APIServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListAll",
			Handler:       _API_ListAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ListStories",
			Handler:       _API_ListStories_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ListJobs",
			Handler:       _API_ListJobs_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api.proto",
}
