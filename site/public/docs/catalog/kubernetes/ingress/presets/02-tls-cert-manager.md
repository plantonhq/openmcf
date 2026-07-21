---
title: "TLS with cert-manager"
description: "This preset serves one host over HTTPS with a certificate issued automatically by cert-manager. The `tls` block names the certificate Secret; the `cert-manager.io/cluster-issuer` annotation tells..."
type: "preset"
rank: "02"
presetSlug: "02-tls-cert-manager"
componentSlug: "ingress"
componentTitle: "Ingress"
provider: "kubernetes"
icon: "package"
order: 2
---

# TLS with cert-manager

This preset serves one host over HTTPS with a certificate issued automatically by cert-manager. The `tls` block names the certificate Secret; the `cert-manager.io/cluster-issuer` annotation tells cert-manager to issue for the listed hosts.

## When to Use

- Production HTTPS exposure with automated certificate issuance and renewal
- Any cluster where cert-manager is installed with a ClusterIssuer (e.g. Let's Encrypt)

## Key Configuration Choices

- **The Secret does not exist yet — and should not.** cert-manager watches annotated Ingresses, runs the ACME flow for the TLS hosts, and creates the `kubernetes.io/tls` Secret under exactly the name written in `secret_name` (`app-example-com-tls`). Do not pre-create it
- **`cluster-issuer` annotation** — references a cluster-scoped ClusterIssuer; use `cert-manager.io/issuer` instead for a namespaced Issuer
- **TLS hosts match rule hosts** — the certificate's SANs must cover the hosts served; keeping the `tls.hosts` and `rules.host` lists aligned is the contract
- **Until the certificate is issued**, most controllers serve the host with a temporary default certificate; issuance typically completes within a minute or two of DNS pointing at the controller

Without cert-manager installed, this manifest still deploys — the hosts are just served with the controller's default certificate until the named Secret appears.

## Values to Replace

| Value | Description |
|---|---|
| `app.example.com` | The public hostname (must be publicly resolvable for ACME HTTP-01) |
| `letsencrypt-prod` | Your ClusterIssuer name (`kubectl get clusterissuer`) |
| `app-example-com-tls` | The Secret name cert-manager will create — any name you like |
| `web` / `web-svc` / `8080` | Namespace, backend Service, and port |

## Related Presets

- **01-single-host** — the same routing without TLS
