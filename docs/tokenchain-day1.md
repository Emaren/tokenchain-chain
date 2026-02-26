# TokenChain Day-1 Scope (Implemented)

## Chain Identity
- Address prefix: `tokenchain`
- Base denom: `utoken`
- Coin type: `118`
- Testnet chain-id target: `tokenchain-testnet-1`

## Economic Defaults
- Total TOKEN supply target: `1,000,000` (in config bootstrap split)
- Founder allocation: `900,000`
- Treasury allocation: `100,000`
- Inflation: disabled at genesis init by forcing `x/mint` inflation fields to `0`

## Enabled Core Modules
- `x/auth`, `x/bank`, `x/staking`, `x/distribution`, `x/slashing`, `x/gov`
- `x/upgrade`
- `x/circuit`
- `x/feegrant`
- `x/authz`
- `x/group`
- IBC transfer + ICA (legacy/manual wiring)

## Loyalty Module Highlights (`x/loyalty`)
- Parameterized creation policy:
  - `admin_only`
  - `allowlisted`
  - `permissionless`
- Authority-gated allowlist management
- Verified token registry:
  - metadata
  - max supply cap
  - recovery policy metadata
- No-seizure default: `seizure_opt_in_default=false`
- Hard cap enforcement in mint path (`mint-verified-token`)
- Reward accrual ledger and claim path:
  - authority can record accruals
  - users claim accrued balances on-chain
- Daily rollup timezone parameter default: `America/Edmonton`
- Fee split parameter defaults: `7000/2000/1000` bps

## Explicitly Deferred
- Native CosmWasm runtime keeper wiring (`x/wasm`)
- Full TokenFactory module parity
- Automated end-block daily reward rollup scheduler
- Osmosis relayer automation and production ops manifests

These deferred items are intentionally separated from the current pass so the baseline chain compiles, tests pass, and the policy layer is enforceable now.
