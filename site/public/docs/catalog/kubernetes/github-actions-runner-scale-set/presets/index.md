---
title: "Presets"
description: "Ready-to-deploy configuration presets for GitHub Actions Runner Scale Set"
type: "preset-list"
componentSlug: "github-actions-runner-scale-set"
componentTitle: "GitHub Actions Runner Scale Set"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-repo-runners"
    rank: "01"
    title: "Repository runners preset"
    excerpt: "Self-hosted runners for one repository, scaled to zero when idle: a queued job creates an ephemeral runner pod, the pod runs exactly that job and is replaced. Workflows target the fleet by NAME —..."
  - slug: "02-org-docker-builds"
    rank: "02"
    title: "Organization Docker-builds preset"
    excerpt: "An organization-wide fleet that can build containers: dind mode runs a privileged Docker daemon beside every runner, so `docker build`, `docker run` and container-based actions work exactly as they..."
  - slug: "03-unprivileged-kubernetes-mode"
    rank: "03"
    title: "Unprivileged Kubernetes-mode preset"
    excerpt: "Container jobs WITHOUT privileged pods: the Kubernetes container hook runs each job container as its own pod, so security-restricted clusters (Pod Security `restricted`, regulated environments) still..."
---

# GitHub Actions Runner Scale Set Presets

Ready-to-deploy configuration presets for GitHub Actions Runner Scale Set. Each preset is a complete manifest you can copy, customize, and deploy.
