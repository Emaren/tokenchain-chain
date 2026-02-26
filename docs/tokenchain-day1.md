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
- `x/wasm` (CosmWasm runtime with permissioned upload default)
- IBC transfer + ICA (legacy/manual wiring)

## Loyalty Module Highlights (`x/loyalty`)
- Parameterized creation policy:
  - `admin_only`
  - `allowlisted`
  - `permissionless`
- Authority defaults to `x/gov`, with optional runtime override via `TOKENCHAIN_LOYALTY_AUTHORITY`
- Authority-gated allowlist management
- Verified token registry:
  - metadata
  - bank denom metadata publication (wallet/explorer-friendly units/name/symbol)
  - max supply cap
  - recovery policy metadata
  - tokenfactory-style denom format: `factory/{issuer}/{subdenom}`
- No-seizure default: `seizure_opt_in_default=false`
- Opt-in recovery execution flow:
  - recovery policy address must exist in `x/group` (not a free-form string)
  - issuer is immutable once token is created
  - seizure opt-in cannot be turned on after minting begins
  - `queue-recovery-transfer` (policy/authority only)
  - timelock enforced on-chain
  - timelock minimum follows network mode:
    - testnet => `testnet_timelock_hours`
    - mainnet => `mainnet_timelock_hours`
  - `execute-recovery-transfer` (policy/authority only, after timelock)
  - `cancel-recovery-transfer` (policy/authority only)
- Hard cap enforcement in mint path (`mint-verified-token`)
- Reward accrual ledger and claim path:
  - authority can record accruals
  - users claim accrued balances on-chain
  - begin-block daily rollup boundary fires once per Edmonton local day and emits `loyalty_daily_rollup`
  - query endpoint exposes rollup status for dashboards: `/tokenchain/loyalty/v1/daily_rollup/status`
  - query endpoint exposes filtered reward accruals by address/denom: `/tokenchain/loyalty/v1/rewardaccruals/filter`
  - query endpoint exposes filtered recovery operations: `/tokenchain/loyalty/v1/recoveryoperations/filter`
- Daily rollup timezone parameter default: `America/Edmonton`
- Fee split parameter defaults: `7000/2000/1000` bps

## Explicitly Deferred
- Full TokenFactory module parity
- Osmosis relayer automation and production ops manifests

These deferred items are intentionally separated from the current pass so the baseline chain compiles, tests pass, and the policy layer is enforceable now.
