# =============================================================================
# DeltaGov Frontend Dockerfile
# Multi-stage build for Angular app with Nginx
# =============================================================================

# -----------------------------------------------------------------------------
# Build Stage
# -----------------------------------------------------------------------------
FROM node:22-alpine AS builder

WORKDIR /build

# Copy package files first for better caching
COPY frontend/package.json frontend/package-lock.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY frontend/ ./

# Build for production
RUN npm run build -- --configuration=production

# -----------------------------------------------------------------------------
# Runtime Stage
# -----------------------------------------------------------------------------
FROM nginx:1.27-alpine

# Remove default nginx config
RUN rm -rf /etc/nginx/conf.d/*

# Copy custom nginx config
COPY deployments/docker/nginx.conf /etc/nginx/nginx.conf

# Copy built Angular app from builder
COPY --from=builder /build/dist/frontend/browser /usr/share/nginx/html

# Copy entrypoint script for runtime config injection
COPY frontend/docker-entrypoint.sh /docker-entrypoint.sh

# Create non-root user setup for Cloud Run compatibility
# Note: nginx.conf uses /tmp/nginx.pid for non-root compatibility
RUN mkdir -p /usr/share/nginx/html/assets && \
    chmod +x /docker-entrypoint.sh && \
    chown nginx:nginx /docker-entrypoint.sh && \
    chown -R nginx:nginx /usr/share/nginx/html && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /var/log/nginx && \
    touch /tmp/nginx.pid && \
    chown nginx:nginx /tmp/nginx.pid

# Use non-root user
USER nginx

# Environment variables for runtime configuration
ENV API_URL=http://localhost:8080/api/v1
ENV CONGRESS_API_KEY=

# Expose port (Cloud Run uses 8080)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/docker-entrypoint.sh"]
