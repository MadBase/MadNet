package objs

import (
	"bytes"

	"github.com/MadBase/MadNet/errorz"

	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/crypto"
	"github.com/MadBase/MadNet/utils"
)

// AtomicSwapOwner describes the necessary information for AtomicSwap object
type AtomicSwapOwner struct {
	SVA            SVA
	HashLock       []byte
	AlternateOwner *AtomicSwapSubOwner
	PrimaryOwner   *AtomicSwapSubOwner
}

// New creates a new AtomicSwapOwner
func (aso *AtomicSwapOwner) New(priOwnerAcct []byte, altOwnerAcct []byte, hashKey []byte) error {
	if aso == nil {
		return errorz.ErrInvalid{}.New("not initialized")
	}
	if len(hashKey) != constants.HashLen {
		return errorz.ErrInvalid{}.New("Error in ASO.New; invalid hashKey")
	}
	if len(priOwnerAcct) != constants.OwnerLen {
		return errorz.ErrInvalid{}.New("Error in ASO.New; invalid primary account length")
	}
	if len(altOwnerAcct) != constants.OwnerLen {
		return errorz.ErrInvalid{}.New("Error in ASO.New; invalid alternate account length")
	}
	aso.SVA = HashedTimelockSVA
	aso.HashLock = crypto.Hasher(hashKey)
	aso.PrimaryOwner = &AtomicSwapSubOwner{
		CurveSpec: constants.CurveSecp256k1,
		Account:   utils.CopySlice(priOwnerAcct),
	}
	aso.AlternateOwner = &AtomicSwapSubOwner{
		CurveSpec: constants.CurveSecp256k1,
		Account:   utils.CopySlice(altOwnerAcct),
	}
	return nil
}

// NewFromOwner creates a new AtomicSwapOwner from Owner objects
func (aso *AtomicSwapOwner) NewFromOwner(priOwner *Owner, altOwner *Owner, hashKey []byte) error {
	if aso == nil {
		return errorz.ErrInvalid{}.New("not initialized")
	}
	if len(hashKey) != constants.HashLen {
		return errorz.ErrInvalid{}.New("Error in ASO.New; invalid hashKey")
	}
	aso.SVA = HashedTimelockSVA
	aso.HashLock = crypto.Hasher(hashKey)
	aso.PrimaryOwner = &AtomicSwapSubOwner{}
	err := aso.PrimaryOwner.NewFromOwner(priOwner)
	if err != nil {
		aso.SVA = 0
		aso.HashLock = nil
		aso.PrimaryOwner = nil
		aso.AlternateOwner = nil
		return err
	}
	aso.AlternateOwner = &AtomicSwapSubOwner{}
	err = aso.AlternateOwner.NewFromOwner(altOwner)
	if err != nil {
		aso.SVA = 0
		aso.HashLock = nil
		aso.PrimaryOwner = nil
		aso.AlternateOwner = nil
		return err
	}
	return nil
}

// MarshalBinary takes the AtomicSwapOwner object and returns the canonical
// byte slice
func (aso *AtomicSwapOwner) MarshalBinary() ([]byte, error) {
	if err := aso.Validate(); err != nil {
		return nil, err
	}
	owner := []byte{}
	owner = append(owner, []byte{uint8(aso.SVA)}...)
	owner = append(owner, utils.CopySlice(aso.HashLock)...)
	priOwner, err := aso.PrimaryOwner.MarshalBinary()
	if err != nil {
		return nil, err
	}
	owner = append(owner, priOwner...)
	altOwner, err := aso.AlternateOwner.MarshalBinary()
	if err != nil {
		return nil, err
	}
	owner = append(owner, altOwner...)
	return owner, nil
}

// UnmarshalBinary takes a byte slice and returns the corresponding
// AtomicSwapOwner object
func (aso *AtomicSwapOwner) UnmarshalBinary(o []byte) error {
	if aso == nil {
		return errorz.ErrInvalid{}.New("not initialized")
	}
	owner := utils.CopySlice(o)
	sva, owner, err := extractSVA(owner)
	if err != nil {
		return err
	}
	aso.SVA = sva
	if err := aso.validateSVA(); err != nil {
		return err
	}
	hashLock, owner, err := extractHash(owner)
	if err != nil {
		return err
	}
	aso.HashLock = hashLock
	priOwner := &AtomicSwapSubOwner{}
	owner, err = priOwner.UnmarshalBinary(owner)
	if err != nil {
		return err
	}
	aso.PrimaryOwner = priOwner
	altOwner := &AtomicSwapSubOwner{}
	owner, err = altOwner.UnmarshalBinary(owner)
	if err != nil {
		return err
	}
	aso.AlternateOwner = altOwner
	if err := extractZero(owner); err != nil {
		return err
	}
	if err := aso.Validate(); err != nil {
		return err
	}
	return nil
}

