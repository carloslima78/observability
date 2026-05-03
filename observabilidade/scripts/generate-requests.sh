#!/bin/bash

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
ERROR_CODES="${ERROR_CODES:-500 400 404}"
SLEEP_SECONDS="${SLEEP_SECONDS:-1}"

while true; do
  for error_code in $ERROR_CODES; do
    curl -s -o /dev/null -X GET "$BASE_URL/metrics/error/$error_code"
    echo "Requisicao enviada para: $BASE_URL/metrics/error/$error_code"

    curl -s -o /dev/null -X GET "$BASE_URL/metrics/latency?seconds=0.3"
    echo "Requisicao enviada para: $BASE_URL/metrics/latency?seconds=0.3"

    curl -s -o /dev/null -X POST "$BASE_URL/metrics/orders"
    echo "Requisicao enviada para: $BASE_URL/metrics/orders"

    sleep "$SLEEP_SECONDS"
  done
done
