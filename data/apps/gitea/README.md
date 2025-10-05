# Gitea Configuration

This Gitea deployment uses a hybrid configuration approach combining environment variables with Gitea's self-managed configuration file.

## Configuration Architecture

### Environment Variables (gitea.env)
Non-secret configuration is stored in `gitea.env` and automatically loaded via kustomize's `configMapGenerator`. This includes:
- Server settings (domain, URLs, ports)
- Database connection details (except password)
- SMTP settings (except password)
- Service settings (registration, notifications)
- Repository and storage paths

### Kubernetes Secrets (gitea-secrets)
Sensitive configuration is stored in the `gitea-secrets` secret and managed by the wild-cloud deployment system:
- `adminPassword` - Gitea admin user password
- `secretKey` - Application secret key
- `jwtSecret` - JWT signing secret
- `dbPassword` - Database password
- `smtpPassword` - SMTP authentication password

Secrets are defined in `secrets.yaml` and listed in `manifest.yaml` under `requiredSecrets`. The `wild-app-deploy` command automatically ensures all required secrets exist in the `gitea-secrets` secret before deployment.

### Persistent Configuration (app.ini)
Gitea manages its own `app.ini` file on persistent storage for:
- Generated security tokens
- Runtime configuration changes made via web UI
- Database migration state
- User-modified settings

## How It Works

1. **Startup**: Kustomize generates a ConfigMap from `gitea.env`
2. **Environment Loading**: Pod loads non-secret config from ConfigMap via `envFrom`
3. **Secret Loading**: Pod loads sensitive config from Kubernetes secrets via `env`
4. **Configuration Merge**: Gitea's environment-to-ini process merges environment variables into `app.ini`
5. **Persistence**: Gitea writes the merged configuration plus generated tokens to persistent storage

## Making Configuration Changes

### Non-Secret Settings
1. Edit `gitea.env` with your changes
2. Run `wild-app-deploy gitea` to apply changes
3. Pod will restart and pick up new configuration

### Secret Settings
1. Edit `secrets.yaml` with your secret values
2. Ensure the secret key is listed in `manifest.yaml` under `requiredSecrets`
3. Run `wild-app-deploy gitea` - this will automatically update the `gitea-secrets` secret and restart the pod

### Web UI Changes
Configuration changes made through Gitea's admin web interface are automatically persisted to the `app.ini` file on persistent storage and will survive pod restarts.

## Configuration Precedence

1. **Kubernetes Secrets** (highest priority)
2. **Environment Variables** (from gitea.env)
3. **Persistent app.ini** (lowest priority)

Environment variables override file settings, and secrets override everything.

## Troubleshooting

### Check Current Configuration
```bash
# View environment variables
kubectl describe pod -n gitea -l app=gitea | grep -A 20 "Environment"

# View current app.ini
kubectl exec -it deployment/gitea -n gitea -- cat /data/gitea/conf/app.ini
```

### Configuration Not Applied
- Verify the ConfigMap was generated: `kubectl get configmap -n gitea`
- Check pod restart: `kubectl get pods -n gitea`
- Review startup logs: `kubectl logs -n gitea -l app=gitea`


## External Dependencies

- **Database**: PostgreSQL instance in `postgres` namespace
- **Storage**: Longhorn distributed storage
- **Ingress**: Traefik with Let's Encrypt certificates
- **DNS**: External-DNS with Cloudflare integration