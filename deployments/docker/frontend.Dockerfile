# =============================================================================
# DeltaGov Frontend Dockerfile
# Multi-stage build with Angular 21 and Nginx Alpine
# =============================================================================
# Build: docker build -f deployments/docker/frontend.Dockerfile -t deltagov-frontend .
# Run: docker run -p 8080:8080 deltagov-frontend
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Dependencies
# -----------------------------------------------------------------------------
FROM node:22-alpine AS deps

WORKDIR /app

# Copy package files for dependency installation
COPY frontend/package.json frontend/package-lock.json ./

# Install dependencies with clean install for reproducibility
RUN npm ci --legacy-peer-deps

# -----------------------------------------------------------------------------
# Stage 2: Build
# -----------------------------------------------------------------------------
FROM node:22-alpine AS builder

WORKDIR /app

# Copy dependencies from deps stage
COPY --from=deps /app/node_modules ./node_modules

# Copy source code
COPY frontend/ .

# Build the Angular application for production
# Output goes to dist/frontend/browser for Angular 17+ with application builder
RUN npm run build -- --configuration=production

# -----------------------------------------------------------------------------
# Stage 3: Runtime (Nginx Alpine)
# -----------------------------------------------------------------------------
FROM nginx:1.27-alpine AS runtime

# Labels for container registry
LABEL org.opencontainers.image.source="https://github.com/drewjst/deltagov"
LABEL org.opencontainers.image.description="DeltaGov Frontend"
LABEL org.opencontainers.image.licenses="AGPL-3.0"

# Replace entire nginx config (not just server block)
# This config is optimized for non-root user
COPY deployments/docker/nginx.conf /etc/nginx/nginx.conf

# Copy built Angular app from builder stage
COPY --from=builder /app/dist/frontend/browser /usr/share/nginx/html

# Create non-root user and set permissions for all nginx paths
RUN adduser -D -g '' nginxuser && \
    chown -R nginxuser:nginxuser /usr/share/nginx/html && \
    chown -R nginxuser:nginxuser /var/cache/nginx && \
    chown -R nginxuser:nginxuser /var/log/nginx && \
    chown -R nginxuser:nginxuser /etc/nginx && \
    mkdir -p /tmp/client_temp /tmp/proxy_temp /tmp/fastcgi_temp /tmp/uwsgi_temp /tmp/scgi_temp && \
    chown -R nginxuser:nginxuser /tmp/*_temp

# Expose port 8080 (non-root compatible, Cloud Run default)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Run as non-root user
USER nginxuser

# Start nginx in foreground
CMD ["nginx", "-g", "daemon off;"]
