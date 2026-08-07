# GitHub Workload Identity Impersonation

This preset grants `roles/iam.workloadIdentityUser` on one service account to the federated principal set of one GitHub repository — the terminal hop of keyless CI/CD. After this grant, workflows in that repository can exchange their GitHub OIDC token for the service account's credentials without any long-lived key.

## When to Use

- A GitHub Actions workflow needs to deploy or operate AS a dedicated service account
- You run keyless authentication (workload identity federation) and need the impersonation hop
- You want the SA-scoped grant instead of a project-wide one, so the repository can impersonate exactly one account

## Key Configuration Choices

- **`valueFrom` serviceAccountId** — referencing the GcpServiceAccount's `name` output keeps the graph honest and survives account recreation
- **`attribute.repository` scoping** — the principal set covers only workflows of the named repository; use `attribute.repository_owner` paths for org-wide trust, or tighten further per branch with a subject-scoped principal
- **Project NUMBER in the principal** — federation principals embed the numeric project, never the project ID

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<service-account-resource-name>` | The Planton resource name of the impersonated GcpServiceAccount | Your GcpServiceAccount manifest's `metadata.name` |
| `<gcp-project-number>` | The numeric project number hosting the pool | `GcpProject` outputs or GCP Console |
| `<pool-id>` | The workload identity pool ID | Your GcpWorkloadIdentityPool manifest |
| `<github-org>` / `<github-repo>` | The trusted GitHub repository | GitHub |

## Related Presets

- **02-token-creator-grant** — Let one account mint short-lived tokens as another
- **03-deployer-act-as** — The actAs permission Cloud Run/GCE deployments require
