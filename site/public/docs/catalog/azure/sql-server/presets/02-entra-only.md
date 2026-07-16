---
title: "Entra-Only Server with Microsoft Defender"
description: "This preset creates a passwordless Azure SQL logical server: Microsoft Entra is the ONLY authentication mechanism (SQL logins are disabled server-wide, so no password exists to leak or rotate), with..."
type: "preset"
rank: "02"
presetSlug: "02-entra-only"
componentSlug: "sql-server"
componentTitle: "SQL Server"
provider: "azure"
icon: "package"
order: 2
---

# Entra-Only Server with Microsoft Defender

This preset creates a passwordless Azure SQL logical server: Microsoft
Entra is the ONLY authentication mechanism (SQL logins are disabled
server-wide, so no password exists to leak or rotate), with Microsoft
Defender's threat detection and express vulnerability assessment on.

## When to Use

- Organizations standardizing database access on Entra identities and
  groups instead of shared SQL passwords
- Security postures that treat static credentials as findings

## Key Configuration Choices

- **`azureadAuthenticationOnly: true`** -- SQL logins are rejected
  server-wide; `administratorLogin`/`administratorPassword` must be
  omitted entirely. Note ARM will not re-enable a password later without
  configuring one explicitly
- **An Entra GROUP as the administrator** -- admits the whole DBA team
  through group membership instead of per-person grants
- **Defender enabled** -- SQL-injection/anomaly/exfiltration detectors
  alert the security team and subscription admins;
  `expressVulnerabilityAssessmentEnabled` adds agentless scanning with
  no storage account to manage

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the server in | The resource group's `status.outputs.resource_group_name` |
| `myorg-entra-sql` | 1-63 lowercase chars, globally unique | Your naming convention |
| `<entra-group-name>` | The Entra group's display name (the login name) | Microsoft Entra ID -> Groups |
| `<entra-group-object-id>` | The Entra group's directory object ID | Microsoft Entra ID -> Groups -> Overview |
| `<security-team-email>` | Where Defender alerts go | Your security team |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Connections authenticate with Entra tokens (no password):

```text
Server={status.outputs.fqdn},1433;Database={database};Authentication=Active Directory Default;Encrypt=True;
```
