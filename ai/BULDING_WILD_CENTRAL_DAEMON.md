# Building the Wild Cloud Central Daemon

## Philosophy

- Like v.PoC, we should only use gomplate templates for distinguishing between cloud instances. However, **within** a cloud instance, there should be no templating. The templates are compiled when being copied into the instances. This allows transparency and simple management by the user.
- Manage state and infrastructure idempotently.
- Cluster state should be the k8s cluster itself, not local files. It should be accesses via kubectl and talosctl.
- All wild cloud state should be stored on the filesystem in easy to read YAML files, and can be edited directly or through the webapp.

- Use talosctl and kubectl wherever possible to leverage existing tools and avoid reinventing the wheel.
- All code should be simple and easy to understand.
  - Avoid unnecessary complexity.
  - Avoid unnecessary dependencies.
  - Avoid unnecessary features.
  - Avoid unnecessary abstractions.
  - Avoid unnecessary comments.
  - Avoid unnecessary configuration options.
- Avoid Helm. Use Kustomize.
- The daemon should be able to run on low-resource devices like a Raspberry Pi 4 (4GB RAM).
- The daemon should be able to manage multiple Wild Cloud instances on the LAN.
- The daemon should include functionality to manage a dnsmasq server on the same device. Currently, this is only used to resolve wild cloud domain names within the LAN to provide for private addresses on the LAN. The LAN router should be configured to use the Wild Central IP as its DNS server.
- The Daemon is configurable to use various providers for:
  - Wild Cloud Apps Directory provider (local FS, git repo, etc)
  - DNS (built-in dnsmasq, external DNS server, etc)

### Coding

- Use a standard Go project structure.
- Use Go modules.
- Use standard Go libraries wherever possible.
- Use popular, well-maintained libraries for common tasks (e.g. gorilla/mux for HTTP routing).
- Write unit tests for all functions and methods.
- Make and use common modules. For example, one module should handle all interactions with talosctl. Another modules should handle all interactions with kubectl. 
- If the code is getting long and complex, break it into smaller modules.
