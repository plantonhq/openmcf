# GCP Cloud Composer User Workloads ConfigMap — Deep Dive

## The problem this resource solves

Airflow DAGs accumulate configuration: which endpoint to call, how large a batch to process, whether the new code path is on. Hard-coding these in DAG files means every tuning change is a code change, a review, and a redeploy. Composer's user workloads ConfigMaps put that configuration in a named Kubernetes ConfigMap in the environment's GKE cluster: DAGs read values by key, and a config change is an in-place apply.

## ConfigMap vs Secret

This kind carries **plain** data — readable by anything with access to the environment's cluster. The decision rule:

| Material | Kind |
|---|---|
| Feature flags, endpoints, batch sizes, timezones, tuning | `GcpCloudComposerUserWorkloadsConfigMap` (this kind) |
| Passwords, API tokens, connection URIs, anything credential-shaped | `GcpCloudComposerUserWorkloadsSecret` |

The two are deliberate siblings with the same shape (region, environment reference, name, data map); the Secret adds base64-encoded values and sensitivity handling.

## Values are strings

Kubernetes ConfigMap values are strings, so YAML scalars that parse as other types must be quoted: `"500"`, `"true"`, and YAML 1.1's `on`/`off`/`yes`/`no` traps. An unquoted `true` arrives as a boolean and fails the string map.

## How DAGs consume it

- **KubernetesPodOperator** — mount the ConfigMap by `config_map_name` as environment variables or files in the task pod.
- **DAG code** — read via the Kubernetes API or mounted files where the task runs in-cluster.

## Mutability profile

| Surface | Mutability |
|---|---|
| `config_map_name`, `environment`, `region`, `project_id` | Immutable |
| `data` | Mutable — tuning changes apply in place |

Deleting this resource deletes the Kubernetes ConfigMap from the environment.

## Composition

The environment is a first-class resource this one composes against by reference:

```yaml
environment:
  valueFrom:
    kind: GcpCloudComposerEnvironment
    name: prod-airflow
    fieldPath: status.outputs.environment_name
```

No API enablement happens here — the Composer API is enabled by the environment the ConfigMap is delivered into (a ConfigMap cannot exist without one).

## 90/10 coverage

The resource is small and the spec models all of it: identity (`project_id`, `region`), the target `environment`, the `config_map_name`, and the plain `data` map (min one entry). Outputs are the fully qualified `name` and the `config_map_name` DAGs reference.

## Recorded skips (with reasons)

| Skipped | Reason |
|---------|--------|
| **`deletion_policy`** | Client-side Terraform lever conflicting with Planton-managed destroy (catalog-wide decision). This is the only skip — the API surface is otherwise fully modeled. |
