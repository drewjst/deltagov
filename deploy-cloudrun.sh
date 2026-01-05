#!/bin/bash
set -e

PROJECT_ID="deltagov"

echo "=== Deploying DeltaGov to Cloud Run ==="

echo ""
echo ">>> Deploying frontend..."
gcloud run deploy deltagov-frontend \
  --project=$PROJECT_ID \
  --image=gcr.io/deltagov/deltagov-frontend:latest \
  --region=us-central1 \
  --platform=managed \
  --allow-unauthenticated

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
  --set-env-vars=DATABASE_URL=postgresql://deltagov_user:PASSWORD@localhost:5432/deltagov,CONGRESS_API_KEY=YOUR_API_KEY \
  --add-cloudsql-instances=deltagov:us-central1:deltagov-postgres

echo ""
echo "=== Deployment complete ==="
