## Minimal Replica Set Cluster (Development)

This example demonstrates a minimal MongoDB Atlas replica set cluster suitable for development environments.

```yaml
apiVersion: atlas.openmcf.org/v1
kind: MongodbAtlas
metadata:
  name: dev-cluster
spec:
  cluster_config:
    project_id: "507f1f77bcf86cd799439011"
    cluster_type: "REPLICASET"
    electable_nodes: 3
    priority: 7
    provider_name: "AWS"
    provider_instance_size_name: "M10"
    mongo_db_major_version: "7.0"
    cloud_backup: true
    auto_scaling_disk_gb_enabled: true
```

## Production Single-Region Cluster

This example shows a production-grade M30 cluster in a single region with cloud backups enabled.

```yaml
apiVersion: atlas.openmcf.org/v1
kind: MongodbAtlas
metadata:
  name: prod-cluster
spec:
  cluster_config:
    project_id: "507f1f77bcf86cd799439011"
    cluster_type: "REPLICASET"
    electable_nodes: 3
    priority: 7
    read_only_nodes: 0
    provider_name: "AWS"
    provider_instance_size_name: "M30"
    mongo_db_major_version: "7.0"
    cloud_backup: true
    auto_scaling_disk_gb_enabled: true
```

## Multi-Region Cluster with Read Replicas

This example demonstrates a multi-region deployment with a primary region and read-only replicas in a secondary region for disaster recovery and read scaling.

```yaml
apiVersion: atlas.openmcf.org/v1
kind: MongodbAtlas
metadata:
  name: multi-region-cluster
spec:
  cluster_config:
    project_id: "507f1f77bcf86cd799439011"
    cluster_type: "REPLICASET"
    # Primary region configuration
    electable_nodes: 3
    priority: 7
    provider_name: "AWS"
    provider_instance_size_name: "M50"
    mongo_db_major_version: "7.0"
    cloud_backup: true
    auto_scaling_disk_gb_enabled: true
```

## Sharded Cluster for High-Throughput Applications

This example shows a sharded cluster configuration for applications requiring horizontal scalability.

```yaml
apiVersion: atlas.openmcf.org/v1
kind: MongodbAtlas
metadata:
  name: sharded-cluster
spec:
  cluster_config:
    project_id: "507f1f77bcf86cd799439011"
    cluster_type: "SHARDED"
    electable_nodes: 3
    priority: 7
    provider_name: "GCP"
    provider_instance_size_name: "M50"
    mongo_db_major_version: "7.0"
    cloud_backup: true
    auto_scaling_disk_gb_enabled: true
```

## Global Cluster for Geographic Distribution

This example demonstrates a geographically distributed cluster for low-latency access across multiple regions.

```yaml
apiVersion: atlas.openmcf.org/v1
kind: MongodbAtlas
metadata:
  name: global-cluster
spec:
  cluster_config:
    project_id: "507f1f77bcf86cd799439011"
    cluster_type: "GEOSHARDED"
    electable_nodes: 3
    priority: 7
    provider_name: "AWS"
    provider_instance_size_name: "M30"
    mongo_db_major_version: "7.0"
    cloud_backup: true
    auto_scaling_disk_gb_enabled: true
```

## Multi-Cloud Cluster for Maximum Resilience

This example shows a multi-cloud deployment spanning AWS and GCP for provider-level disaster recovery.

```yaml
apiVersion: atlas.openmcf.org/v1
kind: MongodbAtlas
metadata:
  name: multi-cloud-cluster
spec:
  cluster_config:
    project_id: "507f1f77bcf86cd799439011"
    cluster_type: "REPLICASET"
    electable_nodes: 3
    priority: 7
    provider_name: "AWS"
    provider_instance_size_name: "M50"
    mongo_db_major_version: "7.0"
    cloud_backup: true
    auto_scaling_disk_gb_enabled: true
```

## High-Performance Cluster with Analytics Nodes

This example includes dedicated analytics nodes to isolate BI/reporting workloads from transactional traffic.

```yaml
apiVersion: atlas.openmcf.org/v1
kind: MongodbAtlas
metadata:
  name: analytics-cluster
spec:
  cluster_config:
    project_id: "507f1f77bcf86cd799439011"
    cluster_type: "REPLICASET"
    electable_nodes: 3
    priority: 7
    read_only_nodes: 2
    provider_name: "AWS"
    provider_instance_size_name: "M60"
    mongo_db_major_version: "7.0"
    cloud_backup: true
    auto_scaling_disk_gb_enabled: true
```

## CLI Workflows

### Validate Manifest

```bash
openmcf validate --manifest mongodb-atlas.yaml
```

### Deploy with Pulumi

```bash
openmcf pulumi up --manifest mongodb-atlas.yaml --stack org/project/stack
```

### Deploy with Terraform

```bash
openmcf tofu apply --manifest mongodb-atlas.yaml --auto-approve
```

### Check Cluster Status

```bash
openmcf get --manifest mongodb-atlas.yaml
```

### Update Cluster Configuration

```bash
# Edit your manifest file with desired changes
openmcf pulumi up --manifest mongodb-atlas.yaml --stack org/project/stack
```

### Destroy Cluster

```bash
openmcf pulumi destroy --manifest mongodb-atlas.yaml --stack org/project/stack
```
