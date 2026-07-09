---
title: "Fork-PR-Gated Open-Source CI"
description: "This preset creates a public-repository CI project hardened for the open-source contribution model: pull requests from forks wait for an approval comment from a maintainer before CodeBuild runs them,..."
type: "preset"
rank: "04"
presetSlug: "04-fork-pr-gated-oss-ci"
componentSlug: "codebuild-project"
componentTitle: "CodeBuild Project"
provider: "aws"
icon: "package"
order: 4
---

# Fork-PR-Gated Open-Source CI

This preset creates a public-repository CI project hardened for the open-source contribution model: pull requests from forks wait for an approval comment from a maintainer before CodeBuild runs them, so untrusted fork code can never reach CI credentials unreviewed. Source access flows through a CodeConnections connection (no stored tokens), a build badge advertises live status, and failed builds retry once automatically.

## When to Use

- Open-source repositories that accept pull requests from forks
- Any repository where CI credentials must be protected from untrusted contributor code
- Teams standardizing on CodeConnections (the modern, token-free source authorization)
- Projects that publish a build badge in their README

## Key Configuration Choices

- **CODECONNECTIONS** (`source.auth.type`) — repository access through a CodeConnections connection instead of a stored OAuth token
- **FORK_PULL_REQUESTS** (`pullRequestBuildPolicy.requiresCommentApproval`) — only fork PRs wait for approval; same-repo PRs build immediately
- **Approver roles** — comments from GitHub users with write, maintain, or admin roles count as approval
- **ARM_CONTAINER + BUILD_GENERAL1_MEDIUM** — Graviton compute, typically the best price/performance for CI
- **`badgeEnabled` + `autoRetryLimit: 1`** — public build badge, one automatic retry absorbs flaky-infrastructure failures

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<github-repo-https-url>` | GitHub repository HTTPS URL (e.g., `https://github.com/org/repo.git`) | GitHub repository settings |
| `<codeconnections-connection-arn>` | CodeConnections connection ARN (must be AVAILABLE after the one-time console handshake) | AWS Developer Tools console → Connections |
| `<codebuild-service-role-arn>` | IAM role ARN granting CodeBuild access to source and logs | AWS IAM console or `AwsIamRole` status outputs |

## Related Presets

- **01-github-ci-linux** — Use instead for private repositories without the fork-approval gate
- **02-docker-build-ecr** — Use instead when the build produces container images
