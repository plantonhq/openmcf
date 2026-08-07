# Authorization Policy on Kubernetes

Defines an Istio AuthorizationPolicy: a namespaced resource that enforces access control on your mesh workloads. It decides whether a request is allowed, denied, or audited based on who the request is from, what it is doing, and additional conditions -- or hands the decision to an external authorizer. This is how you express zero-trust, request-level access rules inside the mesh.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **An AuthorizationPolicy** -- a namespaced Istio policy that applies an action (Allow, Deny, Audit, or Custom) to the requests matched by its rules, scoped to selected workloads or to the whole namespace.
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking.

## What It Controls

- **Scope** -- which workloads the policy protects: every workload in the namespace, a set of pods matched by label, or specific resources (Gateway, Service, ServiceEntry) attached by reference.
- **Action** -- the verdict on a matched request: Allow (permit only what matches), Deny (reject what matches), Audit (log without changing the decision), or Custom (delegate to an external authorizer).
- **Rules** -- which requests the action applies to, expressed as sources (`from`), operations (`to`), and conditions (`when`).

## Important Behavior

The action and the rules work together, and the empty cases are consequential. With the **Allow** action, once any allow policy applies to a workload, only requests that match a rule are permitted -- so a policy with **no rules denies all traffic** to the selected workloads. With **Deny**, a policy with no rules does nothing. Across all policies on a workload, Custom is evaluated first, then Deny, then Allow, so a Deny always overrides an Allow. At most one of a workload selector or target references may be set; both omitted means namespace-wide (or mesh-wide from the Istio root namespace).

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **Istio installed** -- the Istio CRDs (Istio Base CRDs on Kubernetes) must be present and the Istio control plane (istiod) running. The policy is only enforced where istiod is active.
- **Target namespace exists** -- the policy is created in a specific namespace; reference an existing one or create it first.

## Deploy

### Console

Open the deployment store, find **Authorization Policy on Kubernetes**, and click **Deploy**. The creation wizard walks you through the namespace, the scope, the action, and the matching rules, with guidance at each step.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesAuthorizationPolicy
metadata:
  name: allow-frontend-to-api
  org: acme-corp
  env: prod
spec:
  namespace:
    value: prod-apps
  selector:
    matchLabels:
      app: api
  action: ALLOW
  rules:
    - from:
        - source:
            principals:
              - cluster.local/ns/prod-apps/sa/frontend
      to:
        - operation:
            methods: ["GET", "POST"]
            paths: ["/api/*"]
```

```shell
planton apply -f authorization-policy.yaml
```

This allows only the `frontend` service account to call `GET`/`POST` on `/api/*` of the `api` workloads in `prod-apps`; everything else to those workloads is denied.

## Key Configuration

- **Namespace** -- the namespace the policy is enforced in. It is fixed once created; the policy is registered relative to it.
- **Scope** -- namespace-wide, a **workload selector** (pod labels), or **target references** (Gateway/Service/ServiceEntry; required for waypoints). At most one of selector and target references is set.
- **Action** -- **Allow** (the default), **Deny**, **Audit**, or **Custom**. Custom delegates to a named MeshConfig extension provider.
- **Rules** -- each rule combines **from** (source identities, namespaces, or IPs), **to** (hosts, ports, methods, paths), and **when** (attribute conditions such as a JWT claim). A request matches the policy if it matches any rule; within a rule it must match any source, any operation, and all conditions.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Description |
|------------|-------------|
| `namespace` | The namespace the policy is created in. Reference an existing Namespace on Kubernetes or supply the name directly. |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources and operators can reference:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `authorization_policy_name` | Name of the created AuthorizationPolicy (equals `metadata.name`) | Ordering resources that depend on the policy being in place |
| `namespace` | The namespace the policy was created in | Confirming where the policy applies |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

- **Require a valid JWT principal** -- allow only requests that carry a validated request principal, locking down a workload to authenticated callers. Pair with a Request Authentication that validates the token. Start from the **require-jwt-principal** preset.
- **Front a workload with external authorization** -- use the Custom action to delegate the decision to an external authorizer (for example, an ext-authz service) on ingress traffic. Start from the **custom-ext-authz-ingress** preset.

## Works With

AuthorizationPolicy is part of the Istio security family. It requires the Istio Base CRDs on Kubernetes and a running Istio control plane. It pairs naturally with a **Request Authentication** (which validates JWTs so rules can match `request principals` and JWT-claim conditions) and a **Peer Authentication** (which enforces mTLS so rules can match peer `principals`, `namespaces`, and `service accounts`). To order the policy after the workloads it protects within an infra chart, express the dependency through `metadata.relationships`.
