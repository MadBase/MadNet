package peering

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/MadBase/MadNet/interfaces"
	pb "github.com/MadBase/MadNet/proto"
	"github.com/MadBase/MadNet/transport"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
)

type wrappedMock struct {
	*MockP2PClient
	closeChan chan struct{}
}

func (wm *wrappedMock) CloseChan() <-chan struct{} {
	wm.MockP2PClient.CloseChan()
	return wm.closeChan
}

func (wm *wrappedMock) Close() error {
	close(wm.closeChan)
	wm.MockP2PClient.Close()
	return nil
}

func (wm *wrappedMock) GetSnapShotHdrNode(context.Context, *pb.GetSnapShotHdrNodeRequest, ...grpc.CallOption) (*pb.GetSnapShotHdrNodeResponse, error) {
	return nil, nil
}

func TestActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	P2PClientOne := NewMockP2PClient(ctrl)
	P2PClientTwo := NewMockP2PClient(ctrl)

	P2PClientOneChannel := make(chan struct{})
	P2PClientTwoChannel := make(chan struct{})

	clientOne := &wrappedMock{P2PClientOne, P2PClientOneChannel}
	clientTwo := &wrappedMock{P2PClientTwo, P2PClientTwoChannel}

	activePeerStoreObj := activePeerStore{
		canClose:  true,
		store:     make(map[string]interfaces.P2PClient),
		pid:       make(map[string]uint64),
		closeChan: make(chan struct{}),
		closeOnce: sync.Once{},
	}
	randomAddr, err := transport.RandomNodeAddr()
	if err != nil {
		t.Fatal(err)
	}
	clientOne.EXPECT().NodeAddr().Return(randomAddr)
	clientOne.EXPECT().NodeAddr().Return(randomAddr)
	clientOne.EXPECT().NodeAddr().Return(randomAddr)
	clientOne.EXPECT().CloseChan()
	activePeerStoreObj.add(clientOne)
	time.Sleep(3 * time.Second)

	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	clientOne.EXPECT().CloseChan()
	clientTwo.EXPECT().Close()
	activePeerStoreObj.add(clientTwo)
	if len(activePeerStoreObj.store) != 1 {
		t.Fatal("not one")
	}
	if len(activePeerStoreObj.pid) != 1 {
		t.Fatal("not one")
	}
	time.Sleep(3 * time.Second)

	clientOne.EXPECT().NodeAddr().Return(randomAddr)
	clientOne.EXPECT().NodeAddr().Return(randomAddr)
	clientOne.EXPECT().NodeAddr().Return(randomAddr)
	close(P2PClientOneChannel)
	time.Sleep(3 * time.Second)

	if len(activePeerStoreObj.store) != 0 {
		t.Fatal("not zero")
	}
	if len(activePeerStoreObj.pid) != 0 {
		t.Fatal("not zero")
	}

	// reset the close channel
	clientTwo.closeChan = make(chan struct{})

	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	clientTwo.EXPECT().CloseChan()
	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	activePeerStoreObj.add(clientTwo)
	time.Sleep(3 * time.Second)

	clientTwo.EXPECT().Close()
	activePeerStoreObj.del(randomAddr)
	time.Sleep(3 * time.Second)

	if len(activePeerStoreObj.store) != 0 {
		t.Fatal("not zero")
	}
	if len(activePeerStoreObj.pid) != 0 {
		t.Fatal("not zero")
	}

	// reset the close channel
	clientTwo.closeChan = make(chan struct{})

	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	clientTwo.EXPECT().CloseChan()
	activePeerStoreObj.add(clientTwo)
	time.Sleep(3 * time.Second)

	clientTwo.EXPECT().Close()
	activePeerStoreObj.close()
	time.Sleep(3 * time.Second)
}

