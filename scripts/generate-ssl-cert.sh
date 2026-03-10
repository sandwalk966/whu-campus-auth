#!/bin/bash

# Self-Signed SSL Certificate Generation Script
# This script generates self-signed certificates for development/testing

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CERT_DIR="../ssl"
CERT_NAME="server"
DAYS_VALID=365
KEY_SIZE=2048

# Default values
COUNTRY="CN"
STATE="Hubei"
LOCALITY="Wuhan"
ORGANIZATION="WHU Campus Auth"
ORGANIZATIONAL_UNIT="Development"
COMMON_NAME="localhost"

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

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Check if openssl is installed
check_openssl() {
    if ! command -v openssl &> /dev/null; then
        print_error "OpenSSL is not installed. Please install it first."
        echo "  Ubuntu/Debian: sudo apt-get install openssl"
        echo "  CentOS/RHEL: sudo yum install openssl"
        echo "  macOS: brew install openssl"
        exit 1
    fi
    print_info "OpenSSL version: $(openssl version)"
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
    if [ -f "$CERT_DIR/${CERT_NAME}.crt" ] || [ -f "$CERT_DIR/${CERT_NAME}.key" ]; then
        print_warning "Existing certificates found"
        BACKUP_DIR="$CERT_DIR/backup-$(date +%Y%m%d_%H%M%S)"
        mkdir -p "$BACKUP_DIR"
        
        if [ -f "$CERT_DIR/${CERT_NAME}.crt" ]; then
            cp "$CERT_DIR/${CERT_NAME}.crt" "$BACKUP_DIR/"
        fi
        if [ -f "$CERT_DIR/${CERT_NAME}.key" ]; then
            cp "$CERT_DIR/${CERT_NAME}.key" "$BACKUP_DIR/"
        fi
        if [ -f "$CERT_DIR/fullchain.pem" ]; then
            cp "$CERT_DIR/fullchain.pem" "$BACKUP_DIR/"
        fi
        if [ -f "$CERT_DIR/privkey.pem" ]; then
            cp "$CERT_DIR/privkey.pem" "$BACKUP_DIR/"
        fi
        
        print_info "Existing certificates backed up to: $BACKUP_DIR"
    fi
}

# Generate self-signed certificate
generate_certificate() {
    print_step "Generating self-signed certificate..."
    echo ""
    
    # Generate private key and certificate in one command
    openssl req -x509 -nodes -days $DAYS_VALID -newkey rsa:$KEY_SIZE \
        -keyout "$CERT_DIR/${CERT_NAME}.key" \
        -out "$CERT_DIR/${CERT_NAME}.crt" \
        -subj "/C=$COUNTRY/ST=$STATE/L=$LOCALITY/O=$ORGANIZATION/OU=$ORGANIZATIONAL_UNIT/CN=$COMMON_NAME" \
        -addext "subjectAltName=DNS:localhost,IP:127.0.0.1,IP:::1"
    
    print_info "Certificate generated successfully!"
}

# Create combined files for Nginx
create_nginx_files() {
    print_step "Creating Nginx-compatible certificate files..."
    
    # Copy to fullchain.pem and privkey.pem (Let's Encrypt naming convention)
    cp "$CERT_DIR/${CERT_NAME}.crt" "$CERT_DIR/fullchain.pem"
    cp "$CERT_DIR/${CERT_NAME}.key" "$CERT_DIR/privkey.pem"
    
    print_info "Created files:"
    echo "  - $CERT_DIR/fullchain.pem (certificate)"
    echo "  - $CERT_DIR/privkey.pem (private key)"
}

# Set proper permissions
set_permissions() {
    print_step "Setting file permissions..."
    
    # Private key should be readable only by owner
    chmod 600 "$CERT_DIR/${CERT_NAME}.key"
    chmod 600 "$CERT_DIR/privkey.pem"
    
    # Certificate can be world-readable
    chmod 644 "$CERT_DIR/${CERT_NAME}.crt"
    chmod 644 "$CERT_DIR/fullchain.pem"
    
    print_info "Permissions set:"
    echo "  - Private keys: 600 (owner read/write only)"
    echo "  - Certificates: 644 (world readable)"
}

