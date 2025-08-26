#!/usr/bin/env bash

set -euo pipefail

# This wrapper runs the full feature suite once and, if it fails,
# reruns only the failed specs by focusing Ginkgo on their names
# parsed from the generated JUnit XML under the reports directory.

FEATURES="${FEATURES:-}"
if [[ -z "${FEATURES}" ]]; then
  echo "FEATURES env var is required (e.g., FEATURES=operator)"
  exit 1
fi

echo "#### Running initial feature suite: ${FEATURES} ####"
if make test-features; then
  echo "✅ Initial feature run passed"
  exit 0
fi

echo "❌ Feature tests failed. Collecting failed specs from JUnit XML..."
REPORTS_DIR="${REPORTS_DIR:-reports}"
FAILED_NAMES=$(scripts/collect-failed-tests.sh "${REPORTS_DIR}" || true)

if [[ -z "${FAILED_NAMES}" ]]; then
  echo "No failed specs detected in ${REPORTS_DIR}. Exiting with failure."
  exit 1
fi

# Build a safe regex that matches any of the failed spec names
# Cap the number to avoid overlong command lines
max_specs=${MAX_RERUN_SPECS:-30}
count=0
focus_regex=""
while IFS= read -r name; do
  # Escape regex metacharacters
  esc=$(printf '%s' "${name}" | sed -e 's/[\^$.|?*+()\[\]{}]/\\&/g')
  if [[ -z "${focus_regex}" ]]; then
    focus_regex="${esc}"
  else
    focus_regex+="|${esc}"
  fi
  count=$((count+1))
  if [[ ${count} -ge ${max_specs} ]]; then
    break
  fi
done <<< "${FAILED_NAMES}"

if [[ -z "${focus_regex}" ]]; then
  echo "Could not construct a focus regex from failed specs. Exiting with failure."
  exit 1
fi

echo "🔁 Rerunning failed specs only with focus regex: ${focus_regex}"

# Build directories list for the selected feature(s)
dirs=""
for feature in ${FEATURES}; do
  for dir in tests/*; do
    if [[ ${dir} != *"util"* ]] && [[ ${dir} == *"${feature}"* ]]; then
      dirs+=" ${dir}"
    fi
  done
done
dirs=$(echo "${dirs}" | xargs || true)
if [[ -z "${dirs}" ]]; then
  echo "No test directories found for FEATURES='${FEATURES}'"
  exit 1
fi

# Compose flags similar to scripts/run-tests.sh
PFLAG=""
if [[ "${ENABLE_PARALLEL:-false}" == "true" ]]; then
  PFLAG="--procs=16"
fi

FFLAG=""
if [[ "${ENABLE_FLAKY_RETRY:-false}" == "true" ]]; then
  FFLAG="--flake-attempts=2"
fi

SEED_FLAG=""
if [[ -n "${GINKGO_SEED_NUMBER:-}" ]]; then
  SEED_FLAG="--seed=${GINKGO_SEED_NUMBER}"
fi

# Run the focused rerun
# shellcheck disable=SC2086
ginkgo -v ${PFLAG} ${FFLAG} --keep-going ${SEED_FLAG} \
  --output-interceptor-mode=none --timeout=24h --show-node-events --require-suite \
  --focus="${focus_regex}" ${dirs}


