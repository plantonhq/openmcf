---
title: "Tuned Multi-Team Installation"
description: "This preset sizes the External Secrets Operator for clusters where many teams sync many secrets: reconcile concurrency raised to 5, explicit controller resources, two webhook replicas (every ESO..."
type: "preset"
rank: "03"
presetSlug: "03-tuned-multi-team"
componentSlug: "external-secrets-operator"
componentTitle: "External Secrets Operator"
provider: "kubernetes"
icon: "package"
order: 3
---

# Tuned Multi-Team Installation

This preset sizes the External Secrets Operator for clusters where many
teams sync many secrets: reconcile concurrency raised to 5, explicit
controller resources, two webhook replicas (every ESO resource apply goes
through the webhook — it is the availability-sensitive path), and explicit
webhook/cert-controller sizing. No ambient identity — multi-team clusters
should give each store its own identity in its auth block.

## When to Use

- Clusters with hundreds of ExternalSecrets where sync latency at the
  default `concurrent: 1` becomes visible
- Multi-team clusters where the operator is platform infrastructure and
  should be sized/guarded like it
- Environments with resource quotas that require explicit requests/limits

## Key Configuration Choices

- **`concurrent: 5`** — five ExternalSecrets reconcile in parallel; raise
  further if store backends tolerate the API rate
- **Controller sizing** (`resources`) — explicit requests/limits instead of
  the chart's unset defaults
- **Webhook redundancy** (`webhook.replicas: 2`) — a single webhook replica
  being rescheduled blocks every SecretStore/ExternalSecret apply
  cluster-wide
- **cert-controller sizing** — explicit but small; it only bootstraps and
  rotates the webhook certificate
- **No ambient `workloadIdentity`** — per-store identities keep teams'
  secret access isolated

## Placeholders to Replace

No placeholders — adjust the sizing numbers to your cluster's scale.

## Related Presets

- **01-minimal** — chart defaults for smaller clusters
- **02-eks-ambient-identity** — one ambient IAM role when per-store
  isolation is not needed
