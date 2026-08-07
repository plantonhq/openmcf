---
title: "Self-Signed Kubelets (kind / k3s / kubeadm / on-prem)"
description: "This preset installs metrics-server on clusters whose kubelets serve self-signed certificates: kind, k3s, kubeadm without kubelet certificate rotation, and many on-prem setups. `kubeletInsecureTls:..."
type: "preset"
rank: "02"
presetSlug: "02-self-signed-kubelets"
componentSlug: "metrics-server"
componentTitle: "Metrics Server"
provider: "kubernetes"
icon: "package"
order: 2
---

# Self-Signed Kubelets (kind / k3s / kubeadm / on-prem)

This preset installs metrics-server on clusters whose kubelets serve
self-signed certificates: kind, k3s, kubeadm without kubelet certificate
rotation, and many on-prem setups. `kubeletInsecureTls: true` is THE
critical knob here — without it metrics-server can never verify a kubelet,
never completes its first scrape, never passes its readiness probe, and the
atomic install fails with a readiness timeout (loudly, by design — the
alternative is HPAs that silently never scale).

## When to Use

- You run kind, k3s, kubeadm, or a self-managed cluster whose kubelet
  serving certificates are not signed by the cluster CA
- An install without `kubeletInsecureTls` failed its readiness wait

## Key Configuration Choices

- **`kubeletInsecureTls: true`** — skips kubelet serving-certificate
  verification on the scrape side (this is unrelated to the APIService
  serving side, which keeps its own defaults)
- **Dedicated `metrics-server` namespace with `createNamespace: true`** —
  the module creates and owns the namespace, keeping the install's objects
  separable (`kube-system` works too; set `createNamespace: false` there)

## Placeholders to Replace

No placeholders — this preset is directly deployable.

## Related Presets

- **01-managed-cloud** — EKS / AKS, where kubelet certificates verify
  against the cluster CA and this knob stays false
- **03-ha-verified-tls** — production hardening: HA replicas and a verified
  serving-certificate chain via cert-manager
