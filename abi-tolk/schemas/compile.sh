#!/usr/bin/env bash
set -euo pipefail

schemas_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

usage() {
  echo "Usage: compile.sh [--all] | <input.tolk> [output.json]" >&2
}

compile_file() {
  local tolk_path="$1"
  local abi_path="${2:-${tolk_path%.tolk}.json}"
  local out_dir log_file code_boc64 json_file

  out_dir="$(dirname "$abi_path")"
  if [[ -n "$out_dir" && "$out_dir" != "." ]]; then
    mkdir -p "$out_dir"
  fi

  log_file="$(mktemp)"
  if acton compile "$tolk_path" --abi "$abi_path" >"$log_file" 2>&1; then
    cat "$log_file"
  else
    cat "$log_file" >&2
    rm -f "$log_file"
    return 1
  fi

  code_boc64="$(sed -n 's/^Code in base64: //p' "$log_file" | tail -n 1)"
  rm -f "$log_file"

  if [[ -n "$code_boc64" ]]; then
    if ! command -v jq >/dev/null 2>&1; then
      echo "warning: jq not found; code_boc64 was not added to $abi_path" >&2
      return
    fi

    json_file="$(mktemp)"
    jq --arg code_boc64 "$code_boc64" '.code_boc64 = $code_boc64' "$abi_path" >"$json_file"
    mv "$json_file" "$abi_path"
  fi
}

collect_tolk_files() {
  local entry child

  shopt -s nullglob
  for entry in "$schemas_dir"/*; do
    [[ "$(basename "$entry")" == .* ]] && continue

    if [[ -d "$entry" ]]; then
      for child in "$entry"/*.tolk; do
        [[ -f "$child" ]] && printf '%s\n' "$child"
      done
    elif [[ -f "$entry" && "$entry" == *.tolk ]]; then
      printf '%s\n' "$entry"
    fi
  done
}

compile_all() {
  local files=()
  local file out_path

  while IFS= read -r file; do
    files+=("$file")
  done < <(collect_tolk_files)

  if [[ "${#files[@]}" -eq 0 ]]; then
    echo "No .tolk files found."
    return
  fi

  echo "Found ${#files[@]} .tolk file(s)."
  for file in "${files[@]}"; do
    out_path="${file%.tolk}.json"
    compile_file "$file" "$out_path"
    echo "  ${file#"$schemas_dir"/} -> ${out_path#"$schemas_dir"/}"
  done
}

if [[ "$#" -eq 0 || "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  usage
  if [[ "$#" -eq 0 ]]; then
    exit 1
  fi
  exit 0
fi

if [[ "$1" == "--all" ]]; then
  compile_all
  exit 0
fi

tolk_path="$1"
abi_path="${2:-}"
compile_file "$tolk_path" "$abi_path"
echo "ABI written to ${abi_path:-${tolk_path%.tolk}.json}"
