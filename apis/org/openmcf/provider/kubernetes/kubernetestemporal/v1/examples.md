# Temporal on Kubernetes Examples

Below are practical examples demonstrating how to use the `TemporalKubernetes` component within your OpenMCF
deployments. These examples illustrate common configurations and scenarios to help you quickly integrate Temporal
workflows into your Kubernetes environments.

---

## Example 1: Basic Deployment with Default Cassandra Database

This example shows a simple Temporal deployment using an internal Cassandra database and exposing Temporal frontend and
web UI through ingress.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTemporal
metadata:
  name: temporal-basic
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: temporal-basic
  create_namespace: true
  database:
    backend: cassandra
  ingress:
    frontend:
      enabled: true
      grpc_hostname: temporal-frontend.example.com
    web_ui:
      enabled: true
      hostname: temporal-ui.example.com
```

**Use Case:**

- Ideal for quickly setting up Temporal workflows for development or testing purposes without external database
  dependencies.

---

## Example 2: Deployment with External PostgreSQL Database (Using String Value)

Deploy Temporal using an external PostgreSQL database for development/testing with a plain text password.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTemporal
metadata:
  name: temporal-postgres-dev
spec:
  target_cluster:
    cluster_name: dev-gke-cluster
  namespace:
    value: temporal-dev
  create_namespace: true
  database:
    backend: postgresql
    external_database:
      host: postgres-db.example.com
      port: 5432
      username: temporaluser
      password:
        stringValue: securepassword
  ingress:
    frontend:
      enabled: true
      grpc_hostname: temporal-frontend-dev.example.com
    web_ui:
      enabled: true
      hostname: temporal-ui-dev.example.com
```

**Use Case:**

- Suitable for development and testing environments where plain text passwords are acceptable.

---

## Example 2b: Deployment with External PostgreSQL Database (Using Secret Reference - Recommended)

Deploy Temporal using an external PostgreSQL database for production with a password stored in a Kubernetes Secret.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTemporal
metadata:
  name: temporal-postgres-prod
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: temporal-prod
  create_namespace: true
  database:
    backend: postgresql
    external_database:
      host: postgres-db.example.com
      port: 5432
      username: temporaluser
      password:
        secretRef:
          name: temporal-db-credentials
          key: password
  ingress:
    frontend:
      enabled: true
      grpc_hostname: temporal-frontend-prod.example.com
    web_ui:
      enabled: true
      hostname: temporal-ui-prod.example.com
```

**Use Case:**

- Recommended for production environments where credentials should be stored securely in Kubernetes Secrets.
- Enables GitOps workflows where manifests can be safely committed to version control.
- Simplifies password rotation by only updating the Kubernetes Secret.

---

## Example 3: Advanced Observability with External Elasticsearch

Deploy Temporal with advanced observability features using an existing Elasticsearch cluster to power visibility
queries.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTemporal
metadata:
  name: temporal-elasticsearch
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: temporal-observability
  create_namespace: true
  database:
    backend: cassandra
  external_elasticsearch:
    host: elasticsearch.example.com
    port: 9200
    user: elasticuser
    password:
      stringValue: elasticpassword
  ingress:
    frontend:
      enabled: true
      grpc_hostname: temporal-frontend.example.com
    web_ui:
      enabled: true
      hostname: temporal-ui.example.com
```

**Use Case:**

- Recommended for scenarios demanding advanced search and visibility capabilities for workflows, typically in
  monitoring-intensive production deployments.

---

## Example 3b: External Elasticsearch with Secret Reference

Deploy Temporal with external Elasticsearch using credentials stored in Kubernetes Secrets.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTemporal
metadata:
  name: temporal-elasticsearch-prod
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: temporal-observability
  create_namespace: true
  database:
    backend: postgresql
    external_database:
      host: postgres-db.example.com
      port: 5432
      username: temporaluser
      password:
        secretRef:
          name: temporal-db-credentials
          key: password
  external_elasticsearch:
    host: elasticsearch.example.com
    port: 9200
    user: elasticuser
    password:
      secretRef:
        name: es-credentials
        key: password
  ingress:
    frontend:
      enabled: true
      grpc_hostname: temporal-frontend-prod.example.com
    web_ui:
      enabled: true
      hostname: temporal-ui-prod.example.com
```

**Use Case:**

- Production deployments with both external database and Elasticsearch using secure credential management.
- All sensitive values are stored in Kubernetes Secrets, enabling safe GitOps workflows.

---

## Example 4: Minimal Resource Deployment (Web UI Disabled)

Temporal deployment with minimal resource footprint by disabling the Temporal web UI component.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: TemporalKubernetes
metadata:
  name: temporal-minimal
spec:
  database:
    backend: mysql
    externalDatabase:
      host: mysql-db.example.com
      port: 3306
      user: temporaluser
      password: securepassword
  disableWebUi: true
  ingress:
    frontend:
      enabled: true
      hostname: temporal-frontend.example.com
```

