# KubernetesRookCephOperator Examples

This document provides practical examples for deploying the Rook Ceph Operator on Kubernetes clusters using the KubernetesRookCephOperator resource.

## Example 1: Basic Operator Deployment with Defaults

This is the simplest deployment using all default values. Suitable for development and testing.

```yaml
apiVersion: kubernetes.project-planton.org/v1
kind: KubernetesRookCephOperator
metadata:
  name: rook-ceph-operator-basic
spec:
  targetCluster:
    clusterName: "my-k8s-cluster"
  namespace:
    value: "rook-ceph"
  create_namespace: true
  container: {}
```

**What this creates:**
- Deploys Rook Ceph Operator version v1.16.6 (default) with default resource limits
- Enables RBD and CephFS CSI drivers
- Registers all Rook CRDs for Ceph cluster management
- Operator manages all CephCluster resources in the cluster

**When to use:**
- Quick proof-of-concept deployments
- Development environments
- Learning and experimentation

---

## Example 2: Production Operator with Custom Resources

Specify a specific operator version and increase resources for production environments.

```yaml
apiVersion: kubernetes.project-planton.org/v1
kind: KubernetesRookCephOperator
metadata:
  name: rook-ceph-operator-production
spec:
  targetCluster:
    clusterName: "production-cluster"
  namespace:
    value: "rook-ceph"
  create_namespace: true
  operator_version: "v1.16.6"
  crds_enabled: true
  container:
    resources:
      requests:
        cpu: 250m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

**What this creates:**
- Operator with increased baseline resources for production workloads
- Higher burst capacity for intensive reconciliation
- Explicit version pinning for reproducibility

**When to use:**
- Production environments
- Large-scale storage deployments
- High-availability requirements

---

## Example 3: Operator with Full CSI Configuration

Configure all CSI driver options for specific storage requirements.

```yaml
apiVersion: kubernetes.project-planton.org/v1
kind: KubernetesRookCephOperator
metadata:
  name: rook-ceph-operator-csi
spec:
  targetCluster:
    clusterName: "storage-cluster"
  namespace:
    value: "rook-ceph"
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 200m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 512Mi
  csi:
    enable_rbd_driver: true
    enable_cephfs_driver: true
    disable_csi_driver: false
    enable_csi_host_network: true
    provisioner_replicas: 3
    enable_csi_addons: true
    enable_nfs_driver: false
```

**What this creates:**
- Full CSI driver configuration with 3 provisioner replicas for HA
- CSI Addons enabled for additional features (encryption, reclaimspace)
- Host networking enabled for optimal performance

**When to use:**
- Production storage clusters
- Environments requiring CSI Addons features
- High-availability CSI provisioning

---

## Example 4: RBD-Only Deployment

Deploy operator with only block storage (RBD) enabled, no CephFS.

```yaml
apiVersion: kubernetes.project-planton.org/v1
kind: KubernetesRookCephOperator
metadata:
  name: rook-ceph-operator-rbd-only
spec:
  targetCluster:
    clusterName: "block-storage-cluster"
  namespace:
    value: "rook-ceph"
  create_namespace: true
  container: {}
  csi:
    enable_rbd_driver: true
    enable_cephfs_driver: false
    enable_nfs_driver: false
```

**What this creates:**
- RBD CSI driver only for block storage
- CephFS and NFS drivers disabled
- Reduced footprint for block-storage-only use cases

**When to use:**
- Block storage only requirements (databases, VMs)
- Minimal footprint deployments
- Environments not needing shared file systems

---

## Example 5: CephFS-Only Deployment

Deploy operator for file storage only scenarios.

```yaml
apiVersion: kubernetes.project-planton.org/v1
kind: KubernetesRookCephOperator
metadata:
  name: rook-ceph-operator-cephfs-only
spec:
  targetCluster:
    clusterName: "file-storage-cluster"
  namespace:
    value: "rook-ceph"
  create_namespace: true
  container: {}
  csi:
    enable_rbd_driver: false
    enable_cephfs_driver: true
    enable_nfs_driver: false
```

**What this creates:**
- CephFS CSI driver only for file storage
- RBD driver disabled
- Optimized for ReadWriteMany workloads

**When to use:**
- Shared file system requirements
- Web server content storage
- Multi-pod access to same data

---

## Example 6: Using Existing Namespace

Deploy operator to a pre-existing namespace managed separately.

```yaml
apiVersion: kubernetes.project-planton.org/v1
kind: KubernetesRookCephOperator
metadata:
  name: rook-ceph-operator-existing-ns
