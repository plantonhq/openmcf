# AWS Organization Account

A member account of the organization, born in the right OU with the
right contacts on file: the account itself, who AWS calls about
billing/operations/security, the postal address on record, and which
opt-in regions are enabled.

## What Gets Managed

- Account creation (name, root email, the pre-created bootstrap role,
  billing-info access for IAM users).
- Placement in the OU tree — and moves between OUs later.
- Billing, operations, and security alternate contacts.
- Primary contact information.
- Opt-in region enablement (Jakarta, Zurich, ...).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Organizations + Account Management
  permissions, on the organization's MANAGEMENT account.

### AWS Account

- An organization must exist
  ([AWS Organization](/cloud-catalog/aws-organization)).
- The contact and region arms additionally need trusted access for
  Account Management — `account.amazonaws.com` in the organization's
  service-access list.
- The root email must not belong to any existing AWS account.

## Deploy

### Console

Create the resource from the AWS catalog, name the account, set the
root email, pick the OU, fill the contacts, and deploy.

### CLI

```bash
planton apply -f organization-account.yaml
```

## After Deploy

- Account creation runs asynchronously (usually under a minute; the
  module polls to completion). Moves, renames, and contact updates
  apply in place; region enablement takes up to ~60 minutes.
- The account itself is free — what runs inside it is not.
- **Destroy is never a clean delete**: removal leaves a live
  standalone account (it needs its own billing info to leave);
  closure suspends for ~90 days and is quota-limited per rolling 30
  days.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
