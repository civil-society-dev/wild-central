#!/bin/bash

print_info "Collecting NFS configuration..."

prompt_if_unset_config "cloud.nfs.host" "Enter NFS server hostname or IP address" "192.168.1.100"
prompt_if_unset_config "cloud.nfs.mediaPath" "Enter NFS export path for media storage" "/mnt/storage/media"
prompt_if_unset_config "cloud.nfs.storageCapacity" "Enter NFS storage capacity (e.g., 1Ti, 500Gi)" "1Ti"
