# Cloudflare Zero Trust Gateway Policy

One Gateway rule: a filter expression over employee traffic (DNS, HTTP, or network) plus the action to take when it matches. `enabled` defaults to false at Cloudflare -- a policy authored without `enabled: true` deploys disabled and filters nothing.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Gateway policy** -- one `cloudflare_zero_trust_gateway_policy` on the account
- **Rule settings** -- always sent. An empty object when the spec configures nothing (the provider's own drift workaround); otherwise the full settings tree (block page, session check, DNS resolvers, isolation controls, and kin)

## Prerequisites

- **A Cloudflare account with Zero Trust enabled** -- the organization (team name) must already exist
- **A Cloudflare API token** with Account → Zero Trust → Edit
- **The matching add-on** for `isolate` (Browser Isolation) or `egress` (dedicated egress IPs). Those actions fail the apply on an account that lacks the entitlement
- **A CloudflareZeroTrustList** if the `traffic` expression references a list by ID
- **A CloudflareZeroTrustTunnelVirtualNetwork** if `ruleSettings.dnsResolvers` sets `vnetId`

## Quick Start

A DNS block with an explicit enable and precedence:

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustGatewayPolicy
metadata:
  name: block-gambling
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: block-gambling-domains
  action: block
  filter: dns
  enabled: true
  precedence: 1000
  traffic: any(dns.domains[*] == "gambling.example.com")
```

```shell
planton apply -f gateway-policy.yaml
```

DNS queries for that domain are blocked. Without `enabled: true` this rule would deploy and do nothing.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | Owning account. Gateway policies are account-scoped. | Required. 32-character hex. |
| `name` | string | Display name. | Required, min length 1. |
| `action` | string | What happens on a match. | Required. One of `on`, `off`, `allow`, `block`, `scan`, `noscan`, `safesearch`, `ytrestricted`, `isolate`, `noisolate`, `override`, `l4_override`, `egress`, `resolve`, `quarantine`, `redirect`. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `filter` | string | Cloudflare infers from the action | `http`, `dns`, `l4`, `egress`, or `dns_resolver`. The module sends `[filter]`. |
| `enabled` | optional bool | Cloudflare's false | Set `true` or the policy is inert. |
| `precedence` | optional int64 | Cloudflare assigns one | Lower evaluates earlier. |
| `description` | string | unset | Purpose. |
| `traffic` | string | match-all | Wirefilter over the request. The API reformats it; adopt the returned form if the plan drifts. |
| `identity` | string | match-all | Wirefilter over the user. |
| `devicePosture` | string | match-all | Wirefilter over device posture. |
| `expiration` | object | unset | `expiresAt` (RFC3339, required in-object) and optional `duration` (≥5). |
| `schedule` | object | unset | `mon`..`sun` time-interval strings plus `timeZone`. |
| `ruleSettings` | object | empty object (always sent) | Action-specific settings. Avoid `addHeaders` / `overrideIps` if you need a clean re-plan -- they drift on first apply at provider v5.23.0. `dnsResolvers.ipv4[].vnetId` can reference a CloudflareZeroTrustTunnelVirtualNetwork. |

## Examples

### DNS block

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustGatewayPolicy
metadata:
  name: block-gambling
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: block-gambling-domains
  action: block
  filter: dns
  enabled: true
  precedence: 1000
  traffic: any(dns.domains[*] == "gambling.example.com")
  ruleSettings:
    blockPageEnabled: true
    blockReason: Blocked by company policy
```

### HTTP allow with a session check

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustGatewayPolicy
metadata:
  name: allow-intranet
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: allow-intranet-with-reauth
  action: allow
  filter: http
  enabled: true
  precedence: 2000
  traffic: http.request.host == "intranet.example.com"
  ruleSettings:
    checkSession:
      enforce: true
      duration: 24h
    blockPage:
      targetUri: https://block.example.com
      includeContext: true
```

### DNS override with a virtual-network resolver

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustGatewayPolicy
metadata:
  name: resolve-internal
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: resolve-internal-via-vnet
  action: resolve
  filter: dns_resolver
  enabled: true
  precedence: 500
  traffic: dns.fqdn == "db.internal.example.com"
  ruleSettings:
    dnsResolvers:
      ipv4:
        - ip: 10.0.0.53
          vnetId:
            valueFrom:
              kind: CloudflareZeroTrustTunnelVirtualNetwork
              name: corp-vnet
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `policy_id` | string | UUID of the created policy |
| `precedence` | string | Evaluation order (Cloudflare-assigned when the spec left it unset) |

## Related Components

- [Cloudflare Zero Trust List](/docs/catalog/cloudflare/cloudflarezerotrustlist) -- reusable domain/IP/email sets referenced from `traffic`
- [Cloudflare Zero Trust Tunnel Virtual Network](/docs/catalog/cloudflare/cloudflarezerotrusttunnelvirtualnetwork) -- `ruleSettings.dnsResolvers.*.vnetId`
- [Cloudflare Zero Trust Access Identity Provider](/docs/catalog/cloudflare/cloudflarezerotrustaccessidentityprovider) -- identities the `identity` expression matches
- [Cloudflare Ruleset](/docs/catalog/cloudflare/cloudflareruleset) -- website WAF / firewall, not employee Gateway traffic
