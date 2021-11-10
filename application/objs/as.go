package objs

import (
	"github.com/MadBase/MadNet/application/objs/atomicswap"
	mdefs "github.com/MadBase/MadNet/application/objs/capn"
	"github.com/MadBase/MadNet/application/objs/uint256"
	"github.com/MadBase/MadNet/application/wrapper"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/errorz"
	"github.com/MadBase/MadNet/utils"
	capnp "zombiezen.com/go/capnproto2"
)

// AtomicSwap is an atomic swap object based on a time lock hash
type AtomicSwap struct {
	ASPreImage *ASPreImage
	TxHash     []byte
	//
	utxoID []byte
}

// UnmarshalBinary takes a byte slice and returns the corresponding
// AtomicSwap object
func (b *AtomicSwap) UnmarshalBinary(data []byte) error {
	bc, err := atomicswap.Unmarshal(data)
	if err != nil {
		return err
	}
	return b.UnmarshalCapn(bc)
}

// MarshalBinary takes the AtomicSwap object and returns the canonical
// byte slice
func (b *AtomicSwap) MarshalBinary() ([]byte, error) {
	if b == nil {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	bc, err := b.MarshalCapn(nil)
	if err != nil {
		return nil, err
	}
	return atomicswap.Marshal(bc)
}

// UnmarshalCapn unmarshals the capnproto definition of the object
func (b *AtomicSwap) UnmarshalCapn(bc mdefs.AtomicSwap) error {
	if err := atomicswap.Validate(bc); err != nil {
		return err
	}
	b.ASPreImage = &ASPreImage{}
	if err := b.ASPreImage.UnmarshalCapn(bc.ASPreImage()); err != nil {
		return err
	}
	b.TxHash = utils.CopySlice(bc.TxHash())
	return nil
}

// MarshalCapn marshals the object into its capnproto definition
func (b *AtomicSwap) MarshalCapn(seg *capnp.Segment) (mdefs.AtomicSwap, error) {
	if b == nil {
		return mdefs.AtomicSwap{}, errorz.ErrInvalid{}.New("not initialized")
	}
	var bc mdefs.AtomicSwap
	if seg == nil {
		_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		if err != nil {
			return bc, err
		}
		tmp, err := mdefs.NewRootAtomicSwap(seg)
		if err != nil {
			return bc, err
		}
		bc = tmp
	} else {
		tmp, err := mdefs.NewAtomicSwap(seg)
		if err != nil {
			return bc, err
		}
		bc = tmp
	}
	bt, err := b.ASPreImage.MarshalCapn(seg)
	if err != nil {
		return bc, err
	}
	if err := bc.SetTxHash(utils.CopySlice(b.TxHash)); err != nil {
		return bc, err
	}
	if err := bc.SetASPreImage(bt); err != nil {
		return bc, err
	}
	return bc, nil
}

// PreHash returns the PreHash of the object
func (b *AtomicSwap) PreHash() ([]byte, error) {
	if b == nil {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	return b.ASPreImage.PreHash()
}

// UTXOID returns the UTXOID of the object
func (b *AtomicSwap) UTXOID() ([]byte, error) {
	if b == nil || b.ASPreImage == nil {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	if b.utxoID != nil {
		return utils.CopySlice(b.utxoID), nil
	}
	b.utxoID = MakeUTXOID(b.TxHash, b.ASPreImage.TXOutIdx)
	return utils.CopySlice(b.utxoID), nil
}

// TXOutIdx returns the TXOutIdx of the object
func (b *AtomicSwap) TXOutIdx() (uint32, error) {
	if b == nil || b.ASPreImage == nil {
		return 0, errorz.ErrInvalid{}.New("not initialized")
	}
	return b.ASPreImage.TXOutIdx, nil
}

// SetTXOutIdx sets the TXOutIdx of the object
func (b *AtomicSwap) SetTXOutIdx(idx uint32) error {
	if b == nil || b.ASPreImage == nil {
		return errorz.ErrInvalid{}.New("not initialized")
	}
	b.ASPreImage.TXOutIdx = idx
	return nil
}

// SetTxHash sets the TxHash of the object
func (b *AtomicSwap) SetTxHash(txHash []byte) error {
	if b == nil || b.ASPreImage == nil {
		return errorz.ErrInvalid{}.New("not initialized")
	}
	if len(txHash) != constants.HashLen {
		return errorz.ErrInvalid{}.New("Invalid hash length")
	}
	b.TxHash = utils.CopySlice(txHash)
	return nil
}

// Value returns the Value of the object
func (b *AtomicSwap) Value() (*uint256.Uint256, error) {
	if b == nil || b.ASPreImage == nil || b.ASPreImage.Value == nil || b.ASPreImage.Value.IsZero() {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	return b.ASPreImage.Value.Clone(), nil
}

// Fee returns the Fee of the object
func (b *AtomicSwap) Fee() (*uint256.Uint256, error) {
	if b == nil || b.ASPreImage == nil || b.ASPreImage.Fee == nil {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	return b.ASPreImage.Fee.Clone(), nil
}

// ValuePlusFee returns the Value of the object with the associated fee
func (b *AtomicSwap) ValuePlusFee() (*uint256.Uint256, error) {
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

// Owner returns the AtomicSwapOwner object of the AtomicSwap
func (b *AtomicSwap) Owner() (*AtomicSwapOwner, error) {
	if b == nil || b.ASPreImage == nil {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	if err := b.ASPreImage.Owner.Validate(); err != nil {
		return nil, errorz.ErrInvalid{}.New("AtomicSwapOwner invalid")
	}
	return b.ASPreImage.Owner, nil
}

// GenericOwner returns the PrimaryOwner of the AtomicSwap as an Owner object
func (b *AtomicSwap) GenericOwner() (*Owner, error) {
	aso, err := b.Owner()
	if err != nil {
		return nil, err
	}
	onr := &Owner{}
	err = onr.NewFromAtomicSwapOwner(aso)
	if err != nil {
		return nil, err
	}
	return onr, nil
}

// ChainID returns the ChainID of the AtomicSwap
func (b *AtomicSwap) ChainID() (uint32, error) {
	if b == nil || b.ASPreImage == nil || b.ASPreImage.ChainID == 0 {
		return 0, errorz.ErrInvalid{}.New("not initialized")
	}
	return b.ASPreImage.ChainID, nil
}

// Exp returns the epoch after which the AtomicSwap will expire
func (b *AtomicSwap) Exp() (uint32, error) {
	if b == nil || b.ASPreImage == nil || b.ASPreImage.Exp == 0 {
		return 0, errorz.ErrInvalid{}.New("not initialized")
	}
	return b.ASPreImage.Exp, nil
}

// IssuedAt returns the epoch of issuance for the AtomicSwap
func (b *AtomicSwap) IssuedAt() (uint32, error) {
	if b == nil || b.ASPreImage == nil || b.ASPreImage.IssuedAt == 0 {
		return 0, errorz.ErrInvalid{}.New("not initialized")
	}
	return b.ASPreImage.IssuedAt, nil
}

// IsExpired returns true if the current epoch is greater than exp
func (b *AtomicSwap) IsExpired(currentHeight uint32) (bool, error) {
	if b == nil {
		return true, errorz.ErrInvalid{}.New("not initialized")
	}
	return b.ASPreImage.IsExpired(currentHeight)
}

// ValidateFee validates the fee of the object at the time of creation
func (b *AtomicSwap) ValidateFee(storage *wrapper.Storage) error {
	fee, err := b.Fee()
	if err != nil {
		return err
	}
	feeTrue, err := storage.GetAtomicSwapFee()
	if err != nil {
		return err
	}
	if fee.Cmp(feeTrue) != 0 {
		return errorz.ErrInvalid{}.New("invalid fee")
	}
	return nil
}

/*
// ValidateSignature validates the signature of the TXIn against the atomic swap
func (b *AtomicSwap) ValidateSignature(currentHeight uint32, txIn *TXIn) error {
	if b == nil {
		return errorz.ErrInvalid{}.New("not initialized")
	}
	msg, err := txIn.TXInLinker.MarshalBinary()
	if err != nil {
		return err
	}
	sig := &AtomicSwapSignature{}
	if err := sig.UnmarshalBinary(txIn.Signature); err != nil {
		return err
	}
	return b.ASPreImage.ValidateSignature(currentHeight, msg, sig)
}
*/
// // SignAsPrimary signs the object as the user who is the original creator of the AtomicSwap
// func (b *AtomicSwap) SignAsPrimary(txIn *TXIn, signer *crypto.Secp256k1Signer, hashKey []byte) error {
// 	if b == nil {
// 		return errorz.ErrInvalid{}.New("not initialized")
// 	}
// 	msg, err := txIn.TXInLinker.MarshalBinary()
// 	if err != nil {
// 		return err
// 	}
// 	sig, err := b.ASPreImage.SignAsPrimary(msg, signer, hashKey)
// 	if err != nil {
// 		return err
// 	}
// 	sigb, err := sig.MarshalBinary()
// 	if err != nil {
// 		return err
// 	}
// 	txIn.Signature = sigb
// 	return nil
// }

// // SignAsAlternate signs the object as the user who is exchanging in the AtomicSwap
// func (b *AtomicSwap) SignAsAlternate(txIn *TXIn, signer *crypto.Secp256k1Signer, hashKey []byte) error {
// 	if b == nil {
// 		return errorz.ErrInvalid{}.New("not initialized")
// 	}
// 	msg, err := txIn.TXInLinker.MarshalBinary()
// 	if err != nil {
// 		return err
// 	}
// 	sig, err := b.ASPreImage.SignAsAlternate(msg, signer, hashKey)
// 	if err != nil {
// 		return err
// 	}
// 	sigb, err := sig.MarshalBinary()
// 	if err != nil {
// 		return err
// 	}
// 	txIn.Signature = sigb
// 	return nil
// }

// MakeTxIn constructs a TXIn object for the current object
func (b *AtomicSwap) MakeTxIn() (*TXIn, error) {
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
