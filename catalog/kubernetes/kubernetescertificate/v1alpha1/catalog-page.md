# KubernetesCertificate

Requests one signed X.509 certificate from a cert-manager issuer and keeps it renewed. The signed certificate, key, and CA chain land in a TLS Secret that Ingress TLS blocks, Gateway listeners, workloads, and CA issuers reference by name. Covers the complete cert-manager.io/v1 request surface — SAN types through keystore outputs.

## What Gets Created

- **Certificate** — the cert-manager request; issuance and renewal run in-cluster
- **TLS Secret** (by cert-manager) — `spec.secret_name`, holding `tls.crt` / `tls.key` / `ca.crt`, plus optional JKS/PKCS#12 keystores and DER/combined-PEM outputs

## Prerequisites

- cert-manager on the cluster (**KubernetesCertManager**)
- A signing authority: **KubernetesClusterIssuer**, **KubernetesIssuer**, or an external issuer controller

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesCertificate
metadata:
  name: api-cert
spec:
  namespace:
    value: my-namespace
  secretName: api-cert-tls
  issuerRef:
    clusterIssuer:
      name:
        value: letsencrypt-production
  dnsNames:
    - api.example.com
```

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Certificate (and Secret) namespace |
| `certificate_name` | The Certificate resource name |
| `secret_name` | The TLS Secret handle consumers reference |

## Next Steps

Point Ingress `tls.secretName`, Gateway `certificate_refs`, or a CA-backed Issuer's `ca_secret_name` at this certificate's `secret_name` output.
