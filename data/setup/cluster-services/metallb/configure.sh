#!/bin/bash

print_info "Collecting MetalLB configuration..."

prompt_if_unset_config "cluster.ipAddressPool" "Enter IP address pool for MetalLB (CIDR format, e.g., 192.168.1.240-192.168.1.250)" "192.168.1.240-192.168.1.250"
prompt_if_unset_config "cluster.loadBalancerIp" "Enter load balancer IP address" "192.168.1.240"
