# Kubernetes SeaweedFS

Deploys SeaweedFS — an Apache-2.0, S3-compatible object store — from
the official SeaweedFS Helm chart. One resource is one whole store:
master servers coordinate the cluster, volume servers hold the object
bytes, and the filer serves the namespace and the S3 API — each tier
sized and scaled independently. The S3 gateway is on by default with
credentials materialized in a Kubernetes Secret, buckets declared in
the spec are created at install, storage deliberately rides
PersistentVolumeClaims instead of the chart's hostPath default, and
everything stays ClusterIP — external reachability composes from
ingress and gateway kinds.

> **Why an in-cluster object store**: applications need S3 for
> artifacts, backups, datasets and media even where no cloud bucket
> service exists — on-prem fleets, air-gapped clusters, dev
> environments. SeaweedFS serves the S3 API from inside the cluster
> with a deliberately small footprint: three small pods on default
> volumes make a working store.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **Admin auth Secret** (`<name>-admin-auth`, when the console is
  enabled without `existing_auth_secret`) — keys `user`/`password`;
  the console is never installed open
- **Helm release** (official `seaweedfs` chart, pinned 4.40.0, named
  `metadata.name`): the master, volume and filer tiers on
  PersistentVolumeClaims (5Gi/30Gi/10Gi defaults), the S3 gateway
  (embedded on the filer, or its own Deployment with `s3.dedicated`),
  the `<name>-s3` Service (port 8333), the `<name>-s3-secret`
  credentials Secret (when auth is on — the default; stable across
  upgrades, kept on uninstall), the install hook that creates the
  buckets declared in `s3.buckets`, and the admin console when
  enabled

## Prerequisites

- A Kubernetes namespace that already exists, or set
  `create_namespace`
- A StorageClass for the data volumes (most managed clusters provide
  a default; or reference a KubernetesStorageClass per tier)
- For `s3.existing_config_secret`: a Secret carrying the chart's
  `seaweedfs_s3_config` key (the inline JSON identities document)
- For `admin.existing_auth_secret`: a Secret carrying `user` and
  `password` keys
- For `replication`: enough volume servers (and racks / data centers,
  per the placement code) to place every copy
- For `service_monitor_enabled`: the Prometheus Operator CRDs

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesSeaweedFs
metadata:
  name: object-store
spec:
  namespace:
    value: object-store
  create_namespace: true
  s3:
    buckets:
      - name: app-data
```

SDKs point at the exported `s3_endpoint`
(`http://object-store-s3.object-store.svc.cluster.local:8333`,
path-style addressing) with the admin access-key pair from the
`object-store-s3-secret` Secret; a read-only pair lives alongside it
for consumers that must not write.

## Configuration

### Topology and scaling

Empty tiers give 1 master / 1 volume server / 1 filer — a working
single-node store. For HA: 3 masters (a Raft quorum — use an odd
count), volume servers sized for your data (the only tier that grows
with stored bytes), and a `replication` placement code (`"001"` = one
extra copy on another server, requiring at least 2 volume servers).
Keep 1 filer unless a shared external metadata store is wired via
`helm_values` — the embedded leveldb default is per-pod.

### S3, credentials and buckets

The gateway is on with auth on by default — the chart materializes
admin and read-only credential pairs in `<name>-s3-secret`. Declare
`s3.dedicated` to run the gateway as its own Deployment and scale the
API independently of metadata; both shapes expose the same `<name>-s3`
Service. Buckets in `s3.buckets` are typed (name, `anonymous_read`,
`ttl` like `"7d"`, `object_lock` — irreversible, forces versioning —
and `versioning`) and created by the chart's install hook; removing
an entry later does NOT delete the bucket or its data.

### Storage

The chart's out-of-the-box storage is hostPath (bare-metal grain);
this component deliberately maps every data volume to a
PersistentVolumeClaim and logs to emptyDir. Declare `size` and
optionally `storage_class` per tier; empty means the tier default on
the cluster's default class.

### Exposure

Everything is in-cluster only (ClusterIP). Compose a
KubernetesIngress or Gateway API route over the `<name>-s3` Service
for external S3 clients, and over the admin Service for the console.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the store runs in |
| `release_name` | Helm release name (= `metadata.name`) |
| `s3_endpoint` | In-cluster S3 endpoint (`http://<name>-s3.<ns>.svc.cluster.local:8333`); empty when the gateway is disabled |
| `s3_credentials_secret_name` | Secret holding the S3 credentials (`<name>-s3-secret`, or the referenced existing config secret); empty when auth is disabled |
| `filer_service_name` | Filer Service (file namespace HTTP API, port 8888) |
| `master_service_name` | Master Service (cluster coordination, port 9333) |
| `admin_endpoint` | In-cluster admin console endpoint (port 23646); empty when the console is disabled |
| `admin_auth_secret_name` | Secret holding the console credentials (keys `user`/`password`); empty when the console is disabled |
| `port_forward_command` | Workstation S3 access when no exposure is composed |

## Related Components

- [KubernetesNamespace](/docs/catalog/kubernetes/kubernetesnamespace) —
  provides the target namespace via reference
- [KubernetesStorageClass](/docs/catalog/kubernetes/kubernetesstorageclass)
  — backs the per-tier data volumes via reference
- [KubernetesIngress](/docs/catalog/kubernetes/kubernetesingress) —
  composes external S3 or console exposure over the service handles

## Next Steps

Declare the buckets applications need where the store is declared —
the hook creates them at install, and consumers receive a bucket name
plus the credentials Secret, never bucket-creation power. Size the
volume tier for the data you expect and set a `replication` code once
there are enough volume servers to place the copies. Enable the admin
console for operations visibility — it always ships behind
credentials — and compose exposure only where S3 clients live outside
the cluster.
