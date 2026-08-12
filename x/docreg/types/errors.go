// Package types (this file) defines the x/docreg module's sentinel errors.
// Each is registered against ModuleName with a distinct code (1100+) so
// error codes are namespaced per-module and don't collide with codes
// registered by other modules in this chain.
package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

var (
	// ErrInvalidSigner is returned when a message's signer does not match
	// the expected authority (currently only relevant to UpdateParams).
	ErrInvalidSigner = errors.Register(
		ModuleName,
		1100,
		"expected gov account as only signer for proposal message",
	)

	// ErrDocumentAlreadyExists is returned when a document's content hash
	// (not its ID — see ErrDocumentIDExists) is already registered.
	ErrDocumentAlreadyExists = errors.Register(
		ModuleName,
		1101,
		"document hash already exists",
	)

	// ErrDocumentIDExists is returned when a document's caller-supplied ID
	// (not its content hash) is already registered — prevents
	// RegisterDocument from silently overwriting an existing record.
	ErrDocumentIDExists = errors.Register(
		ModuleName,
		1102,
		"document id already exists",
	)

	// ErrInvalidRequest is returned when a message is missing required
	// fields (e.g. RegisterDocument with an empty id/hash/owner/issuer, or
	// TransferDocument with an empty/unchanged new_owner).
	ErrInvalidRequest = errors.Register(
		ModuleName,
		1103,
		"invalid request: missing required fields",
	)

	// ErrDocumentNotFound is returned when a referenced document ID doesn't
	// exist (e.g. TransferDocument targeting an unknown ID).
	ErrDocumentNotFound = errors.Register(
		ModuleName,
		1104,
		"document not found",
	)

	// ErrNotDocumentOwner is returned when the message signer is not the
	// current owner of the document they're trying to act on.
	ErrNotDocumentOwner = errors.Register(
		ModuleName,
		1105,
		"signer is not the document owner",
	)

	// ErrDocumentNotTransferable is returned when TransferDocument targets a
	// document minted as a non-transferable SBT (soulbound token) rather
	// than a transferable NFT — see types.TokenTypeSBT/TokenTypeNFT.
	ErrDocumentNotTransferable = errors.Register(
		ModuleName,
		1106,
		"document is non-transferable (SBT)",
	)
)
