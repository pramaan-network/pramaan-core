package keeper_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"pramaan/x/docreg/types"
)

// stressN returns the number of documents the stress test inserts. Default is
// kept modest so the test stays fast in CI; set STRESS_N to crank it up for a
// heavier soak run (e.g. STRESS_N=20000). Note the per-owner OwnerIndex is a
// single serialized list re-written on every append, so runtime grows roughly
// O(n^2) in n for a single owner — see SECURITY_CHANGELOG.md.
func stressN() int {
	if v := os.Getenv("STRESS_N"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 5000
}

// TestStressSetDocument hammers the docreg state layer with a large batch of
// unique documents that all share a single owner and issuer, then verifies:
//   - every document is retrievable by ID,
//   - the OwnerIndex/IssuerIndex lists stay consistent (one entry per doc, no
//     drops, no duplicates) as they grow,
//   - duplicate-ID and duplicate-hash writes are rejected rather than
//     silently overwriting or corrupting the indexes.
//
// It is deliberately allocation-heavy and index-growth-heavy so it doubles as
// a memory/consistency stress test under `go test -race`.
func TestStressSetDocument(t *testing.T) {
	f := initFixture(t)

	n := stressN()
	const owner = "pramaan1stressowner000000000000000000000000000"
	const issuer = "pramaan1stressissuer00000000000000000000000000"

	for i := 0; i < n; i++ {
		doc := types.Document{
			Id:           fmt.Sprintf("doc-%d", i),
			Hash:         fmt.Sprintf("hash-%d", i),
			Owner:        owner,
			Issuer:       issuer,
			Type:         "education",
			Status:       types.DocumentStatusIssued,
			TokenType:    types.TokenTypeSBT,
			Transferable: false,
		}
		if err := f.keeper.SetDocument(f.ctx, doc); err != nil {
			t.Fatalf("SetDocument(%d) failed: %v", i, err)
		}
	}

	// Spot-check retrieval across the range.
	for _, i := range []int{0, 1, n / 2, n - 1} {
		id := fmt.Sprintf("doc-%d", i)
		if _, found := f.keeper.GetDocumentByID(f.ctx, id); !found {
			t.Fatalf("expected document %s to exist", id)
		}
	}

	// OwnerIndex and IssuerIndex must each hold exactly n IDs with no dupes.
	ownerList, err := f.keeper.OwnerIndex.Get(f.ctx, owner)
	if err != nil {
		t.Fatalf("OwnerIndex.Get: %v", err)
	}
	if len(ownerList.Items) != n {
		t.Fatalf("OwnerIndex length = %d, want %d", len(ownerList.Items), n)
	}
	seen := make(map[string]struct{}, n)
	for _, id := range ownerList.Items {
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate id %s in OwnerIndex", id)
		}
		seen[id] = struct{}{}
	}

	issuerList, err := f.keeper.IssuerIndex.Get(f.ctx, issuer)
	if err != nil {
		t.Fatalf("IssuerIndex.Get: %v", err)
	}
	if len(issuerList.Items) != n {
		t.Fatalf("IssuerIndex length = %d, want %d", len(issuerList.Items), n)
	}

	// Duplicate ID must be rejected.
	dupID := types.Document{Id: "doc-0", Hash: "brand-new-hash", Owner: owner, Issuer: issuer}
	if err := f.keeper.SetDocument(f.ctx, dupID); err == nil {
		t.Fatalf("expected duplicate-ID SetDocument to fail")
	}

	// Duplicate hash (new ID) must be rejected.
	dupHash := types.Document{Id: "brand-new-id", Hash: "hash-0", Owner: owner, Issuer: issuer}
	if err := f.keeper.SetDocument(f.ctx, dupHash); err == nil {
		t.Fatalf("expected duplicate-hash SetDocument to fail")
	}

	// A failed duplicate write must not have grown the indexes.
	ownerList2, err := f.keeper.OwnerIndex.Get(f.ctx, owner)
	if err != nil {
		t.Fatalf("OwnerIndex.Get (post-dup): %v", err)
	}
	if len(ownerList2.Items) != n {
		t.Fatalf("OwnerIndex length changed after rejected writes: %d, want %d", len(ownerList2.Items), n)
	}
}
