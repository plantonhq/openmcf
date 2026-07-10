# AWS SageMaker Domain

Deploys an Amazon SageMaker Domain — the top-level workspace for SageMaker Studio. The domain carries VPC networking, the per-user baseline (JupyterLab, classic Jupyter Server, KernelGateway, Code Editor, TensorBoard, RSession, RStudio, and Canvas app settings), the shared-space baseline, Docker access, idle-timeout cost policies, custom file-system mounts, POSIX identity, Studio UI governance, identity-propagation posture, and the retention decision for the domain's dedicated EFS home-directory file system.

## What Gets Created

When you deploy an AwsSagemakerDomain resource, Planton provisions:

- **SageMaker Domain** — a `sagemaker.Domain` resource placed in the specified VPC and subnets, with the configured authentication mode (IAM or SSO), user and space baselines, and domain-level settings (Docker, security groups, RStudio activation, trusted identity propagation)
- **Dedicated EFS File System** — automatically created by AWS for user home directories (the `home_efs_file_system_id` is exposed as a stack output; `homeEfsRetentionPolicy` decides whether it survives domain deletion)
- **Domain Boundary Security Group** — automatically created by AWS to control cross-app and cross-user traffic within the domain (the `security_group_id_for_domain_boundary` is exposed as a stack output)
- **IAM Identity Center Application** — created only when `authMode` is `SSO`, registers the domain as an SSO application for centralized identity management

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **A VPC** with DNS resolution and DNS hostnames enabled
- **At least one subnet** in the VPC for SageMaker network interfaces (private subnets recommended for production)
- **An IAM execution role** with a trust policy for `sagemaker.amazonaws.com`, granting access to S3, ECR, and other services the ML workloads need
- **AWS IAM Identity Center** configured in the account if using `SSO` authentication mode (also required for trusted identity propagation)
- **A Posit license via AWS License Manager** if activating RStudio Server Pro
- **A security group** allowing outbound traffic for notebook and training workloads

## Quick Start

Create a file `sagemaker-domain.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsSagemakerDomain
metadata:
  name: my-domain
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsSagemakerDomain.my-domain
spec:
  region: us-east-1
  authMode: IAM
  vpcId:
    value: vpc-0123456789abcdef0
  subnetIds:
    - value: subnet-0a1b2c3d4e5f00001
    - value: subnet-0a1b2c3d4e5f00002
  defaultUserSettings:
    executionRoleArn:
      value: arn:aws:iam::123456789012:role/SageMakerExecutionRole
```

Deploy:

```shell
planton apply -f sagemaker-domain.yaml
```

This creates a SageMaker Domain with IAM authentication, public internet access for notebooks, and the SageMaker-provided default JupyterLab environment.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | The AWS region where the resource will be created. | Required |
| `authMode` | `string` | Authentication mode for the domain. ForceNew. | Must be `IAM` or `SSO` |
| `vpcId` | `StringValueOrRef` | VPC where the domain is created. ForceNew. Can reference AwsVpc via `valueFrom`. | Required |
| `subnetIds` | `StringValueOrRef[]` | Subnets for SageMaker network interfaces. ForceNew. Can reference AwsSubnet via `valueFrom`. | 1–16 items |
| `defaultUserSettings.executionRoleArn` | `StringValueOrRef` | IAM role assumed by SageMaker for notebooks, training, and inference. Must trust `sagemaker.amazonaws.com`. Can reference AwsIamRole via `valueFrom`. | Required |

