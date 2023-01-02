// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: api/app.proto

package api

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

// AppClient is the client API for App service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AppClient interface {
	EthAuthorize(ctx context.Context, in *EthAuthorizeRequest, opts ...grpc.CallOption) (*EthAuthorizeReply, error)
	Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*DepositReply, error)
	UserInfo(ctx context.Context, in *UserInfoRequest, opts ...grpc.CallOption) (*UserInfoReply, error)
	RewardList(ctx context.Context, in *RewardListRequest, opts ...grpc.CallOption) (*RewardListReply, error)
	RecommendRewardList(ctx context.Context, in *RecommendRewardListRequest, opts ...grpc.CallOption) (*RecommendRewardListReply, error)
	FeeRewardList(ctx context.Context, in *FeeRewardListRequest, opts ...grpc.CallOption) (*FeeRewardListReply, error)
	WithdrawList(ctx context.Context, in *WithdrawListRequest, opts ...grpc.CallOption) (*WithdrawListReply, error)
	Withdraw(ctx context.Context, in *WithdrawRequest, opts ...grpc.CallOption) (*WithdrawReply, error)
}

type appClient struct {
	cc grpc.ClientConnInterface
}

func NewAppClient(cc grpc.ClientConnInterface) AppClient {
	return &appClient{cc}
}

