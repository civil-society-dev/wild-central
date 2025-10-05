#!/bin/bash

# Detect arm or amd
ARCH=$(uname -m)
if [ "$ARCH" != "aarch64" ] && [ "$ARCH" != "x86_64" ]; then
    echo "Error: Unsupported architecture $ARCH. Only arm64 and amd64 are supported."
    exit 1
fi

ARCH_ABBR="amd64"
if [ "$ARCH" == "aarch64" ]; then
    ARCH_ABBR="arm64"
fi

# Install kubectl
if ! command -v kubectl &> /dev/null; then
    echo "Error: kubectl is not installed. Installing."
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/$ARCH_ABBR/kubectl"
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/$ARCH_ABBR/kubectl.sha256"
    echo "$(cat kubectl.sha256)  kubectl" | sha256sum --check
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
fi

# Install talosctl
if command -v talosctl &> /dev/null; then
    echo "talosctl is already installed."
else
    curl -sL https://talos.dev/install | sh
    echo "talosctl installed successfully."
fi

# Install gomplate
if command -v gomplate &> /dev/null; then
    echo "gomplate is already installed."
else
    curl -sSL https://github.com/hairyhenderson/gomplate/releases/latest/download/gomplate_linux-$ARCH_ABBR -o $HOME/.local/bin/gomplate
    chmod +x $HOME/.local/bin/gomplate
    echo "gomplate installed successfully."
fi

# Install kustomize
if command -v kustomize &> /dev/null; then
    echo "kustomize is already installed."
else
    curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash
    mv kustomize $HOME/.local/bin/
    echo "kustomize installed successfully."
fi

## Install yq
if command -v yq &> /dev/null; then
    echo "yq is already installed."
else
    VERSION=v4.45.4
    BINARY=yq_linux_$ARCH_ABBR
    wget https://github.com/mikefarah/yq/releases/download/${VERSION}/${BINARY}.tar.gz -O - | tar xz
    mv ${BINARY} $HOME/.local/bin/yq
    chmod +x $HOME/.local/bin/yq
    rm yq.1
    echo "yq installed successfully."
fi

## Install restic
if command -v restic &> /dev/null; then
    echo "restic is already installed."
else
    sudo apt-get update
    sudo apt-get install -y restic
    echo "restic installed successfully."
fi

## Install direnv
if command -v direnv &> /dev/null; then
    echo "direnv is already installed."
else
    sudo apt-get update
    sudo apt-get install -y direnv
    echo "direnv installed successfully. Add `eval \"\$(direnv hook bash)\"` to your shell configuration file if not already present."
fi
