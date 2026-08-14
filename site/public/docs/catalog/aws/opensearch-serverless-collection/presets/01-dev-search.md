---
title: "Dev Search Collection"
description: "This preset creates the cheapest usable search collection: standby replicas disabled (half the OCU floor), public network reachability, and a data-access rule granting one application role read/write."
type: "preset"
rank: "01"
presetSlug: "01-dev-search"
componentSlug: "opensearch-serverless-collection"
componentTitle: "OpenSearch Serverless Collection"
provider: "aws"
icon: "package"
order: 1
---

# Dev Search Collection

This preset creates the cheapest usable search collection: standby replicas disabled (half the OCU floor), public network reachability, and a data-access rule granting one application role read/write.

## When to Use

- Development and testing of search features
- Prototyping OpenSearch integrations without sizing domains
- Feature-branch environments where HA is unnecessary

## Key Configuration Choices

- **`standbyReplicas: DISABLED`** — 0.5+0.5 OCU floor instead of 2+2; fixed at create time
- **Public network posture** (the omitted `network` block) — reachability only; every request still needs SigV4 auth plus this data-access rule
- **One data-access rule** — without it the collection is write-proof and read-proof: IAM permissions alone grant nothing in OpenSearch Serverless

## Cost Note

The OCU floor bills while the collection exists — destroy dev collections when idle.
