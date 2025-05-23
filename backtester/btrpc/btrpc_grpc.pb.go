// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: btrpc.proto

package btrpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	BacktesterService_ExecuteStrategyFromFile_FullMethodName   = "/btrpc.BacktesterService/ExecuteStrategyFromFile"
	BacktesterService_ExecuteStrategyFromConfig_FullMethodName = "/btrpc.BacktesterService/ExecuteStrategyFromConfig"
	BacktesterService_ListAllTasks_FullMethodName              = "/btrpc.BacktesterService/ListAllTasks"
	BacktesterService_StartTask_FullMethodName                 = "/btrpc.BacktesterService/StartTask"
	BacktesterService_StartAllTasks_FullMethodName             = "/btrpc.BacktesterService/StartAllTasks"
	BacktesterService_StopTask_FullMethodName                  = "/btrpc.BacktesterService/StopTask"
	BacktesterService_StopAllTasks_FullMethodName              = "/btrpc.BacktesterService/StopAllTasks"
	BacktesterService_ClearTask_FullMethodName                 = "/btrpc.BacktesterService/ClearTask"
	BacktesterService_ClearAllTasks_FullMethodName             = "/btrpc.BacktesterService/ClearAllTasks"
)

// BacktesterServiceClient is the client API for BacktesterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BacktesterServiceClient interface {
	ExecuteStrategyFromFile(ctx context.Context, in *ExecuteStrategyFromFileRequest, opts ...grpc.CallOption) (*ExecuteStrategyResponse, error)
	ExecuteStrategyFromConfig(ctx context.Context, in *ExecuteStrategyFromConfigRequest, opts ...grpc.CallOption) (*ExecuteStrategyResponse, error)
	ListAllTasks(ctx context.Context, in *ListAllTasksRequest, opts ...grpc.CallOption) (*ListAllTasksResponse, error)
	StartTask(ctx context.Context, in *StartTaskRequest, opts ...grpc.CallOption) (*StartTaskResponse, error)
	StartAllTasks(ctx context.Context, in *StartAllTasksRequest, opts ...grpc.CallOption) (*StartAllTasksResponse, error)
	StopTask(ctx context.Context, in *StopTaskRequest, opts ...grpc.CallOption) (*StopTaskResponse, error)
	StopAllTasks(ctx context.Context, in *StopAllTasksRequest, opts ...grpc.CallOption) (*StopAllTasksResponse, error)
	ClearTask(ctx context.Context, in *ClearTaskRequest, opts ...grpc.CallOption) (*ClearTaskResponse, error)
	ClearAllTasks(ctx context.Context, in *ClearAllTasksRequest, opts ...grpc.CallOption) (*ClearAllTasksResponse, error)
}

type backtesterServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBacktesterServiceClient(cc grpc.ClientConnInterface) BacktesterServiceClient {
	return &backtesterServiceClient{cc}
}

