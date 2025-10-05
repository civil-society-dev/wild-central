#!/bin/bash
set -e
set -o pipefail

# Initialize Wild Cloud environment
if [ -z "${WC_ROOT}" ]; then
    print "WC_ROOT is not set."
    exit 1
else
    source "${WC_ROOT}/scripts/common.sh"
    init_wild_env
fi

CLUSTER_SETUP_DIR="${WC_HOME}/setup/cluster-services"
CERT_MANAGER_DIR="${CLUSTER_SETUP_DIR}/cert-manager"

print_header "Setting up cert-manager"

#######################
# # Dependencies
#######################

# Check Traefik dependency
print_info "Verifying Traefik is ready (required for cert-manager)..."
kubectl wait --for=condition=Available deployment/traefik -n traefik --timeout=60s 2>/dev/null || {
    print_warning "Traefik not ready, but continuing with cert-manager installation"
    print_info "Note: cert-manager may not work properly without Traefik"
}

if [ ! -d "${CERT_MANAGER_DIR}/kustomize" ]; then
    print_error "Compiled templates not found. This script should not be run directly. Run with 'wild setup cluster-services cert-manager' instead."
    exit 1
fi

# Validate DNS resolution using temporary test pod
print_info "Validating DNS resolution for ACME challenges..."
domain=$(wild-config cluster.certManager.cloudflare.domain)
print_info "Testing DNS resolution for domain: $domain"

# Create temporary pod with DNS utilities (in default namespace since cert-manager doesn't exist yet)
kubectl run dns-test --image=busybox:1.35 --rm -i --restart=Never -- \
    nslookup -type=SOA "$domain" 1.1.1.1 &>/dev/null && \
    print_success "DNS resolution working for ACME challenges" || \
    print_warning "DNS resolution issues may affect ACME challenges"


########################
# Cloudflare DNS setup
########################

# API token secret setup
print_info "Reading Cloudflare API token secret..."
CLOUDFLARE_API_TOKEN=$(wild-secret cloudflare.token) || exit 1
if [ -z "$CLOUDFLARE_API_TOKEN" ]; then
    print_error "Cloudflare API token not found. Please create it with 'wild secret create cloudflare.token'."
    exit 1
fi

# Validate token
print_info "Validating Cloudflare API token permissions..."
validate_cloudflare_token() {
    local token="$1"
    if ! command -v curl &>/dev/null; then
        print_warning "curl not available, skipping token validation"
        return 0
    fi

    print_info "Testing Cloudflare API token..."
    local response
    response=$(curl -s -H "Authorization: Bearer $token" \
                   "https://api.cloudflare.com/client/v4/zones")

    if echo "$response" | grep -q '"success":true'; then
        print_success "Cloudflare API token is valid and has zone access"
        return 0
    else
        print_error "Cloudflare token validation failed"
        print_info "Response: $response"
        print_info "Please ensure your token has Zone - Zone - Read permission"
        return 1
    fi
}

validate_cloudflare_token "$CLOUDFLARE_API_TOKEN" || {
    print_error "Cloudflare token validation failed. Please check token permissions."
    print_info "Required permissions: Zone - Zone - Read, Zone - DNS - Edit"
    exit 1
}

########################
# Kubernetes components
########################

print_info "Installing cert-manager components..."
# Using stable URL for cert-manager installation
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.1/cert-manager.yaml || \
  kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.13.1/cert-manager.yaml

# Wait for cert-manager to be ready
print_info "Waiting for cert-manager to be ready..."
kubectl wait --for=condition=Available deployment/cert-manager -n cert-manager --timeout=120s
kubectl wait --for=condition=Available deployment/cert-manager-cainjector -n cert-manager --timeout=120s
kubectl wait --for=condition=Available deployment/cert-manager-webhook -n cert-manager --timeout=120s

