package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"tokenchain/testutil/sample"
	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func TestRecordRewardAccrualResponse(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	address := sample.AccAddress()

	resp1, err := srv.RecordRewardAccrual(f.ctx, &types.MsgRecordRewardAccrual{
		Creator: creator,
		Address: address,
		Denom:   "utoken",
		Amount:  100,
		Date:    "2026-02-25",
	})
	require.NoError(t, err)
	require.Equal(t, address+"|utoken", resp1.Key)
	require.Equal(t, address, resp1.Address)
	require.Equal(t, "utoken", resp1.Denom)
	require.EqualValues(t, 100, resp1.AmountAdded)
	require.EqualValues(t, 100, resp1.TotalAmount)
	require.Equal(t, "2026-02-25", resp1.RollupDate)

	resp2, err := srv.RecordRewardAccrual(f.ctx, &types.MsgRecordRewardAccrual{
		Creator: creator,
		Address: address,
		Denom:   "utoken",
		Amount:  25,
		Date:    "2026-02-26",
	})
	require.NoError(t, err)
	require.EqualValues(t, 25, resp2.AmountAdded)
	require.EqualValues(t, 125, resp2.TotalAmount)
	require.Equal(t, "2026-02-26", resp2.RollupDate)
}

func TestClaimRewardResponse(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	address := sample.AccAddress()
	key := address + "|utoken"

	require.NoError(t, f.keeper.Rewardaccrual.Set(f.ctx, key, types.Rewardaccrual{
		Creator:        creator,
		Key:            key,
		Address:        address,
		Denom:          "utoken",
		Amount:         345,
		LastRollupDate: "2026-02-25",
	}))

	coins := sdk.NewCoins(sdk.NewCoin("utoken", sdkmath.NewInt(345)))
	require.NoError(t, f.bankKeeper.MintCoins(f.ctx, types.ModuleName, coins))

	resp, err := srv.ClaimReward(f.ctx, &types.MsgClaimReward{
		Creator: address,
		Denom:   "utoken",
	})
	require.NoError(t, err)
	require.Equal(t, address, resp.Address)
	require.Equal(t, "utoken", resp.Denom)
	require.EqualValues(t, 345, resp.AmountClaimed)

	exists, err := f.keeper.Rewardaccrual.Has(f.ctx, key)
	require.NoError(t, err)
	require.False(t, exists)

	accountBalance := f.bankKeeper.SpendableCoins(f.ctx, sdk.MustAccAddressFromBech32(address))
	require.Equal(t, coins, accountBalance)
}

func TestClaimRewardInsufficientPool(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	address := sample.AccAddress()
	key := address + "|utoken"

	require.NoError(t, f.keeper.Rewardaccrual.Set(f.ctx, key, types.Rewardaccrual{
		Creator:        creator,
		Key:            key,
		Address:        address,
		Denom:          "utoken",
		Amount:         50,
		LastRollupDate: "2026-02-25",
	}))

	_, err := srv.ClaimReward(f.ctx, &types.MsgClaimReward{
		Creator: address,
		Denom:   "utoken",
	})
	require.ErrorIs(t, err, types.ErrRewardPoolInsufficient)

	record, getErr := f.keeper.Rewardaccrual.Get(f.ctx, key)
	require.NoError(t, getErr)
	require.EqualValues(t, 50, record.Amount)
}
