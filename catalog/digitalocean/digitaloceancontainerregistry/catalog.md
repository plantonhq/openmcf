# Container Registry on DigitalOcean

Deploys a private, OCI-compliant container registry on DigitalOcean for storing Docker images and Helm charts. Configures the registry name, subscription tier, and region, optionally mints Docker credentials with a controlled lifetime, and exposes the endpoint as a stack output for downstream workloads to pull images. Integrates with Planton's Provider Connections for DigitalOcean credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Container Registry** -- a `digitalocean_container_registry` resource with the specified name, subscription tier, and region
- **Docker Credentials** -- created only when `dockerCredentials` is set: a `digitalocean_container_registry_docker_credentials` resource minting a read-only or write credential with the configured expiry, exported through the `docker_credentials` stack output (a secret)

DigitalOcean restricts each account to a single container registry, and registry names are globally unique across ALL DigitalOcean accounts. Deploying a second DigitalOceanContainerRegistry resource on the same account will fail.

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **No existing container registry** on the target DigitalOcean account. DigitalOcean allows only one registry per account.
- **A supported region** for container registry storage (e.g., `nyc3`, `sfo3`, `fra1`, `sgp1`). Choose the region nearest to your Kubernetes clusters or CI/CD pipelines for lowest pull latency.

## Deploy

### Console

Open the deployment store, find **Container Registry on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Professional Container Registry** preset in the [Presets](#presets) tab for a production-ready configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanContainerRegistry
metadata:
  name: prod-registry
  org: acme-corp
  env: prod
spec:
  name: prod-registry
  subscriptionTier: professional
  region: nyc3
```

```shell
planton apply -f registry.yaml
```

This creates a professional-tier container registry in NYC3. Images are accessible at `registry.digitalocean.com/prod-registry/<repository>:<tag>`.

## Key Configuration

These are the most important decisions when configuring a container registry. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Subscription tier** -- The `subscriptionTier` field controls storage limits and pricing. `starter` is free with limited storage for development. `basic` provides moderate storage for small teams. `professional` offers the highest storage limits and is recommended for production teams pushing many images.

**Region** -- The `region` field determines where registry data is stored; omit it to let DigitalOcean choose (the chosen slug is reported through the `region` output). Choose the region closest to your DigitalOcean Kubernetes clusters or build pipelines to minimize image pull latency. The region cannot change after creation.

**Docker credentials** -- Set `dockerCredentials` to mint a registry credential: `write: true` for push access (default is read-only pull), and `expirySeconds` for a controlled lifetime (unset means the API maximum, roughly 50 years). No block means no long-lived token -- the secure default. Garbage collection of untagged images is an on-demand DigitalOcean action, not a registry attribute -- nothing in this manifest schedules it.

**Registry naming** -- The `name` field is globally unique across ALL DigitalOcean accounts and is used in image paths (`registry.digitalocean.com/<name>/<image>:<tag>`). Choose a stable name: it is the registry's resource identity, and changing it replaces the registry (all images gone).

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `registry_name` | Name of the created container registry (also its resource identifier) | Image path construction, Kubernetes `imagePullSecret` configuration |
| `server_url` | The registry host, always `registry.digitalocean.com` | Docker login |
| `endpoint` | Full registry endpoint (e.g., `registry.digitalocean.com/prod-registry`) | CI/CD pipeline push targets, App Platform image source |
| `region` | Region slug where the registry is hosted (reported by DigitalOcean) | Verifying data locality alignment with clusters |
| `docker_credentials` | Base64 Docker `config.json` (a secret); empty when credentials are not configured | Kubernetes `imagePullSecret`, CI/CD docker login |
| `credential_expiration_time` | RFC 3339 expiry of the minted credential | Rotation monitoring |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Professional registry for production** -- Professional tier with 30-day push credentials, suitable for teams pushing images frequently from CI/CD pipelines. Provides the highest storage limits. Start from the **Professional Container Registry** preset.

## Works With

This component operates independently and does not reference other components.