package keeper

import (
	"context"
	"errors"

	"tokenchain/x/loyalty/types"

	"cosmossdk.io/collections"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	for _, elem := range genState.CreatorallowlistMap {
		if err := k.Creatorallowlist.Set(ctx, elem.Address, elem); err != nil {
			return err
		}
	}
	for _, elem := range genState.VerifiedtokenMap {
		if err := k.Verifiedtoken.Set(ctx, elem.Denom, elem); err != nil {
			return err
		}
		if err := k.setVerifiedTokenDenomMetadata(ctx, elem); err != nil {
			return err
		}
	}
	for _, elem := range genState.RewardaccrualMap {
		if err := k.Rewardaccrual.Set(ctx, elem.Key, elem); err != nil {
			return err
		}
	}
	for _, elem := range genState.RecoveryoperationList {
		if err := k.Recoveryoperation.Set(ctx, elem.Id, elem); err != nil {
			return err
		}
	}

	if err := k.RecoveryoperationSeq.Set(ctx, genState.RecoveryoperationCount); err != nil {
		return err
	}
	if genState.LastDailyRollupDate != "" {
		if err := k.LastDailyRollupDate.Set(ctx, genState.LastDailyRollupDate); err != nil {
			return err
		}
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	if err := k.Creatorallowlist.Walk(ctx, nil, func(_ string, val types.Creatorallowlist) (stop bool, err error) {
		genesis.CreatorallowlistMap = append(genesis.CreatorallowlistMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	if err := k.Verifiedtoken.Walk(ctx, nil, func(_ string, val types.Verifiedtoken) (stop bool, err error) {
		genesis.VerifiedtokenMap = append(genesis.VerifiedtokenMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	if err := k.Rewardaccrual.Walk(ctx, nil, func(_ string, val types.Rewardaccrual) (stop bool, err error) {
		genesis.RewardaccrualMap = append(genesis.RewardaccrualMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	err = k.Recoveryoperation.Walk(ctx, nil, func(key uint64, elem types.Recoveryoperation) (bool, error) {
		genesis.RecoveryoperationList = append(genesis.RecoveryoperationList, elem)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	genesis.RecoveryoperationCount, err = k.RecoveryoperationSeq.Peek(ctx)
	if err != nil {
		return nil, err
	}
	lastDailyRollupDate, err := k.LastDailyRollupDate.Get(ctx)
	if err == nil {
		genesis.LastDailyRollupDate = lastDailyRollupDate
	} else if !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}

	return genesis, nil
}
