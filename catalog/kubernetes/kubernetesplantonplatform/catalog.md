# Planton Platform

Declares a complete self-hosted Planton platform — control plane, web console, identity server, PostgreSQL, cache, workflow engine, secrets manager, and an in-cluster deployment runner — as a `PlantonPlatform` custom resource the Planton operator reconciles. Zero-config by design: `version` is the only required choice, the built-in gateway serves console + API + sign-in over a single port-forward, and the first console visitor becomes the admin. Several platforms share one cluster, each in its own namespace with its own URL, identity, and databases — all served by one operator.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **PlantonPlatform CR** -- the one declaration; the OPERATOR then creates the platform from it (workloads, Services, Secrets, volumes — all in the platform's namespace, all named from this resource's name)
- **Namespace** (optional) -- created with standard governance labels when `create_namespace` is true

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.

### Cluster Prerequisites

- **The Planton operator** -- a deployed KubernetesPlantonOperator (one per cluster serves every platform).
- For `ingress`: an ingress controller; for `ingress.tls.issuer`: cert-manager.
- A default StorageClass that can actually provision volumes (or set `storage.storage_class_name`) — the operator verifies this before deploying and its status explains any storage problem in plain language.

## Deploy

### Console

Open the deployment store, find **Planton Platform**, and click **Deploy**. The creation wizard walks you through placement, the version, exposure (port-forward vs ingress/TLS), storage, identity and bootstrap seeding, the runner's cloud identity, and the opt-in components. Start from the **Zero Config** preset in the Presets tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesPlantonPlatform
metadata:
  name: planton
  org: acme-corp
  env: platform
spec:
  namespace:
    value: planton
  create_namespace: true
  version: v0.0.40-selfhosted-preview
```

```bash
planton apply -f planton.yaml
```

Watch it come up (`kubectl get plantonplatforms -A` — phase, version, URL), then use the `port_forward_command` output to open the door and the `setup_code_command` output for the first-visit setup page.

## Destroy

Deleting the resource tears the whole platform down — databases included. Every operator-created object is owner-referenced to the declaration, so Kubernetes garbage collection completes the teardown even when the operator itself is already gone, and the database layer removes its volumes and credentials together. Build caches and workflow volumes can survive in the namespace; when this resource owned the namespace (`create_namespace: true`), its deletion sweeps them.