func TestActions(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	P2PClientOne := NewMockP2PClient(ctrl)
	P2PClientTwo := NewMockP2PClient(ctrl)

	P2PClientOneChannel := make(chan struct{})
	P2PClientTwoChannel := make(chan struct{})

	clientOne := &wrappedMock{P2PClientOne, P2PClientOneChannel}
	clientTwo := &wrappedMock{P2PClientTwo, P2PClientTwoChannel}

	activePeerStoreObj := activePeerStore{
		canClose:  true,
		store:     make(map[string]interfaces.P2PClient),
		pid:       make(map[string]uint64),
		closeChan: make(chan struct{}),
		closeOnce: sync.Once{},
	}
	randomAddr, err := transport.RandomNodeAddr()
	if err != nil {
		t.Fatal(err)
	}
	clientOne.EXPECT().NodeAddr().Return(randomAddr)
	clientOne.EXPECT().NodeAddr().Return(randomAddr)
	clientOne.EXPECT().NodeAddr().Return(randomAddr)
	clientOne.EXPECT().CloseChan()
	activePeerStoreObj.add(clientOne)
	time.Sleep(3 * time.Second)

	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	clientOne.EXPECT().CloseChan()
	clientTwo.EXPECT().Close()
	activePeerStoreObj.add(clientTwo)
	if len(activePeerStoreObj.store) != 1 {
		t.Fatal("not one")
	}
	if len(activePeerStoreObj.pid) != 1 {
		t.Fatal("not one")
	}
	time.Sleep(3 * time.Second)

	clientOne.EXPECT().NodeAddr().Return(randomAddr)
	clientOne.EXPECT().NodeAddr().Return(randomAddr)
	clientOne.EXPECT().NodeAddr().Return(randomAddr)
	close(P2PClientOneChannel)
	time.Sleep(3 * time.Second)

	if len(activePeerStoreObj.store) != 0 {
		t.Fatal("not zero")
	}
	if len(activePeerStoreObj.pid) != 0 {
		t.Fatal("not zero")
	}

	// reset the close channel
	clientTwo.closeChan = make(chan struct{})

	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	clientTwo.EXPECT().CloseChan()
	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	activePeerStoreObj.add(clientTwo)
	time.Sleep(3 * time.Second)

	clientTwo.EXPECT().Close()
	activePeerStoreObj.del(randomAddr)
	time.Sleep(3 * time.Second)

	if len(activePeerStoreObj.store) != 0 {
		t.Fatal("not zero")
	}
	if len(activePeerStoreObj.pid) != 0 {
		t.Fatal("not zero")
	}

	// reset the close channel
	clientTwo.closeChan = make(chan struct{})

	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	clientTwo.EXPECT().NodeAddr().Return(randomAddr)
	clientTwo.EXPECT().CloseChan()
	activePeerStoreObj.add(clientTwo)
	time.Sleep(3 * time.Second)

	clientTwo.EXPECT().Close()
	activePeerStoreObj.close()
	time.Sleep(3 * time.Second)

}

func Test_activePeerStore_add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	P2PClientOne := NewMockP2PClient(ctrl)
	P2PClientOneChannel := make(chan struct{})
	clientOne := &wrappedMock{P2PClientOne, P2PClientOneChannel}

	type fields struct {
		RWMutex   sync.RWMutex
		canClose  bool
		store     map[string]interfaces.P2PClient
		pid       map[string]uint64
		closeChan chan struct{}
		closeOnce sync.Once
	}
	type args struct {
		c interfaces.P2PClient
	}
	var tests = []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Add happy path",
			fields: fields{
				canClose:  true,
				store:     make(map[string]interfaces.P2PClient),
				pid:       make(map[string]uint64),
				closeChan: make(chan struct{}),
				closeOnce: sync.Once{},
			},
			args: args{clientOne},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			randomAddr, err := transport.RandomNodeAddr()
			if err != nil {
				t.Fatal(err)
			}
			clientOne.EXPECT().NodeAddr().Return(randomAddr)
			clientOne.EXPECT().NodeAddr().Return(randomAddr)
			clientOne.EXPECT().NodeAddr().Return(randomAddr)
			clientOne.EXPECT().CloseChan()

			ps := &activePeerStore{
				canClose:  tt.fields.canClose,
				store:     tt.fields.store,
				pid:       tt.fields.pid,
				closeChan: tt.fields.closeChan,
			}
			ps.add(tt.args.c)
		})
	}
}

func Test_activePeerStore_close(t *testing.T) {
	type fields struct {
		RWMutex   sync.RWMutex
		canClose  bool
		store     map[string]interfaces.P2PClient
		pid       map[string]uint64
		closeChan chan struct{}
		closeOnce sync.Once
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &activePeerStore{
				RWMutex:   tt.fields.RWMutex,
				canClose:  tt.fields.canClose,
				store:     tt.fields.store,
				pid:       tt.fields.pid,
				closeChan: tt.fields.closeChan,
				closeOnce: tt.fields.closeOnce,
			}
			ps.close()
		})
	}
}

func Test_activePeerStore_contains(t *testing.T) {
	type fields struct {
		RWMutex   sync.RWMutex
		canClose  bool
		store     map[string]interfaces.P2PClient
		pid       map[string]uint64
		closeChan chan struct{}
		closeOnce sync.Once
	}
	type args struct {
		c interfaces.NodeAddr
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &activePeerStore{
				RWMutex:   tt.fields.RWMutex,
				canClose:  tt.fields.canClose,
				store:     tt.fields.store,
				pid:       tt.fields.pid,
				closeChan: tt.fields.closeChan,
				closeOnce: tt.fields.closeOnce,
			}
			if got := ps.contains(tt.args.c); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_activePeerStore_del(t *testing.T) {
	type fields struct {
		RWMutex   sync.RWMutex
		canClose  bool
		store     map[string]interfaces.P2PClient
		pid       map[string]uint64
		closeChan chan struct{}
		closeOnce sync.Once
	}
	type args struct {
		c interfaces.NodeAddr
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &activePeerStore{
				RWMutex:   tt.fields.RWMutex,
				canClose:  tt.fields.canClose,
				store:     tt.fields.store,
				pid:       tt.fields.pid,
				closeChan: tt.fields.closeChan,
				closeOnce: tt.fields.closeOnce,
			}
			ps.del(tt.args.c)
		})
	}
}

