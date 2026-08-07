# Azure Front Door Origin Group

Deploys a Front Door origin group -- a load-balanced backend pool with health probing, session affinity, and traffic-restoration timing. Origins (individual backends) and routes (URL patterns) both reference the group's `origin_group_id`. Renaming the group replaces it and every origin nested under it. The component integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring to the parent profile.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Front Door Origin Group** -- a named child of the profile that owns load-balancing and health-probe settings
- **Load-balancing dials** (optional) -- sample size, successful samples required, and additional latency for latency-based selection
- **Health probe** (optional) -- protocol, interval, request type, and path Front Door uses to decide which origins are healthy
- **Session affinity** -- sticky sessions when enabled (default on)
- **Traffic restoration timer** -- minutes to wait before sending traffic to a healed or newly added origin (0–50, default 10)

## The Origin Group in the Front Door Family

- **AzureFrontDoorProfile** -- the parent container, referenced by `profileId`
- **AzureFrontDoorOrigin** -- one backend inside this group (priority/weight, hostname, optional Private Link)
- **AzureFrontDoorRoute** -- forwards matched requests to this group's `origin_group_id`

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Azure Subscription

- **An Azure Front Door Profile** the group nests under. Reference an AzureFrontDoorProfile via ValueFromRef, or provide the profile ARM ID directly.

## Deploy

### Console

Open the deployment store, find **Azure Front Door Origin Group**, and click **Deploy**. The wizard walks you through the parent profile and name, health probing, then load balancing and session affinity. Start from the **Single Origin** preset when probes would only add load, or **Multi-Region Probed** when failover matters.

### CLI

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorOriginGroup
metadata:
  name: my-app-backends
  org: acme-corp
  env: prod
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: cdn-profile
      fieldPath: status.outputs.profile_id
  originGroupName: app-backends
```

```shell
planton apply -f front-door-origin-group.yaml
```

This creates a group with Azure's load-balancing defaults and no health probe -- the right shape for a single origin. Add origins next, then attach a route.

### InfraChart

```yaml
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: cdn-profile
      fieldPath: status.outputs.profile_id
  originGroupName: app-backends
  healthProbe:
    protocol: HTTPS
    intervalInSeconds: 30
    requestType: HEAD
    path: /health
```

## Key Configuration

**Health probe** -- omit for a single origin (probes only add load with nowhere to fail over). When set, protocol is required (HTTP or HTTPS), interval is 1–255 seconds, and path must start with `/`. Unspecified request type deploys HEAD.

**Load balancing** -- sample size, successful samples required, and additional latency (milliseconds). All optional with Azure defaults (4 / 3 / 50 ms). Tune only when multi-origin latency selection matters.

**Session affinity** -- sticky sessions for stateful apps. Defaults to true when omitted.

**Restore traffic timer** -- 0–50 minutes before a healed or new origin receives traffic. Defaults to 10.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureFrontDoorProfile** | `profileId` | `status.outputs.profile_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `origin_group_id` | ARM resource ID of the group | AzureFrontDoorOrigin.`originGroupId`; AzureFrontDoorRoute.`originGroupId` |
| `origin_group_name` | The group's name | Operator tooling |

## Presets

| Preset | Rank | Description |
|--------|------|-------------|
| Single Origin | 1 | Group without a probe -- single-backend shape |
| Multi-Region Probed | 2 | HTTPS health probe for failover pools |
| Stateful Web App | 3 | Session affinity retained for sticky sessions |

## Related Components

- **AzureFrontDoorProfile** -- the parent container
- **AzureFrontDoorOrigin** -- backends inside this group
- **AzureFrontDoorRoute** -- forwards traffic to this group
