---
title: "ClusterIP for an Application"
description: "This preset creates the default kind of Service: a cluster-internal virtual IP and DNS name in front of an application's pods. Anything inside the cluster reaches the app at..."
type: "preset"
rank: "01"
presetSlug: "01-cluster-ip-app"
componentSlug: "service"
componentTitle: "Service"
provider: "kubernetes"
icon: "package"
order: 1
---

# ClusterIP for an Application

This preset creates the default kind of Service: a cluster-internal virtual IP and DNS name in front of an application's pods. Anything inside the cluster reaches the app at `my-app.<namespace>.svc.cluster.local` (or just `my-app` from the same namespace), while the pods behind it come and go freely.

## When to Use

- Exposing pods managed outside Planton (Helm releases, operators, raw manifests) under a stable in-cluster name
- Giving an Ingress or Gateway a Service backend to point at
- Note: Planton workload kinds (KubernetesDeployment, KubernetesStatefulSet) already create a Service for their own pods — reach for this standalone kind when you need an *additional* or differently-shaped Service

## Key Configuration Choices

- **`type: cluster_ip`** — internal-only; the Service is unreachable from outside the cluster (the safest default). It is also the spec default, so the line is optional
- **`selector`** — traffic goes to pods whose labels match ALL listed pairs. Planton workloads stamp `app: <workload-metadata-name>` on their pods, so a single `app` entry selects a Planton-managed workload's pods
- **`ports`** — clients connect to port `80`; traffic lands on container port `8080`. `target_port` accepts a number (`"8080"`) or a named container port (`"http"`); omit it to reuse `port`
- **Port name** — optional for a single port, required and unique when there are several; named ports let consumers say `http` instead of a number

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace — must be the namespace the selected pods run in (a Service only selects pods in its own namespace) | Your namespace management |
| `<your-app-name>` | Value of the `app` label on the pods to route to; for a Planton workload this is the workload's `metadata.name` | `kubectl get pods --show-labels`, or the workload manifest |

## Related Presets

- **02-public-load-balancer** — expose the same pods to the internet
- **03-headless-statefulset** — per-pod DNS instead of a single virtual IP
