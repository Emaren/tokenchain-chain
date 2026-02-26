# TokenChain

TokenChain is a Cosmos SDK chain for TokenTap loyalty infrastructure.

This repository now includes a working chain baseline with:
- `bech32` prefix: `tokenchain`
- base denom: `utoken` (display `TOKEN`, 6 decimals)
- coin type: `118`
- no inflation at init (mint module inflation fields forced to zero during genesis init)
- IBC transfer + ICA scaffolding enabled
- CosmWasm runtime (`x/wasm`) integrated in app, CLI, and config wiring
- governance-safe ops modules enabled: `x/upgrade`, `x/circuit`, `x/feegrant`, `x/authz`, `x/group`
- optional loyalty authority override via `TOKENCHAIN_LOYALTY_AUTHORITY` (defaults to `x/gov` if unset)

## Loyalty Module (`x/loyalty`)

`x/loyalty` is the first on-chain TokenChain business-logic module and enforces:
- creation mode policy (`admin_only`, `allowlisted`, `permissionless`)
- creator allowlist (authority-gated)
- verified business token registry with metadata and cap
- automatic bank denom metadata publication for verified tokens (create/update/genesis import)
- tokenfactory-style business denom canonicalization (`factory/{issuer}/{subdenom}`)
- no-seizure default (`seizure_opt_in_default=false`)
- optional recovery policy metadata (`recovery_group_policy`, timelock hours)
- recovery policy hardening:
  - `recovery_group_policy` must resolve to an existing `x/group` policy account
  - timelock minimum is enforced per network (`testnet_timelock_hours` vs `mainnet_timelock_hours`)
  - issuer is immutable after token creation
  - seizure/recovery cannot be enabled after minting has begun
- on-chain recovery operation queue with timelock execution:
  - `queue-recovery-transfer`
  - `execute-recovery-transfer` (policy/authority gated)
  - `cancel-recovery-transfer`
- hard-cap mint checks (`mint-verified-token` cannot exceed `max_supply`)
- on-chain accrual ledger (`rewardaccrual`) and user claim flow (`claim-reward`)
- enriched tx responses for `record-reward-accrual` and `claim-reward` (amounts, denom, key/date)
- automatic daily rollup boundary in begin-block using `America/Edmonton`, with on-chain rollup marker persistence
- daily rollup status query (`/tokenchain/loyalty/v1/daily_rollup/status`) for dashboard/indexer consumption
- reward accrual filter query (`/tokenchain/loyalty/v1/rewardaccruals/filter`) by address/denom + pagination
- recovery operations filter query (`/tokenchain/loyalty/v1/recoveryoperations/filter`) with status/token/address + pagination
- daily rollup timezone param default: `America/Edmonton`
- timelock params defaults: testnet `1h`, mainnet `24h`
- fee split params defaults: `7000/2000/1000` bps (validator / token stakers / merchant pool)

## Genesis Defaults

`config.yml` has been updated for the fixed 1,000,000 TOKEN supply split:
- founder: `900000000000utoken` (900,000 TOKEN)
- treasury: `100000000000utoken` (100,000 TOKEN)

## Quick Start

```bash
ignite chain serve
```

Or build/run directly:

```bash
go build ./cmd/tokenchaind
./tokenchaind start
```

## Local Bootstrap Script

Use the included script for deterministic local testnet initialization:

```bash
./scripts/init_local_testnet.sh
```

This script initializes:
- chain-id `tokenchain-testnet-1`
- founder + treasury keys
- founder validator gentx
- fixed-supply genesis accounts
- wasm upload access pinned to founder on local testnet bootstrap

## Notes

- Go `1.24+` is required.
- Default chain genesis policy keeps wasm uploads permissioned (`Nobody`) until explicitly opened by governance.
- If you want founder-operated day-1 loyalty admin flows on testnet, set `TOKENCHAIN_LOYALTY_AUTHORITY=<founder-address>` on every validator.
