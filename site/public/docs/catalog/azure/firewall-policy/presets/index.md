---
title: "Presets"
description: "Ready-to-deploy configuration presets for Firewall Policy"
type: "preset-list"
componentSlug: "firewall-policy"
componentTitle: "Firewall Policy"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-egress-baseline"
    rank: "01"
    title: "Standard Egress Baseline"
    excerpt: "This preset creates the STANDARD-tier policy most hub-spoke deployments start from: Microsoft threat intelligence in DENY (known-malicious destinations are blocked, not just logged) and the DNS proxy..."
  - slug: "02-premium-tls-inspection"
    rank: "02"
    title: "Premium TLS Inspection + IDPS"
    excerpt: "This preset creates the PREMIUM-tier policy for regulated or high-security environments: outbound TLS is decrypted, inspected, and re-encrypted using your intermediate CA from Key Vault, and the IDPS..."
---

# Firewall Policy Presets

Ready-to-deploy configuration presets for Firewall Policy. Each preset is a complete manifest you can copy, customize, and deploy.
