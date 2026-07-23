# KubernetesIstio Terraform Module

Installs the Istio control plane from the official Helm charts (`base`,
`istiod`, and in ambient mode `cni` + `ztunnel`, all at one pinned version
from `https://istio-release.storage.googleapis.com/charts`) as real Helm
releases, in upstream's own order.

## Resources Created

1. **Namespace** (`kubernetes_namespace_v1`, optional) — created and owned
   when `create_namespace` is set, with the standard resource-identity labels
2. **Istio CRDs** (`kubectl_manifest` per CRD, server-side apply) — the
   upstream `crd-all.gen.yaml` bundle at the pinned version, fetched via the
   `http` data source and split into one resource per CRD keyed by the CRD's
   own name (stable state addresses across bundle reorderings)
3. **`helm_release.base`** (`istio-base`) — validation-webhook plumbing, with
   `base.excludedCRDs` covering the entire bundle so the release templates NO
   CRDs
4. **`helm_release.istiod`** (`istiod`, or `istiod-<revision>` for a named
   revision) — the control plane; the module waits for real readiness (a
   control plane whose webhooks are not serving rejects every mesh-config
   apply)
5. **`helm_release.cni`** (`istio-cni`) — the node agent, in ambient mode or
   when `cni.enabled` in sidecar mode
6. **`helm_release.ztunnel`** (`ztunnel`) — the ambient per-node L4 proxy

All releases are atomic with cleanup-on-fail: a failed install never leaves a
half-deployed mesh.

## Why the CRDs Are Module-Owned

Helm refuses to adopt CRDs that already exist without its ownership metadata.
If the `base` release owned the CRDs, a cluster running the CRDs-only
KubernetesIstioBaseCrds kind could never upgrade to the full mesh.
Server-side-applied CRDs are co-ownable by both kinds, making that migration
a plain redeploy. Destroying the module removes the CRDs with everything
else — mesh configuration objects cascade.

## Values Rendering (per-release documents)

The typed spec renders into per-chart values in `locals.tf`
(`base_typed_values`, `istiod_typed_values`, `cni_typed_values`,
`ztunnel_typed_values`). Each release's `helm_values` escape hatch is passed
as a SECOND values document, which the provider merges over the first with
Helm `-f` semantics.

**Parity:** the Pulumi module reaches byte-identical values through its own
rendering and deep-merge — chart identity (repo/name/version) and every typed
mapping are kept in lockstep between the two engines, so both deploy the same
product from one manifest.

## Usage

```bash
planton tofu apply --manifest istio.yaml
```

## Local Development

```bash
terraform init
terraform validate
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Inputs

See `variables.tf` for the full variable specification (generated from the
spec proto).

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Control-plane namespace |
| `istiod_service_name` | istiod Service name (`istiod`, or `istiod-<revision>`) |
| `revision` | Installed control-plane revision (`default` when unnamed) |
| `gateway_class_name` | GatewayClass istiod serves (`istio`) |
| `trust_domain` | The mesh's trust domain |
| `dataplane_mode` | `sidecar` or `ambient` |
