#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${NEW_API_BASE_URL:-}" ]]; then
  echo "NEW_API_BASE_URL is required, e.g. http://127.0.0.1:3000"
  exit 1
fi

if [[ -z "${NEW_API_ADMIN_TOKEN:-}" ]]; then
  echo "NEW_API_ADMIN_TOKEN is required"
  exit 1
fi

if [[ -z "${OPENROUTER_FALLBACK_KEY:-}" ]]; then
  echo "OPENROUTER_FALLBACK_KEY is required"
  exit 1
fi

payload=$(cat <<JSON
{
  "name": "OpenRouter Free Fallback",
  "type": 20,
  "key": "${OPENROUTER_FALLBACK_KEY}",
  "keys": ["${OPENROUTER_FALLBACK_KEY}"],
  "models": "meta-llama/llama-3.1-8b-instruct:free,google/gemma-2-9b-it:free,mistralai/mistral-7b-instruct:free",
  "group": "default",
  "status": 1,
  "priority": 999,
  "weight": 0,
  "is_fallback": true
}
JSON
)

curl -fsSL -X POST "${NEW_API_BASE_URL%/}/api/channel/" \
  -H "Authorization: Bearer ${NEW_API_ADMIN_TOKEN}" \
  -H 'Content-Type: application/json' \
  -d "${payload}"

echo
echo "Fallback channel initialized."
