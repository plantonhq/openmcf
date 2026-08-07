# Production Private

A production-grade Cloud Composer environment with private networking,
multi-zone resilience, VPC-native ranges, control-plane allowlisting,
and scaled workloads.

## When to Use

- Production pipelines that must not expose a public Airflow endpoint
- Workloads requiring high availability across zones
- Networks with pre-planned secondary ranges for GKE pods and services
- Compliance requirements for private-only access

## Key Configuration

- **ENVIRONMENT_SIZE_MEDIUM + HIGH_RESILIENCE** — balanced capacity with
  multi-zone redundancy
- **VPC peering with private endpoint** — the Airflow UI is reachable
  only via private IP; master, Cloud SQL, and Composer components get
  dedicated CIDR ranges
- **ipAllocationPolicy with named secondary ranges** — pods and services
  land on ranges your network team pre-carved on the subnetwork (per
  range, a name and a CIDR are mutually exclusive)
- **masterAuthorizedNetworksConfig** — only listed CIDR blocks reach the
  GKE control plane that runs the Airflow workloads
- **Scaled workloads** — 2 schedulers, 2 triggerers, 2-6 autoscaling
  workers
- **Weekend maintenance window** — 12-hour windows on Saturday/Sunday

## What to Customize

- `projectId`, `nodeConfig.network`, `nodeConfig.subnetwork`,
  `nodeConfig.serviceAccount` — point at your project, VPC resources,
  and node identity (the service account must hold
  `roles/composer.worker`)
- `ipAllocationPolicy` — use your subnetwork's secondary range names, or
  swap to `clusterIpv4CidrBlock`/`servicesIpv4CidrBlock` if GKE should
  carve ranges itself
- `masterAuthorizedNetworksConfig.cidrBlocks` — your corporate/VPN
  ranges
- CIDR blocks under `privateEnvironmentConfig` — must not overlap
  anything else in your network

## Important Notes

- Environment creation takes 25-45 minutes.
- Networking (network, subnetwork, ranges, private config) is immutable
  after creation — plan the CIDRs before deploying.
- VPC peering requires firewall rules that allow the Composer ranges to
  reach your subnetwork.

## Related Presets

- **01-dev-small** — minimal development environment
- **03-enterprise-encrypted** — adds CMEK, data retention, and disaster
  recovery
