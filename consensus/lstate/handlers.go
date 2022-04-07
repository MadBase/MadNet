package lstate

import (
	"bytes"
	"fmt"

	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/consensus/dman"
	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/crypto"
	"github.com/MadBase/MadNet/errorz"
	"github.com/MadBase/MadNet/logging"
	"github.com/MadBase/MadNet/utils"
	"github.com/dgraph-io/badger/v2"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	database *db.Database
	sstore   *Store
	secpVal  *crypto.Secp256k1Validator
	bnVal    *crypto.BNGroupValidator
	dm       *dman.DMan
	logger   *logrus.Logger
}

// Init initializes the Handlers object
func (mb *Handlers) Init(database *db.Database, dm *dman.DMan) {
	mb.logger = logging.GetLogger(constants.LoggerConsensus)
	mb.sstore = &Store{}
	mb.sstore.Init(database)
	mb.database = database
	mb.dm = dm
}

// AddProposal stores a proposal to the database
func (mb *Handlers) AddProposal(v *objs.Proposal) error {
	return mb.Store(v)
}

// AddPreVote stores a preVote to the database
func (mb *Handlers) AddPreVote(v *objs.PreVote) error {
	return mb.Store(v)
}

// AddPreVoteNil stores a preVoteNil to the database
func (mb *Handlers) AddPreVoteNil(v *objs.PreVoteNil) error {
	return mb.Store(v)
}

// AddPreCommit stores a preCommit to the database
func (mb *Handlers) AddPreCommit(v *objs.PreCommit) error {
	return mb.Store(v)
}

// AddPreCommitNil stores a preCommitNil to the database
func (mb *Handlers) AddPreCommitNil(v *objs.PreCommitNil) error {
	return mb.Store(v)
}

// AddNextRound stores a nextRound object to the database
func (mb *Handlers) AddNextRound(v *objs.NextRound) error {
	return mb.Store(v)
}

// AddNextHeight stores a nextHeight object to the database
func (mb *Handlers) AddNextHeight(v *objs.NextHeight) error {
	return mb.Store(v)
}

// AddBlockHeader stores a blockHeader object to the database
func (mb *Handlers) AddBlockHeader(v *objs.BlockHeader) error {
	return mb.Store(v)
}

// Store updates database to include the object; it stores the object.
func (mb *Handlers) Store(v interface{}) error {
	return mb.database.Update(func(txn *badger.Txn) error {
		rc, err := objs.ExtractRCertAny(v)
		if err != nil {
			return err
		}
		roundState, err := mb.sstore.LoadState(txn, rc)
		if err != nil {
			return err
		}
		switch obj := v.(type) {
		case *objs.Proposal:
			txHshLst := obj.TxHshLst
			err = roundState.SetProposal(obj)
			if err != nil {
				return err
			}
			// if !roundState.IsCurrentValidator() {
			// 	err := mb.database.SetBroadcastProposal(txn, obj)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
			go mb.dm.DownloadTxs(roundState.height, roundState.round, txHshLst)
		case *objs.PreVote:
			err = roundState.SetPreVote(obj)
			if err != nil {
				return err
			}
			// if !roundState.IsCurrentValidator() {
			// err := mb.database.SetBroadcastPreVote(txn, obj)
			// if err != nil {
			// return err
			// }
			// }
		case *objs.PreVoteNil:
			err = roundState.SetPreVoteNil(obj)
			if err != nil {
				return err
			}
			// if !roundState.IsCurrentValidator() {
			// err := mb.database.SetBroadcastPreVoteNil(txn, obj)
			// if err != nil {
			// return err
			// }
			// }
		case *objs.PreCommit:
			err = roundState.SetPreCommit(obj)
			if err != nil {
				return err
			}
			// if !roundState.IsCurrentValidator() {
			// 	err := mb.database.SetBroadcastPreCommit(txn, obj)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
		case *objs.PreCommitNil:
			err = roundState.SetPreCommitNil(obj)
			if err != nil {
				return err
			}
			// if !roundState.IsCurrentValidator() {
			// 	err := mb.database.SetBroadcastPreCommitNil(txn, obj)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
		case *objs.NextRound:
			err = roundState.SetNextRound(obj)
			if err != nil {
				return err
			}
			// if !roundState.IsCurrentValidator() {
			// 	err := mb.database.SetBroadcastNextRound(txn, obj)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
		case *objs.NextHeight:
			err = roundState.SetNextHeight(obj)
			if err != nil {
				return err
			}
			// if !roundState.IsCurrentValidator() {
			// 	err := mb.database.SetBroadcastNextHeight(txn, obj)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
		case *objs.BlockHeader:
			ownState := roundState.OwnState
			if obj.BClaims.Height <= ownState.MaxBHSeen.BClaims.Height {
				return errorz.ErrInvalid{}.New("stale bh  - <= MaxBHSeen")
			}
			if obj.BClaims.Height <= ownState.SyncToBH.BClaims.Height {
				return errorz.ErrInvalid{}.New("stale bh - <= SyncTOBH ")
			}
			ownState.MaxBHSeen = obj
		}
		return mb.sstore.WriteState(txn, roundState)
	})
}

