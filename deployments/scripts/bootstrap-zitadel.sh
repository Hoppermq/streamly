#!/bin/bash
set -e

echo "🔐 Bootstrapping Zitadel service accounts..."

ZITADEL_URL="${ZITADEL_URL:-http://zitadel:8080}"
ZITADEL_DOMAIN="${ZITADEL_DOMAIN:-zitadel}"
ZITADEL_PORT="${ZITADEL_PORT:-8080}"
PAT_FILE="${ZITADEL_PAT_FILE:-/machinekey/zitadel-admin.pat}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TERRAFORM_DIR="${SCRIPT_DIR}/../zitadel/terraform"

# Wait for Zitadel
echo "⏳ Waiting for Zitadel..."
until curl -sf "${ZITADEL_URL}/debug/ready" >/dev/null 2>&1; do
  sleep 2
done
echo "✅ Zitadel is ready!"

# Wait for PAT file to be created by Zitadel init
echo "🔑 Waiting for PAT file..."
RETRY=0
until [ -f "$PAT_FILE" ] || [ $RETRY -ge 30 ]; do
  echo "   Waiting for $PAT_FILE... ($RETRY/30)"
  sleep 2
  RETRY=$((RETRY + 1))
done

if [ ! -f "$PAT_FILE" ]; then
  echo "❌ PAT file not found at $PAT_FILE"
  echo "   Zitadel may not have processed FirstInstance config."
  exit 1
fi

PAT=$(cat "$PAT_FILE")
echo "✅ PAT loaded from $PAT_FILE"

# Check for JWT profile file
JWT_FILE="/machinekey/zitadel-admin-sa.json"
if [ ! -f "$JWT_FILE" ]; then
  echo "❌ JWT file not found at $JWT_FILE"
  exit 1
fi
echo "✅ JWT profile file found at $JWT_FILE"

# Fetch the organization ID of the authenticated service account
echo "🔍 Fetching organization ID..."
ORG_RESPONSE=$(curl -s "${ZITADEL_URL}/management/v1/orgs/me" \
  -H "Authorization: Bearer ${PAT}" \
  -H "Content-Type: application/json")

if [ $? -ne 0 ]; then
  echo "❌ Failed to fetch organization ID"
  echo "   Response: $ORG_RESPONSE"
  exit 1
fi

ZITADEL_ORG_ID=$(echo "$ORG_RESPONSE" | jq -r '.org.id // empty')
if [ -z "$ZITADEL_ORG_ID" ]; then
  echo "❌ Failed to extract organization ID from response"
  echo "   Response: $ORG_RESPONSE"
  exit 1
fi

echo "✅ Organization ID: $ZITADEL_ORG_ID"

# Run Terraform
echo "🏗️  Provisioning service accounts with Terraform..."
cd "$TERRAFORM_DIR"

export TF_VAR_zitadel_domain="$ZITADEL_DOMAIN"
export TF_VAR_zitadel_port="$ZITADEL_PORT"
export TF_VAR_zitadel_secure_mode="true"
export TF_VAR_zitadel_jwt_profile_file="$JWT_FILE"
export TF_VAR_organization_id="$ZITADEL_ORG_ID"
export TF_VAR_project_name="local"

echo "🔍 Terraform variables:"
echo "  ZITADEL_DOMAIN: $ZITADEL_DOMAIN"
echo "  ZITADEL_PORT: $ZITADEL_PORT"
echo "  JWT_FILE: $JWT_FILE"
echo "  ORG_ID: $ZITADEL_ORG_ID"

echo "🧹 Cleaning old Terraform state..."
rm -rf .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup

# Configure Terraform plugin cache
export TF_PLUGIN_CACHE_DIR="/root/.terraform.d/plugin-cache"
mkdir -p "$TF_PLUGIN_CACHE_DIR"
echo "📦 Using Terraform plugin cache at $TF_PLUGIN_CACHE_DIR"

# Configure GitHub token for provider downloads if available
if [ -n "${GITHUB_TOKEN}" ]; then
  echo "🔑 Using GitHub token for provider downloads..."
  git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"
fi

# Retry terraform init with exponential backoff
MAX_RETRIES=5
RETRY=0
while [ $RETRY -lt $MAX_RETRIES ]; do
  echo "🔄 Terraform init attempt $((RETRY + 1))/$MAX_RETRIES..."
  if terraform init -reconfigure; then
    echo "✅ Terraform initialized successfully"
    break
  else
    RETRY=$((RETRY + 1))
    if [ $RETRY -lt $MAX_RETRIES ]; then
      WAIT_TIME=$((2 ** RETRY))
      echo "⏳ Retrying in ${WAIT_TIME}s..."
      sleep $WAIT_TIME
    else
      echo "❌ Terraform init failed after $MAX_RETRIES attempts"
      exit 1
    fi
  fi
done
terraform apply -auto-approve

# Save credentials
terraform output -json service_credentials 2>/dev/null | jq -r '
  to_entries[] |
  "# \(.key) service\n\(.key | ascii_upcase)_CLIENT_ID=\(.value.client_id)\n\(.key | ascii_upcase)_CLIENT_SECRET=\(.value.client_secret)\n"
' > "${SCRIPT_DIR}/../.env.zitadel"

echo "✅ Service accounts provisioned!"
echo "📄 Credentials saved to .env.zitadel"
