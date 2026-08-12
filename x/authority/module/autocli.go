// Package authority implements the x/authority Cosmos SDK module: the
// PRAMAAN role registry (ROOT / AUTHORITY / VALIDATOR / ISSUER) that other
// modules (docreg, issuer, validatorreg) consult to authorize privileged
// actions. This file wires up the CLI (autocli) commands for the module.
package authority

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"pramaan/x/authority/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface. It
// declares the `pramaand query authority ...` and `pramaand tx authority ...`
// CLI surface for this module: the query side gets an auto-generated
// `params` command, and on the tx side `UpdateParams` is explicitly skipped
// because it is authority-gated (see keeper/msg_update_params.go — updates
// are disabled entirely for this module).
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
			},
		},
	}
}