func (c *appClient) EthAuthorize(ctx context.Context, in *EthAuthorizeRequest, opts ...grpc.CallOption) (*EthAuthorizeReply, error) {
	out := new(EthAuthorizeReply)
	err := c.cc.Invoke(ctx, "/api.App/EthAuthorize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*DepositReply, error) {
	out := new(DepositReply)
	err := c.cc.Invoke(ctx, "/api.App/deposit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) UserInfo(ctx context.Context, in *UserInfoRequest, opts ...grpc.CallOption) (*UserInfoReply, error) {
	out := new(UserInfoReply)
	err := c.cc.Invoke(ctx, "/api.App/userInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) RewardList(ctx context.Context, in *RewardListRequest, opts ...grpc.CallOption) (*RewardListReply, error) {
	out := new(RewardListReply)
	err := c.cc.Invoke(ctx, "/api.App/RewardList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) RecommendRewardList(ctx context.Context, in *RecommendRewardListRequest, opts ...grpc.CallOption) (*RecommendRewardListReply, error) {
	out := new(RecommendRewardListReply)
	err := c.cc.Invoke(ctx, "/api.App/RecommendRewardList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) FeeRewardList(ctx context.Context, in *FeeRewardListRequest, opts ...grpc.CallOption) (*FeeRewardListReply, error) {
	out := new(FeeRewardListReply)
	err := c.cc.Invoke(ctx, "/api.App/FeeRewardList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) WithdrawList(ctx context.Context, in *WithdrawListRequest, opts ...grpc.CallOption) (*WithdrawListReply, error) {
	out := new(WithdrawListReply)
	err := c.cc.Invoke(ctx, "/api.App/withdrawList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) Withdraw(ctx context.Context, in *WithdrawRequest, opts ...grpc.CallOption) (*WithdrawReply, error) {
	out := new(WithdrawReply)
	err := c.cc.Invoke(ctx, "/api.App/withdraw", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AppServer is the server API for App service.
// All implementations must embed UnimplementedAppServer
// for forward compatibility
type AppServer interface {
	EthAuthorize(context.Context, *EthAuthorizeRequest) (*EthAuthorizeReply, error)
	Deposit(context.Context, *DepositRequest) (*DepositReply, error)
	UserInfo(context.Context, *UserInfoRequest) (*UserInfoReply, error)
	RewardList(context.Context, *RewardListRequest) (*RewardListReply, error)
	RecommendRewardList(context.Context, *RecommendRewardListRequest) (*RecommendRewardListReply, error)
	FeeRewardList(context.Context, *FeeRewardListRequest) (*FeeRewardListReply, error)
	WithdrawList(context.Context, *WithdrawListRequest) (*WithdrawListReply, error)
	Withdraw(context.Context, *WithdrawRequest) (*WithdrawReply, error)
	mustEmbedUnimplementedAppServer()
}

// UnimplementedAppServer must be embedded to have forward compatible implementations.
type UnimplementedAppServer struct {
}

func (UnimplementedAppServer) EthAuthorize(context.Context, *EthAuthorizeRequest) (*EthAuthorizeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EthAuthorize not implemented")
}
func (UnimplementedAppServer) Deposit(context.Context, *DepositRequest) (*DepositReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deposit not implemented")
}
func (UnimplementedAppServer) UserInfo(context.Context, *UserInfoRequest) (*UserInfoReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserInfo not implemented")
}
func (UnimplementedAppServer) RewardList(context.Context, *RewardListRequest) (*RewardListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RewardList not implemented")
}
func (UnimplementedAppServer) RecommendRewardList(context.Context, *RecommendRewardListRequest) (*RecommendRewardListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecommendRewardList not implemented")
}
func (UnimplementedAppServer) FeeRewardList(context.Context, *FeeRewardListRequest) (*FeeRewardListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FeeRewardList not implemented")
}
func (UnimplementedAppServer) WithdrawList(context.Context, *WithdrawListRequest) (*WithdrawListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WithdrawList not implemented")
}
func (UnimplementedAppServer) Withdraw(context.Context, *WithdrawRequest) (*WithdrawReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Withdraw not implemented")
}
func (UnimplementedAppServer) mustEmbedUnimplementedAppServer() {}

// UnsafeAppServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AppServer will
// result in compilation errors.
type UnsafeAppServer interface {
	mustEmbedUnimplementedAppServer()
}

func RegisterAppServer(s grpc.ServiceRegistrar, srv AppServer) {
	s.RegisterService(&App_ServiceDesc, srv)
}

func _App_EthAuthorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EthAuthorizeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).EthAuthorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.App/EthAuthorize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).EthAuthorize(ctx, req.(*EthAuthorizeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_Deposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).Deposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.App/deposit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).Deposit(ctx, req.(*DepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_UserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).UserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.App/userInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).UserInfo(ctx, req.(*UserInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_RewardList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RewardListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).RewardList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.App/RewardList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).RewardList(ctx, req.(*RewardListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_RecommendRewardList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecommendRewardListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).RecommendRewardList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.App/RecommendRewardList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).RecommendRewardList(ctx, req.(*RecommendRewardListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_FeeRewardList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FeeRewardListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).FeeRewardList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.App/FeeRewardList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).FeeRewardList(ctx, req.(*FeeRewardListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_WithdrawList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithdrawListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).WithdrawList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.App/withdrawList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).WithdrawList(ctx, req.(*WithdrawListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_Withdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).Withdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.App/withdraw",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).Withdraw(ctx, req.(*WithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// App_ServiceDesc is the grpc.ServiceDesc for App service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var App_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.App",
	HandlerType: (*AppServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EthAuthorize",
			Handler:    _App_EthAuthorize_Handler,
		},
		{
			MethodName: "deposit",
			Handler:    _App_Deposit_Handler,
		},
		{
			MethodName: "userInfo",
			Handler:    _App_UserInfo_Handler,
		},
		{
			MethodName: "RewardList",
			Handler:    _App_RewardList_Handler,
		},
		{
			MethodName: "RecommendRewardList",
			Handler:    _App_RecommendRewardList_Handler,
		},
		{
			MethodName: "FeeRewardList",
			Handler:    _App_FeeRewardList_Handler,
		},
		{
			MethodName: "withdrawList",
			Handler:    _App_WithdrawList_Handler,
		},
		{
			MethodName: "withdraw",
			Handler:    _App_Withdraw_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/app.proto",
}
