package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"pramaan/x/authority/types"
)

func TestGenesisState_Validate(t *testing.T) {
	validAddr := sdk.AccAddress(make([]byte, 20)).String()
	validAddr2 := sdk.AccAddress([]byte{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	}).String()

	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc: "valid with single ROOT",
			genState: &types.GenesisState{
				Authorities: []*types.Authority{
					{Address: validAddr, Role: types.RoleRoot},
				},
			},
			valid: true,
		},
		{
			desc: "valid with ROOT and AUTHORITY",
			genState: &types.GenesisState{
				Authorities: []*types.Authority{
					{Address: validAddr, Role: types.RoleRoot},
					{Address: validAddr2, Role: types.RoleAuthority},
				},
			},
			valid: true,
		},
		{
			desc:     "empty authorities",
			genState: &types.GenesisState{},
			valid:    false,
		},
		{
			desc: "no ROOT authority",
			genState: &types.GenesisState{
				Authorities: []*types.Authority{
					{Address: validAddr, Role: types.RoleAuthority},
				},
			},
			valid: false,
		},
		{
			desc: "invalid bech32 address",
			genState: &types.GenesisState{
				Authorities: []*types.Authority{
					{Address: "not-a-bech32", Role: types.RoleRoot},
				},
			},
			valid: false,
		},
		{
			desc: "duplicate addresses",
			genState: &types.GenesisState{
				Authorities: []*types.Authority{
					{Address: validAddr, Role: types.RoleRoot},
					{Address: validAddr, Role: types.RoleAuthority},
				},
			},
			valid: false,
		},
		{
			desc: "invalid role",
			genState: &types.GenesisState{
				Authorities: []*types.Authority{
					{Address: validAddr, Role: "SUPERADMIN"},
				},
			},
			valid: false,
		},
		{
			desc: "empty address",
			genState: &types.GenesisState{
				Authorities: []*types.Authority{
					{Address: "", Role: types.RoleRoot},
				},
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
