#!/bin/bash

print_info "Collecting CoreDNS configuration..."

prompt_if_unset_config "cloud.internalDomain" "Enter internal domain name" "local.example.com"
prompt_if_unset_config "cluster.loadBalancerIp" "Enter load balancer IP address" "192.168.1.240"
prompt_if_unset_config "cloud.dns.externalResolver" "Enter external DNS resolver" "8.8.8.8"
