#!/bin/bash

# Let's Encrypt SSL Certificate Management Script
# This script automates the process of obtaining and configuring Let's Encrypt certificates

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
CERTBOT_VERSION="latest"
CERT_DIR="../ssl"
NGINX_DIR="../nginx"
BACKUP_DIR="../ssl/backup"

# Print colored message
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        print_error "Please run as root (use sudo)"
        exit 1
    fi
}

# Check if domain is provided
check_domain() {
    if [ -z "$1" ]; then
        print_error "Domain name is required"
        echo "Usage: sudo $0 your-domain.com"
        exit 1
    fi
}

# Check if certbot is installed
check_certbot() {
    if ! command -v certbot &> /dev/null; then
        print_warning "Certbot is not installed. Installing..."
        install_certbot
    else
        print_info "Certbot is installed: $(certbot --version)"
    fi
}

# Install certbot
install_certbot() {
    if [ -f /etc/debian_version ]; then
        # Debian/Ubuntu
        apt-get update
        apt-get install -y certbot
    elif [ -f /etc/redhat-release ]; then
        # RHEL/CentOS
        yum install -y certbot
    else
        print_error "Unsupported OS. Please install certbot manually."
        exit 1
    fi
    print_info "Certbot installed successfully"
}

# Create certificate directory
create_cert_dir() {
    if [ ! -d "$CERT_DIR" ]; then
        mkdir -p "$CERT_DIR"
        print_info "Created certificate directory: $CERT_DIR"
    fi
}

# Backup existing certificates
backup_existing_certs() {
    if [ -f "$CERT_DIR/fullchain.pem" ] || [ -f "$CERT_DIR/privkey.pem" ]; then
        print_info "Backing up existing certificates..."
        mkdir -p "$BACKUP_DIR"
        TIMESTAMP=$(date +%Y%m%d_%H%M%S)
        
        if [ -f "$CERT_DIR/fullchain.pem" ]; then
            cp "$CERT_DIR/fullchain.pem" "$BACKUP_DIR/fullchain.pem.$TIMESTAMP"
        fi
        if [ -f "$CERT_DIR/privkey.pem" ]; then
            cp "$CERT_DIR/privkey.pem" "$BACKUP_DIR/privkey.pem.$TIMESTAMP"
        fi
        
        print_info "Certificates backed up to: $BACKUP_DIR"
    fi
}

# Obtain certificate using standalone mode (port 80 must be free)
obtain_cert_standalone() {
    DOMAIN=$1
    
    print_info "Stopping Nginx temporarily for standalone mode..."
    docker-compose -f ../docker-compose.yml stop nginx || true
    
    print_info "Obtaining certificate for $DOMAIN using standalone mode..."
    certbot certonly --standalone \
        --email admin@$DOMAIN \
        --agree-tos \
        --no-eff-email \
        --non-interactive \
        -d $DOMAIN \
        --cert-name $DOMAIN \
        --cert-path "$CERT_DIR/cert.pem" \
        --key-path "$CERT_DIR/privkey.pem" \
        --chain-path "$CERT_DIR/chain.pem" \
        --fullchain-path "$CERT_DIR/fullchain.pem"
    
    print_info "Starting Nginx..."
    docker-compose -f ../docker-compose.yml start nginx
}

# Obtain certificate using webroot mode (Nginx stays running)
obtain_cert_webroot() {
    DOMAIN=$1
    
    print_info "Obtaining certificate for $DOMAIN using webroot mode..."
    
    # Create webroot directory for ACME challenge
    WEBROOT_PATH="/var/www/certbot"
    docker-compose -f ../docker-compose.yml exec nginx mkdir -p $WEBROOT_PATH
    
    certbot certonly --webroot \
        --email admin@$DOMAIN \
        --agree-tos \
        --no-eff-email \
        --non-interactive \
        -d $DOMAIN \
        --cert-name $DOMAIN \
        --webroot-path $WEBROOT_PATH \
        --cert-path "$CERT_DIR/cert.pem" \
        --key-path "$CERT_DIR/privkey.pem" \
        --chain-path "$CERT_DIR/chain.pem" \
        --fullchain-path "$CERT_DIR/fullchain.pem"
}

