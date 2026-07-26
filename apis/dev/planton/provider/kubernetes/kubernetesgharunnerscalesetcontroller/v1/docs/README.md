# Kubernetes GHA Runner Scale Set Controller — design notes

## Grain

One resource = one controller install from the official
`gha-runner-scale-set-controller` OCI chart
(ghcr.io/actions/actions-runner-controller-charts; chart and controller
image move in lockstep). The release is named after `metadata.name` and
`fullnameOverride` pins the chart fullname — and with it the created
ServiceAccount's name — to the same value, so `service_account_name`
exports deterministically (the handle fenced scale sets reference).

## The chart surface, typed

Replicas (leader election turns on automatically above one — hot
standby, not sharding), the combined-form image override (the chart
takes `image.repository` as the full mirror path — verified in the
deployment template at the pin), resources, the flags block (log
level/format, the single-namespace watch fence, reconcile concurrency,
`immediate`/`eventual` update strategy, label-propagation exclusions,
API-client rate limits, the workqueue rate limiter — a STRUCTURED chart
block `{name}` behind the spec's plain string — and the health-probe
address that also adds the probes), metrics (declaring the block
ENABLES it; the three addresses wire the controller AND every listener
pod), pull secrets (also handed to listener pods), scheduling, and the
`helm_values` escape hatch.

## The multi-tenancy seam

`flags.watch_single_namespace` fences the controller; scale sets
outside auto-discovery then bind explicitly via
`controller_service_account` on the scale set kind, wired from this
kind's `service_account_name` output. One cluster-wide controller
remains the documented default.

## CRD posture

The chart installs the four `actions.github.com` CRDs release-owned:
they DELETE with the controller, cascade-deleting every runner scale
set. The spec's chart_version comment carries the GitHub
matching-versions support rule; the destroy warning lives on the
README and the verifier asserts the posture both ways.

## Deliberate exclusions

Env entries, volumes, pprof, security contexts and per-resource
metadata — reachable through `helm_values`, never the primary
interface. The legacy community ARC chart (RunnerDeployment CRs,
cert-manager webhooks) is a different product line and not a design
source.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
