# AWS Config Aggregator

One compliance view across every account and region: the Config
aggregator collects resource configurations and rule results into a
single queryable rollup — the screen a security team actually opens.

## What Gets Managed

- The aggregator: an explicit account list or the whole AWS
  Organization, across listed regions or all of them.
- The reciprocal grants source accounts issue so an aggregator may
  collect from them.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Config permissions.

### AWS Account

- Nothing for a same-account rollup. For a cross-account rollup, each
  source account deploys this component with the grant arm. For an
  organization source: the management account (or delegated Config
  administrator) and an IAM role with the
  `AWSConfigRoleForOrganizations` policy
  ([AWS IAM Role](/cloud-catalog/aws-iam-role)).
- Data appears only from accounts/regions with a running Config
  recorder ([AWS Config Recorder](/cloud-catalog/aws-config-recorder))
  — the aggregator itself needs none.

## Deploy

### Console

Create the resource from the AWS catalog, pick the source shape
(accounts or organization), and deploy. In source accounts, deploy the
same kind with just the grant.

### CLI

```bash
planton apply -f config-aggregator.yaml
```

## After Deploy

- The aggregated view fills in as source accounts' recorders deliver
  (sources without a grant show as pending authorization).
- Query it in the Config console's aggregator view or
  `aws configservice select-aggregate-resource-config`.
- Destroying deletes the aggregator/grants; source data is untouched.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
