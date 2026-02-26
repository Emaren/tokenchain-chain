package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/loyalty module sentinel errors
var (
	ErrInvalidSigner       = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidCreationMode = errors.Register(ModuleName, 1101, "invalid token creation mode")
	ErrCreatorNotAllowed   = errors.Register(ModuleName, 1102, "creator is not allowed to create tokens")
	ErrTokenExists         = errors.Register(ModuleName, 1103, "token already exists")
	ErrTokenNotFound       = errors.Register(ModuleName, 1104, "token not found")
	ErrInvalidDenom        = errors.Register(ModuleName, 1105, "invalid denom")
	ErrCapExceeded         = errors.Register(ModuleName, 1106, "mint would exceed max supply cap")
	ErrInvalidCap          = errors.Register(ModuleName, 1107, "invalid max supply cap")
	ErrRecoveryPolicy      = errors.Register(ModuleName, 1108, "invalid admin recovery policy")
	ErrAccrualNotFound     = errors.Register(ModuleName, 1109, "reward accrual not found")
)
