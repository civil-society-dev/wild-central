#!/bin/bash

print_info "Collecting Traefik configuration..."

prompt_if_unset_config "cluster.loadBalancerIp" "Enter load balancer IP address for Traefik" "192.168.1.240"
