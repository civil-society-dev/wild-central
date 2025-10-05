#!/bin/bash

print_info "Collecting Docker Registry configuration..."

prompt_if_unset_config "cloud.dockerRegistryHost" "Enter Docker Registry hostname" "registry.local.example.com"
prompt_if_unset_config "cluster.dockerRegistry.storage" "Enter Docker Registry storage size" "100Gi"
