package lstate

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/MadBase/MadNet/errorz"
	"github.com/MadBase/MadNet/utils"

	"github.com/MadBase/MadNet/consensus/admin"
	"github.com/MadBase/MadNet/consensus/appmock"
	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/consensus/request"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/crypto"
	"github.com/MadBase/MadNet/logging"
	"github.com/dgraph-io/badger/v2"
	"github.com/sirupsen/logrus"
)

// Engine is the consensus algorithm parent object.
type Engine struct {
	ctx       context.Context
	cancelCtx func()

	database db.DatabaseIface
	sstore   *Store

	RequestBus *request.Client
	appHandler appmock.Application

	logger     *logrus.Logger
	secpSigner *crypto.Secp256k1Signer
	bnSigner   *crypto.BNGroupSigner

	AdminBus *admin.Handlers

	fastSync *SnapShotManager

	ethAcct []byte
	EthPubk []byte

	dm *DMan
}

// Init will initialize the Consensus Engine and all sub modules
func (ce *Engine) Init(database db.DatabaseIface, dm *DMan, app appmock.Application, signer *crypto.Secp256k1Signer, adminHandlers *admin.Handlers, publicKey []byte, rbusClient *request.Client) error {
	background := context.Background()
	ctx, cf := context.WithCancel(background)
	ce.cancelCtx = cf
	ce.ctx = ctx
	ce.secpSigner = signer
	ce.database = database
	ce.AdminBus = adminHandlers
	ce.EthPubk = publicKey
	ce.RequestBus = rbusClient
	ce.appHandler = app
	ce.sstore = &Store{}
	err := ce.sstore.Init(database)
	if err != nil {
		return err
	}
	ce.dm = dm
	if len(ce.EthPubk) > 0 {
		ce.ethAcct = crypto.GetAccount(ce.EthPubk)
	}
	ce.logger = logging.GetLogger(constants.LoggerConsensus)
	ce.fastSync = &SnapShotManager{
		appHandler: app,
		requestBus: ce.RequestBus,
	}
	if err := ce.fastSync.Init(database); err != nil {
		return err
	}
	return nil
}

// Status .
func (ce *Engine) Status(status map[string]interface{}) (map[string]interface{}, error) {
	var rs *RoundStates
	err := ce.database.View(func(txn *badger.Txn) error {
		rss, err := ce.sstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		rs = rss
		return nil
	})
	if err != nil {
		return nil, err
	}
	bhsh, err := rs.OwnState.SyncToBH.BlockHash()
	if err != nil {
		return nil, err
	}
	if rs.OwnState.MaxBHSeen.BClaims.Height-rs.OwnState.SyncToBH.BClaims.Height < 2 {
		status[constants.StatusBlkRnd] = fmt.Sprintf("%d/%d", rs.OwnState.SyncToBH.BClaims.Height, rs.OwnRoundState().RCert.RClaims.Round)
		status[constants.StatusBlkHsh] = fmt.Sprintf("%x..%x", bhsh[0:2], bhsh[len(bhsh)-2:])
		status[constants.StatusTxCt] = rs.OwnState.SyncToBH.BClaims.TxCount
		return status, nil
	}
	status[constants.StatusBlkRnd] = fmt.Sprintf("%d/%v", rs.OwnState.MaxBHSeen.BClaims.Height, "-")
	status[constants.StatusBlkHsh] = fmt.Sprintf("%x..%x", bhsh[0:2], bhsh[len(bhsh)-2:])
	status[constants.StatusSyncToBlk] = fmt.Sprintf("%d", rs.OwnState.SyncToBH.BClaims.Height)
	return status, nil
}