func (c *backtesterServiceClient) ExecuteStrategyFromFile(ctx context.Context, in *ExecuteStrategyFromFileRequest, opts ...grpc.CallOption) (*ExecuteStrategyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExecuteStrategyResponse)
	err := c.cc.Invoke(ctx, BacktesterService_ExecuteStrategyFromFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backtesterServiceClient) ExecuteStrategyFromConfig(ctx context.Context, in *ExecuteStrategyFromConfigRequest, opts ...grpc.CallOption) (*ExecuteStrategyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExecuteStrategyResponse)
	err := c.cc.Invoke(ctx, BacktesterService_ExecuteStrategyFromConfig_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backtesterServiceClient) ListAllTasks(ctx context.Context, in *ListAllTasksRequest, opts ...grpc.CallOption) (*ListAllTasksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListAllTasksResponse)
	err := c.cc.Invoke(ctx, BacktesterService_ListAllTasks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backtesterServiceClient) StartTask(ctx context.Context, in *StartTaskRequest, opts ...grpc.CallOption) (*StartTaskResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StartTaskResponse)
	err := c.cc.Invoke(ctx, BacktesterService_StartTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backtesterServiceClient) StartAllTasks(ctx context.Context, in *StartAllTasksRequest, opts ...grpc.CallOption) (*StartAllTasksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StartAllTasksResponse)
	err := c.cc.Invoke(ctx, BacktesterService_StartAllTasks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backtesterServiceClient) StopTask(ctx context.Context, in *StopTaskRequest, opts ...grpc.CallOption) (*StopTaskResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StopTaskResponse)
	err := c.cc.Invoke(ctx, BacktesterService_StopTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backtesterServiceClient) StopAllTasks(ctx context.Context, in *StopAllTasksRequest, opts ...grpc.CallOption) (*StopAllTasksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StopAllTasksResponse)
	err := c.cc.Invoke(ctx, BacktesterService_StopAllTasks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backtesterServiceClient) ClearTask(ctx context.Context, in *ClearTaskRequest, opts ...grpc.CallOption) (*ClearTaskResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClearTaskResponse)
	err := c.cc.Invoke(ctx, BacktesterService_ClearTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backtesterServiceClient) ClearAllTasks(ctx context.Context, in *ClearAllTasksRequest, opts ...grpc.CallOption) (*ClearAllTasksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClearAllTasksResponse)
	err := c.cc.Invoke(ctx, BacktesterService_ClearAllTasks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BacktesterServiceServer is the server API for BacktesterService service.
// All implementations must embed UnimplementedBacktesterServiceServer
// for forward compatibility
type BacktesterServiceServer interface {
	ExecuteStrategyFromFile(context.Context, *ExecuteStrategyFromFileRequest) (*ExecuteStrategyResponse, error)
	ExecuteStrategyFromConfig(context.Context, *ExecuteStrategyFromConfigRequest) (*ExecuteStrategyResponse, error)
	ListAllTasks(context.Context, *ListAllTasksRequest) (*ListAllTasksResponse, error)
	StartTask(context.Context, *StartTaskRequest) (*StartTaskResponse, error)
	StartAllTasks(context.Context, *StartAllTasksRequest) (*StartAllTasksResponse, error)
	StopTask(context.Context, *StopTaskRequest) (*StopTaskResponse, error)
	StopAllTasks(context.Context, *StopAllTasksRequest) (*StopAllTasksResponse, error)
	ClearTask(context.Context, *ClearTaskRequest) (*ClearTaskResponse, error)
	ClearAllTasks(context.Context, *ClearAllTasksRequest) (*ClearAllTasksResponse, error)
	mustEmbedUnimplementedBacktesterServiceServer()
}

// UnimplementedBacktesterServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBacktesterServiceServer struct {
}

func (UnimplementedBacktesterServiceServer) ExecuteStrategyFromFile(context.Context, *ExecuteStrategyFromFileRequest) (*ExecuteStrategyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteStrategyFromFile not implemented")
}
func (UnimplementedBacktesterServiceServer) ExecuteStrategyFromConfig(context.Context, *ExecuteStrategyFromConfigRequest) (*ExecuteStrategyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteStrategyFromConfig not implemented")
}
func (UnimplementedBacktesterServiceServer) ListAllTasks(context.Context, *ListAllTasksRequest) (*ListAllTasksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAllTasks not implemented")
}
func (UnimplementedBacktesterServiceServer) StartTask(context.Context, *StartTaskRequest) (*StartTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartTask not implemented")
}
func (UnimplementedBacktesterServiceServer) StartAllTasks(context.Context, *StartAllTasksRequest) (*StartAllTasksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartAllTasks not implemented")
}
func (UnimplementedBacktesterServiceServer) StopTask(context.Context, *StopTaskRequest) (*StopTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopTask not implemented")
}
func (UnimplementedBacktesterServiceServer) StopAllTasks(context.Context, *StopAllTasksRequest) (*StopAllTasksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopAllTasks not implemented")
}
func (UnimplementedBacktesterServiceServer) ClearTask(context.Context, *ClearTaskRequest) (*ClearTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearTask not implemented")
}
func (UnimplementedBacktesterServiceServer) ClearAllTasks(context.Context, *ClearAllTasksRequest) (*ClearAllTasksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearAllTasks not implemented")
}
func (UnimplementedBacktesterServiceServer) mustEmbedUnimplementedBacktesterServiceServer() {}

// UnsafeBacktesterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BacktesterServiceServer will
// result in compilation errors.
type UnsafeBacktesterServiceServer interface {
	mustEmbedUnimplementedBacktesterServiceServer()
}

func RegisterBacktesterServiceServer(s grpc.ServiceRegistrar, srv BacktesterServiceServer) {
	s.RegisterService(&BacktesterService_ServiceDesc, srv)
}

func _BacktesterService_ExecuteStrategyFromFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteStrategyFromFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktesterServiceServer).ExecuteStrategyFromFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BacktesterService_ExecuteStrategyFromFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktesterServiceServer).ExecuteStrategyFromFile(ctx, req.(*ExecuteStrategyFromFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BacktesterService_ExecuteStrategyFromConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteStrategyFromConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktesterServiceServer).ExecuteStrategyFromConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BacktesterService_ExecuteStrategyFromConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktesterServiceServer).ExecuteStrategyFromConfig(ctx, req.(*ExecuteStrategyFromConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BacktesterService_ListAllTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAllTasksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktesterServiceServer).ListAllTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BacktesterService_ListAllTasks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktesterServiceServer).ListAllTasks(ctx, req.(*ListAllTasksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BacktesterService_StartTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktesterServiceServer).StartTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BacktesterService_StartTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktesterServiceServer).StartTask(ctx, req.(*StartTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BacktesterService_StartAllTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartAllTasksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktesterServiceServer).StartAllTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BacktesterService_StartAllTasks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktesterServiceServer).StartAllTasks(ctx, req.(*StartAllTasksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BacktesterService_StopTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktesterServiceServer).StopTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BacktesterService_StopTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktesterServiceServer).StopTask(ctx, req.(*StopTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BacktesterService_StopAllTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopAllTasksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktesterServiceServer).StopAllTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BacktesterService_StopAllTasks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktesterServiceServer).StopAllTasks(ctx, req.(*StopAllTasksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BacktesterService_ClearTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktesterServiceServer).ClearTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BacktesterService_ClearTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktesterServiceServer).ClearTask(ctx, req.(*ClearTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BacktesterService_ClearAllTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearAllTasksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktesterServiceServer).ClearAllTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BacktesterService_ClearAllTasks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktesterServiceServer).ClearAllTasks(ctx, req.(*ClearAllTasksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BacktesterService_ServiceDesc is the grpc.ServiceDesc for BacktesterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BacktesterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "btrpc.BacktesterService",
	HandlerType: (*BacktesterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExecuteStrategyFromFile",
			Handler:    _BacktesterService_ExecuteStrategyFromFile_Handler,
		},
		{
			MethodName: "ExecuteStrategyFromConfig",
			Handler:    _BacktesterService_ExecuteStrategyFromConfig_Handler,
		},
		{
			MethodName: "ListAllTasks",
			Handler:    _BacktesterService_ListAllTasks_Handler,
		},
		{
			MethodName: "StartTask",
			Handler:    _BacktesterService_StartTask_Handler,
		},
		{
			MethodName: "StartAllTasks",
			Handler:    _BacktesterService_StartAllTasks_Handler,
		},
		{
			MethodName: "StopTask",
			Handler:    _BacktesterService_StopTask_Handler,
		},
		{
			MethodName: "StopAllTasks",
			Handler:    _BacktesterService_StopAllTasks_Handler,
		},
		{
			MethodName: "ClearTask",
			Handler:    _BacktesterService_ClearTask_Handler,
		},
		{
			MethodName: "ClearAllTasks",
			Handler:    _BacktesterService_ClearAllTasks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "btrpc.proto",
}
