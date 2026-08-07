---
title: "Cluster Issuer"
description: "Cluster Issuer deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesclusterissuer"
---

# Kubernetes Cluster Issuer

Creates one cert-manager ClusterIssuer — a cluster-scoped certificate authority front-end that Certificate resources in ANY namespace can request signed certificates from. The ClusterIssuer is named after this resource; Certificates select it by that name, and ingress-shim annotations (`cert-manager.io/cluster-issuer`) use the same name. Four signing backends are modeled: ACME (public CAs like Let's Encrypt), CA (sign with a Secret-held keypair), SelfSigned (bootstrap/dev), and Vault (Vault/OpenBao PKI).

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **ClusterIssuer** -- the cert-manager custom resource, named after `metadata.name`, configured with the selected signing backend
- **Credential Secrets** -- wherever the backend needs a credential (a Cloudflare token, AWS secret key, Vault token, TSIG key), the module materializes it as a Kubernetes Secret in cert-manager's cluster-resource namespace and wires the CR's secretRef to it — you declare the value once, it never appears in rendered manifests
- **ACME Account Key Secret** -- created by cert-manager itself on first ACME registration (ACME backend only)

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with credentials for the target cluster.

### Kubernetes Cluster

- **Cert Manager must be installed first** -- the ClusterIssuer is a cert-manager custom resource. Reference the Kubernetes Cert Manager resource's `cluster_resource_namespace` output for correct-by-construction wiring.
- For keyless DNS-01 (Route53 via IRSA, Cloud DNS via GKE Workload Identity, Azure DNS via AKS Workload Identity), the cert-manager controller's workload identity must be configured on the Kubernetes Cert Manager resource.
- For token-based DNS-01 (Cloudflare, DigitalOcean), have the scoped token ready — for Cloudflare, a token with Zone:Zone:Read and Zone:DNS:Edit on the target zones.

## Deploy

### Console

Open the deployment store, find **Cluster Issuer**, and click **Deploy**. The creation wizard walks you through placement (the cluster-resource namespace), the signing backend, and — for ACME — the challenge solvers. Start from the **Cloudflare** preset for the most common public-TLS shape in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesClusterIssuer
metadata:
  name: letsencrypt-production
  org: acme-corp
  env: prod
spec:
  certManagerNamespace:
    value: cert-manager
  config:
    acme:
      email: platform@example.com
      solvers:
        - dns01:
            cloudflare:
              apiToken:
                token: $secret/cloudflare-dns-token
```

```shell
planton apply -f cluster-issuer.yaml
```

This registers an ACME account with Let's Encrypt production and satisfies challenges by publishing DNS TXT records through Cloudflare — one catch-all solver serving every namespace in the cluster.

## Key Configuration

These are the most important decisions when configuring a Cluster Issuer. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**ClusterIssuer vs Issuer** -- Identical signing capabilities (the config message is shared), different scope. Use a ClusterIssuer for the platform-wide CA (one `letsencrypt-production` serving every team); use a namespace-scoped Kubernetes Issuer when a namespace needs its own CA or its own credential blast radius.

**Pick the signing backend** -- `acme` for browser-trusted public TLS, `ca` to sign with a CA keypair from a Kubernetes Secret (internal PKI, mTLS), `self_signed` to bootstrap a CA chain or for development, `vault` for a centralized Vault/OpenBao PKI engine.

**Stage before production (ACME)** -- Let's Encrypt production enforces strict rate limits that a misconfigured setup can exhaust for a whole domain. Point `acme.server` at the staging endpoint while testing, then switch.

**Solvers decide how ownership is proven** -- DNS-01 (publish a TXT record) is the only challenge type that can issue wildcard certificates and works for non-public services; HTTP-01 (serve a token on port 80) needs public reachability. A single catch-all solver is the common case; add selector-scoped solvers when different DNS zones need different providers.

**Keyless where the platform allows it** -- DNS-01 providers that authenticate ambiently (Route53, Cloud DNS, Azure DNS) inherit the identity configured on the cert-manager controller — leave their static credential fields empty. Token providers (Cloudflare, DigitalOcean) always need their token.

**Credentials are declared, not pre-created** -- Wherever upstream cert-manager expects a `secretRef` to a hand-made Secret, this spec accepts the credential value directly (marked sensitive). The module materializes the Secret and wires the reference.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.certManagerNamespace` | KubernetesCertManager (`status.outputs.cluster_resource_namespace`) | Where credential Secrets are materialized and cert-manager reads them for cluster-scoped resources |
| `spec.config.ca.caSecretName` | KubernetesCertificate (`status.outputs.secret_name`) | CA backend: the Secret holding the CA keypair — the standard CA-chain bootstrap |
| `spec.config.vault.kubernetesAuth.serviceAccountName` | KubernetesServiceAccount (`metadata.name`) | Vault backend: the ServiceAccount whose token authenticates to Vault |
| `spec.config.acme.solvers[].dns01.route53.serviceAccount.serviceAccountName` | KubernetesServiceAccount (`metadata.name`) | Per-issuer IRSA for Route53 without changing the controller identity |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `cluster_issuer_name` | Name of the created ClusterIssuer (equals `metadata.name`) | Kubernetes Certificate `issuerRef.clusterIssuer.name`; the `cert-manager.io/cluster-issuer` ingress-shim annotation |
| `secrets_namespace` | Namespace where this issuer's credential Secrets were materialized | Auditing credential placement |
| `acme_account_key_secret_name` | The ACME account private key Secret cert-manager creates (empty for non-ACME backends) | Migrating the ACME account between clusters |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Public TLS via Cloudflare** -- Let's Encrypt with a DNS-01 catch-all solver authenticated by a scoped API token. Start from the **Cloudflare** preset.

**Keyless on GKE** -- Cloud DNS DNS-01 with no static key; the controller's GKE Workload Identity signs the API calls. Start from the **GCP Cloud DNS** preset.

**Keyless on EKS** -- Route53 DNS-01 through the controller's IRSA identity. Start from the **AWS Route53** preset.

**Enterprise PKI** -- Certificates signed by a Vault PKI engine, authenticated via Vault's Kubernetes auth method (keyless). Start from the **Vault** preset.

## Works With

- **Kubernetes Cert Manager** -- must be installed first; provides the controller, the CRDs, and the cluster-resource namespace this issuer's Secrets live in.
- **Kubernetes Certificate** -- requests certificates from this issuer by name, from any namespace.
- **Kubernetes Ingress Nginx / Kubernetes Gateway** -- serve HTTPS with the issued certificate Secrets; ingress-shim annotations can name this issuer directly.
- **Kubernetes Service Account** -- composes the per-issuer IRSA and Vault Kubernetes-auth identities.
