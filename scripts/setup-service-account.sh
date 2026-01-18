#!/bin/bash

# Script to set up service account key for signed URLs
# Usage: ./scripts/setup-service-account.sh [staging|production]

ENVIRONMENT=${1:-staging}

if [ "$ENVIRONMENT" != "staging" ] && [ "$ENVIRONMENT" != "production" ]; then
    echo "Error: Environment must be 'staging' or 'production'"
    echo "Usage: $0 [staging|production]"
    exit 1
fi

echo "Setting up service account key for $ENVIRONMENT environment..."

# Set project based on environment
if [ "$ENVIRONMENT" = "staging" ]; then
    PROJECT_ID="munene-dev-staging"
elif [ "$ENVIRONMENT" = "production" ]; then
    PROJECT_ID="munene-dev"
fi

# Check if we can access the project
if ! gcloud projects describe "$PROJECT_ID" >/dev/null 2>&1; then
    echo "Error: Cannot access project $PROJECT_ID"
    echo "Make sure you have the necessary permissions and are authenticated"
    echo "Run: gcloud auth login"
    exit 1
fi

# Create keys directory if it doesn't exist
mkdir -p keys

# Service account name
SERVICE_ACCOUNT_NAME="tuneslap-api-$ENVIRONMENT"
SERVICE_ACCOUNT_EMAIL="$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com"
KEY_FILE="keys/$SERVICE_ACCOUNT_NAME-key.json"

echo "Project ID: $PROJECT_ID"
echo "Service Account: $SERVICE_ACCOUNT_EMAIL"
echo "Key File: $KEY_FILE"

# Check if service account exists
if ! gcloud iam service-accounts describe "$SERVICE_ACCOUNT_EMAIL" --project="$PROJECT_ID" >/dev/null 2>&1; then
    echo "Error: Service account $SERVICE_ACCOUNT_EMAIL does not exist in project $PROJECT_ID"
    echo "Make sure you have run 'terraform apply' in the $ENVIRONMENT directory"
    exit 1
fi

# Generate service account key
echo "Generating service account key..."
gcloud iam service-accounts keys create "$KEY_FILE" \
    --iam-account="$SERVICE_ACCOUNT_EMAIL" \
    --project="$PROJECT_ID"

if [ $? -eq 0 ]; then
    echo "Service account key created successfully!"
    echo ""
    
    # Get the absolute path to the key file
    ABSOLUTE_KEY_PATH=$(realpath "$KEY_FILE")
    
    echo "Add these to your environment variables:"
    echo "export GOOGLE_PRIVATE_KEY_PATH=\"$ABSOLUTE_KEY_PATH\""
    echo "# Optional: Override service account email (will use key file email if not set)"
    echo "export GOOGLE_SERVICE_ACCOUNT_EMAIL=\"$SERVICE_ACCOUNT_EMAIL\""
    echo ""
    echo "Or add to your .env file:"
    echo "GOOGLE_PRIVATE_KEY_PATH=$ABSOLUTE_KEY_PATH"
    echo "# Optional: Override service account email (will use key file email if not set)"
    echo "GOOGLE_SERVICE_ACCOUNT_EMAIL=$SERVICE_ACCOUNT_EMAIL"
    echo ""
    echo "Note: Keep the service account key file secure and don't commit it to version control!"
    echo "The key file is located at: $ABSOLUTE_KEY_PATH"
    echo "The service account email ($SERVICE_ACCOUNT_EMAIL) will be automatically extracted from the key file if not set via environment variable."
else
    echo "Error: Failed to generate service account key"
    echo "Make sure you have the necessary permissions in project $PROJECT_ID"
    exit 1
fi 