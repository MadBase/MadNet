package lstate

import (
	"context"
	appObjs "github.com/MadBase/MadNet/application/objs"
	"github.com/MadBase/MadNet/config"
	"github.com/MadBase/MadNet/consensus/admin"
	"github.com/MadBase/MadNet/consensus/appmock"
	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/consensus/dman"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/consensus/request"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/crypto"
	"github.com/MadBase/MadNet/dynamics"
	"github.com/MadBase/MadNet/logging"
	"github.com/MadBase/MadNet/utils"
	"github.com/dgraph-io/badger/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEngine_Status_Ok(t *testing.T) {
	st := make(map[string]interface{})
	engine := initEngine(t)

	os := createOwnState(t, 1)
	rs := createRoundState(t, os)
	vs := createValidatorsSet(os, rs)

	_ = engine.sstore.database.Update(func(txn *badger.Txn) error {
		err := engine.sstore.database.SetOwnState(txn, os)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		err = engine.sstore.database.SetCurrentRoundState(txn, rs)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		err = engine.sstore.database.SetValidatorSet(txn, vs)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		return nil
	})

	st, err := engine.Status(st)
	assert.Nil(t, err)
}

func TestEngine_Status_Error(t *testing.T) {
	st := make(map[string]interface{})
	engine := initEngine(t)

	_, err := engine.Status(st)
	assert.NotNil(t, err)
}

//ce.ethAcct != ownState.VAddr
func TestEngine_UpdateLocalState1(t *testing.T) {
	engine := initEngine(t)
	os := createOwnState(t, 1)
	rs := createRoundState(t, os)
	vs := createValidatorsSet(os, rs)
	rss := createRoundStates(os, rs, vs, &objs.OwnValidatingState{})

	_ = engine.sstore.database.Update(func(txn *badger.Txn) error {
		err := engine.sstore.WriteState(txn, rss)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		err = engine.sstore.database.SetValidatorSet(txn, vs)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		return nil
	})

	isSync, err := engine.UpdateLocalState()
	assert.Nil(t, err)
	assert.True(t, isSync)
}

//ce.ethAcct == ownState.VAddr
//os val GetPrivK not found
func TestEngine_UpdateLocalState2(t *testing.T) {
	engine := initEngine(t)
	os := createOwnState(t, 1)
	rs := createRoundState(t, os)
	vs := createValidatorsSet(os, rs)
	rss := createRoundStates(os, rs, vs, &objs.OwnValidatingState{})
	engine.ethAcct = os.VAddr

	_ = engine.sstore.database.Update(func(txn *badger.Txn) error {
		err := engine.sstore.WriteState(txn, rss)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		err = engine.sstore.database.SetValidatorSet(txn, vs)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		return nil
	})

	isSync, err := engine.UpdateLocalState()
	assert.Nil(t, err)
	assert.True(t, isSync)
}

//ce.ethAcct == ownState.VAddr
//os val GetPrivK found but pubk mismatch
func TestEngine_UpdateLocalState3(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Should have raised panic: pubkey mismatch!")
		}
	}()

	engine := initEngine(t)
	os := createOwnState(t, 1)
	rs := createRoundState(t, os)
	vs := createValidatorsSet(os, rs)
	rss := createRoundStates(os, rs, vs, &objs.OwnValidatingState{})
	engine.ethAcct = os.VAddr

	_ = engine.sstore.database.Update(func(txn *badger.Txn) error {
		err := engine.sstore.WriteState(txn, rss)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		err = engine.sstore.database.SetValidatorSet(txn, vs)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		osVal := vs.Validators[len(vs.Validators)-1]
		es := &objs.EncryptedStore{
			Name: osVal.GroupShare,
		}

		err = es.Encrypt(engine.AdminBus)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		err = engine.database.SetEncryptedStore(txn, es)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		return nil
	})

	_, _ = engine.UpdateLocalState()
}

//ce.ethAcct == ownState.VAddr
//os val GetPrivK not found
//new validators set
func TestEngine_UpdateLocalState4(t *testing.T) {
	engine := initEngine(t)
	os := createOwnState(t, 1)
	rs := createRoundState(t, os)
	vs := createValidatorsSet(os, rs)
	rss := createRoundStates(os, rs, vs, &objs.OwnValidatingState{})
	engine.ethAcct = os.VAddr

	_ = engine.sstore.database.Update(func(txn *badger.Txn) error {
		err := engine.sstore.WriteState(txn, rss)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		err = engine.sstore.database.SetValidatorSet(txn, vs)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		groupSigner := &crypto.BNGroupSigner{}
		err = groupSigner.SetPrivk(crypto.Hasher([]byte("secret123")))
		if err != nil {
			t.Fatal(err)
		}
		groupKey, _ := groupSigner.PubkeyShare()

		vs.GroupKey = groupKey
		vs.NotBefore = 2
		err = engine.database.SetValidatorSetPostApplication(txn, vs, 1)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		return nil
	})

	isSync, err := engine.UpdateLocalState()
	assert.Nil(t, err)
	assert.True(t, isSync)
}

