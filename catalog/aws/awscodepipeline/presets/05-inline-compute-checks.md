# Inline Compute Checks Without a Build Project

This preset creates a three-stage V2 pipeline whose quality gate runs as a Compute action — shell commands executed directly in CodePipeline-managed compute, with no CodeBuild project to create or maintain. The commands lint and test the source, export a coverage percentage as a pipeline variable for downstream actions, and publish the lint/coverage reports as a file-based output artifact. The deploy stage rolls back automatically on failure.

## When to Use

- Lightweight CI steps (lint, unit tests, policy checks) that do not justify a standing CodeBuild project
- Monorepos where each pipeline needs a slightly different check step and per-pipeline build projects would proliferate
- Pipelines that need a computed value (a coverage number, a version string) exported into downstream action configuration
- Anywhere the check step's reports should flow to later stages as a real artifact

## Key Configuration Choices

- **Compute action** (`category: Compute`, `provider: Commands`) — CodePipeline runs the `commands` list in managed compute; CodeBuild logs and permissions are used under the hood, and compute time bills per execution (a few cents per run)
- **`outputVariables` + `namespace`** — variables the commands export become `#{CheckVars.COVERAGE_PCT}` in downstream actions
- **`outputArtifactsForComputeAction`** — Compute actions export named FILE sets instead of the plain `outputArtifacts` every other action category uses (the two are mutually exclusive, enforced at validation)
- **ROLLBACK** (`onFailure.result` on the deploy stage) — automatic revert to the stage's last successful execution state
- **QUEUED** (`executionMode`) — executions line up rather than superseding, so every commit lands in order

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<codepipeline-service-role-arn>` | IAM role CodePipeline assumes | AWS IAM console or `AwsIamRole` status outputs |
| `<pipeline-artifacts-bucket-name>` | Versioned S3 bucket for pipeline artifacts | `AwsS3Bucket` status outputs |
| `<codeconnections-connection-arn>` | CodeConnections connection ARN (AVAILABLE after the one-time console handshake) | AWS Developer Tools console → Connections |
| `<github-org>/<github-repo>` | Full repository identifier | GitHub |
| `<ecs-cluster-name>` / `<ecs-service-name>` | Deploy target | `AwsEcsCluster` / `AwsEcsService` status outputs |

The `commands` list itself is example content — replace it with your repository's real check steps (any shell commands except multi-line formats).

## Related Presets

- **01-github-source-codebuild** — Use instead when the build/check step warrants a full CodeBuild project (custom images, VPC access, caching)
- **04-gated-deploy-with-rollback** — Combine with this preset's Compute checks for deployment-window gating and post-deploy alarm verification