// UpdateLocalState .
func (ce *Engine) UpdateLocalState() (bool, error) {
	var isSync bool
	updateLocalState := true
	// ce.logger.Error("!!! OPEN UpdateLocalState TXN")
	// defer func() { ce.logger.Error("!!! CLOSE UpdateLocalState TXN") }()
	err := ce.database.Update(func(txn *badger.Txn) error {
		ownState, err := ce.database.GetOwnState(txn)
		if err != nil {
			return err
		}
		height, _ := objs.ExtractHR(ownState.SyncToBH)
		vs, err := ce.database.GetValidatorSet(txn, height)
		if err != nil {
			return err
		}
		if !bytes.Equal(vs.GroupKey, ownState.GroupKey) {
			ownState.GroupKey = vs.GroupKey
			err = ce.database.SetOwnState(txn, ownState)
			if err != nil {
				return err
			}
		}
		ownValidatingState, err := ce.database.GetOwnValidatingState(txn)
		if err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}
		}
		if ownValidatingState == nil {
			ownValidatingState = &objs.OwnValidatingState{}
		}
		if !bytes.Equal(ownValidatingState.VAddr, ownState.VAddr) || !bytes.Equal(ownValidatingState.GroupKey, ownState.GroupKey) {
			ovs := &objs.OwnValidatingState{
				VAddr:    ownState.VAddr,
				GroupKey: ownState.GroupKey,
			}
			ovs.SetRoundStarted()
			err := ce.database.SetOwnValidatingState(txn, ovs)
			if err != nil {
				return err
			}
		}
		roundState, err := ce.sstore.LoadLocalState(txn)
		if err != nil {
			return err
		}
		roundState.txn = txn
		if roundState.OwnState.SyncToBH.BClaims.Height%constants.EpochLength == 0 {
			safe, err := ce.database.GetSafeToProceed(txn, roundState.OwnState.SyncToBH.BClaims.Height)
			if err != nil {
				utils.DebugTrace(ce.logger, err)
				return err
			}
			if !safe {
				utils.DebugTrace(ce.logger, nil, "not safe")
				updateLocalState = false
			}
		}
		if roundState.OwnState.SyncToBH.BClaims.Height < roundState.OwnState.MaxBHSeen.BClaims.Height {
			isSync = false
			updateLocalState = false
		}
		if updateLocalState {
			ok, err := ce.updateLocalStateInternal(roundState)
			if err != nil {
				return err
			}
			isSync = ok
		}
		err = ce.sstore.WriteState(roundState)
		if err != nil {
			utils.DebugTrace(ce.logger, err)
			return err
		}
		return nil
	})
	if err != nil {
		e := errorz.ErrInvalid{}.New("")
		if !errors.As(err, &e) && err != errorz.ErrMissingTransactions {
			return false, err
		}
		return false, nil
	}
	err = ce.database.Sync()
	if err != nil {
		return false, err
	}
	return isSync, nil
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// This changes state for local node
// order of ops is as follows:
//
//  check for height jump
//  check for dead block round round jump
//  follow dead block round next round if signed by self
//  follow f+1 other dead block round next round messages
//  do own next height if in dead block round
//  follow a next height from any round in the same height as us
//      this is safe due to how we count next heights to filter
//      dead block round
//  follow a round jump to any non dead block round
//  do a possible next round in same round
//  do a possible precommit/pendingnext in same round
//  do a possible precommit/pendingnext in same round
//  do a possible prevote/pendingprecommit in same round
//  do a possible prevotenil/pendingprecommit in same round
//  do a possible pending prevote
//  do a possible do a proposal if not already proposed and is proposer
//  do nothing if not any of above is true
func (ce *Engine) updateLocalStateInternal(rs *RoundStates) (bool, error) {
	if err := ce.loadValidationKey(rs); err != nil {
		return false, nil
	}
	os := rs.OwnRoundState()

	// extract the round cert for use
	rcert := os.RCert

	// create three vectors that may overlap
	// these vectors sort all current validators by height/round
	// as is determined by their respective rcert
	// these vectors are:
	//  Current height future round
	//  Current height any round
	//  Future height any round
	ChFr := []*objs.RoundState{}
	FH := []*objs.RoundState{}
	for i := 0; i < len(rs.ValidatorSet.Validators); i++ {
		vObj := rs.ValidatorSet.Validators[i]
		vAddr := vObj.VAddr
		vroundState := rs.PeerStateMap[string(vAddr)]
		relationH := objs.RelateH(rcert, vroundState.RCert)
		if relationH == 0 {
			relationHR := objs.RelateHR(rcert, vroundState.RCert)
			if relationHR == -1 {
				ChFr = append(ChFr, vroundState)
			}
		} else if relationH == -1 {
			FH = append(FH, vroundState)
		}
	}

	if len(FH) > 0 {
		var maxHR *objs.RoundState
		var maxHeight uint32
		var vroundState *objs.RoundState
		for i := 0; i < len(FH); i++ {
			vroundState = FH[i]
			if vroundState.RCert.RClaims.Height > maxHeight {
				maxHR = vroundState
				maxHeight = vroundState.RCert.RClaims.Height
			}
		}
		if maxHR != nil {
			rs.maxHR = maxHR
		}
	}

	// var currentHandler handler

	// if there are ANY peers in a future height, try to follow
	// we should try to follow the max height possible
	// currentHandler = fhHandler{ce: ce, txn: txn, rs: rs, FH: FH}
	// if currentHandler.evalCriteria() {
	// 	return currentHandler.evalLogic()
	// }

	// at this point no height jump is possible
	// otherwise we would have followed it

	// check for next height messages
	// if one exists, follow it
	NHs, _, err := rs.GetCurrentNext()
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	NHCurrent := os.NHCurrent(rcert)

	// iterate all possibles from nextRound down to proposal
	// and take that action
	PCurrent := os.PCurrent(rcert)
	PVCurrent := os.PVCurrent(rcert)
	PVNCurrent := os.PVNCurrent(rcert)
	PCCurrent := os.PCCurrent(rcert)
	PCNCurrent := os.PCNCurrent(rcert)
	NRCurrent := os.NRCurrent(rcert)

	// local node cast a precommit nil this round
	// count the precommits
	pcl, _, err := rs.GetCurrentPreCommits()
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}

	// last of all count next round messages from this round only
	_, nrl, err := rs.GetCurrentNext()
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}

	//
	//	Order:
	//		NR, PC, PCN, PV, PVN, ProposalTO, IsProposer

	ceHandlers := []handler{
		// absorbed further handler into fh handler
		&fhHandler{rs: rs, external: ce.fhFunc},
		// &rCertHandler{rs: rs, ChFr: ChFr, external: ce.rCertFunc},
		// split r cert handler into the following two handlers
		&roundJumpUpdateValidValueHandler{rs: rs, ChFr: ChFr, pcl: pcl, external: ce.roundJumpUpdateValidValueFunc},
		&roundJumpSetRCertHandler{rs: rs, ChFr: ChFr, pcl: pcl, external: ce.roundJumpSetRCertFunc},

		// &doNrHandler{rs: rs, rcert: rcert, external: ce.doNrFunc},
		// split do nr handler into the following two
		&doNrsCastNextHeightHandler{rs: rs, rcert: rcert, pcl: pcl, external: ce.doNrsCastNextHeightHandlerFunc},
		&doNrsCastNextRoundHandler{rs: rs, rcert: rcert, pcl: pcl, nrl: nrl, external: ce.doNrsCastNextRoundHandlerFunc},

		// this handler does not have any other further handlers
		&castNhHandler{rs: rs, NHs: NHs, NHCurrent: NHCurrent, external: ce.castNhFunc},

		// &doNextHeightHandler{rs: rs, NHs: NHs, NHCurrent: NHCurrent, external: ce.doNextHeightFunc},
		// flattened handler pair from above into the following
		&castBlockHeaderHandler{rs: rs, NHs: NHs, NHCurrent: NHCurrent, external: ce.castBlockHeaderHandlerFunc},

		// &doRoundJumpHandler{rs: rs, ChFr: ChFr, external: ce.doRoundJumpFunc},
		// split the above into the following two
		&doRoundJumpUpdateValidValueHandler{rs: rs, ChFr: ChFr, pcl: pcl, external: ce.roundJumpUpdateValidValueFunc},
		&doRoundJumpSetRCertHandler{rs: rs, ChFr: ChFr, pcl: pcl, external: ce.roundJumpSetRCertFunc},

		// -----

		// &nrCurrentHandler{rs: rs, NRCurrent: NRCurrent, external: ce.nrCurrentFunc},
		// split the above into the following two
		&nrCurrentCnhHandler{rs: rs, NRCurrent: NRCurrent, rcert: rcert, pcl: pcl, external: ce.doNrsCastNextHeightHandlerFunc},
		&nrCurrentCnrHandler{rs: rs, NRCurrent: NRCurrent, rcert: rcert, pcl: pcl, nrl: nrl, external: ce.doNrsCastNextRoundHandlerFunc},

		&pcCurrentHandler{rs: rs, PCCurrent: PCCurrent, external: ce.pcCurrentFunc},
		&pcnCurrentHandler{rs: rs, PCNCurrent: PCNCurrent, external: ce.pcnCurrentFunc},
		&pvCurrentHandler{rs: rs, PVCurrent: PVCurrent, external: ce.pvCurrentFunc},
		&pvnCurrentHandler{rs: rs, PVNCurrent: PVNCurrent, external: ce.pvnCurrentFunc},
		&ptoExpiredHandler{rs: rs, external: ce.ptoExpiredFunc},
		&validPropHandler{rs: rs, PCurrent: PCurrent, external: ce.validPropFunc},
	}

	for i := 0; i < len(ceHandlers); i++ {
		if ceHandlers[i].evalCriteria() {
			ok, err := ceHandlers[i].evalLogic()
			if err != nil {
				return false, err
			}
			return ok, nil
		}
	}

	return true, nil
}

