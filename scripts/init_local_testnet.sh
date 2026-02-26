#!/usr/bin/env bash

set -euo pipefail

BINARY="${BINARY:-tokenchaind}"
CHAIN_ID="${CHAIN_ID:-tokenchain-testnet-1}"
HOME_DIR="${HOME_DIR:-$HOME/.tokenchain}"
KEYRING_BACKEND="${KEYRING_BACKEND:-test}"
MONIKER="${MONIKER:-founder}"
RESET_HOME="${RESET_HOME:-1}"

if [[ "${RESET_HOME}" == "1" ]]; then
  rm -rf "${HOME_DIR}"
fi

echo "Initializing ${BINARY} home at ${HOME_DIR}"
"${BINARY}" init "${MONIKER}" --chain-id "${CHAIN_ID}" --home "${HOME_DIR}"

echo "Creating keys"
"${BINARY}" keys add founder --keyring-backend "${KEYRING_BACKEND}" --home "${HOME_DIR}" --output json >/dev/null 2>&1 || true
"${BINARY}" keys add treasury --keyring-backend "${KEYRING_BACKEND}" --home "${HOME_DIR}" --output json >/dev/null 2>&1 || true
FOUNDER_ADDR="$("${BINARY}" keys show founder -a --keyring-backend "${KEYRING_BACKEND}" --home "${HOME_DIR}")"

echo "Configuring fixed-supply genesis accounts"
"${BINARY}" genesis add-genesis-account founder 900000000000utoken \
  --keyring-backend "${KEYRING_BACKEND}" --home "${HOME_DIR}"
"${BINARY}" genesis add-genesis-account treasury 100000000000utoken \
  --keyring-backend "${KEYRING_BACKEND}" --home "${HOME_DIR}"

echo "Creating validator gentx"
"${BINARY}" genesis gentx founder 50000000utoken \
  --chain-id "${CHAIN_ID}" \
  --keyring-backend "${KEYRING_BACKEND}" \
  --home "${HOME_DIR}"

"${BINARY}" genesis collect-gentxs --home "${HOME_DIR}"

echo "Setting wasm upload policy for local testnet founder"
GENESIS_FILE="${HOME_DIR}/config/genesis.json"
TMP_GENESIS="$(mktemp)"
jq --arg founder "${FOUNDER_ADDR}" \
  '.app_state.wasm.params.code_upload_access = {"permission":"AnyOfAddresses","addresses":[$founder]} |
   .app_state.wasm.params.instantiate_default_permission = "Everybody"' \
  "${GENESIS_FILE}" >"${TMP_GENESIS}"
mv "${TMP_GENESIS}" "${GENESIS_FILE}"

"${BINARY}" genesis validate-genesis --home "${HOME_DIR}"

TREASURY_ADDR="$("${BINARY}" keys show treasury -a --keyring-backend "${KEYRING_BACKEND}" --home "${HOME_DIR}")"

cat <<EOF
TokenChain local testnet initialized.

Chain ID: ${CHAIN_ID}
Home: ${HOME_DIR}
Founder: ${FOUNDER_ADDR}
Treasury: ${TREASURY_ADDR}

Start node:
  ${BINARY} start --home "${HOME_DIR}"
EOF
