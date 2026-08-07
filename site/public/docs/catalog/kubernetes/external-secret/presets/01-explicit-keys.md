---
title: "Explicit Keys (Application Credentials)"
description: "This preset syncs two named fields (`username`, `password`) from one structured backend entry into a Kubernetes Secret, refreshed hourly. Each mapping is explicit and reviewable — the standard form..."
type: "preset"
rank: "01"
presetSlug: "01-explicit-keys"
componentSlug: "external-secret"
componentTitle: "External Secret"
provider: "kubernetes"
icon: "package"
order: 1
---

# Explicit Keys (Application Credentials)

This preset syncs two named fields (`username`, `password`) from one structured backend entry into a Kubernetes Secret, refreshed hourly. Each mapping is explicit and reviewable — the standard form for application credentials, and the most common ExternalSecret shape. The materialized Secret is named after the resource; workloads reference it in env `valueFrom` / volume `secretName`.

## When to Use

- An application needs specific credentials (database, API keys) from a backend entry
- You want every synced key visible in review — no bulk pulls
- The backend entry is a structured JSON document and you need individual fields from it

## Key Configuration Choices

- **Explicit `data` entries** -- one backend key + `property` per Secret key; nothing lands in the Secret that is not spelled out here
- **`property` extraction** -- pulls one field from a structured entry (a JSON document); drop `property` to sync the whole value
- **Cluster store reference** (`kind: ClusterSecretStore`) -- the platform-wide store grain; switch to `SecretStore` (the default) when the store lives in this namespace
- **Hourly refresh** (`refreshInterval: 1h`) -- the upstream default; backend rotation reaches the cluster within this window
- **Default target** -- no `target` block: the Secret is named after the resource, owned by the operator, retained on delete (the safe posture)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<app-namespace>` | Namespace the ExternalSecret and its Secret live in | `kubectl get namespaces` or `KubernetesNamespace` outputs |
| `<cluster-secret-store-name>` | The store to sync from | `KubernetesClusterSecretStore` `store_name` output |
| `<backend-secret-path>` | The entry's name/path in the backend (e.g. `prod/app/database` in AWS Secrets Manager, `secret/data/app` in Vault KV v2) | Your secret backend's console/CLI |

## Related Presets

- **02-extract-json-document** -- Use to pull ALL fields of a structured entry at once
- **03-docker-registry-template** -- Use to render an image pull secret from synced registry credentials
