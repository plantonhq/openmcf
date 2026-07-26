# Kubernetes Argo Workflows — design notes

## Grain

One resource = one Argo Workflows engine from the official
`argo-workflows` chart (argoproj index). The release is named after
`metadata.name` and `fullnameOverride` pins every child name
(`<name>-workflow-controller`, `<name>-server`), so the exported
outputs are deterministic. Names are capped at 43 characters (63 minus
the longest component suffix); both engines fail loudly instead of
letting the chart truncate silently. Several engines coexist per
cluster via `controller.instance_id` — an instanced controller
reconciles only Workflows carrying its
`workflows.argoproj.io/controller-instanceid` label.

## The composition seam

- **In:** Workflows/WorkflowTemplates/CronWorkflows as plain CRs
  (`KubernetesManifest`, charts, the CLI/UI) once the engine runs.
- **Out:** `server_service` / `server_kube_endpoint` for composed
  exposure, and `workflow_service_account` — the identity to annotate
  for IRSA/workload identity and to grant pipeline RBAC to.
- **Durability by reference:** the artifact endpoint FK-composes a
  `KubernetesSeaweedFs` (whose credentials Secret already carries the
  `accesskey`/`secretkey` pair the chart's selectors expect); the
  archive host FK-composes a `KubernetesPostgres`.

## The runner-identity placement

The chart creates the runner ServiceAccount + Role/RoleBinding in every
`workflow_namespaces` entry PLUS the install namespace. The chart's own
default for that list is `["default"]` — both modules always render the
list (declared entries or `[]`) so a default install never leaks runner
RBAC into the cluster's `default` namespace. The controller's WATCH is
cluster-wide either way; the list places identity, not scope. The chart
also defaults `workflow.serviceAccount.create` to FALSE — an engine
without its runner identity rejects every submission — so the module
always creates it.

## Structured chart blocks behind plain spec fields

`controller.instanceID` is a structured chart value ({enabled,
explicitID} — templates read `.enabled` directly); the spec's plain
`instance_id` string maps onto that shape. The archive engine section
is keyed by name (`postgresql`/`mysql`) with secret SELECTORS the
controller resolves at runtime — credentials never render as values.

## CRD posture

The default arm (`crds.full_schema: true`) installs full-schema CRDs
via a pre-install hook Job that DOWNLOADS them from the chart's GitHub
release (internet-at-install; `crds.base_url` mirrors it for restricted
networks); `false` falls back to chart-templated minified CRDs.
`crds.keep` defaults true either way.

## Cross-engine parity

Both engines render byte-identical chart values. The images are the
SPLIT registry+repository form: `images.tag` is shared, the registry
override moves each component (controller, server/argocli, executor)
to the mirror while upstream repository paths stay — the split-image
discipline.

## Deliberate exclusions

`workflowDefaults` documents, executor tuning, SSO (the server's
`sso` block rides `helm_values`; declare the `sso` auth mode here to
open it), priority classes per component, and the deprecated
`workflow.namespace` — reachable through `helm_values`, never the
primary interface.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
