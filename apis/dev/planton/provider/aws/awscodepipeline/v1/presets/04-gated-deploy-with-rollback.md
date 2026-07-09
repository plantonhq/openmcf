# Gated Production Deploy with Automatic Rollback

This preset creates a three-stage V2 pipeline that treats the production deploy as a guarded, self-healing operation: the deploy stage only opens during business hours (a DeploymentWindow entry gate), a CloudWatch alarm check verifies the deployment after it lands, and any failure in the deploy stage rolls the stage back to its last successful state automatically. The build stage retries its failed actions once before failing the execution.

## When to Use

- Production deployment pipelines where a bad deploy must self-revert without an operator
- Teams that restrict production changes to business hours (change-window policies)
- Services with a reliable post-deploy health alarm (error rate, 5xx count) that should gate success
- Anywhere the "deploy, watch, roll back" loop is currently a human runbook

## Key Configuration Choices

- **DeploymentWindow rule** (`beforeEntry`) — the deploy stage admits executions only Monday-Friday, 9:00-17:00 in the configured time zone; executions arriving outside the window wait
- **CloudWatchAlarm rule** (`onSuccess`) — after the ECS deploy succeeds, the pipeline watches the named alarm for 5 minutes; if it fires, the stage is marked failed (which then triggers the rollback)
- **ROLLBACK** (`onFailure.result` on the deploy stage) — automatic revert to the stage's last successful execution state
- **RETRY / FAILED_ACTIONS** (`onFailure` on the build stage) — transient build failures re-run only the failed actions
- **QUEUED** (`executionMode`) — executions line up rather than superseding, so every commit lands in order

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<codepipeline-service-role-arn>` | IAM role CodePipeline assumes | AWS IAM console or `AwsIamRole` status outputs |
| `<pipeline-artifacts-bucket-name>` | Versioned S3 bucket for pipeline artifacts | `AwsS3Bucket` status outputs |
| `<codeconnections-connection-arn>` | CodeConnections connection ARN (AVAILABLE after the one-time console handshake) | AWS Developer Tools console → Connections |
| `<github-org>/<github-repo>` | Full repository identifier | GitHub |
| `<codebuild-project-name>` | CodeBuild project executing the build stage | `AwsCodeBuildProject` status outputs |
| `<iana-timezone>` | Deployment-window time zone (e.g., `America/Los_Angeles`) | IANA time zone database |
| `<ecs-cluster-name>` / `<ecs-service-name>` | Deploy target | `AwsEcsCluster` / `AwsEcsService` status outputs |
| `<post-deploy-error-alarm-name>` | CloudWatch alarm gating deploy success | `AwsCloudwatchAlarm` status outputs |

## Related Presets

- **01-github-source-codebuild** — Use instead for a plain source-build pipeline without deploy gating
- **02-ecr-ecs-deploy** — Use instead when the source of truth is a container image landing in ECR
