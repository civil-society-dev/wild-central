print_info "Collecting ExternalDNS configuration..."

prompt_if_unset_config "cluster.externalDns.ownerId" "Enter ExternalDNS owner ID (unique identifier for this cluster)" "wild-cloud-$(hostname -s)"
