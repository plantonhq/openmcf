---
title: "Team Workspaces SageMaker Domain"
description: "This preset creates a Studio domain provisioned as a complete team setup in one manifest: the domain itself, one user profile per team member (the per-person plane — each gets a private home..."
type: "preset"
rank: "05"
presetSlug: "05-team-workspaces"
componentSlug: "sagemaker-domain"
componentTitle: "SageMaker Domain"
provider: "aws"
icon: "package"
order: 5
---

# Team Workspaces SageMaker Domain

This preset creates a Studio domain provisioned as a complete team setup in
one manifest: the domain itself, one user profile per team member (the
per-person plane — each gets a private home directory and inherits the
domain's defaults), and a shared space (the collaboration plane — one
JupyterLab runtime and EBS volume that everyone with access edits
together). Adding a teammate is adding one line to `userProfiles`.

## When to Use

- ML teams standing up Studio where membership is known and managed as
  code, not clicked together in the console
- Pair/mob data-science workflows that want a genuinely shared notebook
  runtime rather than per-person copies
- Any domain where "who has a workspace here" should be reviewable in a
  pull request

## Key Configuration Choices

- **Profiles as list entries** (`userProfiles`) — name-keyed satellites:
  adding or removing one person never disturbs the others; removing an
  entry deletes that profile AND its home directory surfaces, so treat
  removals as destructive
- **One shared space** (`spaces` with `sharingType: Shared`) — ownership
  and sharing are declared together (AWS requires the pair); the owner
  must be a profile in the domain, normally one of `userProfiles`
- **Idle shutdown on both planes** — the domain baseline enforces
  auto-shutdown for personal apps (120 min) and the shared space carries a
  tighter 60-minute timeout, since an idle shared runtime bills like any
  other
- **IAM auth** — profiles map to IAM sessions; for IAM Identity Center
  (SSO) domains, add `singleSignOnUserIdentifier: UserName` +
  `singleSignOnUserValue` per profile
- **Example profile names** (`alice`, `bob`) — replace with your team
  roster (1-63 chars, alphanumeric and hyphens; the space's owner must be
  one of them)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<aws-region>` | The AWS region for the domain | Your region strategy |
| `<vpc-id>` | The VPC the domain attaches to | VPC console or your AwsVpc resource outputs |
| `<private-subnet-id-az1>` / `<private-subnet-id-az2>` | Two subnets in distinct AZs | Subnet console or your AwsSubnet outputs |
| `<sagemaker-execution-role-arn>` | The role Studio apps run as (trusts sagemaker.amazonaws.com) | IAM console or your AwsIamRole outputs |

## Related Presets

- **01-basic-jupyter-domain** — the smallest useful domain, no satellites
- **02-production-vpc-only** — the locked-down network posture
- **03-ml-team-with-custom-images** — custom kernel images via
  AppImageConfig
- **04-governed-canvas-workspace** — the Canvas governance arms
