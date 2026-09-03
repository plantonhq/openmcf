---
title: "Self-Hosting"
description: "Run the entire Planton platform on your own Kubernetes cluster with two helm installs — batteries included, no external services, no configuration required"
icon: server
order: 70
tags:
  - Self-Hosting
  - Kubernetes
  - Helm
  - Operator
---

# Self-Hosting Planton

Planton runs entirely on your own Kubernetes cluster: the control plane, the web console, the identity server, the secrets manager, the databases, and an in-cluster runner — installed with two Helm commands (or one Planton CLI command), managed by a Kubernetes operator, reachable in minutes.

## Why self-host

Your infrastructure manifests, deployment history, secrets, and cloud credentials never leave your cluster. Self-hosting fits teams with data-residency requirements, air-gap-adjacent environments, or a platform team that wants Planton inside the same trust boundary as the infrastructure it manages.

## Install

Two commands install everything — the operator first, then the platform:

```bash
helm install planton-operator oci://ghcr.io/plantonhq/charts/planton-operator \
  --namespace planton --create-namespace

helm install planton oci://ghcr.io/plantonhq/charts/planton \
  --namespace planton
```

The first chart installs the Planton operator together with the `PlantonPlatform` definition it serves. The second creates one `PlantonPlatform` resource, and the operator reconciles the whole stack from that single resource: PostgreSQL, the workflow engine, the control plane, the console, the identity server, the secrets manager (OpenBAO, initialized and unsealed automatically), and the in-cluster runner. No license key, no admin account, no database, and no values file are required — a zero-values install works. The Planton CLI's self-hosted install runs both for you.

Watch it converge (typically 7–11 minutes):

```bash
kubectl get plantonplatform -n planton -w
```

Every component reports its own plain-language status on the resource — a stuck component names the problem and the fix in `kubectl describe plantonplatform`.

## First sign-in

The install's `NOTES` print the exact commands. In short:

```bash
kubectl -n planton port-forward svc/planton-gateway 8080:80
```

Open `http://localhost:8080`. The first person to open the console becomes the administrator: enter your email plus the cluster setup code (the `NOTES` print the `kubectl` command that reads it — holding cluster access IS the admin proof), and receive a one-time password.

## Publish at your own URL

The built-in port-forward front door works on every cluster with zero configuration. Going public is a ladder of one field at a time on the `PlantonPlatform` resource:

```yaml
spec:
  ingress:
    enabled: true                      # auto-derives a working URL from your ingress controller
    hostname: planton.example.com      # or serve your own domain
    tls:
      issuer: {name: letsencrypt, kind: ClusterIssuer}   # cert-manager HTTPS
```

With `enabled: true` alone, the operator derives a magic-DNS hostname from your ingress controller's published address — a working URL with zero DNS setup, unique to the platform's name and namespace. The desktop app offers this whole journey as a guided experience: pick a cluster from your kubeconfig, preflight it, choose the front door, and watch the install converge — driving the exact same chart underneath.

## Several Plantons, one cluster

Platforms are namespaced, and one operator serves the whole cluster — it watches every namespace. Teams can run separate Planton platforms side by side (staging and production, or one per team), each fully confined to its own namespace:

```bash
# The operator is already on the cluster; every platform is one more release.
helm install planton oci://ghcr.io/plantonhq/charts/planton \
  --namespace planton-team-b --create-namespace
```

Never install two operators (the operator itself refuses to start beside another and says so), and give each platform its own namespace. Two cluster-level facts are shared by design: build events (the CI event stream Tekton delivers) can feed only one platform per cluster, and all platforms ride the one installed operator version — each platform still pins its own `spec.version`.

## Upgrades and uninstall

Config changes are edits to the `PlantonPlatform` resource; the operator reconciles them. The platform version is `spec.version` on that resource — `helm upgrade planton --set platform.spec.version=<version>` rolls the platform with its data intact. The operator upgrades through its own chart, and that chart carries the `PlantonPlatform` definition with it, so the schema always matches the operator that reads it.

`helm uninstall` removes the platform's workloads; data volumes deliberately survive. To remove everything including data, delete the namespace. The cluster-scoped badge-verification grant (`<namespace>-<name>-control-plane-token-reviewer`) is the one manual cleanup step of a full teardown.

## Requirements

- Kubernetes 1.24+, amd64 nodes (arm64 works only under emulation, e.g. local Docker Desktop)
- A default StorageClass whose storage driver is actually installed (or pin one via `spec.storage.storageClassName`)
- 6 GiB+ allocatable memory is comfortable; smaller evaluation clusters work with resource floors