**Use Case:**

- Optimal for environments where the UI is managed separately or unnecessary, reducing resource consumption.

---

## Example 5: Schema Control and Database Naming

Example demonstrating custom database naming conventions and controlled database schema initialization.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: TemporalKubernetes
metadata:
  name: temporal-custom-schema
spec:
  database:
    backend: postgresql
    externalDatabase:
      host: postgres-db.example.com
      port: 5432
      user: temporaluser
      password: securepassword
    databaseName: custom_temporal_db
    visibilityName: custom_temporal_visibility_db
    disableAutoSchemaSetup: true
  ingress:
    enabled: false
```

**Use Case:**

- Useful when integrating with existing database management procedures, enforcing naming conventions, or preventing
  automatic schema alterations.

---

## Namespace Management

### Example 6: Using Existing Namespace

Deploy Temporal into an existing namespace without creating it:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTemporal
metadata:
  name: temporal-existing-ns
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: existing-temporal-namespace
  create_namespace: false  # Use existing namespace
  database:
    backend: cassandra
  ingress:
    frontend:
      enabled: true
      grpc_hostname: temporal-frontend.example.com
```

**Use Case:**

- When namespace is pre-created with specific policies, quotas, or RBAC
- In environments with strict namespace governance
- When multiple components share the same namespace

### Example 7: Creating New Namespace (Default Behavior)

Create a new namespace for Temporal deployment:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTemporal
metadata:
  name: temporal-new-ns
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: temporal-namespace
  create_namespace: true  # Create new namespace (this is the default)
  database:
    backend: cassandra
  ingress:
    frontend:
      enabled: true
      grpc_hostname: temporal-frontend.example.com
```

**Note:** If `create_namespace` is omitted, it defaults to `true` for backward compatibility. The module will create
the namespace with appropriate labels and metadata.

---

## Advanced Configuration

### Example 8: Configuring Workflow History Limits

Configure dynamic runtime settings to increase workflow history limits for complex workflows with many activities or
large payloads.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTemporal
metadata:
  name: temporal-high-limits
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: temporal-prod
  create_namespace: true
  database:
    backend: postgresql
    external_database:
      host: postgres-db.example.com
      port: 5432
      username: temporaluser
      password:
        secretRef:
          name: temporal-db-credentials
          key: password
  # Increase history limits for complex workflows
  dynamic_config:
    # 100 MB (default is 50 MB)
    history_size_limit_error: 104857600
    # 100K events (default is ~51K)
    history_count_limit_error: 102400
    # Warning at 50 MB
    history_size_limit_warn: 52428800
    # Warning at 50K events
    history_count_limit_warn: 51200
  ingress:
    frontend:
      enabled: true
      grpc_hostname: temporal-frontend.example.com
```

**Use Case:**

- Workflows that orchestrate many activities or child workflows and exceed default history limits
- Long-running workflows that accumulate many events over time
- Workflows with large payloads that increase history size

**Note:** Consider using the `ContinueAsNew` pattern as an alternative to increasing limits indefinitely.

---

### Example 8b: Configuring Blob Size Limits for Large Payloads

Configure blob size limits to allow larger individual payloads in markers, signals, and activity inputs/outputs.
This is different from history limits which control total workflow history size.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTemporal
metadata:
  name: temporal-large-payloads
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: temporal-prod
  create_namespace: true
  database:
    backend: postgresql
    external_database:
      host: postgres-db.example.com
      port: 5432
      username: temporaluser
      password:
        secretRef:
          name: temporal-db-credentials
          key: password
  dynamic_config:
    # Blob size limits - for individual payloads (markers, signals, activity I/O)
    # 10 MB (default is 2 MB) - useful for workflows sending large IaC diffs
    blob_size_limit_error: 10485760
    # Warning at 5 MB
    blob_size_limit_warn: 5242880
  ingress:
    frontend:
      enabled: true
      grpc_hostname: temporal-frontend.example.com
```

**Use Case:**

- IaC runners that send large Pulumi/Terraform diffs as workflow signals or markers
- Workflows that process large documents or data payloads
- Systems where side effects need to store substantial return values

**Note:** If you hit `RecordMarkerCommandAttributes.Details exceeds size limit` errors, increase `blob_size_limit_error`.
This is different from history size limits - blob limits control individual payload sizes, not total history.

---

### Example 9: Configuring History Shards for Scalability

Set the number of history shards at deployment time. This is an **immutable** setting that determines cluster
parallelism and cannot be changed after initial deployment.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTemporal
metadata:
  name: temporal-high-scale
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: temporal-prod
  create_namespace: true
  database:
    backend: cassandra
  # Higher shard count for better parallelism and throughput
  # WARNING: This cannot be changed after deployment!
  num_history_shards: 1024
  ingress:
    frontend:
      enabled: true
      grpc_hostname: temporal-frontend.example.com
    web_ui:
      enabled: true
      hostname: temporal-ui.example.com
```