# Verify certificate
verify_certificate() {
    print_step "Verifying certificate..."
    echo ""
    
    print_info "Certificate details:"
    openssl x509 -in "$CERT_DIR/${CERT_NAME}.crt" -noout -subject -issuer -dates
    
    echo ""
    print_info "Certificate fingerprint (SHA256):"
    openssl x509 -in "$CERT_DIR/${CERT_NAME}.crt" -noout -fingerprint -sha256
    
    echo ""
    print_info "Subject Alternative Names:"
    openssl x509 -in "$CERT_DIR/${CERT_NAME}.crt" -noout -text | grep -A1 "Subject Alternative Name"
}

# Display usage instructions
display_usage() {
    echo ""
    print_info "========================================="
    print_info "✅ Self-signed certificate generation completed!"
    print_info "========================================="
    echo ""
    print_info "Generated files:"
    echo "  $CERT_DIR/"
    echo "  ├── ${CERT_NAME}.crt        (certificate)"
    echo "  ├── ${CERT_NAME}.key        (private key)"
    echo "  ├── fullchain.pem           (certificate for Nginx)"
    echo "  └── privkey.pem             (private key for Nginx)"
    echo ""
    print_warning "⚠️  IMPORTANT: This is a SELF-SIGNED certificate!"
    echo ""
    echo "  - ✅ Suitable for: Development and Testing"
    echo "  - ❌ NOT suitable for: Production"
    echo "  - ⚠️  Browsers will show security warnings"
    echo ""
    print_info "Next steps:"
    echo "  1. Restart Nginx: docker-compose restart nginx"
    echo "  2. Access via HTTPS: https://localhost"
    echo "  3. Accept the browser security warning"
    echo "  4. For production, use Let's Encrypt instead:"
    echo "     sudo ./scripts/letsencrypt.sh your-domain.com"
    echo ""
    print_info "To trust this certificate system-wide (optional):"
    echo "  - Ubuntu/Debian:"
    echo "    sudo cp $CERT_DIR/${CERT_NAME}.crt /usr/local/share/ca-certificates/"
    echo "    sudo update-ca-certificates"
    echo ""
    echo "  - CentOS/RHEL:"
    echo "    sudo cp $CERT_DIR/${CERT_NAME}.crt /etc/pki/ca-trust/source/anchors/"
    echo "    sudo update-ca-trust"
    echo ""
    echo "  - macOS:"
    echo "    sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain $CERT_DIR/${CERT_NAME}.crt"
    echo ""
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--days)
                DAYS_VALID="$2"
                shift 2
                ;;
            -n|--name)
                CERT_NAME="$2"
                shift 2
                ;;
            -c|--cn)
                COMMON_NAME="$2"
                shift 2
                ;;
            -h|--help)
                echo "Usage: $0 [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  -d, --days NUM      Number of days certificate is valid (default: 365)"
                echo "  -n, --name NAME     Certificate name (default: server)"
                echo "  -c, --cn NAME       Common Name for certificate (default: localhost)"
                echo "  -h, --help          Show this help message"
                echo ""
                echo "Example:"
                echo "  $0 -d 730 -n mycert -c mydomain.com"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                echo "Use -h or --help for usage information"
                exit 1
                ;;
        esac
    done
}

# Main function
main() {
    print_info "========================================="
    print_info "Self-Signed SSL Certificate Generator"
    print_info "========================================="
    echo ""
    
    # Parse arguments
    parse_args "$@"
    
    # Pre-flight checks
    check_openssl
    create_cert_dir
    backup_existing_certs
    
    # Generate certificate
    generate_certificate
    
    # Create Nginx files
    create_nginx_files
    
    # Set permissions
    set_permissions
    
    # Verify certificate
    verify_certificate
    
    # Display usage
    display_usage
}

# Run main function
main "$@"
