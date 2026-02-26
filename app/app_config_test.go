package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoyaltyModuleAuthority(t *testing.T) {
	t.Setenv(loyaltyAuthorityEnvVar, "")
	require.Empty(t, loyaltyModuleAuthority())

	t.Setenv(loyaltyAuthorityEnvVar, "  tokenchain1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqf6vsh  ")
	require.Equal(t, "tokenchain1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqf6vsh", loyaltyModuleAuthority())
}
