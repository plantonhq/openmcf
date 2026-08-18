# AWS Planton Runner

Deploys a standing Planton runner appliance inside your AWS VPC -- an always-on, outbound-only worker that receives deploy operations from the control plane and executes them from within the network. It is the piece that makes private endpoints (most notably private EKS cluster APIs) deployable and operable, with subnets, security groups, and the runtime role wired through ValueFromRef.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Secrets Manager secret** -- holds the runner's credentials document; the container reads it at start through native secret injection, so the credentials never appear in the task definition
- **IAM execution role** -- the setup identity: pulls the runner image, writes logs, and reads exactly the one credentials secret
- **IAM runtime role** -- the runner's own AWS identity while executing work; created permissionless, and only when `taskRole` does not reference an existing role
- **CloudWatch log group** -- the audit trail of every operation the runner executes, with configurable retention
- **Security group** -- outbound-only, no inbound rules, created in the VPC derived from the referenced subnets
- **ECS cluster** -- dedicated to this runner, keeping the appliance's blast radius and teardown self-contained
- **ECS task definition and Fargate service** -- the runner container itself, kept at exactly one copy and restarted automatically on exit
- **AWS Tags** -- resource tags derived from the resource's organization, environment, and name

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- the credential used to provision the appliance itself. Required.
- **Runner registration and credentials** -- nothing to create by hand: the platform enrolls the appliance itself, creating the runner registration and minting its identity document into the managed secret the manifest references. Choose a secret slug and reference it as `$secret/<slug>` in the `credentials` field. Plaintext is rejected.

### AWS Account

- **Subnets** -- the VPC placement decision. Prefer private subnets with a NAT route (the runner needs outbound internet for its image and the control plane) across two Availability Zones; public subnets work with `assignPublicIp: true`. Reference `AwsSubnet` resources via ValueFromRef or pass literal subnet IDs.
- **Security groups (optional)** -- only when a private target (a cluster API endpoint, a database) admits traffic by source security group; attach the group that target trusts.
- **IAM role (optional)** -- a role trusted by `ecs-tasks.amazonaws.com` when the runner needs AWS permissions of its own for keyless cloud access.

## Deploy

### Console

Open the deployment store, find **AWS Planton Runner**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **private-vpc-worker** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsPlantonRunner
metadata:
  name: vpc-runner
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  subnets:
    - valueFrom:
        kind: AwsSubnet
        name: private-a
        fieldPath: status.outputs.subnet_id
    - valueFrom:
        kind: AwsSubnet
        name: private-b
        fieldPath: status.outputs.subnet_id
  credentials: $secret/vpc-runner-credentials
```

```shell
planton apply -f runner.yaml
```

This minimal manifest deploys a pull-based (temporal-mode) worker at the default sizing (0.5 vCPU, 1 GiB) tracking the latest runner release, with a permissionless runtime role and 30-day log retention -- sizing, version pinning, execution mode, and the runtime identity are not configured.

### InfraChart

```yaml
spec:
  subnets:
    - valueFrom:
        kind: AwsSubnet
        name: private-a
        fieldPath: status.outputs.subnet_id
  securityGroups:
    - valueFrom:
        kind: AwsSecurityGroup
        name: eks-api-clients
        fieldPath: status.outputs.security_group_id
  taskRole:
    valueFrom:
      kind: AwsIamRole
      name: runner-runtime
      fieldPath: status.outputs.role_arn
```

The InfraPipeline resolves the dependency graph, deploys the subnets, security groups, and IAM role first, then provisions the runner with the resolved values.

## Key Configuration

These are the most important decisions when configuring the runner. Explore the full field reference in the [API Explorer](#api-explorer) tab.

- **Subnet placement** -- `subnets` is the whole point of the appliance: whatever those subnets' route tables can reach, the runner can deploy to. All subnets must belong to the same VPC; two Availability Zones let the service reschedule across an AZ event.
- **Execution mode** -- `executionMode` defaults to `temporal`, a pull-based worker that polls its queue and needs no inbound path at all. `dual` adds the real-time CloudOps channel through an outbound-initiated tunnel (the credentials must carry tunnel material); `grpc` executes no deploy operations and is rarely right for a standing appliance.
- **Sizing** -- `cpu` and `memory` must form one of the fixed serverless pairings (validated up front; AWS would otherwise reject them only at deploy time). The 512/1024 default handles typical IaC operations; memory pressure shows up as failed operations mid-apply, so size memory up before CPU.
- **Runner build** -- empty `runnerVersion` tracks the newest release on every task (re)start; pin a version tag for change control. `imageRepository` is only for air-gapped or mirrored registries hosting a digest-identical copy.
- **Runtime identity** -- leave `taskRole` empty to get a permissionless role (the identity seam always exists, so permissions can be granted later without replacing the runner), or reference an `AwsIamRole` composed with exactly the permissions keyless cloud access needs.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AwsSubnet** | `subnets` | `status.outputs.subnet_id` |
| **AwsSecurityGroup** (optional) | `securityGroups` | `status.outputs.security_group_id` |
| **AwsIamRole** (optional) | `taskRole` | `status.outputs.role_arn` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `service_arn` | The ECS service keeping the runner running | Inspecting the appliance with AWS tooling |
| `service_name` | The service's name | Console and CLI lookups |
| `cluster_arn` | The dedicated ECS cluster's ARN | Scoping AWS operations to the appliance |
| `task_definition_arn` | The running family:revision | Tracking configuration/version rollouts |
| `log_group_name` | The CloudWatch log group | Tailing the operation audit trail (`aws logs tail <name> --follow`) |
| `security_group_id` | The outbound-only group's id | Private targets reference it to trust the runner (e.g. an EKS API's allowed sources) |
| `execution_role_arn` | The setup identity | Auditing image-pull/log/secret access |
| `task_role_arn` | The runtime identity | Granting the runner AWS permissions for keyless operations |
| `credentials_secret_arn` | The credentials secret's ARN | Auditing secret access; rotation tooling |
| `region` | The deployed region | Targeting follow-up AWS operations correctly |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

- **Private VPC worker** -- the standard appliance: a pull-based worker on two private subnets that makes a private-endpoint cluster deployable. Start from the **private-vpc-worker** preset.
- **Dual-mode operations runner** -- adds the real-time CloudOps channel for live resource browsing through the runner, still outbound-only. Start from the **dual-mode** preset.
- **High-capacity deploy worker** -- pinned version and larger sizing for heavy Terraform/Pulumi stacks and concurrent operations. Start from the **high-capacity** preset.

## Works With

- [**AWS Subnet**](/cloud-catalog/aws-subnet) -- where the runner's network interfaces live; the placement that defines what it can reach
- [**AWS Security Group**](/cloud-catalog/aws-security-group) -- extra groups a private target trusts, attached beside the created outbound-only group
- [**AWS IAM Role**](/cloud-catalog/aws-iam-role) -- the runner's runtime identity, composed first-class with exactly the permissions its workloads need
