---
title: "CodeBuild Project"
description: "CodeBuild Project deployment documentation"
icon: "package"
order: 100
componentName: "awscodebuildproject"
---

# AWS CodeBuild Project

Deploys an AWS CodeBuild project with configurable source providers (primary plus up to 12 secondary sources), build environments (containers, Lambda compute, reserved fleets, persistent Docker servers), artifact outputs (primary plus up to 12 secondary), batch builds, EFS mounts, public visibility, a resource-based access policy, and an optional webhook for automatic build triggers.

## What Gets Created

When you deploy an AwsCodeBuildProject resource, Planton provisions:

- **CodeBuild Project** — an `aws_codebuild_project` resource with the specified sources, environment, artifacts, cache, logging, VPC/EFS placement, batch configuration, and visibility
- **Webhook** — created only when `webhook` is configured; registers a webhook with the source provider (or, with `manualCreation`, mints the payload URL + HMAC secret for hand-wiring)
- **Resource Policy** — created only when `resourcePolicy` is configured; attaches a resource-based IAM policy to the project for cross-account access

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **An IAM service role** granting CodeBuild permission to access source code, write artifacts, and publish logs
- **A source repository** accessible to CodeBuild (GitHub via CodeStar Connections, CodeCommit, Bitbucket, GitLab, or S3)
- **An S3 bucket** if using S3 artifacts or S3 cache
- **VPC, subnets, and security groups** if running builds inside a VPC
- **A KMS key** if custom encryption for build artifacts is required

## Quick Start

Create a file `codebuild.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCodeBuildProject
metadata:
  name: my-build
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsCodeBuildProject.my-build
spec:
  region: us-west-2
  source:
    type: GITHUB
    location: https://github.com/my-org/my-app.git
  environment:
    type: LINUX_CONTAINER
    computeType: BUILD_GENERAL1_SMALL
    image: aws/codebuild/amazonlinux2-x86_64-standard:5.0
  artifacts:
    type: NO_ARTIFACTS
  serviceRole:
    value: arn:aws:iam::123456789012:role/codebuild-service-role
```

Deploy:

```shell
planton apply -f codebuild.yaml
```

