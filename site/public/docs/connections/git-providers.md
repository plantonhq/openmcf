---
title: "Git Providers"
description: "Connect GitHub and GitLab organizations to enable automated builds, webhooks, and pull request deployments in Service Hub"
icon: connect
order: 30
tags:
  - Connect
  - Git Providers
  - GitHub
  - GitLab
---

# Git Providers

Service Hub builds and deploys code from your Git repositories. For that to work, Planton needs read access to your repositories and the ability to set up webhooks for commit-triggered pipelines. Git provider connections establish that access at the organization level, so every service you create can pull code without engineers managing individual access tokens.

## Why Git Connections Matter

Without a git connection, Service Hub cannot:

- Clone your repository to run builds
- Set up webhooks to trigger pipelines on push events
- Report build status back to pull requests
- Access private repositories

A single git connection covers your entire GitHub or GitLab organization. Every service that references a repository under that organization automatically uses the connection — no per-repository configuration needed.

## GitHub

GitHub is the primary git provider in Planton. The connection uses GitHub's App installation model — an OAuth-based flow that grants Planton access to your organization's repositories without requiring personal access tokens.

### How the Connection Works

Unlike cloud provider connections where you paste credentials, the GitHub connection uses an OAuth flow:

1. You initiate the connection from the Planton web console.
2. Planton redirects you to GitHub to authorize the Planton GitHub App.
3. You choose which GitHub organization (or personal account) to install the app on.
4. GitHub redirects you back to Planton with the installation details.
5. The connection is created automatically — no credentials to copy or paste.

This model has several advantages over token-based access:

- **Scoped permissions** — The GitHub App requests only the permissions it needs (repository read, webhook management, commit status).
- **Organization-level** — One installation covers all repositories in the organization.
- **No token rotation** — GitHub manages the app's authentication tokens. You don't need to rotate credentials manually.
- **Audit trail** — GitHub logs all API calls made by the app installation.

### Connecting via the Web Console

1. Navigate to **Connections** and click the **GitHub** card under DevOps Pipeline.
2. Planton redirects you to GitHub's authorization page.
3. Select the GitHub organization (or personal account) you want to connect.
4. Review the permissions and approve the installation.
5. GitHub redirects you back to Planton. The connection appears immediately.

<!-- SCREENSHOT: GitHub OAuth authorization page
  Page: GitHub.com OAuth authorization flow
  Action: Show the GitHub App installation approval screen
  Focus: The organization selector and permission list
  Alt: GitHub App installation page showing organization selection and permission approval
-->

### Connecting via the CLI

The CLI provides a utility to verify your GitHub connection by listing accessible repositories:

```bash
# List repositories accessible through a GitHub connection
planton connect github list-repos my-github-org
```

### GitHub Enterprise (Self-Hosted)

If your organization uses GitHub Enterprise Server (self-hosted), the connection supports a custom host URL. During setup, specify your GitHub Enterprise endpoint instead of the default `https://github.com`.

### Account Types

GitHub connections support both organization accounts and personal accounts. Organization accounts are the typical choice for teams — they cover all repositories under the organization. Personal accounts are useful for individual developers or for connecting repositories that aren't under an organization.

---

## GitLab

GitLab connections use an access token and refresh token pair to authenticate. Unlike GitHub's OAuth App model, GitLab connections require you to provide a group ID and tokens.

### Connecting via the Web Console

1. Navigate to **Connections** and click the **GitLab** card under DevOps Pipeline.
2. **Name your connection** and provide the required credentials.
3. **Create the connection**.

> **Note**: GitLab integration is currently being finalized in the web console. The connection type is available in the API, but the web console wizard may show limited functionality. Check back for updates as GitLab support is expanded.

### What You Need

| Field | Description |
|-------|-------------|
| GitLab Host | Your GitLab instance URL (defaults to `https://gitlab.com`) |
| Group ID | The GitLab group whose repositories Planton will access |
| Access Token | A GitLab personal access token or group access token |
| Refresh Token | The refresh token for token renewal |

### GitLab Self-Hosted

GitLab connections support self-hosted GitLab instances. Provide your instance's URL as the GitLab host. This is useful for organizations that run GitLab on their own infrastructure for compliance or data residency requirements.

### Creating GitLab Tokens

In GitLab:

1. Navigate to your group's **Settings > Access Tokens** (for group tokens) or **User Settings > Access Tokens** (for personal tokens).
2. Create a token with the `api` scope for full API access, or scope it down to `read_repository` and `write_repository` for minimal permissions.
3. Copy both the access token and the refresh token.

---

## How Git Connections Are Used

Once a git connection is established, Service Hub uses it automatically when you create a service:

1. **Service creation** — When you configure a service's source repository, Planton uses the git connection to verify the repository exists and is accessible.
2. **Pipeline execution** — When a pipeline runs, Planton clones the repository using the git connection. For monorepos, sparse checkout is used to pull only the relevant directory.
3. **Webhook management** — Planton registers webhooks on the repository to trigger pipelines on push events, pull requests, and tag creation.
4. **Status reporting** — Build and deployment results are reported back to the git provider as commit statuses or check runs.

You don't need to specify which git connection to use on each service — Planton resolves it automatically based on the repository's organization.

## Related Documentation

- [Connections Overview](/docs/connections) — Understanding the Connect system
- [CI/CD](/docs/ci-cd) — How services use git connections for builds and deployments
- [Container Registries](/docs/connections/container-registries) — Where build artifacts are stored
