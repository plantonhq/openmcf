# AWS Organizational Unit

A folder in the organization's account tree: group member accounts by
function (Workloads, Sandbox, Security), attach guardrail policies to
the group, and let accounts inherit them.

## What Gets Managed

- The OU itself, under the organization root or nested under another
  OU.
- Its display name (spaces are fine — renames apply in place).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Organizations permissions, on the
  organization's MANAGEMENT account.

### AWS Account

- An organization must exist — deploy
  [AWS Organization](/cloud-catalog/aws-organization) first (its
  `root_id` output is the default parent wiring).

## Deploy

### Console

Create the resource from the AWS catalog, name the OU, pick the parent
(the organization's root by default), and deploy.

### CLI

```bash
planton apply -f organizational-unit.yaml
```

## After Deploy

- Place accounts in the OU with
  [AWS Organization Account](/cloud-catalog/aws-organization-account)
  (`parentId` → this OU's `ou_id` output).
- Attach guardrails with
  [AWS Organization Policy](/cloud-catalog/aws-organization-policy)
  (an attachment targeting `ou_id`).
- OUs are free; deleting one requires it to be empty (AWS refuses
  otherwise).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
