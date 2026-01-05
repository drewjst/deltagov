#!/bin/sh
set -e

# Generate runtime config from environment variables
cat > /usr/share/nginx/html/assets/config.json << EOF
{
  "apiUrl": "${API_URL:-http://localhost:8080/api/v1}",
  "congressApiKey": "${CONGRESS_API_KEY:-}"
}
EOF

# Start nginx
exec nginx -g 'daemon off;'
