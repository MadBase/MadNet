package peering

import (
	"context"
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
	NodeAddrOne, err := transport.RandomNodeAddr()
	if err != nil {
		t.Fatal(err)
	}
	clientOne.EXPECT().NodeAddr().Return(NodeAddrOne)
	clientOne.EXPECT().NodeAddr().Return(NodeAddrOne)
	clientOne.EXPECT().NodeAddr().Return(NodeAddrOne)
	clientOne.EXPECT().CloseChan()
	activePeerStoreObj.add(clientOne)
	time.Sleep(3 * time.Second)

	clientTwo.EXPECT().NodeAddr().Return(NodeAddrOne)
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

	clientOne.EXPECT().NodeAddr().Return(NodeAddrOne)
	clientOne.EXPECT().NodeAddr().Return(NodeAddrOne)
	clientOne.EXPECT().NodeAddr().Return(NodeAddrOne)
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

	clientTwo.EXPECT().NodeAddr().Return(NodeAddrOne)
	clientTwo.EXPECT().NodeAddr().Return(NodeAddrOne)
	clientTwo.EXPECT().NodeAddr().Return(NodeAddrOne)
	clientTwo.EXPECT().CloseChan()
	clientTwo.EXPECT().NodeAddr().Return(NodeAddrOne)
	activePeerStoreObj.add(clientTwo)
	time.Sleep(3 * time.Second)

	clientTwo.EXPECT().Close()
	activePeerStoreObj.del(NodeAddrOne)
	time.Sleep(3 * time.Second)

	if len(activePeerStoreObj.store) != 0 {
		t.Fatal("not zero")
	}
	if len(activePeerStoreObj.pid) != 0 {
		t.Fatal("not zero")
	}

	// reset the close channel
	clientTwo.closeChan = make(chan struct{})

	clientTwo.EXPECT().NodeAddr().Return(NodeAddrOne)
	clientTwo.EXPECT().NodeAddr().Return(NodeAddrOne)
	clientTwo.EXPECT().NodeAddr().Return(NodeAddrOne)
	clientTwo.EXPECT().CloseChan()
	activePeerStoreObj.add(clientTwo)
	time.Sleep(3 * time.Second)

	clientTwo.EXPECT().Close()
	activePeerStoreObj.close()
	time.Sleep(3 * time.Second)
}
