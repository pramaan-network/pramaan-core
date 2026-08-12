// Package docreg (this file) implements the x/docreg CLI command wiring
// (autocli): what shows up under `pramaand query docreg ...` and
// `pramaand tx docreg ...`.
package docreg

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"pramaan/x/docreg/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface,
// declaring the query commands (params, get-document — note the richer
// query RPCs like DocumentsByOwner/DocumentsByIssuer/GetDocumentByHash/
// QueryDocuments are not exposed here, only reachable via gRPC/REST
// directly) and tx commands (register-document, transfer-document;
// update-params is skipped since it's authority-gated).
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
				{
					RpcMethod:      "GetDocument",
					Use:            "get-document [id]",
					Short:          "Query get-document",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
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
					RpcMethod:      "RegisterDocument",
					Use:            "register-document [id] [hash] [owner] [issuer] [doc-type] [metadata]",
					Short:          "Send a register-document tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "hash"}, {ProtoField: "owner"}, {ProtoField: "issuer"}, {ProtoField: "doc_type"}, {ProtoField: "metadata"}},
				},
				{
					RpcMethod:      "TransferDocument",
					Use:            "transfer-document [id] [new-owner]",
					Short:          "Send a transfer-document tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "new_owner"}},
				},
			},
		},
	}
}
