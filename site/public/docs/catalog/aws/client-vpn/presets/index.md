---
title: "Presets"
description: "Ready-to-deploy configuration presets for Client VPN"
type: "preset-list"
componentSlug: "client-vpn"
componentTitle: "Client VPN"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-certificate-split-tunnel"
    rank: "01"
    title: "Certificate Split-Tunnel VPN"
    excerpt: "This preset creates a Client VPN endpoint with mutual-TLS authentication and split-tunnel routing — the standard corp-access posture: only traffic destined for the VPC goes through the VPN,..."
  - slug: "02-certificate-full-tunnel"
    rank: "02"
    title: "Certificate Full-Tunnel VPN"
    excerpt: "This preset creates a Client VPN endpoint that routes ALL client traffic through AWS — internet included — for postures where every packet must egress through inspected, NAT-ed infrastructure."
  - slug: "03-saml-sso"
    rank: "03"
    title: "SAML SSO VPN with Self-Service Portal"
    excerpt: "This preset creates a Client VPN endpoint with SAML 2.0 single sign-on and group-scoped authorization — users authenticate with their existing identity provider (Okta, Entra ID, ...), download their..."
---

# Client VPN Presets

Ready-to-deploy configuration presets for Client VPN. Each preset is a complete manifest you can copy, customize, and deploy.
