---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cert Manager Certificate"
type: "preset-list"
componentSlug: "cert-manager-certificate"
componentTitle: "Cert Manager Certificate"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-single-domain-dns"
    rank: "01"
    title: "Single Domain DNS-Validated Certificate"
    excerpt: "This preset provisions an ACM certificate for a single domain using automated DNS validation via Route53. DNS validation is the recommended method because it requires no manual intervention -- ACM..."
  - slug: "02-wildcard-domain"
    rank: "02"
    title: "Wildcard Domain DNS-Validated Certificate"
    excerpt: "This preset provisions a wildcard ACM certificate that covers all subdomains of a domain, plus the apex domain itself as a Subject Alternative Name (SAN). A single wildcard certificate can secure..."
  - slug: "03-external-dns"
    rank: "03"
    title: "DNS-Validated Certificate with External DNS"
    excerpt: "This preset requests an ACM certificate for a domain whose DNS lives outside Route53 (Cloudflare, a registrar's DNS, an on-prem zone). The deployment creates the certificate and finishes without..."
  - slug: "04-private-ca-internal-tls"
    rank: "04"
    title: "Private CA Certificate for Internal TLS"
    excerpt: "This preset issues a certificate from your AWS Private Certificate Authority (ACM-PCA) with managed early renewal -- the shape for internal TLS where clients trust your private root instead of a..."
---

# Cert Manager Certificate Presets

Ready-to-deploy configuration presets for Cert Manager Certificate. Each preset is a complete manifest you can copy, customize, and deploy.