type handler interface {
	evalCriteria() bool
	evalLogic() (bool, error)
}

type fhHandler struct {
	rs       *RoundStates
	external func(*RoundStates) (bool, error)
}

func (fhh *fhHandler) evalCriteria() bool {
	rcert := fhh.rs.maxHR.RCert
	cond1 := rcert.RClaims.Height <= fhh.rs.Height()+1
	cond2 := fhh.rs.ValidValueCurrent()
	return fhh.rs.maxHR != nil && cond1 && cond2
}

func (fhh *fhHandler) evalLogic() (bool, error) {
	return fhh.external(fhh.rs)
}

func (ce *Engine) fhFunc(rs *RoundStates) (bool, error) {
	rcert := rs.maxHR.RCert

	// get the last element of the sorted future height
	// if we can just jump up to this height, do so.
	// if the height is only one more, we can simply move to this
	// height if everything looks okay

	// if we have a valid value, check if the valid value
	// matches the previous blockhash of block n+1
	// if so, form the block and jump up to this level
	// this is safe because we call isValid on all values
	// before storing in a lock

	bhsh, err := rs.ValidValue().PClaims.BClaims.BlockHash()
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	if bytes.Equal(bhsh, rcert.RClaims.PrevBlock) && rcert.RClaims.Round == 1 {
		vv := rs.ValidValue()
		err := ce.castNewCommittedBlockFromProposalAndRCert(rs, vv, rcert)
		if err != nil {
			var e *errorz.ErrInvalid
			if err != errorz.ErrMissingTransactions && !errors.As(err, &e) {
				utils.DebugTrace(ce.logger, err)
				return false, err
			}
		}
		rs.OwnValidatingState.ValidValue = nil
		rs.OwnValidatingState.LockedValue = nil
	}

	// we can not do anything from here without a ton of work
	// so easier to just wait for the next block header to unsync us

	return true, nil
}

type rCertHandler struct {
	rs       *RoundStates
	maxRCert *objs.RCert
	ChFr     []*objs.RoundState
	external func(*RoundStates, *objs.RCert) (bool, error)
}

