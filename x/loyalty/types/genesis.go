package types

import (
	"fmt"
	"time"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:                 DefaultParams(),
		CreatorallowlistMap:    []Creatorallowlist{},
		VerifiedtokenMap:       []Verifiedtoken{},
		RewardaccrualMap:       []Rewardaccrual{},
		MerchantallocationMap:  []Merchantallocation{},
		RecoveryoperationList:  []Recoveryoperation{},
		RecoveryoperationCount: 0,
	}
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
	merchantallocationIndexMap := make(map[string]struct{})
	for _, elem := range gs.MerchantallocationMap {
		index := fmt.Sprint(elem.Key)
		if _, ok := merchantallocationIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for merchantallocation")
		}
		merchantallocationIndexMap[index] = struct{}{}
		if elem.Date != "" {
			if _, err := time.Parse("2006-01-02", elem.Date); err != nil {
				return fmt.Errorf("invalid merchantallocation date: %w", err)
			}
		}
	}
	recoveryoperationIdMap := make(map[uint64]bool)
	recoveryoperationCount := gs.GetRecoveryoperationCount()
	for _, elem := range gs.RecoveryoperationList {
		if _, ok := recoveryoperationIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for recoveryoperation")
		}
		if elem.Id >= recoveryoperationCount {
			return fmt.Errorf("recoveryoperation id should be lower or equal than the last id")
		}
		recoveryoperationIdMap[elem.Id] = true
	}
	if gs.LastDailyRollupDate != "" {
		if _, err := time.Parse("2006-01-02", gs.LastDailyRollupDate); err != nil {
			return fmt.Errorf("invalid last daily rollup date: %w", err)
		}
	}

	return gs.Params.Validate()
}
