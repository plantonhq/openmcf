---
title: "Ingress NGINX"
description: "Ingress NGINX deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesingressnginx"
---

# Ingress NGINX on Kubernetes

Installs the ingress-nginx controller — the cluster's HTTP and HTTPS entry
point. It watches Ingress resources of its class and programs an embedded NGINX
to route external traffic to Services.

## What this component is, and is not

This installs the **controller only**. Routing rules are separate first-class
resources: create Ingress objects naming this controller's published class.
Certificates come from cert-manager through each Ingress's TLS secret, or
cluster-wide through the default certificate below.

## Several controllers per cluster is the normal pattern

A public entry point and an internal one is the common split. Each instance gets
its own resource and its own **distinct ingress class** — release naming,
controller resource names and leader-election identity all derive from the
resource name, so instances never collide. Two controllers sharing a class both
program the same Ingresses and fight over them.

## The cloud integration surface

What the host cloud provisions is driven entirely by the controller Service's
**annotations**. The controller itself never calls a cloud API.

Which controller *reads* those annotations decides whether anything happens at
all. On EKS the built-in cloud controller handles
`service.beta.kubernetes.io/aws-load-balancer-type: "nlb"`, while the
`external` family belongs to the AWS Load Balancer Controller. Set the second
without that controller installed and the Service simply never receives an
address — no error surfaces anywhere, and the module's readiness wait times out.
Both paths are verified live.

Internal placement is likewise annotation-driven: GCP uses
`networking.gke.io/load-balancer-type`, Azure uses
`service.beta.kubernetes.io/azure-load-balancer-internal`, AWS uses the scheme
annotation.

## Deployment shapes

| Shape | Controller kind | Service | Networking |
|---|---|---|---|
| Managed cloud (default) | Deployment | LoadBalancer | Cluster network |
| Bare metal / edge | DaemonSet | NodePort | Host ports or host network |
| Internal only | Deployment | ClusterIP or internal LB | Cluster network |

Host networking and host ports are alternatives — host networking already binds
every listener on the node.

## Availability

An ingress controller is a single point of failure for everything behind it. Two
or more replicas gives zero-drop rollouts and node-failure tolerance, and above
one replica the chart adds a PodDisruptionBudget automatically. Autoscaling is
chart-managed and needs metrics-server; while it is on, it owns the replica
count.

Upstream recommends leaving the CPU limit unset: throttling an entry point turns
a traffic spike into cluster-wide latency.

## Two settings worth a second look

**Snippet annotations** let any Ingress author inject raw NGINX directives into
the shared controller. Upstream disabled them by default after CVE-2021-25742;
enable only where every Ingress author is trusted at controller level.

**The admission webhook** rejects a broken Ingress at apply time instead of at
NGINX reload time — which is the difference between one developer seeing an
error and every service behind the controller going down. The certgen Job is how
it bootstraps its certificate, so the only real reason to disable it is a
cluster policy forbidding hook Jobs.

## Works with

| Kind | Relationship |
|---|---|
| KubernetesNamespace | Where the controller runs |
| KubernetesIngress | The routing rules this controller serves, selected by class name |
| KubernetesCertificate | Supplies the cluster-wide default TLS certificate |
| KubernetesService | The backends Ingresses route to, and raw TCP/UDP targets |
| KubernetesPriorityClass | Ranks controller pods above ordinary workloads |
| KubernetesMetricsServer | Required for autoscaling utilization targets |

## Outputs

| Output | Description |
|---|---|
| `namespace` | The namespace the controller runs in |
| `release_name` | The Helm release name |
| `ingress_class_name` | The class an Ingress names to be served by this controller |
| `controller_service_name` | The controller's Service |
| `internal_service_name` | The second internal Service, when enabled |
| `load_balancer_ip` | The provisioned address, on providers exposing an IP |
| `load_balancer_hostname` | The provisioned address, on providers exposing a hostname |
