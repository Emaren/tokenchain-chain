package keeper

import (
	"context"

	"tokenchain/x/loyalty/types"
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
	}
	for _, elem := range genState.RewardaccrualMap {
		if err := k.Rewardaccrual.Set(ctx, elem.Key, elem); err != nil {
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

	return genesis, nil
}
