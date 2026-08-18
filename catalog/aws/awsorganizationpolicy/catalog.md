# AWS Organization Policy

The guardrail document for the account tree: a service control policy
(or tag, backup, AI-services-opt-out, and nine more 2026 types)
written once and attached to the root, OUs, or individual accounts —
everything beneath inherits it.

## What Gets Managed

- The policy document (structured JSON in the type's own syntax).
- Its type (SCP by default; RCP, tag, declarative EC2, SecurityHub,
  Bedrock, ...).
- Attachments to any mix of the root, OUs, and member accounts.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Organizations permissions, on the
  organization's MANAGEMENT account.

### AWS Account

- An organization with ALL features
  ([AWS Organization](/cloud-catalog/aws-organization)), with this
  policy's type in its `enabledPolicyTypes`.

## Deploy

### Console

Create the resource from the AWS catalog, write the document, pick the
targets, and deploy.

### CLI

```bash
planton apply -f organization-policy.yaml
```

## After Deploy

- Effects are inheritance-scoped: attached to the root it governs
  everything; attached to an OU it governs that subtree; attached to
  an account, just that account.
- SCPs never grant — they cap what IAM can allow. AWS's own
  FullAWSAccess stays attached unless you deliberately replace it.
- Policies and attachments are free.
- Destroy detaches every target, then deletes the policy.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
