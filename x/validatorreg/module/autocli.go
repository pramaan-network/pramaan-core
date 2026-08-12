// Package validatorreg (this file) implements the x/validatorreg CLI
// command wiring (autocli).
package validatorreg

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"pramaan/x/validatorreg/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface. Note
// the ApplyValidator/ApproveValidator/ActivateValidator proposal-flow
// messages (see keeper/msg_server.go) have no hand-tuned CLI entries here —
// only AddValidator/RemoveValidator (the direct, ROOT-only path) and the
// skipped UpdateParams do. The proposal-flow messages fall back to
// autocli's default auto-generated commands.
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
				{
					RpcMethod:      "AddValidator",
					Use:            "add-validator [address]",
					Short:          "Send a add-validator tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod:      "RemoveValidator",
					Use:            "remove-validator [address]",
					Short:          "Send a remove-validator tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
			},
		},
	}
}
