# GCP DNS Record

Deploys an individual DNS record set within an existing Google Cloud DNS Managed Zone. A record set answers queries either with static values (round-robin) or with exactly one routing policy — weighted round robin, geolocation, or primary/backup failover with health-checked targets. All record types the Cloud DNS API supports are accepted (A, AAAA, CNAME, MX, TXT, SRV, NS, PTR, CAA, SOA, HTTPS, SVCB, and more).

## What Gets Created

When you deploy a GcpDnsRecord resource, Planton provisions:

- **DNS Record Set** — a `google_dns_record_set` resource in the specified managed zone, with the given type, FQDN, TTL, and either static values or a routing policy
- **Cloud DNS API enablement** — `dns.googleapis.com` is enabled on the target project (never disabled on destroy)

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId`
- **An existing Cloud DNS Managed Zone** — referenced via `managedZone`, either by direct name or as a foreign key to a GcpDnsZone resource
- **IAM permissions** to create and manage DNS record sets in the target managed zone

## Quick Start

Create a file `dns-record.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDnsRecord
metadata:
  name: app-a-record
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.GcpDnsRecord.app-a-record
spec:
  projectId: my-gcp-project-123
  managedZone: example-zone
  type: A
  name: app.example.com.
  values:
    - 203.0.113.10
```

Deploy:

```shell
planton apply -f dns-record.yaml
```

This creates an A record for `app.example.com.` pointing to `203.0.113.10` with the default TTL of 300 seconds.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `managedZone` | `StringValueOrRef` | Name of the Cloud DNS Managed Zone (the zone RESOURCE name, not the DNS name). Can reference a GcpDnsZone resource via `valueFrom`. | Required |
| `type` | `string` | DNS record type, uppercase (e.g., `A`, `AAAA`, `CNAME`, `MX`, `TXT`, `HTTPS`, `SVCB`). Any type the Cloud DNS API supports is accepted. | Required, uppercase alphanumeric |
| `name` | `string` | Fully qualified domain name for the record. Must end with a trailing dot (e.g., `www.example.com.`). | Required, must match valid FQDN pattern |

Exactly one of `values` or `routingPolicy` must be set:

| Field | Type | Description |
|-------|------|-------------|
| `values` | `string[]` | Static record values (RRDATA). For A records: IPv4 addresses. For AAAA: IPv6 addresses. For CNAME: target hostname with trailing dot. Multiple values create a round-robin record set. |
| `routingPolicy` | `object` | Query steering: exactly one of `wrr` (weighted round robin), `geo` (geolocation), or `primaryBackup` (failover). See below. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | GCP project ID where the managed zone exists. Can reference a GcpProject resource via `valueFrom`. If omitted, the provider's default project is used. |
| `ttlSeconds` | `int32` | `300` | Time to live in seconds — how long resolvers cache this record. Common values: 60 (fast failover), 300, 3600, 86400; NS records conventionally use 172800. |

### Routing Policy

| Field | Type | Description |
|-------|------|-------------|
| `routingPolicy.wrr[]` | `object[]` | Weighted round robin entries: `weight` (required, ≥ 0; traffic splits by weight ratio; 0 stages an entry with no traffic), `values`, and/or `healthCheckedTargets`. |
| `routingPolicy.geo[]` | `object[]` | Geolocation entries: `location` (a GCP location such as `us-east1`), `values`, and/or `healthCheckedTargets`. |
| `routingPolicy.enableGeoFencing` | `bool` | Geo routing only: unhealthy locations keep answering instead of failing over to the next-closest location. |
| `routingPolicy.primaryBackup` | `object` | Failover: `primary` (health-checked targets, required), `backupGeo[]` (regional fallback, required), `trickleRatio` (0.0–1.0 traffic to backups while primaries are healthy), `enableGeoFencingForBackups`. |
| `routingPolicy.healthCheck` | `StringValueOrRef` | Health check for public-IP (external endpoint) targets. Can reference a GcpHealthCheck. Internal load balancer targets carry their own implicit health signal. |

`healthCheckedTargets` (A/AAAA records only) holds `internalLoadBalancers[]` (each with `ipAddress` — referenceable to GcpAddress — `ipProtocol` tcp/udp, `networkUrl` — referenceable to GcpVpcNetwork — `port`, `project`, optional `loadBalancerType` and `region`) and/or `externalEndpoints[]` (public IPs, requiring `healthCheck`).

## Examples

### Simple A Record

An A record pointing a subdomain to a single IP address:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDnsRecord
metadata:
  name: web-a-record
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.GcpDnsRecord.web-a-record
spec:
  projectId: my-gcp-project-123
  managedZone: example-zone
  type: A
  name: www.example.com.
  values:
    - 203.0.113.10
  ttlSeconds: 300
```