# Now that cert-manager namespace exists, create the Cloudflare API token secret
print_info "Creating Cloudflare API token secret..."
kubectl create secret generic cloudflare-api-token \
  --namespace cert-manager \
  --from-literal=api-token="${CLOUDFLARE_API_TOKEN}" \
  --dry-run=client -o yaml | kubectl apply -f -

# Ensure webhook is fully operational
print_info "Verifying cert-manager webhook is fully operational..."
until kubectl get validatingwebhookconfigurations cert-manager-webhook &>/dev/null; do
    print_info "Waiting for cert-manager webhook to register..."
    sleep 5
done

# Test webhook connectivity before proceeding
print_info "Testing webhook connectivity..."
kubectl auth can-i create certificates.cert-manager.io --as=system:serviceaccount:cert-manager:cert-manager

# Configure cert-manager to use external DNS for challenge verification
print_info "Configuring cert-manager to use external DNS servers..."
kubectl patch deployment cert-manager -n cert-manager --patch '
spec:
  template:
    spec:
      dnsPolicy: None
      dnsConfig:
        nameservers:
          - "1.1.1.1"
          - "8.8.8.8"
        searches:
          - cert-manager.svc.cluster.local
          - svc.cluster.local
          - cluster.local
        options:
          - name: ndots
            value: "5"'

# Wait for cert-manager to restart with new DNS config
print_info "Waiting for cert-manager to restart with new DNS configuration..."
kubectl rollout status deployment/cert-manager -n cert-manager --timeout=120s

########################
# Create issuers and certificates
########################

# Apply Let's Encrypt issuers and certificates using kustomize
print_info "Creating Let's Encrypt issuers and certificates..."
kubectl apply -k ${CERT_MANAGER_DIR}/kustomize

# Wait for issuers to be ready
print_info "Waiting for Let's Encrypt issuers to be ready..."
kubectl wait --for=condition=Ready clusterissuer/letsencrypt-prod --timeout=60s || print_warning "Production issuer not ready, proceeding anyway..."
kubectl wait --for=condition=Ready clusterissuer/letsencrypt-staging --timeout=60s || print_warning "Staging issuer not ready, proceeding anyway..."

# Give cert-manager a moment to process the certificates
sleep 5

######################################
# Fix stuck certificates and cleanup
######################################

needs_restart=false

# STEP 1: Fix certificates stuck with 404 errors FIRST (before cleaning up orders)
print_info "Checking for certificates with failed issuance attempts..."
stuck_certs=$(kubectl get certificates --all-namespaces -o json 2>/dev/null | \
    jq -r '.items[] | select(.status.conditions[]? | select(.type=="Issuing" and .status=="False" and (.message | contains("404")))) | "\(.metadata.namespace) \(.metadata.name)"')

if [ -n "$stuck_certs" ]; then
    print_warning "Found certificates stuck with non-existent orders, recreating them..."
    echo "$stuck_certs" | while read ns name; do
        print_info "Recreating certificate $ns/$name..."
        # Get just the spec
        cert_spec=$(kubectl get certificate "$name" -n "$ns" -o json | jq '.spec')
        # Delete the certificate
        kubectl delete certificate "$name" -n "$ns"
        # Recreate with clean state
        echo "{\"apiVersion\":\"cert-manager.io/v1\",\"kind\":\"Certificate\",\"metadata\":{\"name\":\"$name\",\"namespace\":\"$ns\"},\"spec\":$cert_spec}" | kubectl apply -f -
    done
    needs_restart=true
    # Give cert-manager time to process the recreated certificates
    sleep 5
else
    print_success "No certificates stuck with failed orders"
fi

# STEP 2: Clean up orphaned orders (after fixing certificates)
print_info "Checking for orphaned ACME orders..."

# Check logs for 404 errors
orphaned_orders=$(kubectl logs -n cert-manager deployment/cert-manager --tail=200 2>/dev/null | \
    grep -E "failed to retrieve the ACME order.*404" 2>/dev/null | \
    sed -n 's/.*resource_name="\([^"]*\)".*/\1/p' | \
    sort -u || true)

