---
title: "Presets"
description: "Ready-to-deploy configuration presets for SQL Server"
type: "preset-list"
componentSlug: "sql-server"
componentTitle: "SQL Server"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-sql-auth"
    rank: "01"
    title: "Standard SQL-Auth Server"
    excerpt: "This preset creates the simplest useful Azure SQL logical server: SQL authentication, the public endpoint with the Azure-services firewall rule, and Azure's defaults everywhere else (TLS 1.2 floor,..."
  - slug: "02-entra-only"
    rank: "02"
    title: "Entra-Only Server with Microsoft Defender"
    excerpt: "This preset creates a passwordless Azure SQL logical server: Microsoft Entra is the ONLY authentication mechanism (SQL logins are disabled server-wide, so no password exists to leak or rotate), with..."
  - slug: "03-private-hardened"
    rank: "03"
    title: "Private Hardened Server with CMK and Auditing"
    excerpt: "This preset creates a compliance-oriented Azure SQL logical server: no public endpoint (private endpoints only), a customer-managed transparent-data-encryption key unwrapped through a user-assigned..."
---

# SQL Server Presets

Ready-to-deploy configuration presets for SQL Server. Each preset is a complete manifest you can copy, customize, and deploy.
