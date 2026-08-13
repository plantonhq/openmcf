# Git Source Web App

This preset deploys a DigitalOcean App Platform application with one HTTP service built from a public Git repository. App Platform clones the repo, detects the language, and builds it. HTTPS and a default `ondigitalocean.app` hostname come with the app.

The sample uses DigitalOcean's public `sample-nodejs` repo so the preset applies without placeholders. Point `git.repoCloneUrl` at your own public clone URL when you are ready. Use `github` (owner/repo) instead of `git` only when the DigitalOcean account has GitHub connected in the control panel -- `deployOnPush` needs that connection.

## When to Use

- A web app or API whose source is in a public Git repository
- Getting from source to a running URL without a container build pipeline
- Accounts that do not have GitHub/GitLab/Bitbucket linked yet (plain `git` clone URL)

## Key Configuration Choices

- **App name** (`appName`) -- 2-32 characters, unique in the DigitalOcean account. This is `spec.name` in Terraform.
- **Git clone URL** (`services[].git`) -- App Platform clones over HTTPS. No linked GitHub account required.
- **Instance size** (`instanceSizeSlug: basic-xxs`) -- smallest paid tier. Size slugs are free-form strings; new sizes work without a catalog change.
- **Single instance** (`instanceCount: 1`) -- raise this, or switch to `autoscaling`, for production traffic.
- **Environment variables** (`envs`) -- `plaintext` for ordinary values; `secret` for credentials (stored in App Platform's secret store).

## Related Presets

- **02-container-image** -- use when you already have a container image (Docker Hub, GHCR, or DigitalOcean Container Registry)
