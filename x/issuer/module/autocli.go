// Package issuer (this file) implements the x/issuer CLI command wiring
// (autocli).
package issuer

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"pramaan/x/issuer/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface. Note
// CreateIssuer/RevokeIssuer aren't given explicit CLI command options here
// (only Params query and the skipped UpdateParams tx are), so they fall
// back to autocli's default auto-generated command from the proto
// definition rather than a hand-tuned one.
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
