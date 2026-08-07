# KubernetesExternalSecretsOperator Pulumi Module

Installs the External Secrets Operator from the official Helm chart
(`external-secrets` at `https://charts.external-secrets.io`) as a real Helm
release. The typed spec renders into chart values in `module/values.go`;
the `helm_values` escape hatch merges LAST over them with Helm `-f`
semantics (maps deep-merge, later document wins, lists replace) — the exact
semantic twin of the Terraform module's `helm_release` with
`values = [typed, helm_values]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance labels
   when `create_namespace` is true; otherwise the namespace must exist
2. **Helm Release** — always named `external-secrets` (one installation per
   cluster is an upstream architectural constraint: cluster-scoped CRDs and
   webhook configuration). The chart owns ALL operator objects — the three
   Deployments (controller, webhook, cert-controller), ServiceAccounts,
   RBAC, and CRDs. The controller ServiceAccount name is pinned to
   `external-secrets` and carries the workload-identity annotations (plus
   the `azure.workload.identity/use` pod label on the AKS arm)

Unlike the ExternalDNS module, no credential Secrets are materialized here:
store credentials belong to the store kinds
(KubernetesClusterSecretStore / KubernetesSecretStore), not the operator
install.

## CRD Handling

`installCRDs` renders from `crds.install` (default true). With
`crds.keep_on_uninstall` (default true) the module renders the
`helm.sh/resource-policy: keep` annotation onto the CRDs via the chart's
`crds.annotations` — the chart itself has no keep knob and Helm would
otherwise DELETE the CRDs on uninstall, cascading to every
ExternalSecret/SecretStore object cluster-wide.

## Wait / Atomic Posture

The release installs with `Atomic` + `CleanupOnFail` and waits (600s
timeout) for all three Deployments to become Available. The webhook
validates every ESO resource at admission — an operator whose webhook is
not ready rejects every SecretStore/ExternalSecret apply, so a premature
"success" would just move the failure downstream.

## Usage

```shell
planton pulumi up --manifest e2e/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator is installed in |
| `release_name` | Helm release name (always `external-secrets`) |
| `controller_service_account` | Controller ServiceAccount — bind cloud-side for ambient keyless store access |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → Helm release
- `module/values.go`: typed-spec → chart values rendering (CRDs, controller
  scaling, identity, per-component tuning), escape-hatch merge
- `module/locals.go`: resolved names (release, ServiceAccount) — kept in
  lockstep with the Terraform module's `locals.tf`
- `module/vars.go`: chart identity and the pinned default chart version
