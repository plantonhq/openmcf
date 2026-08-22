# AWS IAM Group

The IAM container that gives a set of users one shared permission set
— membership declared as a list, permissions as managed-policy
attachments and inline documents.

## What Gets Managed

- The group and its IAM path.
- Declarative membership: the users list is authoritative — additions
  made outside it are removed on the next apply.
- Permissions: managed-policy attachments (yours or AWS-managed) and
  inline policies unique to this group.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with IAM permissions.

### AWS Account

- The users in the membership list must already exist
  ([AWS IAM User](/cloud-catalog/aws-iam-user)); custom managed
  policies come from [AWS IAM Policy](/cloud-catalog/aws-iam-policy).

## Deploy

### Console

Create the resource from the AWS catalog, add users and policies, and
deploy.

### CLI

```bash
planton apply -f group.yaml
```

## After Deploy

- Users in the group inherit its policies immediately.
- Outputs publish the group's ARN, name, and stable unique ID.