func (rch *rCertHandler) evalCriteria() bool {
	if len(rch.ChFr) > 0 {
		var vroundState *objs.RoundState
		for i := 0; i < len(rch.ChFr); i++ {
			vroundState = rch.ChFr[i]
			if vroundState.RCert.RClaims.Round == constants.DEADBLOCKROUND {
				rch.maxRCert = vroundState.RCert
				return true
			}
		}
	}
	return false
}

func (rch *rCertHandler) evalLogic() (bool, error) {
	return rch.external(rch.rs, rch.maxRCert)
}

func (ce *Engine) rCertFunc(rs *RoundStates, maxRCert *objs.RCert) (bool, error) {
	err := ce.doRoundJump(rs, maxRCert)
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	return true, nil
}

type roundJumpUpdateValidValueHandler struct {
	rs       *RoundStates
	maxRCert *objs.RCert
	ChFr     []*objs.RoundState
	pcl      objs.PreCommitList
	external func(*RoundStates, *objs.RCert, objs.PreCommitList) (bool, error)
}

func (rch *roundJumpUpdateValidValueHandler) evalCriteria() bool {
	cond := len(rch.pcl) > rch.rs.GetCurrentThreshold()
	if len(rch.ChFr) > 0 {
		var vroundState *objs.RoundState
		for i := 0; i < len(rch.ChFr); i++ {
			vroundState = rch.ChFr[i]
			if vroundState.RCert.RClaims.Round == constants.DEADBLOCKROUND {
				rch.maxRCert = vroundState.RCert
				if cond {
					return true
				}
			}
		}
	}
	return false
}

func (rch *roundJumpUpdateValidValueHandler) evalLogic() (bool, error) {
	return rch.external(rch.rs, rch.maxRCert, rch.pcl)
}

func (ce *Engine) roundJumpUpdateValidValueFunc(rs *RoundStates, maxRCert *objs.RCert, pcl objs.PreCommitList) (bool, error) {

	p, err := pcl.GetProposal()
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	if err := ce.updateValidValue(rs, p); err != nil {
		var e *errorz.ErrInvalid
		if err != errorz.ErrMissingTransactions && !errors.As(err, &e) {
			utils.DebugTrace(ce.logger, err)
			return false, err
		}
	}
	if err := ce.setMostRecentRCert(rs, maxRCert); err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}

	return true, nil
}

type roundJumpSetRCertHandler struct {
	rs       *RoundStates
	maxRCert *objs.RCert
	ChFr     []*objs.RoundState
	pcl      objs.PreCommitList
	external func(*RoundStates, *objs.RCert, objs.PreCommitList) (bool, error)
}

func (rch *roundJumpSetRCertHandler) evalCriteria() bool {
	cond := len(rch.pcl) <= rch.rs.GetCurrentThreshold()
	if len(rch.ChFr) > 0 {
		var vroundState *objs.RoundState
		for i := 0; i < len(rch.ChFr); i++ {
			vroundState = rch.ChFr[i]
			if vroundState.RCert.RClaims.Round == constants.DEADBLOCKROUND {
				rch.maxRCert = vroundState.RCert
				if cond {
					return true
				}
			}
		}
	}
	return false
}

func (rch *roundJumpSetRCertHandler) evalLogic() (bool, error) {
	return rch.external(rch.rs, rch.maxRCert, rch.pcl)
}

func (ce *Engine) roundJumpSetRCertFunc(rs *RoundStates, maxRCert *objs.RCert, pcl objs.PreCommitList) (bool, error) {

	if err := ce.setMostRecentRCert(rs, maxRCert); err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}

	return true, nil
}

type doNrHandler struct {
	rs       *RoundStates
	rcert    *objs.RCert
	external func(*RoundStates, *objs.RCert) (bool, error)
}

func (dnrh *doNrHandler) evalCriteria() bool {
	if dnrh.rs.OwnRoundState().NextRound != nil {
		if dnrh.rs.OwnRoundState().NRCurrent(dnrh.rcert) {
			return dnrh.rcert.RClaims.Round == constants.DEADBLOCKROUNDNR
		}
	}
	return false
}

func (dnrh *doNrHandler) evalLogic() (bool, error) {
	return dnrh.external(dnrh.rs, dnrh.rcert)
}

func (ce *Engine) doNrFunc(rs *RoundStates, rcert *objs.RCert) (bool, error) {
	err := ce.doNextRoundStep(rs)
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	return true, nil
}

type doNrsCastNextHeightHandler struct {
	rs       *RoundStates
	rcert    *objs.RCert
	pcl      objs.PreCommitList
	external func(*RoundStates, *objs.RCert, objs.PreCommitList) (bool, error)
}

func (dnrh *doNrsCastNextHeightHandler) evalCriteria() bool {
	cond := len(dnrh.pcl) >= dnrh.rs.GetCurrentThreshold()
	if dnrh.rs.OwnRoundState().NextRound != nil {
		if dnrh.rs.OwnRoundState().NRCurrent(dnrh.rcert) {
			return dnrh.rcert.RClaims.Round == constants.DEADBLOCKROUNDNR && cond
		}
	}
	return false
}

