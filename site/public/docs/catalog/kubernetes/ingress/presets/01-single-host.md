---
title: "Single Host"
description: "This preset creates the simplest useful Ingress: one hostname routed to one Service backend on a prefix path covering everything (`/`). It is the standard \"expose this app on this domain\" shape."
type: "preset"
rank: "01"
presetSlug: "01-single-host"
componentSlug: "ingress"
componentTitle: "Ingress"
provider: "kubernetes"
icon: "package"
order: 1
---

# Single Host

This preset creates the simplest useful Ingress: one hostname routed to one Service backend on a prefix path covering everything (`/`). It is the standard "expose this app on this domain" shape.

## When to Use

- Exposing a single web application or API under one stable hostname
- The first Ingress for a new workload — add TLS (see the cert-manager preset) or more paths later
- Any case where all traffic for the host goes to one Service

## Key Configuration Choices

- **`path: /` with `path_type: prefix`** — matches every request path for the host; prefix is the default and the one matching semantics every controller implements identically
- **`ingress_class_name: nginx`** — pins the serving controller explicitly rather than relying on the cluster's default class; change it to whatever `kubectl get ingressclass` lists in your cluster
- **Backend by port number** — routes to port 8080 of the Service; if your Service names its ports, `port_name` is the more change-resilient reference
- **Same-namespace backend** — the Service must live in the Ingress's namespace (`web` here); this is a Kubernetes API constraint

The Ingress is created without waiting for a controller — `load_balancer_ip`/`load_balancer_hostname` outputs stay empty until one reconciles it.

## Values to Replace

| Value | Description |
|---|---|
| `app.example.com` | The public hostname this Ingress serves |
| `web` | Namespace of the Ingress AND the backend Service |
| `web-svc` | Name of the backend Service (or a KubernetesService reference) |
| `8080` | The Service port receiving traffic |

## Related Presets

- **02-tls-cert-manager** — the same shape with HTTPS via cert-manager
- **03-fanout-paths** — one host, multiple paths to different backends
