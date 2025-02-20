#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Setup logging
timestamp() {
    date "+%Y-%m-%d %H:%M:%S"
}

log() {
    echo -e "[$(timestamp)] $1"
}

log_error() {
    echo -e "${RED}[$(timestamp)] ERROR: $1${NC}"
}

log_success() {
    echo -e "${GREEN}[$(timestamp)] SUCCESS: $1${NC}"
}

log_info() {
    echo -e "${BLUE}[$(timestamp)] INFO: $1${NC}"
}

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    log_error "GitHub CLI (gh) is not installed. Please install it first."
    log_info "Visit: https://cli.github.com/"
    exit 1
fi

# Check if logged in to GitHub
if ! gh auth status &> /dev/null; then
    log_info "Please login to GitHub CLI:"
    gh auth login
    log "GitHub CLI authentication completed"
fi

# Create a new token with necessary scopes
log_info "Creating new Personal Access Token..."
gh auth refresh --scopes "repo,write:packages"
TOKEN=$(gh auth token)

if [ -z "$TOKEN" ]; then
    log_error "Failed to generate token"
    exit 1
fi

# Get current repository
REPO=$(gh repo view --json nameWithOwner -q .nameWithOwner)

if [ -z "$REPO" ]; then
    log_error "Not in a GitHub repository directory"
    exit 1
fi

# Set the secret
log_info "Setting GITHUB_TOKEN secret for $REPO..."
echo "$TOKEN" | gh secret set GITHUB_TOKEN

if [ $? -eq 0 ]; then
    log_success "Successfully set GITHUB_TOKEN secret!"
    log_info "Note: For most GitHub Actions workflows, you don't need to set this manually"
    log_info "as GitHub automatically provides GITHUB_TOKEN."
else
    log_error "Failed to set secret"
    exit 1
fi 