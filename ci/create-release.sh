#!/usr/bin/env sh

set -euo pipefail

RELEASE=$1

git tag -s v${RELEASE} -m "v${RELEASE}"