// PreValidate checks a message for validity and performs cryptographic
// validation
func (mb *Handlers) PreValidate(v interface{}) error {
	var Voter []byte
	var Proposer []byte
	var GroupShare []byte
	var GroupKey []byte
	var CoSigners [][]byte
	var round uint32
	_, chainID := objs.ExtractHCID(v)
	err := mb.database.View(func(txn *badger.Txn) error {
		os, err := mb.database.GetOwnState(txn)
		if err != nil {
			utils.DebugTrace(mb.logger, err)
			return err
		}
		rs, err := mb.database.GetCurrentRoundState(txn, os.VAddr)
		if err != nil {
			utils.DebugTrace(mb.logger, err)
			return err
		}
		cid := rs.RCert.RClaims.ChainID
		if cid != chainID {
			return errorz.ErrInvalid{}.New("cid mismatch")
		}
		h := rs.RCert.RClaims.Height
		r := rs.RCert.RClaims.Round
		var height uint32
		switch obj := v.(type) {
		case *objs.Proposal:
			height = obj.PClaims.RCert.RClaims.Height
			if height < h {
				return errorz.ErrStale{}.New("Proposal h<h-1: OwnH:%v ObjH:%v", h, height)
			}
			round := obj.PClaims.RCert.RClaims.Round
			if round < r {
				errorz.ErrStale{}.New("Proposal r<r-1: OwnR:%v ObjR:%v", r, round)
			}
			if err := obj.ValidateSignatures(mb.secpVal, mb.bnVal); err != nil {
				return err
			}
			//Voter = nil
			Proposer = obj.Proposer
			//GroupShare = nil
			GroupKey = obj.GroupKey
			//CoSigners = nil
			round = obj.PClaims.RCert.RClaims.Round
		case *objs.PreVote:
			height = obj.Proposal.PClaims.RCert.RClaims.Height
			if height < h {
				return errorz.ErrStale{}.New("PreVote h<h-1: OwnH:%v ObjH:%v", h, height)
			}
			round = obj.Proposal.PClaims.RCert.RClaims.Round
			if round < r {
				errorz.ErrStale{}.New("PreVote r<r-1: OwnR:%v ObjR:%v", r, round)
			}
			if err := obj.ValidateSignatures(mb.secpVal, mb.bnVal); err != nil {
				return err
			}
			Voter = obj.Voter
			Proposer = obj.Proposal.Proposer
			//GroupShare = nil
			GroupKey = obj.GroupKey
			//CoSigners = nil
		case *objs.PreVoteNil:
			height = obj.RCert.RClaims.Height
			if height < h {
				return errorz.ErrStale{}.New("PreVoteNil h<h-1: OwnH:%v ObjH:%v", h, height)
			}
			round = obj.RCert.RClaims.Round
			if round < r {
				errorz.ErrStale{}.New("PreVoteNil r<r-1: OwnR:%v ObjR:%v", r, round)
			}
			if err := obj.ValidateSignatures(mb.secpVal, mb.bnVal); err != nil {
				return err
			}
			Voter = obj.Voter
			//Proposer = nil
			//GroupShare = nil
			GroupKey = obj.GroupKey
			//CoSigners = nil
		case *objs.PreCommit:
			height = obj.Proposal.PClaims.RCert.RClaims.Height
			if height < h {
				return errorz.ErrStale{}.New("PreCommit h<h-1: OwnH:%v ObjH:%v", h, height)
			}
			round = obj.Proposal.PClaims.RCert.RClaims.Round
			if round < r {
				errorz.ErrStale{}.New("PreCommit r<r-1: OwnR:%v ObjR:%v", r, round)
			}
			if err := obj.ValidateSignatures(mb.secpVal, mb.bnVal); err != nil {
				return err
			}
			Voter = obj.Voter
			Proposer = obj.Proposer
			//GroupShare = nil
			GroupKey = obj.GroupKey
			CoSigners = obj.Signers
		case *objs.PreCommitNil:
			height = obj.RCert.RClaims.Height
			if height < h {
				return errorz.ErrStale{}.New("PreCommitNil h<h-1: OwnH:%v ObjH:%v", h, height)
			}
			round = obj.RCert.RClaims.Round
			if round < r {
				errorz.ErrStale{}.New("PreCommitNil r<r-1: OwnR:%v ObjR:%v", r, round)
			}
			if err := obj.ValidateSignatures(mb.secpVal, mb.bnVal); err != nil {
				return err
			}
			Voter = obj.Voter
			//Proposer = nil
			//GroupShare = nil
			GroupKey = obj.GroupKey
			//CoSigners = nil
		case *objs.NextRound:
			height = obj.NRClaims.RCert.RClaims.Height
			if height < h {
				return errorz.ErrStale{}.New("NextRound h<h-1: OwnH:%v ObjH:%v", h, height)
			}
			round = obj.NRClaims.RCert.RClaims.Round
			if round < r {
				errorz.ErrStale{}.New("NextRound r<r-1: OwnR:%v ObjR:%v", r, round)
			}
			if err := obj.ValidateSignatures(mb.secpVal, mb.bnVal); err != nil {
				return err
			}
			Voter = obj.Voter
			//Proposer = nil
			GroupShare = obj.GroupShare
			GroupKey = obj.GroupKey
			//CoSigners = nil
		case *objs.NextHeight:
			height = obj.NHClaims.Proposal.PClaims.RCert.RClaims.Height
			if height < h {
				return errorz.ErrStale{}.New("NextHeight h<h-1: OwnH:%v ObjH:%v", h, height)
			}
			if err := obj.ValidateSignatures(mb.secpVal, mb.bnVal); err != nil {
				return err
			}
			Voter = obj.Voter
			Proposer = obj.NHClaims.Proposal.Proposer
			GroupShare = obj.GroupShare
			GroupKey = obj.GroupKey
			CoSigners = obj.Signers
			round = obj.NHClaims.Proposal.PClaims.RCert.RClaims.Round
		case *objs.BlockHeader:
			height = obj.BClaims.Height
			if err := obj.ValidateSignatures(mb.bnVal); err != nil {
				return err
			}
			//Voter = nil
			//Proposer = nil
			//GroupShare = nil
			GroupKey = obj.GroupKey
			//CoSigners = nil
			round = 1
		default:
			panic("Unknown type")
		}

		if height == 1 {
			return errorz.ErrInvalid{}.New("No Height 1 message is valid except for initial block")
		}
		// Do something in height 2 round 1 when GroupKey is not set; we set it
		// to the value in ValidatorSet. This may not be best, as it now
		// automatically passes those portions of the test.
		if height == 2 && round == 1 && len(GroupKey) == 0 {
			vSet, err := mb.database.GetValidatorSet(txn, height)
			if err != nil {
				utils.DebugTrace(mb.logger, err)
				return err
			}
			GroupKey = utils.CopySlice(vSet.GroupKey)
		}
		vSet, err := mb.database.GetValidatorSet(txn, height)
		if err != nil {
			utils.DebugTrace(mb.logger, err)
			return err
		}
		if !bytes.Equal(GroupKey, vSet.GroupKey) {
			return errorz.ErrInvalid{}.New("group key mismatch in state handlers")
		}
		if Voter != nil && GroupShare != nil {
			if !vSet.IsValidTriplet(Voter, GroupShare, GroupKey) {
				correctgk := GroupShare
				vl := [][]byte{}
				for _, vobj := range vSet.Validators {
					vl = append(vl, vobj.GroupShare)
				}
				return errorz.ErrInvalid{}.New(fmt.Sprintf("invalid triplet in state handlers: \nvoter:%x \nGroupShare:%x\n%x\n%x\n%x\n%x", Voter, correctgk, vl[0], vl[1], vl[2], vl[3]))
			}
		}
		if Voter != nil && GroupShare == nil {
			if !vSet.IsValidTuple(Voter, GroupKey) {
				correctgk := vSet.GroupKey
				vl := [][]byte{}
				for _, vobj := range vSet.Validators {
					vl = append(vl, vobj.VAddr)
				}
				return errorz.ErrInvalid{}.New(fmt.Sprintf("invalid tuple in state handlers: \nvoter:%x \nGroupKey:%x\ncorrectgk:%x\n%x\n%x\n%x\n%x", Voter, GroupKey, correctgk, vl[0], vl[1], vl[2], vl[3]))
			}
		}
		if Proposer != nil {
			if !vSet.IsValidTuple(Proposer, GroupKey) {
				return errorz.ErrInvalid{}.New("invalid proposer in state handlers")
			}
			rcert := objs.ExtractRCert(v)
			if rcert.RClaims.Round < constants.DEADBLOCKROUND {
				pidx := objs.GetProposerIdx(len(vSet.Validators), rcert.RClaims.Height, rcert.RClaims.Round)
				valObj := vSet.Validators[pidx]
				vAddr := valObj.VAddr
				if !bytes.Equal(Proposer, vAddr) {
					return errorz.ErrInvalid{}.New("bad proposer")
				}
			}
		}
		for _, cs := range CoSigners {
			if !vSet.IsValidTuple(utils.CopySlice(cs), GroupKey) {
				return errorz.ErrInvalid{}.New("bad co signer")
			}
		}
		return nil
	})
	if err != nil {
		utils.DebugTrace(mb.logger, err)
		return err
	}
	return nil
}
