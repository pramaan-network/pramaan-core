// Package keeper (this file) implements the gRPC QueryServer for x/docreg:
// single-document lookups (by ID or hash), owner/issuer-indexed listings,
// module params, and a general filtered QueryDocuments query.
package keeper

import (
	"context"
	"errors"
	"pramaan/x/docreg/types"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// queryServer adapts a Keeper to the generated types.QueryServer interface.
type queryServer struct {
	k Keeper
}

// queryServer must be used as a pointer receiver type (methods below are
// defined on *queryServer), so the interface assertion and constructor both
// return a pointer.
var _ types.QueryServer = &queryServer{}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// backed by the given Keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{k: k}
}

// GetDocument handles the QueryGetDocumentRequest gRPC query: look up a
// single document by its ID. A not-found ID returns an empty response
// rather than an error (a normal "no such document" answer, not a fault).
func (q *queryServer) GetDocument(
	goCtx context.Context,
	req *types.QueryGetDocumentRequest,
) (*types.QueryGetDocumentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	doc, found := q.k.GetDocumentByID(ctx, req.Id)
	if !found {
		return &types.QueryGetDocumentResponse{}, nil
	}
	return &types.QueryGetDocumentResponse{Document: &doc}, nil
}

// Params handles the QueryParamsRequest gRPC query, returning the module's
// current on-chain parameters.
func (q *queryServer) Params(
	goCtx context.Context,
	req *types.QueryParamsRequest,
) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params, err := q.k.ParamsStore.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryParamsResponse{Params: &params}, nil
}

// DocumentsByOwner handles the QueryDocumentsByOwnerRequest gRPC query:
// returns every document currently owned by the given address, via the
// OwnerIndex maintained by SetDocument/TransferDocument.
func (q *queryServer) DocumentsByOwner(
	goCtx context.Context,
	req *types.QueryDocumentsByOwnerRequest,
) (*types.QueryDocumentsByOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var docs []*types.Document
	ownerList, err := q.k.OwnerIndex.Get(ctx, req.Owner)
	if err != nil {
		// Use errors.Is rather than == here: the collections API can wrap
		// ErrNotFound with additional context, so a direct equality check
		// (as this used to be) can miss a wrapped not-found error and
		// misreport a normal "no documents for this owner" as an Internal
		// gRPC error instead of an empty result. Same class of bug as
		// SECURITY_CHANGELOG.md #6 (docreg keeper.go), fixed here too.
		if errors.Is(err, collections.ErrNotFound) {
			return &types.QueryDocumentsByOwnerResponse{}, nil
		}
		return nil, err
	}
	for _, id := range ownerList.Items {
		doc, found := q.k.GetDocumentByID(ctx, id)
		if found {
			d := doc
			docs = append(docs, &d)
		}
	}
	return &types.QueryDocumentsByOwnerResponse{Documents: docs}, nil
}

// DocumentsByIssuer handles the QueryDocumentsByIssuerRequest gRPC query:
// returns every document issued by the given issuer address, via the
// IssuerIndex maintained by SetDocument.
func (q *queryServer) DocumentsByIssuer(
	goCtx context.Context,
	req *types.QueryDocumentsByIssuerRequest,
) (*types.QueryDocumentsByIssuerResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	var docs []*types.Document

	issuerList, err := q.k.IssuerIndex.Get(ctx, req.Issuer)
	if err != nil {
		// See the errors.Is note in DocumentsByOwner above — same fix.
		if errors.Is(err, collections.ErrNotFound) {
			return &types.QueryDocumentsByIssuerResponse{}, nil
		}
		return nil, err
	}

	for _, id := range issuerList.Items {
		doc, found := q.k.GetDocumentByID(ctx, id)
		if found {
			d := doc
			docs = append(docs, &d)
		}
	}

	return &types.QueryDocumentsByIssuerResponse{
		Documents: docs,
	}, nil
}

