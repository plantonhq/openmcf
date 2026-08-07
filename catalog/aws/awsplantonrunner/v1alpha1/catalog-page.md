# AWS Planton Runner

Runs a standing Planton runner appliance inside your AWS network on ECS
Fargate: an always-on, outbound-only worker that executes deploy
operations and cloud operations from within the VPC -- the piece that
makes private endpoints (most notably private Kubernetes cluster APIs)
deployable and operable.

## What Gets Created

When you deploy an AwsPlantonRunner resource, Planton provisions:

- **Secrets Manager secret** — holds the runner's credentials document;
  the container reads it at start through AWS's native secret injection,
  so the credentials never appear in the task definition
- **ECS task execution role** — the setup identity: pulls the runner
  image, writes logs, and reads exactly the one credentials secret
- **Runtime IAM role** — the runner's own AWS identity while executing
  work (skipped when you reference your own via `taskRole`)
- **CloudWatch log group** — the audit trail of every operation the
  runner executes, with configurable retention
- **Security group** — outbound-only, no inbound rules, created in the
  VPC derived from the referenced subnets
- **ECS cluster** — a dedicated, per-runner cluster (clusters are free
  scheduling boundaries; a dedicated one keeps the appliance's blast
  radius and teardown self-contained)
- **ECS task definition + Fargate service** — the runner container
  itself, kept at exactly one copy, restarted automatically on exit

The subnets (and any extra security groups or the runtime role) are
referenced resources -- the appliance never creates or mutates them.

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless SSO/OIDC).
- **A runner registration and its credentials document** — created with
  `planton runner generate-credentials <runner-name>` and stored as a
  managed secret the manifest references.
- **Subnets** (`AwsSubnet`) in the VPC whose private endpoints the runner
  must reach — private subnets with a NAT route recommended (the runner
  needs outbound internet for its image and the control plane); public
  subnets require `assignPublicIp: true`.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsPlantonRunner
metadata:
  name: vpc-runner
spec:
  region: us-west-2
  subnets:
    - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
    - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
  credentials: $secret/vpc-runner-credentials
```

```shell
planton apply -f runner.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
| --- | --- | --- | --- |
| `region` | `string` | AWS region; deploy the runner beside the private endpoints it must reach. | Required; non-empty |
| `subnets` | `(string \| valueFrom)[]` | The subnets the runner's network interfaces are placed in; all must belong to one VPC. Two AZs recommended. | Required |
| `credentials` | `string` | The runner's credentials document (from `planton runner generate-credentials`), supplied as a managed-secret reference. | Required; secret |

### Optional Fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `securityGroups` | `(string \| valueFrom)[]` | `[]` | Extra groups beside the created outbound-only one — attach the group a private target trusts. |
| `assignPublicIp` | `bool` | `false` | Assign public IPv4 to the runner — required in public subnets without a NAT gateway. |
| `cpu` | `int32` | `512` | CPU units (1024 = 1 vCPU); one of 256–16384 (Fargate sizes). |
| `memory` | `int32` | `1024` | Memory in MiB; must form a valid pairing with `cpu` (validated up front). |
| `runnerVersion` | `string` | `latest` | The runner image tag to run; pin for change control. |
| `imageRepository` | `string` | `ghcr.io/plantonhq/planton/runner` | Override only for air-gapped/mirrored registries. |
| `executionMode` | `string` | `temporal` | `temporal` (pull-based deploy worker), `dual` (adds the real-time CloudOps channel via outbound tunnel), or `grpc` (CloudOps only). |
| `taskRole` | `string \| valueFrom` | created | The runner's runtime IAM role — reference an `AwsIamRole` when the runner needs AWS permissions of its own. |
| `logRetentionDays` | `int32` | `30` | CloudWatch retention for the runner's operation logs. |

## Examples

### Runner for a private EKS cluster (recommended posture)

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsPlantonRunner
metadata:
  name: prod-vpc-runner
spec:
  region: us-west-2
  subnets:
    - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
    - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
  credentials: $secret/prod-vpc-runner-credentials
  cpu: 1024
  memory: 4096
  runnerVersion: v0.4.0
```

The private cluster's API endpoint admits traffic from the runner's
security group -- reference this runner's `security_group_id` output from
the cluster's allowed sources.

### Dual-mode runner with its own AWS identity

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsPlantonRunner
metadata:
  name: ops-runner
spec:
  region: us-west-2
  subnets:
    - valueFrom: { kind: AwsSubnet, name: private-a, fieldPath: status.outputs.subnet_id }
    - valueFrom: { kind: AwsSubnet, name: private-b, fieldPath: status.outputs.subnet_id }
  credentials: $secret/ops-runner-credentials
  executionMode: dual
  taskRole:
    valueFrom: { kind: AwsIamRole, name: ops-runner-runtime, fieldPath: status.outputs.role_arn }
```

## Stack Outputs

| Output | Description |
| --- | --- |
| `service_arn` | The ECS service keeping the runner running |
| `service_name` | The service's name (metadata.name) |
| `cluster_arn` | The dedicated ECS cluster's ARN |
| `task_definition_arn` | The running task-definition revision |
| `log_group_name` | The CloudWatch log group — the runner's operation audit trail |
| `security_group_id` | The outbound-only group; private targets reference it to trust the runner |
| `execution_role_arn` | The setup identity (image pull, logs, secret read) |
| `task_role_arn` | The runner's runtime identity — grant it permissions for keyless cloud access |
| `credentials_secret_arn` | The Secrets Manager secret holding the credentials document |
| `region` | The region the runner was deployed in |

## Related Components

- [AwsSubnet](/docs/catalog/aws/awssubnet) — where the runner's network interfaces live
- [AwsSecurityGroup](/docs/catalog/aws/awssecuritygroup) — extra groups a private target trusts
- [AwsIamRole](/docs/catalog/aws/awsiamrole) — the runner's runtime identity, composed first-class
- [AwsEksCluster](/docs/catalog/aws/awsekscluster) — the private-endpoint cluster the runner typically serves
- [AwsVpc](/docs/catalog/aws/awsvpc) — the network the runner is placed into
