---
title: "Network Policy"
description: "Network Policy deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesnetworkpolicy"
---

# Kubernetes Network Policy

Deploys a Kubernetes NetworkPolicy — the in-cluster firewall — to a target cluster through a single declarative manifest, covering the complete `networking/v1` surface: pod and namespace selectors, IP blocks with exceptions, TCP/UDP/SCTP ports, named ports, and port ranges. The IaC module handles label merging, namespace resolution, and policy-type inference automatically.

## What Gets Created

When you deploy a KubernetesNetworkPolicy resource, Planton provisions:

- **NetworkPolicy** — a `networking/v1` NetworkPolicy selecting pods via `podSelector` and declaring the allowed ingress and egress traffic
- **Labels** — standard Planton tracking labels merged with any user-provided labels
- **Annotations** — user-provided annotations applied to the NetworkPolicy metadata

**Policies are additive allows — there is no deny rule.** A pod not selected by any policy accepts all traffic; once selected in a direction, only traffic allowed by some policy in that direction gets through. Isolation comes from selecting pods while allowing little.

## Prerequisites

- **A NetworkPolicy-enforcing CNI** — Calico, Cilium, or a cloud CNI with policy enforcement enabled. On clusters whose CNI ignores NetworkPolicies (including default kind clusters), the object exists but all traffic still flows
- **Kubernetes credentials** configured via environment variables or Planton provider config
- **A Kubernetes namespace** that already exists, or a `KubernetesNamespace` resource referenced from `spec.namespace` so both deploy in one run

## Quick Start

Create a file `network-policy.yaml` — this is the namespace-wide default-deny baseline:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesNetworkPolicy
metadata:
  name: default-deny-all
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesNetworkPolicy.default-deny-all
spec:
  name: default-deny-all
  namespace:
    value: backend
  pod_selector: {}
  policy_types:
    - ingress
    - egress
```

Deploy:

```shell
planton apply -f network-policy.yaml
```

The empty `pod_selector` selects ALL pods in the `backend` namespace; with both directions governed and no rules, all traffic to and from those pods is denied. Every subsequent policy is a targeted exception that only widens.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `spec.name` | `string` | Name of the NetworkPolicy (`metadata.name` in the cluster). | 1–253 characters, matches `^[a-z0-9]([a-z0-9.-]{0,251}[a-z0-9])?$` |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.namespace` | `StringValueOrRef` | `default` | Namespace the policy lives in — a NetworkPolicy governs only pods in its own namespace. Accepts a literal name (`{ value: my-namespace }`) or a reference to a `KubernetesNamespace` resource. |
| `spec.pod_selector` | `LabelSelector` | empty | The pods this policy applies to. **An empty selector selects ALL pods in the namespace** — the default-deny building block. Target one Planton workload with `match_labels: {app: <workload-name>}`. |
| `spec.policy_types` | `list(ingress\|egress)` | inferred | The governed directions. When omitted, Kubernetes infers: ingress always, egress only when egress rules are present. Egress-only and deny-all-egress policies MUST set it explicitly. |
| `spec.ingress_rules` | `list(rule)` | `[]` | Allow rules for inbound traffic. Each rule has `from` (peers, OR'd) and `ports` (OR'd); peers AND ports must both match. Empty list with ingress governed = deny all inbound. |
| `spec.egress_rules` | `list(rule)` | `[]` | Allow rules for outbound traffic. Same shape with `to`. Denying all egress also blocks DNS — pair with an allow-DNS rule. |
| `spec.labels` / `spec.annotations` | `map<string, string>` | `{}` | Merged with standard Planton labels / applied to the object. |

### Peers

Each entry in `from`/`to` is one of three forms:

| Form | Meaning |
|------|---------|
| `pod_selector` alone | Pods matching the selector in the policy's OWN namespace |
| `namespace_selector` alone | ALL pods in namespaces matching the selector (e.g. the automatic `kubernetes.io/metadata.name` label) |
| both in ONE peer | AND: pods matching `pod_selector` in namespaces matching `namespace_selector` — different from two separate peers, which OR |
| `ip_block` | A CIDR with optional `except` carve-outs, for traffic outside the cluster. Cannot be combined with selectors in the same peer |

### Ports

| Field | Description |
|-------|-------------|
| `protocol` | `TCP` (default), `UDP`, or `SCTP` |
| `port` | A number (`"5432"`) or a named container port (`"metrics"`); omit to match all ports for the protocol |
| `end_port` | Inclusive upper bound of a range starting at a numeric `port` (e.g. `"30000"` + `32767`) |

## Examples

### Allow One Workload's Callers

Only the `frontend` workload's pods may reach `backend-api` pods, and only on port 8080. All other inbound traffic to `backend-api` is denied:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesNetworkPolicy
metadata:
  name: backend-api-ingress
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesNetworkPolicy.backend-api-ingress
spec:
  name: backend-api-ingress
  namespace:
    value: backend
  pod_selector:
    match_labels:
      app: backend-api
  policy_types:
    - ingress
  ingress_rules:
    - from:
        - pod_selector:
            match_labels:
              app: frontend
      ports:
        - protocol: TCP
          port: "8080"
```

### Allow Monitoring From Another Namespace

Prometheus pods in the `monitoring` namespace — and only those — may scrape the metrics port. Note the single AND'd peer:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesNetworkPolicy
metadata:
  name: allow-prometheus-scrape
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesNetworkPolicy.allow-prometheus-scrape
spec:
  name: allow-prometheus-scrape
  namespace:
    value: backend
  pod_selector:
    match_labels:
      app: backend-api
  policy_types:
    - ingress
  ingress_rules:
    - from:
        - namespace_selector:
            match_labels:
              kubernetes.io/metadata.name: monitoring
          pod_selector:
            match_labels:
              app: prometheus
      ports:
        - protocol: TCP
          port: metrics
```

### Lock Down Egress: DNS Plus One External CIDR

Deny all egress from `worker` pods except DNS and one external service range. `policy_types: [egress]` is explicit — an egress-only policy must say so, or inference adds ingress:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesNetworkPolicy
metadata:
  name: worker-egress
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesNetworkPolicy.worker-egress
spec:
  name: worker-egress
  namespace:
    value: backend
  pod_selector:
    match_labels:
      app: worker
  policy_types:
    - egress
  egress_rules:
    - to:
        - namespace_selector:
            match_labels:
              kubernetes.io/metadata.name: kube-system
      ports:
        - protocol: UDP
          port: "53"
        - protocol: TCP
          port: "53"
    - to:
        - ip_block:
            cidr: 203.0.113.0/24
      ports:
        - protocol: TCP
          port: "443"
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `networkPolicyName` | `string` | Name of the NetworkPolicy object as created in the cluster |
| `namespace` | `string` | Namespace the NetworkPolicy was created in |
| `policyTypes` | `string` | Governed directions as deployed — `"Ingress"`, `"Egress"`, or `"Ingress,Egress"` — including inferred directions when the spec omitted `policy_types` |

## Related Components

- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) — provides the target namespace; reference it from `spec.namespace` to deploy both in one run
- [KubernetesDeployment](/docs/catalog/kubernetes/deployment) — workload whose pods carry the `app: <workload-name>` label that `pod_selector` targets
- [KubernetesHttpEndpoint](/docs/catalog/kubernetes/kuberneteshttpendpoint) — external exposure; NetworkPolicies govern the pod-level traffic behind it