# Update Nginx configuration
update_nginx_config() {
    DOMAIN=$1
    
    print_info "Updating Nginx configuration..."
    
    # Create backup of current config
    if [ -f "$NGINX_DIR/nginx.conf" ]; then
        cp "$NGINX_DIR/nginx.conf" "$NGINX_DIR/nginx.conf.backup.$(date +%Y%m%d_%H%M%S)"
    fi
    
    # Update server_name in nginx.conf if it exists
    if [ -f "$NGINX_DIR/nginx.conf" ]; then
        sed -i "s/server_name localhost;/server_name $DOMAIN;/g" "$NGINX_DIR/nginx.conf"
        print_info "Updated server_name in nginx.conf to $DOMAIN"
    fi
    
    print_info "Nginx configuration updated"
}

# Restart Nginx
restart_nginx() {
    print_info "Restarting Nginx to apply new certificates..."
    docker-compose -f ../docker-compose.yml restart nginx
    print_info "Nginx restarted successfully"
}

# Verify certificate
verify_certificate() {
    DOMAIN=$1
    
    print_info "Verifying certificate..."
    
    if [ -f "$CERT_DIR/fullchain.pem" ] && [ -f "$CERT_DIR/privkey.pem" ]; then
        print_info "Certificate files exist:"
        ls -lh "$CERT_DIR"/*.pem
        print_info "Certificate will expire on:"
        openssl x509 -in "$CERT_DIR/fullchain.pem" -noout -dates
    else
        print_error "Certificate files not found!"
        exit 1
    fi
}

# Setup auto-renewal
setup_auto_renewal() {
    DOMAIN=$1
    
    print_info "Setting up automatic renewal..."
    
    # Create renewal script
    RENEWAL_SCRIPT="../scripts/renew-cert.sh"
    cat > "$RENEWAL_SCRIPT" << 'EOF'
#!/bin/bash
set -e

CERT_DIR="../ssl"
NGINX_DIR="../nginx"

echo "Renewing Let's Encrypt certificate..."
certbot renew --quiet

echo "Restarting Nginx..."
docker-compose -f ../docker-compose.yml restart nginx

echo "Certificate renewed successfully at $(date)"
EOF
    
    chmod +x "$RENEWAL_SCRIPT"
    
    # Add to crontab (renew twice daily)
    if ! crontab -l | grep -q "renew-cert.sh"; then
        (crontab -l 2>/dev/null; echo "0 0,12 * * * cd /opt/whu-campus-auth && $RENEWAL_SCRIPT") | crontab -
        print_info "Auto-renewal cron job added (runs at midnight and noon)"
    else
        print_info "Auto-renewal cron job already exists"
    fi
    
    # Test renewal process
    print_info "Testing renewal process..."
    certbot renew --dry-run
}

# Main function
main() {
    DOMAIN=$1
    
    print_info "========================================="
    print_info "Let's Encrypt Certificate Setup"
    print_info "Domain: $DOMAIN"
    print_info "========================================="
    
    # Pre-flight checks
    check_root
    check_domain "$DOMAIN"
    check_certbot
    create_cert_dir
    backup_existing_certs
    
    # Obtain certificate (try webroot first, fallback to standalone)
    if ! obtain_cert_webroot "$DOMAIN"; then
        print_warning "Webroot mode failed, trying standalone mode..."
        obtain_cert_standalone "$DOMAIN"
    fi
    
    # Update configuration
    update_nginx_config "$DOMAIN"
    
    # Verify and restart
    verify_certificate "$DOMAIN"
    restart_nginx
    
    # Setup auto-renewal
    setup_auto_renewal "$DOMAIN"
    
    print_info "========================================="
    print_info "✅ Certificate setup completed!"
    print_info "Domain: $DOMAIN"
    print_info "Certificate directory: $CERT_DIR"
    print_info "Auto-renewal: Enabled"
    print_info "========================================="
    
    # Display next steps
    echo ""
    print_info "Next steps:"
    echo "1. Verify your site is accessible via HTTPS: https://$DOMAIN"
    echo "2. Check certificate expiration: openssl x509 -in $CERT_DIR/fullchain.pem -noout -dates"
    echo "3. Test auto-renewal: sudo certbot renew --dry-run"
    echo ""
}

# Run main function
main "$@"
