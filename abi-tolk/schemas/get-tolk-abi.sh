#!/usr/bin/env bash

TOLK_PATH=$1
ABI_PATH=$2

TEMPFILE=$(mktemp)
pnpm dlx "file:$(pwd)/ton-tolk-js-tmp-206.tgz" --output-json "$TEMPFILE" "$TOLK_PATH"

jq '.abiJson' "$TEMPFILE" > "$ABI_PATH"
rm "$TEMPFILE"