spec:
  targetCluster:
    clusterName: "my-cluster"
  namespace:
    value: "rook-ceph"
  create_namespace: false  # Namespace must already exist
  container:
    resources:
      requests:
        cpu: 200m
        memory: 128Mi
      limits:
        cpu: 500m
        memory: 512Mi
```

**What this creates:**
- Deploys operator to existing `rook-ceph` namespace
- Namespace must be created beforehand
- Ideal for GitOps workflows

**When to use:**
- Namespace managed by platform team
- Using KubernetesNamespace resource for namespace lifecycle
- Multi-component deployments

---

## Example 7: High-Performance CSI Configuration

Maximum performance configuration for demanding workloads.

```yaml
apiVersion: kubernetes.project-planton.org/v1
kind: KubernetesRookCephOperator
metadata:
  name: rook-ceph-operator-high-perf
spec:
  targetCluster:
    clusterName: "high-perf-cluster"
  namespace:
    value: "rook-ceph"
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 500m
        memory: 512Mi
      limits:
        cpu: 2000m
        memory: 2Gi
  csi:
    enable_rbd_driver: true
    enable_cephfs_driver: true
    disable_csi_driver: false
    enable_csi_host_network: true
    provisioner_replicas: 3
    enable_csi_addons: true
```

**What this creates:**
- High-resource operator for large-scale management
- Triple CSI provisioner replicas for availability
- Host networking for reduced latency

**When to use:**
- Enterprise storage deployments
- High-throughput requirements
- Large number of PVCs to manage

---

## Post-Deployment: Creating Ceph Clusters

After deploying the operator, create Ceph storage clusters:

### Minimal CephCluster Example

```yaml
apiVersion: ceph.rook.io/v1
kind: CephCluster
metadata:
  name: rook-ceph
  namespace: rook-ceph
spec:
  cephVersion:
    image: quay.io/ceph/ceph:v19.2.0
    allowUnsupported: false
  dataDirHostPath: /var/lib/rook
  skipUpgradeChecks: false
  continueUpgradeAfterChecksEvenIfNotHealthy: false
  mon:
    count: 3
    allowMultiplePerNode: false
  mgr:
    count: 2
    allowMultiplePerNode: false
  dashboard:
    enabled: true
    ssl: true
  storage:
    useAllNodes: true
    useAllDevices: true
```

### Production CephCluster with Specific Devices

```yaml
apiVersion: ceph.rook.io/v1
kind: CephCluster
metadata:
  name: rook-ceph
  namespace: rook-ceph
spec:
  cephVersion:
    image: quay.io/ceph/ceph:v19.2.0
  dataDirHostPath: /var/lib/rook
  mon:
    count: 3
    allowMultiplePerNode: false
  mgr:
    count: 2
    modules:
      - name: pg_autoscaler
        enabled: true
  dashboard:
    enabled: true
    ssl: true
  storage:
    useAllNodes: false
    useAllDevices: false
    nodes:
      - name: "node1"
        devices:
          - name: "sdb"
          - name: "sdc"
      - name: "node2"
        devices:
          - name: "sdb"
          - name: "sdc"
      - name: "node3"
        devices:
          - name: "sdb"
          - name: "sdc"
  resources:
    osd:
      limits:
        cpu: "2"
        memory: "4Gi"
      requests:
        cpu: "500m"
        memory: "2Gi"
```

### Creating Block Storage (RBD Pool)

```yaml
apiVersion: ceph.rook.io/v1
kind: CephBlockPool
metadata:
  name: replicapool
  namespace: rook-ceph
spec:
  failureDomain: host
  replicated:
    size: 3
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: rook-ceph-block
provisioner: rook-ceph.rbd.csi.ceph.com
parameters:
  clusterID: rook-ceph
  pool: replicapool
  imageFormat: "2"
  imageFeatures: layering
  csi.storage.k8s.io/provisioner-secret-name: rook-csi-rbd-provisioner
  csi.storage.k8s.io/provisioner-secret-namespace: rook-ceph
  csi.storage.k8s.io/controller-expand-secret-name: rook-csi-rbd-provisioner
  csi.storage.k8s.io/controller-expand-secret-namespace: rook-ceph
  csi.storage.k8s.io/node-stage-secret-name: rook-csi-rbd-node
  csi.storage.k8s.io/node-stage-secret-namespace: rook-ceph
  csi.storage.k8s.io/fstype: ext4
