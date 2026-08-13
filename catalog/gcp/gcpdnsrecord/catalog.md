# DNS Record on Google Cloud

Deploys a DNS record set in an existing Google Cloud DNS Managed Zone — one (name, type) pair answering either with static values (round-robin) or with exactly one routing policy that steers each query by weight, caller geography, or target health. Any record type the Cloud DNS API supports passes through (A, AAAA, CNAME, MX, TXT, SRV, NS, CAA, and newer types like HTTPS and SVCB). The name and values are reference-capable, so certificate-validation records compose directly from a GcpCertManagerDnsAuthorization's outputs. Integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to GCP projects, DNS zones, health checks, and reserved addresses.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cloud DNS Record Set** -- a record set in the specified Managed Zone, configured with the chosen record type, fully qualified domain name, and TTL
- **Static values (round-robin)** -- when `values` are provided, a multi-value record set answering every resolver identically
- **Routing policy** -- when `routingPolicy` is set instead, weighted round-robin, geolocation, or primary/backup failover steering, optionally health-checking internal load balancer frontends and external endpoints

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** containing the Cloud DNS Managed Zone where the record will be created. Provide the project ID directly or reference a GcpProject Cloud Resource via ValueFromRef.
- **An existing Cloud DNS Managed Zone** for the target domain. Provide the zone name directly or reference a GcpDnsZone Cloud Resource via ValueFromRef.
- **Cloud DNS API** (`dns.googleapis.com`) enabled in the target project.

## Deploy

### Console

Open the deployment store, find **DNS Record on Google Cloud**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **A Record** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDnsRecord
metadata:
  name: api-a-record
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  managedZone:
    value: "example-zone"
  type: A
  name: "api.example.com."
  values:
    - "34.120.0.1"
```

```shell
planton apply -f dns-record.yaml
```

This creates an A record pointing `api.example.com.` to `34.120.0.1` with the default 300-second TTL. No round-robin or custom TTL is configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the DNS record to a GCP project and DNS zone deployed in the same InfraPipeline:

```yaml
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: production-project
      fieldPath: status.outputs.project_id
  managedZone:
    valueFrom:
      kind: GcpDnsZone
      name: example-zone
      fieldPath: status.outputs.zone_name
```

The InfraPipeline resolves the dependency graph, deploys the project and DNS zone first, then provisions the DNS record with the resolved values.

## Key Configuration

These are the most important decisions when configuring a DNS record. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Record type** -- Uppercase, e.g. `A`, `CNAME`, `TXT` via the `type` field. Cloud DNS accepts any type its API supports — the string passes through as-is, so newer types (HTTPS, SVCB, TLSA) need no platform change. The type determines the format of values: IPv4 addresses for A records, hostnames with trailing dots for CNAME records, priority-prefixed servers for MX records.

**Fully qualified domain name** -- The `name` field must be a valid FQDN ending with a trailing dot (e.g., `www.example.com.`). Wildcard records are supported with a `*.` prefix. The name is reference-capable — point it at a GcpCertManagerDnsAuthorization's `dns_record_name` output to compose certificate-validation records with zero hand-copying.

**Static values XOR routing policy** -- Exactly one. `values` entries answer as a round-robin set (each entry a literal or a reference, e.g. an authorization's `dns_record_data` output). `routingPolicy` instead steers each query: `wrr` splits by weight (a weight of 0 stages an entry — the DNS canary lever), `geo` answers with the entry nearest the caller, `primaryBackup` serves primaries while healthy and falls back to regional backups (with an optional `trickleRatio` keeping the backup path warm). Health-checked internal load balancer targets carry their own health signal; external endpoints need the policy-level `healthCheck` reference.

**TTL** -- `ttlSeconds` defaults to 300 seconds when omitted; an explicit 0 disables caching. Use lower values (60s) during migrations and for failover policies, and higher values (3600s) for stable production records to reduce DNS query load.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** | `projectId` | `status.outputs.project_id` |
| **GcpDnsZone** | `managedZone` | `status.outputs.zone_name` |
| **GcpCertManagerDnsAuthorization** | `name` / `values[]` | `status.outputs.dns_record_name` / `status.outputs.dns_record_data` |
| **GcpHealthCheck** | `routingPolicy.healthCheck` | `status.outputs.self_link` |
| **GcpAddress** | `routingPolicy.*.healthCheckedTargets.internalLoadBalancers[].ipAddress` | `status.outputs.address` |
| **GcpVpcNetwork** | `routingPolicy.*.healthCheckedTargets.internalLoadBalancers[].networkUrl` | `status.outputs.network_self_link` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `fqdn` | Fully qualified domain name of the created record | Application configuration, monitoring |
| `record_type` | DNS record type that was created | Audit, documentation |
| `managed_zone` | Name of the managed zone containing this record | Related record creation |
| `project_id` | GCP project ID where the record was created | Cross-project references |
| `ttl_seconds` | TTL in seconds for the DNS record | Cache behavior validation |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**A record** -- Maps a hostname to an IPv4 address with a 5-minute TTL. The most common DNS record type for pointing domains to load balancers, VMs, or external services. Start from the **A Record** preset.

**CNAME record** -- Aliases one hostname to another, allowing the target IP to change without updating the record. Use for `www` subdomains, CDN frontends, or SaaS provider endpoints. Start from the **CNAME Record** preset.

**Weighted canary** -- A `wrr` routing policy splitting traffic 90/10 between stable and canary answers; shift by editing weights, stage new entries at weight 0. Start from the **Weighted Canary** preset.

**Certificate validation** -- A CNAME whose `name` and `values` reference a GcpCertManagerDnsAuthorization's outputs — the three-node composition (authorization → this record → GcpCertManagerCert) validates certificates before any load balancer serves the domain.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project containing the DNS Managed Zone
- [**DNS Zone on Google Cloud**](/cloud-catalog/gcp-dns-zone) -- provides the Managed Zone where the record is created
- [**GCP Cert Manager DNS Authorization**](/cloud-catalog/gcp-cert-manager-dns-authorization) -- exports the validation CNAME this record serves
- [**GCP Health Check**](/cloud-catalog/gcp-health-check) -- probes external endpoints for routing-policy health withdrawal