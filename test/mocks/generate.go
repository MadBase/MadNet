package mocks

//go:generate go-mockgen -f -i Ethereum -i Contracts -i AdminHandler -i TxnQueue -i GethClient -o interfaces.mockgen.go github.com/MadBase/MadNet/blockchain/interfaces
//go:generate go-mockgen -f -i IETHDKG -i IGovernance -i IMadByte -i IMadToken -i IMadnetFactory -i IPublicStaking -i ISnapshots -i IValidatorPool -i IValidatorStaking  -o bindings.mockgen.go github.com/MadBase/MadNet/bridge/bindings
