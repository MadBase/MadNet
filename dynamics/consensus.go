package dynamics

import (
	"time"
)

// Original from constants/consensus.go
const (
	maxBytes        = 3000000
	maxProposalSize = maxBytes // Parameterize: equal to maxBytes
	msgTimeout      = 4 * time.Second
	srvrMsgTimeout  = (3 * msgTimeout) / 4 // Parameterize: 0.75*MsgTimeout
	proposalStepTO  = 4 * time.Second
	preVoteStepTO   = 3 * time.Second
	preCommitStepTO = 3 * time.Second
	dBRNRTO         = (5 * (proposalStepTO + preVoteStepTO + preCommitStepTO)) / 2 // dead block round next round TO // Parameterize: make 2.5 times Prop, PV, PC timeouts
	downloadTO      = proposalStepTO + preVoteStepTO + preCommitStepTO             // Parameterize: sum of Prop, PV, PC timeouts
)

/* block time lower bound: 2*proposalStepTO
 * block time upper bound: proposalStepTO + 2*dBRNRTO
 */
