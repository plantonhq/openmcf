# Service on Kubernetes

The stable network identity in front of a set of pods. A Service gives its
backends one durable name and address while the pods behind it come and go.

## Why a standalone Service

Workload kinds create a Service for their own pods already. This kind covers
everything else: exposing pods managed outside Planton, LoadBalancer and
NodePort exposure with cloud annotations, ExternalName aliases to hosts outside
the cluster, headless services for custom discovery, dual-stack addressing, and
selectorless services fronting manually-managed endpoints.

It is also the unit every other networking construct composes on — Ingress
backends, Gateway API route backends, DestinationRule hosts and BackendTLSPolicy
targets all point at a Service.

## Choosing a type

| Type | Reachable from | Builds on |
|---|---|---|
| ClusterIP (default) | Inside the cluster only | — |
| NodePort | A static port on every node | ClusterIP |
| LoadBalancer | A cloud-provisioned external address | NodePort |
| ExternalName | Nothing — it is a DNS alias | — |

ExternalName is a different mechanism, not a fourth exposure level: cluster DNS
returns a CNAME and no traffic passes through Kubernetes. Selectors, cluster
IPs, external IPs and ports are all rejected on it.

**Headless** is a modifier rather than a type. It removes the virtual IP so DNS
returns pod addresses directly — the tool for StatefulSet peer discovery and for
clients that must reach each pod individually. It cannot combine with NodePort
or LoadBalancer, which need a virtual IP to build on.

## The selector is the contract

Traffic reaches a pod only if that pod carries **every** label in the selector.
The most common failure here is silent: a selector matching nothing produces a
healthy-looking Service with zero endpoints, so DNS resolves and connections
hang.

Planton workloads stamp `app: <workload-name>` on their pods and export the full
set as their `selector_labels` output, so that single pair is enough to select
one.

An empty selector is a real posture, not an omission — a selectorless Service
whose endpoints are managed by hand or by a controller.

## Ports

Each port maps what clients connect to onto what the pods listen on. Naming the
target port rather than numbering it is the resilient choice: the workload can
change its listening port without the Service changing. Once a Service exposes
more than one port, every port needs a unique name.

## Cloud behaviour comes from annotations

For LoadBalancer services the annotation set is the portable way to tune the
provisioned balancer — AWS NLB selection, internal placement on GCP and Azure,
external-dns hostnames. Which controller answers matters: on EKS the built-in
cloud controller handles `aws-load-balancer-type: "nlb"`, while the `external`
family is handled only by the AWS Load Balancer Controller. With those set and
that controller absent, the Service never receives an address and no error
surfaces.

The deprecated `loadBalancerIP` field is deliberately omitted — upstream
deprecated it as under-specified, and every cloud pins an address through its
own annotation instead.

## Works with

| Kind | Relationship |
|---|---|
| KubernetesNamespace | Placement (optional — omitted means the `default` namespace) |
| KubernetesDeployment, KubernetesStatefulSet | The workloads whose pods it selects |
| KubernetesIngress | Routes HTTP traffic to this Service |
| KubernetesHttpRoute, KubernetesGrpcRoute | Gateway API backends pointing here |
| KubernetesBackendTlsPolicy | Secures the gateway-to-Service hop |
| KubernetesDestinationRule | Applies traffic policy to this host |
| KubernetesNetworkPolicy | Governs which pods may reach the selected pods |

## Outputs

| Output | Description |
|---|---|
| `service_name` | The Service object's name in the cluster |
| `namespace` | The namespace it was created in |
| `type` | The service type as deployed |
| `cluster_ip` | The assigned virtual IP — empty for headless and ExternalName |
| `load_balancer_ip` | The provisioned address, on providers exposing an IP |
| `load_balancer_hostname` | The provisioned address, on providers exposing a hostname |
| `kube_endpoint` | The in-cluster DNS endpoint other resources connect to |
| `port_forward_command` | A ready-to-run command for reaching it from a laptop |