// GetDocumentByHash handles the QueryGetDocumentByHashRequest gRPC query:
// resolves a content hash to a document ID via HashIndex, then fetches the
// full document.
func (q *queryServer) GetDocumentByHash(
	goCtx context.Context,
	req *types.QueryGetDocumentByHashRequest,
) (*types.QueryGetDocumentByHashResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Step 1: resolve the hash to a document ID via the hash index.
	docID, err := q.k.HashIndex.Get(ctx, req.Hash)
	if err != nil {
		// See the errors.Is note in DocumentsByOwner above — same fix.
		if errors.Is(err, collections.ErrNotFound) {
			return &types.QueryGetDocumentByHashResponse{}, nil
		}
		return nil, err
	}

	// Step 2: fetch the full document using the resolved ID.
	doc, found := q.k.GetDocumentByID(ctx, docID)
	if !found {
		return &types.QueryGetDocumentByHashResponse{}, nil
	}

	return &types.QueryGetDocumentByHashResponse{
		Document: &doc,
	}, nil
}

// QueryDocuments handles the QueryDocumentsRequest gRPC query: the general
// filterable document search. If Owner or Issuer is set, it narrows the
// base dataset via the corresponding index first (cheaper than a full
// scan); otherwise it falls back to scanning every document. Type/Status/
// Owner/Issuer filters are then applied in-memory over that base dataset.
//
// NOTE: as with DocumentsByOwner/DocumentsByIssuer, a real error from
// OwnerIndex.Get/IssuerIndex.Get here (not just "not found") is currently
// swallowed into an empty response rather than propagated — lower severity
// than the direct-equality bug fixed above (this already treats every error
// as empty, so it doesn't misreport not-found specifically), but still
// means a transient store error would look like "zero documents" to a
// client instead of a query failure. Left as-is: fixing it well means
// distinguishing not-found from real errors here too, which changes this
// method's behavior on genuine backend errors — worth a deliberate follow-up
// rather than a drive-by change.
func (q *queryServer) QueryDocuments(
	goCtx context.Context,
	req *types.QueryDocumentsRequest,
) (*types.QueryDocumentsResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	var docs []*types.Document
	var candidateIDs []string

	// Step 1: choose the base dataset — narrow by owner or issuer index if
	// requested, otherwise fall back to a full scan of all documents.
	if req.Owner != "" {
		ownerList, err := q.k.OwnerIndex.Get(ctx, req.Owner)
		if err != nil {
			return &types.QueryDocumentsResponse{}, nil
		}
		candidateIDs = ownerList.Items

	} else if req.Issuer != "" {
		issuerList, err := q.k.IssuerIndex.Get(ctx, req.Issuer)
		if err != nil {
			return &types.QueryDocumentsResponse{}, nil
		}
		candidateIDs = issuerList.Items

	} else {
		// Fallback: no owner/issuer filter given, so scan every document.
		iter, err := q.k.Documents.Iterate(ctx, nil)
		if err != nil {
			return nil, err
		}
		defer iter.Close()

		for ; iter.Valid(); iter.Next() {
			doc, err := iter.Value()
			if err != nil {
				return nil, err
			}
			d := doc
			docs = append(docs, &d)
		}

		// Full-scan path already populated docs directly; skip the
		// ID-based fetch step below and go straight to filtering.
		goto FILTER
	}

	// Step 2: resolve the candidate IDs (from whichever index was used
	// above) into full documents.
	for _, id := range candidateIDs {
		doc, found := q.k.GetDocumentByID(ctx, id)
		if found {
			d := doc
			docs = append(docs, &d)
		}
	}

FILTER:

	// Step 3: apply the remaining Type/Status/Owner/Issuer filters over
	// whichever base dataset step 1/2 produced.
	var filtered []*types.Document

	for _, d := range docs {

		if req.Type != "" && d.Type != req.Type {
			continue
		}

		if req.Status != "" && d.Status != req.Status {
			continue
		}

		if req.Owner != "" && d.Owner != req.Owner {
			continue
		}

		if req.Issuer != "" && d.Issuer != req.Issuer {
			continue
		}

		filtered = append(filtered, d)
	}

	return &types.QueryDocumentsResponse{
		Documents: filtered,
	}, nil
}
