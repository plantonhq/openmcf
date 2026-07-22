# AWS CI/CD Pipeline

Push-to-production delivery on AWS's native tooling: a CodePipeline V2
that triggers on every push to your branch, builds a Docker image in a
CodeBuild project driven by the buildspec in your repository, pushes it to
a private ECR repository with immutable tags — and, when you flip the
toggle, rolls an ECS service to the new image with no scripts and no
long-lived CI credentials. Everything runs on IAM roles scoped to this
pipeline's own resources.

The one thing no chart can automate is the Git connection: AWS requires a
human to approve the app installation on your repository host. That
one-time handshake is the first step below; everything after it is this
chart.

## Architecture

```
 git push ──▶ CodeConnections connection (one-time human handshake)
                │  V2 push trigger, branch-filtered
        AwsCodePipeline (QUEUED — every commit gets its own run)
                │
   Source ──▶ Build ─────────────────▶ Deploy (ecs_deploy_enabled)
                │                          │
        AwsCodeBuildProject          ECS rolling deploy
          (Docker, privileged)       (imagedefinitions.json →
                │                     new task-def revision →
                ├─▶ AwsEcrRepo        service rolls with its own
                │   (immutable tags)  circuit breaker)
                │
        AwsS3Bucket (artifacts, encrypted, 30-day expiry)
        AwsIamRole ×2 (CodeBuild / CodePipeline, least-privilege)
```

## Included Cloud Resources

| Resource | Kind | Purpose |
|----------|------|---------|
| Pipeline | `AwsCodePipeline` | V2 pipeline: git trigger, queued executions, the stage graph |
| Build project | `AwsCodeBuildProject` | Docker-enabled build driven by your repo's buildspec |
| Image repository | `AwsEcrRepo` | Immutable-tag registry the build pushes to, scan-on-push |
| Artifact bucket | `AwsS3Bucket` | Stage-to-stage artifact store — encrypted, private, self-pruning |
| CodeBuild role | `AwsIamRole` | Logs + artifact exchange + push to this repo, nothing else |
| CodePipeline role | `AwsIamRole` | Artifacts + this connection + this project (+ ECS deploy when toggled) |

## Parameters

| Parameter | Description | Default | Type |
|-----------|-------------|---------|------|
| `aws_region` | Region for everything the chart creates | `us-east-1` | string |
| `pipeline_name` | Name prefix and the ECR repository name | `app` | string |
| `connection_arn` | The CodeConnections ARN from the one-time handshake | placeholder | string |
| `repository_id` | The watched repository, owner/name | `my-org/my-service` | string |
| `branch` | The branch whose pushes trigger runs | `main` | string |
| `artifacts_bucket_name` | Globally unique artifact bucket name | placeholder | string |
| `build_compute_type` | CodeBuild fleet size (the cost dial) | `BUILD_GENERAL1_SMALL` | string |
| `build_image` | The managed build image | amazonlinux2 standard:5.0 | string |
| `buildspec_path` | Buildspec location inside your repository | `buildspec.yml` | string |
| `ecs_deploy_enabled` | Add the ECS rolling-deploy stage | `false` | bool |
| `ecs_cluster_name` / `ecs_service_name` | The service to roll (fargate-web-service pairing) | `web-cluster` / `web` | string |

## Setup

1. **Create the connection (one-time, human-approved).** In the AWS
   console: Developer Tools → Settings → Connections → Create connection →
   pick your provider (GitHub, GitLab, Bitbucket) → approve the app
   installation on your organization. Paste the resulting ARN into
   `connection_arn`. A connection is reusable by every pipeline in the
   account — one handshake per Git organization, ever.
2. **Add the buildspec to your repository** (at `buildspec_path`). This one
   builds the image, pushes it with an immutable per-commit tag, and emits
   the deploy artifact:

```yaml
version: 0.2
phases:
  pre_build:
    commands:
      - AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
      - ECR_URL=$AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com/$IMAGE_REPO_NAME
      - aws ecr get-login-password | docker login --username AWS --password-stdin $ECR_URL
      - IMAGE_TAG=${CODEBUILD_RESOLVED_SOURCE_VERSION:0:12}
  build:
    commands:
      - docker build -t $ECR_URL:$IMAGE_TAG .
  post_build:
    commands:
      - docker push $ECR_URL:$IMAGE_TAG
      - printf '[{"name":"%s","imageUri":"%s"}]' "$CONTAINER_NAME" "$ECR_URL:$IMAGE_TAG" > imagedefinitions.json
artifacts:
  files:
    - imagedefinitions.json
```

   `IMAGE_REPO_NAME` and `CONTAINER_NAME` arrive as environment variables
   from the build project; the account id is looked up at runtime so it is
   never baked into configuration.

3. **Deploy the chart.** Push to the branch — the pipeline runs, and the
   image lands in ECR tagged with the commit SHA.

## Deploying to ECS

Flip `ecs_deploy_enabled` with the cluster and service names (with the
fargate-web-service chart: `<service_name>-cluster` and `<service_name>`).
The deploy action reads `imagedefinitions.json`, registers a new
task-definition revision pointing at the fresh image, and rolls the
service — which deploys under its own safety net (circuit breaker,
health-check gating, rollback).

## After deploying

- `AwsCodePipeline` → `status.outputs.pipeline_arn`
- `AwsCodeBuildProject` → `status.outputs.project_name` (CloudWatch logs
  live under `/aws/codebuild/<project_name>`)
- `AwsEcrRepo` → `status.outputs.repository_url` (what production task
  definitions reference)

## Day-2 guidance

- **Tighten the account wildcards**: the roles' resource ARNs wildcard the
  account segment because the account id is unknown at render time.
  Replacing `*` with your account id in the two roles' inline policies is
  the one-line hardening pass.
- **Manual approval gate**: add an Approval-category action (provider
  Manual) as a stage before Deploy for environments where a human signs
  off on releases; the pipeline pauses until approved.
- **Tests in the build**: add a test phase to the buildspec — a failing
  phase fails the build and stops the pipeline, which is the entire CI
  contract.
- **Fork-PR / OSS posture**: for public repositories, keep this pipeline
  on push triggers to protected branches only; PR-triggered builds from
  forks belong in a separate, unprivileged build project with
  `pull_request_build_policy` fork gating.
- **More stages**: staging-then-production is two ECS deploy actions with
  different cluster/service names and an approval between them — extend
  the `stages` list.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
