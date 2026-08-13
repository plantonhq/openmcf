# Overview

The **Azure VPN Gateway API Resource** provides a consistent and standardized interface for deploying and managing Virtual WAN VPN Gateways -- the managed site-to-site VPN terminators that live inside virtual hubs (one per hub). Branches described by VPN Sites connect to it through VPN Gateway Connections; the hub's routing distributes the branch routes.

## Purpose

We developed this API resource so the hub's branch on-ramp -- and the NAT machinery overlapping branches need -- is one first-class, versioned object:

- **Capacity as one number**: scale units (500 Mbps each) across a managed active-active pair -- no SKU matrix
- **NAT rules composed in**: each rule an ARM child of the gateway, its ID published in the name-keyed `nat_rule_ids` output for connections to opt into
- **Chart-ready wiring**: the hub is a typed reference; the gateway's ARM ID and instance public IPs surface where connections and branch devices consume them

## Key Features

- **Consistent Interface**: aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Instance-Aware BGP**: custom APIPA peering addresses per gateway instance, the AWS-interop shape
- **Cost Honesty**: the gateway bills from creation and provisions in tens of minutes -- stated where decisions are made

## Use Cases

- **Terminate branch site-to-site VPN** at a Virtual WAN hub
- **Overlapping branch address spaces**: static/dynamic NAT rules that tunnels opt into
- **Multi-cloud interop**: APIPA BGP peering for AWS site-to-site VPN into Azure

## Future Enhancements

Future updates will include:

- **Gateway diagnostics surfaces**: packet capture and health integration

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
