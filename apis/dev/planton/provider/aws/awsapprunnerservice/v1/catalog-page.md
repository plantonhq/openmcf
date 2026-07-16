# AWS App Runner Service

Deploys an AWS App Runner service from a container image or a source repository with automatic HTTPS, concurrency-based auto scaling, custom domains, optional VPC egress, WAF protection, and customer-managed KMS encryption. The component supports two mutually exclusive source types: an ECR image or a connected code repository.

## What Gets Created

When you deploy an AwsAppRunnerService resource, Planton provisions:

- **App Runner Service** — the runtime itself: source configuration (image or code), instance sizing, health checks, networking, encryption, and observability settings
- **Custom Domain Associations** — one per `customDomains` entry; App Runner issues and renews each domain's TLS certificate, and the validation CNAME records surface as stack outputs
- **WAF Web ACL Association** — created only when `webAclArn` references a REGIONAL web ACL, placing WAF inspection in front of every request

The service's shared companions are separate first-class resources referenced by ARN, never created here: `AwsAppRunnerAutoScalingConfiguration` (scaling posture), `AwsAppRunnerVpcConnector` (VPC egress), and `AwsAppRunnerObservabilityConfiguration` (X-Ray tracing). Each is designed by AWS to be shared across services and tuned in one place.

## Prerequisites

- **AWS credentials** configured via Planton provider config
- **A container image in ECR or ECR Public** if using image-based deployment
- **An IAM access role** assumable by `build.apprunner.amazonaws.com` with ECR pull permissions if using a private ECR image (`imageRepositoryType: ECR`)
- **An App Runner connection** (created in the AWS console with a one-time OAuth handshake) if using code-based deployment
- **An `AwsAppRunnerVpcConnector`** if the service needs to reach VPC resources
- **A customer-managed KMS key** if encrypting with your own key (changing it replaces the service)

## Quick Start

Create a file `app-runner.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAppRunnerService
metadata:
  name: my-api
spec:
  region: us-east-1
  imageSource:
    imageIdentifier: public.ecr.aws/aws-containers/hello-app-runner:latest
    imageRepositoryType: ECR_PUBLIC
  port: "8000"
```

Deploy:

```shell
planton apply -f app-runner.yaml
```