if [ -n "$orphaned_orders" ]; then
    print_warning "Found orphaned ACME orders from logs"
    for order in $orphaned_orders; do
        print_info "Deleting orphaned order: $order"
        # Find and delete the order in whatever namespace it exists
        orders_found=$(kubectl get orders --all-namespaces 2>/dev/null | grep "$order" 2>/dev/null || true)
        if [ -n "$orders_found" ]; then
            echo "$orders_found" | while read ns name rest; do
                kubectl delete order "$name" -n "$ns" 2>/dev/null || true
            done
        fi
    done
    needs_restart=true
else
    print_success "No orphaned orders found in logs"
fi

# Check for errored state orders
errored_orders=$(kubectl get orders --all-namespaces -o json 2>/dev/null | \
    jq -r '.items[] | select(.status.state == "errored") | "\(.metadata.namespace) \(.metadata.name)"')

if [ -n "$errored_orders" ]; then
    print_warning "Found errored ACME orders"
    echo "$errored_orders" | while read ns name; do
        print_info "Deleting errored order: $ns/$name"
        kubectl delete order "$name" -n "$ns" 2>/dev/null || true
    done
    needs_restart=true
else
    print_success "No errored orders found"
fi

# STEP 3: Clean up bad challenges
print_info "Checking for stuck ACME challenges..."

# Delete expired, invalid, or errored challenges
bad_challenges=$(kubectl get challenges --all-namespaces -o json 2>/dev/null | \
    jq -r '.items[] | select(.status.state == "expired" or .status.state == "invalid" or .status.state == "errored") | "\(.metadata.namespace) \(.metadata.name) \(.status.state)"')

if [ -n "$bad_challenges" ]; then
    print_warning "Found stuck ACME challenges"
    echo "$bad_challenges" | while read ns name state; do
        print_info "Deleting $state challenge: $ns/$name"
        kubectl delete challenge "$name" -n "$ns" 2>/dev/null || true
    done
    needs_restart=true
else
    print_success "No stuck challenges found"
fi

# Delete very old challenges (over 1 hour) - only if they exist
all_challenges=$(kubectl get challenges --all-namespaces -o json 2>/dev/null | jq '.items | length' || echo 0)
if [ "$all_challenges" -gt 0 ]; then
    old_challenges=$(kubectl get challenges --all-namespaces -o json 2>/dev/null | \
        jq -r --arg cutoff "$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || date -u -v-1H '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null)" \
        '.items[] | select(.metadata.creationTimestamp < $cutoff) | "\(.metadata.namespace) \(.metadata.name)"')

    if [ -n "$old_challenges" ]; then
        print_warning "Found old challenges (over 1 hour)"
        echo "$old_challenges" | while read ns name; do
            print_info "Deleting old challenge: $ns/$name"
            kubectl delete challenge "$name" -n "$ns" 2>/dev/null || true
        done
        needs_restart=true
    fi
fi

# STEP 4: Check for DNS errors
dns_errors=$(kubectl logs -n cert-manager deployment/cert-manager --tail=50 2>/dev/null | \
    grep "Could not route to /client/v4/zones/dns_records" | wc -l | tr -d '\n' || echo "0")
dns_errors=${dns_errors:-0}

if [ "$dns_errors" -gt 0 ]; then
    print_warning "Cert-manager has DNS record cleanup errors"
    needs_restart=true
fi

# STEP 5: Single restart if anything needs cleaning
if [ "$needs_restart" = true ]; then
    print_info "Restarting cert-manager once to clear all internal state..."
    kubectl rollout restart deployment cert-manager -n cert-manager
    kubectl rollout status deployment/cert-manager -n cert-manager --timeout=120s
    # Give cert-manager time to reinitialize
    sleep 10
else
    print_success "No restart needed - cert-manager state is clean"
fi


##################################
# Handle certificate renewal
##################################

# Check for expired or near-expiry certificates and trigger renewal
print_info "Checking certificate expiration status..."
current_date=$(date +%s)

