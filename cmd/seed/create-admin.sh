#!/bin/bash

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
USERNAME="admin"
EMAIL="admin@quizme.com"
PASSWORD="admin123"
FULLNAME="Administrator"
UPDATE_FLAG=""

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    -u|--username)
      USERNAME="$2"
      shift 2
      ;;
    -e|--email)
      EMAIL="$2"
      shift 2
      ;;
    -p|--password)
      PASSWORD="$2"
      shift 2
      ;;
    -n|--fullname)
      FULLNAME="$2"
      shift 2
      ;;
    --update)
      UPDATE_FLAG="-update"
      shift
      ;;
    -h|--help)
      show_help
      exit 0
      ;;
    *)
      echo -e "${RED}Unknown option: $1${NC}"
      show_help
      exit 1
      ;;
  esac
done

show_help() {
  cat << EOF
${BLUE}QuizMe Admin User Seeding Script${NC}

Usage: ./create-admin.sh [options]

Options:
  -u, --username USERNAME    Admin username (default: admin)
  -e, --email EMAIL          Admin email (default: admin@quizme.com)
  -p, --password PASSWORD    Admin password (default: admin123)
  -n, --fullname FULLNAME    Admin full name (default: Administrator)
  --update                   Update existing admin user if exists
  -h, --help                 Show this help message

Examples:
  # Create admin with default values
  ./create-admin.sh

  # Create admin with custom values
  ./create-admin.sh -u superadmin -e super@quizme.com -p Pass123

  # Update existing admin password
  ./create-admin.sh -u admin -p NewPassword123 --update

EOF
}

# Display info
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  QuizMe Admin User Seeding            ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

echo -e "${YELLOW}Configuration:${NC}"
echo "  Username:  $USERNAME"
echo "  Email:     $EMAIL"
echo "  Password:  $(printf '*%.0s' $(seq 1 ${#PASSWORD}))"
echo "  Full Name: $FULLNAME"
if [ -n "$UPDATE_FLAG" ]; then
  echo "  Mode:      UPDATE (if exists)"
else
  echo "  Mode:      CREATE (new)"
fi
echo ""

# Get the script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BACKEND_DIR="$( cd "$SCRIPT_DIR/.." && pwd )"

# Change to backend directory
cd "$BACKEND_DIR" || { echo -e "${RED}Failed to change directory${NC}"; exit 1; }

echo -e "${YELLOW}Running seed command...${NC}"
echo ""

# Build command
CMD="go run cmd/seed/main.go -username=$USERNAME -email=$EMAIL -password=$PASSWORD -fullname=$FULLNAME"
if [ -n "$UPDATE_FLAG" ]; then
  CMD="$CMD $UPDATE_FLAG"
fi

# Run the command
eval "$CMD"

EXIT_CODE=$?

if [ $EXIT_CODE -eq 0 ]; then
  echo ""
  echo -e "${GREEN}✓ Admin user setup completed successfully!${NC}"
  echo -e "${BLUE}You can now log in with:${NC}"
  echo "  Username/Email: $USERNAME (or $EMAIL)"
  echo "  Password:       $(printf '*%.0s' $(seq 1 ${#PASSWORD}))"
  echo ""
else
  echo ""
  echo -e "${RED}✗ Admin user setup failed (Exit code: $EXIT_CODE)${NC}"
  exit $EXIT_CODE
fi
