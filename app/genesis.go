// Package app (this file) defines the top-level GenesisState type: a raw
// map from module name to that module's own genesis JSON blob, keyed
// exactly as each AppModuleBasic.Name() returns.
package app

import (
	"encoding/json"
)

// GenesisState of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage
