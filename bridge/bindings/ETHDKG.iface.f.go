// Generated by ifacemaker. DO NOT EDIT.

package bindings

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// IETHDKGFilterer ...
type IETHDKGFilterer interface {
	// FilterAddressRegistered is a free log retrieval operation binding the contract event 0x7f1304057ec61140fbf2f5f236790f34fcafe123d3eb0d298d92317c97da500d.
	//
	// Solidity: event AddressRegistered(address account, uint256 index, uint256 nonce, uint256[2] publicKey)
	FilterAddressRegistered(opts *bind.FilterOpts) (*ETHDKGAddressRegisteredIterator, error)
	// WatchAddressRegistered is a free log subscription operation binding the contract event 0x7f1304057ec61140fbf2f5f236790f34fcafe123d3eb0d298d92317c97da500d.
	//
	// Solidity: event AddressRegistered(address account, uint256 index, uint256 nonce, uint256[2] publicKey)
	WatchAddressRegistered(opts *bind.WatchOpts, sink chan<- *ETHDKGAddressRegistered) (event.Subscription, error)
	// ParseAddressRegistered is a log parse operation binding the contract event 0x7f1304057ec61140fbf2f5f236790f34fcafe123d3eb0d298d92317c97da500d.
	//
	// Solidity: event AddressRegistered(address account, uint256 index, uint256 nonce, uint256[2] publicKey)
	ParseAddressRegistered(log types.Log) (*ETHDKGAddressRegistered, error)
	// FilterGPKJSubmissionComplete is a free log retrieval operation binding the contract event 0x87bfe600b78cad9f7cf68c99eb582c1748f636b3269842b37d5873b0e069f628.
	//
	// Solidity: event GPKJSubmissionComplete(uint256 blockNumber)
	FilterGPKJSubmissionComplete(opts *bind.FilterOpts) (*ETHDKGGPKJSubmissionCompleteIterator, error)
	// WatchGPKJSubmissionComplete is a free log subscription operation binding the contract event 0x87bfe600b78cad9f7cf68c99eb582c1748f636b3269842b37d5873b0e069f628.
	//
	// Solidity: event GPKJSubmissionComplete(uint256 blockNumber)
	WatchGPKJSubmissionComplete(opts *bind.WatchOpts, sink chan<- *ETHDKGGPKJSubmissionComplete) (event.Subscription, error)
	// ParseGPKJSubmissionComplete is a log parse operation binding the contract event 0x87bfe600b78cad9f7cf68c99eb582c1748f636b3269842b37d5873b0e069f628.
	//
	// Solidity: event GPKJSubmissionComplete(uint256 blockNumber)
	ParseGPKJSubmissionComplete(log types.Log) (*ETHDKGGPKJSubmissionComplete, error)
	// FilterKeyShareSubmissionComplete is a free log retrieval operation binding the contract event 0x522cec98f6caa194456c44afa9e8cef9ac63eecb0be60e20d180ce19cfb0ef59.
	//
	// Solidity: event KeyShareSubmissionComplete(uint256 blockNumber)
	FilterKeyShareSubmissionComplete(opts *bind.FilterOpts) (*ETHDKGKeyShareSubmissionCompleteIterator, error)
	// WatchKeyShareSubmissionComplete is a free log subscription operation binding the contract event 0x522cec98f6caa194456c44afa9e8cef9ac63eecb0be60e20d180ce19cfb0ef59.
	//
	// Solidity: event KeyShareSubmissionComplete(uint256 blockNumber)
	WatchKeyShareSubmissionComplete(opts *bind.WatchOpts, sink chan<- *ETHDKGKeyShareSubmissionComplete) (event.Subscription, error)
	// ParseKeyShareSubmissionComplete is a log parse operation binding the contract event 0x522cec98f6caa194456c44afa9e8cef9ac63eecb0be60e20d180ce19cfb0ef59.
	//
	// Solidity: event KeyShareSubmissionComplete(uint256 blockNumber)
	ParseKeyShareSubmissionComplete(log types.Log) (*ETHDKGKeyShareSubmissionComplete, error)
	// FilterKeyShareSubmitted is a free log retrieval operation binding the contract event 0x6162e2d11398e4063e4c8565dafc4fb6755bbead93747ea836a5ef73a594aaf7.
	//
	// Solidity: event KeyShareSubmitted(address account, uint256 index, uint256 nonce, uint256[2] keyShareG1, uint256[2] keyShareG1CorrectnessProof, uint256[4] keyShareG2)
	FilterKeyShareSubmitted(opts *bind.FilterOpts) (*ETHDKGKeyShareSubmittedIterator, error)
	// WatchKeyShareSubmitted is a free log subscription operation binding the contract event 0x6162e2d11398e4063e4c8565dafc4fb6755bbead93747ea836a5ef73a594aaf7.
	//
	// Solidity: event KeyShareSubmitted(address account, uint256 index, uint256 nonce, uint256[2] keyShareG1, uint256[2] keyShareG1CorrectnessProof, uint256[4] keyShareG2)
	WatchKeyShareSubmitted(opts *bind.WatchOpts, sink chan<- *ETHDKGKeyShareSubmitted) (event.Subscription, error)
	// ParseKeyShareSubmitted is a log parse operation binding the contract event 0x6162e2d11398e4063e4c8565dafc4fb6755bbead93747ea836a5ef73a594aaf7.
	//
	// Solidity: event KeyShareSubmitted(address account, uint256 index, uint256 nonce, uint256[2] keyShareG1, uint256[2] keyShareG1CorrectnessProof, uint256[4] keyShareG2)
	ParseKeyShareSubmitted(log types.Log) (*ETHDKGKeyShareSubmitted, error)
	// FilterMPKSet is a free log retrieval operation binding the contract event 0x71b1ebd27be320895a22125d6458e3363aefa6944a312ede4bf275867e6d5a71.
	//
	// Solidity: event MPKSet(uint256 blockNumber, uint256 nonce, uint256[4] mpk)
	FilterMPKSet(opts *bind.FilterOpts) (*ETHDKGMPKSetIterator, error)
	// WatchMPKSet is a free log subscription operation binding the contract event 0x71b1ebd27be320895a22125d6458e3363aefa6944a312ede4bf275867e6d5a71.
	//
	// Solidity: event MPKSet(uint256 blockNumber, uint256 nonce, uint256[4] mpk)
	WatchMPKSet(opts *bind.WatchOpts, sink chan<- *ETHDKGMPKSet) (event.Subscription, error)
	// ParseMPKSet is a log parse operation binding the contract event 0x71b1ebd27be320895a22125d6458e3363aefa6944a312ede4bf275867e6d5a71.
	//
	// Solidity: event MPKSet(uint256 blockNumber, uint256 nonce, uint256[4] mpk)
	ParseMPKSet(log types.Log) (*ETHDKGMPKSet, error)
	// FilterRegistrationComplete is a free log retrieval operation binding the contract event 0x833013b96b786b4eca83baac286920e5e53956c21ff3894f1d9f02e97d6ed764.
	//
	// Solidity: event RegistrationComplete(uint256 blockNumber)
	FilterRegistrationComplete(opts *bind.FilterOpts) (*ETHDKGRegistrationCompleteIterator, error)
	// WatchRegistrationComplete is a free log subscription operation binding the contract event 0x833013b96b786b4eca83baac286920e5e53956c21ff3894f1d9f02e97d6ed764.
	//
	// Solidity: event RegistrationComplete(uint256 blockNumber)
	WatchRegistrationComplete(opts *bind.WatchOpts, sink chan<- *ETHDKGRegistrationComplete) (event.Subscription, error)
	// ParseRegistrationComplete is a log parse operation binding the contract event 0x833013b96b786b4eca83baac286920e5e53956c21ff3894f1d9f02e97d6ed764.
	//
	// Solidity: event RegistrationComplete(uint256 blockNumber)
	ParseRegistrationComplete(log types.Log) (*ETHDKGRegistrationComplete, error)
	// FilterRegistrationOpened is a free log retrieval operation binding the contract event 0xbda431b9b63510f1398bf33d700e013315bcba905507078a1780f13ea5b354b9.
	//
	// Solidity: event RegistrationOpened(uint256 startBlock, uint256 numberValidators, uint256 nonce, uint256 phaseLength, uint256 confirmationLength)
	FilterRegistrationOpened(opts *bind.FilterOpts) (*ETHDKGRegistrationOpenedIterator, error)
	// WatchRegistrationOpened is a free log subscription operation binding the contract event 0xbda431b9b63510f1398bf33d700e013315bcba905507078a1780f13ea5b354b9.
	//
	// Solidity: event RegistrationOpened(uint256 startBlock, uint256 numberValidators, uint256 nonce, uint256 phaseLength, uint256 confirmationLength)
	WatchRegistrationOpened(opts *bind.WatchOpts, sink chan<- *ETHDKGRegistrationOpened) (event.Subscription, error)
	// ParseRegistrationOpened is a log parse operation binding the contract event 0xbda431b9b63510f1398bf33d700e013315bcba905507078a1780f13ea5b354b9.
	//
	// Solidity: event RegistrationOpened(uint256 startBlock, uint256 numberValidators, uint256 nonce, uint256 phaseLength, uint256 confirmationLength)
	ParseRegistrationOpened(log types.Log) (*ETHDKGRegistrationOpened, error)
	// FilterShareDistributionComplete is a free log retrieval operation binding the contract event 0xbfe94ffef5ddde4d25ac7b652f3f67686ea63f9badbfe1f25451e26fc262d11c.
	//
	// Solidity: event ShareDistributionComplete(uint256 blockNumber)
	FilterShareDistributionComplete(opts *bind.FilterOpts) (*ETHDKGShareDistributionCompleteIterator, error)
	// WatchShareDistributionComplete is a free log subscription operation binding the contract event 0xbfe94ffef5ddde4d25ac7b652f3f67686ea63f9badbfe1f25451e26fc262d11c.
	//
	// Solidity: event ShareDistributionComplete(uint256 blockNumber)
	WatchShareDistributionComplete(opts *bind.WatchOpts, sink chan<- *ETHDKGShareDistributionComplete) (event.Subscription, error)
	// ParseShareDistributionComplete is a log parse operation binding the contract event 0xbfe94ffef5ddde4d25ac7b652f3f67686ea63f9badbfe1f25451e26fc262d11c.
	//
	// Solidity: event ShareDistributionComplete(uint256 blockNumber)
	ParseShareDistributionComplete(log types.Log) (*ETHDKGShareDistributionComplete, error)
	// FilterSharesDistributed is a free log retrieval operation binding the contract event 0xf0c8b0ef2867c2b4639b404a0296b6bbf0bf97e20856af42144a5a6035c0d0d2.
	//
	// Solidity: event SharesDistributed(address account, uint256 index, uint256 nonce, uint256[] encryptedShares, uint256[2][] commitments)
	FilterSharesDistributed(opts *bind.FilterOpts) (*ETHDKGSharesDistributedIterator, error)
	// WatchSharesDistributed is a free log subscription operation binding the contract event 0xf0c8b0ef2867c2b4639b404a0296b6bbf0bf97e20856af42144a5a6035c0d0d2.
	//
	// Solidity: event SharesDistributed(address account, uint256 index, uint256 nonce, uint256[] encryptedShares, uint256[2][] commitments)
	WatchSharesDistributed(opts *bind.WatchOpts, sink chan<- *ETHDKGSharesDistributed) (event.Subscription, error)
	// ParseSharesDistributed is a log parse operation binding the contract event 0xf0c8b0ef2867c2b4639b404a0296b6bbf0bf97e20856af42144a5a6035c0d0d2.
	//
	// Solidity: event SharesDistributed(address account, uint256 index, uint256 nonce, uint256[] encryptedShares, uint256[2][] commitments)
	ParseSharesDistributed(log types.Log) (*ETHDKGSharesDistributed, error)
	// FilterValidatorMemberAdded is a free log retrieval operation binding the contract event 0x09b90b08bbc3dbe22e9d2a0bc9c2c7614c7511cd0ad72177727a1e762115bf06.
	//
	// Solidity: event ValidatorMemberAdded(address account, uint256 index, uint256 nonce, uint256 epoch, uint256 share0, uint256 share1, uint256 share2, uint256 share3)
	FilterValidatorMemberAdded(opts *bind.FilterOpts) (*ETHDKGValidatorMemberAddedIterator, error)
	// WatchValidatorMemberAdded is a free log subscription operation binding the contract event 0x09b90b08bbc3dbe22e9d2a0bc9c2c7614c7511cd0ad72177727a1e762115bf06.
	//
	// Solidity: event ValidatorMemberAdded(address account, uint256 index, uint256 nonce, uint256 epoch, uint256 share0, uint256 share1, uint256 share2, uint256 share3)
	WatchValidatorMemberAdded(opts *bind.WatchOpts, sink chan<- *ETHDKGValidatorMemberAdded) (event.Subscription, error)
	// ParseValidatorMemberAdded is a log parse operation binding the contract event 0x09b90b08bbc3dbe22e9d2a0bc9c2c7614c7511cd0ad72177727a1e762115bf06.
	//
	// Solidity: event ValidatorMemberAdded(address account, uint256 index, uint256 nonce, uint256 epoch, uint256 share0, uint256 share1, uint256 share2, uint256 share3)
	ParseValidatorMemberAdded(log types.Log) (*ETHDKGValidatorMemberAdded, error)
	// FilterValidatorSetCompleted is a free log retrieval operation binding the contract event 0xd7237b781669fa700ecf77be6cd8fa0f4b98b1a24ac584a9b6b44c509216718a.
	//
	// Solidity: event ValidatorSetCompleted(uint256 validatorCount, uint256 nonce, uint256 epoch, uint256 ethHeight, uint256 aliceNetHeight, uint256 groupKey0, uint256 groupKey1, uint256 groupKey2, uint256 groupKey3)
	FilterValidatorSetCompleted(opts *bind.FilterOpts) (*ETHDKGValidatorSetCompletedIterator, error)
	// WatchValidatorSetCompleted is a free log subscription operation binding the contract event 0xd7237b781669fa700ecf77be6cd8fa0f4b98b1a24ac584a9b6b44c509216718a.
	//
	// Solidity: event ValidatorSetCompleted(uint256 validatorCount, uint256 nonce, uint256 epoch, uint256 ethHeight, uint256 aliceNetHeight, uint256 groupKey0, uint256 groupKey1, uint256 groupKey2, uint256 groupKey3)
	WatchValidatorSetCompleted(opts *bind.WatchOpts, sink chan<- *ETHDKGValidatorSetCompleted) (event.Subscription, error)
	// ParseValidatorSetCompleted is a log parse operation binding the contract event 0xd7237b781669fa700ecf77be6cd8fa0f4b98b1a24ac584a9b6b44c509216718a.
	//
	// Solidity: event ValidatorSetCompleted(uint256 validatorCount, uint256 nonce, uint256 epoch, uint256 ethHeight, uint256 aliceNetHeight, uint256 groupKey0, uint256 groupKey1, uint256 groupKey2, uint256 groupKey3)
	ParseValidatorSetCompleted(log types.Log) (*ETHDKGValidatorSetCompleted, error)
}