func Test_activePeerStore_get(t *testing.T) {
	type fields struct {
		RWMutex   sync.RWMutex
		canClose  bool
		store     map[string]interfaces.P2PClient
		pid       map[string]uint64
		closeChan chan struct{}
		closeOnce sync.Once
	}
	type args struct {
		c interfaces.NodeAddr
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interfaces.P2PClient
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &activePeerStore{
				RWMutex:   tt.fields.RWMutex,
				canClose:  tt.fields.canClose,
				store:     tt.fields.store,
				pid:       tt.fields.pid,
				closeChan: tt.fields.closeChan,
				closeOnce: tt.fields.closeOnce,
			}
			got, got1 := ps.get(tt.args.c)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_activePeerStore_getPeers(t *testing.T) {
	type fields struct {
		RWMutex   sync.RWMutex
		canClose  bool
		store     map[string]interfaces.P2PClient
		pid       map[string]uint64
		closeChan chan struct{}
		closeOnce sync.Once
	}
	tests := []struct {
		name   string
		fields fields
		want   []interfaces.P2PClient
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &activePeerStore{
				RWMutex:   tt.fields.RWMutex,
				canClose:  tt.fields.canClose,
				store:     tt.fields.store,
				pid:       tt.fields.pid,
				closeChan: tt.fields.closeChan,
				closeOnce: tt.fields.closeOnce,
			}
			got, got1 := ps.getPeers()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPeers() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getPeers() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_activePeerStore_len(t *testing.T) {
	type fields struct {
		RWMutex   sync.RWMutex
		canClose  bool
		store     map[string]interfaces.P2PClient
		pid       map[string]uint64
		closeChan chan struct{}
		closeOnce sync.Once
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &activePeerStore{
				RWMutex:   tt.fields.RWMutex,
				canClose:  tt.fields.canClose,
				store:     tt.fields.store,
				pid:       tt.fields.pid,
				closeChan: tt.fields.closeChan,
				closeOnce: tt.fields.closeOnce,
			}
			if got := ps.len(); got != tt.want {
				t.Errorf("len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_activePeerStore_onExit(t *testing.T) {
	type fields struct {
		RWMutex   sync.RWMutex
		canClose  bool
		store     map[string]interfaces.P2PClient
		pid       map[string]uint64
		closeChan chan struct{}
		closeOnce sync.Once
	}
	type args struct {
		pid uint64
		c   interfaces.P2PClient
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &activePeerStore{
				RWMutex:   tt.fields.RWMutex,
				canClose:  tt.fields.canClose,
				store:     tt.fields.store,
				pid:       tt.fields.pid,
				closeChan: tt.fields.closeChan,
				closeOnce: tt.fields.closeOnce,
			}
			ps.onExit(tt.args.pid, tt.args.c)
		})
	}
}

func Test_activePeerStore_random(t *testing.T) {
	type fields struct {
		RWMutex   sync.RWMutex
		canClose  bool
		store     map[string]interfaces.P2PClient
		pid       map[string]uint64
		closeChan chan struct{}
		closeOnce sync.Once
	}
	tests := []struct {
		name   string
		fields fields
		want   string
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &activePeerStore{
				RWMutex:   tt.fields.RWMutex,
				canClose:  tt.fields.canClose,
				store:     tt.fields.store,
				pid:       tt.fields.pid,
				closeChan: tt.fields.closeChan,
				closeOnce: tt.fields.closeOnce,
			}
			got, got1 := ps.random()
			if got != tt.want {
				t.Errorf("random() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("random() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_activePeerStore_randomClient(t *testing.T) {
	type fields struct {
		RWMutex   sync.RWMutex
		canClose  bool
		store     map[string]interfaces.P2PClient
		pid       map[string]uint64
		closeChan chan struct{}
		closeOnce sync.Once
	}
	tests := []struct {
		name   string
		fields fields
		want   interfaces.P2PClient
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &activePeerStore{
				RWMutex:   tt.fields.RWMutex,
				canClose:  tt.fields.canClose,
				store:     tt.fields.store,
				pid:       tt.fields.pid,
				closeChan: tt.fields.closeChan,
				closeOnce: tt.fields.closeOnce,
			}
			got, got1 := ps.randomClient()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("randomClient() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("randomClient() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