**Use Case:**

- High-throughput production environments expecting significant workflow volume
- Deployments where horizontal scalability is a priority
- Large-scale enterprise deployments

**Warning:** Choose this value carefully. The default (512) is safe for most production workloads. Higher values (1024,
2048) enable more parallelism but require more resources. This setting cannot be changed after initial deployment.

---

### Example 10: Per-Service Resource Configuration

Configure replicas and resources for each Temporal service independently. This allows fine-tuning based on workload
characteristics.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTemporal
metadata:
  name: temporal-tuned
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: temporal-prod
  create_namespace: true
  database:
    backend: postgresql
    external_database:
      host: postgres-db.example.com
      port: 5432
      username: temporaluser
      password:
        secretRef:
          name: temporal-db-credentials
          key: password
  services:
    # Frontend: API gateway - moderate resources, 2 replicas for HA
    frontend:
      replicas: 2
      resources:
        limits:
          cpu: "1000m"
          memory: "2Gi"
        requests:
          cpu: "200m"
          memory: "512Mi"
    # History: Most resource-intensive - higher resources, 3 replicas
    history:
      replicas: 3
      resources:
        limits:
          cpu: "2000m"
          memory: "4Gi"
        requests:
          cpu: "500m"
          memory: "1Gi"
    # Matching: Task dispatch - moderate resources
    matching:
      replicas: 2
      resources:
        limits:
          cpu: "1000m"
          memory: "2Gi"
        requests:
          cpu: "200m"
          memory: "512Mi"
    # Worker: Internal workflows - lighter resources
    worker:
      replicas: 1
      resources:
        limits:
          cpu: "500m"
          memory: "1Gi"
        requests:
          cpu: "100m"
          memory: "256Mi"
  ingress:
    frontend:
      enabled: true
      grpc_hostname: temporal-frontend.example.com
    web_ui:
      enabled: true
      hostname: temporal-ui.example.com
```

**Use Case:**

- Production deployments requiring fine-tuned resource allocation
- Environments with specific resource quotas or constraints
- Optimizing costs by right-sizing each service based on actual usage

**Note:** The history service is typically the most resource-intensive as it manages workflow state. Allocate more
resources to history for high-throughput deployments.

---

### Example 11: Complete Production Configuration

A comprehensive production deployment combining all advanced features.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesTemporal
metadata:
  name: temporal-production
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: temporal-prod
  create_namespace: true
  database:
    backend: postgresql
    external_database:
      host: postgres-rds.example.com
      port: 5432
      username: temporal_prod
      password:
        secretRef:
          name: temporal-db-credentials
          key: password
    database_name: temporal_prod
    visibility_name: temporal_visibility_prod
  external_elasticsearch:
    host: elasticsearch.example.com
    port: 9200
    user: temporal
    password:
      secretRef:
        name: es-credentials
        key: password
  # Scalability: 1024 shards for high throughput
  num_history_shards: 1024
  # Relaxed history and blob size limits for complex workflows
  dynamic_config:
    # History limits - for total workflow history
    history_size_limit_error: 104857600   # 100 MB
    history_count_limit_error: 102400     # 100K events
    history_size_limit_warn: 52428800     # 50 MB
    history_count_limit_warn: 51200       # 50K events
    # Blob size limits - for individual payloads (markers, signals, activity I/O)
    blob_size_limit_error: 10485760       # 10 MB (for large IaC diffs)
    blob_size_limit_warn: 5242880         # 5 MB warning
  # Production-grade service configuration
  services:
    frontend:
      replicas: 3
      resources:
        limits:
          cpu: "2000m"
          memory: "4Gi"
        requests:
          cpu: "500m"
          memory: "1Gi"
    history:
      replicas: 5
      resources:
        limits:
          cpu: "4000m"
          memory: "8Gi"
        requests:
          cpu: "1000m"
          memory: "2Gi"
    matching:
      replicas: 3
      resources:
        limits:
          cpu: "2000m"
          memory: "4Gi"
        requests:
          cpu: "500m"
          memory: "1Gi"
    worker:
      replicas: 2
      resources:
        limits:
          cpu: "1000m"
          memory: "2Gi"
        requests:
          cpu: "200m"
          memory: "512Mi"
  enable_monitoring_stack: true
  ingress:
    frontend:
      enabled: true
      grpc_hostname: temporal-grpc.example.com
      http_hostname: temporal-http.example.com
    web_ui:
      enabled: true
      hostname: temporal-ui.example.com
```

**Use Case:**

- Enterprise production deployments requiring high availability and throughput
- Mission-critical workflow orchestration systems
- Deployments expecting significant scale and needing fine-grained control

