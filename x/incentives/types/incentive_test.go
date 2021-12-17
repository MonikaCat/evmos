package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/tharsis/ethermint/tests"
)

type IncentiveTestSuite struct {
	suite.Suite
}

func TestIncentiveSuite(t *testing.T) {
	suite.Run(t, new(IncentiveTestSuite))
}

func (suite *IncentiveTestSuite) TestIncentiveNew() {
	testCases := []struct {
		name        string
		contract    common.Address
		allocations sdk.DecCoins
		epochs      uint32
		expectPass  bool
	}{
		{
			"Register incentive - pass",
			tests.GenerateAddress(),
			sdk.DecCoins{sdk.NewDecCoin("aphoton", sdk.NewInt(1))},
			10,
			true,
		},
		// TODO: should it be allowed to create empty allocation?
		// {
		// 	"Register incentive - empty allocation",
		// 	tests.GenerateAddress(),
		// 	sdk.DecCoins{},
		// 	10,
		// 	false,
		// },
		// TODO how to test without panicking?
		// {
		// 	"Register incentive - invalid allocation denom",
		// 	tests.GenerateAddress(),
		// 	sdk.DecCoins{sdk.NewDecCoin("(photon", sdk.NewInt(1))},
		// 	10,
		// 	false,
		// },
		{
			"Register incentive - invalid allocation amount",
			tests.GenerateAddress(),
			sdk.DecCoins{sdk.NewDecCoin("aphoton", sdk.NewInt(5))},
			10,
			false,
		},
		{
			"Register incentive - zero epochs",
			tests.GenerateAddress(),
			sdk.DecCoins{sdk.NewDecCoin("aphoton", sdk.NewInt(1))},
			0,
			false,
		},
	}

	for _, tc := range testCases {
		tp := NewIncentive(tc.contract, tc.allocations, tc.epochs)
		err := tp.Validate()

		if tc.expectPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}

func (suite *IncentiveTestSuite) TestIncentive() {
	testCases := []struct {
		msg        string
		incentive  Incentive
		expectPass bool
	}{
		{
			"Register token pair - invalid address (no hex)",
			Incentive{
				"0x5dCA2483280D9727c80b5518faC4556617fb19ZZ",
				sdk.DecCoins{sdk.NewDecCoin("aphoton", sdk.NewInt(1))},
				10,
				time.Now(),
			},
			false,
		},
		{
			"Register token pair - invalid address (invalid length 1)",
			Incentive{
				"0x5dCA2483280D9727c80b5518faC4556617fb19",
				sdk.DecCoins{sdk.NewDecCoin("aphoton", sdk.NewInt(1))},
				10,
				time.Now(),
			},
			false,
		},
		{
			"Register token pair - invalid address (invalid length 2)",
			Incentive{
				"0x5dCA2483280D9727c80b5518faC4556617fb194FFF",
				sdk.DecCoins{sdk.NewDecCoin("aphoton", sdk.NewInt(1))},
				10,
				time.Now(),
			},
			false,
		},
		{
			"pass",
			Incentive{
				tests.GenerateAddress().String(),
				sdk.DecCoins{sdk.NewDecCoin("aphoton", sdk.NewInt(1))},
				10,
				time.Now(),
			},
			true,
		},
	}

	for _, tc := range testCases {
		err := tc.incentive.Validate()

		if tc.expectPass {
			suite.Require().NoError(err, tc.msg)
		} else {
			suite.Require().Error(err, tc.msg)
		}
	}
}

func (suite *IncentiveTestSuite) TestIsActive() {
	testCases := []struct {
		name       string
		incentive  Incentive
		expectPass bool
	}{
		{
			"pass",
			Incentive{
				tests.GenerateAddress().String(),
				sdk.DecCoins{sdk.NewDecCoin("aphoton", sdk.NewInt(1))},
				10,
				time.Now(),
			},
			true,
		},
		{
			"epoch is zero",
			Incentive{
				tests.GenerateAddress().String(),
				sdk.DecCoins{sdk.NewDecCoin("aphoton", sdk.NewInt(1))},
				0,
				time.Now(),
			},
			false,
		},
	}
	for _, tc := range testCases {
		res := tc.incentive.IsActive()
		if tc.expectPass {
			suite.Require().True(res, tc.name)
		} else {
			suite.Require().False(res, tc.name)
		}
	}
}
