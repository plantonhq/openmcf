# Kubernetes GHA Runner Scale Set — design notes

## Grain

One resource = one AutoscalingRunnerSet from the official
`gha-runner-scale-set` OCI chart — one GitHub registration
(repo/org/enterprise) per resource, many per cluster. The Helm release
is named after `metadata.name`; the GitHub-visible fleet name is
`runner_scale_set_name` falling back to `metadata.name`, rendered
explicitly into values so the exported `runs-on` handle and the chart
agree by construction. Both engines fail loudly past the 45-character
name budget (the chart's own template `fail` — a GitHub registration
limit) instead of erroring mid-apply.

## The secret discipline

`githubConfigSecret` always renders as a Secret NAME (the chart's
pre-defined-secret form). The existing-Secret arm references the
user's own Secret; the declared PAT / GitHub App arms materialize
`<name>-github-auth` with the chart's key contract (`github_token` /
`github_app_*`) BEFORE the release — inline credential material is
proto-sensitive and never lands in rendered values. The re-pin after
the escape-hatch merge keeps an override from moving the contract.

## No Helm wait — the verifier owns readiness

The chart's workload is a custom resource the controller reconciles
AFTER the release returns, and the listener needs a valid GitHub
credential to come up — Helm `--wait` would pass trivially while
proving nothing. Both engines skip the wait (the CR-kind precedent);
the E2E verifier owns the reconcile-attempt proof on every lane and
the listener-registered proof on the live lane.

## Container modes

`dind` XOR `kubernetes` XOR `kubernetes-novolume` behind one mode
field; `kubernetes` requires the ephemeral work-volume block
(StorageClass reference + size — spec CEL enforces the pairing both
ways). The runner container override re-states the chart's own
container contract (name `runner`, the run.sh command) because Helm
values LISTS replace — rendering image or resources alone would drop
the fields the mode wiring keys on.

## Composition seams

The work volume's StorageClass and the TLS CA ConfigMap are reference
fields; `controller_service_account` binds a fleet to a
namespace-fenced controller (wired from that kind's
`service_account_name` output); `github_config_url` and the fleet name
export for workflow authors and dashboards.

## Deliberate exclusions

The full runner pod `template` PodSpec, `listenerTemplate` sidecars,
`listenerMetrics` histogram tuning, `scaleSetLabels`, Azure Key Vault
credential sourcing, and per-resource `resourceMeta` — reachable
through `helm_values`, never the primary interface.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