func (dnrh *doNrsCastNextHeightHandler) evalLogic() (bool, error) {
	return dnrh.external(dnrh.rs, dnrh.rcert, dnrh.pcl)
}

func (ce *Engine) doNrsCastNextHeightHandlerFunc(rs *RoundStates, rcert *objs.RCert, pcl objs.PreCommitList) (bool, error) {
	p, err := pcl.GetProposal()
	if err != nil {
		return false, err
	}
	errorFree := true
	if err := ce.updateValidValue(rs, p); err != nil {
		var e *errorz.ErrInvalid
		if err != errorz.ErrMissingTransactions && !errors.As(err, &e) {
			utils.DebugTrace(ce.logger, err)
			return false, err
		}
		errorFree = false
	}
	if errorFree {
		if err := ce.castNextHeightFromPreCommits(rs, pcl); err != nil {
			var e *errorz.ErrInvalid
			if err != errorz.ErrMissingTransactions && !errors.As(err, &e) {
				utils.DebugTrace(ce.logger, err)
				return false, err
			}
			errorFree = false
		}
	}
	if errorFree {
		return true, nil
	}

	// last of all count next round messages from this round only
	_, nrl, err := rs.GetCurrentNext()
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}

	// form a new round cert if we have enough
	if len(nrl) >= rs.GetCurrentThreshold() {
		if err := ce.castNextRoundRCert(rs, nrl); err != nil {
			utils.DebugTrace(ce.logger, err)
			return false, err
		}
	}
	// if we do not have enough yet,
	// do nothing and wait for more votes
	return true, nil
}

type doNrsCastNextRoundHandler struct {
	rs       *RoundStates
	rcert    *objs.RCert
	pcl      objs.PreCommitList
	nrl      objs.NextRoundList
	external func(*RoundStates, *objs.RCert, objs.NextRoundList) (bool, error)
}

func (dnrh *doNrsCastNextRoundHandler) evalCriteria() bool {
	cond1 := len(dnrh.pcl) < dnrh.rs.GetCurrentThreshold()
	cond2 := len(dnrh.nrl) >= dnrh.rs.GetCurrentThreshold()
	if cond1 && cond2 {
		if dnrh.rs.OwnRoundState().NextRound != nil {
			if dnrh.rs.OwnRoundState().NRCurrent(dnrh.rcert) {
				return dnrh.rcert.RClaims.Round == constants.DEADBLOCKROUNDNR
			}
		}
	}
	return false
}

func (dnrh *doNrsCastNextRoundHandler) evalLogic() (bool, error) {
	return dnrh.external(dnrh.rs, dnrh.rcert, dnrh.nrl)
}

func (ce *Engine) doNrsCastNextRoundHandlerFunc(rs *RoundStates, rcert *objs.RCert, nrl objs.NextRoundList) (bool, error) {
	if err := ce.castNextRoundRCert(rs, nrl); err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	return true, nil
}

type castNhHandler struct {
	rs        *RoundStates
	NHs       objs.NextHeightList
	NHCurrent bool
	external  func(*RoundStates, objs.NextHeightList) (bool, error)
}

func (cnhh *castNhHandler) evalCriteria() bool {
	return len(cnhh.NHs) > 0 && !cnhh.NHCurrent
}

func (cnhh *castNhHandler) evalLogic() (bool, error) {
	return cnhh.external(cnhh.rs, cnhh.NHs)
}

func (ce *Engine) castNhFunc(rs *RoundStates, NHs objs.NextHeightList) (bool, error) {
	// seems to be good already
	err := ce.castNextHeightFromNextHeight(rs, NHs[0])
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	return true, nil
}

type doNextHeightHandler struct {
	rs        *RoundStates
	NHs       objs.NextHeightList
	NHCurrent bool
	external  func(*RoundStates, objs.NextHeightList) (bool, error)
}

func (dnhh *doNextHeightHandler) evalCriteria() bool {
	return len(dnhh.NHs) > 0 && dnhh.NHCurrent
}

func (dnhh *doNextHeightHandler) evalLogic() (bool, error) {
	return dnhh.external(dnhh.rs, dnhh.NHs)
}

func (ce *Engine) doNextHeightFunc(rs *RoundStates, NHs objs.NextHeightList) (bool, error) {
	err := ce.doNextHeightStep(rs)
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	return true, nil
}

type castBlockHeaderHandler struct {
	rs        *RoundStates
	NHs       objs.NextHeightList
	NHCurrent bool
	external  func(*RoundStates, objs.NextHeightList) (bool, error)
}

func (dnhh *castBlockHeaderHandler) evalCriteria() bool {
	cond := len(dnhh.NHs) >= dnhh.rs.GetCurrentThreshold()
	return len(dnhh.NHs) > 0 && dnhh.NHCurrent && cond
}

func (dnhh *castBlockHeaderHandler) evalLogic() (bool, error) {
	return dnhh.external(dnhh.rs, dnhh.NHs)
}

func (ce *Engine) castBlockHeaderHandlerFunc(rs *RoundStates, nhl objs.NextHeightList) (bool, error) {
	if err := ce.castNewCommittedBlockHeader(rs, nhl); err != nil {
		utils.DebugTrace(ce.logger, err)
		var e *errorz.ErrInvalid
		if err != errorz.ErrMissingTransactions && !errors.As(err, &e) {
			return false, err
		}
	}
	return true, nil
}

