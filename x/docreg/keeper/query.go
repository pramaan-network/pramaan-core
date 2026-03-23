package keeper

import (
	"context"
	"pramaan/x/docreg/types"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type queryServer struct {
	k Keeper
}

// Aa line ma pointer (&) hovo jaruri chhe
var _ types.QueryServer = &queryServer{}

// Aa constructor pointer return karvo joie
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{k: k}
}

// --- GetDocument ---
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

// --- Params ---
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

// --- DocumentsByOwner ---
func (q *queryServer) DocumentsByOwner(
	goCtx context.Context,
	req *types.QueryDocumentsByOwnerRequest,
) (*types.QueryDocumentsByOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var docs []*types.Document
	ownerList, err := q.k.OwnerIndex.Get(ctx, req.Owner)
	if err != nil {
		if err == collections.ErrNotFound {
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

func (q *queryServer) DocumentsByIssuer(
	goCtx context.Context,
	req *types.QueryDocumentsByIssuerRequest,
) (*types.QueryDocumentsByIssuerResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	var docs []*types.Document

	issuerList, err := q.k.IssuerIndex.Get(ctx, req.Issuer)
	if err != nil {
		if err == collections.ErrNotFound {
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

func (q *queryServer) GetDocumentByHash(
	goCtx context.Context,
	req *types.QueryGetDocumentByHashRequest,
) (*types.QueryGetDocumentByHashResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// 🔥 Step 1: Get document ID from hash index
	docID, err := q.k.HashIndex.Get(ctx, req.Hash)
	if err != nil {
		if err == collections.ErrNotFound {
			return &types.QueryGetDocumentByHashResponse{}, nil
		}
		return nil, err
	}

	// 🔥 Step 2: Get full document using ID
	doc, found := q.k.GetDocumentByID(ctx, docID)
	if !found {
		return &types.QueryGetDocumentByHashResponse{}, nil
	}

	return &types.QueryGetDocumentByHashResponse{
		Document: &doc,
	}, nil
}

func (q *queryServer) QueryDocuments(
	goCtx context.Context,
	req *types.QueryDocumentsRequest,
) (*types.QueryDocumentsResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	var docs []*types.Document
	var candidateIDs []string

	// 🔥 STEP 1: Choose base dataset
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
		// 🔥 fallback: scan all documents
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

		// apply filters later
		goto FILTER
	}

	// 🔥 STEP 2: fetch documents from IDs
	for _, id := range candidateIDs {
		doc, found := q.k.GetDocumentByID(ctx, id)
		if found {
			d := doc
			docs = append(docs, &d)
		}
	}

FILTER:

	// 🔥 STEP 3: apply filters
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
