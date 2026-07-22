---
title: "GCP Secret Manager with Workload Identity (Keyless)"
description: "This preset creates a cluster-wide store connected to GCP Secret Manager, authenticating through a referenced ServiceAccount whose GKE Workload Identity binding authorizes the reads. No..."
type: "preset"
rank: "02"
presetSlug: "02-gcp-secret-manager-workload-identity"
componentSlug: "kubernetesclustersecretstore"
componentTitle: "KubernetesClusterSecretStore"
provider: "kubernetes"
icon: "package"
order: 2
---

# GCP Secret Manager with Workload Identity (Keyless)

This preset creates a cluster-wide store connected to GCP Secret Manager, authenticating through a referenced ServiceAccount whose GKE Workload Identity binding authorizes the reads. No service-account key touches the cluster — the production posture on GKE.

## When to Use

- Your secrets live in GCP Secret Manager and your cluster runs on GKE
- You want keyless authentication — no exported service-account keys
- Every team should sync from one platform-wide store (add `conditions` to fence namespaces if not)

## Key Configuration Choices

- **Project reference** (`projectId`) -- the GCP project the secrets live in; accepts a literal project ID or a `valueFrom` reference to a `GcpProject` resource
- **Keyless ServiceAccount auth** (`serviceAccountName`) -- references a ServiceAccount with a Workload Identity binding; create it with `KubernetesServiceAccount` (its workload-identity arm carries the GCP service-account binding) and reference its `service_account_name` output
- **Explicit ServiceAccount namespace** -- required on a ClusterSecretStore; cluster scope has no home namespace to default to
- **Global endpoint** -- no `location` set; add one (e.g. `us-central1`) only for regional secrets

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project ID holding the secrets | GCP console or `GcpProject` outputs |
| `<eso-reader-service-account>` | ServiceAccount with a Workload Identity binding allowing `secretmanager.versions.access` | `KubernetesServiceAccount` outputs or `kubectl get sa -n external-secrets` |

## Related Presets

- **01-aws-secrets-manager-irsa** -- Use when secrets live in AWS Secrets Manager
- **03-azure-key-vault-workload-identity** -- Use when secrets live in Azure Key Vault
- **04-vault-kubernetes-auth** -- Use when secrets live in Vault or OpenBao
