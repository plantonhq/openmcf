---
title: "Application Gateway"
description: "Application Gateway deployment documentation"
icon: "package"
order: 100
componentName: "azureapplicationgateway"
---

# Azure Application Gateway

Creates an Azure Application Gateway -- the Layer 7 load balancer and reverse proxy that routes by host name and URI path, terminates TLS (including mutual TLS) with Key Vault certificates that renew in place, rewrites requests and responses, proxies raw TCP/TLS, and enforces a Web Application Firewall policy on the WAF_v2 SKU.

## What Gets Created

When you deploy an AzureApplicationGateway resource, Planton provisions:

- **Application Gateway** -- an `azurerm_application_gateway` carrying every sub-object as one atomic resource: frontend IP configurations and ports, HTTP(S) and TCP/TLS listeners, backend pools and settings, routing rules and URL path maps, health probes, certificates and SSL profiles, redirects, rewrite rule sets, Private Link configurations, and custom error pages

Sub-objects wire to each other by name inside the spec. What other resources need to reach is exported as name-keyed map outputs: NICs and scale sets join pools through `backend_address_pool_ids`, and frontends chain through `frontend_ip_configuration_ids`. The WAF policy composes as its own `AzureWebApplicationFirewallPolicy` resource, referenced at three levels (gateway, listener, path rule).

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A DEDICATED subnet** for the gateway (no other resource types; /24 recommended -- an `AzureSubnet` in composed environments)
- **A Standard static public IP** (`AzurePublicIp`) for public frontends
- **For Key Vault certificates**: a user-assigned identity with GET on the vault's secrets, granted before the gateway deploys
- **Time**: applies run 15-25 minutes -- Azure's slowest networking resource

## Quick Start

Create a file `gateway.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureApplicationGateway
metadata:
  name: web-gateway
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureApplicationGateway.web-gateway
spec:
  region: eastus
  resourceGroup:
    value: app-rg
  name: web-gateway
  subnetId:
    valueFrom:
      name: gateway-subnet
  sku: STANDARD_V2
  capacity: 2
  frontendIpConfigurations:
    - name: public
      publicIpAddressId:
        valueFrom:
          name: gateway-ip
  frontendPorts:
    - name: http
      port: 80
  backendAddressPools:
    - name: web
      ipAddresses:
        - 10.0.1.4
  backendHttpSettings:
    - name: web-settings
      port: 8080
      protocol: HTTP
  httpListeners:
    - name: http-listener
      frontendIpConfigurationName: public
      frontendPortName: http
      protocol: HTTP
  requestRoutingRules:
    - name: http-rule
      ruleType: BASIC_ROUTING
      httpListenerName: http-listener
      priority: 100
      backendAddressPoolName: web
      backendHttpSettingsName: web-settings
```

Deploy:

