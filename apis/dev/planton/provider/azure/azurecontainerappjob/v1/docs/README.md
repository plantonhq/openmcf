# AzureContainerAppJob -- Design Research

## The Resource

A Container App Job (`Microsoft.App/jobs`) is the run-to-completion
sibling of the Container App: each execution runs the job's container
template to completion (bounded by a hard replica timeout), driven by
exactly one of three trigger models. The component maps onto
`azurerm_container_app_job` (azurerm v4.x,
`internal/services/containerapps/container_app_job_resource.go` +
`helpers/container_app_job.go`), parity-verified against pulumi-azure v6
(`containerapp.Job`).

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `job_name` | The provider's name contract (lowercase alphanumerics/hyphens, no dots -- unlike app names, <=32) mirrored as a CEL. ForceNew |
| `location` | `region` | Jobs carry their own location (unlike apps, which inherit the environment's). ForceNew |
| `resource_group_name` | `resource_group` | FK to `AzureResourceGroup`. ForceNew |
| `container_app_environment_id` | `container_app_environment_id` | Single parent ARM-id FK. ForceNew |
| `replica_timeout_in_seconds` | `replica_timeout_in_seconds` | Required; the hard per-replica deadline -- a replica killed by it counts as failed |
| `replica_retry_limit` | `replica_retry_limit` | Optional; retries per failed replica before the execution is marked failed |
| `workload_profile_name` | `workload_profile_name` | Optional; omitted runs on Consumption |
| `template.container/init_container/volume` | `containers`/`init_containers`/`volumes` | The Container App container model minus continuous-serving concerns; the job template has NO scale block -- executions are the scaling unit |
| `manual_trigger_config` / `schedule_trigger_config` / `event_trigger_config` | `manual_trigger` / `schedule_trigger` / `event_trigger` | The provider's `ExactlyOneOf` contract front-loaded as a message CEL; all trigger blocks ForceNew. Each carries parallelism + replica_completion_count (defaults 1/1) |
| `event_trigger_config.scale` | `event_trigger.scale` | max/min executions (100/0), polling interval (30), and KEDA rules -- the rules reuse the provider's custom-scale-rule schema, so the same KEDA scaler allowlist is mirrored as a CEL, with optional `identity_id` for workload-identity scalers |
| `secret` | `secrets[]` | The app-level secret schema: plain value XOR Key Vault reference + identity pairing (CELs) |
| `registry` | `registries[]` | Identity XOR username+password_secret_name (CEL) |
| `identity` | `identity` | SystemAssigned / UserAssigned / both; ids-match-type CEL |
| `tags` | `tags` | User tags merged over metadata-derived tags |

Computed: `outbound_ip_addresses`, `event_stream_endpoint` -- both
exported as stack outputs alongside `job_id`/`job_name` and
`identity_principal_id`.

## Decomposition Decision

First-class kind beside `AzureContainerApp`, never a mode on it: ARM
models jobs as a separate resource type with a disjoint lifecycle
surface (triggers and executions vs revisions and ingress), and the
provider models them as separate resources. A single kind with a
"job mode" would bury two different operational models in one spec.

## Front-Loaded Contracts

- **Exactly one trigger** (message CEL) -- mirrors the provider's
  `ExactlyOneOf` across the three trigger blocks.
- **Per-probe-type contracts** on containers -- success threshold is
  readiness-only; failure ceilings 30/48/240 (liveness/readiness/
  startup); the same container-level CELs as the app kind.
- **KEDA scaler vocabulary** on event scale rules -- the provider
  validates `custom_rule_type` against its pinned allowlist (the job's
  rules reuse the app's custom-scale-rule schema); mirrored as a CEL.
- **Volume storage_name pairing** -- AZURE_FILE / NFS_AZURE_FILE volumes
  require the environment-storage registration name; EMPTY_DIR / SECRET
  must omit it.

## Recorded Skips (with reasons)

Nothing skipped: the spec models the provider's full job surface.

## Operational Behavior Worth Knowing

- **The replica timeout is a hard kill.** Size
  `replica_timeout_in_seconds` to the slowest legitimate execution;
  Azure terminates and fails anything that exceeds it.
- **Event-triggered jobs poll.** KEDA evaluates the scaler every
  `polling_interval_in_seconds`; queue latency is bounded by the poll,
  not instantaneous.
- **Manual executions start via the ARM jobs API** (`az containerapp
  job start`, SDKs, or pipelines) using the exported `job_id`.
