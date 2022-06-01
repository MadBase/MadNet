package bindings

type IAToken interface {
	IATokenCaller
	IATokenTransactor
	IATokenFilterer
}

type IATokenBurner interface {
	IATokenBurnerCaller
	IATokenBurnerTransactor
	IATokenBurnerFilterer
}

type IATokenMinter interface {
	IATokenMinterCaller
	IATokenMinterTransactor
	IATokenMinterFilterer
}

type IAliceNetFactory interface {
	IAliceNetFactoryCaller
	IAliceNetFactoryTransactor
	IAliceNetFactoryFilterer
}

type IBToken interface {
	IBTokenCaller
	IBTokenTransactor
	IBTokenFilterer
}

type IBTokenErrorCodes interface {
	IBTokenErrorCodesCaller
	IBTokenErrorCodesTransactor
	IBTokenErrorCodesFilterer
}

type IETHDKG interface {
	IETHDKGCaller
	IETHDKGTransactor
	IETHDKGFilterer
}

type IETHDKGErrorCodes interface {
	IETHDKGErrorCodesCaller
	IETHDKGErrorCodesTransactor
	IETHDKGErrorCodesFilterer
}

type IGovernance interface {
	IGovernanceCaller
	IGovernanceTransactor
	IGovernanceFilterer
}

type IGovernanceErrorCodes interface {
	IGovernanceErrorCodesCaller
	IGovernanceErrorCodesTransactor
	IGovernanceErrorCodesFilterer
}

type IPublicStaking interface {
	IPublicStakingCaller
	IPublicStakingTransactor
	IPublicStakingFilterer
}

type ISnapshots interface {
	ISnapshotsCaller
	ISnapshotsTransactor
	ISnapshotsFilterer
}

type ISnapshotsErrorCodes interface {
	ISnapshotsErrorCodesCaller
	ISnapshotsErrorCodesTransactor
	ISnapshotsErrorCodesFilterer
}

type IValidatorPool interface {
	IValidatorPoolCaller
	IValidatorPoolTransactor
	IValidatorPoolFilterer
}

type IValidatorPoolErrorCodes interface {
	IValidatorPoolErrorCodesCaller
	IValidatorPoolErrorCodesTransactor
	IValidatorPoolErrorCodesFilterer
}

type IValidatorStaking interface {
	IValidatorStakingCaller
	IValidatorStakingTransactor
	IValidatorStakingFilterer
}
