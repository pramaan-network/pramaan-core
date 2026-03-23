package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidSigner = errors.Register(
		ModuleName,
		1100,
		"expected gov account as only signer for proposal message",
	)

	ErrDocumentAlreadyExists = errors.Register(
		ModuleName,
		1101,
		"document hash already exists",
	)
)
