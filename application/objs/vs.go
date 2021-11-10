package objs

import (
	mdefs "github.com/MadBase/MadNet/application/objs/capn"
	"github.com/MadBase/MadNet/application/objs/uint256"
	"github.com/MadBase/MadNet/application/objs/valuestore"
	"github.com/MadBase/MadNet/application/wrapper"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/errorz"
	"github.com/MadBase/MadNet/utils"
	capnp "zombiezen.com/go/capnproto2"
)

// ValueStore stores value in a UTXO
type ValueStore struct {
	VSPreImage *VSPreImage
	TxHash     []byte
	//
	utxoID []byte
}

// New creates a new ValueStore
func (b *ValueStore) New(chainID uint32, value *uint256.Uint256, fee *uint256.Uint256, acct []byte, curveSpec constants.CurveSpec, txHash []byte) error {
	if b == nil {
		return errorz.ErrInvalid{}.New("not initialized")
	}
	if value == nil || value.IsZero() {
		return errorz.ErrInvalid{}.New("invalue value: nil or zero")
	}
	if fee == nil {
		return errorz.ErrInvalid{}.New("invalue fee: nil")
	}
	vsowner := &ValueStoreOwner{}
	vsowner.New(acct, curveSpec)
	// if err := vsowner.Validate(); err != nil {
	// 	return err
	// }
	if chainID == 0 {
		return errorz.ErrInvalid{}.New("Error in ValueStore.New: invalid chainID")
	}
	if len(txHash) != constants.HashLen {
		return errorz.ErrInvalid{}.New("Error in ValueStore.New: invalid txHash")
	}
	vsp := &VSPreImage{
		ChainID:  chainID,
		Value:    value.Clone(),
		TXOutIdx: constants.MaxUint32,
		Owner:    vsowner,
		Fee:      fee.Clone(),
	}
	b.VSPreImage = vsp
	b.TxHash = utils.CopySlice(txHash)
	return nil
}

// NewFromDeposit creates a new ValueStore from a deposit event
func (b *ValueStore) NewFromDeposit(chainID uint32, value *uint256.Uint256, acct []byte, nonce []byte) error {
	vsowner := &ValueStoreOwner{}
	vsowner.New(acct, constants.CurveSecp256k1)
	// if err := vsowner.Validate(); err != nil {
	// 	return err
	// }
	if chainID == 0 {
		return errorz.ErrInvalid{}.New("Error in ValueStore.NewFromDeposit: invalid chainID")
	}
	if len(nonce) != constants.HashLen {
		return errorz.ErrInvalid{}.New("Error in ValueStore.NewFromDeposit: invalid nonce")
	}
	vsp := &VSPreImage{
		ChainID:  chainID,
		Value:    value,
		TXOutIdx: constants.MaxUint32,
		Owner:    vsowner,
	}
	b.VSPreImage = vsp
	b.TxHash = utils.CopySlice(nonce)
	return nil
}

// UnmarshalBinary takes a byte slice and returns the corresponding
// ValueStore object
func (b *ValueStore) UnmarshalBinary(data []byte) error {
	if b == nil {
		return errorz.ErrInvalid{}.New("not initialized")
	}
	bc, err := valuestore.Unmarshal(data)
	if err != nil {
		return err
	}
	return b.UnmarshalCapn(bc)
}

