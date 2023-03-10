// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: proto/service.proto

package id_generator

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// IdServiceClient is the example-client API for IdService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IdServiceClient interface {
	GenerateId(ctx context.Context, in *GenerateIdRequest, opts ...grpc.CallOption) (*GenerateIdResponse, error)
	GenerateMultipleIds(ctx context.Context, in *GenerateMultipleIdsRequest, opts ...grpc.CallOption) (*GenerateMultipleIdsResponse, error)
}

type idServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIdServiceClient(cc grpc.ClientConnInterface) IdServiceClient {
	return &idServiceClient{cc}
}

func (c *idServiceClient) GenerateId(ctx context.Context, in *GenerateIdRequest, opts ...grpc.CallOption) (*GenerateIdResponse, error) {
	out := new(GenerateIdResponse)
	err := c.cc.Invoke(ctx, "/IdService/GenerateId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *idServiceClient) GenerateMultipleIds(ctx context.Context, in *GenerateMultipleIdsRequest, opts ...grpc.CallOption) (*GenerateMultipleIdsResponse, error) {
	out := new(GenerateMultipleIdsResponse)
	err := c.cc.Invoke(ctx, "/IdService/GenerateMultipleIds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IdServiceServer is the server API for IdService service.
// All implementations must embed UnimplementedIdServiceServer
// for forward compatibility
type IdServiceServer interface {
	GenerateId(context.Context, *GenerateIdRequest) (*GenerateIdResponse, error)
	GenerateMultipleIds(context.Context, *GenerateMultipleIdsRequest) (*GenerateMultipleIdsResponse, error)
	mustEmbedUnimplementedIdServiceServer()
}

// UnimplementedIdServiceServer must be embedded to have forward compatible implementations.
type UnimplementedIdServiceServer struct {
}

func (UnimplementedIdServiceServer) GenerateId(context.Context, *GenerateIdRequest) (*GenerateIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateId not implemented")
}
func (UnimplementedIdServiceServer) GenerateMultipleIds(context.Context, *GenerateMultipleIdsRequest) (*GenerateMultipleIdsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateMultipleIds not implemented")
}
func (UnimplementedIdServiceServer) mustEmbedUnimplementedIdServiceServer() {}

// UnsafeIdServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IdServiceServer will
// result in compilation errors.
type UnsafeIdServiceServer interface {
	mustEmbedUnimplementedIdServiceServer()
}

func RegisterIdServiceServer(s grpc.ServiceRegistrar, srv IdServiceServer) {
	s.RegisterService(&IdService_ServiceDesc, srv)
}

func _IdService_GenerateId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdServiceServer).GenerateId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IdService/GenerateId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdServiceServer).GenerateId(ctx, req.(*GenerateIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdService_GenerateMultipleIds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateMultipleIdsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdServiceServer).GenerateMultipleIds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IdService/GenerateMultipleIds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdServiceServer).GenerateMultipleIds(ctx, req.(*GenerateMultipleIdsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// IdService_ServiceDesc is the grpc.ServiceDesc for IdService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IdService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "IdService",
	HandlerType: (*IdServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenerateId",
			Handler:    _IdService_GenerateId_Handler,
		},
		{
			MethodName: "GenerateMultipleIds",
			Handler:    _IdService_GenerateMultipleIds_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/service.proto",
}
