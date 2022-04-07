package localrpc

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/MadBase/MadNet/blockchain/dkg/dtest"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/proto"
	"google.golang.org/grpc"
)

var address string = "127.0.0.1:8884"
var timeout time.Duration = time.Second * 100

func TestMain(m *testing.M) {
	flag.Set("test.timeout", "30m0s")
	fmt.Println(flag.Lookup("test.timeout").Value.String())
	rootPath := dtest.GetMadnetRootPath()
	scriptPath := append(rootPath, "scripts")
	scriptPath = append(scriptPath, "main.sh")
	scriptPathJoined := filepath.Join(scriptPath...)
	fmt.Println("scriptPathJoined2: ", scriptPathJoined)
	os.Setenv("SKIP_REGISTRATION", "1")
	fmt.Println("Deploying contracts. This could take minutes...")
	deploy := exec.Cmd{
		Path:   scriptPathJoined,
		Args:   []string{scriptPathJoined, "deploy"},
		Dir:    filepath.Join(rootPath...),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err := deploy.Run()
	if err != nil {
		panic("could not deploy contracts")
	}
	fmt.Println("Contracts deployed")

	fmt.Println("Starting validator")
	validator := exec.Cmd{
		Path:   scriptPathJoined,
		Args:   []string{scriptPathJoined, "validator", "1"},
		Dir:    filepath.Join(rootPath...),
		Stdout: os.Stdout,
		Stderr: nil,
	}

	err = validator.Run()
	if err != nil {
		panic("could not run validator node")
	}
	fmt.Println("Validator Started")

	//Now that we've got validator running let's test
	exitVal := m.Run()

	//After tests close validator
	validator.Process.Kill()
	os.Exit(exitVal)
}
func TestClient_Connect(t *testing.T) {
	type fields struct {
		Mutex       sync.Mutex
		closeChan   chan struct{}
		closeOnce   sync.Once
		Address     string
		TimeOut     time.Duration
		conn        *grpc.ClientConn
		client      proto.LocalStateClient
		wg          sync.WaitGroup
		isConnected bool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: constants.LoggerApp,
			fields: fields{
				Mutex:       sync.Mutex{},
				closeChan:   nil,
				closeOnce:   sync.Once{},
				Address:     address,
				TimeOut:     timeout,
				conn:        nil,
				client:      nil,
				wg:          sync.WaitGroup{},
				isConnected: false,
			},
		},
	}
	for _, tt := range tests {
		tt.args.ctx = context.Background()
		t.Run(tt.name, func(t *testing.T) {
			lrpc := &Client{
				Mutex:       tt.fields.Mutex,
				closeChan:   tt.fields.closeChan,
				closeOnce:   tt.fields.closeOnce,
				Address:     tt.fields.Address,
				TimeOut:     tt.fields.TimeOut,
				conn:        tt.fields.conn,
				client:      tt.fields.client,
				wg:          tt.fields.wg,
				isConnected: tt.fields.isConnected,
			}
			if err := lrpc.Connect(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	type fields struct {
		Mutex       sync.Mutex
		closeChan   chan struct{}
		closeOnce   sync.Once
		Address     string
		TimeOut     time.Duration
		conn        *grpc.ClientConn
		client      proto.LocalStateClient
		wg          sync.WaitGroup
		isConnected bool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: constants.LoggerApp,
			fields: fields{
				Mutex:       sync.Mutex{},
				closeChan:   nil,
				closeOnce:   sync.Once{},
				Address:     address,
				TimeOut:     timeout,
				conn:        nil,
				client:      nil,
				wg:          sync.WaitGroup{},
				isConnected: false,
			},
		},
	}
	for _, tt := range tests {
		tt.args.ctx = context.Background()
		t.Run(tt.name, func(t *testing.T) {
			lrpc := &Client{
				Mutex:       tt.fields.Mutex,
				closeChan:   tt.fields.closeChan,
				closeOnce:   tt.fields.closeOnce,
				Address:     tt.fields.Address,
				TimeOut:     tt.fields.TimeOut,
				conn:        tt.fields.conn,
				client:      tt.fields.client,
				wg:          tt.fields.wg,
				isConnected: tt.fields.isConnected,
			}
			if err := lrpc.Connect(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := lrpc.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetBlockHeader(t *testing.T) {
	type fields struct {
		Mutex       sync.Mutex
		closeChan   chan struct{}
		closeOnce   sync.Once
		Address     string
		TimeOut     time.Duration
		conn        *grpc.ClientConn
		client      proto.LocalStateClient
		wg          sync.WaitGroup
		isConnected bool
	}
	type args struct {
		ctx    context.Context
		height uint32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *objs.BlockHeader
		wantErr bool
	}{
		{name: constants.LoggerApp,
			fields: fields{
				Mutex:       sync.Mutex{},
				closeChan:   nil,
				closeOnce:   sync.Once{},
				Address:     address,
				TimeOut:     timeout,
				conn:        nil,
				client:      nil,
				wg:          sync.WaitGroup{},
				isConnected: false,
			},
		},
	}
	for _, tt := range tests {
		tt.args.ctx = context.Background()
		tt.args.height = 1
		t.Run(tt.name, func(t *testing.T) {
			lrpc := &Client{
				Mutex:       tt.fields.Mutex,
				closeChan:   tt.fields.closeChan,
				closeOnce:   tt.fields.closeOnce,
				Address:     tt.fields.Address,
				TimeOut:     tt.fields.TimeOut,
				conn:        tt.fields.conn,
				client:      tt.fields.client,
				wg:          tt.fields.wg,
				isConnected: tt.fields.isConnected,
			}
			if err := lrpc.Connect(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := lrpc.GetBlockHeader(tt.args.ctx, tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.BClaims.Height, tt.args.height) {
				t.Errorf("GetBlockHeader() got = %v, want %v", got, tt.want)
			}
		})
	}
}

/* func TestClient_GetBlockHeightForTx(t *testing.T) {
	type fields struct {
		Mutex       sync.Mutex
		closeChan   chan struct{}
		closeOnce   sync.Once
		Address     string
		TimeOut     time.Duration
		conn        *grpc.ClientConn
		client      proto.LocalStateClient
		wg          sync.WaitGroup
		isConnected bool
	}
	type args struct {
		ctx    context.Context
		txHash []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint32
		wantErr bool
	}{
		{name: constants.LoggerApp,
			fields: fields{
				Mutex:       sync.Mutex{},
				closeChan:   nil,
				closeOnce:   sync.Once{},
				Address:     address,
				TimeOut:     timeout,
				conn:        nil,
				client:      nil,
				wg:          sync.WaitGroup{},
				isConnected: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lrpc := &Client{
				Mutex:       tt.fields.Mutex,
				closeChan:   tt.fields.closeChan,
				closeOnce:   tt.fields.closeOnce,
				Address:     tt.fields.Address,
				TimeOut:     tt.fields.TimeOut,
				conn:        tt.fields.conn,
				client:      tt.fields.client,
				wg:          tt.fields.wg,
				isConnected: tt.fields.isConnected,
			}
			if err := lrpc.Connect(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := lrpc.GetBlockHeightForTx(tt.args.ctx, tt.args.txHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockHeightForTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBlockHeightForTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
/* func TestClient_GetBlockNumber(t *testing.T) {
    type fields struct {
        Mutex       sync.Mutex
        closeChan   chan struct{}
        closeOnce   sync.Once
        Address     string
        TimeOut     time.Duration
        conn        *grpc.ClientConn
        client      proto.LocalStateClient
        wg          sync.WaitGroup
        isConnected bool
    }
    type args struct {
        ctx context.Context
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    uint32
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            lrpc := &Client{
                Mutex:       tt.fields.Mutex,
                closeChan:   tt.fields.closeChan,
                closeOnce:   tt.fields.closeOnce,
                Address:     tt.fields.Address,
                TimeOut:     tt.fields.TimeOut,
                conn:        tt.fields.conn,
                client:      tt.fields.client,
                wg:          tt.fields.wg,
                isConnected: tt.fields.isConnected,
            }
            got, err := lrpc.GetBlockNumber(tt.args.ctx)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetBlockNumber() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("GetBlockNumber() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestClient_GetData(t *testing.T) {
    type fields struct {
        Mutex       sync.Mutex
        closeChan   chan struct{}
        closeOnce   sync.Once
        Address     string
        TimeOut     time.Duration
        conn        *grpc.ClientConn
        client      proto.LocalStateClient
        wg          sync.WaitGroup
        isConnected bool
    }
    type args struct {
        ctx       context.Context
        curveSpec constants.CurveSpec
        account   []byte
        index     []byte
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    []byte
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            lrpc := &Client{
                Mutex:       tt.fields.Mutex,
                closeChan:   tt.fields.closeChan,
                closeOnce:   tt.fields.closeOnce,
                Address:     tt.fields.Address,
                TimeOut:     tt.fields.TimeOut,
                conn:        tt.fields.conn,
                client:      tt.fields.client,
                wg:          tt.fields.wg,
                isConnected: tt.fields.isConnected,
            }
            got, err := lrpc.GetData(tt.args.ctx, tt.args.curveSpec, tt.args.account, tt.args.index)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetData() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetData() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestClient_GetEpochNumber(t *testing.T) {
    type fields struct {
        Mutex       sync.Mutex
        closeChan   chan struct{}
        closeOnce   sync.Once
        Address     string
        TimeOut     time.Duration
        conn        *grpc.ClientConn
        client      proto.LocalStateClient
        wg          sync.WaitGroup
        isConnected bool
    }
    type args struct {
        ctx context.Context
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    uint32
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            lrpc := &Client{
                Mutex:       tt.fields.Mutex,
                closeChan:   tt.fields.closeChan,
                closeOnce:   tt.fields.closeOnce,
                Address:     tt.fields.Address,
                TimeOut:     tt.fields.TimeOut,
                conn:        tt.fields.conn,
                client:      tt.fields.client,
                wg:          tt.fields.wg,
                isConnected: tt.fields.isConnected,
            }
            got, err := lrpc.GetEpochNumber(tt.args.ctx)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetEpochNumber() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("GetEpochNumber() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestClient_GetMinedTransaction(t *testing.T) {
    type fields struct {
        Mutex       sync.Mutex
        closeChan   chan struct{}
        closeOnce   sync.Once
        Address     string
        TimeOut     time.Duration
        conn        *grpc.ClientConn
        client      proto.LocalStateClient
        wg          sync.WaitGroup
        isConnected bool
    }
    type args struct {
        ctx    context.Context
        txHash []byte
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    *aobjs.Tx
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            lrpc := &Client{
                Mutex:       tt.fields.Mutex,
                closeChan:   tt.fields.closeChan,
                closeOnce:   tt.fields.closeOnce,
                Address:     tt.fields.Address,
                TimeOut:     tt.fields.TimeOut,
                conn:        tt.fields.conn,
                client:      tt.fields.client,
                wg:          tt.fields.wg,
                isConnected: tt.fields.isConnected,
            }
            got, err := lrpc.GetMinedTransaction(tt.args.ctx, tt.args.txHash)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetMinedTransaction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetMinedTransaction() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestClient_GetPendingTransaction(t *testing.T) {
    type fields struct {
        Mutex       sync.Mutex
        closeChan   chan struct{}
        closeOnce   sync.Once
        Address     string
        TimeOut     time.Duration
        conn        *grpc.ClientConn
        client      proto.LocalStateClient
        wg          sync.WaitGroup
        isConnected bool
    }
    type args struct {
        ctx    context.Context
        txHash []byte
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    *aobjs.Tx
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            lrpc := &Client{
                Mutex:       tt.fields.Mutex,
                closeChan:   tt.fields.closeChan,
                closeOnce:   tt.fields.closeOnce,
                Address:     tt.fields.Address,
                TimeOut:     tt.fields.TimeOut,
                conn:        tt.fields.conn,
                client:      tt.fields.client,
                wg:          tt.fields.wg,
                isConnected: tt.fields.isConnected,
            }
            got, err := lrpc.GetPendingTransaction(tt.args.ctx, tt.args.txHash)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetPendingTransaction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetPendingTransaction() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestClient_GetTxFees(t *testing.T) {
    type fields struct {
        Mutex       sync.Mutex
        closeChan   chan struct{}
        closeOnce   sync.Once
        Address     string
        TimeOut     time.Duration
        conn        *grpc.ClientConn
        client      proto.LocalStateClient
        wg          sync.WaitGroup
        isConnected bool
    }
    type args struct {
        ctx context.Context
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    []string
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            lrpc := &Client{
                Mutex:       tt.fields.Mutex,
                closeChan:   tt.fields.closeChan,
                closeOnce:   tt.fields.closeOnce,
                Address:     tt.fields.Address,
                TimeOut:     tt.fields.TimeOut,
                conn:        tt.fields.conn,
                client:      tt.fields.client,
                wg:          tt.fields.wg,
                isConnected: tt.fields.isConnected,
            }
            got, err := lrpc.GetTxFees(tt.args.ctx)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetTxFees() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetTxFees() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestClient_GetUTXO(t *testing.T) {
    type fields struct {
        Mutex       sync.Mutex
        closeChan   chan struct{}
        closeOnce   sync.Once
        Address     string
        TimeOut     time.Duration
        conn        *grpc.ClientConn
        client      proto.LocalStateClient
        wg          sync.WaitGroup
        isConnected bool
    }
    type args struct {
        ctx     context.Context
        utxoIDs [][]byte
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    objs.Vout
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            lrpc := &Client{
                Mutex:       tt.fields.Mutex,
                closeChan:   tt.fields.closeChan,
                closeOnce:   tt.fields.closeOnce,
                Address:     tt.fields.Address,
                TimeOut:     tt.fields.TimeOut,
                conn:        tt.fields.conn,
                client:      tt.fields.client,
                wg:          tt.fields.wg,
                isConnected: tt.fields.isConnected,
            }
            got, err := lrpc.GetUTXO(tt.args.ctx, tt.args.utxoIDs)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetUTXO() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetUTXO() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestClient_GetValueForOwner(t *testing.T) {
    type fields struct {
        Mutex       sync.Mutex
        closeChan   chan struct{}
        closeOnce   sync.Once
        Address     string
        TimeOut     time.Duration
        conn        *grpc.ClientConn
        client      proto.LocalStateClient
        wg          sync.WaitGroup
        isConnected bool
    }
    type args struct {
        ctx       context.Context
        curveSpec constants.CurveSpec
        account   []byte
        minValue  *uint256.Uint256
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    [][]byte
        want1   *uint256.Uint256
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            lrpc := &Client{
                Mutex:       tt.fields.Mutex,
                closeChan:   tt.fields.closeChan,
                closeOnce:   tt.fields.closeOnce,
                Address:     tt.fields.Address,
                TimeOut:     tt.fields.TimeOut,
                conn:        tt.fields.conn,
                client:      tt.fields.client,
                wg:          tt.fields.wg,
                isConnected: tt.fields.isConnected,
            }
            got, got1, err := lrpc.GetValueForOwner(tt.args.ctx, tt.args.curveSpec, tt.args.account, tt.args.minValue)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetValueForOwner() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetValueForOwner() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("GetValueForOwner() got1 = %v, want %v", got1, tt.want1)
            }
        })
    }
}

func TestClient_PaginateDataStoreUTXOByOwner(t *testing.T) {
    type fields struct {
        Mutex       sync.Mutex
        closeChan   chan struct{}
        closeOnce   sync.Once
        Address     string
        TimeOut     time.Duration
        conn        *grpc.ClientConn
        client      proto.LocalStateClient
        wg          sync.WaitGroup
        isConnected bool
    }
    type args struct {
        ctx        context.Context
        curveSpec  constants.CurveSpec
        account    []byte
        num        uint8
        startIndex []byte
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    []*aobjs.PaginationResponse
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            lrpc := &Client{
                Mutex:       tt.fields.Mutex,
                closeChan:   tt.fields.closeChan,
                closeOnce:   tt.fields.closeOnce,
                Address:     tt.fields.Address,
                TimeOut:     tt.fields.TimeOut,
                conn:        tt.fields.conn,
                client:      tt.fields.client,
                wg:          tt.fields.wg,
                isConnected: tt.fields.isConnected,
            }
            got, err := lrpc.PaginateDataStoreUTXOByOwner(tt.args.ctx, tt.args.curveSpec, tt.args.account, tt.args.num, tt.args.startIndex)
            if (err != nil) != tt.wantErr {
                t.Errorf("PaginateDataStoreUTXOByOwner() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("PaginateDataStoreUTXOByOwner() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestClient_SendTransaction(t *testing.T) {
    type fields struct {
        Mutex       sync.Mutex
        closeChan   chan struct{}
        closeOnce   sync.Once
        Address     string
        TimeOut     time.Duration
        conn        *grpc.ClientConn
        client      proto.LocalStateClient
        wg          sync.WaitGroup
        isConnected bool
    }
    type args struct {
        ctx context.Context
        tx  *aobjs.Tx
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    []byte
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            lrpc := &Client{
                Mutex:       tt.fields.Mutex,
                closeChan:   tt.fields.closeChan,
                closeOnce:   tt.fields.closeOnce,
                Address:     tt.fields.Address,
                TimeOut:     tt.fields.TimeOut,
                conn:        tt.fields.conn,
                client:      tt.fields.client,
                wg:          tt.fields.wg,
                isConnected: tt.fields.isConnected,
            }
            got, err := lrpc.SendTransaction(tt.args.ctx, tt.args.tx)
            if (err != nil) != tt.wantErr {
                t.Errorf("SendTransaction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("SendTransaction() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestClient_contextGuard(t *testing.T) {
    type fields struct {
        Mutex       sync.Mutex
        closeChan   chan struct{}
        closeOnce   sync.Once
        Address     string
        TimeOut     time.Duration
        conn        *grpc.ClientConn
        client      proto.LocalStateClient
        wg          sync.WaitGroup
        isConnected bool
    }
    type args struct {
        ctx context.Context
    }
    tests := []struct {
        name   string
        fields fields
        args   args
        want   context.Context
        want1  func()
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            lrpc := &Client{
                Mutex:       tt.fields.Mutex,
                closeChan:   tt.fields.closeChan,
                closeOnce:   tt.fields.closeOnce,
                Address:     tt.fields.Address,
                TimeOut:     tt.fields.TimeOut,
                conn:        tt.fields.conn,
                client:      tt.fields.client,
                wg:          tt.fields.wg,
                isConnected: tt.fields.isConnected,
            }
            got, got1 := lrpc.contextGuard(tt.args.ctx)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("contextGuard() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("contextGuard() got1 = %v, want %v", got1, tt.want1)
            }
        })
    }
}

func TestClient_entrancyGuard(t *testing.T) {
    type fields struct {
        Mutex       sync.Mutex
        closeChan   chan struct{}
        closeOnce   sync.Once
        Address     string
        TimeOut     time.Duration
        conn        *grpc.ClientConn
        client      proto.LocalStateClient
        wg          sync.WaitGroup
        isConnected bool
    }
    tests := []struct {
        name    string
        fields  fields
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            lrpc := &Client{
                Mutex:       tt.fields.Mutex,
                closeChan:   tt.fields.closeChan,
                closeOnce:   tt.fields.closeOnce,
                Address:     tt.fields.Address,
                TimeOut:     tt.fields.TimeOut,
                conn:        tt.fields.conn,
                client:      tt.fields.client,
                wg:          tt.fields.wg,
                isConnected: tt.fields.isConnected,
            }
            if err := lrpc.entrancyGuard(); (err != nil) != tt.wantErr {
                t.Errorf("entrancyGuard() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
*/
