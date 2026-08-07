# AliCloud RAM Policy

Deploys a custom RAM policy on Alibaba Cloud with a JSON permission document, optional version rotation, and tag management. Custom policies define fine-grained permissions beyond what system-managed policies provide and can be attached to RAM roles via AliCloudRamRole's `policyAttachments` field.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **RAM Policy** -- an `alicloud_ram_policy` with the specified JSON policy document, version rotation strategy, and force-deletion behavior
- **AliCloud Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with user-provided `tags`

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- A valid JSON policy document conforming to the Alibaba Cloud RAM policy structure (Version, Statement, Effect, Action, Resource). Maximum 6144 bytes.
- Custom policies referenced by `policyType: Custom` in AliCloudRamRole `policyAttachments` must exist before the role is deployed.

## Deploy

### Console

Open the deployment store, find **AliCloud RAM Policy**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Scoped OSS Access** preset in the [Presets](#presets) tab to pre-populate a bucket-scoped read/write policy.

### CLI

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudRamPolicy
metadata:
  name: app-data-bucket-rw
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  policyName: app-data-bucket-rw
  policyDocument: |
    {
      "Version": "1",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": ["oss:GetObject", "oss:PutObject", "oss:DeleteObject", "oss:ListObjects"],
          "Resource": ["acs:oss:*:*:app-data-prod", "acs:oss:*:*:app-data-prod/*"]
        }
      ]
    }
```

```shell
planton apply -f ram-policy.yaml
```

This creates a custom RAM policy granting read/write access to a single OSS bucket. Version rotation and force-deletion are not configured.

## Key Configuration

These are the most important decisions when configuring a RAM policy. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Policy document structure** -- The `policyDocument` field accepts a JSON string with Alibaba Cloud's IAM policy format. Each statement specifies an Effect (Allow/Deny), a list of Actions (API operations), and Resource ARNs. Scope resources as narrowly as possible -- `acs:oss:*:*:my-bucket/*` is safer than `acs:oss:*:*:*`.

**Version rotation strategy** -- Alibaba Cloud allows up to 5 versions per policy. Set `rotateStrategy` to `DeleteOldestNonDefaultVersionWhenLimitExceeded` for policies that are updated frequently (e.g., CI/CD pipeline policies). Leave it at the default `None` for stable policies that rarely change -- updates will fail at the 5-version limit, which is a useful guard against unintended policy drift.

**Force deletion** -- Set `force: true` to allow deleting the policy even when it is still attached to roles, users, or groups. When `false` (default), deletion fails if the policy is attached anywhere, preventing accidental permission loss.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `policy_name` | The RAM policy name as created | AliCloudRamRole `policyAttachments` with `policyType: Custom` |
| `policy_type` | Always `Custom` for user-created policies | AliCloudRamRole `policyAttachments` require both name and type |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Scoped OSS bucket access** -- A policy granting read/write access to a single OSS bucket with automatic version rotation. Start from the **Scoped OSS Access** preset.

**CI/CD deploy pipeline** -- A cross-service policy combining Container Registry image push, ACK cluster deployment, and SLS log writing permissions for CI/CD pipelines. Start from the **CI/CD Deploy Pipeline** preset.

## Works With

This component operates independently and does not reference other deployment components.