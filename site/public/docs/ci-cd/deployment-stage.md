---
title: "Deployment Stage"
description: "How pipelines turn built artifacts into running services — manifest resolution, Kustomize overlays, environment mapping, and local development."
icon: deployment
order: 45
tags:
  - Deployment
  - Kustomize
  - Service Hub
---

# Deployment Stage

After the build stage produces a container image or worker script, the deployment stage takes over. It resolves deployment manifests, creates a deployment task for each target environment, and provisions each one through a Stack Job. This page explains what happens during that process — how manifests are produced, how environments are matched, and how you can use the same system for local development.

For the high-level pipeline model (triggers, stages, cancellation, approval gates), see [Pipelines](/docs/ci-cd/pipelines). For supported platforms and the choice between Git-based and inline configuration, see [Deployment Targets](/docs/ci-cd/deployment-targets).

## How Manifests Are Resolved

The deployment stage produces cloud resource manifests through one of two paths, depending on how the service is configured:

**Git-based** (default): The pipeline reads deployment manifests from a `_kustomize` directory in your repository. During the build stage, a kustomize-build task processes each overlay directory and stores the merged manifests. The deployment stage reads those manifests and creates one deployment task per environment.

**Inline** (UI-based): Deployment targets are defined directly in the Service configuration. The deployment stage reads them from the service spec, substitutes template variables (such as `{{ .Image }}` for the built container image reference and `{{ .CommitSHA }}` for the Git commit), and creates one deployment task per target.

Both paths converge at the same point: a set of cloud resource manifests, one per environment, ready to be provisioned through Planton's infrastructure layer.

## The Kustomize Model

Kustomize is a YAML patching tool built into `kubectl`. Planton uses it because it works with plain YAML patches — no templating language, no chart packaging, no extra toolchain. You write a base manifest and patch it per environment using overlays. The output is deterministic and version-controlled.

### Directory Structure

The `_kustomize` directory lives in your repository (by default at the project root, configurable in service settings):

```
_kustomize/
  base/
    kustomization.yaml
    service.yaml
  overlays/
    local/
      kustomization.yaml
      service.yaml
    dev/
      kustomization.yaml
      service.yaml
    production/
      kustomization.yaml
      service.yaml
```

### Base Configuration

The base defines your service's default resource specification — the settings shared across all environments:

```yaml
# _kustomize/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - service.yaml
```

```yaml
# _kustomize/base/service.yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesDeployment
metadata:
  name: my-service
  org: my-org
spec:
  container:
    app:
      env:
        variables:
          PORT: "8080"
      resources:
        requests:
          cpu: 50m
          memory: 100Mi
        limits:
          cpu: 500m
          memory: 500Mi
      ports:
        - name: rest-api
          appProtocol: http
          networkProtocol: TCP
          servicePort: 80
          containerPort: 8080
          isIngressPort: true
  availability:
    minReplicas: 1
  version: main
```

### Environment Overlays

Each overlay directory patches the base with environment-specific values. An overlay inherits everything from the base and overrides only what differs:

```yaml
# _kustomize/overlays/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../base
patches:
  - path: service.yaml
```

```yaml
# _kustomize/overlays/production/service.yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesDeployment
metadata:
  name: my-service
  env: app-production
spec:
  container:
    app:
      resources:
        requests:
          cpu: 200m
          memory: 500Mi
        limits:
          cpu: "2"
          memory: 2Gi
      env:
        variables:
          LOG_LEVEL: warn
        secrets:
          DATABASE_PASSWORD: $secret/@production/database/password
  availability:
    minReplicas: 3
    horizontalPodAutoscaling:
      isEnabled: true
      maxReplicas: 10
```

The `metadata.env` field determines which Planton environment receives this deployment — not the overlay directory name. You can name directories however you like (`prod/`, `us-east/`, `blue/`), but the `env` field must match a real environment in your organization.

### Variables and Secrets in Overlays

Overlay manifests can reference organization-scoped secrets and variables using a substitution syntax:

- `$secret/<slug>` (or `$secret/@<env>/<slug>/<key>`) — Resolved at deployment time. The value is read just-in-time in the Runner from your secret backend, never stored in the platform's database.
- `$variables-group/<group>/<key>` — Resolved at deployment time. Variables support literal values or dynamic references to infrastructure outputs.

See [Secrets](/docs/secrets) for the full scoping model, backends, and lifecycle.

## The Local Overlay

The `local` overlay is a special case — it is never deployed. The deployment stage skips it entirely.

Its purpose is local development: generating `.env` files so developers can run services locally with the same variable and secret configuration used in deployed environments.

```bash
# Generate .env and .env_export files from the local overlay
planton service dot-env

# Override specific values for local testing
planton service dot-env --set API_KEY=test-key --set DEBUG=true
```

This keeps local development configuration version-controlled alongside deployment configuration, without any risk of it accidentally deploying to a real environment.

## How Deployment Tasks Execute

For each resolved manifest (excluding the local overlay), the deployment stage creates a deployment task:

1. **Environment matching** — If the service has deployment environment filters configured, only matching overlays produce tasks. See [Deployment Environments](/docs/ci-cd/deployment-environments).
2. **Ordering** — Tasks execute sequentially following the organization's promotion policy (for example: dev, then staging, then production).
3. **Manual gates** — If a deployment target requires manual approval, the pipeline pauses at that task until a team member approves or rejects. See [Pipelines](/docs/ci-cd/pipelines#manual-approval-gates).
4. **Stack Job creation** — Each task provisions the cloud resource manifest through a [Stack Job](/docs/infrastructure/stack-jobs). The Stack Job applies the infrastructure changes and reports completion.
5. **Failure handling** — If a deployment task fails, all subsequent tasks are cancelled. No partial rollouts across environments.

<!-- SCREENSHOT: Pipeline deployment stage
  Page: /{org}/service/{slug}/pipelines/runs/{id}
  Action: Show a pipeline run with the deployment stage in progress or completed
  Focus: The deployment stage section showing environment-specific deployment tasks
  Alt: Pipeline run detail showing the deployment stage with per-environment deployment tasks and their status
-->

## CLI Reference

```bash
# Initialize a new _kustomize directory
planton service kustomize init --new

# Initialize from an existing cloud resource
planton service kustomize init <resource-kind> <resource-id>

# Build and inspect the resolved manifests
planton service kustomize build

# Generate .env files from the local overlay
planton service dot-env

# Deploy using the _kustomize configuration
planton service deploy --project .
```

## Related Documentation

- [Pipelines](/docs/ci-cd/pipelines) — The full pipeline model including build stage and triggers
- [Deployment Targets](/docs/ci-cd/deployment-targets) — Supported platforms and the Git-based vs inline choice
- [Deployment Environments](/docs/ci-cd/deployment-environments) — Controlling which environments a service deploys to
- [Stack Jobs](/docs/infrastructure/stack-jobs) — How infrastructure changes are provisioned
- [Secrets](/docs/secrets) — Variable and secret management
