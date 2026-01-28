# Create using CLI

Create a yaml using the example shown below. After the yaml is created, use the below command to apply.

```shell
planton apply -f <yaml-path>
```

# Basic Example

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: todo-list-api
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: cronjobs
  create_namespace: true
  schedule: "0 0 * * *"
  image:
    repo: nginx
    tag: latest
  resources:
    requests:
      cpu: 100m
      memory: 100Mi
    limits:
      cpu: 2000m
      memory: 2Gi
```

# Example w/ Environment Variables

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: todo-list-api
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: cronjobs
  create_namespace: true
  schedule: "0 0 * * *"
  env:
    variables:
      DATABASE_NAME: todo
  image:
    repo: nginx
    tag: latest
  resources:
    requests:
      cpu: 100m
      memory: 100Mi
    limits:
      cpu: 2000m
      memory: 2Gi
```

# Example w/ Environment Secrets (Direct String Value)

Provide secrets as direct string values (suitable for development/testing):

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: todo-list-api
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: cronjobs
  create_namespace: true
  schedule: "0 0 * * *"
  env:
    secrets:
      DATABASE_PASSWORD:
        value: my-secret-password
    variables:
      DATABASE_NAME: todo
  image:
    repo: nginx
    tag: latest
  resources:
    requests:
      cpu: 100m
      memory: 100Mi
    limits:
      cpu: 2000m
      memory: 2Gi
```

# Example w/ Environment Secrets (Kubernetes Secret Reference)

Reference existing Kubernetes Secrets for production deployments:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: todo-list-api
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: cronjobs
  create_namespace: true
  schedule: "0 0 * * *"
  env:
    secrets:
      DATABASE_PASSWORD:
        secretRef:
          name: my-app-secrets       # Name of existing K8s Secret
          key: db-password           # Key within the Secret
    variables:
      DATABASE_NAME: todo
  image:
    repo: nginx
    tag: latest
  resources:
    requests:
      cpu: 100m
      memory: 100Mi
    limits:
      cpu: 2000m
      memory: 2Gi
```

# Example w/ Mixed Secret Types

Use both direct values and secret references together:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: todo-list-api
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: cronjobs
  create_namespace: true
  schedule: "0 0 * * *"
  env:
    secrets:
      # Dev secret - direct value
      DEBUG_TOKEN:
        value: debug-only-token
      # Production secret - external reference
      DATABASE_PASSWORD:
        secretRef:
          name: postgres-credentials
          key: password
    variables:
      DATABASE_NAME: todo
  image:
    repo: nginx
    tag: latest
  resources:
    requests:
      cpu: 100m
      memory: 100Mi
    limits:
      cpu: 2000m
      memory: 2Gi
```

# Example w/ ConfigMap and Volume Mounts

Deploy a database backup CronJob with a custom backup script mounted as a ConfigMap:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: db-backup
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: database
  create_namespace: true
  schedule: "0 2 * * *"
  configMaps:
    backup-script: |
      #!/bin/bash
      echo "Starting backup at $(date)"
      pg_dump -h $DB_HOST -U $DB_USER $DB_NAME > /backup/dump.sql
      gzip /backup/dump.sql
      echo "Backup completed"
  image:
    repo: postgres
    tag: "15"
  command: ["/bin/bash", "/scripts/backup.sh"]
  volumeMounts:
    - name: backup-script
      mountPath: /scripts/backup.sh
      configMap:
        name: backup-script
        key: backup-script
        path: backup.sh
        defaultMode: 493  # 0755 - executable
    - name: backup-data
      mountPath: /backup
      emptyDir:
        sizeLimit: 1Gi
  env:
    variables:
      DB_HOST: postgres.database.svc
      DB_NAME: myapp
    secrets:
      DB_USER:
        value: admin
  resources:
    requests:
      cpu: 100m
      memory: 256Mi
    limits:
      cpu: 500m
      memory: 512Mi
  concurrencyPolicy: Forbid
  successfulJobsHistoryLimit: 3
  failedJobsHistoryLimit: 1
```

# Example w/ Multiple ConfigMaps

Deploy a CronJob with multiple configuration files:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: data-processor
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: batch-jobs
  create_namespace: true
  schedule: "0 */6 * * *"
  configMaps:
    processor-script: |
      #!/bin/bash
      source /config/settings.env
      python /scripts/process.py --config /config/processor.yaml
    processor-config: |
      input_bucket: s3://data-input
      output_bucket: s3://data-output
      batch_size: 1000
      compression: gzip
    settings-env: |
      export AWS_REGION=us-west-2
      export LOG_LEVEL=info
  image:
    repo: python
    tag: "3.11-slim"
  command: ["/bin/bash", "/scripts/run.sh"]
  volumeMounts:
    - name: processor-script
      mountPath: /scripts/run.sh
      configMap:
        name: processor-script
        key: processor-script
        path: run.sh
        defaultMode: 493
    - name: processor-config
      mountPath: /config/processor.yaml
      configMap:
        name: processor-config
        key: processor-config
        path: processor.yaml
    - name: settings-env
      mountPath: /config/settings.env
      configMap:
        name: settings-env
        key: settings-env
        path: settings.env
    - name: temp-work
      mountPath: /tmp/work
      emptyDir:
        medium: Memory
        sizeLimit: 512Mi
  resources:
    requests:
      cpu: 500m
      memory: 1Gi
    limits:
      cpu: 2000m
      memory: 4Gi
