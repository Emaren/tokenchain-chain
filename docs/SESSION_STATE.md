# TokenChain Session State

Updated: 2026-02-26
Workspace: /Users/tonyblum/projects/TokenChain

## Completed
- Cosmos SDK chain scaffolded with Ignite.
- Toolchain compatibility fixed in `go.mod`:
  - `github.com/bufbuild/buf => v1.50.0`
  - `github.com/golangci/golangci-lint => v1.64.8`
- Custom `x/loyalty` module added and wired.
- Implemented in `x/loyalty`:
  - creation mode params and authority checks
  - creator allowlist (authority-gated)
  - verified token registry
  - no-seizure default policy handling
  - strict full-denom minting (`factory/{issuer}/{subdenom}`) with on-chain hard cap enforcement
  - recovery policy hardening: `x/group` policy existence checks + network-aware timelock minimums
  - recovery execute authorization tightened (policy/authority only)
  - accrual and claim ledger
  - begin-block daily rollup boundary scheduler (`America/Edmonton`) with persisted last-rollup marker
  - daily rollup status query endpoint (`/tokenchain/loyalty/v1/daily_rollup/status`)
  - richer tx responses for create/mint/recovery operations (denom, minted supply, operation IDs/status/timestamps)
- Inflation suppression hook added in app init flow.
- CosmWasm runtime integration completed:
  - wasm keeper + module registration in app runtime (`app/wasm.go`)
  - wasm client-side module registration in CLI root wiring
  - wasm config template + start flags added (`app.toml` + start command)
  - pinned-code initialization on node load
- Genesis/testnet bootstrap updates:
  - fixed supply split in `config.yml`
  - local init script `scripts/init_local_testnet.sh`
  - local founder wasm upload allowlist patch in bootstrap flow
- Documentation updates:
  - `readme.md`
  - `docs/tokenchain-day1.md`

## In Progress
- TokenFactory parity work:
  - completed tokenfactory-style denom canonicalization + strict full-denom enforcement
  - remaining: native tokenfactory module parity and admin flow integration.

## Resume Plan
1. Design and wire TokenFactory module path that preserves locked Day-1 policy rules.
2. Add query/indexer surfaces for daily rollup status + decoded tx response helpers.
3. Run `go build ./...` and `go test ./...`.
4. Commit clean milestone and push.

## Checkpoint Process (enforced)
- Commit at each milestone.
- Push immediately after each milestone once remote exists.
- Keep this file updated before ending a session.
