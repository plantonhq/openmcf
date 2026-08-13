# AzurePrivateDnsResolver Terraform Module

## Overview

Creates an Azure DNS Private Resolver -- the managed DNS proxy that resolves names across the hybrid boundary -- and its inbound/outbound endpoints as composed child resources. Inbound endpoints give on-premises DNS forwarders a private IP inside the network to query; outbound endpoints carry queries out of Azure, steered by attached forwarding rulesets.

## Resources Created

- `azurerm_private_dns_resolver` -- the resolver, anchored to one virtual network
- `azurerm_private_dns_resolver_inbound_endpoint` -- one per `spec.inbound_endpoints` entry, keyed by endpoint name
- `azurerm_private_dns_resolver_outbound_endpoint` -- one per `spec.outbound_endpoints` entry, keyed by endpoint name

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzurePrivateDnsResolverSpec fields; the resource group, virtual network, and endpoint subnet references arrive as resolved literals; the allocation method arrives as the enum NAME and is mapped to the ARM wire value in `locals.tf`

## Outputs

- `dns_resolver_id` -- the resolver's ARM resource ID
- `dns_resolver_name` -- the resolver's name
- `inbound_endpoint_ip` -- the FIRST declared inbound endpoint's private IP (what on-premises forwarders point at); empty when no inbound endpoints are declared
- `inbound_endpoint_ips` -- ALL inbound endpoint IPs, keyed by endpoint name
- `outbound_endpoint_id` -- the FIRST declared outbound endpoint's ARM id (what forwarding rulesets reference by default); empty when no outbound endpoints are declared
- `outbound_endpoint_ids` -- ALL outbound endpoint ids, keyed by endpoint name

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **One resolver per virtual network** -- Azure enforces it at deploy time; a second create against the same network is rejected.
- **Every endpoint needs its own dedicated subnet**, delegated to `Microsoft.Network/dnsResolvers`, sized /28 to /24, hosting nothing else -- ARM validates all of this at deploy time (the delegation lives on the referenced subnet, so no offline gate can check it).
- **Everything except tags is create-only** -- the resolver and its endpoints have no ARM update surface beyond tags; endpoint edits replace that endpoint (keyed by name, so siblings are untouched).
- **The inbound Static/Dynamic contract is spec-validated** (STATIC requires an address; DYNAMIC forbids one) -- mirroring the provider's own pre-flight check.
- **Endpoint deletes outlive the API's first answer** -- the provider polls until the endpoint is verifiably gone (up to a few minutes each).
- **Endpoints bill hourly** from the moment they provision; the resolver object itself is free. Azure allows up to 5 endpoints each way per resolver (a service quota, deliberately not mirrored as validation).

## Required Permissions

The deploying principal needs `Microsoft.Network/dnsResolvers/*` plus join permissions on the delegated subnets (Network Contributor on the resource group covers both).