// MarshalBinary takes the ValueStore object and returns the canonical
// byte slice
func (b *ValueStore) MarshalBinary() ([]byte, error) {
	if b == nil {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	bc, err := b.MarshalCapn(nil)
	if err != nil {
		return nil, err
	}
	return valuestore.Marshal(bc)
}

// UnmarshalCapn unmarshals the capnproto definition of the object
func (b *ValueStore) UnmarshalCapn(bc mdefs.ValueStore) error {
	if err := valuestore.Validate(bc); err != nil {
		return err
	}
	b.VSPreImage = &VSPreImage{}
	if err := b.VSPreImage.UnmarshalCapn(bc.VSPreImage()); err != nil {
		return err
	}
	b.TxHash = utils.CopySlice(bc.TxHash())
	return nil
}

// MarshalCapn marshals the object into its capnproto definition
func (b *ValueStore) MarshalCapn(seg *capnp.Segment) (mdefs.ValueStore, error) {
	if b == nil {
		return mdefs.ValueStore{}, errorz.ErrInvalid{}.New("not initialized")
	}
	var bc mdefs.ValueStore
	if seg == nil {
		_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		if err != nil {
			return bc, err
		}
		tmp, err := mdefs.NewRootValueStore(seg)
		if err != nil {
			return bc, err
		}
		bc = tmp
	} else {
		tmp, err := mdefs.NewValueStore(seg)
		if err != nil {
			return bc, err
		}
		bc = tmp
	}
	seg = bc.Struct.Segment()
	bt, err := b.VSPreImage.MarshalCapn(seg)
	if err != nil {
		return bc, err
	}
	if err := bc.SetVSPreImage(bt); err != nil {
		return bc, err
	}
	if err := bc.SetTxHash(utils.CopySlice(b.TxHash)); err != nil {
		return bc, err
	}
	return bc, nil
}

// PreHash calculates the PreHash of the object
func (b *ValueStore) PreHash() ([]byte, error) {
	if b == nil {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	return b.VSPreImage.PreHash()
}

// UTXOID calculates the UTXOID of the object
func (b *ValueStore) UTXOID() ([]byte, error) {
	if b == nil || b.VSPreImage == nil || len(b.TxHash) != constants.HashLen {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	if b.utxoID != nil {
		return utils.CopySlice(b.utxoID), nil
	}
	b.utxoID = MakeUTXOID(b.TxHash, b.VSPreImage.TXOutIdx)
	return utils.CopySlice(b.utxoID), nil
}

// TXOutIdx returns the TXOutIdx of the object
func (b *ValueStore) TXOutIdx() (uint32, error) {
	if b == nil || b.VSPreImage == nil {
		return 0, errorz.ErrInvalid{}.New("not initialized")
	}
	return b.VSPreImage.TXOutIdx, nil
}

// SetTXOutIdx sets the TXOutIdx of the object
func (b *ValueStore) SetTXOutIdx(idx uint32) error {
	if b == nil || b.VSPreImage == nil {
		return errorz.ErrInvalid{}.New("not initialized")
	}
	b.VSPreImage.TXOutIdx = idx
	return nil
}

// SetTxHash sets the TxHash of the object
func (b *ValueStore) SetTxHash(txHash []byte) error {
	if b == nil || b.VSPreImage == nil {
		return errorz.ErrInvalid{}.New("not initialized")
	}
	if len(txHash) != constants.HashLen {
		return errorz.ErrInvalid{}.New("Invalid hash length")
	}
	b.TxHash = utils.CopySlice(txHash)
	return nil
}

// ChainID returns the ChainID of the object
func (b *ValueStore) ChainID() (uint32, error) {
	if b == nil || b.VSPreImage == nil || b.VSPreImage.ChainID == 0 {
		return 0, errorz.ErrInvalid{}.New("not initialized")
	}
	return b.VSPreImage.ChainID, nil
}

// Value returns the Value of the object
func (b *ValueStore) Value() (*uint256.Uint256, error) {
	if b == nil || b.VSPreImage == nil || b.VSPreImage.Value == nil || b.VSPreImage.Value.IsZero() {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	return b.VSPreImage.Value.Clone(), nil
}

// Fee returns the Fee of the object
func (b *ValueStore) Fee() (*uint256.Uint256, error) {
	if b == nil || b.VSPreImage == nil || b.VSPreImage.Fee == nil {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	return b.VSPreImage.Fee.Clone(), nil
}

// ValuePlusFee returns the Value of the object with the associated fee
func (b *ValueStore) ValuePlusFee() (*uint256.Uint256, error) {
	value, err := b.Value()
	if err != nil {
		return nil, err
	}
	fee, err := b.Fee()
	if err != nil {
		return nil, err
	}
	total, err := new(uint256.Uint256).Add(value, fee)
	if err != nil {
		return nil, err
	}
	return total, nil
}

// IsDeposit returns true if the object is a deposit
func (b *ValueStore) IsDeposit() bool {
	if b == nil || b.VSPreImage == nil {
		return false
	}
	return b.VSPreImage.TXOutIdx == constants.MaxUint32
}

// Owner returns the ValueStoreOwner of the ValueStore
func (b *ValueStore) Owner() (*ValueStoreOwner, error) {
	if b == nil || b.VSPreImage == nil {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	// if err := b.VSPreImage.Owner.Validate(); err != nil {
	// 	return nil, errorz.ErrInvalid{}.New("ValueStoreOwner invalid")
	// }
	return b.VSPreImage.Owner, nil
}

// GenericOwner returns the Owner of the ValueStore
func (b *ValueStore) GenericOwner() (*Owner, error) {
	vso, err := b.Owner()
	if err != nil {
		return nil, err
	}
	onr := &Owner{}
	if err := onr.NewFromValueStoreOwner(vso); err != nil {
		return nil, err
	}
	return onr, nil
}

// Sign generates the signature for a ValueStore at the time of consumption
// func (b *ValueStore) Sign(txIn *TXIn, s Signer) error {
// msg, err := txIn.TXInLinker.MarshalBinary()
// if err != nil {
// return err
// }
// owner, err := b.Owner()
// if err != nil {
// return err
// }
// sig, err := owner.Sign(msg, s)
// if err != nil {
// return err
// }
// sigb, err := sig.MarshalBinary()
// if err != nil {
// return err
// }
// txIn.Signature = sigb
// return nil
// }

// ValidateFee validates the fee of the object at the time of creation
func (b *ValueStore) ValidateFee(storage *wrapper.Storage) error {
	fee, err := b.Fee()
	if err != nil {
		return err
	}
	if b.IsDeposit() {
		if !fee.IsZero() {
			return errorz.ErrInvalid{}.New("vs: invalid fee; deposits should have fee equal zero")
		}
		return nil
	}
	feeTrue, err := storage.GetValueStoreFee()
	if err != nil {
		return err
	}
	if fee.Cmp(feeTrue) != 0 {
		return errorz.ErrInvalid{}.New("vs: invalid fee")
	}
	return nil
}

// ValidateSignature validates the signature of the ValueStore at the time of
// consumption
// func (b *ValueStore) ValidateSignature(txIn *TXIn) error {
// if b == nil {
// return errorz.ErrInvalid{}.New("not initialized")
// }
// msg, err := txIn.TXInLinker.MarshalBinary()
// if err != nil {
// return err
// }
// sig := &ValueStoreSignature{}
// if err := sig.UnmarshalBinary(txIn.Signature); err != nil {
// return err
// }
// return b.VSPreImage.ValidateSignature(msg, sig)
// }

// MakeTxIn constructs a TXIn object for the current object
func (b *ValueStore) MakeTxIn() (*TXIn, error) {
	txOutIdx, err := b.TXOutIdx()
	if err != nil {
		return nil, err
	}
	cid, err := b.ChainID()
	if err != nil {
		return nil, err
	}
	if len(b.TxHash) != constants.HashLen {
		return nil, errorz.ErrInvalid{}.New("invalid TxHash")
	}
	return &TXIn{
		TXInLinker: &TXInLinker{
			TXInPreImage: &TXInPreImage{
				ConsumedTxIdx:  txOutIdx,
				ConsumedTxHash: utils.CopySlice(b.TxHash),
				ChainID:        cid,
			},
		},
	}, nil
}
