// Package app (this file) sets this chain's global SDK config at process
// startup: the default bond denom and every Bech32 address prefix
// (account/validator-operator/validator-consensus, each with its matching
// pubkey prefix), all derived from AccountAddressPrefix ("pramaan", see
// app.go). Config.Seal() locks these in for the lifetime of the process —
// they cannot be changed after init() runs.
package app

import sdk "github.com/cosmos/cosmos-sdk/types"

// init sets the chain's bond denom and Bech32 prefixes before anything else
// in the app package runs, since address encoding/decoding throughout the
// rest of the app depends on this global config being set first.
func init() {
	// Set bond denom

	sdk.DefaultBondDenom = "stake"

	// Set address prefixes
	accountPubKeyPrefix := AccountAddressPrefix + "pub"
	validatorAddressPrefix := AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := AccountAddressPrefix + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetCoinType(ChainCoinType)
	config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()
}
