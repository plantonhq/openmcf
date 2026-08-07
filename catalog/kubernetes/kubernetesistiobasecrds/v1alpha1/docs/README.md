# KubernetesIstioBaseCrds — Design Notes

## Purpose

`KubernetesIstioBaseCrds` exists to decouple **CRD installation** from **mesh installation**
for the Istio API family. The typed Istio components (DestinationRule, ServiceEntry,
PeerAuthentication, RequestAuthentication, AuthorizationPolicy, Telemetry, EnvoyFilter)
only need the Istio CRDs present on the cluster to be applied and server-side validated —
they do not need a running control plane. Standing up full `istiod` just to use the
policy/config APIs (or per E2E run) is heavy and unnecessary for validating the API
surface.

This mirrors the Gateway API family's `KubernetesGatewayApiCrds`.

## Version Coupling (single source of truth)

Planton ships **typed** resources for the Istio kinds, generated from a pinned Istio
release. The custom resources Planton emits are therefore frozen to that schema version.
For those objects to apply cleanly (no silent field pruning, no server-side rejection),
the **CRDs installed on the cluster must be the same version**.

Consequently this component has **no user-facing version knob**. The CRD bundle version
is a build-time constant in both IaC modules, bumped in lockstep with the SDK pin. This
deliberately avoids the latent inconsistency in `KubernetesGatewayApiCrds`, whose
user-facing `version` field can default below its typed SDK version.

The pin is always an exact release **tag** (e.g. `1.30.3`), never a release branch. A
branch ref moves as patches land, so the same deployed resource would install different
CRD schemas at different times; a tag keeps installs reproducible and exactly matched to
the generated SDK.

## Relationship to `KubernetesIstio`

| | `KubernetesIstioBaseCrds` | `KubernetesIstio` |
|---|---|---|
| Installs | Istio CRDs only | The control plane: module-owned CRDs + Helm releases `base` and `istiod` (plus `cni` and `ztunnel` in ambient mode) |
| Mechanism | server-side apply of `crd-all.gen.yaml` | server-side-applied CRDs + Helm releases |
| Version field | none (pinned to the SDK) | user-facing `version` (one pin for all charts) |
| Role | prerequisite for the typed Istio components | the running mesh |

Both kinds apply the CRDs via **server-side apply**, so they co-own rather than
conflict. That is a deliberate design point: `KubernetesIstio` keeps its CRDs outside
its Helm release precisely so a CRDs-only cluster can later upgrade to the full mesh
with a plain redeploy (Helm refuses to adopt pre-existing CRDs it does not own). Once a
mesh is installed, `KubernetesIstio` owns everything, CRDs included.

## E2E

`e2e/prerequisite.yaml` is installed automatically by the harness for any of the seven
consuming kinds. `e2e/scenarios/minimal.yaml` exercises this component standalone, and
the harness verifies the Istio CRDs are present on the cluster.
