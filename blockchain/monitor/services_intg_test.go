package monitor_test

import (
	"context"
	"github.com/MadBase/MadNet/blockchain"
	"github.com/MadBase/MadNet/blockchain/dkg/dtest"
	"github.com/MadBase/MadNet/blockchain/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math"
	"math/big"
	"testing"
	"time"
)

type ServicesSuite struct {
	suite.Suite
	eth interfaces.Ethereum
}

func (s *ServicesSuite) SetupTest() {
	t := s.T()

	privateKeys, _ := dtest.InitializePrivateKeysAndAccounts(4)
	eth, err := blockchain.NewEthereumSimulator(
		privateKeys,
		3,
		2*time.Second,
		5*time.Second,
		0,
		big.NewInt(9223372036854775807),
		50,
		math.MaxInt64,
		5*time.Second,
		30*time.Second)

	assert.Nil(t, err, "Error creating Ethereum simulator")

	s.eth = eth
}

func (s *ServicesSuite) TestRegistrationOpenEvent() {
	t := s.T()
	eth := s.eth
	c := eth.Contracts()
	assert.NotNil(t, c, "Need a *Contracts")

	height, err := s.eth.GetCurrentHeight(context.TODO())
	assert.Nil(t, err, "could not get height")
	assert.Equal(t, uint64(0), height, "Height should be 0")

	s.eth.Commit()

	height, err = s.eth.GetCurrentHeight(context.TODO())
	assert.Nil(t, err, "could not get height")
	assert.Equal(t, uint64(1), height, "Height should be 1")
}

func TestServicesSuite(t *testing.T) {
	suite.Run(t, new(ServicesSuite))
}
