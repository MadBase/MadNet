package lstate

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/MadBase/MadNet/consensus/appmock"
	"github.com/MadBase/MadNet/consensus/db"
	objs "github.com/MadBase/MadNet/consensus/objs"
	"github.com/dgraph-io/badger/v2"
	"github.com/golang/mock/gomock"

	mcrypto "github.com/MadBase/MadNet/crypto"

	bn256 "github.com/MadBase/MadNet/crypto/bn256/cloudflare"
)

func TestMockeddb(t *testing.T) {
	ctr := gomock.NewController(t)

	db := db.NewMockDatabaseIface(ctr)

	if db == nil {
		t.Fatal("db should not be nil")
	}
}

func getBdb(t *testing.T) *badger.DB {
	// Open the DB.
	dir, err := ioutil.TempDir("", "badger-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
	}()
	opts := badger.DefaultOptions(dir)
	bdb, err := badger.Open(opts)
	if err != nil {
		t.Fatal(err)
	}

	return bdb
}

func TestFhFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{2, 3}, TxRoot: txRoot, StateRoot: []byte{2},
		HeaderRoot: []byte{3}}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	msg, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}
	grpSig, _, _, _ := getGroupSig(msg)

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims, SigGroup: grpSig}}
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetCommittedBlockHeader(gomock.Any(), gomock.Any()).Return(nil)
		mdb.EXPECT().SetBroadcastBlockHeader(gomock.Any(), gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		roundStates.maxHR = roundState

		booleanValue, err := stateHandler.fhFunc(roundStates)
		if err != nil {
			fmt.Println("err is", err)
		}

		if booleanValue != true {
			t.Fatal("output value for fhfunc is not correct")
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}

	updateFunc = func(txn *badger.Txn) error {

		otherrClaims := &objs.RClaims{Height: 2, Round: 1}
		rCert := &objs.RCert{RClaims: otherrClaims}
		otherbClaims = &objs.BClaims{ChainID: 42, Height: 2, TxCount: 53, PrevBlock: []byte{2, 3}, TxRoot: txRoot, StateRoot: []byte{2},
			HeaderRoot: []byte{3}}
		pClaims = &objs.PClaims{RCert: rCert, BClaims: otherbClaims}

		bhsh, err := pClaims.BClaims.BlockHash()
		if err != nil {
			log.Fatal(err)
		}

		rClaims = &objs.RClaims{ChainID: 42, Height: 2, Round: 1, PrevBlock: bhsh}
		roundState = &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
			RClaims: rClaims}}

		validValue = &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
		ownValState = &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetCommittedBlockHeader(gomock.Any(), gomock.Any()).Return(nil)
		mdb.EXPECT().SetBroadcastBlockHeader(gomock.Any(), gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		msg, err := pClaims.BClaims.BlockHash()
		if err != nil {
			return err
		}
		grpSig, _, _, _ := getGroupSig(msg)

		roundStates.maxHR = roundState

		roundStates.maxHR.RCert.SigGroup = grpSig

		booleanValue, err := stateHandler.fhFunc(roundStates)
		if err != nil {
			fmt.Println("err is", err)
		}

		// should we be checking the boolean value for the handlers

		if booleanValue != true || roundStates.OwnState.MaxBHSeen.BClaims.Height != 2 {
			t.Fatal("max bh seen height should have a value of 2")
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRoundJumpUpdateValidValueFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, TxRoot: txRoot, StateRoot: []byte{173, 233, 94, 109, 13, 42, 99,
		22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147},
		HeaderRoot: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21,
			229, 112, 18, 48, 147}}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims, GroupKey: []byte{3, 5, 4}}
	otherValidValue := &objs.Proposal{Signature: []byte{3, 3, 8}, PClaims: pClaims, GroupKey: []byte{3, 5, 4}}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: otherValidValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	// shouldn't the validators have different vaddrs ? seems to not allow precommit updates if they do
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: []byte{1}}, {VAddr: []byte{2}, GroupShare: []byte{2}},
			{VAddr: []byte{3}, GroupShare: []byte{3}}, {VAddr: []byte{1}, GroupShare: []byte{4}}, {VAddr: []byte{2}, GroupShare: []byte{5}}}}

	updateFunc = func(txn *badger.Txn) error {

		h := uint32(1)
		mdb.EXPECT().GetHeaderTrieRoot(txn, h).Return([]byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120,
			194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, nil)

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		roundStates.PeerStateMap[string([]byte{1})].PreCommit = &objs.PreCommit{Proposal: validValue, Voter: []byte{1, 2}}
		roundStates.PeerStateMap[string([]byte{2})].PreCommit = &objs.PreCommit{Proposal: validValue, Voter: []byte{1, 3}}
		roundStates.PeerStateMap[string([]byte{3})].PreCommit = &objs.PreCommit{Proposal: validValue, Voter: []byte{1, 4}}

		msg, err := pClaims.BClaims.BlockHash()
		if err != nil {
			return err
		}
		grpSig, _, _, _ := getGroupSig(msg)

		rs := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 2, 3}, RClaims: rClaims,
			SigGroup: grpSig}}

		pcl, _, err := roundStates.GetCurrentPreCommits()
		if err != nil {
			log.Fatal(err)
		}

		booleanValue, err := stateHandler.roundJumpUpdateValidValueFunc(roundStates, rs.RCert, pcl)
		if err != nil {
			fmt.Println("err is", err)
		}

		if booleanValue != true {
			t.Fatal("output value for r cert func is not correct")
		}

		expectedSig := []byte{3, 3, 3}
		for i := 0; i < len(expectedSig); i++ {
			if expectedSig[i] != roundStates.OwnValidatingState.ValidValue.Signature[i] {
				t.Fatal("incorrect value for signature")
			}
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRoundJumpSetRCertFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, TxRoot: txRoot, StateRoot: []byte{173, 233, 94, 109, 13, 42, 99,
		22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147},
		HeaderRoot: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21,
			229, 112, 18, 48, 147}}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims, GroupKey: []byte{3, 5, 4}}
	otherValidValue := &objs.Proposal{Signature: []byte{3, 3, 8}, PClaims: pClaims, GroupKey: []byte{3, 5, 4}}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: otherValidValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	// shouldn't the validators have different vaddrs ? seems to not allow precommit updates if they do
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: []byte{1}}, {VAddr: []byte{2}, GroupShare: []byte{2}},
			{VAddr: []byte{3}, GroupShare: []byte{3}}, {VAddr: []byte{1}, GroupShare: []byte{4}}, {VAddr: []byte{2}, GroupShare: []byte{5}}}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		roundStates.PeerStateMap[string([]byte{1})].PreCommit = &objs.PreCommit{Proposal: validValue, Voter: []byte{1, 2}}
		roundStates.PeerStateMap[string([]byte{2})].PreCommit = &objs.PreCommit{Proposal: validValue, Voter: []byte{1, 3}}
		roundStates.PeerStateMap[string([]byte{3})].PreCommit = &objs.PreCommit{Proposal: validValue, Voter: []byte{1, 4}}

		msg, err := pClaims.BClaims.BlockHash()
		if err != nil {
			return err
		}
		grpSig, _, _, _ := getGroupSig(msg)

		rs := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 2, 3}, RClaims: rClaims,
			SigGroup: grpSig}}

		pcl, _, err := roundStates.GetCurrentPreCommits()
		if err != nil {
			log.Fatal(err)
		}

		booleanValue, err := stateHandler.roundJumpSetRCertFunc(roundStates, rs.RCert, pcl)
		if err != nil {
			fmt.Println("err is", err)
		}

		if booleanValue != true {
			t.Fatal("output value for r cert func is not correct")
		}

		expectedSig := []byte{3, 3, 8}
		for i := 0; i < len(expectedSig); i++ {
			if expectedSig[i] != roundStates.OwnValidatingState.ValidValue.Signature[i] {
				t.Fatal("incorrect value for signature")
			}
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoNrsCastNextHeightHandlerFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, ChainID: 42}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{2, 3}, TxRoot: txRoot, StateRoot: []byte{2},
		HeaderRoot: []byte{3}}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	msg, err := pClaims.BClaims.BlockHash()
	if err != nil {
		t.Fatal(err)
	}
	_, gs, gpk, _ := getGroupSig(msg)

	stateHandler.bnSigner.SetGroupPubk(gpk)

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{
		GroupKey:          gpk,
		ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: gs[0]}, {VAddr: []byte{2}, GroupShare: gs[1]},
			{VAddr: []byte{3}, GroupShare: gs[2]}, {VAddr: []byte{1}, GroupShare: gs[3]}}}

	updateFunc = func(txn *badger.Txn) error {

		h := uint32(1)
		mdb.EXPECT().GetHeaderTrieRoot(txn, h).Return([]byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120,
			194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, nil)

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastNextHeight(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		msg, err := pClaims.BClaims.BlockHash()
		if err != nil {
			return err
		}

		txr := []byte{197, 210, 70, 1, 134, 247, 35, 60, 146, 126, 125, 178, 220, 199, 3, 192, 229, 0, 182, 83, 202, 130, 39, 59, 123, 250, 216, 4, 93, 133, 164, 112}
		bs := []byte{137, 158, 164, 26, 219, 131, 151, 198, 183, 30, 184, 92, 126, 36, 84, 26, 33, 2, 95, 173, 235, 114, 104, 193, 225, 73, 193, 104, 229, 123, 61, 37, 111, 25, 109, 229, 148, 232, 96, 32, 23, 29, 116, 208, 88, 123, 82, 228, 215, 71, 195, 127, 104, 209, 148, 7, 41, 209, 77, 220, 127, 177, 247, 214, 0}
		roundStates.PeerStateMap[string([]byte{1})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{2})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{3})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}

		roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].RCert = &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}

		pcl, _, err := roundStates.GetCurrentPreCommits()
		if err != nil {
			log.Fatal(err)
		}

		booleanValue, err := stateHandler.doNrsCastNextHeightHandlerFunc(roundStates, roundState.RCert, pcl)
		if err != nil {
			t.Fatal("err is", err)
		}

		if booleanValue != true {
			t.Fatal("value of the output from do nr func seems to not be correct")
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoNrsCastNextRoundHandlerFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, PrevBlock: txRoot, ChainID: 42}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: txRoot, TxRoot: txRoot, StateRoot: txRoot,
		HeaderRoot: txRoot}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 2, Round: 2, PrevBlock: bhsh}

	msg, err := rClaims.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	grpSig, gs, gpk, ss := getGroupSig(msg)

	stateHandler.bnSigner.SetGroupPubk(gpk)

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: gs[0]}, {VAddr: []byte{2}, GroupShare: gs[1]},
			{VAddr: []byte{3}, GroupShare: gs[2]}, {VAddr: []byte{5, 5, 5}, GroupShare: gs[3]}}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		roundState1 := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
			RClaims: rClaims}}
		roundState1.NextRound = &objs.NextRound{Signature: ss[0],
			NRClaims: &objs.NRClaims{RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Round: 2, Height: 2, PrevBlock: bhsh}},
				RClaims: &objs.RClaims{ChainID: 42, Round: 2, Height: 2, PrevBlock: bhsh}, SigShare: ss[0]}}
		roundState2 := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
			RClaims: rClaims}}
		roundState2.NextRound = &objs.NextRound{Signature: ss[1],
			NRClaims: &objs.NRClaims{RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Round: 2, Height: 2, PrevBlock: bhsh}},
				RClaims: &objs.RClaims{ChainID: 42, Round: 2, Height: 2, PrevBlock: bhsh}, SigShare: ss[1]}}
		roundState3 := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
			RClaims: rClaims}}
		roundState3.NextRound = &objs.NextRound{Signature: ss[2],
			NRClaims: &objs.NRClaims{RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Round: 2, Height: 2, PrevBlock: bhsh}},
				RClaims: &objs.RClaims{ChainID: 42, Round: 2, Height: 2, PrevBlock: bhsh}, SigShare: ss[2]}}

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState1, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState2, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState3, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{5, 5, 5}).Return(roundState, nil)

		mdb.EXPECT().SetBroadcastRCert(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		roundStates.PeerStateMap[string([]byte{1})].NextRound = roundState1.NextRound
		roundStates.PeerStateMap[string([]byte{2})].NextRound = roundState2.NextRound
		roundStates.PeerStateMap[string([]byte{3})].NextRound = roundState3.NextRound

		roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].RCert = &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 2, Round: 2, PrevBlock: bhsh}}
		roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].NextRound = &objs.NextRound{Signature: ss[3],
			NRClaims: &objs.NRClaims{RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Round: 2, Height: 2, PrevBlock: bhsh}},
				RClaims: &objs.RClaims{ChainID: 42, Round: 1, Height: 2, PrevBlock: bhsh}, SigShare: ss[3]}}

		rs := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 2, 3}, RClaims: rClaims,
			SigGroup: grpSig}}

		_, nrl, err := roundStates.GetCurrentNext()
		if err != nil {
			log.Fatal(err)
		}

		booleanValue, err := stateHandler.doNrsCastNextRoundHandlerFunc(roundStates, rs.RCert, nrl)
		if err != nil {
			fmt.Println("err is", err)
		}

		if booleanValue != true {
			t.Fatal("value of the output from do nr func seems to not be correct")
		}

		if roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].RCert.RClaims.Round != 2 {
			t.Fatal("value for the round is probably not correct")
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCastBlockHeaderHandlerFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: txRoot, TxRoot: txRoot, StateRoot: txRoot,
		HeaderRoot: txRoot}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	msg, err := pClaims.BClaims.BlockHash()
	if err != nil {
		t.Fatal(err)
	}
	_, gs, gpk, sigs := getGroupSig(msg)

	stateHandler.bnSigner.SetGroupPubk(gpk)

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	roundState1 := &objs.RoundState{VAddr: []byte{1}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	roundState2 := &objs.RoundState{VAddr: []byte{2}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	roundState3 := &objs.RoundState{VAddr: []byte{3}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{
		GroupKey:          gpk,
		ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: gs[0]}, {VAddr: []byte{2}, GroupShare: gs[1]},
			{VAddr: []byte{3}, GroupShare: gs[2]}, {VAddr: []byte{1}, GroupShare: gs[3]}}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, gomock.Any()).Return(roundState1, nil)
		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, gomock.Any()).Return(roundState3, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, gomock.Any()).Return(roundState1, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, gomock.Any()).Return(roundState2, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, gomock.Any()).Return(roundState3, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, gomock.Any()).Return(roundState2, nil)

		mdb.EXPECT().SetCommittedBlockHeader(gomock.Any(), gomock.Any()).Return(nil)
		mdb.EXPECT().SetBroadcastBlockHeader(gomock.Any(), gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		nhs := objs.NextHeightList{&objs.NextHeight{NHClaims: &objs.NHClaims{SigShare: sigs[0], Proposal: &objs.Proposal{PClaims: pClaims, GroupKey: gpk}}},
			&objs.NextHeight{NHClaims: &objs.NHClaims{SigShare: sigs[1], Proposal: &objs.Proposal{PClaims: pClaims, GroupKey: gpk}}},
			&objs.NextHeight{NHClaims: &objs.NHClaims{SigShare: sigs[2], Proposal: &objs.Proposal{PClaims: pClaims, GroupKey: gpk}}},
			&objs.NextHeight{NHClaims: &objs.NHClaims{SigShare: sigs[3], Proposal: &objs.Proposal{PClaims: pClaims, GroupKey: gpk}}}}

		_, err = stateHandler.castBlockHeaderHandlerFunc(roundStates, nhs)
		if err != nil {
			fmt.Println("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPNCastNextHeightFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{2, 3}, TxRoot: txRoot, StateRoot: []byte{2},
		HeaderRoot: []byte{3}}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	msg, err := pClaims.BClaims.BlockHash()
	if err != nil {
		t.Fatal(err)
	}
	_, gs, gpk, _ := getGroupSig(msg)

	stateHandler.bnSigner.SetGroupPubk(gpk)

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{
		GroupKey:          gpk,
		ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: gs[0]}, {VAddr: []byte{2}, GroupShare: gs[1]},
			{VAddr: []byte{3}, GroupShare: gs[2]}, {VAddr: []byte{1}, GroupShare: gs[3]}}}

	updateFunc = func(txn *badger.Txn) error {

		h := uint32(1)
		mdb.EXPECT().GetHeaderTrieRoot(txn, h).Return([]byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120,
			194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, nil)

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastNextHeight(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		msg, err := pClaims.BClaims.BlockHash()
		if err != nil {
			return err
		}

		txr := []byte{197, 210, 70, 1, 134, 247, 35, 60, 146, 126, 125, 178, 220, 199, 3, 192, 229, 0, 182, 83, 202, 130, 39, 59, 123, 250, 216, 4, 93, 133, 164, 112}
		bs := []byte{137, 158, 164, 26, 219, 131, 151, 198, 183, 30, 184, 92, 126, 36, 84, 26, 33, 2, 95, 173, 235, 114, 104, 193, 225, 73, 193, 104, 229, 123, 61, 37, 111, 25, 109, 229, 148, 232, 96, 32, 23, 29, 116, 208, 88, 123, 82, 228, 215, 71, 195, 127, 104, 209, 148, 7, 41, 209, 77, 220, 127, 177, 247, 214, 0}
		roundStates.PeerStateMap[string([]byte{1})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{2})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{3})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}

		roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].RCert = &objs.RCert{RClaims: &objs.RClaims{Height: 1, Round: 1, PrevBlock: bhsh}}

		pcl, _, err := roundStates.GetCurrentPreCommits()
		if err != nil {
			log.Fatal(err)
		}

		err = stateHandler.dPNCastNextHeightFunc(roundStates, pcl)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPNCastNextRoundFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, ChainID: 42}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{2, 3}, TxRoot: txRoot, StateRoot: []byte{2},
		HeaderRoot: []byte{3}}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	msg, err := pClaims.BClaims.BlockHash()
	if err != nil {
		t.Fatal(err)
	}
	_, gs, gpk, _ := getGroupSig(msg)

	stateHandler.bnSigner.SetGroupPubk(gpk)

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{
		GroupKey:          gpk,
		ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: gs[0]}, {VAddr: []byte{2}, GroupShare: gs[1]},
			{VAddr: []byte{3}, GroupShare: gs[2]}, {VAddr: []byte{1}, GroupShare: gs[3]}}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastNextRound(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		msg, err := pClaims.BClaims.BlockHash()
		if err != nil {
			return err
		}

		txr := []byte{197, 210, 70, 1, 134, 247, 35, 60, 146, 126, 125, 178, 220, 199, 3, 192, 229, 0, 182, 83, 202, 130, 39, 59, 123, 250, 216, 4, 93, 133, 164, 112}
		bs := []byte{137, 158, 164, 26, 219, 131, 151, 198, 183, 30, 184, 92, 126, 36, 84, 26, 33, 2, 95, 173, 235, 114, 104, 193, 225, 73, 193, 104, 229, 123, 61, 37, 111, 25, 109, 229, 148, 232, 96, 32, 23, 29, 116, 208, 88, 123, 82, 228, 215, 71, 195, 127, 104, 209, 148, 7, 41, 209, 77, 220, 127, 177, 247, 214, 0}
		roundStates.PeerStateMap[string([]byte{1})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{2})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{3})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}

		roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].RCert = &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}

		err = stateHandler.dPNCastNextRoundFunc(roundStates)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPCSCastNHFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, ChainID: 42}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{2, 3}, TxRoot: txRoot, StateRoot: []byte{2},
		HeaderRoot: []byte{3}}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	msg, err := pClaims.BClaims.BlockHash()
	if err != nil {
		t.Fatal(err)
	}
	_, gs, gpk, _ := getGroupSig(msg)

	stateHandler.bnSigner.SetGroupPubk(gpk)

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{
		GroupKey:          gpk,
		ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: gs[0]}, {VAddr: []byte{2}, GroupShare: gs[1]},
			{VAddr: []byte{3}, GroupShare: gs[2]}, {VAddr: []byte{1}, GroupShare: gs[3]}}}

	updateFunc = func(txn *badger.Txn) error {

		h := uint32(1)
		mdb.EXPECT().GetHeaderTrieRoot(txn, h).Return([]byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120,
			194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, nil)

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastNextHeight(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		msg, err := pClaims.BClaims.BlockHash()
		if err != nil {
			return err
		}

		txr := []byte{197, 210, 70, 1, 134, 247, 35, 60, 146, 126, 125, 178, 220, 199, 3, 192, 229, 0, 182, 83, 202, 130, 39, 59, 123, 250, 216, 4, 93, 133, 164, 112}
		bs := []byte{137, 158, 164, 26, 219, 131, 151, 198, 183, 30, 184, 92, 126, 36, 84, 26, 33, 2, 95, 173, 235, 114, 104, 193, 225, 73, 193, 104, 229, 123, 61, 37, 111, 25, 109, 229, 148, 232, 96, 32, 23, 29, 116, 208, 88, 123, 82, 228, 215, 71, 195, 127, 104, 209, 148, 7, 41, 209, 77, 220, 127, 177, 247, 214, 0}
		roundStates.PeerStateMap[string([]byte{1})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{2})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{3})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}

		roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].RCert = &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}

		pcl, _, err := roundStates.GetCurrentPreCommits()
		if err != nil {
			log.Fatal(err)
		}

		err = stateHandler.dPCSCastNHFunc(roundStates, pcl)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPCNSCastNRFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, ChainID: 42}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{2, 3}, TxRoot: txRoot, StateRoot: []byte{2},
		HeaderRoot: []byte{3}}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	msg, err := pClaims.BClaims.BlockHash()
	if err != nil {
		t.Fatal(err)
	}
	_, gs, gpk, _ := getGroupSig(msg)

	stateHandler.bnSigner.SetGroupPubk(gpk)

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{
		GroupKey:          gpk,
		ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: gs[0]}, {VAddr: []byte{2}, GroupShare: gs[1]},
			{VAddr: []byte{3}, GroupShare: gs[2]}, {VAddr: []byte{1}, GroupShare: gs[3]}}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastNextRound(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		msg, err := pClaims.BClaims.BlockHash()
		if err != nil {
			return err
		}

		txr := []byte{197, 210, 70, 1, 134, 247, 35, 60, 146, 126, 125, 178, 220, 199, 3, 192, 229, 0, 182, 83, 202, 130, 39, 59, 123, 250, 216, 4, 93, 133, 164, 112}
		bs := []byte{137, 158, 164, 26, 219, 131, 151, 198, 183, 30, 184, 92, 126, 36, 84, 26, 33, 2, 95, 173, 235, 114, 104, 193, 225, 73, 193, 104, 229, 123, 61, 37, 111, 25, 109, 229, 148, 232, 96, 32, 23, 29, 116, 208, 88, 123, 82, 228, 215, 71, 195, 127, 104, 209, 148, 7, 41, 209, 77, 220, 127, 177, 247, 214, 0}
		roundStates.PeerStateMap[string([]byte{1})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{2})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{3})].PreCommit = &objs.PreCommit{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: msg, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}

		roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].RCert = &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}

		err = stateHandler.dPCNSCastNRFunc(roundStates)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPPCCastPCFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, TxRoot: txRoot, StateRoot: []byte{2},
		HeaderRoot: []byte{3}}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: []byte{2}}, {VAddr: []byte{2}, GroupShare: []byte{1}},
			{VAddr: []byte{3}, GroupShare: []byte{3}}, {VAddr: []byte{1}, GroupShare: []byte{4}}}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastPreCommit(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		msg, err := pClaims.BClaims.BlockHash()
		if err != nil {
			return err
		}

		txr := []byte{197, 210, 70, 1, 134, 247, 35, 60, 146, 126, 125, 178, 220, 199, 3, 192, 229, 0, 182, 83, 202, 130, 39, 59, 123, 250, 216, 4, 93, 133, 164, 112}
		bs := []byte{137, 158, 164, 26, 219, 131, 151, 198, 183, 30, 184, 92, 126, 36, 84, 26, 33, 2, 95, 173, 235, 114, 104, 193, 225, 73, 193, 104, 229,
			123, 61, 37, 111, 25, 109, 229, 148, 232, 96, 32, 23, 29, 116, 208, 88, 123, 82, 228, 215, 71, 195, 127, 104, 209, 148, 7, 41, 209, 77, 220, 127,
			177, 247, 214, 0}
		roundStates.PeerStateMap[string([]byte{1})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{2})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{3})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}

		roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].RCert = &objs.RCert{RClaims: &objs.RClaims{Height: 1, Round: 1, PrevBlock: bhsh}}

		pvl, _, err := roundStates.GetCurrentPreVotes()
		if err != nil {
			log.Fatal(err)
		}

		err = stateHandler.dPPCCastPCFunc(roundStates, pvl)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPPCUpdateVVFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, ChainID: 42}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{2, 3}, TxRoot: txRoot, StateRoot: []byte{2},
		HeaderRoot: []byte{3}}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	msg, err := pClaims.BClaims.BlockHash()
	if err != nil {
		t.Fatal(err)
	}
	_, gs, gpk, _ := getGroupSig(msg)

	stateHandler.bnSigner.SetGroupPubk(gpk)

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{
		GroupKey:          gpk,
		ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: gs[0]}, {VAddr: []byte{2}, GroupShare: gs[1]},
			{VAddr: []byte{3}, GroupShare: gs[2]}, {VAddr: []byte{1}, GroupShare: gs[3]}}}

	updateFunc = func(txn *badger.Txn) error {

		h := uint32(1)
		mdb.EXPECT().GetHeaderTrieRoot(txn, h).Return([]byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120,
			194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, nil)

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastPreCommitNil(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		msg, err := pClaims.BClaims.BlockHash()
		if err != nil {
			return err
		}

		txr := []byte{197, 210, 70, 1, 134, 247, 35, 60, 146, 126, 125, 178, 220, 199, 3, 192, 229, 0, 182, 83, 202, 130, 39, 59, 123, 250, 216, 4, 93, 133, 164, 112}
		bs := []byte{137, 158, 164, 26, 219, 131, 151, 198, 183, 30, 184, 92, 126, 36, 84, 26, 33, 2, 95, 173, 235, 114, 104, 193, 225, 73, 193, 104, 229, 123, 61, 37, 111, 25, 109, 229, 148, 232, 96, 32, 23, 29, 116, 208, 88, 123, 82, 228, 215, 71, 195, 127, 104, 209, 148, 7, 41, 209, 77, 220, 127, 177, 247, 214, 0}
		roundStates.PeerStateMap[string([]byte{1})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{2})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{3})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}

		roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].RCert = &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}

		pvl, pvnl, err := roundStates.GetCurrentPreVotes()
		if err != nil {
			log.Fatal(err)
		}

		err = stateHandler.dPPCUpdateVVFunc(roundStates, pvl, pvnl)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPPCNotDBRFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, TxRoot: txRoot, StateRoot: []byte{2},
		HeaderRoot: []byte{3}}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: []byte{2}}, {VAddr: []byte{2}, GroupShare: []byte{1}},
			{VAddr: []byte{3}, GroupShare: []byte{3}}, {VAddr: []byte{1}, GroupShare: []byte{4}}}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastPreCommitNil(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		msg, err := pClaims.BClaims.BlockHash()
		if err != nil {
			return err
		}

		txr := []byte{197, 210, 70, 1, 134, 247, 35, 60, 146, 126, 125, 178, 220, 199, 3, 192, 229, 0, 182, 83, 202, 130, 39, 59, 123, 250, 216, 4, 93, 133, 164, 112}
		bs := []byte{137, 158, 164, 26, 219, 131, 151, 198, 183, 30, 184, 92, 126, 36, 84, 26, 33, 2, 95, 173, 235, 114, 104, 193, 225, 73, 193, 104, 229, 123, 61, 37, 111, 25, 109, 229, 148, 232, 96, 32, 23, 29, 116, 208, 88, 123, 82, 228, 215, 71, 195, 127, 104, 209, 148, 7, 41, 209, 77, 220, 127, 177, 247, 214, 0}
		roundStates.PeerStateMap[string([]byte{1})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{2})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{3})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}

		roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].RCert = &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}

		err = stateHandler.dPPCNotDBRFunc(roundStates)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPVSCastPCFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, TxRoot: txRoot, StateRoot: []byte{2},
		HeaderRoot: []byte{3}}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: []byte{2}}, {VAddr: []byte{2}, GroupShare: []byte{1}},
			{VAddr: []byte{3}, GroupShare: []byte{3}}, {VAddr: []byte{1}, GroupShare: []byte{4}}}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastPreCommit(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		msg, err := pClaims.BClaims.BlockHash()
		if err != nil {
			return err
		}

		txr := []byte{197, 210, 70, 1, 134, 247, 35, 60, 146, 126, 125, 178, 220, 199, 3, 192, 229, 0, 182, 83, 202, 130, 39, 59, 123, 250, 216, 4, 93, 133, 164, 112}
		bs := []byte{137, 158, 164, 26, 219, 131, 151, 198, 183, 30, 184, 92, 126, 36, 84, 26, 33, 2, 95, 173, 235, 114, 104, 193, 225, 73, 193, 104, 229, 123, 61, 37, 111, 25, 109, 229, 148, 232, 96, 32, 23, 29, 116, 208, 88, 123, 82, 228, 215, 71, 195, 127, 104, 209, 148, 7, 41, 209, 77, 220, 127, 177, 247, 214, 0}
		roundStates.PeerStateMap[string([]byte{1})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{2})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{3})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}

		roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].RCert = &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}

		pvl, _, err := roundStates.GetCurrentPreVotes()
		if err != nil {
			log.Fatal(err)
		}

		err = stateHandler.dPVSCastPCFunc(roundStates, pvl)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPVNSUpdateVVFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, ChainID: 42}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{2, 3}, TxRoot: txRoot, StateRoot: []byte{2},
		HeaderRoot: []byte{3}}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	msg, err := pClaims.BClaims.BlockHash()
	if err != nil {
		t.Fatal(err)
	}
	_, gs, gpk, _ := getGroupSig(msg)

	stateHandler.bnSigner.SetGroupPubk(gpk)

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{
		GroupKey:          gpk,
		ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: gs[0]}, {VAddr: []byte{2}, GroupShare: gs[1]},
			{VAddr: []byte{3}, GroupShare: gs[2]}, {VAddr: []byte{1}, GroupShare: gs[3]}}}

	updateFunc = func(txn *badger.Txn) error {

		h := uint32(1)
		mdb.EXPECT().GetHeaderTrieRoot(txn, h).Return([]byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120,
			194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, nil)

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		msg, err := pClaims.BClaims.BlockHash()
		if err != nil {
			return err
		}

		txr := []byte{197, 210, 70, 1, 134, 247, 35, 60, 146, 126, 125, 178, 220, 199, 3, 192, 229, 0, 182, 83, 202, 130, 39, 59, 123, 250, 216, 4, 93, 133, 164, 112}
		bs := []byte{137, 158, 164, 26, 219, 131, 151, 198, 183, 30, 184, 92, 126, 36, 84, 26, 33, 2, 95, 173, 235, 114, 104, 193, 225, 73, 193, 104, 229, 123, 61, 37, 111, 25, 109, 229, 148, 232, 96, 32, 23, 29, 116, 208, 88, 123, 82, 228, 215, 71, 195, 127, 104, 209, 148, 7, 41, 209, 77, 220, 127, 177, 247, 214, 0}
		roundStates.PeerStateMap[string([]byte{1})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{2})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{3})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}

		roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].RCert = &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}

		pvl, _, err := roundStates.GetCurrentPreVotes()
		if err != nil {
			log.Fatal(err)
		}

		err = stateHandler.dPVNSUpdateVVFunc(roundStates, pvl)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPVNSCastPCNFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, TxRoot: txRoot, StateRoot: []byte{2},
		HeaderRoot: []byte{3}}
	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: []byte{2}}, {VAddr: []byte{2}, GroupShare: []byte{1}},
			{VAddr: []byte{3}, GroupShare: []byte{3}}, {VAddr: []byte{1}, GroupShare: []byte{4}}}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastPreCommitNil(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		msg, err := pClaims.BClaims.BlockHash()
		if err != nil {
			return err
		}

		txr := []byte{197, 210, 70, 1, 134, 247, 35, 60, 146, 126, 125, 178, 220, 199, 3, 192, 229, 0, 182, 83, 202, 130, 39, 59, 123, 250, 216, 4, 93, 133, 164, 112}
		bs := []byte{137, 158, 164, 26, 219, 131, 151, 198, 183, 30, 184, 92, 126, 36, 84, 26, 33, 2, 95, 173, 235, 114, 104, 193, 225, 73, 193, 104, 229, 123, 61, 37, 111, 25, 109, 229, 148, 232, 96, 32, 23, 29, 116, 208, 88, 123, 82, 228, 215, 71, 195, 127, 104, 209, 148, 7, 41, 209, 77, 220, 127, 177, 247, 214, 0}
		roundStates.PeerStateMap[string([]byte{1})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{2})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}
		roundStates.PeerStateMap[string([]byte{3})].PreVote = &objs.PreVote{Proposal: &objs.Proposal{Signature: bs, PClaims: &objs.PClaims{BClaims: &objs.BClaims{ChainID: 42, Height: 1, PrevBlock: bhsh, StateRoot: msg, TxRoot: txr, HeaderRoot: msg},
			RCert: &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}}}, Signature: bs}

		roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].RCert = &objs.RCert{RClaims: &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}}

		err = stateHandler.dPVNSCastPCNFunc(roundStates)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPPVSDeadBlockRoundFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: txRoot, TxRoot: txRoot, StateRoot: txRoot,
		HeaderRoot: txRoot}

	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1, StateRoot: txRoot}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		ValidatorVAddrMap: map[string]int{string(ownState.VAddr): 1},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: []byte{2}}, {VAddr: []byte{2}, GroupShare: []byte{1}},
			{VAddr: []byte{3}, GroupShare: []byte{3}}, {VAddr: []byte{1}, GroupShare: []byte{4}}}}

	updateFunc = func(txn *badger.Txn) error {

		h := uint32(1)
		mdb.EXPECT().GetHeaderTrieRoot(txn, h).Return([]byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120,
			194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, nil)

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastPreVote(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		err = stateHandler.dPPVSDeadBlockRoundFunc(roundStates)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPPVSPreVoteNewHandler(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, ChainID: 42}
	rCert := &objs.RCert{RClaims: otherrClaims}
	otherrClaims2 := &objs.RClaims{Height: 2, Round: 1, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, ChainID: 42}
	rCert2 := &objs.RCert{RClaims: otherrClaims2}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 2, TxCount: 53, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, TxRoot: txRoot, StateRoot: []byte{173, 233, 94, 109, 13, 42, 99,
		22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147},
		HeaderRoot: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21,
			229, 112, 18, 48, 147}}

	// otherbClaims2 := &objs.BClaims{ChainID: 42, Height: 2, TxCount: 53, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
	// 	98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, TxRoot: txRoot, StateRoot: []byte{173, 233, 94, 109, 13, 42, 99,
	// 	22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147},
	// 	HeaderRoot: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21,
	// 		229, 112, 18, 48, 147}}

	bs := []byte{137, 158, 164, 26, 219, 131, 151, 198, 183, 30, 184, 92, 126, 36, 84, 26, 33, 2, 95, 173, 235, 114, 104, 193, 225, 73, 193, 104, 229, 123, 61, 37, 111, 25, 109, 229, 148, 232, 96, 32, 23, 29, 116, 208, 88, 123, 82, 228, 215, 71, 195, 127, 104, 209, 148, 7, 41, 209, 77, 220, 127, 177, 247, 214, 0}

	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	pClaims2 := &objs.PClaims{RCert: rCert2, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: bs, PClaims: pClaims}
	validValue2 := &objs.Proposal{Signature: bs, PClaims: pClaims2}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 2, Round: 1, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}}

	rClaims2 := &objs.RClaims{ChainID: 42, Height: 2, Round: 1, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}}

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}, Proposal: validValue}

	roundState2 := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims2}, Proposal: validValue}
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		ValidatorVAddrMap: map[string]int{string(ownState.VAddr): 1},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: []byte{2}}, {VAddr: []byte{2}, GroupShare: []byte{1}},
			{VAddr: []byte{3}, GroupShare: []byte{3}}, {VAddr: []byte{1}, GroupShare: []byte{4}}}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState2, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState2, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastPreVote(txn, gomock.Any()).Return(nil)

		hr := []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114,
			27, 21, 229, 112, 18, 48, 147}
		mdb.EXPECT().GetHeaderTrieRoot(txn, gomock.Any()).Return(hr, nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		roundStates.PeerStateMap["1"] = roundState

		roundStates.OwnValidatingState.ValidValue = validValue

		tempVAddr := roundStates.ValidatorSet.Validators[0].VAddr
		roundStates.PeerStateMap[string(tempVAddr)] = roundState2
		roundStates.PeerStateMap[string(tempVAddr)].Proposal = validValue2

		err = stateHandler.dPPVSPreVoteNewFunc(roundStates, validValue2)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPPVSPreVoteValid(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, ChainID: 42}
	rCert := &objs.RCert{RClaims: otherrClaims}

	// otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{2, 3}, TxRoot: txRoot, StateRoot: []byte{2},
	// 	HeaderRoot: []byte{3}}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, TxRoot: txRoot, StateRoot: []byte{173, 233, 94, 109, 13, 42, 99,
		22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147},
		HeaderRoot: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21,
			229, 112, 18, 48, 147}}

	bs := []byte{137, 158, 164, 26, 219, 131, 151, 198, 183, 30, 184, 92, 126, 36, 84, 26, 33, 2, 95, 173, 235, 114, 104, 193, 225, 73, 193, 104, 229, 123, 61, 37, 111, 25, 109, 229, 148, 232, 96, 32, 23, 29, 116, 208, 88, 123, 82, 228, 215, 71, 195, 127, 104, 209, 148, 7, 41, 209, 77, 220, 127, 177, 247, 214, 0}

	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: bs, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}}

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}, Proposal: validValue}
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		ValidatorVAddrMap: map[string]int{string(ownState.VAddr): 1},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: []byte{2}}, {VAddr: []byte{2}, GroupShare: []byte{1}},
			{VAddr: []byte{3}, GroupShare: []byte{3}}, {VAddr: []byte{1}, GroupShare: []byte{4}}}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastPreVote(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		roundStates.PeerStateMap["1"] = roundState

		err = stateHandler.dPPVSPreVoteValidFunc(roundStates, validValue)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPPVSPreVoteLocked(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, ChainID: 42}
	rCert := &objs.RCert{RClaims: otherrClaims}

	// otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{2, 3}, TxRoot: txRoot, StateRoot: []byte{2},
	// 	HeaderRoot: []byte{3}}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, TxRoot: txRoot, StateRoot: []byte{173, 233, 94, 109, 13, 42, 99,
		22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147},
		HeaderRoot: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21,
			229, 112, 18, 48, 147}}

	bs := []byte{137, 158, 164, 26, 219, 131, 151, 198, 183, 30, 184, 92, 126, 36, 84, 26, 33, 2, 95, 173, 235, 114, 104, 193, 225, 73, 193, 104, 229, 123, 61, 37, 111, 25, 109, 229, 148, 232, 96, 32, 23, 29, 116, 208, 88, 123, 82, 228, 215, 71, 195, 127, 104, 209, 148, 7, 41, 209, 77, 220, 127, 177, 247, 214, 0}

	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: bs, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, LockedValue: validValue}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137,
		98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}}

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}, Proposal: validValue}
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		ValidatorVAddrMap: map[string]int{string(ownState.VAddr): 1},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: []byte{2}}, {VAddr: []byte{2}, GroupShare: []byte{1}},
			{VAddr: []byte{3}, GroupShare: []byte{3}}, {VAddr: []byte{1}, GroupShare: []byte{4}}}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastPreVote(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		roundStates.PeerStateMap["1"] = roundState

		err = stateHandler.dPPVSPreVoteLockedFunc(roundStates, validValue)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPPVSPreVoteNilFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1}
	rCert := &objs.RCert{RClaims: otherrClaims}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{2, 3}, TxRoot: txRoot, StateRoot: []byte{2},
		HeaderRoot: []byte{3}}

	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bhsh, err := pClaims.BClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true},
		ValidatorVAddrMap: map[string]int{string(ownState.VAddr): 1},
		Validators: []*objs.Validator{{VAddr: []byte{1}, GroupShare: []byte{2}}, {VAddr: []byte{2}, GroupShare: []byte{1}},
			{VAddr: []byte{3}, GroupShare: []byte{3}}, {VAddr: []byte{1}, GroupShare: []byte{4}}}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)

		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{2}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{3}).Return(roundState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, []byte{1}).Return(roundState, nil)

		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastPreVoteNil(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		err = stateHandler.dPPVSPreVoteNilFunc(roundStates)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDPPSProposeNewFunc(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{254, 177, 174, 234, 241, 105, 181, 186, 254, 176, 225, 70, 95,
		60, 83, 73, 39, 227, 200, 167, 155, 247, 162, 137, 163, 156, 155, 59, 112, 14, 33, 172}, TxRoot: txRoot, StateRoot: []byte{173, 233, 94, 109, 13, 42, 99,
		22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147},
		HeaderRoot: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21,
			229, 112, 18, 48, 147}}

	bhsh, err := otherbClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, PrevBlock: []byte{254, 177, 174, 234, 241, 105, 181, 186, 254, 176, 225, 70, 95,
		60, 83, 73, 39, 227, 200, 167, 155, 247, 162, 137, 163, 156, 155, 59, 112, 14, 33, 172}, ChainID: 42}
	rCert := &objs.RCert{RClaims: otherrClaims}

	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, ValidValue: validValue}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{SigGroup: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true}}

	updateFunc = func(txn *badger.Txn) error {

		h := uint32(1)
		mdb.EXPECT().GetHeaderTrieRoot(txn, h).Return([]byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120,
			194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, nil)
		mdb.EXPECT().GetHeaderTrieRoot(txn, h).Return([]byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120,
			194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147}, nil)

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastProposal(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		// had to do this because otherwise the valid value was nil and was causing a nil error type thing
		stateHandler.appHandler.(*appmock.MockApplication).SetNextValidValue(validValue)

		roundStates.PeerStateMap[string(roundStates.OwnState.VAddr)].RCert = &objs.RCert{SigGroup: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			RClaims: &objs.RClaims{Height: 1, Round: 1}}

		err = stateHandler.dPPSProposeNewFunc(roundStates)
		if err != nil {
			t.Fatal("err is", err)
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCastProposalFromValue(t *testing.T) {

	bdb := getBdb(t)
	defer bdb.Close()

	ctr := gomock.NewController(t)
	defer ctr.Finish()
	mdb := db.NewMockDatabaseIface(ctr)

	var updateFunc db.TxnFunc

	stateHandler := getStateHandler(t, mdb)

	msstore := NewMockStore(mdb)

	txRoot, err := objs.MakeTxRoot([][]byte{})
	if err != nil {
		log.Fatal(err)
	}

	otherbClaims := &objs.BClaims{ChainID: 42, Height: 1, TxCount: 53, PrevBlock: []byte{254, 177, 174, 234, 241, 105, 181, 186, 254, 176, 225, 70, 95,
		60, 83, 73, 39, 227, 200, 167, 155, 247, 162, 137, 163, 156, 155, 59, 112, 14, 33, 172}, TxRoot: txRoot, StateRoot: []byte{173, 233, 94, 109, 13, 42, 99,
		22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21, 229, 112, 18, 48, 147},
		HeaderRoot: []byte{173, 233, 94, 109, 13, 42, 99, 22, 95, 251, 120, 194, 241, 137, 98, 59, 27, 223, 219, 43, 28, 200, 41, 191, 114, 27, 21,
			229, 112, 18, 48, 147}}

	bhsh, err := otherbClaims.BlockHash()
	if err != nil {
		log.Fatal(err)
	}

	otherrClaims := &objs.RClaims{Height: 1, Round: 1, PrevBlock: []byte{254, 177, 174, 234, 241, 105, 181, 186, 254, 176, 225, 70, 95,
		60, 83, 73, 39, 227, 200, 167, 155, 247, 162, 137, 163, 156, 155, 59, 112, 14, 33, 172}, ChainID: 42}
	rCert := &objs.RCert{RClaims: otherrClaims}

	pClaims := &objs.PClaims{RCert: rCert, BClaims: otherbClaims}
	validValue := &objs.Proposal{Signature: []byte{3, 3, 3}, PClaims: pClaims}
	ownValState := &objs.OwnValidatingState{VAddr: []byte{5, 5, 5}, LockedValue: validValue}

	bClaims := &objs.BClaims{ChainID: 42, Height: 1}
	maxBlockHeightSeen := &objs.BlockHeader{BClaims: bClaims}
	nextbClaims := &objs.BClaims{ChainID: 42, Height: 1}
	syncToBH := &objs.BlockHeader{BClaims: nextbClaims}

	ownState := &objs.OwnState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, MaxBHSeen: maxBlockHeightSeen, SyncToBH: syncToBH}

	rClaims := &objs.RClaims{ChainID: 42, Height: 1, Round: 1, PrevBlock: bhsh}

	roundState := &objs.RoundState{VAddr: []byte{5, 5, 5}, GroupKey: []byte{4, 4, 4}, RCert: &objs.RCert{GroupKey: []byte{1, 1, 1},
		RClaims: rClaims}}
	valSet := &objs.ValidatorSet{GroupKey: []byte{5, 4, 3}, ValidatorVAddrSet: map[string]bool{string(ownState.VAddr): true}}

	updateFunc = func(txn *badger.Txn) error {

		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetValidatorSet(txn, roundState.RCert.RClaims.Height).Return(valSet, nil)
		mdb.EXPECT().GetOwnState(txn).Return(ownState, nil)
		mdb.EXPECT().GetCurrentRoundState(txn, ownState.VAddr).Return(roundState, nil)
		mdb.EXPECT().GetOwnValidatingState(txn).Return(ownValState, nil)

		mdb.EXPECT().SetBroadcastProposal(txn, gomock.Any()).Return(nil)

		roundStates, err := msstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundStates.txn = txn

		err = stateHandler.castProposalFromValue(roundStates, validValue)
		if err != nil {
			t.Fatal("err is", err)
		}

		if roundStates.OwnState.MaxBHSeen.BClaims.Height != 1 {
			t.Fatal("incorrect value for one of the output values")
		}

		return nil
	}

	err = bdb.Update(updateFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func getGroupSig(msg []byte) ([]byte, [][]byte, []byte, [][]byte) {
	s := new(mcrypto.BNGroupSigner)

	secret1 := big.NewInt(100)
	secret2 := big.NewInt(101)
	secret3 := big.NewInt(102)
	secret4 := big.NewInt(103)

	msk := big.NewInt(0)
	msk.Add(msk, secret1)
	msk.Add(msk, secret2)
	msk.Add(msk, secret3)
	msk.Add(msk, secret4)
	msk.Mod(msk, bn256.Order)
	mpk := new(bn256.G2).ScalarBaseMult(msk)

	big1 := big.NewInt(1)
	big2 := big.NewInt(2)

	privCoefs1 := []*big.Int{secret1, big1, big2}
	privCoefs2 := []*big.Int{secret2, big1, big2}
	privCoefs3 := []*big.Int{secret3, big1, big2}
	privCoefs4 := []*big.Int{secret4, big1, big2}

	share1to1 := bn256.PrivatePolyEval(privCoefs1, 1)
	share1to2 := bn256.PrivatePolyEval(privCoefs1, 2)
	share1to3 := bn256.PrivatePolyEval(privCoefs1, 3)
	share1to4 := bn256.PrivatePolyEval(privCoefs1, 4)
	share2to1 := bn256.PrivatePolyEval(privCoefs2, 1)
	share2to2 := bn256.PrivatePolyEval(privCoefs2, 2)
	share2to3 := bn256.PrivatePolyEval(privCoefs2, 3)
	share2to4 := bn256.PrivatePolyEval(privCoefs2, 4)
	share3to1 := bn256.PrivatePolyEval(privCoefs3, 1)
	share3to2 := bn256.PrivatePolyEval(privCoefs3, 2)
	share3to3 := bn256.PrivatePolyEval(privCoefs3, 3)
	share3to4 := bn256.PrivatePolyEval(privCoefs3, 4)
	share4to1 := bn256.PrivatePolyEval(privCoefs4, 1)
	share4to2 := bn256.PrivatePolyEval(privCoefs4, 2)
	share4to3 := bn256.PrivatePolyEval(privCoefs4, 3)
	share4to4 := bn256.PrivatePolyEval(privCoefs4, 4)

	groupShares := make([][]byte, 4)
	for k := 0; k < len(groupShares); k++ {
		groupShares[k] = make([]byte, len(mpk.Marshal()))
	}

	listOfSS1 := []*big.Int{share1to1, share2to1, share3to1, share4to1}
	gsk1 := bn256.GenerateGroupSecretKeyPortion(listOfSS1)
	gpk1 := new(bn256.G2).ScalarBaseMult(gsk1)
	groupShares[0] = gpk1.Marshal()
	s1 := new(mcrypto.BNGroupSigner)
	s1.SetPrivk(gsk1.Bytes())
	sig1, err := s1.Sign(msg)
	if err != nil {
		log.Fatal(err)
	}

	listOfSS2 := []*big.Int{share1to2, share2to2, share3to2, share4to2}
	gsk2 := bn256.GenerateGroupSecretKeyPortion(listOfSS2)
	gpk2 := new(bn256.G2).ScalarBaseMult(gsk2)
	groupShares[1] = gpk2.Marshal()
	s2 := new(mcrypto.BNGroupSigner)
	s2.SetPrivk(gsk2.Bytes())
	sig2, err := s2.Sign(msg)
	if err != nil {
		log.Fatal(err)
	}

	listOfSS3 := []*big.Int{share1to3, share2to3, share3to3, share4to3}
	gsk3 := bn256.GenerateGroupSecretKeyPortion(listOfSS3)
	gpk3 := new(bn256.G2).ScalarBaseMult(gsk3)
	groupShares[2] = gpk3.Marshal()
	s3 := new(mcrypto.BNGroupSigner)
	s3.SetPrivk(gsk3.Bytes())
	sig3, err := s3.Sign(msg)
	if err != nil {
		log.Fatal(err)
	}

	listOfSS4 := []*big.Int{share1to4, share2to4, share3to4, share4to4}
	gsk4 := bn256.GenerateGroupSecretKeyPortion(listOfSS4)
	gpk4 := new(bn256.G2).ScalarBaseMult(gsk4)
	groupShares[3] = gpk4.Marshal()
	s4 := new(mcrypto.BNGroupSigner)
	s4.SetPrivk(gsk4.Bytes())
	sig4, err := s4.Sign(msg)
	if err != nil {
		log.Fatal(err)
	}

	sigs := make([][]byte, 4)
	for k := 0; k < len(sigs); k++ {
		sigs[k] = make([]byte, 192)
	}
	sigs[0] = sig1
	sigs[1] = sig2
	sigs[2] = sig3
	sigs[3] = sig4

	// Attempt with empty GroupShares
	emptyShares := make([][]byte, 4)
	_, err = s.Aggregate(sigs, emptyShares)
	if err == nil {
		log.Fatal("Error should have been raised for invalid groupShares!")
	}

	// Attempt without groupPubk
	_, err = s.Aggregate(sigs, groupShares)
	if err != mcrypto.ErrPubkeyGroupNotSet {
		log.Fatal("Error should have been raised for no PubkeyGroup!")
	}
	err = s.SetGroupPubk(mpk.Marshal())
	if err != nil {
		log.Fatal(err)
	}

	// Finally submit signature
	grpsig, err := s.Aggregate(sigs, groupShares)
	if err != nil {
		log.Fatal(err)
	}

	// leaving this check here for now since so it could fail earlier if it is not working correctly - could probably remove this
	// at some point
	v := new(mcrypto.BNGroupValidator)
	_, err = v.Validate(msg, grpsig)
	if err != nil {
		log.Fatal(err)
	}

	return grpsig, groupShares, mpk.Marshal(), sigs
}