type doRoundJumpHandler struct {
	rs       *RoundStates
	ChFr     []*objs.RoundState
	external func(*RoundStates, []*objs.RoundState) (bool, error)
}

func (drjh *doRoundJumpHandler) evalCriteria() bool {
	return len(drjh.ChFr) > 0
}

func (drjh *doRoundJumpHandler) evalLogic() (bool, error) {
	return drjh.external(drjh.rs, drjh.ChFr)
}

func (ce *Engine) doRoundJumpFunc(rs *RoundStates, ChFr []*objs.RoundState) (bool, error) {
	var maxRCert *objs.RCert
	var vroundState *objs.RoundState
	for i := 0; i < len(ChFr); i++ {
		vroundState = ChFr[i]
		if maxRCert == nil {
			maxRCert = vroundState.RCert
			continue
		}
		if vroundState.RCert.RClaims.Round > maxRCert.RClaims.Round {
			maxRCert = vroundState.RCert
		}
	}
	err := ce.doRoundJump(rs, maxRCert)
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	return true, nil
}

type doRoundJumpUpdateValidValueHandler struct {
	rs       *RoundStates
	ChFr     []*objs.RoundState
	pcl      objs.PreCommitList
	external func(*RoundStates, *objs.RCert, objs.PreCommitList) (bool, error)
}

func (drjh *doRoundJumpUpdateValidValueHandler) evalCriteria() bool {
	cond := len(drjh.pcl) > drjh.rs.GetCurrentThreshold()
	return len(drjh.ChFr) > 0 && cond
}

func (drjh *doRoundJumpUpdateValidValueHandler) evalLogic() (bool, error) {
	var maxRCert *objs.RCert
	var vroundState *objs.RoundState
	for i := 0; i < len(drjh.ChFr); i++ {
		vroundState = drjh.ChFr[i]
		if maxRCert == nil {
			maxRCert = vroundState.RCert
			continue
		}
		if vroundState.RCert.RClaims.Round > maxRCert.RClaims.Round {
			maxRCert = vroundState.RCert
		}
	}
	return drjh.external(drjh.rs, maxRCert, drjh.pcl)
}

type doRoundJumpSetRCertHandler struct {
	rs       *RoundStates
	ChFr     []*objs.RoundState
	pcl      objs.PreCommitList
	external func(*RoundStates, *objs.RCert, objs.PreCommitList) (bool, error)
}

func (drjh *doRoundJumpSetRCertHandler) evalCriteria() bool {
	cond := len(drjh.pcl) <= drjh.rs.GetCurrentThreshold()
	return len(drjh.ChFr) > 0 && cond
}

func (drjh *doRoundJumpSetRCertHandler) evalLogic() (bool, error) {
	var maxRCert *objs.RCert
	var vroundState *objs.RoundState
	for i := 0; i < len(drjh.ChFr); i++ {
		vroundState = drjh.ChFr[i]
		if maxRCert == nil {
			maxRCert = vroundState.RCert
			continue
		}
		if vroundState.RCert.RClaims.Round > maxRCert.RClaims.Round {
			maxRCert = vroundState.RCert
		}
	}
	return drjh.external(drjh.rs, maxRCert, drjh.pcl)
}

type nrCurrentHandler struct {
	rs        *RoundStates
	NRCurrent bool
	external  func(*RoundStates) (bool, error)
}

func (nrch *nrCurrentHandler) evalCriteria() bool {
	return nrch.NRCurrent
}

func (nrch *nrCurrentHandler) evalLogic() (bool, error) {
	return nrch.external(nrch.rs)
}

func (ce *Engine) nrCurrentFunc(rs *RoundStates) (bool, error) {
	// seems to have two sub handlers
	err := ce.doNextRoundStep(rs)
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	return true, nil
}

type nrCurrentCnhHandler struct {
	rs        *RoundStates
	NRCurrent bool
	rcert     *objs.RCert
	pcl       objs.PreCommitList
	external  func(*RoundStates, *objs.RCert, objs.PreCommitList) (bool, error)
}

func (nrch *nrCurrentCnhHandler) evalCriteria() bool {
	cond := len(nrch.pcl) >= nrch.rs.GetCurrentThreshold()
	return nrch.NRCurrent && cond
}

func (nrch *nrCurrentCnhHandler) evalLogic() (bool, error) {
	return nrch.external(nrch.rs, nrch.rcert, nrch.pcl)
}

type nrCurrentCnrHandler struct {
	rs        *RoundStates
	NRCurrent bool
	rcert     *objs.RCert
	pcl       objs.PreCommitList
	nrl       objs.NextRoundList
	external  func(*RoundStates, *objs.RCert, objs.NextRoundList) (bool, error)
}

func (nrch *nrCurrentCnrHandler) evalCriteria() bool {
	cond1 := len(nrch.pcl) < nrch.rs.GetCurrentThreshold()
	cond2 := len(nrch.nrl) >= nrch.rs.GetCurrentThreshold()
	return nrch.NRCurrent && cond1 && cond2
}

func (nrch *nrCurrentCnrHandler) evalLogic() (bool, error) {
	return nrch.external(nrch.rs, nrch.rcert, nrch.nrl)
}

