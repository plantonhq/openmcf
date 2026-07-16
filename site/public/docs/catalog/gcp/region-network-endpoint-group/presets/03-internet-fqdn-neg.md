---
title: "External Internet Origin"
description: "An internet NEG that lets a Google Cloud load balancer front an external origin — an on-prem service or a third-party API reached by FQDN — so it can sit behind Cloud CDN, Cloud Armor, and a single..."
type: "preset"
rank: "03"
presetSlug: "03-internet-fqdn-neg"
componentSlug: "region-network-endpoint-group"
componentTitle: "Region Network Endpoint Group"
provider: "gcp"
icon: "package"
order: 3
---

# External Internet Origin

An internet NEG that lets a Google Cloud load balancer front an external origin — an on-prem service or a third-party API reached by FQDN — so it can sit behind Cloud CDN, Cloud Armor, and a single anycast IP alongside your GCP backends.

## When to Use

- Hybrid setups where some backends live outside GCP (on-prem, another cloud)
- Putting Cloud CDN or Cloud Armor in front of a third-party origin

## Remix Notes

- Use `INTERNET_FQDN_PORT` for a hostname origin and `INTERNET_IP_PORT` for a fixed IP origin.
- The actual endpoint (the FQDN:port or IP:port) is attached to the NEG as a network endpoint after the NEG is created.
- Reference a `GcpVpcNetwork` under `network.valueFrom` to wire the NEG to a network Planton manages.