# Track if we found any issues
found_expired=false
found_expiring_soon=false
all_certs_valid=true

# Process certificates and collect their status
while IFS= read -r line; do
    ns=$(echo "$line" | awk '{print $1}')
    name=$(echo "$line" | awk '{print $2}')
    secret=$(echo "$line" | awk '{print $3}')
    expiry=$(echo "$line" | awk '{print $4}')

    if [ "$expiry" != "unknown" ] && [ "$expiry" != "null" ] && [ "$expiry" != "" ]; then
        expiry_ts=$(date -d "$expiry" +%s 2>/dev/null || date -j -f "%Y-%m-%dT%H:%M:%SZ" "$expiry" +%s 2>/dev/null || echo 0)
        if [ "$expiry_ts" -gt 0 ]; then
            days_until_expiry=$(( (expiry_ts - current_date) / 86400 ))

            if [ "$days_until_expiry" -lt 0 ]; then
                print_warning "Certificate $ns/$name has EXPIRED (expired ${days_until_expiry#-} days ago)"
                if [ -n "$secret" ] && [ "$secret" != "unknown" ] && [ "$secret" != "null" ]; then
                    print_info "Deleting secret $secret to trigger renewal..."
                    kubectl delete secret "$secret" -n "$ns" 2>/dev/null || true
                    found_expired=true
                    all_certs_valid=false
                fi
            elif [ "$days_until_expiry" -lt 7 ]; then
                print_warning "Certificate $ns/$name expires in $days_until_expiry days"
                if [ "$days_until_expiry" -lt 3 ]; then
                    # Force renewal for certificates expiring very soon
                    if [ -n "$secret" ] && [ "$secret" != "unknown" ] && [ "$secret" != "null" ]; then
                        print_info "Forcing renewal by deleting secret $secret..."
                        kubectl delete secret "$secret" -n "$ns" 2>/dev/null || true
                        found_expiring_soon=true
                        all_certs_valid=false
                    fi
                else
                    print_info "Will renew automatically when closer to expiry"
                fi
            elif [ "$days_until_expiry" -lt 30 ]; then
                print_info "Certificate $ns/$name expires in $days_until_expiry days (renewal not needed yet)"
            else
                print_success "Certificate $ns/$name is valid for $days_until_expiry days"
            fi
        fi
    else
        # Certificate has no expiry (being issued)
        print_info "Certificate $ns/$name is currently being issued..."
    fi
done < <(kubectl get certificates --all-namespaces -o json 2>/dev/null | jq -r '.items[] | "\(.metadata.namespace) \(.metadata.name) \(.spec.secretName) \(.status.notAfter // "unknown")"')

if [ "$all_certs_valid" = true ]; then
    print_success "All certificates are valid - no renewals needed"
fi


#########################
# Final checks
#########################

# Wait for the certificates to be issued (with a timeout)
print_info "Waiting for wildcard certificates to be ready (this may take several minutes)..."
kubectl wait --for=condition=Ready certificate wildcard-internal-wild-cloud -n cert-manager --timeout=300s || true
kubectl wait --for=condition=Ready certificate wildcard-wild-cloud -n cert-manager --timeout=300s || true

# Final health check
print_info "Performing final cert-manager health check..."
failed_certs=$(kubectl get certificates --all-namespaces -o json 2>/dev/null | jq -r '.items[] | select(.status.conditions[]? | select(.type=="Ready" and .status!="True")) | "\(.metadata.namespace)/\(.metadata.name)"' | wc -l)
if [ "$failed_certs" -gt 0 ]; then
    print_warning "Found $failed_certs certificates not in Ready state"
    print_info "Check certificate status with: kubectl get certificates --all-namespaces"
    print_info "Check cert-manager logs with: kubectl logs -n cert-manager deployment/cert-manager"
else
    print_success "All certificates are in Ready state"
fi

print_success "cert-manager setup complete!"