### CNAME Record with Foreign Key References

A CNAME record that references Planton-managed GcpProject and GcpDnsZone resources instead of hardcoding identifiers:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDnsRecord
metadata:
  name: docs-cname
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpDnsRecord.docs-cname
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: my-project
      fieldPath: status.outputs.project_id
  managedZone:
    valueFrom:
      kind: GcpDnsZone
      name: example.com
      fieldPath: status.outputs.zone_name
  type: CNAME
  name: docs.example.com.
  values:
    - example.github.io.
  ttlSeconds: 3600
```

### Round-Robin A Record with Multiple IPs

An A record with multiple values for basic load distribution across servers:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDnsRecord
metadata:
  name: api-round-robin
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpDnsRecord.api-round-robin
spec:
  projectId: my-prod-project-456
  managedZone: example-zone
  type: A
  name: api.example.com.
  values:
    - 203.0.113.10
    - 203.0.113.11
    - 203.0.113.12
  ttlSeconds: 60
```

### MX Record for Email Routing

An MX record configuring mail delivery with primary and backup mail servers:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDnsRecord
metadata:
  name: mail-mx
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpDnsRecord.mail-mx
spec:
  projectId: my-prod-project-456
  managedZone: example-zone
  type: MX
  name: example.com.
  values:
    - "10 mail.example.com."
    - "20 mail2.example.com."
  ttlSeconds: 3600
```

### TXT Record for SPF and Domain Verification

A TXT record used for email sender policy and domain ownership verification:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDnsRecord
metadata:
  name: spf-txt
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpDnsRecord.spf-txt
spec:
  projectId: my-prod-project-456
  managedZone: example-zone
  type: TXT
  name: example.com.
  values:
    - "v=spf1 include:_spf.google.com ~all"
  ttlSeconds: 3600
```

### Weighted Canary Rollout

A weighted round-robin policy sending 5% of traffic to a new backend:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDnsRecord
metadata:
  name: api-canary
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpDnsRecord.api-canary
spec:
  projectId: my-prod-project-456
  managedZone: example-zone
  type: A
  name: api.example.com.
  routingPolicy:
    wrr:
      - weight: 95
        values: ["203.0.113.10"]
      - weight: 5
        values: ["203.0.113.20"]
  ttlSeconds: 60
```

### Geolocation Routing

Answers each query from the location nearest the caller:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpDnsRecord
metadata:
  name: app-geo
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpDnsRecord.app-geo
spec:
  projectId: my-prod-project-456
  managedZone: example-zone
  type: A
  name: app.example.com.
  routingPolicy:
    geo:
      - location: us-east1
        values: ["203.0.113.10"]
      - location: europe-west3
        values: ["198.51.100.10"]
  ttlSeconds: 300
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `fqdn` | `string` | The fully qualified domain name of the created DNS record (e.g., `www.example.com.`) |
| `record_type` | `string` | The DNS record type that was created (e.g., `A`, `CNAME`, `TXT`) |
| `managed_zone` | `string` | The name of the managed zone containing this record |
| `project_id` | `string` | The GCP project ID where the record was created |
| `ttl_seconds` | `int32` | The TTL (time to live) in seconds for the DNS record |

## Related Components

- [GcpDnsZone](/docs/catalog/gcp/gcpdnszone) — creates the Cloud DNS Managed Zone where records are hosted
- [GcpProject](/docs/catalog/gcp/gcpproject) — provides the GCP project referenced by `projectId`
- [GcpServiceAccount](/docs/catalog/gcp/gcpserviceaccount) — creates service accounts that can be granted DNS management permissions
- [GcpGkeCluster](/docs/catalog/gcp/gcpgkecluster) — deploys GKE clusters whose ingress endpoints are commonly referenced by A or CNAME records
