---
title: "Keycloak Operator"
description: "Keycloak Operator deployment documentation"
icon: "package"
order: 100
componentName: "kuberneteskeycloakoperator"
---

# Keycloak Operator

The lifecycle manager for Keycloak on Kubernetes: the official
operator that turns `KubernetesKeycloak` declarations into running,
clustered Keycloak servers. Keycloak ships no official Helm chart —
the operator IS the first-party distribution, and this component
installs it faithfully from the release manifests.

## Highlights

- **The manager, cleanly separated** — this resource installs the
  operator alone; servers are `KubernetesKeycloak` declarations it
  reconciles. Installing the operator deploys no Keycloak.
- **First-party distribution, pinned** — the keycloak-k8s-resources
  release manifests at an exact tag: 16 plain-YAML documents plus 4
  CRDs, no webhooks, no cert-manager dependency, no install hooks.
- **Watch scope as one field** — namespaced by default (isolated
  per-team operator+Keycloak stacks side by side), cluster-wide when
  one operator should serve every namespace (at most one per
  cluster).
- **No version knob, on purpose** — the declaration kind renders its
  CR against this bundle's CRD schema; a selectable version would
  drift them apart. Upgrades arrive as module updates.
- **Air-gap ready** — operator-image and default-server-image
  (`RELATED_IMAGE_KEYCLOAK`) overrides for private mirrors, keeping
  the pinned tags.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
