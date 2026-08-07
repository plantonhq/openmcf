# Valkey

Deploys Valkey — the Linux Foundation's Redis-compatible in-memory data store (the open-source successor every Redis client library speaks natively) — from the official Valkey Helm chart. Two topologies, chosen by the presence of the `replication` block: STANDALONE (one instance, the right shape for caches and development) or PRIMARY/REPLICA (one primary plus N streaming replicas with a dedicated read Service — read scaling, not automated failover). Several instances coexist per cluster: the release and every Service name are pinned to `metadata.name`.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Helm Release** (`<metadata.name>`) -- a Deployment (standalone) or StatefulSet (replication) running Valkey, with the typed `config` block rendered into the chart's valkey.conf string deterministically on both engines
- **Services** -- the write Service (`<name>`), and in replication mode the read Service (`<name>-read`, load balancing reads across all pods) and the headless Service (`<name>-headless`, direct pod discovery)
- **Auth Secret** (optional) -- when ACL users are declared, their passwords are materialized as the `<name>-auth` Kubernetes Secret (one key per username) the chart consumes; they never appear in rendered chart values
- **PersistentVolumeClaims** (optional) -- the dataset volume (standalone), or one per pod (replication, where persistence is required)
- **Metrics sidecar + ServiceMonitor** (optional) -- the redis_exporter sidecar with its metrics Service, and a ServiceMonitor for the Prometheus operator
- **Namespace** (optional) -- created with standard governance labels when `create_namespace` is true

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.

### Cluster

- A StorageClass for the persistence volume (or accept the cluster default). Replication mode REQUIRES persistence — replicas bootstrap by syncing the primary's on-disk dataset.
- For TLS: an existing `kubernetes.io/tls` Secret — issue a `KubernetesCertificate` and wire its Secret name through the reference (the cert-manager seam).
- With the Prometheus ServiceMonitor: the Prometheus operator CRDs — the release fails to install without them.

## Deploy

### Console

Open the deployment store, find **Valkey**, and click **Deploy**. The creation wizard walks you through placement, the chart pin, the standalone-vs-replication topology decision, persistence, the valkey.conf memory and durability dials, ACL authentication, TLS, exposure, observability, and scheduling. Start from the **Single Instance** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesValkey
metadata:
  name: sessions-cache
  org: acme-corp
  env: prod
spec:
  namespace:
    value: my-app
  persistence:
    size: 5Gi
  config:
    appendOnly: true
    maxMemory: 256mb
    maxMemoryPolicy: allkeys-lru
  auth:
    users:
      - name: default
        password: $secret/valkey-default-password
```

```shell
planton apply -f valkey.yaml
```

Applications then reach the store at the exported in-cluster endpoint, authenticating as `default` with the password from the `sessions-cache-auth` Secret.

## Key Configuration

These are the most important decisions when configuring Valkey. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Topology by presence** -- omitting `replication` means standalone; declaring it means one primary plus N replicas, a read Service, and REQUIRED persistence. Note what replication is NOT: automated failover. The chart ships no Sentinel — if the primary pod dies, Kubernetes restarts it and replicas re-attach, but no replica is promoted. Durability through a restart comes from persistence, not promotion.

**Authentication is declared, never defaulted** -- the chart ships with auth OFF: anyone who can reach the Service can read and write every key. Declaring ACL users turns it on, and the `default` user must be among them — without it, unauthenticated clients keep full access. Passwords are managed-secret references materialized into the `<name>-auth` Secret.

**The memory ceiling and eviction pair** -- `config.maxMemory` bounds the dataset; without it the pod's memory limit is the only bound, and hitting THAT is an OOM kill of the whole store, not an eviction. Pair the ceiling with an eviction policy (`allkeys-lru` for caches); `noeviction` — the server default — makes writes fail at the ceiling instead, which is right only when losing a key would be data loss.

**The durability posture** -- `appendOnly` logs every write and replays it on restart: paired with a persistence volume, restarts are lossless. RDB snapshots run on the server's built-in schedule, on custom save points, or not at all (`snapshotsDisabled`) — save points and the disable flag are mutually exclusive.

**Write safety in replication** -- `minReplicasToWrite` makes the primary refuse writes unless that many replicas are in sync, so a partitioned primary stops accepting writes replicas would never see. Applications see errors instead of silent divergence.

**Keep-on-uninstall** -- standalone only: keep the PVC (and the dataset) when the release is uninstalled. Off — the chart default — means the data dies with the resource.

**Exposure is composed** -- the store is in-cluster plumbing reachable at the exported endpoint; the LoadBalancer arm exists for managed-cloud recipes carried by the Service annotations. Compose a first-class exposure kind for anything else, and never expose an unauthenticated store.

**The escape hatch** -- `helm_values` carries additional chart values as a YAML document, merged LAST — never the substitute for typed fields, and never a place for secrets.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | The namespace the instance is deployed into |
| `spec.persistence.storage_class` / `spec.replication.persistence.storage_class` | KubernetesStorageClass (`status.outputs.storage_class_name`) | The StorageClass backing the dataset volume |
| `spec.tls.certificate_secret` | KubernetesCertificate (`status.outputs.secret_name`) | The kubernetes.io/tls Secret serving client TLS |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | Namespace the instance runs in | Debugging and composition |
| `service` | The write Service name (`<metadata.name>`) | Where applications send writes |
| `read_service` | The read Service name (`<metadata.name>-read`); empty standalone or when disabled | Read-heavy application paths |
| `headless_service` | The headless Service name (`<metadata.name>-headless`); replication mode only | Direct pod discovery |
| `kube_endpoint` | `<service>.<namespace>.svc.cluster.local:<port>` | The connection string applications consume |
| `port_forward_command` | kubectl port-forward one-liner | Reaching the store from a workstation |
| `username` | The ACL username applications authenticate with; empty when auth is off | Client configuration |
| `password_secret` | The Kubernetes Secret name + key holding that user's password | Mounting the credential without copying its value |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Single-instance cache** -- one instance with append-only persistence, a memory ceiling with LRU eviction, and ACL auth. Start from the **Single Instance** preset.

**Production read scaling** -- a primary plus two replicas with per-pod persistence, write safety, the read Service, and a PodDisruptionBudget. Start from the **Persistent With Replicas** preset.

## Works With

- **Kubernetes Namespace** -- the placement target; deploy the cache beside the application that uses it.
- **Kubernetes Certificate** -- issues and rotates the kubernetes.io/tls Secret the TLS block references.
- **Kubernetes Deployment / StatefulSet** -- the applications that consume the exported endpoint and the auth Secret.
- **Metrics / monitoring stack** -- the ServiceMonitor needs the Prometheus operator CRDs and turns cache health (memory, evictions, hit ratio, replica lag) into alerts.
