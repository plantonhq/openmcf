# AzurePrivateDnsRecord Terraform Module

## Overview

Creates one DNS record set (A, AAAA, CNAME, MX, PTR, SRV, or TXT) in an Azure PRIVATE DNS zone -- name resolution visible only inside the virtual networks linked to the zone. The record type is whichever typed payload the spec carries (validation guarantees exactly one), so exactly one count-gated resource materializes.

## Resources Created

Exactly one of:

- `azurerm_private_dns_a_record` -- when `spec.a` carries addresses
- `azurerm_private_dns_aaaa_record` -- when `spec.aaaa` carries addresses
- `azurerm_private_dns_cname_record` -- when `spec.cname` is set
- `azurerm_private_dns_mx_record` -- when `spec.mx` carries entries
- `azurerm_private_dns_ptr_record` -- when `spec.ptr` carries hostnames
- `azurerm_private_dns_srv_record` -- when `spec.srv` carries entries
- `azurerm_private_dns_txt_record` -- when `spec.txt` carries values

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzurePrivateDnsRecordSpec fields; the zone reference, CNAME target, and TXT values arrive as resolved literals

## Outputs

- `record_id` -- the record set's ARM resource ID (coalesced per attribute across the seven variants)
- `fqdn` -- the record set's fully qualified name with a trailing dot, resolvable only from networks linked to the zone

## Behavior Notes

- **Private record sets address the zone by ARM ID** (`private_dns_zone_id`) -- unlike the public DNS resources, which address by (resource group, zone name).
- **No alias records and no CAA/NS types** -- alias is a public-DNS concept; private zones cannot delegate subdomains (create a separate zone per suffix) and never face public certificate authorities.
- **One record set per (name, type)**: a second AzurePrivateDnsRecord with the same name and type in the same zone conflicts rather than merges.
- **The TTL is always sent explicitly** -- 300 seconds (5 minutes) when the spec leaves it unset -- so plans stay deterministic across engines.
- **Azure caps an A record set at 20 addresses** (spec-validated) and each TXT value at 1,024 characters (documented; ARM enforces).
- **Tags land as ARM record-set METADATA** -- record sets carry no ARM tags proper; the provider maps its `tags` argument onto that map.
- **Name, zone, and type are fixed at creation** -- changing any replaces the record set; values and TTL update in place. Record sets are free at rest (private DNS bills per zone and per million queries).

## Required Permissions

The deploying principal needs `Microsoft.Network/privateDnsZones/*/write` on the zone's record types (Private DNS Zone Contributor on the resource group covers it).