This creates a publicly accessible App Runner service with default settings: 1 vCPU, 2 GB memory, and the account's default auto scaling configuration.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region for the service (e.g., `us-east-1`). | Required; non-empty |
| `imageSource` or `codeSource` | `object` | Deployment source. Exactly one must be provided. | CEL: `exactly_one_source` |
| `imageSource.imageIdentifier` | `string` | Full container image URI including tag or digest. | Min length 1 |
| `imageSource.imageRepositoryType` | `string` | `ECR` (private) or `ECR_PUBLIC`. | Closed set |
| `imageSource.accessRoleArn` | `ref` | IAM role App Runner uses to pull from private ECR. Required when type is `ECR`. | CEL: `ecr_requires_access_role` |
| `codeSource.repositoryUrl` | `string` | Repository URL (e.g., `https://github.com/owner/repo`). | Min length 1 |
| `codeSource.branch` | `string` | Branch to deploy from. | Min length 1 |
| `codeSource.connectionArn` | `ref` | ARN of the App Runner connection authorizing repository access. | Required |
| `codeSource.configurationSource` | `string` | `API` (build config inline) or `REPOSITORY` (apprunner.yaml in the repo). | Closed set |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `port` | `string` | `8080` | Port the application listens on. |
| `startCommand` | `string` | — | Override the container start command. |
| `environmentVariables` | `map<string,string>` | `{}` | Plaintext environment variables. `AWSAPPRUNNER`-prefixed keys are reserved. |
| `environmentSecrets` | `map<string,string>` | `{}` | Values are Secrets Manager / SSM Parameter Store ARNs; resolved at deploy time. The instance role must be allowed to read them. |
| `cpu` | `string` | `1024` | `256`–`4096` or `0.25 vCPU`–`4 vCPU`. |
| `memory` | `string` | `2048` | `512`–`12288` (MB) or `0.5 GB`–`12 GB`. Not all CPU/memory pairings are valid. |
| `instanceRoleArn` | `ref` | — | IAM role instances assume at runtime (references `AwsIamRole`). |
| `healthCheck.protocol` | `string` | `TCP` | `TCP` (port open) or `HTTP` (GET expecting success). |
| `healthCheck.path` | `string` | `/` | Path for HTTP checks; ignored for TCP. |
| `healthCheck.interval` | `int` | `5` | Seconds between checks (1–20). |
| `healthCheck.timeout` | `int` | `2` | Seconds to wait for a response (1–20). |
| `healthCheck.healthyThreshold` | `int` | `1` | Consecutive successes to mark healthy (1–20). |
| `healthCheck.unhealthyThreshold` | `int` | `5` | Consecutive failures before replacement (1–20). |
| `autoScalingConfigurationArn` | `ref` | account default | References `AwsAppRunnerAutoScalingConfiguration`. |
| `vpcConnectorArn` | `ref` | — | References `AwsAppRunnerVpcConnector` for VPC egress. |
| `observabilityConfigurationArn` | `ref` | — | References `AwsAppRunnerObservabilityConfiguration`; the reference enables tracing. |
| `isPubliclyAccessible` | `bool` | `true` | When `false`, reachable only through a VPC Ingress Connection. |
| `ipAddressType` | `string` | `IPV4` | `IPV4` or `DUAL_STACK`. |
| `kmsKeyArn` | `ref` | AWS-managed key | Customer-managed encryption key (references `AwsKmsKey`). **ForceNew.** |
| `autoDeploymentsEnabled` | `bool` | `false` | Auto-deploy on source changes. Private ECR and code repos only — **AWS rejects it for ECR_PUBLIC** (CEL-validated). |
| `customDomains[]` | `list` | `[]` | Custom domains keyed by `domainName`; `enableWwwSubdomain` defaults true. |
| `webAclArn` | `ref` | — | REGIONAL `AwsWafWebAcl` to associate. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `service_arn` | Full ARN — the join key for VPC Ingress Connections and deployment triggers. |
| `service_id` | AWS-assigned service identifier. |
| `service_url` | Default HTTPS endpoint (scheme-less). |
| `service_name` | The service name (metadata.name). |
| `service_status` | Lifecycle status at the end of deployment. |
| `custom_domains` | Per-domain `dns_target` + certificate-validation CNAMEs — compose into `AwsRoute53DnsRecord`. |

## Composing the Family

```yaml
spec:
  autoScalingConfigurationArn:
    valueFrom:
      kind: AwsAppRunnerAutoScalingConfiguration
      name: prod-scaling
      fieldPath: status.outputs.configuration_arn
  vpcConnectorArn:
    valueFrom:
      kind: AwsAppRunnerVpcConnector
      name: private-backend-access
      fieldPath: status.outputs.vpc_connector_arn
  observabilityConfigurationArn:
    valueFrom:
      kind: AwsAppRunnerObservabilityConfiguration
      name: xray-tracing
      fieldPath: status.outputs.configuration_arn
  webAclArn:
    valueFrom:
      kind: AwsWafWebAcl
      name: prod-edge-acl
      fieldPath: status.outputs.web_acl_arn
```

Because the companion ARNs carry revisions, registering a new configuration revision rolls every referencing service on its next deployment — fleet-wide tuning through the resource graph.

## Deliberately Omitted

- **App Runner connections** (`aws_apprunner_connection`): require an out-of-band OAuth handshake in the console to become usable; compose by literal ARN.
- **VPC Ingress Connections** (`aws_apprunner_vpc_ingress_connection`): the inbound private-access plane; composes against the exported `service_arn`.
- **The deployment trigger** (`aws_apprunner_deployment`): a one-shot operation, not infrastructure.
- **Per-kind tags**: identity tags derive from metadata.
