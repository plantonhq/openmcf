# KubernetesRookCephCluster Examples

This document provides practical examples for deploying Rook Ceph storage clusters on Kubernetes.

## Example 1: Minimal Cluster with Block Storage

Deploy a minimal Ceph cluster with one block storage pool. Suitable for development and testing.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesRookCephCluster
metadata:
  name: ceph-cluster-dev
spec:
  namespace:
    value: "rook-ceph"
  create_namespace: true
  block_pools:
    - name: "ceph-blockpool"
      storage_class:
        name: "ceph-block"
        is_default: true
```

**What this creates:**
- CephCluster using all nodes and all devices (defaults)
- Single CephBlockPool with 3x replication
- StorageClass `ceph-block` set as cluster default
- Dashboard enabled by default

**When to use:**
- Development environments
- Quick proof-of-concept
- Testing Ceph storage

---

## Example 2: Production Cluster with All Storage Types

Full-featured production deployment with block, file, and object storage.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesRookCephCluster
metadata:
  name: ceph-cluster-production
spec:
  namespace:
    value: "rook-ceph"
  create_namespace: true
  helm_chart_version: "v1.16.6"
  ceph_image:
    repository: "quay.io/ceph/ceph"
    tag: "v19.2.3"
  cluster:
    data_dir_host_path: "/var/lib/rook"
    mon:
      count: 3
      allow_multiple_per_node: false
    mgr:
      count: 2
      allow_multiple_per_node: false
    storage:
      use_all_nodes: true
      use_all_devices: true
    network:
      enable_encryption: false
      enable_compression: false
  block_pools:
    - name: "ceph-blockpool"
      failure_domain: "host"
      replicated_size: 3
      storage_class:
        name: "ceph-block"
        is_default: true
        reclaim_policy: "Delete"
        allow_volume_expansion: true
  filesystems:
    - name: "ceph-filesystem"
      metadata_pool_replicated_size: 3
      data_pool_replicated_size: 3
      active_mds_count: 1
      active_standby: true
      storage_class:
        name: "ceph-filesystem"
        reclaim_policy: "Delete"
  object_stores:
    - name: "ceph-objectstore"
      gateway_port: 80
      gateway_instances: 2
      preserve_pools_on_delete: true
      storage_class:
        name: "ceph-bucket"
  enable_dashboard: true
  enable_toolbox: true
  enable_monitoring: true
```

**What this creates:**
- CephCluster with 3 MON + 2 MGR daemons
- Block storage pool with StorageClass
- CephFS filesystem with MDS for shared storage
- Object store with RGW for S3-compatible storage
- Toolbox pod for debugging
- Prometheus monitoring enabled

**When to use:**
- Production environments
- Multi-storage-type requirements
- Enterprise deployments

---

## Example 3: Block Storage Only

Optimized deployment for block storage use cases like databases.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesRookCephCluster
metadata:
  name: ceph-block-only
spec:
  namespace:
    value: "rook-ceph"
  create_namespace: true
  cluster:
    mon:
      count: 3
    mgr:
      count: 2
    storage:
      use_all_nodes: true
      use_all_devices: true
  block_pools:
    - name: "fast-block"
      failure_domain: "host"
      replicated_size: 3
      storage_class:
        name: "ceph-block-fast"
        is_default: true
        volume_binding_mode: "WaitForFirstConsumer"
    - name: "standard-block"
      replicated_size: 2
      storage_class:
        name: "ceph-block-standard"
  enable_dashboard: true
```

**What this creates:**
- Two block pools with different replication factors
- Fast pool (3x replication) for critical workloads
- Standard pool (2x replication) for less critical data
- No CephFS or object storage overhead

**When to use:**
- Database workloads (PostgreSQL, MySQL, MongoDB)
- Stateful applications needing block storage
- Performance-focused deployments

---

## Example 4: CephFS for Shared Storage

Deployment focused on shared filesystem for ReadWriteMany workloads.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesRookCephCluster
metadata:
  name: ceph-filesystem-cluster
spec:
  namespace:
    value: "rook-ceph"
  create_namespace: true
  cluster:
    mon:
      count: 3
    mgr:
      count: 2
    storage:
      use_all_nodes: true
      use_all_devices: true
  filesystems:
    - name: "shared-fs"
      metadata_pool_replicated_size: 3
      data_pool_replicated_size: 3
      failure_domain: "host"
      active_mds_count: 2
      active_standby: true
      mds_resources:
        limits:
          cpu: "2000m"
          memory: "4Gi"
        requests:
          cpu: "1000m"
          memory: "2Gi"
      storage_class:
        name: "ceph-filesystem"
        is_default: true
  enable_dashboard: true
```

**What this creates:**
- CephFS filesystem with 2 active MDS + standby
- High-resource MDS for performance
- StorageClass for ReadWriteMany volumes

**When to use:**
- Web server content storage
- Shared ML training data
- Multi-pod shared storage requirements

---

## Example 5: Object Storage for S3-Compatible Workloads

Deployment focused on S3-compatible object storage.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesRookCephCluster
metadata:
  name: ceph-object-cluster
spec:
  namespace:
    value: "rook-ceph"
  create_namespace: true
  cluster:
    mon:
      count: 3
    mgr:
      count: 2
    storage:
      use_all_nodes: true
      use_all_devices: true
  object_stores:
    - name: "s3-store"
      metadata_pool_replicated_size: 3
      data_pool_erasure_data_chunks: 2
      data_pool_erasure_coding_chunks: 1
      failure_domain: "host"
      gateway_port: 80
      gateway_instances: 3
      preserve_pools_on_delete: true
      gateway_resources:
        limits:
          cpu: "2000m"
          memory: "2Gi"
        requests:
          cpu: "1000m"
          memory: "1Gi"
      storage_class:
        name: "ceph-bucket"
  enable_dashboard: true