```

# Example w/ Secret Volume Mount

Mount TLS certificates for secure connections:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: secure-sync
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: sync-jobs
  create_namespace: true
  schedule: "*/30 * * * *"
  image:
    repo: curlimages/curl
    tag: latest
  command: ["curl"]
  args:
    - "--cacert"
    - "/certs/ca.crt"
    - "--cert"
    - "/certs/tls.crt"
    - "--key"
    - "/certs/tls.key"
    - "https://secure-api.example.com/sync"
  volumeMounts:
    - name: tls-certs
      mountPath: /certs
      readOnly: true
      secret:
        name: api-tls-certs
  resources:
    requests:
      cpu: 50m
      memory: 64Mi
    limits:
      cpu: 200m
      memory: 128Mi
```

# Advanced Example w/ grpc api call

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: grpc-invoker
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: cronjobs
  create_namespace: true
  
  # The container image to use. In this case, a public image with grpcurl.
  image:
    repo: fullstorydev/grpcurl
    tag: latest

  # The schedule for the cron job. Runs at midnight (00:00) every day.
  schedule: "0 0 * * *"

  # How concurrency is handled. "Forbid" is the default if omitted.
  concurrencyPolicy: "Forbid"

  # Environment variables (optional). Adjust or remove as needed.
  env:
    variables:
      FOO: "BAR"
    secrets: {}

  # The commands/arguments for grpcurl. Calls "my.package.Service/Method" on "my-grpc-service:50051".
  command:
    - "grpcurl"
  args:
    - "-plaintext"
    - "my-grpc-service:50051"
    - "my.package.Service/Method"

  # Basic resource limits and requests. Adjust these to match your usage.
  resources:
    requests:
      cpu: "50m"
      memory: "100Mi"
    limits:
      cpu: "200m"
      memory: "256Mi"

  # Typically "Never" is recommended for CronJobs to avoid restarts. 
  restartPolicy: "Never"

  # Keep the last 3 successful runs and the last 1 failure in history.
  successfulJobsHistoryLimit: 3
  failedJobsHistoryLimit: 1

  # Number of retries if the job fails before considering it permanently failed.
  backoffLimit: 6

  # If you need to temporarily pause further runs, set this to true.
  # suspend: false

  # If you miss a schedule, how long (in seconds) to still try to start. 0 = no deadline enforced.
  # startingDeadlineSeconds: 0
```

# Namespace Management

The `create_namespace` field controls whether the namespace is automatically created or if an existing namespace should be used.

## Creating a New Namespace

Set `create_namespace: true` to automatically create the namespace:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: isolated-job
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: my-cronjobs
  create_namespace: true
  schedule: "0 2 * * *"
  image:
    repo: busybox
    tag: latest
  resources:
    requests:
      cpu: "100m"
      memory: "128Mi"
    limits:
      cpu: "500m"
      memory: "256Mi"
```

**Use this when:**
- You want dedicated namespace isolation for the CronJob
- The namespace doesn't exist yet
- You want the CronJob to manage its namespace lifecycle

## Using an Existing Namespace

Set `create_namespace: false` to deploy into a pre-existing namespace:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: shared-job
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: shared-batch-jobs
  create_namespace: false
  schedule: "0 4 * * *"
  image:
    repo: busybox
    tag: latest
  resources:
    requests:
      cpu: "100m"
      memory: "128Mi"
    limits:
      cpu: "500m"
      memory: "256Mi"
```

**Use this when:**
- Multiple CronJobs share the same namespace
- The namespace is managed separately (e.g., by a platform team)
- You're following a GitOps pattern where namespaces are pre-created
- **Important:** Ensure the namespace exists before deploying, or the deployment will fail

## Best Practices

- **Isolated workloads**: Use `create_namespace: true` for dedicated, single-purpose CronJobs
- **Shared namespaces**: Use `create_namespace: false` when multiple CronJobs need to coexist in the same namespace
- **GitOps workflows**: Set `create_namespace: false` and manage namespaces separately for better lifecycle control
- **Development environments**: Use `create_namespace: true` for quick iterations and easy cleanup
```
