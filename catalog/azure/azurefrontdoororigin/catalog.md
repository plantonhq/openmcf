# Azure Front Door Origin

Deploys a Front Door origin -- one backend inside an origin group, with hostname, ports, priority/weight, certificate name check, and optional Private Link. Private Link requires a PREMIUM profile and certificate name check enabled; Azure rejects the combination at apply otherwise. The component integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring to the parent origin group and backend hostnames.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Front Door Origin** -- a named child of an origin group pointing at a backend hostname
- **Ports and traffic role** -- HTTP/HTTPS ports (defaults 80/443), priority (1–5), weight (1–1000), and enabled switch
- **Certificate name check** -- whether Front Door validates the origin certificate against the host header (default on; required when Private Link is set)
- **Private Link** (optional, PREMIUM profile only) -- private connectivity to an App Service, storage, Container Apps environment, or Private Link Service

## The Origin in the Front Door Family

- **AzureFrontDoorOriginGroup** -- the parent pool, referenced by `originGroupId`
- **AzureFrontDoorRoute** -- may list this origin's `origin_id` for deploy sequencing (Azure never receives the list; it ensures the group is non-empty before the route applies)
- **Backend hosts** -- App Service, storage, Container Apps, or any hostname via ValueFromRef or a literal

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Azure Subscription

- **An Azure Front Door Origin Group** the origin nests under.
- **A PREMIUM Front Door profile** when using Private Link (Azure enforces this server-side).
- **A backend hostname** -- App Service default hostname, storage endpoint, or custom host.

## Deploy

### Console

Open the deployment store, find **Azure Front Door Origin**, and click **Deploy**. The wizard walks you through the parent group and name, the backend address (hostname and host header), traffic role (ports, priority, weight), and optional Private Link. Start from the **App Service Origin** preset for the common web-app pairing.

### CLI

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorOrigin
metadata:
  name: my-app-origin
  org: acme-corp
  env: prod
spec:
  originGroupId:
    valueFrom:
      kind: AzureFrontDoorOriginGroup
      name: my-app-backends
      fieldPath: status.outputs.origin_group_id
  originName: primary-app
  hostName:
    valueFrom:
      kind: AzureLinuxWebApp
      name: my-web-app
      fieldPath: status.outputs.default_hostname
```

```shell
planton apply -f front-door-origin.yaml
```

This creates a public-HTTPS origin tracking the App Service hostname. Leaving `originHostHeader` unset sends the same hostname -- the correct pairing for multi-tenant Azure backends.

### InfraChart

```yaml
spec:
  originGroupId:
    valueFrom:
      kind: AzureFrontDoorOriginGroup
      name: my-app-backends
      fieldPath: status.outputs.origin_group_id
  originName: primary-app
  hostName:
    valueFrom:
      kind: AzureLinuxWebApp
      name: my-web-app
      fieldPath: status.outputs.default_hostname
  priority: 1
  weight: 500
```

## Key Configuration

**Hostname** -- required. Reference a Cloud Resource output (App Service `default_hostname`, storage endpoint) or enter a literal FQDN. No default kind: name the target kind explicitly on each reference.

**Certificate name check** -- default on. Private Link requires it true; the wizard blocks the illegal pairing.

**Priority / weight** -- priority 1–5 (lower preferred); weight 1–1000 within a priority band. Defaults 1 / 500.

**Private Link** -- PREMIUM-only. Requires location, an ARM target ID (`/subscriptions/...`), and a target type unless the ID is a Private Link Service. Request message max 140 characters.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureFrontDoorOriginGroup** | `originGroupId` | `status.outputs.origin_group_id` |
| Backend hostname (any kind) | `hostName` / `originHostHeader` | Kind-specific hostname outputs |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `origin_id` | ARM resource ID of the origin | AzureFrontDoorRoute.`originIds` (deploy sequencing) |
| `origin_name` | The origin's name | Operator tooling |

## Presets

| Preset | Rank | Description |
|--------|------|-------------|
| App Service Origin | 1 | Public origin tracking an App Service hostname |
| Weighted Canary | 2 | Secondary origin with lower weight for canary traffic |
| Private Link Origin | 3 | PREMIUM Private Link to a private backend |

## Related Components

- **AzureFrontDoorOriginGroup** -- the parent pool
- **AzureFrontDoorRoute** -- sequences after origins exist
- **AzureLinuxWebApp** / **AzureStorageAccount** -- common hostname sources