```

**What this creates:**
- Object store with erasure coding for efficiency
- 3 RGW instances for high availability
- StorageClass for ObjectBucketClaims

**When to use:**
- Backup storage
- Log aggregation
- Media/asset storage
- Cloud-native applications using S3 API

---

## Example 6: Specific Node Storage Configuration

Deploy on specific nodes with specific devices.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesRookCephCluster
metadata:
  name: ceph-specific-nodes
spec:
  namespace:
    value: "rook-ceph"
  create_namespace: true
  cluster:
    data_dir_host_path: "/var/lib/rook"
    mon:
      count: 3
    mgr:
      count: 2
    storage:
      use_all_nodes: false
      use_all_devices: false
      nodes:
        - name: "storage-node-1"
          devices:
            - "sdb"
            - "sdc"
            - "sdd"
        - name: "storage-node-2"
          devices:
            - "sdb"
            - "sdc"
            - "sdd"
        - name: "storage-node-3"
          devices:
            - "sdb"
            - "sdc"
            - "sdd"
  block_pools:
    - name: "ceph-blockpool"
      storage_class:
        name: "ceph-block"
        is_default: true
  enable_dashboard: true
```

**What this creates:**
- OSDs only on specified storage nodes
- Uses only specified devices on each node
- Isolates storage from compute workloads

**When to use:**
- Dedicated storage nodes
- Hybrid compute/storage clusters
- Specific disk allocation requirements

---

## Example 7: Device Filter Configuration

Use device filter patterns instead of explicit device lists.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesRookCephCluster
metadata:
  name: ceph-device-filter
spec:
  namespace:
    value: "rook-ceph"
  create_namespace: true
  cluster:
    mon:
      count: 3
    mgr:
      count: 2
    storage:
      use_all_nodes: true
      use_all_devices: false
      device_filter: "^sd[b-z]$"
  block_pools:
    - name: "ceph-blockpool"
      storage_class:
        name: "ceph-block"
        is_default: true
  enable_dashboard: true
```

**What this creates:**
- Uses devices matching pattern (sdb, sdc, etc.)
- Excludes sda (typically OS disk)
- Applies to all nodes

**When to use:**
- Uniform disk naming across nodes
- Exclude OS disks automatically
- Simplified configuration for homogeneous hardware

---

## Example 8: High-Availability Configuration

Maximum availability with zone-aware failure domains.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesRookCephCluster
metadata:
  name: ceph-ha-cluster
spec:
  namespace:
    value: "rook-ceph"
  create_namespace: true
  cluster:
    mon:
      count: 5
      allow_multiple_per_node: false
    mgr:
      count: 2
      allow_multiple_per_node: false
    storage:
      use_all_nodes: true
      use_all_devices: true
    network:
      enable_encryption: true
      require_msgr2: true
    resources:
      mon:
        limits:
          cpu: "2000m"
          memory: "2Gi"
        requests:
          cpu: "1000m"
          memory: "1Gi"
      mgr:
        limits:
          cpu: "1000m"
          memory: "1Gi"
        requests:
          cpu: "500m"
          memory: "512Mi"
      osd:
        limits:
          cpu: "2000m"
          memory: "4Gi"
        requests:
          cpu: "1000m"
          memory: "4Gi"
  block_pools:
    - name: "ha-blockpool"
      failure_domain: "zone"
      replicated_size: 3
      storage_class:
        name: "ceph-block-ha"
        is_default: true
        volume_binding_mode: "WaitForFirstConsumer"
  filesystems:
    - name: "ha-filesystem"
      failure_domain: "zone"
      active_mds_count: 2
      active_standby: true
      storage_class:
        name: "ceph-filesystem-ha"
  enable_dashboard: true
  enable_monitoring: true
  enable_toolbox: true
```

**What this creates:**
- 5 MON pods for stronger quorum
- Zone-level failure domains
- Encrypted network traffic
- Higher resource allocations
- Full monitoring and toolbox

**When to use:**
- Production critical workloads
- Multi-zone/multi-rack clusters
- Strict availability requirements

---

## Common Operations

### Accessing the Ceph Dashboard

```bash
# Port-forward to dashboard service
kubectl port-forward svc/rook-ceph-mgr-dashboard -n rook-ceph 7000:7000

# Get admin password
kubectl -n rook-ceph get secret rook-ceph-dashboard-password \
  -o jsonpath="{['data']['password']}" | base64 -d

# Access at https://localhost:7000
```

### Using the Toolbox

```bash
# Access toolbox pod
kubectl -n rook-ceph exec -it deploy/rook-ceph-tools -- bash

# Inside toolbox:
ceph status
ceph osd status
ceph df
ceph health detail
```

### Creating a PVC

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: my-block-pvc
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: ceph-block
  resources:
    requests:
      storage: 10Gi
```

### Creating an Object Bucket Claim

```yaml
apiVersion: objectbucket.io/v1alpha1
kind: ObjectBucketClaim
metadata:
  name: my-bucket
spec:
  storageClassName: ceph-bucket
  generateBucketName: my-bucket
```

---

## Additional Resources

- **Component README**: [README.md](README.md)
- **Research Documentation**: [docs/README.md](docs/README.md)
- **Rook Documentation**: https://rook.io/docs/rook/latest/
- **Ceph Documentation**: https://docs.ceph.com/
