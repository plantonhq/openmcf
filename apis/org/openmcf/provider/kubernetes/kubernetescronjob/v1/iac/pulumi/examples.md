# CronJobKubernetes API - Example Configurations

This document shows a series of examples for creating and configuring Kubernetes CronJobs using **CronJobKubernetes**.
Each example highlights different aspects such as scheduling, environment variables, secrets, concurrency policies, and
retry strategies.

> **Tip:** To apply a configuration, save the YAML content to a file (e.g., `cronjob.yaml`) and run:

> ```shell
> planton apply -f cronjob.yaml
> ```

---

## 1. Minimal Configuration

In this basic example, the CronJob runs a simple container image with default concurrency and retry settings, executing
once a day at midnight.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: daily-backup
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: my-namespace
  create_namespace: true
  schedule: "0 0 * * *"
  image:
    repo: busybox
    tag: latest
  resources:
    limits:
      cpu: "500m"
      memory: "256Mi"
    requests:
      cpu: "100m"
      memory: "128Mi"
  env: { }
```

- **schedule**: Defines when the job will run (cron format). Here, it runs daily at midnight.
- **image**: A minimal container image (`busybox`) for demonstration.
- **resources**: Basic CPU/memory requests and limits.

---

## 2. Setting Environment Variables

This example demonstrates how to pass environment variables to the CronJob container. Environment variables might be
used for passing non-sensitive configuration values, such as database hostnames or feature flags.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: weekly-report
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: my-namespace
  create_namespace: true
  schedule: "0 9 * * 1"
  image:
    repo: org/reports-generator
    tag: v1.0.0
  env:
    variables:
      REPORT_TYPE: "weekly"
      S3_BUCKET_NAME: "my-report-bucket"
  resources:
    limits:
      cpu: "1"
      memory: "1Gi"
    requests:
      cpu: "100m"
      memory: "256Mi"
```

- **schedule**: Runs every Monday at 9:00 AM.
- **env.variables**: Adds two environment variables, `REPORT_TYPE` and `S3_BUCKET_NAME`.

---

## 3. Using Secret Environment Variables

Sensitive data should be stored in secrets rather than environment variables. This example shows two approaches for providing secrets to the CronJob container.

### Option 1: Direct String Value (Development/Testing)

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: db-maintenance
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: my-namespace
  create_namespace: true
  schedule: "0 3 * * *"
  image:
    repo: org/maintenance-tool
    tag: stable
  env:
    secrets:
      DB_PASSWORD:
        value: my-secret-password
    variables:
      DB_HOST: db-server.prod.svc.cluster.local
      DB_NAME: myappdb
  resources:
    limits:
      cpu: "1"
      memory: "512Mi"
    requests:
      cpu: "100m"
      memory: "128Mi"
```

### Option 2: Kubernetes Secret Reference (Production)

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: db-maintenance
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: my-namespace
  create_namespace: true
  schedule: "0 3 * * *"
  image:
    repo: org/maintenance-tool
    tag: stable
  env:
    secrets:
      DB_PASSWORD:
        secretRef:
          name: postgres-credentials   # Name of existing K8s Secret
          key: password                # Key within the Secret
    variables:
      DB_HOST: db-server.prod.svc.cluster.local
      DB_NAME: myappdb
  resources:
    limits:
      cpu: "1"
      memory: "512Mi"
    requests:
      cpu: "100m"
      memory: "128Mi"
```

- **env.secrets.value**: Direct string value for development/testing scenarios.
- **env.secrets.secretRef**: Reference to an existing Kubernetes Secret (recommended for production).
- **env.variables**: Non-sensitive parameters like `DB_HOST` and `DB_NAME` remain in `variables`.

---

## 4. Tuning Concurrency & Retry Behavior

You can customize concurrency policies and retry logic to control how CronJobs run under load or when they fail.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: heavy-lift-job
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: my-namespace
  create_namespace: true
  schedule: "*/30 * * * *"
  concurrencyPolicy: "Forbid"
  backoffLimit: 3
  successfulJobsHistoryLimit: 2
  failedJobsHistoryLimit: 2
  image:
    repo: org/heavy-lift
    tag: v2.1
  resources:
    limits:
      cpu: "2"
      memory: "2Gi"
    requests:
      cpu: "200m"
      memory: "256Mi"
  env: { }
