package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"pramaan/x/authority/keeper"
	"pramaan/x/authority/types"
)

// TestMsgUpdateParams asserts that UpdateParams is intentionally disabled for
// the authority module: every call must return an error containing "disabled".
func TestMsgUpdateParams(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	params := types.DefaultParams()

	authorityStr, err := f.addressCodec.BytesToString(f.authority)
	require.NoError(t, err)

	testCases := []struct {
		name  string
		input *types.MsgUpdateParams
	}{
		{
			name: "disabled for invalid authority",
			input: &types.MsgUpdateParams{
				Authority: "invalid",
				Params:    params,
			},
		},
		{
			name: "disabled for valid authority",
			input: &types.MsgUpdateParams{
				Authority: authorityStr,
				Params:    params,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.UpdateParams(f.ctx, tc.input)
			require.Error(t, err)
			require.Contains(t, err.Error(), "disabled")
		})
	}
}