reclaimPolicy: Delete
allowVolumeExpansion: true
```

### Creating File Storage (CephFS)

```yaml
apiVersion: ceph.rook.io/v1
kind: CephFilesystem
metadata:
  name: myfs
  namespace: rook-ceph
spec:
  metadataPool:
    replicated:
      size: 3
  dataPools:
    - name: replicated
      failureDomain: host
      replicated:
        size: 3
  metadataServer:
    activeCount: 1
    activeStandby: true
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: rook-cephfs
provisioner: rook-ceph.cephfs.csi.ceph.com
parameters:
  clusterID: rook-ceph
  fsName: myfs
  pool: myfs-replicated
  csi.storage.k8s.io/provisioner-secret-name: rook-csi-cephfs-provisioner
  csi.storage.k8s.io/provisioner-secret-namespace: rook-ceph
  csi.storage.k8s.io/controller-expand-secret-name: rook-csi-cephfs-provisioner
  csi.storage.k8s.io/controller-expand-secret-namespace: rook-ceph
  csi.storage.k8s.io/node-stage-secret-name: rook-csi-cephfs-node
  csi.storage.k8s.io/node-stage-secret-namespace: rook-ceph
reclaimPolicy: Delete
allowVolumeExpansion: true
```

---

## Common Operations

### Checking Operator Status

```bash
# Get operator pod
kubectl get pods -n rook-ceph -l app=rook-ceph-operator

# View operator logs
kubectl logs -n rook-ceph -l app=rook-ceph-operator

# Check installed CRDs
kubectl get crds | grep ceph
```

### Managing Ceph Clusters

```bash
# List all CephClusters
kubectl get cephcluster -A

# Check cluster health
kubectl -n rook-ceph get cephcluster rook-ceph -o jsonpath='{.status.ceph.health}'

# Access Ceph toolbox for debugging
kubectl -n rook-ceph exec -it deploy/rook-ceph-tools -- ceph status
```

### Monitoring Storage

```bash
# Check OSD status
kubectl -n rook-ceph exec -it deploy/rook-ceph-tools -- ceph osd status

# View pool usage
kubectl -n rook-ceph exec -it deploy/rook-ceph-tools -- ceph df

# Check PG status
kubectl -n rook-ceph exec -it deploy/rook-ceph-tools -- ceph pg stat
```

---

## Resource Planning

### Operator Resource Usage

| Scenario | CPU Request | Memory Request | CPU Limit | Memory Limit |
|----------|-------------|----------------|-----------|--------------|
| Small    | 200m        | 128Mi          | 500m      | 512Mi        |
| Medium   | 250m        | 256Mi          | 1000m     | 1Gi          |
| Large    | 500m        | 512Mi          | 2000m     | 2Gi          |

### Ceph Cluster Resource Requirements

| Component | Min Replicas | CPU Request | Memory Request |
|-----------|--------------|-------------|----------------|
| MON       | 3            | 500m        | 1Gi            |
| MGR       | 2            | 500m        | 512Mi          |
| OSD       | 3            | 500m        | 2Gi            |
| MDS       | 2            | 500m        | 1Gi            |
| RGW       | 2            | 500m        | 512Mi          |

---

## Best Practices

1. **Three Monitors**: Always run 3 MON pods for quorum and high availability
2. **Dedicated Storage Nodes**: Use nodes with local SSDs/NVMe for OSDs
3. **Failure Domains**: Configure failure domains across racks/zones
4. **Resource Monitoring**: Set up Prometheus monitoring for Ceph metrics
5. **Regular Backups**: Implement backup strategy for critical data
6. **Version Pinning**: Pin operator and Ceph versions for stability
7. **Gradual Rollouts**: Upgrade one component at a time

---

## Additional Resources

- **Rook Documentation**: https://rook.io/docs/rook/latest/
- **Ceph Documentation**: https://docs.ceph.com/
- **Rook GitHub**: https://github.com/rook/rook
- **Research Documentation**: [docs/README.md](docs/README.md) - Detailed deployment patterns
- **Component README**: [README.md](README.md) - Full API reference