### Optional Fields (top level)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `kmsKeyId` | `StringValueOrRef` | AWS-managed key | KMS key for encrypting the domain's EFS volume. ForceNew. Can reference AwsKmsKey via `valueFrom`. |
| `appNetworkAccessType` | `string` | `PublicInternetOnly` | Network access for notebooks. `VpcOnly` keeps all traffic in the VPC (recommended for production). |
| `appSecurityGroupManagement` | `string` | — | `Service` or `Customer` — who manages app-ENI security groups. Only honored when RStudio is configured. |
| `tagPropagation` | `string` | `DISABLED` | `ENABLED` copies the domain's tags onto apps, spaces, and user profiles created inside it — recommended for cost allocation. |
| `homeEfsRetentionPolicy` | `string` | `Retain` | `Delete` removes the domain's auto-created EFS file system with the domain (ephemeral/test domains); `Retain` preserves it. ForceNew. |
| `domainSecurityGroupIds` | `StringValueOrRef[]` | `[]` | Domain-level security groups for shared resources. Maximum 3. ForceNew. Can reference AwsSecurityGroup via `valueFrom`. |
| `dockerSettings.enableDockerAccess` | `string` | — | Enable Docker in notebooks/terminals. Valid values: `ENABLED`, `DISABLED`. |
| `dockerSettings.vpcOnlyTrustedAccounts` | `string[]` | `[]` | AWS account IDs allowed for Docker image pulls in VpcOnly mode. Maximum 20. |
| `executionRoleIdentityConfig` | `string` | — | `USER_PROFILE_NAME` stamps each session's `sts:SourceIdentity` with the acting user profile for per-user CloudTrail attribution; `DISABLED` leaves sessions role-only. |
| `rStudioServerProDomainSettings` | `object` | — | Activates RStudio (Posit) Workbench: `domainExecutionRoleArn` (required, ref), `rStudioConnectUrl`, `rStudioPackageManagerUrl`, `defaultResourceSpec`. |
| `trustedIdentityPropagationStatus` | `string` | — | `ENABLED` forwards the Identity Center user identity to downstream analytics services. Requires `authMode: SSO`. |
| `defaultSpaceSettings` | `object` | — | Baseline inherited by shared spaces: own `executionRoleArn` (required, ref), `securityGroupIds`, JupyterLab/JupyterServer/KernelGateway settings, `spaceStorageSettings`, `customFileSystemConfigs`, `customPosixUserConfig`. |

### Optional Fields (defaultUserSettings)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `securityGroupIds` | `StringValueOrRef[]` | `[]` | User-level security groups for notebook ENIs. Maximum 5. Can reference AwsSecurityGroup via `valueFrom`. |
| `defaultLandingUri` | `string` | — | Default app opened on domain access. Common: `studio::relative/JupyterLab`. |
| `studioWebPortal` | `string` | `ENABLED` | Whether the Studio web portal is accessible. Valid values: `ENABLED`, `DISABLED`. |
| `autoMountHomeEfs` | `string` | — | `Enabled` auto-mounts each user's EFS home directory into apps; `Disabled` starts apps without it. |
| `jupyterLabAppSettings` | `object` | — | JupyterLab IDE: `defaultResourceSpec`, `lifecycleConfigArns`, `builtInLifecycleConfigArn`, `customImages` (max 200), `codeRepositories` (max 10), `idleSettings`, `emrSettings` (EMR connect/runtime role refs). |
| `jupyterServerAppSettings` | `object` | — | Classic Jupyter Server (Studio Classic): `defaultResourceSpec`, `lifecycleConfigArns`, `codeRepositories`. |
| `kernelGatewayAppSettings` | `object` | — | Custom kernels: `defaultResourceSpec`, `lifecycleConfigArns`, `customImages` (max 200). |
| `codeEditorAppSettings` | `object` | — | Code Editor (VS Code): `defaultResourceSpec`, `lifecycleConfigArns`, `builtInLifecycleConfigArn`, `customImages`, `idleSettings`. |
| `tensorBoardAppSettings` | `object` | — | TensorBoard: `defaultResourceSpec`. |
| `rSessionAppSettings` | `object` | — | RSession (R kernels): `defaultResourceSpec`, `customImages`. Requires RStudio on the domain. |
| `rStudioServerProAppSettings` | `object` | — | Per-user RStudio access: `accessStatus` (`ENABLED`/`DISABLED`), `userGroup` (`R_STUDIO_ADMIN`/`R_STUDIO_USER`, only with `accessStatus: ENABLED`). |
| `canvasAppSettings` | `object` | — | Canvas no-code ML: `directDeployStatus`, `emrServerlessSettings`, `generativeAiBedrockRoleArn` (ref), `identityProviderOauthSettings` (max 20; Secrets Manager `secretArn` per connector), `kendraSettingsStatus`, `modelRegisterSettings`, `timeSeriesForecastingSettings`, `workspaceSettings`. |
| `sharingSettings` | `object` | — | Notebook output sharing: `notebookOutputOption` (`Allowed`/`Disabled`), `s3OutputPath` (required when Allowed), `s3KmsKeyId` (ref). |
| `spaceStorageSettings` | `object` | — | `defaultEbsVolumeSizeInGb` + `maximumEbsVolumeSizeInGb` (max ≥ default). |
| `customFileSystemConfigs` | `object[]` | `[]` | Additional file-system mounts: each entry's `efsFileSystemConfig` carries `fileSystemId` (ref to AwsElasticFileSystem) + `fileSystemPath`. |
| `customPosixUserConfig` | `object` | — | POSIX identity for file-system operations: `uid` (≥ 10000) + `gid` (≥ 1001). |
| `studioWebPortalSettings` | `object` | — | Hide Studio UI surface: `hiddenAppTypes`, `hiddenInstanceTypes`, `hiddenMlTools`. |

