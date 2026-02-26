package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:              DefaultParams(),
		CreatorallowlistMap: []Creatorallowlist{}, VerifiedtokenMap: []Verifiedtoken{}, RewardaccrualMap: []Rewardaccrual{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	creatorallowlistIndexMap := make(map[string]struct{})

	for _, elem := range gs.CreatorallowlistMap {
		index := fmt.Sprint(elem.Address)
		if _, ok := creatorallowlistIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for creatorallowlist")
		}
		creatorallowlistIndexMap[index] = struct{}{}
	}
	verifiedtokenIndexMap := make(map[string]struct{})

	for _, elem := range gs.VerifiedtokenMap {
		index := fmt.Sprint(elem.Denom)
		if _, ok := verifiedtokenIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for verifiedtoken")
		}
		verifiedtokenIndexMap[index] = struct{}{}
	}
	rewardaccrualIndexMap := make(map[string]struct{})

	for _, elem := range gs.RewardaccrualMap {
		index := fmt.Sprint(elem.Key)
		if _, ok := rewardaccrualIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for rewardaccrual")
		}
		rewardaccrualIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
