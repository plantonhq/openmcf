# KubernetesCertificate Guide

The judgment this guide carries: certificates on this platform arrive two
ways — a cert-manager annotation on an Ingress, or this explicit kind —
and the choice decides whether the certificate is a visible, reusable
node or a side effect of one Ingress.

## Annotation or explicit Certificate?

- **The annotation path** (`cert-manager.io/cluster-issuer` on an
  Ingress, with the `tls` block naming a not-yet-existing Secret) is the
  low-ceremony default for ONE Ingress needing ONE certificate:
  cert-manager creates and renews the Secret invisibly.
- **This kind** earns its node when the certificate is its own artifact:
  a wildcard shared by many consumers (the ingress controller's
  `defaultTlsCertificate` references this kind's secret output), a
  Gateway listener's certificate, an mTLS/SPIFFE workload identity, an
  internal CA's own certificate, or any certificate whose parameters
  (lifetime, key algorithm, usages, output formats — the full surface on
  [reference.md](reference.md)) must be controlled rather than defaulted.

What breaks when chosen wrong: annotation-issued certs are bound to their
Ingress — reusing one elsewhere means fishing a Secret name out of an
annotation side effect nothing else references; an explicit Certificate
for a single simple Ingress is just one more manifest.

## Scope pairing with the issuer

A ClusterIssuer signs from any namespace; a namespace-scoped
KubernetesIssuer only signs Certificates beside it (the scope judgment
lives in the [ClusterIssuer guide](../../kubernetesclusterissuer/v1alpha1/GUIDE.md)
and the Issuer's own reference page). The Certificate and its output
Secret always land in THIS resource's namespace — consumers must live
there too.

## On the diagram

An explicit Certificate is a node consumers reference through its secret
output — the wildcard-default seam into KubernetesIngressNginx renders as
a real edge. Annotation-issued certificates render as nothing, which is
acceptable exactly when nothing else will ever consume them.

## Pairs well with

- KubernetesClusterIssuer / KubernetesIssuer — who signs (wired via
  `issuerRef`).
- KubernetesIngressNginx — the `defaultTlsCertificate` wildcard seam.
- KubernetesGateway — listeners consuming the certificate Secret.
