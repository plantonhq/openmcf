---
title: "API Gateway"
description: "API Gateway deployment documentation"
icon: "package"
order: 100
componentName: "ociapigateway"
---

# API Gateway on OCI

Deploys an Oracle Cloud Infrastructure API Gateway bundled with a single API deployment. The gateway provides the managed network endpoint (public or private) while the deployment defines routes, backends (HTTP services, OCI Functions, or stock responses), optional JWT authentication, CORS, rate limiting, and per-route authorization. Integrates with Planton's Provider Connections for OCI credential management and ValueFromRef for wiring to compartments, subnets, and security groups.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **API Gateway** -- the managed network endpoint in the specified compartment and subnet with configurable endpoint type (public or private), optional TLS certificate, and optional NSG bindings
- **API Deployment** -- a deployment on the gateway with a path prefix, route definitions mapping URL paths to backends, and optional request policies (JWT authentication, CORS, rate limiting). The deployment is always created in the same compartment as the gateway.
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the gateway

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the gateway and deployment in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A subnet for the gateway. For public gateways, this must be a public subnet. Provide the subnet OCID directly or reference an OciSubnet Cloud Resource via ValueFromRef. Changing the subnet forces recreation.
- A JWKS endpoint URL or static public keys (for JWT authentication) -- required only when enabling token-based authentication on the deployment.
- OCI Functions function OCIDs (for Functions backends) -- required only when routing requests to serverless functions.
- An OCI Certificates service certificate OCID (optional) -- for TLS termination on public gateways.

## Deploy

### Console

Open the deployment store, find **API Gateway on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Public HTTP Proxy** preset in the [Presets](#presets) tab to pre-populate a public gateway proxying requests to an HTTP backend with CORS and access logging.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciApiGateway
metadata:
  name: my-api
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  endpointType: endpoint_type_public
  subnetId:
    value: "ocid1.subnet.oc1..example"
  deployment:
    pathPrefix: /api/v1
    routes:
      - path: /health
        methods:
          - GET
        backend:
          type: stock_response
          status: 200
          body: '{"status":"ok"}'
      - path: /{path*}
        methods:
          - GET
          - POST
          - PUT
          - DELETE
        backend:
          type: http
          url: "https://backend.example.com:8080"
```

```shell
planton apply -f api-gateway.yaml
```

This creates a public API gateway with a health check route returning a stock response and a catch-all route proxying to an HTTP backend. No JWT authentication, CORS, or rate limiting is configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the gateway to a compartment, subnet, and security group deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: api-compartment
      fieldPath: status.outputs.compartmentId
  subnetId:
    valueFrom:
      kind: OciSubnet
      name: public-subnet
      fieldPath: status.outputs.subnetId
  networkSecurityGroupIds:
    - valueFrom:
        kind: OciSecurityGroup
        name: gateway-nsg
        fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, subnet, and security group first, then provisions the API gateway with the resolved values.

## Key Configuration

These are the most important decisions when configuring an API gateway. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Endpoint type** -- Set `endpointType` to `endpoint_type_public` for internet-facing APIs or `endpoint_type_private` for VCN-internal APIs. Public gateways require a public subnet. The endpoint type is immutable after creation.

**Backend types** -- Each route maps to one of three backend types: `http` for proxying to upstream services (with configurable connect/read/send timeouts and TLS verification), `oracle_functions` for invoking OCI Functions by OCID, or `stock_response` for returning a fixed HTTP status and body (health checks, maintenance pages). Routes are evaluated in order; first match wins.

**JWT authentication** -- Configure `deployment.requestPolicies.authentication` to validate JWT tokens before requests reach backends. Specify `issuers` (allowed iss claims), `audiences` (allowed aud claims), and `publicKeys` sourced from a remote JWKS endpoint or static PEM/JWK keys. Set `isAnonymousAccessAllowed: true` to allow unauthenticated requests to reach routes that enforce authorization individually.

**Per-route authorization** -- When deployment-level authentication is enabled, each route can set `authorization.type` to `anonymous` (no token required), `authentication_only` (valid token required, no scope check), or `any_of` with `allowedScope` (token must contain at least one listed OAuth2 scope). This enables mixed public/private API surfaces on a single gateway.

**Rate limiting and CORS** -- Configure `deployment.requestPolicies.rateLimiting` with `rateInRequestsPerSecond` and `rateKey` (`client_ip` for per-IP limits, `total` for aggregate). Configure `deployment.requestPolicies.cors` with `allowedOrigins`, `allowedMethods`, and credential/preflight cache settings for browser-based API consumers.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciSubnet** | `subnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `networkSecurityGroupIds` | `status.outputs.networkSecurityGroupId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `gatewayId` | OCID of the API gateway | Monitoring, IAM policy scoping, resource management |
| `hostname` | Hostname assigned to the gateway by OCI | DNS CNAME records (OciDnsRecord), client configuration |
| `deploymentEndpoint` | Full endpoint URL (gateway hostname + deployment path prefix) | API client base URL, integration testing |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Public HTTP proxy** -- A public gateway proxying all requests to an HTTP backend with CORS configured for a frontend domain, access and execution logging enabled, and a health check stock-response route. The standard pattern for exposing backend services as a managed API. Start from the **Public HTTP Proxy** preset.

**JWT-authenticated API** -- A public gateway with JWT authentication via remote JWKS, scope-based per-route authorization (anonymous health check, authenticated user routes, admin routes requiring admin scope), CORS, and rate limiting. Start from the **JWT Authenticated API** preset.

**Private Functions backend** -- A private VCN-internal gateway routing requests to OCI Functions for serverless API backends. No authentication or CORS (internal traffic only). Access and execution logging enabled. Start from the **Private Functions Backend** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this gateway and its deployment
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the subnet where the gateway endpoint is deployed
- [**Network Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides network security rules applied to the gateway