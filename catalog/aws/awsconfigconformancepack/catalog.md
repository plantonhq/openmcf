# AWS Config Conformance Pack

Compliance as one deployable unit: a conformance pack turns a ruleset
— CIS, PCI, or your own — into a single template that deploys, scores,
and unwinds together, per account or across the whole organization.

## What Gets Managed

- The pack and its template (inline or from S3), with parameter
  values.
- Account scope (the default) or organization scope with account
  exclusions.
- Optional S3 delivery of evaluation results.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Config permissions.

### AWS Account

- A running Config recorder in the region
  ([AWS Config Recorder](/cloud-catalog/aws-config-recorder)) — AWS
  rejects pack deployment without one.
- For organization scope: the management account or a delegated
  Config administrator; delivery buckets must be named with the
  `awsconfigconforms` prefix.

## Deploy

### Console

Create the resource from the AWS catalog, paste or point at the
template, pick the scope, and deploy.

### CLI

```bash
planton apply -f conformance-pack.yaml
```

## After Deploy

- The pack reaches CREATE_COMPLETE once its rules materialize; scores
  appear in the Config console's conformance packs view.
- Rule names are prefixed with the pack name — they are pack-owned,
  not standalone rules.
- Destroying deletes the pack and its rules (organization packs
  unwind from member accounts over several minutes).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
