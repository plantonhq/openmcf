# AWS Organization

The umbrella for multi-account AWS: creating it makes the deploying
account the management account, and everything org-wide hangs off it —
which AWS services get trusted access, which policy types are enabled,
which member accounts administer which services, and the org's own
delegation policy.

## What Gets Managed

- The organization's feature level (ALL, or consolidated billing
  only).
- Trusted service access per AWS service (CloudTrail, Config, Account
  Management, ...).
- Enabled policy types (SCPs and their twelve 2026 siblings).
- Delegated administrator registrations (a security account running
  GuardDuty org-wide, ...).
- The organization's single resource-based delegation policy.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Organizations permissions.

### AWS Account

- The account must NOT already belong to another organization —
  creating this resource makes it a MANAGEMENT account.
- Advanced arms (service access, policy types, delegated admins, the
  resource policy) require the ALL feature set — validation enforces
  it.

## Deploy

### Console

Create the resource from the AWS catalog, keep the ALL feature set,
name the service principals and policy types you need, and deploy.

### CLI

```bash
planton apply -f organization.yaml
```

## After Deploy

- The `root_id` output is the OU tree's top — wire
  [AWS Organizational Unit](/cloud-catalog/aws-organizational-unit)
  parents and root-scoped
  [AWS Organization Policy](/cloud-catalog/aws-organization-policy)
  attachments to it.
- Organizations, OUs, and policies are free.
- Deleting this resource deletes the ENTIRE organization — members,
  OUs, and policies must be removed first (AWS refuses otherwise).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
