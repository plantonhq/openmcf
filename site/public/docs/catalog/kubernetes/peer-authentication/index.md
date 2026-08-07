---
title: "Peer Authentication"
description: "Peer Authentication deployment documentation"
icon: "package"
order: 100
componentName: "kubernetespeerauthentication"
---

# Peer Authentication on Kubernetes

Defines an Istio PeerAuthentication: a namespaced policy that controls whether incoming traffic to your workloads must arrive over a mutual-TLS (mTLS) tunnel. Use it to require encrypted, authenticated service-to-service traffic across a namespace, tighten a single workload, or carve out a specific port that must stay plaintext.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **A PeerAuthentication policy** -- a namespaced Istio policy that tells the mesh how to treat incoming connections to the workloads it selects: require mTLS (Strict), accept either plaintext or mTLS (Permissive), refuse mTLS (Disable), or inherit the surrounding default (Unset).
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking.

## What It Controls

- **Namespace-wide mTLS** -- with no workload selector, the policy sets the default for every workload in its namespace. In the mesh root namespace, it becomes the mesh-wide default.
- **Targeted workloads** -- with a selector, it applies only to the pods whose labels match, overriding any looser namespace default for them.
- **Per-port exceptions** -- individual workload ports can override the workload-level mode, e.g. keep a metrics or health-check port plaintext while everything else is Strict.

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **Istio installed** -- the Istio CRDs (Istio Base CRDs on Kubernetes) must be present and the Istio control plane (istiod) running. The policy is only enforced where istiod is active and the selected workloads are part of the mesh.
- **Target namespace exists** -- the policy is created in a specific namespace; reference an existing one or create it first.
- **Callers on the mesh before Strict** -- enabling Strict rejects plaintext, so confirm every client of the selected workloads connects over mTLS first. Permissive is the safe migration mode.

## Deploy

### Console

Open the deployment store, find **Peer Authentication on Kubernetes**, and click **Deploy**. The creation wizard walks you through the namespace, the optional workload selector, the mTLS mode, and any per-port overrides, with guidance at each step.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesPeerAuthentication
metadata:
  name: default
  org: acme-corp
  env: prod
spec:
  namespace:
    value: prod-apps
  mtls:
    mode: STRICT
```

```shell
planton apply -f peer-authentication.yaml
```

This requires mutual TLS for every workload in the `prod-apps` namespace.

## Key Configuration

- **Namespace** -- the namespace the policy governs. It is fixed once created; the policy's scope is defined relative to it.
- **Workload selector** -- optional pod labels that narrow the policy to specific workloads. Leave it empty to cover the whole namespace. It is a plain label match, not a reference to another resource, so it creates no ordering dependency.
- **mTLS mode** -- **Strict** requires an mTLS tunnel, **Permissive** accepts plaintext or mTLS (use while migrating onto the mesh), **Disable** forces plaintext, and **Unset** inherits the namespace or mesh default.
- **Per-port overrides** -- override the mode for specific workload ports (the container's ports, not Service ports). Overrides only take effect when a workload selector is set.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Description |
|------------|-------------|
| `namespace` | The namespace the policy is created in. Reference an existing Namespace on Kubernetes or supply the name directly. |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources and operators can reference:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `peer_authentication_name` | Name of the created PeerAuthentication (equals `metadata.name`) | Ordering resources that depend on the policy being in place |
| `namespace` | The namespace the policy was created in | Confirming where the policy applies |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

- **Namespace-wide Strict mTLS** -- require mutual TLS for every workload in a namespace. The canonical hardening baseline. Start from the **namespace-strict-mtls** preset.
- **Strict workload with a plaintext port** -- require mTLS for one selected workload while exempting a single port (e.g. a metrics scrape or health check that a non-mesh client probes). Start from the **workload-strict-with-plaintext-port** preset.

## Works With

PeerAuthentication is part of the Istio security family. It requires the Istio Base CRDs on Kubernetes and a running Istio control plane, and it commonly pairs with Authorization Policy (who may call a workload) and Request Authentication (end-user JWT validation). To order the policy after the workload it protects within an infra chart, express the dependency through `metadata.relationships`.
