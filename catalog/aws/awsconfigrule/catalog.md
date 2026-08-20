# AWS Config Rule

One compliance check over the region's recorded configurations —
AWS-managed, custom Lambda, or CloudFormation-Guard policy — with
optional organization-wide deployment and SSM auto-remediation.

## What Gets Managed

- The rule: one of an AWS-managed identifier, a Lambda evaluator, or
  a Guard policy.
- Scope (resource types, one pinned resource, or a tag) and
  evaluation modes (detective / proactive).
- Organization-wide deployment with per-account exclusions.
- SSM remediation, manual or automatic with a retry contract.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AWS Config permissions.

### AWS Account

- A RUNNING configuration recorder in the region
  ([AWS Config Recorder](/cloud-catalog/aws-config-recorder)) — AWS
  rejects rule creation without one.
- For `custom_lambda`: the function
  ([AWS Lambda](/cloud-catalog/aws-lambda)) with
  `config.amazonaws.com` invoke permission.
- For organization rules: the management or Config delegated-admin
  account.

## Deploy

### Console

Create the resource from the AWS catalog, pick the source arm, scope
it, and deploy.

### CLI

```bash
planton apply -f config-rule.yaml
```

## After Deploy

- Evaluations appear under the rule in the AWS Config console as
  recorded resources change (or on the periodic schedule).
- Non-compliant resources remediate per the SSM document when
  remediation is configured.
- Destroy deletes the rule and its remediation; evaluation history
  ages out with Config retention.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