This creates a CodeBuild project that pulls from a GitHub repository, runs builds in a standard Linux container, and produces no stored artifacts (CI-only).

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region where the project will be created (e.g., `us-west-2`). | Required |
| `source.type` | `string` | Source provider: `GITHUB`, `BITBUCKET`, `CODECOMMIT`, `CODEPIPELINE`, `GITHUB_ENTERPRISE`, `GITLAB`, `GITLAB_SELF_MANAGED`, `NO_SOURCE`, `S3` | Must be one of the listed values |
| `source.location` | `string` | Repository URL or S3 path | Required unless type is `CODEPIPELINE` or `NO_SOURCE` |
| `environment.type` | `string` | Container type: `LINUX_CONTAINER`, `LINUX_GPU_CONTAINER`, `ARM_CONTAINER`, `WINDOWS_CONTAINER`, `WINDOWS_SERVER_2019_CONTAINER`, `WINDOWS_SERVER_2022_CONTAINER`, `LINUX_LAMBDA_CONTAINER`, `ARM_LAMBDA_CONTAINER`, `LINUX_EC2`, `ARM_EC2`, `WINDOWS_EC2`, `MAC_ARM` | Must be one of the listed values |
| `environment.computeType` | `string` | Compute capacity: `BUILD_GENERAL1_SMALL` through `BUILD_GENERAL1_2XLARGE`, `BUILD_LAMBDA_1GB` through `BUILD_LAMBDA_10GB`, or the fleet-driven `ATTRIBUTE_BASED_COMPUTE` / `CUSTOM_INSTANCE_TYPE` | Must be one of the listed values |
| `environment.image` | `string` | Docker image for the build environment (e.g., `aws/codebuild/amazonlinux2-x86_64-standard:5.0`) | Required, non-empty |
| `artifacts.type` | `string` | Artifact output type: `NO_ARTIFACTS`, `S3`, `CODEPIPELINE` | Must be one of the listed values. Must be `CODEPIPELINE` when source is `CODEPIPELINE` |
| `serviceRole` | `StringValueOrRef` | IAM role ARN for CodeBuild. Can reference an AwsIamRole resource via `valueFrom`. | Required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `description` | `string` | — | Project description (max 255 characters) |
| `encryptionKey` | `StringValueOrRef` | AWS-managed key | KMS key ARN for artifact encryption. Can reference AwsKmsKey via `valueFrom`. |
| `buildTimeout` | `int` | `60` | Build timeout in minutes (5-2160) |
| `queuedTimeout` | `int` | `480` | Queue timeout in minutes (5-480) |
| `concurrentBuildLimit` | `int` | Unlimited | Maximum concurrent builds (minimum 1) |
| `sourceVersion` | `string` | — | Default branch, tag, or commit to build |
| `autoRetryLimit` | `int` | `0` | Additional automatic retries after a failed build (max 10) |
| `badgeEnabled` | `bool` | `false` | Publish a build badge (repository sources only; exported as `badge_url`) |
| `projectVisibility` | `string` | `PRIVATE` | `PRIVATE` or `PUBLIC_READ` (world-readable build results) |
| `resourceAccessRole` | `StringValueOrRef` | — | Role CodeBuild uses to read public logs/artifacts. Required with `PUBLIC_READ`. |
| `resourcePolicy` | `object` | — | Resource-based IAM policy document for cross-account access |
| `secondarySources` | `list` | `[]` | Up to 12 extra sources, each with a required `sourceIdentifier` (checked out at `$CODEBUILD_SRC_DIR_<identifier>`) |
| `secondarySourceVersions` | `list` | `[]` | Per-secondary-source branch/tag/commit pins (`sourceIdentifier` + `sourceVersion`) |
| `secondaryArtifacts` | `list` | `[]` | Up to 12 extra outputs, each with a required `artifactIdentifier` |
| `source.buildspec` | `string` | `buildspec.yml` | Build specification file path or inline YAML. Required when type is `NO_SOURCE`. |
| `source.gitCloneDepth` | `int` | Full clone | Git clone depth (0 = full clone) |
| `source.gitSubmodulesConfig.fetchSubmodules` | `bool` | — | Fetch Git submodules (BITBUCKET/CODECOMMIT/GITHUB/GITHUB_ENTERPRISE only) |
| `source.insecureSsl` | `bool` | `false` | Skip TLS verification (self-hosted providers with private CAs only) |
| `source.reportBuildStatus` | `bool` | `false` | Report build status back to the source provider |
| `source.buildStatusConfig` | `object` | — | Custom commit-status `context` and `targetUrl` |
| `source.auth` | `object` | — | Per-source authorization: `type` (`CODECONNECTIONS`, `SECRETS_MANAGER`, `OAUTH`) + `resource` ARN |
| `environment.certificate` | `string` | — | S3 path to a trusted certificate bundle (must end `.pem` or `.zip`) |
| `environment.privilegedMode` | `bool` | `false` | Enable Docker daemon access inside the build container (not for Lambda types) |
| `environment.imagePullCredentialsType` | `string` | `CODEBUILD` | Image pull credentials: `CODEBUILD` or `SERVICE_ROLE` |
| `environment.environmentVariables` | `list` | `[]` | Build environment variables with `name`, `value`, and optional `type` (`PLAINTEXT`, `PARAMETER_STORE`, `SECRETS_MANAGER`) |
| `environment.registryCredential` | `object` | — | Private registry credentials (`credential` ARN, `credentialProvider`: `SECRETS_MANAGER`); requires `SERVICE_ROLE` pulls |
| `environment.dockerServer` | `object` | — | Persistent dedicated Docker server (`computeType` + optional `securityGroupIds`); layer state survives across builds |
| `environment.fleetArn` | `string` | — | Reserved-capacity fleet the project joins (required for `MAC_ARM` / EC2 types) |
| `artifacts.location` | `StringValueOrRef` | — | S3 bucket name. Required when type is `S3`. Can reference AwsS3Bucket via `valueFrom`. |
| `artifacts.name` | `string` | — | Artifact output name |
| `artifacts.path` | `string` | — | S3 prefix path for artifacts |
| `artifacts.packaging` | `string` | `NONE` | Packaging type: `NONE` or `ZIP` |
| `artifacts.namespaceType` | `string` | `NONE` | Namespace: `NONE` or `BUILD_ID` |
| `artifacts.encryptionDisabled` | `bool` | `false` | Disable artifact encryption |
| `artifacts.overrideArtifactName` | `bool` | `false` | Let the buildspec override the artifact name per build |
| `artifacts.bucketOwnerAccess` | `string` | `NONE` | `NONE`, `READ_ONLY`, or `FULL` (cross-account artifact buckets) |
| `cache.type` | `string` | `NO_CACHE` | Cache type: `NO_CACHE`, `S3`, or `LOCAL` |
| `cache.location` | `StringValueOrRef` | — | S3 cache location. Required when type is `S3`. Can reference AwsS3Bucket via `valueFrom`. |
| `cache.modes` | `string[]` | `[]` | Local cache modes: `LOCAL_SOURCE_CACHE`, `LOCAL_DOCKER_LAYER_CACHE`, `LOCAL_CUSTOM_CACHE` |
| `cache.cacheNamespace` | `string` | — | Scopes S3 cache keys so projects/branches can share one cache bucket |
| `logsConfig.cloudwatchLogs.status` | `string` | `ENABLED` | CloudWatch logging: `ENABLED` or `DISABLED` |
| `logsConfig.cloudwatchLogs.groupName` | `StringValueOrRef` | Auto-generated | Log group name. Can reference AwsCloudwatchLogGroup via `valueFrom`. |
| `logsConfig.cloudwatchLogs.streamName` | `string` | Auto-generated | Log stream name prefix |
| `logsConfig.s3Logs.status` | `string` | `DISABLED` | S3 logging: `ENABLED` or `DISABLED` |
| `logsConfig.s3Logs.location` | `StringValueOrRef` | — | S3 bucket and prefix for logs. Can reference AwsS3Bucket via `valueFrom`. |
| `logsConfig.s3Logs.encryptionDisabled` | `bool` | `false` | Disable log file encryption |
| `logsConfig.s3Logs.bucketOwnerAccess` | `string` | `NONE` | Bucket-owner access for centralized cross-account log buckets |
| `vpcConfig.vpcId` | `StringValueOrRef` | — | VPC ID. Can reference AwsVpc via `valueFrom`. Required if vpcConfig is set. |
| `vpcConfig.subnetIds` | `StringValueOrRef[]` | — | VPC subnets (max 16). Can reference AwsSubnet via `valueFrom`. Required if vpcConfig is set. |
| `vpcConfig.securityGroupIds` | `StringValueOrRef[]` | — | Security groups (max 5). Can reference AwsSecurityGroup via `valueFrom`. Required if vpcConfig is set. |
| `fileSystemLocations` | `list` | `[]` | EFS mounts: `identifier`, `location` (`<fs-dns>:/<path>`), `mountPoint`, optional `mountOptions` (requires `vpcConfig`) |
| `buildBatchConfig.serviceRole` | `StringValueOrRef` | — | Role for launching batch child builds. Required if buildBatchConfig is set. |
| `buildBatchConfig.combineArtifacts` | `bool` | `false` | Merge child-build artifacts into one location |
| `buildBatchConfig.timeoutInMins` | `int` | — | Whole-batch timeout (5-2160) |
| `buildBatchConfig.restrictions` | `object` | — | `computeTypesAllowed` + `maximumBuildsAllowed` (1-100) bounds for child builds |
| `webhook.buildType` | `string` | `BUILD` | Webhook build type: `BUILD`, `BUILD_BATCH`, or `RUNNER_BUILDKITE_BUILD` |
| `webhook.manualCreation` | `bool` | `false` | Mint payload URL + secret without registering; wire the repository webhook by hand |
| `webhook.filterGroups` | `list` | `[]` | Filter groups (OR'd). Each group contains filters (AND'd) with `type` (`EVENT`, `BASE_REF`, `HEAD_REF`, `ACTOR_ACCOUNT_ID`, `FILE_PATH`, `COMMIT_MESSAGE`, `WORKFLOW_NAME`, `TAG_NAME`, `RELEASE_NAME`, `REPOSITORY_NAME`, `ORGANIZATION_NAME`), `pattern`, and optional `excludeMatchedPattern`. |
| `webhook.scopeConfiguration` | `object` | — | Organization/group-wide webhook: `name`, `scope` (`GITHUB_ORGANIZATION`, `GITHUB_GLOBAL`, `GITLAB_GROUP`), optional `domain` |
| `webhook.pullRequestBuildPolicy` | `object` | — | Comment-approval gate: `requiresCommentApproval` (`DISABLED`, `FORK_PULL_REQUESTS`, `ALL_PULL_REQUESTS`) + optional `approverRoles` |

