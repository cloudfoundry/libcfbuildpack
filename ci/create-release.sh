#!/usr/bin/env bash

set -euo pipefail

RELEASE=$1

git tag -s v${RELEASE} -m "v${RELEASE}"
