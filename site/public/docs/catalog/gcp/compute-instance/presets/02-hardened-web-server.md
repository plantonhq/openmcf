---
title: "Hardened Web Server"
description: "The production security posture for an internet-facing workload: a Shielded VM on a custom-mode subnetwork with no external IP, a dedicated least-privilege service account, OS Login instead of static..."
type: "preset"
rank: "02"
presetSlug: "02-hardened-web-server"
componentSlug: "compute-instance"
componentTitle: "Compute Instance"
provider: "gcp"
icon: "package"
order: 2
---

# Hardened Web Server

The production security posture for an internet-facing workload: a
Shielded VM on a custom-mode subnetwork with no external IP, a dedicated
least-privilege service account, OS Login instead of static SSH keys,
and firewall targeting via network tags.

## What this preset creates

An `e2-standard-2` Debian VM with secure boot verifying the boot chain,
placed in a referenced `GcpSubnetwork`, authenticating as a referenced
`GcpServiceAccount` with the single `cloud-platform` scope (so IAM
roles — not OAuth scopes — decide what it can do). No `accessConfigs`
means no external IP: traffic arrives through a load balancer, egress
through Cloud NAT, and operators through IAP.
`allowStoppingForUpdate: true` lets machine-type and service-account
changes apply via a brief stop/restart instead of failing.

## Prerequisites

- A `GcpSubnetwork` named `app-subnet` whose region contains the
  instance's zone (replace the reference with your own).
- A `GcpServiceAccount` named `web-server` granted exactly the IAM roles
  the workload needs (replace the reference with your own).
- Firewall rules selecting the `web-server` and `allow-iap-ssh` tags
  (80/443 from the load balancer, 22 from IAP's range).

## Why no external IP

A VM without an external IP is unreachable from the internet regardless
of firewall mistakes. The `external_ip` output is `""` for this VM;
consumers reference `internal_ip` instead.

## Remix ideas

- Reserve a static internal IP (`GcpAddress`, address type INTERNAL) and
  reference it from `networkIp` so the address survives VM rebuilds.
- Add `confidentialInstanceConfig` (with
  `scheduling.onHostMaintenance: TERMINATE` and a supported machine
  family) for hardware memory encryption.
- Add a CMEK boot disk by referencing a `GcpKmsKey` from
  `bootDisk.kmsKey`.
