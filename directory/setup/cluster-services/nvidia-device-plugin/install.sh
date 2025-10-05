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
NVIDIA_PLUGIN_DIR="${CLUSTER_SETUP_DIR}/nvidia-device-plugin"

print_header "Setting up NVIDIA Device Plugin"

# Check if we have NVIDIA GPUs in the cluster
print_info "Checking for NVIDIA GPUs in the cluster..."

# Check if any worker nodes exist (device plugin only runs on worker nodes)
WORKER_NODES=$(kubectl get nodes --selector='!node-role.kubernetes.io/control-plane' -o name | wc -l)
if [ "$WORKER_NODES" -eq 0 ]; then
    print_error "No worker nodes found in cluster. NVIDIA Device Plugin requires worker nodes."
    exit 1
fi

print_info "Found $WORKER_NODES worker node(s)"

# Templates should already be compiled by wild-cluster-services-generate
echo "Using pre-compiled NVIDIA Device Plugin templates..."
if [ ! -d "${NVIDIA_PLUGIN_DIR}/kustomize" ]; then
    echo "Error: Compiled templates not found. Run 'wild-cluster-services-generate' first."
    exit 1
fi

print_info "Deploying NVIDIA Device Plugin..."
kubectl apply -k ${NVIDIA_PLUGIN_DIR}/kustomize

print_info "Waiting for NVIDIA Device Plugin DaemonSet to be ready..."
kubectl rollout status daemonset/nvidia-device-plugin-daemonset -n kube-system --timeout=120s

print_success "NVIDIA Device Plugin installed successfully"
echo ""
echo "To verify the installation:"
echo "  kubectl get pods -n kube-system | grep nvidia"
echo "  kubectl get nodes -o json | jq '.items[].status.capacity | select(has(\"nvidia.com/gpu\"))'"
echo ""
echo "GPU nodes should now be labeled with GPU product information:"
echo "  kubectl get nodes --show-labels | grep nvidia"