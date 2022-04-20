package localrpc

import (
	"context"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/MadBase/MadNet/application/objs"
	"github.com/MadBase/MadNet/application/objs/uint256"
	"github.com/MadBase/MadNet/crypto"
	mncrypto "github.com/MadBase/MadNet/crypto"
	pb "github.com/MadBase/MadNet/proto"

	"github.com/MadBase/MadNet/constants"
)

func TestHandlers_HandleLocalStateGetBlockNumber(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.BlockNumberRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.BlockNumberResponse
		wantErr bool
	}{
		{name: constants.LoggerApp,
			args: args{
				ctx: ctx,
				req: &pb.BlockNumberRequest{},
			},
			want: &pb.BlockNumberResponse{BlockHeight: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srpc.HandleLocalStateGetBlockNumber(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetBlockNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateGetBlockNumber() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestHandlers_HandleLocalStateGetBlockHeader(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.BlockHeaderRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.BlockHeaderResponse
		wantErr bool
	}{
		{name: constants.LoggerApp,
			args: args{
				ctx: ctx,
				req: &pb.BlockHeaderRequest{Height: 1},
			},
			want: &pb.BlockHeaderResponse{
				BlockHeader: &pb.BlockHeader{
					BClaims: &pb.BClaims{
						ChainID:    1337,
						Height:     1,
						PrevBlock:  "41dd7c959793d4228a3c1c90d308ec31c9dd5d907c1f90afabdd38308fb5f3c8",
						StateRoot:  "2eca01388b3218b366daa6e88cb5d86b71200b428ccb06a4e3bb0065e76f1056",
						TxRoot:     "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
						HeaderRoot: fmt.Sprintf("%064d", 0),
					},
					// SigGroup: fmt.Sprintf("%0384d", 0),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srpc.HandleLocalStateGetBlockHeader(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetBlockNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.BlockHeader.BClaims, tt.want.BlockHeader.BClaims) {
				t.Errorf("HandleLocalStateGetBlockNumber() got = %v, want %v", got.BlockHeader.BClaims, tt.want.BlockHeader.BClaims)
			}
		})
	}
}

func TestHandlers_HandleLocalStateGetChainID(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.ChainIDRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.ChainIDResponse
		wantErr bool
	}{
		{name: constants.LoggerApp,
			args: args{
				ctx: ctx,
				req: &pb.ChainIDRequest{},
			},
			want: &pb.ChainIDResponse{ChainID: 1337},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srpc.HandleLocalStateGetChainID(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetChainID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateGetChainID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
func TestHandlers_HandleLocalStateGetData(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.GetDataRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GetDataResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := srpc.HandleLocalStateGetData(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateGetData() got = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
func TestHandlers_HandleLocalStateGetEpochNumber(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.EpochNumberRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.EpochNumberResponse
		wantErr bool
	}{
		{name: constants.LoggerApp,
			args: args{
				ctx: ctx,
				req: &pb.EpochNumberRequest{},
			},
			want: &pb.EpochNumberResponse{Epoch: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := srpc.HandleLocalStateGetEpochNumber(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetEpochNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateGetEpochNumber() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlers_HandleLocalStateGetFees(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.FeeRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.FeeResponse
		wantErr bool
	}{
		{name: constants.LoggerApp,
			args: args{
				ctx: ctx,
				req: &pb.FeeRequest{},
			},
			want: &pb.FeeResponse{
				MinTxFee:      fmt.Sprintf("%064d", 2),
				ValueStoreFee: fmt.Sprintf("%064d", 2),
				DataStoreFee:  fmt.Sprintf("%064d", 0),
				AtomicSwapFee: fmt.Sprintf("%064d", 1),
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srpc.HandleLocalStateGetFees(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetFees() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateGetFees() got = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
func TestHandlers_HandleLocalStateGetMinedTransaction(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.MinedTransactionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.MinedTransactionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := srpc.HandleLocalStateGetMinedTransaction(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetMinedTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateGetMinedTransaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlers_HandleLocalStateGetPendingTransaction(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.PendingTransactionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.PendingTransactionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := srpc.HandleLocalStateGetPendingTransaction(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetPendingTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateGetPendingTransaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlers_HandleLocalStateGetRoundStateForValidator(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.RoundStateForValidatorRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.RoundStateForValidatorResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := srpc.HandleLocalStateGetRoundStateForValidator(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetRoundStateForValidator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateGetRoundStateForValidator() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlers_HandleLocalStateGetTransactionStatus(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.TransactionStatusRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.TransactionStatusResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := srpc.HandleLocalStateGetTransactionStatus(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetTransactionStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateGetTransactionStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlers_HandleLocalStateGetTxBlockNumber(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.TxBlockNumberRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.TxBlockNumberResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := srpc.HandleLocalStateGetTxBlockNumber(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetTxBlockNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateGetTxBlockNumber() got = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
func TestHandlers_HandleLocalStateGetUTXO(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.UTXORequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.UTXOResponse
		wantErr bool
	}{
		{name: constants.LoggerApp,
			args: args{
				ctx: ctx,
				req: &pb.UTXORequest{
					UTXOIDs: []string{
						hex.EncodeToString(utxoIDs[0]),
					},
				},
			},
			want: &pb.UTXOResponse{
				UTXOs: []*pb.TXOut{
					&pb.TXOut{
						Utxo: &pb.TXOut_ValueStore{
							ValueStore: &pb.ValueStore{
								VSPreImage: &pb.VSPreImage{
									ChainID:  1337,
									Value:    fmt.Sprintf("%064d", 6),
									TXOutIdx: 0,
									Owner:    "0101" + hex.EncodeToString(account),
									Fee:      fmt.Sprintf("%064d", 0),
								},
								TxHash: hex.EncodeToString(txHash),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srpc.HandleLocalStateGetUTXO(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetUTXO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateGetUTXO() got = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
func TestHandlers_HandleLocalStateGetValidatorSet(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.ValidatorSetRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.ValidatorSetResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := srpc.HandleLocalStateGetValidatorSet(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetValidatorSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateGetValidatorSet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlers_HandleLocalStateGetValueForOwner(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.GetValueRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.GetValueResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := srpc.HandleLocalStateGetValueForOwner(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetValueForOwner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateGetValueForOwner() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlers_HandleLocalStateIterateNameSpace(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.IterateNameSpaceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.IterateNameSpaceResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := srpc.HandleLocalStateIterateNameSpace(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateIterateNameSpace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateIterateNameSpace() got = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
var hash []byte
var signature []byte
var newValueStore *objs.ValueStore
var vsValue *uint256.Uint256 = uint256.One()
var vsFee *uint256.Uint256 = uint256.Two()
var chainID uint32 = 1337

func TestHandlers_HandleLocalStateSendTransaction(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.TransactionData
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.TransactionDetails
		wantErr bool
	}{
		{name: constants.LoggerApp,
			args: args{
				ctx: ctx,
				req: getTransactionRequest(),
			},
			want: &pb.TransactionDetails{
				TxHash: hex.EncodeToString(hash),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srpc.HandleLocalStateSendTransaction(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateSendTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateSendTransaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTransactionRequest() *pb.TransactionData {
	pubKey, _ := signer.Pubkey()
	v := &objs.ValueStore{
		VSPreImage: &objs.VSPreImage{
			TXOutIdx: 0,
			Value:    vsValue,
			ChainID:  chainID,
			Owner: &objs.ValueStoreOwner{
				SVA:       objs.ValueStoreSVA,
				CurveSpec: constants.CurveSecp256k1,
				Account:   mncrypto.GetAccount(pubKey),
			},
			Fee: vsFee,
		},
		TxHash: make([]byte, 32),
	}
	txin := &objs.TXIn{
		TXInLinker: &objs.TXInLinker{
			TXInPreImage: &objs.TXInPreImage{
				ChainID:        chainID,
				ConsumedTxIdx:  0,
				ConsumedTxHash: txHash,
			},
			TxHash: make([]byte, 32),
		},
	}
	tx = &objs.Tx{}
	tx.Vin = []*objs.TXIn{txin}
	newValueStore = &objs.ValueStore{
		VSPreImage: &objs.VSPreImage{
			ChainID:  chainID,
			Value:    vsValue,
			TXOutIdx: 0,
			Fee:      vsFee,
			Owner: &objs.ValueStoreOwner{
				SVA:       objs.ValueStoreSVA,
				CurveSpec: constants.CurveSecp256k1,
				Account:   crypto.GetAccount(pubKey)},
		},
		TxHash: make([]byte, constants.HashLen),
	}
	newUTXO := &objs.TXOut{}
	err = newUTXO.NewValueStore(newValueStore)
	tx.Vout = append(tx.Vout, newUTXO)
	tx.Fee, _ = new(uint256.Uint256).FromUint64(3)
	err = tx.SetTxHash()
	err = v.Sign(tx.Vin[0], signer)
	hash, _ = tx.TxHash()
	signature = txin.Signature
	fmt.Printf("Hash2 %x \n", hash)
	fmt.Printf("Signature %x \n", signature)
	transactionData := *&pb.TransactionData{
		Tx: &pb.Tx{
			Vin: []*pb.TXIn{
				&pb.TXIn{
					TXInLinker: &pb.TXInLinker{
						TXInPreImage: &pb.TXInPreImage{
							ChainID:        txin.TXInLinker.TXInPreImage.ChainID,
							ConsumedTxIdx:  txin.TXInLinker.TXInPreImage.ConsumedTxIdx,
							ConsumedTxHash: hex.EncodeToString(txin.TXInLinker.TXInPreImage.ConsumedTxHash),
						},
						TxHash: hex.EncodeToString(hash),
					},
					Signature: hex.EncodeToString(signature),
				},
			},
			Vout: []*pb.TXOut{
				&pb.TXOut{
					Utxo: &pb.TXOut_ValueStore{
						ValueStore: &pb.ValueStore{
							VSPreImage: &pb.VSPreImage{
								ChainID:  newValueStore.VSPreImage.ChainID,
								TXOutIdx: newValueStore.VSPreImage.TXOutIdx,
								Value:    newValueStore.VSPreImage.Value.String(),
								Owner:    "0101" + hex.EncodeToString(account),
								Fee:      newValueStore.VSPreImage.Fee.String(),
							},
							TxHash: hex.EncodeToString(hash),
						},
					},
				},
			},
			Fee: tx.Fee.String(),
		},
	}
	fmt.Println(transactionData)
	return &transactionData
}

/* func TestHandlers_HandleLocalStateGetTransactionStatus(t *testing.T) {
	getTransactionRequest()
	fmt.Println("hash", hash)
	txin := &objs.TXIn{
		TXInLinker: &objs.TXInLinker{
			TXInPreImage: &objs.TXInPreImage{
				ChainID:        chainID,
				ConsumedTxIdx:  0,
				ConsumedTxHash: txHash,
			},
			TxHash: make([]byte, 32),
		},
	}
	type args struct {
		ctx context.Context
		req *pb.TransactionStatusRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.TransactionStatusResponse
		wantErr bool
	}{
		{name: constants.LoggerApp,
			args: args{
				ctx: ctx,
				req: &pb.TransactionStatusRequest{
					TxHash:   hex.EncodeToString(hash),
					ReturnTx: true,
				},
			},
			want: &pb.TransactionStatusResponse{IsMined: false,
				Tx: &pb.Tx{
					Vin: []*pb.TXIn{
						&pb.TXIn{
							TXInLinker: &pb.TXInLinker{
								TXInPreImage: &pb.TXInPreImage{
									ChainID:        txin.TXInLinker.TXInPreImage.ChainID,
									ConsumedTxIdx:  txin.TXInLinker.TXInPreImage.ConsumedTxIdx,
									ConsumedTxHash: hex.EncodeToString(txin.TXInLinker.TXInPreImage.ConsumedTxHash),
								},
								TxHash: hex.EncodeToString(hash),
							},
							Signature: hex.EncodeToString(signature),
						},
					},
					Vout: []*pb.TXOut{
						&pb.TXOut{
							Utxo: &pb.TXOut_ValueStore{
								ValueStore: &pb.ValueStore{
									VSPreImage: &pb.VSPreImage{
										ChainID:  chainID,
										TXOutIdx: 0,
										Value:    vsValue.String(),
										Owner:    "0101" + hex.EncodeToString(account),
										Fee:      vsFee.String(),
									},
									TxHash: hex.EncodeToString(hash),
								},
							},
						},
					},
					Fee: tx.Fee.String(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srpc.HandleLocalStateGetTransactionStatus(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLocalStateGetTransactionStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLocalStateGetTransactionStatus() got = %v, want %v", got, tt.want)
			}
		})
	}

} */

/* func TestHandlers_Init(t *testing.T) {
	consensuslessValidator("../scripts/generated/stateDBs/validator1")

	type args struct {
		database *db.Database
		app      *application.Application
		gh       *gossip.Handlers
		pubk     []byte
		safe     func() bool
		storage  dynamics.StorageGetter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: constants.LoggerApp,
			fields: fields{
				ctx:         nil,
				cancelCtx:   nil,
				database:    stateDB,
				sstore:      nil,
				AppHandler:  nil,
				GossipBus:   nil,
				Storage:     nil,
				logger:      nil,
				ethAcct:     nil,
				EthPubk:     nil,
				safeHandler: func() bool { return true },
				safecount:   0,
			},
			args: args{
				database: stateDB,
				app:      nil,
				gh:       nil,
				pubk:     nil,
				safe:     func() bool { return true },
				storage:  nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			srpc.Init(tt.args.database, tt.args.app, tt.args.gh, tt.args.pubk, tt.args.safe, tt.args.storage)
		})
	}
} */

/*
func TestHandlers_SafeMonitor(t *testing.T) {

	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			srpc.SafeMonitor()
		})
	}
}
*/

/*
func TestHandlers_notReady(t *testing.T) {

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := srpc.notReady(); (err != nil) != tt.wantErr {
				t.Errorf("notReady() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandlers_safe(t *testing.T) {

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := srpc.safe(); got != tt.want {
				t.Errorf("safe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bigIntToString(t *testing.T) {
	type args struct {
		b *big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bigIntToString(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("bigIntToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("bigIntToString() got = %v, want %v", got, tt.want)
			}
		})
	}

}
*/
