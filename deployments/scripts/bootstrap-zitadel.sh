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
until curl -sf -H "Host: ${ZITADEL_HOST}" "${ZITADEL_URL}/debug/ready" >/dev/null 2>&1; do
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

# Set project name for checking
PROJECT_NAME="${TF_VAR_project_name:-local}"

# Check if resources already exist by querying the project
echo "🔍 Checking if resources already exist in Zitadel..."
PROJECT_CHECK=$(curl -s "${ZITADEL_URL}/management/v1/projects/_search" \
  -H "Authorization: Bearer ${PAT}" \
  -H "Content-Type: application/json" \
  -d "{\"queries\":[{\"nameQuery\":{\"name\":\"$PROJECT_NAME\",\"method\":\"TEXT_QUERY_METHOD_EQUALS\"}}]}" | jq -r '.result[0].id // empty' 2>/dev/null)

if [ -n "$PROJECT_CHECK" ]; then
  echo "✅ Bootstrap already completed - project '$PROJECT_NAME' exists (ID: $PROJECT_CHECK)"
  echo "   Bootstrap skipped - services already configured in Zitadel"
  echo "   To force re-run: Delete project in Zitadel UI or run 'docker-compose down -v' to clean everything"
  exit 0
fi

# Check Terraform state as secondary check
if [ -f "terraform.tfstate" ]; then
  RESOURCE_COUNT=$(cat terraform.tfstate | jq -r '.resources | length' 2>/dev/null || echo "0")
  if [ "$RESOURCE_COUNT" -gt 0 ]; then
    echo "✅ Terraform state has $RESOURCE_COUNT resources"
  else
    echo "⚠️  Terraform state exists but empty - will attempt to provision"
  fi
else
  echo "🆕 No Terraform state found - first run"
fi

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

# Note: We keep Terraform state to make this idempotent
# Only clean state if explicitly needed (e.g., corrupted state)
if [ "$CLEAN_STATE" = "true" ]; then
  echo "🧹 Cleaning Terraform state (CLEAN_STATE=true)..."
  rm -rf .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup
fi

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

# Save service credentials
terraform output -json service_credentials 2>/dev/null | jq -r '
  to_entries[] |
  "# \(.key) service\n\(.key | ascii_upcase)_CLIENT_ID=\(.value.client_id)\n\(.key | ascii_upcase)_CLIENT_SECRET=\(.value.client_secret)\n"
' > "${SCRIPT_DIR}/../.env.zitadel"

# Save root admin credentials
echo "🔑 Exporting root admin credentials..."
ROOT_USER_ID=$(terraform output -json root_admin_credentials 2>/dev/null | jq -r '.user_id')
ROOT_USER_PAT=$(terraform output -json root_admin_credentials 2>/dev/null | jq -r '.pat')

if [ -n "$ROOT_USER_ID" ] && [ "$ROOT_USER_ID" != "null" ]; then
  echo "$ROOT_USER_ID" > "${SCRIPT_DIR}/../zitadel/machinekey/root-user.id"
  echo "✅ Root user ID saved to zitadel/machinekey/root-user.id"
else
  echo "⚠️  Could not export root user ID"
fi

if [ -n "$ROOT_USER_PAT" ] && [ "$ROOT_USER_PAT" != "null" ]; then
  echo "$ROOT_USER_PAT" > "${SCRIPT_DIR}/../zitadel/machinekey/root-user.pat"
  echo "✅ Root user PAT saved to zitadel/machinekey/root-user.pat"
else
  echo "⚠️  Could not export root user PAT"
fi

echo "✅ Service accounts provisioned!"
echo "📄 Service credentials saved to .env.zitadel"
echo "📄 Root admin credentials saved to zitadel/machinekey/"
echo "🎉 Bootstrap complete! Terraform state saved."
