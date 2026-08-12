package keeper_test

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	"pramaan/x/authority/types"
)

func TestGenesis(t *testing.T) {
	f := initFixture(t)

	rootAddrStr, err := f.addressCodec.BytesToString(authtypes.NewModuleAddress("root-test"))
	require.NoError(t, err)

	params := types.DefaultParams()
	genesisState := types.GenesisState{
		Params: &params,
		Authorities: []*types.Authority{
			{Address: rootAddrStr, Role: types.RoleRoot},
		},
	}

	f.keeper.InitGenesis(f.ctx, genesisState)
	got := f.keeper.ExportGenesis(f.ctx)
	require.NotNil(t, got)

	require.EqualExportedValues(t, *genesisState.Params, *got.Params)
}
