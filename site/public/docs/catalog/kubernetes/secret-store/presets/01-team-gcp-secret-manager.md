---
title: "Team-Scoped GCP Secret Manager (Keyless)"
description: "This preset creates a namespaced store connected to GCP Secret Manager, authenticating through a ServiceAccount in the team's own namespace whose GKE Workload Identity binding authorizes the reads...."
type: "preset"
rank: "01"
presetSlug: "01-team-gcp-secret-manager"
componentSlug: "secret-store"
componentTitle: "Secret Store"
provider: "kubernetes"
icon: "package"
order: 1
---

# Team-Scoped GCP Secret Manager (Keyless)

This preset creates a namespaced store connected to GCP Secret Manager, authenticating through a ServiceAccount in the team's own namespace whose GKE Workload Identity binding authorizes the reads. The store, its identity, and its blast radius all end at the namespace boundary — the per-team posture, with each team carrying its OWN cloud identity.

## When to Use

- A team's secrets live in GCP Secret Manager and only that team's namespace should sync them
- Different teams need DIFFERENT cloud identities (per-team least privilege), not one shared platform identity
- You want keyless authentication — no exported service-account keys

## Key Configuration Choices

- **Namespaced grain** -- only ExternalSecrets in `<team-namespace>` can use this store; no cluster-scoped fencing needed
- **Keyless ServiceAccount auth** (`serviceAccountName`) -- references a ServiceAccount with a Workload Identity binding; create it with `KubernetesServiceAccount` (its workload-identity arm carries the GCP service-account binding) and reference its `service_account_name` output
- **ServiceAccount namespace defaulted** -- on a SecretStore the referenced ServiceAccount defaults to the store's own namespace; no explicit namespace needed
- **Global endpoint** -- no `location` set; add one (e.g. `us-central1`) only for regional secrets

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<team-namespace>` | The team's namespace (store, ServiceAccount, and ExternalSecrets all live here) | `kubectl get namespaces` or `KubernetesNamespace` outputs |
| `<gcp-project-id>` | GCP project ID holding the team's secrets | GCP console or `GcpProject` outputs |
| `<team-secrets-reader-service-account>` | ServiceAccount with a Workload Identity binding allowing `secretmanager.versions.access` | `KubernetesServiceAccount` outputs |

## Related Presets

- **02-vault-approle** -- Use when the team's secrets live in Vault/OpenBao and the cluster has no identity federation with it
- **03-fake-sandbox** -- Use for pipelines and sandboxes, without any external account
