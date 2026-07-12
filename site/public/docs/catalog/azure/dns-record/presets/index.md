---
title: "Presets"
description: "Ready-to-deploy configuration presets for DNS Record"
type: "preset-list"
componentSlug: "dns-record"
componentTitle: "DNS Record"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-web-app-a-record"
    rank: "01"
    title: "Web App A Record"
    excerpt: "The 30-second DNS record: point a hostname (www) at an IPv4 address. Add more addresses to the list for round-robin distribution."
  - slug: "02-apex-alias-record"
    rank: "02"
    title: "Apex Alias Record"
    excerpt: "Point the zone apex (your-domain.com itself) at an Azure Public IP with an alias record. Azure keeps the answer in sync with the resource -- when the IP's address changes, the record follows..."
  - slug: "03-mail-mx-records"
    rank: "03"
    title: "Mail MX Records"
    excerpt: "The domain's mail-exchange record set: a primary and a backup mail server, each with its own delivery preference in one record set (Azure stores all MX values for the apex together)."
  - slug: "04-domain-verification-txt"
    rank: "04"
    title: "Domain Verification TXT"
    excerpt: "The ownership-proof TXT record that custom-domain flows require before they bind a hostname: Container Apps checks `asuid.{host}`, Front Door checks `_dnsauth.{host}`, and most SaaS domain..."
---

# DNS Record Presets

Ready-to-deploy configuration presets for DNS Record. Each preset is a complete manifest you can copy, customize, and deploy.
