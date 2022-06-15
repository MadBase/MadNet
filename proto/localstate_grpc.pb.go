// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// LocalStateClient is the client API for LocalState service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LocalStateClient interface {
	// Get only the raw data from a datastore UTXO that has been mined into chain
	GetData(ctx context.Context, in *GetDataRequest, opts ...grpc.CallOption) (*GetDataResponse, error)
	// Get a list of UTXOs that sum to at least a minimum of some value where each
	// UTXO has a common owner
	GetValueForOwner(ctx context.Context, in *GetValueRequest, opts ...grpc.CallOption) (*GetValueResponse, error)
	// Iterate all datastores in a namespace defined by an owner
	IterateNameSpace(ctx context.Context, in *IterateNameSpaceRequest, opts ...grpc.CallOption) (*IterateNameSpaceResponse, error)
	// Get a mined transaction by hash
	GetMinedTransaction(ctx context.Context, in *MinedTransactionRequest, opts ...grpc.CallOption) (*MinedTransactionResponse, error)
	// Get blockheader by hash or blocknumber
	GetBlockHeader(ctx context.Context, in *BlockHeaderRequest, opts ...grpc.CallOption) (*BlockHeaderResponse, error)
	// Get a raw UTXO by TxHash and index or by UTXOID
	GetUTXO(ctx context.Context, in *UTXORequest, opts ...grpc.CallOption) (*UTXOResponse, error)
	// Get transaction status by hash
	GetTransactionStatus(ctx context.Context, in *TransactionStatusRequest, opts ...grpc.CallOption) (*TransactionStatusResponse, error)
	// Get a pending transaction by hash
	GetPendingTransaction(ctx context.Context, in *PendingTransactionRequest, opts ...grpc.CallOption) (*PendingTransactionResponse, error)
	// Get the round state object for a specified round for a specified validator
	// This allows tracing the consensus flow.
	GetRoundStateForValidator(ctx context.Context, in *RoundStateForValidatorRequest, opts ...grpc.CallOption) (*RoundStateForValidatorResponse, error)
	// Get the set of validators for a specified block height
	GetValidatorSet(ctx context.Context, in *ValidatorSetRequest, opts ...grpc.CallOption) (*ValidatorSetResponse, error)
	// Get the current block number
	GetBlockNumber(ctx context.Context, in *BlockNumberRequest, opts ...grpc.CallOption) (*BlockNumberResponse, error)
	// Get the current ChainID of the node
	GetChainID(ctx context.Context, in *ChainIDRequest, opts ...grpc.CallOption) (*ChainIDResponse, error)
	// Send a transaction to the node
	SendTransaction(ctx context.Context, in *TransactionData, opts ...grpc.CallOption) (*TransactionDetails, error)
	// Get the current block number
	GetEpochNumber(ctx context.Context, in *EpochNumberRequest, opts ...grpc.CallOption) (*EpochNumberResponse, error)
	// Get the current block number
	GetTxBlockNumber(ctx context.Context, in *TxBlockNumberRequest, opts ...grpc.CallOption) (*TxBlockNumberResponse, error)
	GetFees(ctx context.Context, in *FeeRequest, opts ...grpc.CallOption) (*FeeResponse, error)
}

type localStateClient struct {
	cc grpc.ClientConnInterface
}

func NewLocalStateClient(cc grpc.ClientConnInterface) LocalStateClient {
	return &localStateClient{cc}
}