// PrimaryAccount returns the account byte slice of the PrimaryOwner
func (aso *AtomicSwapOwner) PrimaryAccount() ([]byte, error) {
	if aso == nil || aso.PrimaryOwner == nil {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	return utils.CopySlice(aso.PrimaryOwner.Account), nil
}

// AlternateAccount returns the account byte slice of the AlternateOwner
func (aso *AtomicSwapOwner) AlternateAccount() ([]byte, error) {
	if aso == nil || aso.AlternateOwner == nil {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	return utils.CopySlice(aso.AlternateOwner.Account), nil
}

// Validate validates the AtomicSwapOwner
func (aso *AtomicSwapOwner) Validate() error {
	if aso == nil {
		return errorz.ErrInvalid{}.New("object is nil")
	}
	if err := aso.validateSVA(); err != nil {
		return err
	}
	if err := utils.ValidateHash(aso.HashLock); err != nil {
		return err
	}
	if err := aso.AlternateOwner.Validate(); err != nil {
		return err
	}
	if err := aso.PrimaryOwner.Validate(); err != nil {
		return err
	}
	return nil
}

// ValidateSignature validates the signature
func (aso *AtomicSwapOwner) ValidateSignature(msg []byte, sig *AtomicSwapSignature, isExpired bool) error {
	if err := aso.Validate(); err != nil {
		return errorz.ErrInvalid{}.New("invalid AtomicSwapOwner")
	}
	if err := sig.Validate(); err != nil {
		return errorz.ErrInvalid{}.New("invalid AtomicSwapSignature")
	}
	if aso.SVA != sig.SVA {
		return errorz.ErrInvalid{}.New("incorrect SVA")
	}
	hsh := crypto.Hasher(sig.HashKey)
	if !bytes.Equal(aso.HashLock, hsh) {
		return errorz.ErrInvalid{}.New("incorrect hash key")
	}
	switch sig.SignerRole {
	case PrimarySignerRole:
		if !isExpired {
			return errorz.ErrInvalid{}.New("PrimaryOwner can not sign before expiration")
		}
		// if err := aso.PrimaryOwner.ValidateSignature(msg, sig); err != nil {
		// 	return err
		// }
		return nil
	case AlternateSignerRole:
		if isExpired {
			return errorz.ErrInvalid{}.New("AlternateOwner can not sign after expiration")
		}
		// if err := aso.AlternateOwner.ValidateSignature(msg, sig); err != nil {
		// 	return err
		// }
		return nil
	default:
		return errorz.ErrInvalid{}.New("Invalid signerRole")
	}
}

// validateSVA validates the Signature Verification Algorithm
func (aso *AtomicSwapOwner) validateSVA() error {
	if aso == nil {
		return errorz.ErrInvalid{}.New("not initialized")
	}
	if aso.SVA != HashedTimelockSVA {
		return errorz.ErrInvalid{}.New("signature verification algorithm invalid for AtomicSwapOwner")
	}
	return nil
}

// // SignAsPrimary ...
// func (aso *AtomicSwapOwner) SignAsPrimary(msg []byte, signer *crypto.Secp256k1Signer, hashKey []byte) (*AtomicSwapSignature, error) {
// 	if aso == nil {
// 		return nil, errorz.ErrInvalid{}.New("not initialized")
// 	}
// 	sig, err := signer.Sign(msg)
// 	if err != nil {
// 		return nil, err
// 	}
// 	hsh := crypto.Hasher(hashKey)
// 	if !bytes.Equal(hsh, aso.HashLock) {
// 		return nil, errorz.ErrInvalid{}.New("invalid hash key")
// 	}
// 	s := &AtomicSwapSignature{
// 		SVA:        HashedTimelockSVA,
// 		CurveSpec:  constants.CurveSecp256k1,
// 		SignerRole: PrimarySignerRole,
// 		HashKey:    hashKey,
// 		Signature:  sig,
// 	}
// 	if err := s.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return s, nil
// }

// // SignAsAlternate ...
// func (aso *AtomicSwapOwner) SignAsAlternate(msg []byte, signer *crypto.Secp256k1Signer, hashKey []byte) (*AtomicSwapSignature, error) {
// 	if aso == nil {
// 		return nil, errorz.ErrInvalid{}.New("not initialized")
// 	}
// 	sig, err := signer.Sign(msg)
// 	if err != nil {
// 		return nil, err
// 	}
// 	hsh := crypto.Hasher(hashKey)
// 	if !bytes.Equal(hsh, aso.HashLock) {
// 		return nil, errorz.ErrInvalid{}.New("invalid hash key")
// 	}
// 	s := &AtomicSwapSignature{
// 		SVA:        HashedTimelockSVA,
// 		CurveSpec:  constants.CurveSecp256k1,
// 		SignerRole: AlternateSignerRole,
// 		HashKey:    hashKey,
// 		Signature:  sig,
// 	}
// 	if err := s.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return s, nil
// }

// AtomicSwapSubOwner ...
type AtomicSwapSubOwner struct {
	CurveSpec constants.CurveSpec
	Account   []byte
}

// NewFromOwner takes an Owner object and creates the corresponding
// AtomicSwapSubOwner
func (asso *AtomicSwapSubOwner) NewFromOwner(o *Owner) error {
	if asso == nil {
		return errorz.ErrInvalid{}.New("not initialized")
	}
	if err := o.Validate(); err != nil {
		return err
	}
	asso.CurveSpec = o.CurveSpec
	asso.Account = utils.CopySlice(o.Account)
	if err := asso.Validate(); err != nil {
		asso.CurveSpec = 0
		asso.Account = nil
		return err
	}
	return nil
}

// MarshalBinary takes the AtomicSwapSubOwner object and returns the canonical
// byte slice
func (asso *AtomicSwapSubOwner) MarshalBinary() ([]byte, error) {
	if err := asso.Validate(); err != nil {
		return nil, err
	}
	var owner []byte
	owner = append(owner, []byte{uint8(asso.CurveSpec)}...)
	owner = append(owner, utils.CopySlice(asso.Account)...)
	return owner, nil
}

// UnmarshalBinary takes a byte slice and returns the corresponding
// AtomicSwapSubOwner object
func (asso *AtomicSwapSubOwner) UnmarshalBinary(o []byte) ([]byte, error) {
	owner := utils.CopySlice(o)
	curveSpec, owner, err := extractCurveSpec(owner)
	if err != nil {
		return nil, err
	}
	account, owner, err := extractAccount(owner)
	if err != nil {
		return nil, err
	}
	asso.CurveSpec = curveSpec
	asso.Account = account
	if err := asso.Validate(); err != nil {
		return nil, err
	}
	return owner, nil
}

// ValidateSignature validates the signature of the AtomicSwapSignature object
// func (asso *AtomicSwapSubOwner) ValidateSignature(msg []byte, sig *AtomicSwapSignature) error {
// if err := asso.Validate(); err != nil {
// return errorz.ErrInvalid{}.New("invalid AtomicSwapSubOwner")
// }
// if asso.CurveSpec != sig.CurveSpec {
// return errorz.ErrInvalid{}.New("mismatched curve spec")
// }
// val := crypto.Secp256k1Validator{}
// pk, err := val.Validate(msg, sig.Signature)
// if err != nil {
// return err
// }
// account := crypto.GetAccount(pk)
// if !bytes.Equal(account, asso.Account) {
// return errorz.ErrInvalid{}.New("invalid sig for account")
// }
// return nil
// }

// validateCurveSpec validates the curve specification for AtomicSwapSubOwner
func (asso *AtomicSwapSubOwner) validateCurveSpec() error {
	if asso.CurveSpec != constants.CurveSecp256k1 {
		return errorz.ErrInvalid{}.New("Invalid curveSpec for AtomicSwapSubOwner")
	}
	return nil
}

// validateAccount validates the account for AtomicSwapSubOwner
func (asso *AtomicSwapSubOwner) validateAccount() error {
	if len(asso.Account) != constants.OwnerLen {
		return errorz.ErrInvalid{}.New("account length wrong")
	}
	return nil
}

// Validate validates the AtomicSwapSubOwner object
func (asso *AtomicSwapSubOwner) Validate() error {
	if asso == nil {
		return errorz.ErrInvalid{}.New("object is nil")
	}
	if err := asso.validateCurveSpec(); err != nil {
		return err
	}
	if err := asso.validateAccount(); err != nil {
		return err
	}
	return nil
}
