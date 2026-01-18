#!/bin/bash
# Wrapper script to check for keys and generate if missing

set -e

KEY_PATH=${GOOGLE_PRIVATE_KEY_PATH:?GOOGLE_PRIVATE_KEY_PATH is not set}

if [ -f "$KEY_PATH" ]; then
    echo "Key file $KEY_PATH already exists. Skipping generation."
else
    echo "Key file $KEY_PATH not found. Attempting to generate..."
    
    # Run the setup script
    # We assume we are in the root or where scripts folder is accessible
    # The setup script expects to run from project root usually
    
    ./scripts/setup-service-account.sh staging
fi