## Examples

### GitHub CI with Webhook

A standard CI project triggered by pushes and pull requests on the main branch:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCodeBuildProject
metadata:
  name: app-ci
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsCodeBuildProject.app-ci
spec:
  region: us-west-2
  source:
    type: GITHUB
    location: https://github.com/my-org/my-app.git
    reportBuildStatus: true
  environment:
    type: LINUX_CONTAINER
    computeType: BUILD_GENERAL1_SMALL
    image: aws/codebuild/amazonlinux2-x86_64-standard:5.0
  artifacts:
    type: NO_ARTIFACTS
  serviceRole:
    value: arn:aws:iam::123456789012:role/codebuild-role
  webhook:
    filterGroups:
      - filters:
          - type: EVENT
            pattern: PUSH, PULL_REQUEST_CREATED, PULL_REQUEST_UPDATED
          - type: HEAD_REF
            pattern: ^refs/heads/main$
```

### Docker Build with Layer Caching

A privileged build project for Docker image builds with local layer caching:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCodeBuildProject
metadata:
  name: docker-builder
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsCodeBuildProject.docker-builder
spec:
  region: us-west-2
  source:
    type: GITHUB
    location: https://github.com/my-org/my-service.git
    gitCloneDepth: 1
    reportBuildStatus: true
  environment:
    type: LINUX_CONTAINER
    computeType: BUILD_GENERAL1_LARGE
    image: aws/codebuild/amazonlinux2-x86_64-standard:5.0
    privilegedMode: true
    environmentVariables:
      - name: ECR_REPO
        value: 123456789012.dkr.ecr.us-east-1.amazonaws.com/my-service
      - name: DOCKER_BUILDKIT
        value: "1"
  artifacts:
    type: NO_ARTIFACTS
  serviceRole:
    value: arn:aws:iam::123456789012:role/codebuild-docker-role
  buildTimeout: 30
  cache:
    type: LOCAL
    modes:
      - LOCAL_DOCKER_LAYER_CACHE
      - LOCAL_SOURCE_CACHE
  webhook:
    filterGroups:
      - filters:
          - type: EVENT
            pattern: PUSH
          - type: HEAD_REF
            pattern: ^refs/heads/main$
```

