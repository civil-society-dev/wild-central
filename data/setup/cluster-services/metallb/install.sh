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
METALLB_DIR="${CLUSTER_SETUP_DIR}/metallb"

print_header "Setting up MetalLB"

# Templates should already be compiled by wild-cluster-services-generate
echo "Using pre-compiled MetalLB templates..."
if [ ! -d "${METALLB_DIR}/kustomize" ]; then
    echo "Error: Compiled templates not found. Run 'wild-cluster-services-generate' first."
    exit 1
fi

echo "Deploying MetalLB..."
kubectl apply -k ${METALLB_DIR}/kustomize/installation

echo "Waiting for MetalLB to be deployed..."
kubectl wait --for=condition=Available deployment/controller -n metallb-system --timeout=60s
sleep 10 # Extra buffer for webhook initialization

echo "Customizing MetalLB..."
kubectl apply -k ${METALLB_DIR}/kustomize/configuration

echo "✅ MetalLB installed and configured"
echo ""
echo "To verify the installation:"
echo "  kubectl get pods -n metallb-system"
echo "  kubectl get ipaddresspools.metallb.io -n metallb-system"
