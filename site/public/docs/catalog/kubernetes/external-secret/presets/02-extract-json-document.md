---
title: "Extract a JSON Document (Bulk Pull with Rewrite)"
description: "This preset pulls ALL properties of one structured backend entry — a JSON document of related credentials — into a Kubernetes Secret in a single `dataFrom.extract`, with a regex rewrite that strips..."
type: "preset"
rank: "02"
presetSlug: "02-extract-json-document"
componentSlug: "external-secret"
componentTitle: "External Secret"
provider: "kubernetes"
icon: "package"
order: 2
---

# Extract a JSON Document (Bulk Pull with Rewrite)

This preset pulls ALL properties of one structured backend entry — a JSON document of related credentials — into a Kubernetes Secret in a single `dataFrom.extract`, with a regex rewrite that strips the backend path prefix so the Secret keys come out as bare names. Use it when an application's whole credential set lives in one backend document and enumerating fields one-by-one would just duplicate it.

## When to Use

- The backend entry is a JSON document holding a related set of credentials (host, username, password, ...)
- New fields added to the document should flow into the Secret without editing the manifest
- You accept the bulk-pull trade-off: what lands in the Secret is whatever the document holds, not an explicit list

## Key Configuration Choices

- **`dataFrom.extract`** -- each property of the structured entry becomes one Secret key; contrast with explicit `data` entries, which name every mapping
- **`rewrite`** -- ordered regex steps applied to the pulled keys; the shipped example strips a `prod/app/` prefix (`^prod/app/(.*)$` → `$1`). Adjust the pattern to your backend's path convention, or drop the block to keep keys as-is
- **Cluster store reference** (`kind: ClusterSecretStore`) -- the platform-wide store grain; switch to `SecretStore` (the default) when the store lives in this namespace
- **Default refresh and target** -- hourly refresh (upstream default) into an operator-owned Secret named after the resource, retained on delete

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<app-namespace>` | Namespace the ExternalSecret and its Secret live in | `kubectl get namespaces` or `KubernetesNamespace` outputs |
| `<cluster-secret-store-name>` | The store to sync from | `KubernetesClusterSecretStore` `store_name` output |
| `<backend-secret-path>` | The structured entry's name/path in the backend | Your secret backend's console/CLI |

The `rewrite` regex (`^prod/app/(.*)$`) is an example — match it to your backend's key naming, or remove the block entirely.

## Related Presets

- **01-explicit-keys** -- Use when every synced key should be spelled out and reviewable
- **03-docker-registry-template** -- Use to render an image pull secret from synced registry credentials
