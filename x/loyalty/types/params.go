package types

import (
	"fmt"
	"time"
)

const (
	CreationModeAdminOnly      = "admin_only"
	CreationModeAllowlisted    = "allowlisted"
	CreationModePermissionless = "permissionless"

	TotalBPS uint64 = 10_000
)

// DefaultCreationMode represents the CreationMode default value.
var DefaultCreationMode string = CreationModeAdminOnly

// DefaultDailyRollupTimezone represents the DailyRollupTimezone default value.
var DefaultDailyRollupTimezone string = "America/Edmonton"

// DefaultTestnetTimelockHours represents the TestnetTimelockHours default value.
var DefaultTestnetTimelockHours uint64 = 1

// DefaultMainnetTimelockHours represents the MainnetTimelockHours default value.
var DefaultMainnetTimelockHours uint64 = 24

// DefaultFeeSplitValidatorBps represents the FeeSplitValidatorBps default value.
var DefaultFeeSplitValidatorBps uint64 = 7000

// DefaultFeeSplitTokenStakersBps represents the FeeSplitTokenStakersBps default value.
var DefaultFeeSplitTokenStakersBps uint64 = 2000

// DefaultFeeSplitMerchantPoolBps represents the FeeSplitMerchantPoolBps default value.
var DefaultFeeSplitMerchantPoolBps uint64 = 1000

// DefaultSeizureOptInDefault represents the SeizureOptInDefault default value.
var DefaultSeizureOptInDefault bool = false

// DefaultMerchantIncentiveStakersBps represents the default per-token share of Bucket C routed to token stakers.
var DefaultMerchantIncentiveStakersBps uint64 = 5000

// DefaultMerchantIncentiveTreasuryBps represents the default per-token share of Bucket C routed to merchant treasury.
var DefaultMerchantIncentiveTreasuryBps uint64 = 5000

// NewParams creates a new Params instance.
func NewParams(
	creationMode string,
	dailyRollupTimezone string,
	testnetTimelockHours uint64,
	mainnetTimelockHours uint64,
	feeSplitValidatorBps uint64,
	feeSplitTokenStakersBps uint64,
	feeSplitMerchantPoolBps uint64,
	seizureOptInDefault bool,
) Params {
	return Params{
		CreationMode:            creationMode,
		DailyRollupTimezone:     dailyRollupTimezone,
		TestnetTimelockHours:    testnetTimelockHours,
		MainnetTimelockHours:    mainnetTimelockHours,
		FeeSplitValidatorBps:    feeSplitValidatorBps,
		FeeSplitTokenStakersBps: feeSplitTokenStakersBps,
		FeeSplitMerchantPoolBps: feeSplitMerchantPoolBps,
		SeizureOptInDefault:     seizureOptInDefault,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultCreationMode,
		DefaultDailyRollupTimezone,
		DefaultTestnetTimelockHours,
		DefaultMainnetTimelockHours,
		DefaultFeeSplitValidatorBps,
		DefaultFeeSplitTokenStakersBps,
		DefaultFeeSplitMerchantPoolBps,
		DefaultSeizureOptInDefault,
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := validateCreationMode(p.CreationMode); err != nil {
		return err
	}

	if err := validateDailyRollupTimezone(p.DailyRollupTimezone); err != nil {
		return err
	}

	if err := validateTestnetTimelockHours(p.TestnetTimelockHours); err != nil {
		return err
	}

	if err := validateMainnetTimelockHours(p.MainnetTimelockHours); err != nil {
		return err
	}

	if err := validateFeeSplitValidatorBps(p.FeeSplitValidatorBps); err != nil {
		return err
	}

	if err := validateFeeSplitTokenStakersBps(p.FeeSplitTokenStakersBps); err != nil {
		return err
	}

	if err := validateFeeSplitMerchantPoolBps(p.FeeSplitMerchantPoolBps); err != nil {
		return err
	}

	if err := validateSeizureOptInDefault(p.SeizureOptInDefault); err != nil {
		return err
	}

	if p.MainnetTimelockHours < p.TestnetTimelockHours {
		return fmt.Errorf("mainnet timelock must be greater than or equal to testnet timelock")
	}

	if p.FeeSplitValidatorBps+p.FeeSplitTokenStakersBps+p.FeeSplitMerchantPoolBps != TotalBPS {
		return fmt.Errorf("fee split bps must total %d", TotalBPS)
	}

	return nil
}

// validateCreationMode validates the CreationMode parameter.
func validateCreationMode(v string) error {
	switch v {
	case CreationModeAdminOnly, CreationModeAllowlisted, CreationModePermissionless:
		return nil
	default:
		return fmt.Errorf("invalid creation mode: %s", v)
	}
}

// validateDailyRollupTimezone validates the DailyRollupTimezone parameter.
func validateDailyRollupTimezone(v string) error {
	if v == "" {
		return fmt.Errorf("daily rollup timezone cannot be empty")
	}
	if _, err := time.LoadLocation(v); err != nil {
		return fmt.Errorf("invalid timezone %q: %w", v, err)
	}
	return nil
}

// validateTestnetTimelockHours validates the TestnetTimelockHours parameter.
func validateTestnetTimelockHours(v uint64) error {
	if v == 0 {
		return fmt.Errorf("testnet timelock must be greater than zero")
	}
	return nil
}

// validateMainnetTimelockHours validates the MainnetTimelockHours parameter.
func validateMainnetTimelockHours(v uint64) error {
	if v == 0 {
		return fmt.Errorf("mainnet timelock must be greater than zero")
	}
	return nil
}

// validateFeeSplitValidatorBps validates the FeeSplitValidatorBps parameter.
func validateFeeSplitValidatorBps(v uint64) error {
	if v > TotalBPS {
		return fmt.Errorf("validator fee split exceeds 100%%")
	}
	return nil
}

// validateFeeSplitTokenStakersBps validates the FeeSplitTokenStakersBps parameter.
func validateFeeSplitTokenStakersBps(v uint64) error {
	if v > TotalBPS {
		return fmt.Errorf("token stakers fee split exceeds 100%%")
	}
	return nil
}

// validateFeeSplitMerchantPoolBps validates the FeeSplitMerchantPoolBps parameter.
func validateFeeSplitMerchantPoolBps(v uint64) error {
	if v > TotalBPS {
		return fmt.Errorf("merchant pool fee split exceeds 100%%")
	}
	return nil
}

// validateSeizureOptInDefault validates the SeizureOptInDefault parameter.
func validateSeizureOptInDefault(v bool) error {
	_ = v
	return nil
}

// ValidateMerchantIncentiveRouting validates per-token Bucket C routing split.
func ValidateMerchantIncentiveRouting(stakersBps, treasuryBps uint64) error {
	if stakersBps > TotalBPS {
		return fmt.Errorf("merchant incentive stakers bps exceeds 100%%")
	}
	if treasuryBps > TotalBPS {
		return fmt.Errorf("merchant incentive treasury bps exceeds 100%%")
	}
	if stakersBps+treasuryBps != TotalBPS {
		return fmt.Errorf("merchant incentive routing bps must total %d", TotalBPS)
	}
	return nil
}
