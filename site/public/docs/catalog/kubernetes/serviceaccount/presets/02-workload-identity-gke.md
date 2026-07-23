---
title: "GKE Workload Identity"
description: "This preset creates a ServiceAccount bound to a GCP service account via GKE Workload Identity. Pods running as this identity call Google Cloud APIs keylessly — the cluster's OIDC issuer vouches for..."
type: "preset"
rank: "02"
presetSlug: "02-workload-identity-gke"
componentSlug: "serviceaccount"
componentTitle: "ServiceAccount"
provider: "kubernetes"
icon: "package"
order: 2
---

# GKE Workload Identity

This preset creates a ServiceAccount bound to a GCP service account via GKE Workload Identity. Pods running as this identity call Google Cloud APIs keylessly — the cluster's OIDC issuer vouches for the ServiceAccount and GCP exchanges that token for credentials; no keys are stored in the cluster.

## When to Use

- Pods on GKE that need GCP APIs (Cloud Storage, Pub/Sub, Cloud DNS, Secret Manager, ...)
- Replacing exported GCP service account key files mounted as secrets — the pattern this mechanism exists to eliminate

## Key Configuration Choices

- **`workloadIdentity.gke.serviceAccountEmail`** — the GCP service account to act as; the module emits it as the `iam.gke.io/gcp-service-account` annotation, the exact key GKE's webhook expects
- **Both halves are required** — the annotation alone grants nothing. The GCP service account must carry a `roles/iam.workloadIdentityUser` binding for member `serviceAccount:<project>.svc.id.goog[<namespace>/<ksa-name>]`, and the cluster must have Workload Identity enabled
- **Name and namespace are part of the trust** — the GCP-side member string embeds both; renaming or moving this ServiceAccount silently breaks the federation until the binding is updated

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Namespace for the ServiceAccount (embedded in the GCP-side trust member) | Your namespace management |
| `<your-gcp-service-account-email>` | GCP service account email, e.g. `app@my-project.iam.gserviceaccount.com` | GCP Console → IAM & Admin → Service Accounts, or your GcpServiceAccount resource's outputs |

Also rename `gcp-app-identity` to match your workload — the name is embedded in the GCP-side trust member too.

## Related Presets

- **01-basic** — identity with no cloud federation
- **03-workload-identity-eks-irsa** — the AWS equivalent
- **04-image-pull-secrets** — private registry credentials and automount hardening
