# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Wild-cloud Central is a bash-based infrastructure setup system for creating inexpensive Linux machines that serve as the central component of a "wild-cloud". The system currently focuses on setting up a dnsmasq server that provides DNS and PXE services for clustering, with plans to evolve into a web service with JSON API.

## Architecture

The project consists of two main phases:

1. **Local preparation phase** (`install-dnsmasq`): Run on operator's computer
   - Reads configuration from `.wildcloud/config.yaml` and `.wildcloud/secrets.yaml`
   - Uses gomplate templating to generate customized setup files
   - Downloads and caches Talos Linux PXE boot assets from factory.talos.dev
   - Copies generated files to target server via SCP

2. **Remote setup phase** (`ref/dnsmasq/setup.sh`): Run on the wild-central server
   - Installs dnsmasq and nginx packages
   - Configures dnsmasq for DHCP, DNS, and PXE boot services
   - Sets up nginx to serve iPXE boot files and Talos Linux images
   - Handles network configuration conflicts (systemd-resolved)

## Key Components

- **Configuration templating**: Uses gomplate with YAML config files for dynamic configuration
- **PXE boot chain**: iPXE → nginx-served boot script → Talos Linux kernel/initramfs
- **Network services**: dnsmasq provides DHCP, DNS forwarding, and TFTP for PXE
- **Asset management**: Downloads and caches Talos factory images based on configuration

## Development Commands

This is a bash script-based project with no build system. Key operations:

- Run local setup: `./ref/install-dnsmasq` (requires `.wildcloud/config.yaml`)
- Run remote setup: `./ref/dnsmasq/setup.sh` (run on target server)

## Dependencies

- `gomplate` - for template processing
- `yq` - for YAML parsing
- `jq` - for JSON processing
- `wget`, `curl` - for downloading assets
- `ssh`, `scp` - for remote operations

## Future Architecture

The project plans to transition from bash scripts to a web service installable via `sudo apt install wild-cloud-central` with a JSON API, indicating a move toward more structured service architecture.