type pcCurrentHandler struct {
	rs        *RoundStates
	PCCurrent bool
	external  func(*RoundStates, bool) (bool, error)
}

func (pcch *pcCurrentHandler) evalCriteria() bool {
	return pcch.PCCurrent
}

func (pcch *pcCurrentHandler) evalLogic() (bool, error) {
	PCTOExpired := pcch.rs.OwnValidatingState.PCTOExpired()
	return pcch.external(pcch.rs, PCTOExpired)
}

func (ce *Engine) pcCurrentFunc(rs *RoundStates, PCTOExpired bool) (bool, error) {
	if PCTOExpired {
		err := ce.doPendingNext(rs)
		if err != nil {
			utils.DebugTrace(ce.logger, err)
			return false, err
		}
		return true, nil
	}
	err := ce.doPreCommitStep(rs)
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	return true, nil
}

type pcnCurrentHandler struct {
	rs         *RoundStates
	PCNCurrent bool
	external   func(*RoundStates, bool) (bool, error)
}

func (pcnch *pcnCurrentHandler) evalCriteria() bool {
	return pcnch.PCNCurrent
}

func (pcnch *pcnCurrentHandler) evalLogic() (bool, error) {
	PCTOExpired := pcnch.rs.OwnValidatingState.PCTOExpired()
	return pcnch.external(pcnch.rs, PCTOExpired)
}

func (ce *Engine) pcnCurrentFunc(rs *RoundStates, PCTOExpired bool) (bool, error) {
	if PCTOExpired {
		err := ce.doPendingNext(rs)
		if err != nil {
			utils.DebugTrace(ce.logger, err)
			return false, err
		}
		return true, nil
	}
	err := ce.doPreCommitNilStep(rs)
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	return true, nil
}

type pvCurrentHandler struct {
	rs        *RoundStates
	PVCurrent bool
	external  func(*RoundStates, bool) (bool, error)
}

func (pvch *pvCurrentHandler) evalCriteria() bool {
	return pvch.PVCurrent
}

func (pvch *pvCurrentHandler) evalLogic() (bool, error) {
	PVTOExpired := pvch.rs.OwnValidatingState.PVTOExpired()
	return pvch.external(pvch.rs, PVTOExpired)
}

func (ce *Engine) pvCurrentFunc(rs *RoundStates, PVTOExpired bool) (bool, error) {
	if PVTOExpired {
		err := ce.doPendingPreCommit(rs)
		if err != nil {
			utils.DebugTrace(ce.logger, err)
			return false, err
		}
		return true, nil
	}
	err := ce.doPreVoteStep(rs)
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	return true, nil
}

type pvnCurrentHandler struct {
	rs         *RoundStates
	PVNCurrent bool
	external   func(*RoundStates, bool) (bool, error)
}

func (pvnch *pvnCurrentHandler) evalCriteria() bool {
	return pvnch.PVNCurrent
}

func (pvnch *pvnCurrentHandler) evalLogic() (bool, error) {
	PVTOExpired := pvnch.rs.OwnValidatingState.PVTOExpired()
	return pvnch.external(pvnch.rs, PVTOExpired)
}

func (ce *Engine) pvnCurrentFunc(rs *RoundStates, PVTOExpired bool) (bool, error) {
	if PVTOExpired {
		err := ce.doPendingPreCommit(rs)
		if err != nil {
			utils.DebugTrace(ce.logger, err)
			return false, err
		}
		return true, nil
	}
	err := ce.doPreVoteNilStep(rs)
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	return true, nil
}

type ptoExpiredHandler struct {
	rs       *RoundStates
	external func(*RoundStates) (bool, error)
}

func (ptoeh *ptoExpiredHandler) evalCriteria() bool {
	PTOExpired := ptoeh.rs.OwnValidatingState.PTOExpired()
	return PTOExpired
}

func (ptoeh *ptoExpiredHandler) evalLogic() (bool, error) {
	return ptoeh.external(ptoeh.rs)
}

func (ce *Engine) ptoExpiredFunc(rs *RoundStates) (bool, error) {
	err := ce.doPendingPreVoteStep(rs)
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	return true, nil
}

type validPropHandler struct {
	rs       *RoundStates
	PCurrent bool
	external func(*RoundStates) (bool, error)
}

func (vph *validPropHandler) evalCriteria() bool {
	IsProposer := vph.rs.LocalIsProposer()
	return (IsProposer && !vph.PCurrent && vph.rs.OwnRoundState().RCert.RClaims.Round < constants.DEADBLOCKROUND)
}

func (vph *validPropHandler) evalLogic() (bool, error) {
	return vph.external(vph.rs)
}

func (ce *Engine) validPropFunc(rs *RoundStates) (bool, error) {
	err := ce.doPendingProposalStep(rs)
	if err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	return true, nil
}

