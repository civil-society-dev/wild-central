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
DOCKER_REGISTRY_DIR="${CLUSTER_SETUP_DIR}/docker-registry"

print_header "Setting up Docker Registry"

# Templates should already be compiled by wild-cluster-services-generate
echo "Using pre-compiled Docker Registry templates..."
if [ ! -d "${DOCKER_REGISTRY_DIR}/kustomize" ]; then
    echo "Error: Compiled templates not found. Run 'wild-cluster-services-generate' first."
    exit 1
fi

# Apply the docker registry manifests using kustomize
kubectl apply -k "${DOCKER_REGISTRY_DIR}/kustomize"

echo "Waiting for Docker Registry to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/docker-registry -n docker-registry

echo "Docker Registry setup complete!"

# Show deployment status
kubectl get pods -n docker-registry
kubectl get services -n docker-registry