package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "vcoa/testutil/keeper"
	"vcoa/x/vnft/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.VnftKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
