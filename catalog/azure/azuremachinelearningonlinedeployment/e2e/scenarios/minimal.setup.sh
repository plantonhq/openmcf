#!/usr/bin/env bash
# SETUP-phase seed for the minimal scenario (see planton.dev/e2e-setup-script
# in minimal.yaml and the SETUP section of e2e/README.md).
#
# Azure ML will not provision a managed online deployment without a model:
# the ARM schema marks the model reference optional, but the SERVICE rejects
# a model-less create at deployment time (live-proven boundary; the exact
# error is recorded on the spec's `model` field). A registered model is
# data-plane content -- files uploaded into the fixture workspace -- which no
# catalog kind can create, so this script registers one here.
#
# The model is Microsoft's own canonical no-code-deployment sample (a ~750 B
# scikit-learn diabetes regressor in MLflow format), fetched from the
# azureml-examples repository PINNED to a commit and verified by sha256, so
# the lane's inputs can never drift. Azure builds the serving container from
# the model's own conda.yaml -- no scoring code, no environment asset.
#
# Idempotent: skips registration when the fixture workspace already holds the
# model version. Everything seeded here lives inside the fixture workspace
# and dies with it at DEPENDENCIES-DOWN.
set -euo pipefail

# The workspace's CLOUD-SIDE name (its fixture profile's spec.name), not the
# manifest metadata.name -- az ml addresses ARM objects.
RESOURCE_GROUP="planton-oss-e2e-azure-fixture-rg"
WORKSPACE="planton-oss-e2e-fixture-mlw"
MODEL_NAME="planton-oss-e2e-sklearn"
MODEL_VERSION="1"

# Pin: Azure/azureml-examples @ 2026-08 head of cli/endpoints/online/ncd/sklearn-diabetes/model
PIN_COMMIT="385f0590ec96839473d797f5663ce99fded920bf"
BASE_URL="https://raw.githubusercontent.com/Azure/azureml-examples/${PIN_COMMIT}/cli/endpoints/online/ncd/sklearn-diabetes/model"

sha256_of() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | cut -d' ' -f1
  else
    shasum -a 256 "$1" | cut -d' ' -f1
  fi
}

# Fetch one model file and verify its sha256. Written as a function
# (not an associative array) because the runner invokes `bash` and
# macOS ships bash 3.2, which has no `declare -A` -- a `set -u` run
# treats `[MLmodel]=` as an unbound variable and dies before any fetch.
fetch_and_verify() {
  local f="$1"
  local want="$2"
  curl -fsSL -o "${WORK_DIR}/${f}" "${BASE_URL}/${f}"
  local got
  got="$(sha256_of "${WORK_DIR}/${f}")"
  if [ "${got}" != "${want}" ]; then
    echo "[setup] checksum mismatch for ${f}: got ${got}" >&2
    exit 1
  fi
}

# The ml extension ships separately from the az CLI core.
if ! az extension show --name ml >/dev/null 2>&1; then
  echo "[setup] installing the az ml extension"
  az extension add --name ml --yes >/dev/null
fi

if az ml model show --name "${MODEL_NAME}" --version "${MODEL_VERSION}" \
  --resource-group "${RESOURCE_GROUP}" --workspace-name "${WORKSPACE}" >/dev/null 2>&1; then
  echo "[setup] model ${MODEL_NAME}:${MODEL_VERSION} already registered in ${WORKSPACE} -- nothing to seed"
  exit 0
fi

WORK_DIR="$(mktemp -d /tmp/planton-e2e-mlflow-model-XXXXXX)"
trap 'rm -rf "${WORK_DIR}"' EXIT

fetch_and_verify MLmodel e1bc4576b3e892a908919dcbf179219d135bed5eea01626323db48f8fde07802
fetch_and_verify conda.yaml 894424fba5a27f55fccddf335d69b073f9029a4efe6199023e43ee3ddb7ba849
fetch_and_verify model.pkl f1dd186b2cbeba89cef192493d6d7df9c55f629f9a5ce107e55fd2bdca79fe24
fetch_and_verify requirements.txt 83f19c94f2f64485f7c8dc306e870e17e52113cbd6b4bc304a6352e60bfe66f7
echo "[setup] fetched and verified 4 model files (pin ${PIN_COMMIT:0:8})"

az ml model create \
  --name "${MODEL_NAME}" \
  --version "${MODEL_VERSION}" \
  --type mlflow_model \
  --path "${WORK_DIR}" \
  --resource-group "${RESOURCE_GROUP}" \
  --workspace-name "${WORKSPACE}" \
  --only-show-errors >/dev/null

echo "[setup] registered ${MODEL_NAME}:${MODEL_VERSION} in ${WORKSPACE}"