## Examples

### IAM Authentication with JupyterLab Idle Shutdown

Cost-optimized domain with automatic shutdown of idle JupyterLab instances after 2 hours, and the EFS home file system deleted with the domain:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsSagemakerDomain
metadata:
  name: ml-team
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsSagemakerDomain.ml-team
spec:
  region: us-east-1
  authMode: IAM
  homeEfsRetentionPolicy: Delete
  vpcId:
    value: vpc-0123456789abcdef0
  subnetIds:
    - value: subnet-0a1b2c3d4e5f00001
    - value: subnet-0a1b2c3d4e5f00002
  defaultUserSettings:
    executionRoleArn:
      value: arn:aws:iam::123456789012:role/SageMakerExecutionRole
    defaultLandingUri: "studio::relative/JupyterLab"
    jupyterLabAppSettings:
      defaultResourceSpec:
        instanceType: ml.t3.medium
      idleSettings:
        idleTimeoutInMinutes: 120
        minIdleTimeoutInMinutes: 60
        maxIdleTimeoutInMinutes: 480
```

### SSO Authentication with VPC-Only Networking and Identity Propagation

Enterprise domain using IAM Identity Center, VPC-only access, per-user CloudTrail attribution, and trusted identity propagation into downstream analytics services:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsSagemakerDomain
metadata:
  name: enterprise-ml
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsSagemakerDomain.enterprise-ml
spec:
  region: us-east-1
  authMode: SSO
  appNetworkAccessType: VpcOnly
  tagPropagation: ENABLED
  executionRoleIdentityConfig: USER_PROFILE_NAME
  trustedIdentityPropagationStatus: ENABLED
  vpcId:
    value: vpc-0123456789abcdef0
  subnetIds:
    - value: subnet-private-az1
    - value: subnet-private-az2
    - value: subnet-private-az3
  kmsKeyId:
    value: arn:aws:kms:us-east-1:123456789012:key/mrk-abc123
  defaultUserSettings:
    executionRoleArn:
      value: arn:aws:iam::123456789012:role/SageMakerExecutionRole
    securityGroupIds:
      - value: sg-user-notebooks
    jupyterLabAppSettings:
      defaultResourceSpec:
        instanceType: ml.m5.large
      idleSettings:
        idleTimeoutInMinutes: 120
        minIdleTimeoutInMinutes: 60
        maxIdleTimeoutInMinutes: 480
    studioWebPortalSettings:
      hiddenInstanceTypes:
        - ml.p3.16xlarge
  domainSecurityGroupIds:
    - value: sg-domain-boundary
```

### Governed Canvas for Business Teams

