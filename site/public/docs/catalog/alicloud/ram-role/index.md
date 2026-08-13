---
title: "RAM Role"
description: "RAM Role deployment documentation"
icon: "package"
order: 100
componentName: "alicloudramrole"
---

# AliCloud RAM Role

Deploys an Alibaba Cloud RAM role with bundled policy attachments and a configurable trust policy document. The component provisions the role and its policy attachments as a single atomic unit, ensuring the role is always created with its intended permissions. The role ARN is consumed by other AliCloud components (Function Compute, ECS, ACK) via ValueFromRef.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **RAM Role** -- an `alicloud_ram_role` with the specified trust policy, session duration, force-deletion behavior, and tags
- **Policy Attachments** -- one `alicloud_ram_role_policy_attachment` per entry in `policyAttachments`, granting the role permissions defined by system-managed or custom policies
- **AliCloud Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with user-provided `tags`

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- A valid JSON trust policy document specifying which principals (services, accounts, or federated identities) can assume the role via STS.
- Custom policies referenced with `policyType: Custom` must already exist -- create them with AliCloudRamPolicy first.
- System-managed policies (e.g., `AliyunOSSFullAccess`, `AliyunLogFullAccess`) are provided by Alibaba Cloud and do not need to be created.

## Deploy

### Console

Open the deployment store, find **AliCloud RAM Role**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **ECS Service Role** preset in the [Presets](#presets) tab to pre-populate a role for ECS instances with OSS, CloudMonitor, and Log Service access.

### CLI

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudRamRole
metadata:
  name: fc-execution-role
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  roleName: fc-execution-role
  assumeRolePolicyDocument: |
    {
      "Statement": [{
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {"Service": ["fc.aliyuncs.com"]}
      }],
      "Version": "1"
    }
  policyAttachments:
    - policyName: AliyunLogFullAccess
```

```shell
planton apply -f ram-role.yaml
```

This creates a RAM role that Function Compute can assume, with Log Service full access for writing invocation logs. Session duration defaults to 3600 seconds (1 hour) and force-deletion is disabled.

## Key Configuration

These are the most important decisions when configuring a RAM role. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Trust policy document** -- The `assumeRolePolicyDocument` field defines which principals can call `sts:AssumeRole`. For service roles, specify a Service principal (e.g., `fc.aliyuncs.com` for Function Compute, `ecs.aliyuncs.com` for ECS). For cross-account access, specify a RAM principal with the trusted account ID (e.g., `acs:ram::1234567890123456:root`).

**Policy attachments** -- Each entry in `policyAttachments` creates a policy attachment resource. Set `policyType` to `System` (default) for Alibaba Cloud managed policies or `Custom` for policies created with AliCloudRamPolicy. A role without policies can authenticate via STS but has zero permissions.

**Session duration** -- Set `maxSessionDuration` between 3600 (1 hour, default) and 43200 (12 hours). Longer durations suit CI/CD pipelines and batch workloads that should not need to re-assume the role mid-execution.

**Force deletion** -- Set `force: true` to allow deleting the role even when policies are still attached. When `false` (default), deletion fails if any policies remain, preventing accidental permission removal.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `role_id` | The RAM role ID assigned by Alibaba Cloud | Internal reference, audit logs |
| `role_name` | The RAM role name as created | Direct role name references in other configurations |
| `arn` | The role ARN (`acs:ram::<account-id>:role/<role-name>`) | AliCloudFunction `role`, ECS instance profiles, ACK service authentication |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**ECS service role** -- A role for ECS instances with OSS, CloudMonitor, and Log Service access. Also serves as a starting point for ACK worker node roles. Start from the **ECS Service Role** preset.

**FC execution role** -- A minimal role for Function Compute functions with Log Service access for invocation logging. Add service-specific policies based on what the function accesses. Start from the **FC Execution Role** preset.

**Cross-account audit role** -- A role that another Alibaba Cloud account can assume for read-only security auditing with billing, log, and ActionTrail access. Uses a 12-hour session duration and force-deletion. Start from the **Cross-Account Audit** preset.

## Works With

This component operates independently and does not reference other components.