### CodePipeline Build Stage

A build project designed as a stage in AWS CodePipeline with Secrets Manager variables:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCodeBuildProject
metadata:
  name: pipeline-build
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsCodeBuildProject.pipeline-build
spec:
  region: us-west-2
  source:
    type: CODEPIPELINE
    buildspec: buildspec.yml
  environment:
    type: LINUX_CONTAINER
    computeType: BUILD_GENERAL1_MEDIUM
    image: aws/codebuild/amazonlinux2-x86_64-standard:5.0
    environmentVariables:
      - name: STAGE
        value: production
      - name: DB_CONNECTION_STRING
        value: prod/db-connection-string
        type: SECRETS_MANAGER
  artifacts:
    type: CODEPIPELINE
  serviceRole:
    value: arn:aws:iam::123456789012:role/codebuild-pipeline-role
  buildTimeout: 20
  concurrentBuildLimit: 3
```

### Using Foreign Key References

Reference other Planton-managed resources instead of hardcoding IDs:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCodeBuildProject
metadata:
  name: connected-build
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsCodeBuildProject.connected-build
spec:
  region: us-west-2
  source:
    type: GITHUB
    location: https://github.com/my-org/my-app.git
  environment:
    type: LINUX_CONTAINER
    computeType: BUILD_GENERAL1_SMALL
    image: aws/codebuild/amazonlinux2-x86_64-standard:5.0
  artifacts:
    type: NO_ARTIFACTS
  serviceRole:
    valueFrom:
      kind: AwsIamRole
      name: codebuild-role
      field: status.outputs.role_arn
  encryptionKey:
    valueFrom:
      kind: AwsKmsKey
      name: build-key
      field: status.outputs.key_arn
  vpcConfig:
    vpcId:
      valueFrom:
        kind: AwsVpc
        name: main-vpc
        field: status.outputs.vpc_id
    subnetIds:
      - valueFrom:
          kind: AwsSubnet
          name: main-private-subnet-a
          fieldPath: status.outputs.subnet_id
      - valueFrom:
          kind: AwsSubnet
          name: main-private-subnet-b
          fieldPath: status.outputs.subnet_id
    securityGroupIds:
      - valueFrom:
          kind: AwsSecurityGroup
          name: codebuild-sg
          field: status.outputs.security_group_id
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `project_arn` | `string` | ARN of the CodeBuild project, used for IAM policies and cross-resource references |
| `project_name` | `string` | Name of the CodeBuild project (what CodePipeline Build actions reference) |
| `service_role_arn` | `string` | IAM service role ARN used by the project |
| `badge_url` | `string` | Build badge URL (empty when `badgeEnabled` is false) |
| `public_project_alias` | `string` | Public alias in the public build results URL (empty for private projects) |
| `webhook_url` | `string` | Webhook URL from the source provider (empty when no webhook is configured) |
| `webhook_payload_url` | `string` | URL that receives webhook payloads — with `manualCreation`, register it on the repository by hand (empty when no webhook is configured) |
| `webhook_secret` | `string` | Webhook HMAC signing secret, only minted at creation — sensitive, treat as a credential (empty when no webhook is configured) |

## Related Components

- [AwsIamRole](/docs/catalog/aws/iam-role) — provides the service role for CodeBuild
- [AwsVpc](/docs/catalog/aws/vpc) — provides VPC and subnets for builds that access private resources
- [AwsSecurityGroup](/docs/catalog/aws/security-group) — controls network access for VPC-enabled builds
- [AwsS3Bucket](/docs/catalog/aws/s3-bucket) — stores build artifacts and cache
- [AwsCloudwatchLogGroup](/docs/catalog/aws/cloudwatch-log-group) — hosts build logs
