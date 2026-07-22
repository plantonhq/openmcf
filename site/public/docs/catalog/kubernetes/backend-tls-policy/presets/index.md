---
title: "Presets"
description: "Ready-to-deploy configuration presets for Backend TLS Policy"
type: "preset-list"
componentSlug: "backend-tls-policy"
componentTitle: "Backend TLS Policy"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-internal-ca-configmap"
    rank: "01"
    title: "Internal CA via ConfigMap"
    excerpt: "The most common BackendTLSPolicy: the gateway originates TLS to an internal backend whose serving certificate is signed by YOUR OWN CA — the private-PKI posture behind most cert-manager..."
  - slug: "02-public-ca-system"
    rank: "02"
    title: "Public CA via System Trust Store"
    excerpt: "This preset secures the gateway-to-backend hop for backends serving PUBLICLY-issued certificates (Let's Encrypt, a commercial CA): instead of bringing a CA bundle, the policy trusts the gateway..."
  - slug: "03-spiffe-mtls-backend"
    rank: "03"
    title: "SPIFFE mTLS Backend"
    excerpt: "This preset handles the mesh-identity case: the backend's certificate does not carry the DNS hostname the gateway dials — it carries a SPIFFE ID (the URI-SAN identity pattern of SPIRE, Istio, and..."
---

# Backend TLS Policy Presets

Ready-to-deploy configuration presets for Backend TLS Policy. Each preset is a complete manifest you can copy, customize, and deploy.
