# Overview

The **Azure ExpressRoute Port API Resource** provides a consistent and standardized interface for deploying and managing ExpressRoute Ports -- your own pair of physical ports on a Microsoft Enterprise Edge router at a peering location (ExpressRoute Direct). Where an ordinary circuit rents capacity through a connectivity provider, a port IS the capacity: order 10 or 100 Gbps of dual physical links, arrange the cross-connects with the facility, and carve circuits from the port's bandwidth.

## Purpose

We developed this API resource so the largest tier of hybrid connectivity is a first-class, versioned object:

- **The physical facts, typed**: peering location, bandwidth, encapsulation, and billing model as explicit, validated choices
- **Link-pair management**: admin state and MACsec (cipher, Key Vault-held CAK/CKN, SCI) on the fixed physical pair
- **Letter-of-authorization outputs**: each link's router, interface, patch panel, and rack surface as outputs -- the facts the facility's cross-connect order needs
- **Authorization issuance**: declare named authorizations whose ARM-generated keys let circuits in OTHER subscriptions ride this port

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **MACsec Contracts Enforced Upfront**: the CAK/CKN pairing and the user-assigned-identity requirement are validated in seconds, not at deploy time
- **Sensitive-by-Default Keys**: every issued authorization key is secret-typed in both deployment engines
- **Cost Honesty**: the spec and docs are explicit that the port bills its full monthly rate from creation

## Use Cases

- **Massive-scale interconnect**: 10/100 Gbps dedicated capacity for the largest hybrid estates
- **MACsec-encrypted links**: layer-2 encryption between your edge and Microsoft's, keyed from your own Key Vault
- **Multi-subscription capacity sharing**: issue authorizations so other subscriptions' circuits ride one port

## Future Enhancements

Future updates will include:

- **Circuit composition views**: the port with its carved circuits as one capacity story
- **Facility handoff tracking**: console surfacing of the LOA facts and cross-connect workflow

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
