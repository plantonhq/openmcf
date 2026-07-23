# KubernetesExternalDns Pulumi Module

Installs ExternalDNS from the official Helm chart (`external-dns` at
`https://kubernetes-sigs.github.io/external-dns/`) as a real Helm release.
The typed spec renders into chart values in `module/values.go`; the
`helm_values` escape hatch merges LAST over them with Helm `-f` semantics
(maps deep-merge, later document wins, lists replace) — the exact semantic
twin of the Terraform module's `helm_release` with
`values = [typed, helm_values]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance labels
   when `create_namespace` is true; otherwise the namespace must exist
2. **Credential Secrets** (provider-dependent) — declared static
   credentials materialize as deterministically-named Kubernetes Secrets so
   the credential never appears in chart values or pod specs:
   - Cloudflare: `<name>-cloudflare-credentials` (`api-token`), consumed as
     `CF_API_TOKEN`
   - AWS static keys: `<name>-aws-credentials`, consumed as
     `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY`
   - GCP key: `<name>-gcp-credentials`, mounted with
     `GOOGLE_APPLICATION_CREDENTIALS` pointing at it
   - Azure: `<name>-azure-config` carrying a rendered `azure.json`
     (identity selection: service principal > Workload Identity > managed
     identity), mounted at the controller's default config path
   Keyless installs (workload identity / ambient) and the webhook/in-memory
   arms materialize nothing
3. **Helm Release** — named after `metadata.name` (NOT a fixed chart name:
   multiple instances per cluster are a first-class pattern). The chart's
   `fullnameOverride` is pinned to the release name so every chart object
   carries a deterministic, manifest-derived name; the controller
   ServiceAccount name is pinned the same way and carries the
   workload-identity annotations (plus the `azure.workload.identity/use`
   pod label on the AKS arm)

## Wait / Atomic Posture

The release installs with `Atomic` + `CleanupOnFail` and waits (300s
timeout) for the controller Deployment to become Available — a controller
that never starts fails THIS deploy, not the first record sync. Note the
controller validates provider credentials at first zone sync, not at
startup: an install with wrong credentials still goes green and surfaces in
controller logs, matching upstream behavior.

## Usage

```shell
planton pulumi up --manifest hack/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace ExternalDNS is installed in |
| `release_name` | Helm release name (equals `metadata.name`) |
| `service_account_name` | Controller ServiceAccount — bind cloud-side for keyless provider access |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → credential Secrets → Helm release
- `module/values.go`: typed-spec → chart values rendering, provider arm
  wiring (extraArgs/env/volumes), escape-hatch merge
- `module/secrets.go`: credential Secret materialization, `azure.json`
  rendering
- `module/locals.go`: resolved names (release, ServiceAccount, Secrets) —
  kept in lockstep with the Terraform module's `locals.tf`
- `module/vars.go`: chart identity and the pinned default chart version