Canvas no-code ML with direct endpoint deployment disabled (models go through the registry), forecasting enabled, and a pinned artifact workspace:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsSagemakerDomain
metadata:
  name: canvas-analytics
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsSagemakerDomain.canvas-analytics
spec:
  region: us-east-1
  authMode: IAM
  vpcId:
    value: vpc-0123456789abcdef0
  subnetIds:
    - value: subnet-private-az1
    - value: subnet-private-az2
  defaultUserSettings:
    executionRoleArn:
      value: arn:aws:iam::123456789012:role/SageMakerExecutionRole
    canvasAppSettings:
      directDeployStatus: DISABLED
      modelRegisterSettings:
        status: ENABLED
      timeSeriesForecastingSettings:
        amazonForecastRoleArn:
          value: arn:aws:iam::123456789012:role/CanvasForecastRole
        status: ENABLED
      workspaceSettings:
        s3ArtifactPath: s3://canvas-workspace/artifacts/
```

### Using Foreign Key References

Reference Planton-managed VPC, subnets, IAM role, security groups, and KMS key instead of hardcoding IDs:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsSagemakerDomain
metadata:
  name: ref-domain
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsSagemakerDomain.ref-domain
spec:
  region: us-east-1
  authMode: IAM
  vpcId:
    valueFrom:
      kind: AwsVpc
      name: ml-vpc
      fieldPath: status.outputs.vpc_id
  subnetIds:
    - valueFrom:
        kind: AwsSubnet
        name: ml-private-subnet-a
        fieldPath: status.outputs.subnet_id
    - valueFrom:
        kind: AwsSubnet
        name: ml-private-subnet-b
        fieldPath: status.outputs.subnet_id
  kmsKeyId:
    valueFrom:
      kind: AwsKmsKey
      name: ml-key
      fieldPath: status.outputs.key_arn
  defaultUserSettings:
    executionRoleArn:
      valueFrom:
        kind: AwsIamRole
        name: sagemaker-exec
        fieldPath: status.outputs.role_arn
    securityGroupIds:
      - valueFrom:
          kind: AwsSecurityGroup
          name: notebook-sg
          fieldPath: status.outputs.security_group_id
    customFileSystemConfigs:
      - efsFileSystemConfig:
          fileSystemId:
            valueFrom:
              kind: AwsElasticFileSystem
              name: shared-datasets
              fieldPath: status.outputs.file_system_id
          fileSystemPath: /shared/datasets
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `domain_id` | `string` | Unique identifier of the SageMaker Domain, used when creating user profiles and spaces |
| `domain_arn` | `string` | Amazon Resource Name for IAM policies and CloudWatch metrics |
| `domain_url` | `string` | HTTPS URL for accessing the SageMaker Studio web interface |
| `home_efs_file_system_id` | `string` | ID of the EFS file system created for user home directories |
| `security_group_id_for_domain_boundary` | `string` | ID of the AWS-managed security group controlling cross-app traffic |
| `single_sign_on_application_arn` | `string` | IAM Identity Center application ARN. Only populated when `authMode` is `SSO`. |
| `single_sign_on_managed_application_instance_id` | `string` | Identity Center managed application instance ID. Only populated when `authMode` is `SSO`. |

## Related Components

- [AwsVpc](/docs/catalog/aws/awsvpc) — provides the VPC and subnets for domain networking
- [AwsIamRole](/docs/catalog/aws/awsiamrole) — execution roles for the user plane, space plane, RStudio, Canvas capabilities, and EMR access
- [AwsSecurityGroup](/docs/catalog/aws/awssecuritygroup) — controls network access for domain and user-level resources
- [AwsKmsKey](/docs/catalog/aws/awskmskey) — customer-managed encryption for the EFS volume, shared outputs, and Canvas artifacts
- [AwsElasticFileSystem](/docs/catalog/aws/awselasticfilesystem) — additional shared file systems mounted into apps
- [AwsS3Bucket](/docs/catalog/aws/awss3bucket) — stores notebook outputs and Canvas artifacts
