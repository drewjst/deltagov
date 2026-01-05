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

# Create non-root user setup for Cloud Run compatibility
RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /var/log/nginx && \
    touch /var/run/nginx.pid && \
    chown -R nginx:nginx /var/run/nginx.pid

# Use non-root user
USER nginx

# Expose port (Cloud Run uses 8080)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["nginx", "-g", "daemon off;"]
