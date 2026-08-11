# Overview

The **Azure ExpressRoute Circuit API Resource** provides a consistent and standardized interface for deploying and managing ExpressRoute circuits -- the dedicated PRIVATE connection between your infrastructure (on-premises or colocation) and Microsoft, bought either through a connectivity provider or carved from your own ExpressRoute Direct port. The circuit is the billing and identity object; routing starts once peerings are configured on it.

## Purpose

We developed this API resource so the private-connectivity foundation of a hybrid estate is a first-class, versioned object:

- **Two provisioning modes**: name a connectivity provider, its peering location, and Mbps bandwidth -- or reference your ExpressRoute Direct port with Gbps bandwidth
- **Authorization issuance**: declare named authorizations whose ARM-generated keys let gateways in OTHER subscriptions connect to this circuit
- **Honest lifecycle**: the service key, provider provisioning state, and per-authorization keys surface as outputs the rest of the graph consumes

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Provisioning-Mode Contract**: the service-provider trio and the ExpressRoute Direct pair are mutually exclusive and internally co-required -- enforced at validation time, not forty minutes into a deployment
- **SKU Clarity**: tier (LOCAL / STANDARD / PREMIUM reach) and billing family (metered / unlimited) as explicit, documented choices
- **Sensitive-by-Default Keys**: the service key and every issued authorization key are secret-typed in both deployment engines

## Use Cases

- **Datacenter interconnect**: a provider-provisioned circuit from your colocation cage to Azure, private peering into your hub VNet
- **Metro-local connectivity**: a LOCAL-tier circuit for egress-fee-free connectivity within the circuit's metro
- **ExpressRoute Direct**: carve 10/100 Gbps circuits from your own port pair for the largest estates
- **Multi-subscription topologies**: issue authorizations so spoke subscriptions' gateways connect to one shared circuit

## Future Enhancements

Future updates will include:

- **Provider handoff tracking**: console surfacing of the provisioning state with provider-specific next steps
- **Peering composition views**: the circuit with its private/Microsoft peerings as one connectivity story

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
