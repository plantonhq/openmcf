---
title: "Workload Least-Privilege Role"
description: "This preset defines a custom role with exactly the permissions one workload needs — the standard replacement for granting an over-broad predefined role like `roles/storage.admin` to a service account..."
type: "preset"
rank: "01"
presetSlug: "01-workload-least-privilege"
componentSlug: "iam-custom-role-on-google-cloud"
componentTitle: "IAM Custom Role on Google Cloud"
provider: "gcp"
icon: "package"
order: 1
---

# Workload Least-Privilege Role

This preset defines a custom role with exactly the permissions one workload needs — the standard replacement for granting an over-broad predefined role like `roles/storage.admin` to a service account that only writes objects. Pair it with a GcpProjectIamMember grant referencing this role's `name` output.

## When to Use

- A service account needs a handful of specific permissions and every predefined role over-grants
- Security review flagged broad predefined roles on workload identities
- You want permission changes to happen in one place and propagate to every grant instantly

## Key Configuration Choices

- **camelCase `roleId`** — the API rejects hyphens in role IDs; `logBucketWriter`, not `log-bucket-writer`
- **`stage: GA`** — the right label for production roles; it is informational except `DISABLED`, which is an IAM kill switch
- **Two example permissions** — replace with the exact set your workload uses; discover valid values with `gcloud iam list-testable-permissions <resource>`

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project that owns the role | GCP Console or `GcpProject` outputs |
| `logBucketWriter` (roleId) | Unique role ID (3-64 chars; letters, digits, underscores, periods) | Replace with a descriptive camelCase name |
| `Log Bucket Writer` (title) | Console-visible title (max 100 chars) | Human-readable version of the role ID |

## Related Presets

- **02-readonly-auditor** — A read-only role for auditors and dashboards
- **03-ci-cd-deployer** — A deployment role for CI/CD pipelines
