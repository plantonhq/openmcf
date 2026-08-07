# KubernetesOtelCollector Pulumi Module

Deploys one operator-managed OpenTelemetry Collector: the optional
namespace and the single `opentelemetry.io/v1beta1`
`OpenTelemetryCollector` custom resource — the collector workload
(Deployment/DaemonSet/StatefulSet per mode, or sidecar registration),
the `<name>-collector` Service with receiver-derived ports, the
headless and monitoring Services, and the rendered config ConfigMap
are all operator-created from it. The CR renders through
`apiextensions.CustomResource` with an untyped spec map assembled in
one place, in byte lockstep with the Terraform twin's
`collector_spec`.

Prerequisite at deploy time: a KubernetesOtelOperator on the cluster
(it watches every namespace).

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance
   labels when `create_namespace` is true
2. **OpenTelemetryCollector CR** — the one declaration the operator
   reconciles; everything else (workload, Services, ConfigMap) is
   operator-created from it

## Rendering Notes

- **The config document is parsed before anything applies** — the
  v1beta1 CR's `config` is a STRUCTURED object (not the v1alpha1
  string), so `config_yaml` unmarshals and embeds as an object; an
  unparseable document fails loudly before preview. The operator's
  admission webhook validates the collector semantics at apply.
- **The name budget fails loudly** — `Resources()` rejects
  `metadata.name` past 42 characters: the operator derives child names
  by suffixing (`-collector-monitoring` is the longest stable suffix
  at 21 characters) and Kubernetes caps names at 63.
- **BACKGROUND deletion, explicitly** — the OPERATOR owns the
  collector CR's cascade; its ownership references tear down the
  workload, Services and ConfigMap. Foreground propagation would block
  the delete on children the operator keeps reconciling, so the CR
  carries the `pulumi.com/deletionPropagationPolicy: background`
  annotation.
- **No wait on the CR, deliberately** — collector readiness depends on
  the operator (webhook admission, image injection, workload rollout),
  which is not part of applying the resource; the verifier owns
  readiness.
- **Unset optionals are omitted** — every key renders only when the
  spec declares it, so the operator's defaulting stays authoritative
  (mode defaults to deployment, an empty image gets the operator's
  default collector image, an empty service account gets an
  operator-created one). Env vars render sorted by name for
  deterministic CR bodies on both engines.
- **Scaling renders only where it means something** — `replicas` and
  `autoscaler` render in deployment/statefulset modes only; the
  middleware-defaulted `replicas: 1` in daemonset/sidecar modes is
  deliberately ignored.
- **One volumes entry renders both halves** — each spec `volumes`
  entry splits into the CR's `volumes` (pod volume source) and
  `volumeMounts` (container mount), so the two lists can never
  disagree.

## Usage

```shell
planton pulumi up --manifest e2e/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the collector runs in |
| `collector_name` | Name of the OpenTelemetryCollector resource (equals `metadata.name`) |
| `service` | Collector Service (`<name>-collector`) with receiver-derived ports — empty in sidecar mode |
| `otlp_grpc_endpoint` | In-cluster OTLP gRPC ingest endpoint (`<service>:4317`) — valid when the config declares the standard `otlp` receiver; empty in sidecar mode |
| `otlp_http_endpoint` | In-cluster OTLP HTTP ingest endpoint (`http://<service>:4318`) — same contract; empty in sidecar mode |
| `headless_service` | Headless Service (`<name>-collector-headless`) for per-pod addressing — empty in sidecar mode |
| `monitoring_service` | Monitoring Service (`<name>-collector-monitoring`, port 8888) — empty in sidecar mode |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: name-budget guard → namespace → collector CR →
  output exports
- `module/locals.go`: resource identity, labels, the effective mode
  and the workload-mode rule — kept in lockstep with the Terraform
  module's `locals.tf`
- `module/otel_collector_cr.go`: the untyped CR and `collectorSpecBody`
  (config parsing, scaling, env/envFrom, the volume split, extra
  ports, resources, scheduling)
- `module/namespace.go`: optional namespace creation
- `module/outputs.go`: the exported handles (empty Service-derived
  handles in sidecar mode)
- `module/vars.go`: apiVersion/kind and the 42-character name budget

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
