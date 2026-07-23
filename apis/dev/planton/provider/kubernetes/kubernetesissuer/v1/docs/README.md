# KubernetesIssuer: Research and Design

## Introduction

An Issuer is the namespace-scoped sibling of ClusterIssuer: the same signing
machinery, readable only by Certificates in its own namespace, with every
Secret it needs (credentials, CA keypairs) confined to that namespace.

## Design Authority and the Shared-Config Guarantee

Upstream defines Issuer and ClusterIssuer with an IDENTICAL spec
(`cert-manager.io_issuers.yaml` differs from the clusterissuers schema only
in scope). The Planton kinds encode that identity structurally: both specs
embed the same shared `CertManagerIssuerConfig` proto, both Pulumi modules
render through one shared spec builder, and the Terraform locals are exact
namespace-scoped siblings. Full backend surface — ACME (all solvers/providers),
CA, self-signed, Vault — with Venafi deliberately excluded; see the
KubernetesClusterIssuer research doc for the per-backend rationale, which
applies here verbatim.

## Why the Namespace Scope Is a Feature

- **Credential blast radius**: a team's DNS token or CA private key is
  readable only inside the team's namespace — a ClusterIssuer's credentials
  sit in cert-manager's cluster-resource namespace, trusted cluster-wide.
- **The internal-PKI bootstrap** lives entirely in one namespace:
  self-signed Issuer → root CA Certificate (`is_ca`) → CA Issuer signing
  with the root's Secret (the `ca_secret_name` FK defaults to a
  KubernetesCertificate's secret output) → leaf certificates.

## Differences from ClusterIssuer (all scope-derived)

| Aspect | Issuer | ClusterIssuer |
|---|---|---|
| Serves | Same-namespace Certificates | Every namespace |
| Credential Secrets | The Issuer's own namespace | cert-manager's cluster-resource namespace |
| Namespace field | `spec.namespace` (KubernetesNamespace FK) | `spec.cert_manager_namespace` (KubernetesCertManager FK) |

## Engine Mechanics

Identical to KubernetesClusterIssuer: Pulumi via the shared builder +
CustomResource; Terraform via `kubectl_manifest` (plannable before the CRDs
exist); credential Secrets created before the CR so its secretRefs never
dangle; neither engine waits for Ready.

## E2E

Self-signed and CA-chain lanes run deterministically on the kind cluster —
the CA-chain scenario resolves a live root-CA Certificate fixture through
the FK and proves cert-manager accepts the chain (Ready requires reading a
valid CA keypair from the referenced Secret).
