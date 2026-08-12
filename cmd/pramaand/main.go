// Command pramaand is the PRAMAAN chain's node daemon and CLI: it starts
// the validator/full node, and provides the query/tx/keys/genesis
// subcommands used to interact with the chain. See cmd/root.go for command
// tree construction and app/app.go for the application itself.
package main

import (
	"fmt"
	"os"

	clienthelpers "cosmossdk.io/client/v2/helpers"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"pramaan/app"
	"pramaan/cmd/pramaand/cmd"
)

// main builds the root cobra command and executes it, exiting non-zero on
// error (printed to the root command's configured stderr).
func main() {
	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, clienthelpers.EnvPrefix, app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
