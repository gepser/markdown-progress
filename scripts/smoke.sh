#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${BASE_URL:-}" ]]; then
  echo "BASE_URL is required (example: https://geps.dev)"
  exit 1
fi

BASE_URL="${BASE_URL%/}"
PROGRESS_PATH="${PROGRESS_PATH:-/progress}"

if [[ -n "$PROGRESS_PATH" && "$PROGRESS_PATH" != /* ]]; then
  PROGRESS_PATH="/${PROGRESS_PATH}"
fi

progress_url() {
  local suffix="$1"
  if [[ -z "$PROGRESS_PATH" ]]; then
    echo "${BASE_URL}${suffix}"
    return
  fi

  echo "${BASE_URL}${PROGRESS_PATH}${suffix}"
}

expect_status() {
  local expected="$1"
  local method="$2"
  local url="$3"

  local code
  code="$(curl -sS -o /dev/null -w '%{http_code}' -X "$method" "$url")"
  if [[ "$code" != "$expected" ]]; then
    echo "Expected $expected for $method $url, got $code"
    exit 1
  fi
}

expect_header_contains() {
  local url="$1"
  local header="$2"
  local expected_value="$3"

  local headers
  headers="$(curl -sSI "$url" | tr -d '\r')"
  if ! grep -Eiq "^${header}:" <<<"$headers"; then
    echo "Header check failed for $url"
    echo "Expected header: ${header}"
    echo "Got headers:"
    echo "$headers"
    exit 1
  fi

  if ! grep -Ei "^${header}:" <<<"$headers" | grep -Fqi "$expected_value"; then
    echo "Header check failed for $url"
    echo "Expected: ${header} contains '${expected_value}'"
    echo "Got headers:"
    echo "$headers"
    exit 1
  fi
}

expect_body_contains() {
  local url="$1"
  local expected="$2"

  local body
  body="$(curl -sS "$url")"
  if ! grep -Fq "$expected" <<<"$body"; then
    echo "Body check failed for $url"
    echo "Expected body to contain: $expected"
    exit 1
  fi
}

echo "Running smoke tests against ${BASE_URL}"

expect_status "200" "GET" "$(progress_url "/76")"
expect_header_contains "$(progress_url "/76")" "Content-Type" "image/svg+xml"
expect_header_contains "$(progress_url "/76")" "Cache-Control" "max-age=300"
expect_body_contains "$(progress_url "/76")" "76%"

expect_status "400" "GET" "$(progress_url "/not-a-number")"
expect_status "400" "GET" "$(progress_url "/50?successColor=nothex")"
expect_status "405" "POST" "$(progress_url "/50")"

expect_status "200" "GET" "$(progress_url "/150")"
expect_body_contains "$(progress_url "/150")" "100%"

echo "Smoke tests passed."
