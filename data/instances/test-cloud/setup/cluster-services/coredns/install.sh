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
COREDNS_DIR="${CLUSTER_SETUP_DIR}/coredns"

print_header "Setting up CoreDNS"

# Templates should already be compiled by wild-cluster-services-generate
echo "Using pre-compiled CoreDNS templates..."
if [ ! -d "${COREDNS_DIR}/kustomize" ]; then
    echo "Error: Compiled templates not found. Run 'wild-cluster-services-generate' first."
    exit 1
fi

# Apply the custom DNS override
# TODO: Is this needed now that we are no longer on k3s?
echo "Applying CoreDNS custom override configuration..."
kubectl apply -f "${COREDNS_DIR}/kustomize/coredns-custom-config.yaml"

# Restart CoreDNS pods to apply the changes
echo "Restarting CoreDNS pods to apply changes..."
kubectl rollout restart deployment/coredns -n kube-system
kubectl rollout status deployment/coredns -n kube-system

echo "CoreDNS setup complete!"
echo
echo "To verify the installation:"
echo "  kubectl get pods -n kube-system"
echo "  kubectl get svc -n kube-system coredns"
echo "  kubectl describe svc -n kube-system coredns"
echo "  kubectl logs -n kube-system -l k8s-app=kube-dns -f"
