# Multi-stage build for wild-cloud-central
FROM golang:1.21-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wild-cloud-central .

# Final stage with dnsmasq and nginx
FROM debian:bookworm-slim

# Install required packages
RUN apt-get update && apt-get install -y \
    dnsmasq \
    nginx \
    systemctl \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create necessary directories
RUN mkdir -p /etc/wild-cloud-central \
    /var/lib/wild-cloud-central \
    /var/log/wild-cloud-central \
    /var/www/html/talos \
    /var/ftpd

# Create wildcloud user
RUN useradd -r -s /bin/false wildcloud

# Copy built application
COPY --from=builder /app/wild-cloud-central /usr/bin/wild-cloud-central
COPY --from=builder /app/config.yaml /etc/wild-cloud-central/config.yaml
COPY --from=builder /app/static /var/www/html/wild-central
COPY --from=builder /app/wild-cloud-central.service /etc/systemd/system/

# Set permissions
RUN chown -R wildcloud:wildcloud /var/lib/wild-cloud-central /var/log/wild-cloud-central
RUN chmod +x /usr/bin/wild-cloud-central

# Configure nginx
COPY nginx-container.conf /etc/nginx/sites-available/wild-central
RUN ln -s /etc/nginx/sites-available/wild-central /etc/nginx/sites-enabled/ && \
    rm -f /etc/nginx/sites-enabled/default

# Expose ports
EXPOSE 8081 53/udp 67/udp 69/udp 80

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8081/api/v1/health || exit 1

# Start script
COPY start-container.sh /start-container.sh
RUN chmod +x /start-container.sh

CMD ["/start-container.sh"]