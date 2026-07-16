---
title: "Workload Identity Service Account"
description: "This preset creates a GCP service account designed for GKE Workload Identity. No JSON key is generated -- pods authenticate via KSA-to-GSA binding instead. The account is granted logging, monitoring,..."
type: "preset"
rank: "01"
presetSlug: "01-workload-identity"
componentSlug: "service-account"
componentTitle: "Service Account"
provider: "gcp"
icon: "package"
order: 1
---

# Workload Identity Service Account

This preset creates a GCP service account designed for GKE Workload Identity. No JSON key is generated -- pods authenticate via KSA-to-GSA binding instead. The account is granted logging, monitoring, and secret access roles for typical application workloads.

## When to Use

- GKE application pods that need to access GCP services (Secret Manager, Cloud SQL, GCS)
- Any workload using Workload Identity where a JSON key is not needed
- Service accounts following the least-privilege principle for logging, monitoring, and secrets

## Key Configuration Choices

- **No key** (`createKey: false`) -- Workload Identity eliminates exported keys, reducing secret sprawl
- **Logging** (`roles/logging.logWriter`) -- write application logs to Cloud Logging
- **Monitoring** (`roles/monitoring.metricWriter`) -- write custom metrics to Cloud Monitoring
- **Secrets** (`roles/secretmanager.secretAccessor`) -- read secrets from Secret Manager
- **Pair with GcpGkeWorkloadIdentityBinding** -- bind this GSA to a KSA to enable pod access

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `my-app-workload` (serviceAccountId) | Service account ID (6-30 chars, lowercase) | Replace with a descriptive ID for your workload |
| `<gcp-project-id>` | GCP project ID | `GcpProject` outputs |

## Related Presets

- **02-ci-cd-pipeline** -- Use for CI/CD service accounts that need a JSON key for external authentication
- **03-identity-with-first-class-grants** -- Pure identity whose access lives entirely in GcpProjectIamMember nodes