//updateLoadedObjects = OK
//updateLocalStateInternal = OK
func TestEngine_UpdateLocalState5(t *testing.T) {
	engine := initEngine(t)
	os := createOwnState(t, 1)
	rs := createRoundState(t, os)
	vs := createValidatorsSet(os, rs)
	rss := createRoundStates(os, rs, vs, &objs.OwnValidatingState{})
	engine.ethAcct = os.VAddr
	os.GroupKey = vs.GroupKey

	_ = engine.sstore.database.Update(func(txn *badger.Txn) error {
		err := engine.sstore.WriteState(txn, rss)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		err = engine.sstore.database.SetValidatorSet(txn, vs)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		osVal := vs.Validators[len(vs.Validators)-1]
		es := &objs.EncryptedStore{
			Name: osVal.GroupShare,
		}

		err = es.Encrypt(engine.AdminBus)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		err = engine.database.SetEncryptedStore(txn, es)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		signer := &crypto.BNGroupSigner{}
		pk := utils.CopySlice(es.ClearText)
		err = signer.SetPrivk(pk)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}
		engine.bnSigner = signer

		return nil
	})

	isSync, err := engine.UpdateLocalState()
	assert.Nil(t, err)
	assert.True(t, isSync)
}

//updateLoadedObjects = OK
//updateLocalStateInternal = OK
//bHeight = 1024 and not safe to proceed
func TestEngine_UpdateLocalState6(t *testing.T) {
	engine := initEngine(t)
	os := createOwnState(t, 1024)
	rs := createRoundState(t, os)
	vs := createValidatorsSet(os, rs)
	rss := createRoundStates(os, rs, vs, &objs.OwnValidatingState{})
	engine.ethAcct = os.VAddr
	os.GroupKey = vs.GroupKey

	_ = engine.sstore.database.Update(func(txn *badger.Txn) error {
		err := engine.sstore.WriteState(txn, rss)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		err = engine.sstore.database.SetValidatorSet(txn, vs)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		osVal := vs.Validators[len(vs.Validators)-1]
		es := &objs.EncryptedStore{
			Name: osVal.GroupShare,
		}

		err = es.Encrypt(engine.AdminBus)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		err = engine.database.SetEncryptedStore(txn, es)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		signer := &crypto.BNGroupSigner{}
		pk := utils.CopySlice(es.ClearText)
		err = signer.SetPrivk(pk)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}
		engine.bnSigner = signer

		return nil
	})

	isSync, err := engine.UpdateLocalState()
	assert.Nil(t, err)
	assert.True(t, isSync)
}

//updateLoadedObjects = OK
//updateLocalStateInternal = OK
//bHeight = 1024 and safe to proceed
func TestEngine_UpdateLocalState7(t *testing.T) {
	engine := initEngine(t)
	os := createOwnState(t, 1024)
	rs := createRoundState(t, os)
	vs := createValidatorsSet(os, rs)
	rss := createRoundStates(os, rs, vs, &objs.OwnValidatingState{})
	engine.ethAcct = os.VAddr
	os.GroupKey = vs.GroupKey

	_ = engine.sstore.database.Update(func(txn *badger.Txn) error {
		err := engine.sstore.WriteState(txn, rss)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		err = engine.sstore.database.SetValidatorSet(txn, vs)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		osVal := vs.Validators[len(vs.Validators)-1]
		es := &objs.EncryptedStore{
			Name: osVal.GroupShare,
		}

		err = es.Encrypt(engine.AdminBus)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		err = engine.database.SetEncryptedStore(txn, es)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		signer := &crypto.BNGroupSigner{}
		pk := utils.CopySlice(es.ClearText)
		err = signer.SetPrivk(pk)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}
		engine.bnSigner = signer

		err = engine.database.SetSafeToProceed(txn, 1025, true)
		if err != nil {
			t.Fatalf("Shouldn't have raised error: %v", err)
		}

		return nil
	})

	isSync, err := engine.UpdateLocalState()
	assert.Nil(t, err)
	assert.True(t, isSync)
}

func initEngine(t *testing.T) *Engine {
	ctx := context.Background()
	logger := logging.GetLogger("test")

	rawEngineDb, err := utils.OpenBadger(ctx.Done(), "", true)
	if err != nil {
		t.Fatal(err)
	}
	engineDb := &db.Database{}
	engineDb.Init(rawEngineDb)

	p2pClientMock := &request.P2PClientMock{}
	client := &request.Client{}
	client.Init(p2pClientMock, &dynamics.Storage{})

	app := appmock.New()

	secpSigner := &crypto.Secp256k1Signer{}
	err = secpSigner.SetPrivk(crypto.Hasher([]byte("secret")))
	if err != nil {
		t.Fatal(err)
	}

	bnSigner := &crypto.BNGroupSigner{}
	err = bnSigner.SetPrivk(crypto.Hasher([]byte("secret2")))
	if err != nil {
		t.Fatal(err)
	}

	adminBus := initAdminBus(t, logger, engineDb)
	storage := appObjs.MakeMockStorageGetter()
	reqBusViewMock := &dman.ReqBusViewMock{}
	dMan := &dman.DMan{}
	dMan.Init(engineDb, app, reqBusViewMock)

	engine := &Engine{}
	engine.Init(engineDb, dMan, app, secpSigner, adminBus, make([]byte, constants.HashLen), client, storage)

	return engine
}

func initAdminBus(t *testing.T, logger *logrus.Logger, db *db.Database) *admin.Handlers {
	app := appmock.New()
	s := initStorage(t, logger)

	handler := &admin.Handlers{}
	handler.Init(1, db, crypto.Hasher([]byte(config.Configuration.Validator.SymmetricKey)), app, make([]byte, constants.HashLen), s, nil)

	return handler
}

func initStorage(t *testing.T, logger *logrus.Logger) *dynamics.Storage {
	s := &dynamics.Storage{}
	err := s.Init(&dynamics.MockRawDB{}, logger)
	if err != nil {
		t.Fatal(err)
	}

	return s
}
