# Overview

The **Azure Private Link Service API Resource** provides a consistent and standardized interface for deploying and managing Private Link Services -- the PROVIDER side of Azure Private Link. A service you run (behind a Standard internal load balancer, or at one fixed destination IP) becomes privately consumable from other virtual networks, subscriptions, and even other Entra tenants, through private endpoints -- traffic never leaves the Microsoft backbone, and nothing is exposed publicly.

## Purpose

We developed this API resource so service teams can publish internal services across network boundaries without VNet peering, address-space negotiation, or public exposure:

- **Private publication**: consumers connect through private endpoints using the service's globally unique alias -- no RBAC on your subscription needed
- **Controlled discoverability**: visibility lists gate who can find the service; auto-approval lists decide whose connections skip the manual queue
- **Overlap-proof networking**: consumer address spaces never meet yours -- traffic is source-NATed through the service's NAT configurations

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Destination Contract**: exactly one of a Standard load balancer frontend set or a fixed destination IP, enforced at validation time
- **Single-Primary NAT Contract**: ARM's one-primary-NAT-configuration rule is spec-enforced upfront, with up to 8 configurations for NAT-port headroom
- **PROXY Protocol Support**: optionally carry the consumer's original source IP to backends that parse PROXY v2 headers

## Use Cases

- **Internal SaaS**: publish a shared platform service (APIs, databases behind a load balancer) to product teams in other VNets
- **Cross-tenant delivery**: let customers in their own Entra tenants consume your service through its alias
- **Partner connectivity**: expose one service to a partner's subscription without peering networks or opening firewalls
- **Managed-service backends**: the pattern Azure's own first-party services use to reach into customer VNets, applied to yours

## Future Enhancements

Future updates will include:

- **Connection inventory**: console surfacing of pending/approved private-endpoint connections with one-click approval
- **Consumer-side pairing**: guided creation of the consuming AzurePrivateEndpoint from a service's alias

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
