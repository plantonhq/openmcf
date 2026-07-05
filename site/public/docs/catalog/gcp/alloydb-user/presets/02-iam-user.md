---
title: "IAM User (ALLOYDB_IAM_USER)"
description: "This preset maps a GCP IAM principal (here, a service account email) to an AlloyDB database user with no stored password."
type: "preset"
rank: "02"
presetSlug: "02-iam-user"
componentSlug: "alloydb-user"
componentTitle: "AlloyDB User"
provider: "gcp"
icon: "package"
order: 2
---

# IAM User (ALLOYDB_IAM_USER)

This preset maps a GCP IAM principal (here, a service account email) to an AlloyDB database user with no stored password.

## When to Use

- Workloads that already authenticate through AlloyDB Auth Proxy or Language Connectors with IAM
- Eliminating long-lived database passwords from application configuration

## Key Configuration Choices

- **userId is the IAM principal email** — must match the identity presented at connect time
- **No password** — spec CEL rejects passwords on ALLOYDB_IAM_USER

## Related Components

- [GcpServiceAccount](/docs/catalog/gcp/gcpserviceaccount) — the identity behind IAM users
