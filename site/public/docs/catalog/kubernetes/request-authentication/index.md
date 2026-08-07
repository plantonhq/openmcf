---
title: "Request Authentication"
description: "Request Authentication deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesrequestauthentication"
---

# Request Authentication on Kubernetes

Defines an Istio RequestAuthentication: a namespaced policy that decides which JSON Web Tokens (JWTs) are accepted on your workloads. Use it to validate end-user or service tokens from one or more identity providers at the mesh proxy, extract the verified identity for authorization, and reject requests that present an invalid token.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **A RequestAuthentication policy** -- a namespaced Istio policy that tells the mesh how to validate JWTs presented to the workloads it selects: which issuers to trust, where to read the token from, and what to do with the verified claims.
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking.

## What It Controls

- **Trusted issuers** -- one JWT rule per issuer, each with its own signing keys (discovered via OpenID Connect, pinned to a JWKS URL, or supplied inline) and allowed audiences.
- **Where the policy applies** -- namespace-wide, narrowed to workloads by label selector, or attached to specific resources (Gateway, Service, ServiceEntry) via target references.
- **Token extraction** -- read the token from the standard Authorization header, custom headers (with an optional prefix), query parameters, or cookies.
- **Identity output** -- forward the original token upstream, emit the whole payload to a header, or copy individual verified claims (groups, roles, tenant) into request headers the application can read.

## Important Behavior

RequestAuthentication validates a token **when one is present**. A request carrying a valid token gets an authenticated principal; a request with an invalid token is rejected; a request with **no** token is allowed through. To actually require a token, pair this with an Authorization Policy that demands a request principal.

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **Istio installed** -- the Istio CRDs (Istio Base CRDs on Kubernetes) must be present and the Istio control plane (istiod) running. The policy is only enforced where istiod is active and the selected workloads are part of the mesh.
- **Target namespace exists** -- the policy is created in a specific namespace; reference an existing one or create it first.
- **Reachable signing keys** -- istiod must be able to fetch the issuer's JWKS (via discovery or the configured URL) to verify token signatures.

## Deploy

### Console

Open the deployment store, find **Request Authentication on Kubernetes**, and click **Deploy**. The creation wizard walks you through the namespace, the scope (namespace-wide, a workload selector, or target references), and the JWT rules, with guidance at each step.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesRequestAuthentication
metadata:
  name: jwt-auth
  org: acme-corp
  env: prod
spec:
  namespace:
    value: prod-apps
  jwt_rules:
    - issuer: https://accounts.google.com
      audiences:
        - my-app.acme.com
```

```shell
planton apply -f request-authentication.yaml
```

This validates Google-issued JWTs for every workload in the `prod-apps` namespace.

## Key Configuration

- **Namespace** -- the namespace the policy is created in. It is fixed once created; the policy's scope is defined relative to it.
- **Scope** -- namespace-wide by default, or narrowed to a **workload selector** (pod labels) or **target references** (a Gateway, Service, or ServiceEntry). At most one of a selector and target references can be set. Waypoint proxies require target references.
- **JWT rules** -- one per trusted issuer. Each names the **issuer** (matching the token's `iss` claim), optional **audiences**, the **JWKS source** (OIDC discovery, a URL, or inline keys), the **token locations** (headers, query params, cookies), and **claim output** (forward the token, emit the payload, or map claims to headers).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Description |
|------------|-------------|
| `namespace` | The namespace the policy is created in. Reference an existing Namespace on Kubernetes or supply the name directly. |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources and operators can reference:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `request_authentication_name` | Name of the created RequestAuthentication (equals `metadata.name`) | Ordering resources that depend on the policy being in place |
| `namespace` | The namespace the policy was created in | Confirming where the policy applies |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

- **Validate JWTs across a namespace** -- accept tokens from a single trusted issuer for every workload in a namespace. The canonical starting point. Start from the **namespace-jwt-validation** preset.
- **Workload JWT validation with claim forwarding** -- validate tokens for one selected workload, read the token from a custom header, restrict it to an audience, and forward a verified claim to the backend as an HTTP header. Start from the **workload-jwt-with-claim-headers** preset.

## Works With

RequestAuthentication is part of the Istio security family. It requires the Istio Base CRDs on Kubernetes and a running Istio control plane, and it pairs naturally with Authorization Policy -- RequestAuthentication establishes *who the caller is* (a verified JWT principal), and AuthorizationPolicy decides *whether that principal may proceed*. To order the policy after the workload or resource it protects within an infra chart, express the dependency through `metadata.relationships`.
