#!/usr/bin/env bash
set -euo pipefail

# Benchmark /api/v1/sentinel/projects/:id/details with simple timing summary.
#
# Usage examples:
#   ./api/scripts/bench_project_details.sh --base-url http://localhost:8080/api/v1 --project mims-hdmap-main --token "$TOKEN"
#   ./api/scripts/bench_project_details.sh --project mims-hdmap-main --token "$TOKEN" --runs 25 --tasks-limit 600
#
# Before/after workflow:
#   1) Run this script on old code and save output.
#   2) Deploy new code/migrations.
#   3) Run again with same args and compare p50/p95/avg.

BASE_URL="http://localhost:8080/api/v1"
PROJECT_ID_OR_CODE=""
TOKEN=""
RUNS=15
TASKS_LIMIT=600

while [[ $# -gt 0 ]]; do
  case "$1" in
    --base-url)
      BASE_URL="$2"; shift 2 ;;
    --project)
      PROJECT_ID_OR_CODE="$2"; shift 2 ;;
    --token)
      TOKEN="$2"; shift 2 ;;
    --runs)
      RUNS="$2"; shift 2 ;;
    --tasks-limit)
      TASKS_LIMIT="$2"; shift 2 ;;
    *)
      echo "Unknown arg: $1" >&2
      exit 1 ;;
  esac
done

if [[ -z "$PROJECT_ID_OR_CODE" ]]; then
  echo "Missing required --project <id-or-code>" >&2
  exit 1
fi
if [[ -z "$TOKEN" ]]; then
  echo "Missing required --token <jwt>" >&2
  exit 1
fi
if ! [[ "$RUNS" =~ ^[0-9]+$ ]] || [[ "$RUNS" -le 0 ]]; then
  echo "--runs must be positive integer" >&2
  exit 1
fi
if ! [[ "$TASKS_LIMIT" =~ ^[0-9]+$ ]] || [[ "$TASKS_LIMIT" -le 0 ]]; then
  echo "--tasks-limit must be positive integer" >&2
  exit 1
fi

URL="${BASE_URL}/sentinel/projects/${PROJECT_ID_OR_CODE}/details?tasks_limit=${TASKS_LIMIT}"

echo "Benchmarking: $URL"
echo "Runs: $RUNS"

declare -a TIMES_MS=()

for ((i=1; i<=RUNS; i++)); do
  OUT=$(curl -sS -o /tmp/bench_project_details_body.json -w "%{http_code} %{time_total}" \
    -H "Authorization: Bearer ${TOKEN}" \
    "$URL")

  HTTP_CODE=$(echo "$OUT" | awk '{print $1}')
  SEC=$(echo "$OUT" | awk '{print $2}')
  MS=$(awk -v s="$SEC" 'BEGIN { printf "%.3f", s*1000 }')

  if [[ "$HTTP_CODE" != "200" ]]; then
    echo "Run #$i failed: HTTP $HTTP_CODE" >&2
    head -c 1000 /tmp/bench_project_details_body.json >&2 || true
    echo >&2
    exit 1
  fi

  TIMES_MS+=("$MS")
  echo "#${i}: ${MS} ms"
done

SORTED=$(printf "%s\n" "${TIMES_MS[@]}" | sort -n)
AVG=$(printf "%s\n" "${TIMES_MS[@]}" | awk '{sum+=$1} END {if (NR==0) print 0; else printf "%.3f", sum/NR}')
MIN=$(printf "%s\n" "$SORTED" | head -n 1)
MAX=$(printf "%s\n" "$SORTED" | tail -n 1)

# Percentile helper (nearest-rank)
percentile() {
  local p="$1"
  local n
  n=$(printf "%s\n" "$SORTED" | wc -l | tr -d ' ')
  local rank
  rank=$(awk -v n="$n" -v p="$p" 'BEGIN { r=int((p/100)*n); if (r<1) r=1; if (r>n) r=n; print r }')
  printf "%s\n" "$SORTED" | sed -n "${rank}p"
}

P50=$(percentile 50)
P95=$(percentile 95)

echo ""
echo "Summary (ms):"
echo "  min:  $MIN"
echo "  p50:  $P50"
echo "  p95:  $P95"
echo "  avg:  $AVG"
echo "  max:  $MAX"

echo ""
echo "Tip: run the script before/after migration with same args and compare p50/p95/avg."