func (c *localStateClient) GetData(ctx context.Context, in *GetDataRequest, opts ...grpc.CallOption) (*GetDataResponse, error) {
	out := new(GetDataResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/GetData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) GetValueForOwner(ctx context.Context, in *GetValueRequest, opts ...grpc.CallOption) (*GetValueResponse, error) {
	out := new(GetValueResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/GetValueForOwner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) IterateNameSpace(ctx context.Context, in *IterateNameSpaceRequest, opts ...grpc.CallOption) (*IterateNameSpaceResponse, error) {
	out := new(IterateNameSpaceResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/IterateNameSpace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) GetMinedTransaction(ctx context.Context, in *MinedTransactionRequest, opts ...grpc.CallOption) (*MinedTransactionResponse, error) {
	out := new(MinedTransactionResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/GetMinedTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) GetBlockHeader(ctx context.Context, in *BlockHeaderRequest, opts ...grpc.CallOption) (*BlockHeaderResponse, error) {
	out := new(BlockHeaderResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/GetBlockHeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) GetUTXO(ctx context.Context, in *UTXORequest, opts ...grpc.CallOption) (*UTXOResponse, error) {
	out := new(UTXOResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/GetUTXO", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) GetTransactionStatus(ctx context.Context, in *TransactionStatusRequest, opts ...grpc.CallOption) (*TransactionStatusResponse, error) {
	out := new(TransactionStatusResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/GetTransactionStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) GetPendingTransaction(ctx context.Context, in *PendingTransactionRequest, opts ...grpc.CallOption) (*PendingTransactionResponse, error) {
	out := new(PendingTransactionResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/GetPendingTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) GetRoundStateForValidator(ctx context.Context, in *RoundStateForValidatorRequest, opts ...grpc.CallOption) (*RoundStateForValidatorResponse, error) {
	out := new(RoundStateForValidatorResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/GetRoundStateForValidator", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) GetValidatorSet(ctx context.Context, in *ValidatorSetRequest, opts ...grpc.CallOption) (*ValidatorSetResponse, error) {
	out := new(ValidatorSetResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/GetValidatorSet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) GetBlockNumber(ctx context.Context, in *BlockNumberRequest, opts ...grpc.CallOption) (*BlockNumberResponse, error) {
	out := new(BlockNumberResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/GetBlockNumber", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) GetChainID(ctx context.Context, in *ChainIDRequest, opts ...grpc.CallOption) (*ChainIDResponse, error) {
	out := new(ChainIDResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/GetChainID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) SendTransaction(ctx context.Context, in *TransactionData, opts ...grpc.CallOption) (*TransactionDetails, error) {
	out := new(TransactionDetails)
	err := c.cc.Invoke(ctx, "/proto.LocalState/SendTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) GetEpochNumber(ctx context.Context, in *EpochNumberRequest, opts ...grpc.CallOption) (*EpochNumberResponse, error) {
	out := new(EpochNumberResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/GetEpochNumber", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) GetTxBlockNumber(ctx context.Context, in *TxBlockNumberRequest, opts ...grpc.CallOption) (*TxBlockNumberResponse, error) {
	out := new(TxBlockNumberResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/GetTxBlockNumber", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localStateClient) GetFees(ctx context.Context, in *FeeRequest, opts ...grpc.CallOption) (*FeeResponse, error) {
	out := new(FeeResponse)
	err := c.cc.Invoke(ctx, "/proto.LocalState/GetFees", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LocalStateServer is the server API for LocalState service.
// All implementations should embed UnimplementedLocalStateServer
// for forward compatibility
type LocalStateServer interface {
	// Get only the raw data from a datastore UTXO that has been mined into chain
	GetData(context.Context, *GetDataRequest) (*GetDataResponse, error)
	// Get a list of UTXOs that sum to at least a minimum of some value where each
	// UTXO has a common owner
	GetValueForOwner(context.Context, *GetValueRequest) (*GetValueResponse, error)
	// Iterate all datastores in a namespace defined by an owner
	IterateNameSpace(context.Context, *IterateNameSpaceRequest) (*IterateNameSpaceResponse, error)
	// Get a mined transaction by hash
	GetMinedTransaction(context.Context, *MinedTransactionRequest) (*MinedTransactionResponse, error)
	// Get blockheader by hash or blocknumber
	GetBlockHeader(context.Context, *BlockHeaderRequest) (*BlockHeaderResponse, error)
	// Get a raw UTXO by TxHash and index or by UTXOID
	GetUTXO(context.Context, *UTXORequest) (*UTXOResponse, error)
	// Get transaction status by hash
	GetTransactionStatus(context.Context, *TransactionStatusRequest) (*TransactionStatusResponse, error)
	// Get a pending transaction by hash
	GetPendingTransaction(context.Context, *PendingTransactionRequest) (*PendingTransactionResponse, error)
	// Get the round state object for a specified round for a specified validator
	// This allows tracing the consensus flow.
	GetRoundStateForValidator(context.Context, *RoundStateForValidatorRequest) (*RoundStateForValidatorResponse, error)
	// Get the set of validators for a specified block height
	GetValidatorSet(context.Context, *ValidatorSetRequest) (*ValidatorSetResponse, error)
	// Get the current block number
	GetBlockNumber(context.Context, *BlockNumberRequest) (*BlockNumberResponse, error)
	// Get the current ChainID of the node
	GetChainID(context.Context, *ChainIDRequest) (*ChainIDResponse, error)
	// Send a transaction to the node
	SendTransaction(context.Context, *TransactionData) (*TransactionDetails, error)
	// Get the current block number
	GetEpochNumber(context.Context, *EpochNumberRequest) (*EpochNumberResponse, error)
	// Get the current block number
	GetTxBlockNumber(context.Context, *TxBlockNumberRequest) (*TxBlockNumberResponse, error)
	GetFees(context.Context, *FeeRequest) (*FeeResponse, error)
}

// UnimplementedLocalStateServer should be embedded to have forward compatible implementations.
type UnimplementedLocalStateServer struct {
}

func (UnimplementedLocalStateServer) GetData(context.Context, *GetDataRequest) (*GetDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetData not implemented")
}
func (UnimplementedLocalStateServer) GetValueForOwner(context.Context, *GetValueRequest) (*GetValueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetValueForOwner not implemented")
}
func (UnimplementedLocalStateServer) IterateNameSpace(context.Context, *IterateNameSpaceRequest) (*IterateNameSpaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IterateNameSpace not implemented")
}
func (UnimplementedLocalStateServer) GetMinedTransaction(context.Context, *MinedTransactionRequest) (*MinedTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMinedTransaction not implemented")
}
func (UnimplementedLocalStateServer) GetBlockHeader(context.Context, *BlockHeaderRequest) (*BlockHeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockHeader not implemented")
}
func (UnimplementedLocalStateServer) GetUTXO(context.Context, *UTXORequest) (*UTXOResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUTXO not implemented")
}
func (UnimplementedLocalStateServer) GetTransactionStatus(context.Context, *TransactionStatusRequest) (*TransactionStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransactionStatus not implemented")
}
func (UnimplementedLocalStateServer) GetPendingTransaction(context.Context, *PendingTransactionRequest) (*PendingTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPendingTransaction not implemented")
}
func (UnimplementedLocalStateServer) GetRoundStateForValidator(context.Context, *RoundStateForValidatorRequest) (*RoundStateForValidatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoundStateForValidator not implemented")
}
func (UnimplementedLocalStateServer) GetValidatorSet(context.Context, *ValidatorSetRequest) (*ValidatorSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetValidatorSet not implemented")
}
func (UnimplementedLocalStateServer) GetBlockNumber(context.Context, *BlockNumberRequest) (*BlockNumberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockNumber not implemented")
}
func (UnimplementedLocalStateServer) GetChainID(context.Context, *ChainIDRequest) (*ChainIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChainID not implemented")
}
func (UnimplementedLocalStateServer) SendTransaction(context.Context, *TransactionData) (*TransactionDetails, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendTransaction not implemented")
}
func (UnimplementedLocalStateServer) GetEpochNumber(context.Context, *EpochNumberRequest) (*EpochNumberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEpochNumber not implemented")
}
func (UnimplementedLocalStateServer) GetTxBlockNumber(context.Context, *TxBlockNumberRequest) (*TxBlockNumberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTxBlockNumber not implemented")
}
func (UnimplementedLocalStateServer) GetFees(context.Context, *FeeRequest) (*FeeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFees not implemented")
}

// UnsafeLocalStateServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LocalStateServer will
// result in compilation errors.
type UnsafeLocalStateServer interface {
	mustEmbedUnimplementedLocalStateServer()
}

func RegisterLocalStateServer(s grpc.ServiceRegistrar, srv LocalStateServer) {
	s.RegisterService(&LocalState_ServiceDesc, srv)
}

func _LocalState_GetData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).GetData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/GetData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).GetData(ctx, req.(*GetDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_GetValueForOwner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetValueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).GetValueForOwner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/GetValueForOwner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).GetValueForOwner(ctx, req.(*GetValueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_IterateNameSpace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IterateNameSpaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).IterateNameSpace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/IterateNameSpace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).IterateNameSpace(ctx, req.(*IterateNameSpaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_GetMinedTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MinedTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).GetMinedTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/GetMinedTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).GetMinedTransaction(ctx, req.(*MinedTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_GetBlockHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockHeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).GetBlockHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/GetBlockHeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).GetBlockHeader(ctx, req.(*BlockHeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_GetUTXO_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UTXORequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).GetUTXO(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/GetUTXO",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).GetUTXO(ctx, req.(*UTXORequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_GetTransactionStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransactionStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).GetTransactionStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/GetTransactionStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).GetTransactionStatus(ctx, req.(*TransactionStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_GetPendingTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PendingTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).GetPendingTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/GetPendingTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).GetPendingTransaction(ctx, req.(*PendingTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_GetRoundStateForValidator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoundStateForValidatorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).GetRoundStateForValidator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/GetRoundStateForValidator",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).GetRoundStateForValidator(ctx, req.(*RoundStateForValidatorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_GetValidatorSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidatorSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).GetValidatorSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/GetValidatorSet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).GetValidatorSet(ctx, req.(*ValidatorSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_GetBlockNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockNumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).GetBlockNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/GetBlockNumber",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).GetBlockNumber(ctx, req.(*BlockNumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_GetChainID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChainIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).GetChainID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/GetChainID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).GetChainID(ctx, req.(*ChainIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_SendTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransactionData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).SendTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/SendTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).SendTransaction(ctx, req.(*TransactionData))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_GetEpochNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EpochNumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).GetEpochNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/GetEpochNumber",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).GetEpochNumber(ctx, req.(*EpochNumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_GetTxBlockNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TxBlockNumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).GetTxBlockNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/GetTxBlockNumber",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).GetTxBlockNumber(ctx, req.(*TxBlockNumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalState_GetFees_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FeeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalStateServer).GetFees(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LocalState/GetFees",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalStateServer).GetFees(ctx, req.(*FeeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LocalState_ServiceDesc is the grpc.ServiceDesc for LocalState service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LocalState_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.LocalState",
	HandlerType: (*LocalStateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetData",
			Handler:    _LocalState_GetData_Handler,
		},
		{
			MethodName: "GetValueForOwner",
			Handler:    _LocalState_GetValueForOwner_Handler,
		},
		{
			MethodName: "IterateNameSpace",
			Handler:    _LocalState_IterateNameSpace_Handler,
		},
		{
			MethodName: "GetMinedTransaction",
			Handler:    _LocalState_GetMinedTransaction_Handler,
		},
		{
			MethodName: "GetBlockHeader",
			Handler:    _LocalState_GetBlockHeader_Handler,
		},
		{
			MethodName: "GetUTXO",
			Handler:    _LocalState_GetUTXO_Handler,
		},
		{
			MethodName: "GetTransactionStatus",
			Handler:    _LocalState_GetTransactionStatus_Handler,
		},
		{
			MethodName: "GetPendingTransaction",
			Handler:    _LocalState_GetPendingTransaction_Handler,
		},
		{
			MethodName: "GetRoundStateForValidator",
			Handler:    _LocalState_GetRoundStateForValidator_Handler,
		},
		{
			MethodName: "GetValidatorSet",
			Handler:    _LocalState_GetValidatorSet_Handler,
		},
		{
			MethodName: "GetBlockNumber",
			Handler:    _LocalState_GetBlockNumber_Handler,
		},
		{
			MethodName: "GetChainID",
			Handler:    _LocalState_GetChainID_Handler,
		},
		{
			MethodName: "SendTransaction",
			Handler:    _LocalState_SendTransaction_Handler,
		},
		{
			MethodName: "GetEpochNumber",
			Handler:    _LocalState_GetEpochNumber_Handler,
		},
		{
			MethodName: "GetTxBlockNumber",
			Handler:    _LocalState_GetTxBlockNumber_Handler,
		},
		{
			MethodName: "GetFees",
			Handler:    _LocalState_GetFees_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/localstate.proto",
}
