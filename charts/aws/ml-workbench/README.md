# AWS ML Workbench

SageMaker Studio ready for a team in one deploy: a domain with idle
auto-shutdown so forgotten notebooks stop billing, a least-privilege
execution role, a versioned artifacts bucket for datasets and models, and a
dedicated no-internet-path VPC. One toggle locks the whole workbench into
VPC-only mode — and ships the complete set of VPC endpoints Studio actually
needs to function without internet, because a locked-down domain whose apps
hang on their first AWS call is not a security posture, it is an outage.

This is the on-ramp for a data-science team: JupyterLab and Code Editor
spaces per user, shared home storage, training jobs and inference endpoints
launched from notebooks — with the two costs that bite (idle compute and
unbounded artifact versions) already governed.

## Architecture

```
                    [AwsSagemakerDomain]  auth: IAM
                      idle auto-shutdown (JupyterLab + Code Editor)
                      homeEfsRetentionPolicy: Retain (param)
                     /        |                    \
        default execution     |                     appNetworkAccessType:
        [AwsIamRole]          |                     PublicInternetOnly (default)
        SageMakerFullAccess   |                     -- or, vpc_only_enabled --
        + artifacts bucket    |                     VpcOnly + studio SG
        inline grants         |
                              v
   [AwsVpc] {{ 10.64.0.0/16 }}, DNS support + hostnames ON
     ├── [AwsSubnet] x2 (routeless -> VPC main route table, no IGW/NAT)
     └── vpc_only_enabled additionally renders:
           [AwsSecurityGroup] studio-sg      (inter-app 8192-65535 + NFS 2049, self)
           [AwsSecurityGroup] endpoints-sg   (443 from studio-sg only)
           [AwsVpcEndpoint] x7 Interface     (sagemaker.api, sagemaker.runtime,
                                              sts, logs, ecr.api, ecr.dkr,
                                              servicecatalog; private DNS on)
           [AwsVpcEndpoint] Gateway s3       (on the VPC main route table)

   [AwsS3Bucket] artifacts: versioned, SSE, private, version-pruning lifecycle
```

## Included Cloud Resources

| Resource | Kind | Purpose |
|---|---|---|
| Studio domain | `AwsSagemakerDomain` | The workbench: user profiles inherit the execution role, idle shutdown, and network posture. |
| Execution role | `AwsIamRole` | `AmazonSageMakerFullAccess` starter + explicit artifacts-bucket grants. |
| Artifacts bucket | `AwsS3Bucket` | Versioned, encrypted, private; noncurrent versions pruned on a schedule. |
| Workbench VPC + subnets | `AwsVpc`, `AwsSubnet` x2 | Dedicated, internet-path-free network across two AZs. |
| Studio + endpoint SGs | `AwsSecurityGroup` x2 | Inter-app/EFS contract and endpoint HTTPS scoping. VPC-only arm. |
| Service endpoints | `AwsVpcEndpoint` x8 | The full documented VpcOnly set (7 interface + S3 gateway). VPC-only arm. |

## Parameters

| Name | Description | Default | Required |
|---|---|---|---|
| `aws_region` | Region — keep it where the training data lives. | `us-east-1` | yes |
| `workbench_name` | Domain name and companion-resource prefix. | `ml-workbench` | yes |
| `artifacts_bucket_name` | Globally unique artifacts bucket name. | `my-org-ml-artifacts` | yes |
| `availability_zones` | Two AZs; entry N pairs with `subnet_cidrs` entry N. | `us-east-1a/b` | yes |
| `subnet_cidrs` | Subnet CIDRs (every VPC-only app consumes an address). | `10.64.1.0/24`, `10.64.2.0/24` | yes |
| `vpc_cidr` | The workbench VPC's CIDR — pick a non-overlapping range. | `10.64.0.0/16` | yes |
| `idle_shutdown_minutes` | Inactivity before apps auto-stop. | `120` | yes |
| `home_efs_retention` | `Retain` or `Delete` for home directories on domain delete. | `Retain` | yes |
| `vpc_only_enabled` | Lock traffic in-VPC and render the endpoint set. | `false` | no |
| `force_destroy` | Allow teardown of a non-empty artifacts bucket. | `false` | no |

## Onboarding the team (post-deploy)

The domain ships without user profiles — one per teammate is the unit of
identity, home directory, and CloudTrail attribution:

```bash
aws sagemaker create-user-profile \
  --domain-id <status.outputs.domain_id> \
  --user-profile-name ada
```

Each user opens Studio from the SageMaker console (or a presigned URL) and
lands in their own home directory with the chart's defaults applied. Point
notebooks at the artifacts bucket (`s3://<artifacts_bucket_name>/...`) —
the execution role already carries read/write on exactly that bucket.

## The two postures, honestly

- **PublicInternetOnly (default)**: notebooks reach PyPI, GitHub, and
  Docker Hub directly; outbound rides the SageMaker service network, so
  the chart's VPC needs no NAT and costs nothing. The friction-free shape
  for teams without a compliance driver.
- **VpcOnly (`vpc_only_enabled`)**: no direct internet from any app; every
  AWS call rides the chart's private endpoints; flows are visible to VPC
  Flow Logs; required for HIPAA/SOC2/PCI-class postures. Two honest costs:
  the seven interface endpoints bill per AZ-hour (~$60/month across two
  AZs), and `pip install` from public PyPI stops working — stand up a
  CodeArtifact repository (with its own endpoint) or bake dependencies
  into custom images pulled through the ECR endpoints.

Note AWS treats the domain's network mode as create-time: flipping the
toggle on an already-deployed workbench replaces the domain (home
directories survive per `home_efs_retention`). Choose the posture before
the team moves in.

## Day-2 guidance

- **Tighten the execution role.** `AmazonSageMakerFullAccess` is the
  supported starter, not the destination: once usage patterns settle,
  replace it with scoped grants (the role's `managedPolicyArns` accepts
  references to your own `AwsIamPolicy` nodes) and split per-team roles
  across user profiles.
- **Customer-managed keys.** For KMS-everything compliance, deploy an
  `AwsKmsKey` and set it as the domain's `kmsKeyId` (home EFS) and the
  bucket's `aws:kms` encryption — one key policy governing both stores.
- **Docker in Studio.** Container builds require VpcOnly mode plus the
  domain's `dockerSettings` (enable access, pin
  `vpcOnlyTrustedAccounts`) — the ECR endpoints this chart ships are the
  transport for those pulls.
- **Shared spaces.** Real-time collaboration uses the domain's
  `defaultSpaceSettings` (the spec models it); add it when the team wants
  shared JupyterLab spaces rather than per-user apps.
- **Cost visibility.** Idle shutdown caps runaway notebooks; for the rest,
  alarm on the `AWS/SageMaker` `InstanceMemoryUtilization`/CPU metrics per
  app and review the artifacts bucket's storage class after the first few
  training cycles — big cold datasets belong in the data-lakehouse chart's
  tiered lake, not in hot artifact storage.
