# GcpDataprocAutoscalingPolicy — Deep Dive

## The problem this resource solves

Hand-sizing a Dataproc cluster forces a bad choice: provision for peak (idle cost) or provision for average (queued jobs at peak). Dataproc's answer is metric-driven autoscaling — but unlike inline autoscalers (AlloyDB, Bigtable), Dataproc factors the scaling behavior out into a **separate, shareable policy resource**. The cluster only holds a reference.

That separation is the operational win: a platform team encodes its scaling philosophy once ("modest on-demand base, aggressive spot burst, hour-long drain") and every batch cluster in the region attaches the same policy. Re-tuning the policy re-tunes the whole fleet in place — no cluster is touched, no cluster is recreated.

## Where it sits in the composition

- **GcpDataprocAutoscalingPolicy** — this resource: bounds, weights, and the YARN algorithm. Regional.
- **GcpDataprocCluster** — attaches the policy through `clusterConfig.autoscalingPolicyUri`, a `StringValueOrRef` whose default reference resolves this policy's `status.outputs.name` (the full resource name `projects/{p}/locations/{l}/autoscalingPolicies/{id}`). The attachment is an in-place cluster update; so are swap and detach.

A cluster can only attach a policy in its own region, so a multi-region estate carries one policy per region (typically identical manifests differing only in `location`).

## Lifecycle contract

| Property | Behavior |
|---|---|
| `policyId`, `location`, `projectId` | Immutable (ForceNew) |
| Everything else (bounds, weights, algorithm) | Mutable in place — every attached cluster re-tunes on the next evaluation |
| Deletion | The API **refuses to delete a policy while any cluster references it**; detach or destroy the clusters first |

Teardown ordering follows: clusters first, policy last. In dependency-graph terms the policy behaves like a network or a KMS key — long-lived shared infrastructure beneath ephemeral consumers.

## How the autoscaler works

Every `cooldown_period` (default `120s`; bounds 2 minutes to 1 day), the autoscaler samples the cluster's YARN memory metrics:

- **Pending memory** — demand YARN cannot place. Scale-up signal.
- **Available memory** — idle capacity. Scale-down signal.

It computes the worker delta those metrics suggest, then applies the configured fraction of it:

```
delta_up   = pending_memory_workers   × scale_up_factor
delta_down = available_memory_workers × scale_down_factor
```

### Factor semantics

`scale_up_factor` / `scale_down_factor` (both required, 0.0-1.0) are the aggressiveness dials:

- `1.0` — act on the full recommendation every evaluation. Fastest reaction; fits retryable batch.
- `0.5` — halve each step. Balanced default.
- `0.05` — creep in ~5% steps. Smooth but slow; fits latency-tolerant steady workloads.

The `scale_up_min_worker_fraction` / `scale_down_min_worker_fraction` filters (0.0-1.0, default 0.0) discard recommendations smaller than the given fraction of current cluster size — the anti-noise lever. At `0.1`, a 20-worker cluster ignores any recommendation to change by fewer than 2 workers.

### The drain window

`graceful_decommission_timeout` (required; `^[0-9]+s$`; 0s to 1 day) is how long a scale-down waits for a worker's running tasks before forcing removal. Short windows (`"300s"`) suit retryable batch; long windows (`"7200s"`) protect long-task production jobs. This is the policy's own drain setting for autoscaler-initiated removals — distinct from the cluster-level `gracefulDecommissionTimeout`, which applies to manual worker-count changes.

## Primary vs. secondary: weight semantics

The policy bounds two groups independently and splits new capacity between them by weight:

- **Primary workers** (`worker_config`) carry HDFS DataNodes — the stable base. `max_instances` is required (>= 1). `min_instances` is 0-or->=2: the Dataproc API requires at least 2 primary workers on an autoscaled cluster, and `0` simply accepts that default.
- **Secondary workers** (`secondary_worker_config`, optional) carry no HDFS data, so the autoscaler can add and remove them aggressively. All bounds default to 0; `min_instances: 0` gives a scale-to-zero burst arm.

Weights are proportional shares: primary weight 1 + secondary weight 3 sends ~75% of new capacity to the secondary group. The canonical cost-optimized shape is a small fixed-ish primary group (weight 1) with a large, spot-priced, scale-to-zero secondary group (weight 3-4) — pair it with `preemptibility: SPOT` and an `instanceFlexibilityPolicy` on the cluster side.

## Tuning archetypes

| Archetype | Factors | Weights (pri:sec) | Drain | Filter fractions |
|---|---|---|---|---|
| Balanced default | 0.5 / 0.5 | 1:1 | `1800s` | — |
| Aggressive batch | 1.0 / 1.0 | 1:4, secondary scale-to-zero | `300s` | — |
| Conservative production | 0.2 / 0.1 | 2:1 | `7200s` | 0.1 / 0.1 |

These are shipped as the three presets.

## Failure modes an operator will actually meet

- **Policy delete fails with "in use"** — a cluster still references it. Detach (`autoscalingPolicyUri` removal is an in-place cluster update) or destroy the clusters first.
- **Cluster attach fails across regions** — the policy must live in the cluster's region.
- **Autoscaled cluster stuck at 2 primaries** — `min_instances: 0` accepts the API's default floor of 2; set an explicit >= 2 floor to raise it. A floor of 1 is rejected pre-deploy.
- **Cluster never scales down** — `scale_down_factor` near 0 disables shrink in practice; use the cluster's `lifecycleConfig.idleDeleteTtl` as the cost backstop.
- **Scaling feels twitchy** — raise the min-worker fractions and/or lengthen `cooldown_period` before shrinking the factors.

## 90/10 Coverage

The spec models the complete `google_dataproc_autoscaling_policy` surface: identity (`policy_id`, `location`, ambient `project_id`), both worker-group bounds with weights, and the basic algorithm (cooldown + the five YARN knobs).

### Recorded skips (with reasons)

- **`deletion_policy`** — a client-side lever that conflicts with Planton-managed destroy (catalog-wide skip).

Nothing else on the resource is skipped — this kind reaches full API coverage.

## Implementation landscape

Both modules enable the Dataproc API with `disable_on_destroy=false` (tearing down one policy never disables the API project-wide) and export the same three outputs. Zero-valued optional fields (weights, floors, fractions, cooldown) are withheld from the API so GCP's own defaults apply, identically on both engines. The Dataproc autoscaling-policy API has no labels surface, so no platform attribution labels are stamped — one of the few catalog kinds where that is true.

- **Terraform**: `iac/tf/` — `google_dataproc_autoscaling_policy` on the plain `google` provider (`~> 6.0`); every modeled field is GA on the released line.
- **Pulumi**: `iac/pulumi/module/` — `dataproc.AutoscalingPolicy` (pulumi-gcp v9).

Outputs: `name` (the full resource name — the cluster-side composition handle), `policy_id`, `location` (echoed from the spec so verifiers can address the policy without parsing paths).

## References

- [Autoscaling clusters](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/autoscaling)
- [AutoscalingPolicies REST API](https://cloud.google.com/dataproc/docs/reference/rest/v1/projects.locations.autoscalingPolicies)
- [Terraform google_dataproc_autoscaling_policy](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/dataproc_autoscaling_policy)
- [Pulumi dataproc.AutoscalingPolicy](https://www.pulumi.com/registry/packages/gcp/api-docs/dataproc/autoscalingpolicy/)
