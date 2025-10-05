#!/bin/bash

print_info "Collecting Kubernetes Dashboard configuration..."

prompt_if_unset_config "cloud.internalDomain" "Enter internal domain name (for dashboard URL)" "local.example.com"
