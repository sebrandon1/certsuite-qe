#!/usr/bin/env bash

# Collect failed Ginkgo test names from JUnit XML reports.
# Usage: scripts/collect-failed-tests.sh <reports_dir>

set -euo pipefail

REPORTS_DIR="${1:-reports}"

if [[ ! -d "${REPORTS_DIR}" ]]; then
  exit 0
fi

# Extract testcase name attributes that have a <failure> in the same testcase block.
# Portable awk approach (no PCRE dependency):
awk '
  /<testcase/{
    name=""
    if (match($0, /name="([^"]+)"/, m)) { name=m[1] }
    open=1
  }
  /<failure/ && open==1 {
    if (name != "") { print name }
    open=0
    name=""
  }
  /<\/testcase/ { open=0; name="" }
' "${REPORTS_DIR}"/*.xml 2>/dev/null | sort -u




