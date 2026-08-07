# GCP Cloud Composer User Workloads Secret — Deep Dive

## The problem this resource solves

Airflow DAGs need credentials — database passwords, API tokens, connection URIs. The anti-patterns are well known: hard-coding them in DAG files (they land in version control), or stuffing them into environment variables (they leak into logs and are visible environment-wide). Composer's user workloads Secrets deliver credentials as a Kubernetes Secret in the environment's GKE cluster: DAGs consume them by name, the material lives in exactly one governed place, and rotating a credential is an in-place data update.

## The base64 contract

Values follow Kubernetes Secret semantics: every entry MUST be base64-encoded, and the API rejects raw strings (validated at manifest time, before any deploy):

```shell
echo -n 'postgresql://appuser:s3cr3t@10.20.0.12:5432/orders' | base64
```

The `-n` matters — a trailing newline becomes part of the decoded secret and produces maddening authentication failures.

## Where the secret material lives (and doesn't)

| Surface | Exposure |
|---|---|
| Manifest `spec.data` | Base64-encoded (obscured, not encrypted — treat manifests as sensitive) |
| Stack outputs | Never — only `name` and `secret_name` are exported |
| Terraform plan/state | Attribute marked sensitive; plans redact it; state is the engine's secret boundary |
| Pulumi state | Whole map wrapped with `ToSecret` |
| The environment's GKE cluster | A regular Kubernetes Secret, consumed by workloads |

## How DAGs consume it

- **KubernetesPodOperator** — mount the Secret by `secret_name` as environment variables or files in the task pod.
- **Airflow connections** — store a connection URI as an entry and point Airflow's secrets backend at it, keeping connections out of the Airflow UI/database.

## Mutability profile

| Surface | Mutability |
|---|---|
| `secret_name`, `environment`, `region`, `project_id` | Immutable |
| `data` | Mutable — rotation is an in-place update |

Deleting this resource deletes the Kubernetes Secret from the environment.

## Composition

The environment is a first-class resource this one composes against by reference:

```yaml
environment:
  valueFrom:
    kind: GcpCloudComposerEnvironment
    name: prod-airflow
    fieldPath: status.outputs.environment_name
```

No API enablement happens here — the Composer API is enabled by the environment the Secret is delivered into (a Secret cannot exist without one). For non-secret configuration (feature flags, endpoints, tuning), use the sibling `GcpCloudComposerUserWorkloadsConfigMap` instead.

## 90/10 coverage

The resource is small and the spec models all of it: identity (`project_id`, `region`), the target `environment`, the `secret_name`, and the base64 `data` map (min one entry, sensitive). Outputs are the fully qualified `name` and the `secret_name` DAGs reference.

## Recorded skips (with reasons)

| Skipped | Reason |
|---------|--------|
| **`deletion_policy`** | Client-side Terraform lever conflicting with Planton-managed destroy (catalog-wide decision). This is the only skip — the API surface is otherwise fully modeled. |