// Sync .
func (ce *Engine) Sync() (bool, error) {
	// see if sync is done
	// if yes exit
	syncDone := false
	// ce.logger.Error("!!! OPEN SYNC TXN")
	// defer func() { ce.logger.Error("!!! CLOSE SYNC TXN") }()
	err := ce.database.Update(func(txn *badger.Txn) error {
		rs, err := ce.sstore.LoadLocalState(txn)
		if err != nil {
			utils.DebugTrace(ce.logger, err)
			return err
		}
		rs.txn = txn
		// begin handling logic
		if rs.OwnState.MaxBHSeen.BClaims.Height == rs.OwnState.SyncToBH.BClaims.Height {
			syncDone = true
			return nil
		}
		if rs.OwnState.MaxBHSeen.BClaims.Height > constants.EpochLength*2 {
			if rs.OwnState.SyncToBH.BClaims.Height <= rs.OwnState.MaxBHSeen.BClaims.Height-constants.EpochLength*2 {
				// Guard against the short first epoch causing errors in the sync logic
				// by escaping early and just waiting for the MaxBHSeen to increase.
				if rs.OwnState.MaxBHSeen.BClaims.Height%constants.EpochLength == 0 {
					return nil
				}
				epochOfMaxBHSeen := utils.Epoch(rs.OwnState.MaxBHSeen.BClaims.Height)
				canonicalEpoch := epochOfMaxBHSeen - 2
				canonicalSnapShotHeight := canonicalEpoch * constants.EpochLength
				mrcbh, err := ce.database.GetMostRecentCommittedBlockHeaderFastSync(txn)
				if err != nil {
					utils.DebugTrace(ce.logger, err)
					return err
				}
				csbh, err := ce.database.GetSnapshotBlockHeader(txn, canonicalSnapShotHeight)
				if err != nil {
					utils.DebugTrace(ce.logger, err)
					return err
				}
				canonicalBlockHash, err := csbh.BlockHash()
				if err != nil {
					utils.DebugTrace(ce.logger, err)
					return err
				}
				fastSyncDone, err := ce.fastSync.Update(txn, csbh.BClaims.Height, mrcbh.BClaims.Height, csbh.BClaims.StateRoot, csbh.BClaims.HeaderRoot, canonicalBlockHash)
				if err != nil {
					utils.DebugTrace(ce.logger, err)
					return err
				}
				if fastSyncDone {
					if err := ce.setMostRecentBlockHeaderFastSync(rs, csbh); err != nil {
						utils.DebugTrace(ce.logger, err)
						return err
					}
					err = ce.sstore.WriteState(rs)
					if err != nil {
						utils.DebugTrace(ce.logger, err)
						return err
					}
					return nil
				}
				return nil
			}
		}
		ce.logger.Debugf("SyncOneBH:  MBHS:%v  STBH:%v", rs.OwnState.MaxBHSeen.BClaims.Height, rs.OwnState.SyncToBH.BClaims.Height)
		txs, bh, err := ce.dm.SyncOneBH(txn, rs)
		if err != nil {
			utils.DebugTrace(ce.logger, err)
			return err
		}
		ok, err := ce.isValid(rs, bh.BClaims.ChainID, bh.BClaims.StateRoot, bh.BClaims.HeaderRoot, txs)
		if err != nil {
			utils.DebugTrace(ce.logger, err)
			return err
		}
		if !ok {
			return nil
		}
		err = ce.setMostRecentBlockHeader(rs, bh)
		if err != nil {
			utils.DebugTrace(ce.logger, err)
			return err
		}
		err = ce.sstore.WriteState(rs)
		if err != nil {
			utils.DebugTrace(ce.logger, err)
			return err
		}
		return nil
	})
	if err != nil {
		e := errorz.ErrInvalid{}.New("")
		if errors.As(err, &e) {
			utils.DebugTrace(ce.logger, err)
			return false, nil
		}
		if err == errorz.ErrMissingTransactions {
			utils.DebugTrace(ce.logger, err)
			return false, nil
		}
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	if err := ce.database.Sync(); err != nil {
		utils.DebugTrace(ce.logger, err)
		return false, err
	}
	return syncDone, nil
}

func (ce *Engine) loadValidationKey(rs *RoundStates) error {
	if rs.IsCurrentValidator() {
		if !bytes.Equal(rs.ValidatorSet.GroupKey, rs.OwnValidatingState.GroupKey) || ce.bnSigner == nil {
			for i := 0; i < len(rs.ValidatorSet.Validators); i++ {
				v := rs.ValidatorSet.Validators[i]
				if bytes.Equal(v.VAddr, rs.OwnState.VAddr) {
					name := make([]byte, len(v.GroupShare))
					copy(name[:], v.GroupShare)
					pk, err := ce.AdminBus.GetPrivK(name)
					if err != nil {
						utils.DebugTrace(ce.logger, err)
						return nil // TODO: are we supposed to swallow this error?
					}
					signer := &crypto.BNGroupSigner{}
					signer.SetPrivk(pk)
					err = signer.SetGroupPubk(rs.ValidatorSet.GroupKey)
					if err != nil {
						utils.DebugTrace(ce.logger, err)
						return err
					}
					ce.bnSigner = signer
					pubk, err := ce.bnSigner.PubkeyShare()
					if err != nil {
						return err
					}
					if !bytes.Equal(name, pubk) {
						utils.DebugTrace(ce.logger, nil, "name and public key do not match")
						return err // TODO: err == nil; should return an errorz.ErrInvalid?;
					}
					break
				}
			}
			rs.OwnValidatingState.GroupKey = rs.ValidatorSet.GroupKey
		}
	}
	return nil
}