```shell
planton apply -f gateway.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region (must match the subnet, IPs, and WAF policy). | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name. | Required |
| `name` | `string` | The gateway's name, unique within the resource group. | 1-80 chars |
| `subnetId` | `StringValueOrRef` | The DEDICATED gateway subnet (references `AzureSubnet`). | Required |
| `sku` | `enum` | `BASIC`, `STANDARD_V2`, or `WAF_V2`. | Required |
| `frontendIpConfigurations` | `list` | Public (`publicIpAddressId`) XOR private (`subnetId` + allocation) frontends. | min 1 |
| `frontendPorts` | `list` | Named ports listeners bind to. | min 1 |
| `backendAddressPools` | `list` | Pools (FQDNs/IPs, or joined member-side). | min 1 |

Exactly one of `capacity` (1-125; 1-2 on BASIC) and `autoscale` must be set, and at least one listener + routing rule (L7 `httpListeners`/`requestRoutingRules` or L4 `listeners`/`routingRules`) plus matching backend settings.

### The L7 Surface

| Field | Description |
|-------|-------------|
| `backendHttpSettings[]` | Port, protocol (HTTP/HTTPS), cookie affinity, path prefix, timeouts, probe, host-header handling, connection draining, backend-CA trust, dedicated connections. |
| `httpListeners[]` | Frontend + port + protocol, wildcard `hostNames`, TLS certificate, SNI, SSL profile, per-listener WAF policy, custom error pages. |
| `requestRoutingRules[]` | `BASIC_ROUTING` (listener -> pool/redirect) or `PATH_BASED_ROUTING` (listener -> url path map); `priority` 1-20000 required. |
| `urlPathMaps[]` | Per-path rules with their own pools, redirects, rewrites, and WAF policies, plus defaults for unmatched paths. |
| `probes[]` | HTTP(S) GET probes (path, host, match criteria) or TCP/TLS connection probes. |
| `redirectConfigurations[]` | 301/302/303/307 redirects to a listener or external URL. |
| `rewriteRuleSets[]` | Condition-driven header edits and URL rewrites (URL rewriting not on BASIC). |

### TLS

| Field | Description |
|-------|-------------|
| `sslCertificates[]` | Key Vault secret reference (renewals propagate) XOR inline PFX `data`+`password`. |
| `trustedRootCertificates[]` | Private-CA bundles for backend HTTPS trust. |
| `trustedClientCertificates[]` + `sslProfiles[]` | Mutual TLS: client-CA bundles, issuer-DN and OCSP checks, per-profile TLS policy. |
| `sslPolicy` | Gateway-wide policy: `PREDEFINED` (e.g. `AppGwSslPolicy20220101S`), `CUSTOM`/`CUSTOM_V2` (min protocol + ciphers), or `disabledProtocols`. |

### Layer 4, WAF, and the Rest

| Field | Description |
|-------|-------------|
| `listeners[]` / `backends[]` / `routingRules[]` | The TCP/TLS proxy trio. |
| `firewallPolicyId` | The `AzureWebApplicationFirewallPolicy` reference (WAF_V2 only) + `forceFirewallPolicyAssociation`. |
| `privateLinkConfigurations[]` | Private Link NAT allocations exposing frontends to other networks. |
| `customErrorConfigurations[]` | Gateway-wide custom error pages. |
| `http2Enabled` / `fipsEnabled` / `globalConfiguration` | HTTP/2, FIPS mode, request/response buffering. |
| `zones` / `identity` / `tags` | Zone redundancy, managed identity, user tags. |

## Examples

### WAF Gateway with Path-Based Routing

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureApplicationGateway
metadata:
  name: waf-gateway
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: app-rg
  name: waf-gateway
  subnetId:
    valueFrom:
      name: gateway-subnet
  sku: WAF_V2
  autoscale:
    minCapacity: 2
    maxCapacity: 10
  firewallPolicyId:
    valueFrom:
      kind: AzureWebApplicationFirewallPolicy
      name: waf-baseline
      fieldPath: status.outputs.policy_id
  identity:
    type: USER_ASSIGNED
    identityIds:
      - valueFrom:
          name: gateway-identity
  frontendIpConfigurations:
    - name: public
      publicIpAddressId:
        valueFrom:
          name: gateway-ip
  frontendPorts:
    - name: https
      port: 443
  backendAddressPools:
    - name: web
    - name: api
  backendHttpSettings:
    - name: web-settings
      port: 8080
      protocol: HTTP
    - name: api-settings
      port: 8443
      protocol: HTTPS
      pickHostNameFromBackendAddress: true
  httpListeners:
    - name: https-listener
      frontendIpConfigurationName: public
      frontendPortName: https
      protocol: HTTPS
      sslCertificateName: tls
      requireSni: true
  requestRoutingRules:
    - name: path-routing
      ruleType: PATH_BASED_ROUTING
      httpListenerName: https-listener
      priority: 100
      urlPathMapName: paths
  urlPathMaps:
    - name: paths
      defaultBackendAddressPoolName: web
      defaultBackendHttpSettingsName: web-settings
      pathRules:
        - name: api
          paths: ["/api/*"]
          backendAddressPoolName: api
          backendHttpSettingsName: api-settings
  sslCertificates:
    - name: tls
      keyVaultSecretId:
        valueFrom:
          kind: AzureKeyVaultCertificate
          name: web-tls
          fieldPath: status.outputs.versionless_secret_id
```

### Join a Pool from a Network Interface

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkInterface
metadata:
  name: app-nic
spec:
  ipConfigurations:
    - name: primary
      subnetId:
        valueFrom:
          name: app-subnet
      applicationGatewayBackendAddressPoolIds:
        - valueFrom:
            kind: AzureApplicationGateway
            name: waf-gateway
            fieldPath: status.outputs.backend_address_pool_ids.web
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `application_gateway_id` | `string` | The gateway's ARM ID -- referenced by AKS AGIC and management tooling |
| `application_gateway_name` | `string` | The gateway's name |
| `backend_address_pool_ids` | `map(string)` | Each pool's ARM ID, keyed by pool name -- the membership seam for NICs and scale sets |
| `frontend_ip_configuration_ids` | `map(string)` | Each frontend's ARM ID, keyed by frontend name |
| `private_ip_address` | `string` | The first private frontend's address (internal DNS target); empty when all frontends are public |
| `private_ip_addresses` | `list(string)` | All private frontends' addresses |

## Related Components

- [AzureWebApplicationFirewallPolicy](/docs/catalog/azure/web-application-firewall-policy) — the WAF rule set the gateway enforces
- [AzureSubnet](/docs/catalog/azure/subnet) — the dedicated gateway subnet
- [AzurePublicIp](/docs/catalog/azure/public-ip) — public frontend addresses (and where their IPs live)
- [AzureKeyVaultCertificate](/docs/catalog/azure/key-vault-certificate) — TLS certificates that renew in place
- [AzureUserAssignedIdentity](/docs/catalog/azure/user-assigned-identity) — the vault-access identity
- [AzureNetworkInterface](/docs/catalog/azure/network-interface) — backend pool members
- [AzureVirtualMachineScaleSet](/docs/catalog/azure/virtual-machine-scale-set) — fleet pool members
- [AzureResourceGroup](/docs/catalog/azure/resource-group) — provides the resource group for gateway placement
