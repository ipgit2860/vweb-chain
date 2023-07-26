package vnft_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "vcoa/testutil/keeper"
	"vcoa/testutil/nullify"
	"vcoa/x/vnft"
	"vcoa/x/vnft/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.VnftKeeper(t)
	vnft.InitGenesis(ctx, *k, genesisState)
	got := vnft.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
