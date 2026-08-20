---
title: "Simpler Runner Setup"
date: 2026-03-05
category: improvement
tags:
  - runner
  - cli
excerpt: "Runner deployment now uses a single credentials file instead of multiple certificates and flags, and recovers faster after restarts."
author:
  - name: Swarup Donepudi
    title: Founder
---

Setting up and deploying a runner used to involve passing multiple certificates, flags, and environment variables separately. You had to manage a CA certificate, an agent certificate, a private key, a channel identifier, a tunnel host, and a tunnel port — each as its own configuration value. That entire model has been replaced by a single credentials file that contains everything the runner needs.

## One File, Every Deployment Target

When you generate credentials for a runner, the platform now produces a single `credentials.json` file. This file bundles the runner's identity, certificates, API key, and platform endpoint into one document. Whether you deploy to Kubernetes, AWS ECS, GCP Cloud Run, or Azure Container Apps, the setup follows the same pattern: pass one file.

**Kubernetes via Helm:**

```bash
helm install my-runner planton-runner \
  --set-file credentials.content=credentials.json
```

The Helm chart creates a Kubernetes Secret from the file and mounts it into the runner pod. If you use sealed-secrets, external-secrets, or a GitOps pipeline, you can reference an existing Secret instead:

```bash
helm install my-runner planton-runner \
  --set credentials.existingSecret=my-runner-creds
```

**Cloud deployments via the CLI:**

```bash
planton runner deploy my-runner --target aws-ecs
planton runner deploy my-runner --target gcp-cloud-run
planton runner deploy my-runner --target azure-container-apps
```

Each deployer stores the credentials file as a single secret in your cloud provider's native secret manager (AWS Secrets Manager, GCP Secret Manager, or Azure Container App secrets) and injects it into the runner container at startup.

Credential rotation is also simpler now — regenerate credentials, run `helm upgrade` (Kubernetes) or redeploy (cloud), and the runner picks up the new file. For Kubernetes, the Helm chart automatically triggers a rolling restart when the credentials content changes.

## Easier Local Startup

Running a runner on your development machine is now more convenient. If you've previously generated credentials, the runner can discover them automatically:

```bash
# Start by name and org — no file path needed
planton-runner start my-runner --org acme

# Start by name — runner searches across all orgs
planton-runner start my-runner

# Start with no arguments — interactive selection
planton-runner start
```

When you run `planton-runner start` without arguments in a terminal, it scans your local credential store, lists all saved runners with their org and connection details, and lets you pick one interactively. On headless systems (containers, CI), it exits with a clear error message explaining exactly how to provide credentials.

Platform operators also have access to `planton-runner operator` subcommands for inspecting and managing local credentials:

```bash
planton-runner operator config list     # Show all saved runners
planton-runner operator config show x   # Inspect credentials (sensitive fields masked)
planton-runner operator config delete x # Remove saved credentials
```

## Faster Recovery After Restarts

Previously, when a runner restarted, the platform could take up to five minutes to detect that the old connection was stale. During that window, every operation routed to that runner would fail with "Runner Unreachable." Now, the platform detects stale connections on the first failed request and immediately establishes a fresh connection. Recovery time dropped from up to five minutes to a single request.

The runner also handles transient network issues more gracefully at startup. If the control plane is temporarily unreachable when the runner boots, it retries its identity verification with exponential backoff instead of failing immediately. And once running, the runner continuously monitors its tunnel connection health, logging clear state transitions — "healthy," "connection lost," "connection restored" — so you can see exactly what's happening without parsing low-level network logs.
