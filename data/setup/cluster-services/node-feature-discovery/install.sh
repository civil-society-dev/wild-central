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
NFD_DIR="${CLUSTER_SETUP_DIR}/node-feature-discovery"

print_header "Setting up Node Feature Discovery"

# Templates should already be compiled by wild-cluster-services-generate
info "Using pre-compiled Node Feature Discovery templates..."
if [ ! -d "${NFD_DIR}/kustomize" ]; then
    error "Compiled templates not found. Run 'wild-cluster-services-configure node-feature-discovery' first."
    exit 1
fi

info "Deploying Node Feature Discovery..."
kubectl apply -k "${NFD_DIR}/kustomize"

info "Waiting for Node Feature Discovery DaemonSet to be ready..."
kubectl rollout status daemonset/node-feature-discovery-worker -n node-feature-discovery --timeout=300s

success "Node Feature Discovery installed successfully"

echo ""
echo "To verify the installation:"
echo "  kubectl get pods -n node-feature-discovery"
echo "  kubectl get nodes --show-labels | grep feature.node.kubernetes.io"
echo ""
echo "GPU nodes should now be labeled with GPU device information:"
echo "  kubectl get nodes --show-labels | grep pci-10de"