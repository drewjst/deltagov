#!/bin/bash
set -e

PROJECT_ID="deltagov"

echo "=== Deploying DeltaGov to Cloud Run ==="

BACKEND_URL="https://deltagov-backend-222145696718.us-central1.run.app"

echo ""
echo ">>> Deploying frontend..."
gcloud run deploy deltagov-frontend \
  --project=$PROJECT_ID \
  --image=gcr.io/deltagov/deltagov-frontend:latest \
  --region=us-central1 \
  --platform=managed \
  --allow-unauthenticated \
  --set-env-vars="API_URL=${BACKEND_URL}/api/v1,CONGRESS_API_KEY=${CONGRESS_API_KEY:-}"

echo ""
echo ">>> Deploying backend..."
gcloud run deploy deltagov-backend \
  --project=$PROJECT_ID \
  --image=gcr.io/deltagov/deltagov-backend:latest \
  --region=us-central1 \
  --platform=managed \
  --allow-unauthenticated \
  --memory=2Gi \
  --cpu=2 \
  --set-env-vars="DATABASE_URL=${DATABASE_URL:-postgresql://deltagov_user:PASSWORD@localhost:5432/deltagov},CONGRESS_API_KEY=${CONGRESS_API_KEY:-}" \
  --add-cloudsql-instances=deltagov:us-central1:deltagov-postgres

echo ""
echo "=== Deployment complete ==="
