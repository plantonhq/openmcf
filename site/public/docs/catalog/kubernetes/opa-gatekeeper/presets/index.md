---
title: "Presets"
description: "Ready-to-deploy configuration presets for OPA Gatekeeper"
type: "preset-list"
componentSlug: "opa-gatekeeper"
componentTitle: "OPA Gatekeeper"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-audit-first"
    rank: "01"
    title: "Audit-first preset"
    excerpt: "Gatekeeper as it ships: three webhook replicas, the policy webhook fail-OPEN (`failurePolicy: Ignore`), and the audit controller re-checking existing resources every 60 seconds. This is the right..."
  - slug: "02-production-enforce"
    rank: "02"
    title: "Production enforce preset"
    excerpt: "Gatekeeper with teeth. The policy webhook goes fail-CLOSED (`failurePolicy: Fail`): a request the engine cannot evaluate is REJECTED, closing the window an attacker could time against an engine..."
  - slug: "03-cert-manager-tls"
    rank: "03"
    title: "cert-manager TLS preset"
    excerpt: "Gatekeeper serving its webhook with a certificate issued by cert-manager instead of the engine's embedded rotator. Organizations standardizing certificate issuance (one CA, one audit trail, one..."
---

# OPA Gatekeeper Presets

Ready-to-deploy configuration presets for OPA Gatekeeper. Each preset is a complete manifest you can copy, customize, and deploy.
