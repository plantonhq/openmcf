---
title: "How to Connect Your GitHub Account to Planton"
date: "2026-04-02"
author:
  - name: "Planton Team"
    title: "Platform Engineering"
    bio: "Helping teams deploy infrastructure and services without the DevOps bottleneck"
tags:
  - "github"
  - "connect"
  - "scm-connection"
  - "getting-started"
category: "connect"
excerpt: "Connect your GitHub repositories to Planton so the platform can clone code, receive push events, and deploy your services."
---

# How to Connect Your GitHub Account to Planton

Before Planton can build and deploy your code, it needs access to your GitHub repositories. This tutorial walks you through installing the Planton GitHub App on your GitHub organization (or personal account) to establish that connection. Once connected, Planton can clone your repositories, receive webhook notifications when you push code, and report build status back to GitHub as check runs.

> **Note**: This connection is created through the Planton web console because it requires completing GitHub's App installation flow in your browser. Unlike provider connections (AWS, GCP), there is no CLI-only path for the initial setup. Once created, the connection is fully manageable via CLI.

## What You Will Learn

- What an SCM connection is and how it differs from provider connections
- How the GitHub App model works (no personal access tokens, no stored secrets)
- How to install the Planton GitHub App on your GitHub organization or personal account
- How to verify the connection and list accessible repositories from the CLI

## Prerequisites

- [ ] A Planton account with an organization already created
- [ ] The `planton` CLI installed and authenticated (run `planton auth login` if you have not done this yet)
- [ ] A GitHub account with **owner** or **admin** access to the organization you want to connect (or a personal account)

## How GitHub Connections Work

Planton connects to GitHub using a **GitHub App** -- not personal access tokens. The App uses short-lived installation tokens (one hour expiry), declares granular permissions you review during installation, and is scoped to your organization rather than any individual user. Once installed, Planton uses this connection to clone repositories, receive push webhooks, and report build status as GitHub check runs. For more details, see the [Git provider connections documentation](/docs/connections/git-providers).

## Step 1: Start the GitHub App Installation from the Console

Open the Planton web console and navigate to the **Connect** section of your organization. Select **Add Connection**, then choose **GitHub Connection**.

The console presents a creation form where you provide a **name** for the connection (for example, `my-github` or `acme-github`). This name becomes the connection's slug, which you will reference when creating Services later.

After entering a name, the console directs you to install the Planton GitHub App. This opens GitHub's App installation page in your browser.

## Step 2: Install the Planton GitHub App on GitHub

On GitHub's installation page, you choose:

1. **Which account** to install the App on -- your personal account or one of your GitHub organizations
2. **Which repositories** to grant access to -- all repositories, or a selection of specific repositories

For most teams, installing on the organization account and granting access to **all repositories** is the simplest approach. You can change this scope later from your GitHub organization settings without affecting the Planton connection.

After you approve the installation, GitHub redirects you back to the Planton console with the installation metadata. The console creates the `GithubConnection` resource automatically.

## Step 3: Verify the Connection via CLI

Once the console confirms the connection was created, verify it from the CLI:

```bash
planton get github-connection my-github
```

This returns the connection resource, confirming the `installation_id`, `account_type`, and `account_id` were captured correctly.

## Step 4: Verify Repository Access

Confirm that Planton can see the repositories available through this connection:

```bash
planton service github list-repos my-github
```

This lists the GitHub repositories accessible to the Planton GitHub App installation. The list reflects the repository scope you selected during installation (all repositories, or a specific set).

If you do not see a repository you expect, check the App's repository access settings in your GitHub organization: **Settings > Integrations > GitHub Apps > Planton > Configure**.

## Connecting GitHub Enterprise Server

If your organization uses GitHub Enterprise Server instead of github.com, the process is the same with one additional configuration: the `github_connection_host` field in the spec points to your Enterprise Server URL instead of the default `https://github.com`.

```yaml
apiVersion: connect.planton.ai/v1
kind: GithubConnection
metadata:
  name: ghe-internal
  org: your-org
spec:
  github_connection_host:
    value: "https://github.enterprise.corp"
  app_install_info:
    installation_id: 87654321
    account_type: organization
    account_id: "platform-team"
```

The Planton GitHub App must be registered on your GitHub Enterprise Server instance for this to work. Contact your Planton account team for Enterprise Server configuration.

## What to Do Next

- **Connect a container registry** so Planton can push the images it builds: [How to Connect a Container Registry to Planton](/tutorials/how-to-connect-a-container-registry-to-planton)
- **Deploy your first service** with zero-config CI/CD: [How to Deploy Your First Service with Zero-Config CI/CD](/tutorials/how-to-deploy-your-first-service-with-zero-config-cicd)
