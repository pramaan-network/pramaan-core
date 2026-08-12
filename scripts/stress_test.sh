#!/usr/bin/env bash
#
# PRAMAAN stress / stability harness.
#
# Exercises the chain code paths under adversarial and high-volume conditions
# and with the Go race detector enabled, to surface panics, data races, and
# state-index corruption before they reach a live network. Safe to run
# repeatedly; touches no chain data, keys, Docker, or git state.
#
# Usage:
#   scripts/stress_test.sh                 # default soak
#   STRESS_N=50000 scripts/stress_test.sh  # heavier document soak
#   ITER=5 scripts/stress_test.sh          # repeat the whole suite N times
#
set -euo pipefail

export PATH="$PATH:/usr/local/go/bin"
cd "$(dirname "$0")/.."

ITER="${ITER:-1}"
STRESS_N="${STRESS_N:-20000}"
export STRESS_N

echo "== PRAMAAN stress harness =="
echo "repo:      $(pwd)"
echo "go:        $(go version)"
echo "ITER:      $ITER"
echo "STRESS_N:  $STRESS_N"
echo

echo "== [1/4] build =="
go build ./...
echo "build OK"
echo

echo "== [2/4] vet =="
go vet ./...
echo "vet OK"
echo

for i in $(seq 1 "$ITER"); do
  echo "== [3/4] race-enabled full test suite (iteration $i/$ITER) =="
  go test -race -count=1 ./...
  echo "iteration $i OK"
  echo
done

echo "== [4/4] high-volume docreg state soak (STRESS_N=$STRESS_N, race) =="
go test -race -count=1 -timeout 30m -run TestStressSetDocument ./x/docreg/keeper/ -v
echo

echo "== stress harness completed with no failures =="
