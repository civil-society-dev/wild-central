#!/bin/bash

print_info "Collecting cert-manager configuration..."

prompt_if_unset_config "cloud.domain" "Enter main domain name" "example.com"
domain=$(wild-config "cloud.domain")
baseDomain=$(wild-config "cloud.baseDomain")
prompt_if_unset_config "cloud.internalDomain" "Enter internal domain name" "local.${domain}"
prompt_if_unset_config "operator.email" "Enter operator email address (for Let's Encrypt)" ""
prompt_if_unset_config "cluster.certManager.cloudflare.domain" "Enter Cloudflare domain" "${baseDomain}"
prompt_if_unset_config "cluster.certManager.cloudflare.zoneID" "Enter Cloudflare zone ID" ""
prompt_if_unset_secret "cloudflare.token" "Enter Cloudflare API token" ""
