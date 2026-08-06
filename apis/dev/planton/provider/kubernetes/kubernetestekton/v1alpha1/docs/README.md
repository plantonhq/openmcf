# Kubernetes Tekton — design notes

## Grain

One resource = the cluster's `TektonConfig` — the singleton declaration
the Tekton Operator reconciles into running components via
TektonInstallerSets. The operator's admission webhook enforces both the
singleton and its fixed name `config`; `metadata.name` keys the Planton
resource and the state identity only. The CR renders through the
untyped-CustomResource pattern (the CRD types its `options` override
blocks with preserve-unknown-fields, which typed codegen cannot carry)
in byte lockstep across both engines.

## The two-kind grain is the destroy fix

Deleting a TektonConfig makes the operator remove every component it
installed; the InstallerSet finalizers are processed by the RUNNING
operator. With the operator and the configuration as separate kinds,
teardown order is structural: this resource destroys first (both
modules block on the finalizer with a 15-minute delete budget), the
operator after. The single-kind shape this replaced deadlocked exactly
here — the operator died in the same destroy as its own finalizers'
processor.

## What the spec models

The operator API's own JSON keys, rendered only when declared so
upstream defaulting stays authoritative: profile + immutable
targetNamespace (+ namespace metadata), cluster-wide component
placement (`config`), the pipelines surface (the cluster-global
CloudEvents sink, api-fields gate, execution defaults, tri-state
feature flags, resolver toggles, metrics shape, the
replicas+buckets performance block — Tekton's real HA mechanism),
triggers, the read-only dashboard knob, Chains (disabled /
generateSigningSecret), and the pruner (keep XOR keep-since — the
webhook's own rule, mirrored as spec CEL). `additional_params` is the
operator-param escape surface.

## Deliberate exclusions

The per-component `options` override blocks (free-form manifest
patching — an operator power tool, not a configuration surface), Tekton
Hub/Results (not installed by any profile the catalog models),
OpenShift-only platform blocks, and the deprecated
`send-cloudevents-for-runs` flag (superseded upstream; the sink alone
governs delivery).

## Outputs

`namespace` (the resolved target namespace), `profile`, and — on
profile `all` — the dashboard Service handles for composed exposure
plus the port-forward one-liner.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