```

- **schedule**: Runs every 30 minutes.
- **concurrencyPolicy**: `"Forbid"` means a new run cannot start if a previous run is still in progress.
- **backoffLimit**: The job can retry up to three times on failure.
- **successfulJobsHistoryLimit** / **failedJobsHistoryLimit**: Retain logs for the 2 most recent successes or failures.

---

## 5. Advanced Scheduling & Optional Settings

This example highlights some additional parameters like `suspend`, `startingDeadlineSeconds`, and a more complex cron
schedule.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: CronJobKubernetes
metadata:
  name: end-of-month-report
spec:
  target_cluster:
    cluster_name: "my-gke-cluster"
  namespace:
    value: my-namespace
  create_namespace: true
  # Run at 23:30 on the last day of each month
  schedule: "30 23 28-31 * *"
  # This ensures if a schedule is missed for over 2 hours, it won't run retroactively
  startingDeadlineSeconds: 7200
  # Temporarily suspend the schedule
  suspend: false
  concurrencyPolicy: "Replace"
  backoffLimit: 1
  successfulJobsHistoryLimit: 5
  failedJobsHistoryLimit: 5
  image:
    repo: org/report-service
    tag: monthly-latest
  resources:
    limits:
      cpu: "2"
      memory: "2Gi"
    requests:
      cpu: "500m"
      memory: "512Mi"
  env:
    variables:
      REPORT_TZ: "UTC"
      SEND_EMAIL_NOTIFICATIONS: "true"
```

- **schedule**: Runs at 23:30 on day 28, 29, 30, or 31 (last valid day of each month).
- **startingDeadlineSeconds**: If a job can’t start within 2 hours of its scheduled time, it’s skipped.
- **suspend**: Set to `true` to stop scheduling new runs; set back to `false` to resume.
- **concurrencyPolicy**: `"Replace"` cancels the currently running job and starts a new run if it’s time again.

---

## 6. Namespace Management

The `create_namespace` field controls whether the CronJob creates a new namespace or uses an existing one.

### Creating a New Namespace

Set `create_namespace: true` to automatically create the namespace with the CronJob:

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
    limits:
      cpu: "500m"
      memory: "256Mi"
    requests:
      cpu: "100m"
      memory: "128Mi"
  env: {}
```

- **create_namespace: true**: The Pulumi module will create the namespace `my-cronjobs` if it doesn't exist.
- **Use case**: Ideal when you want dedicated namespace isolation for this CronJob.

### Using an Existing Namespace

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
    limits:
      cpu: "500m"
      memory: "256Mi"
    requests:
      cpu: "100m"
      memory: "128Mi"
  env: {}
```

- **create_namespace: false**: The Pulumi module will reference the existing `shared-batch-jobs` namespace.
- **Use case**: Ideal for multi-tenant scenarios where multiple CronJobs share a namespace.
- **Important**: Ensure the namespace exists before deploying, or the deployment will fail.

### Best Practices

- **Isolated workloads**: Use `create_namespace: true` for dedicated, single-purpose CronJobs.
- **Shared namespaces**: Use `create_namespace: false` when multiple CronJobs need to coexist in the same namespace.
- **GitOps workflows**: Set `create_namespace: false` and manage namespaces separately for better lifecycle control.
- **Development environments**: Use `create_namespace: true` for quick iterations and easy cleanup.

---

## Summary

- **Cron Schedules**: Use the `schedule` field for standard cron expressions like `"0 0 * * *"` or more advanced
  patterns.
- **Resource & Env Config**: The `resources` field controls CPU/memory requests/limits, while `env` can hold both
  non-sensitive `variables` and secret-based `secrets`.
- **Concurrency & Retry**: Tweak `concurrencyPolicy`, `backoffLimit`, and the job history limits for advanced behavior.
- **Suspend & Deadline**: `suspend` can pause job creation; `startingDeadlineSeconds` avoids flooding the system with
  missed jobs.

**Next Steps**

1. Write your CronJobKubernetes YAML file using the examples above as a guide.
2. Run `planton apply -f <your-cronjob-file.yaml>` to apply the resource on your cluster.
3. Monitor logs and verify that your tasks run according to schedule.

For detailed reference, consult the [CronJobKubernetes API documentation](#) (link placeholder). If you encounter any
issues or need advanced features (e.g., custom commands/args), feel free to open a ticket with our support team.
