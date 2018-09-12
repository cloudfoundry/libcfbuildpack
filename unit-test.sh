#!/usr/bin/env bash

set -euo pipefail

[[ -d $PWD/go-build && ! -d $HOME/.cache/go-build ]] && mkdir -p $HOME/.cache && ln -s $PWD/go-build $HOME/.cache/go-build

cd libjavabuildpack
go test -v
