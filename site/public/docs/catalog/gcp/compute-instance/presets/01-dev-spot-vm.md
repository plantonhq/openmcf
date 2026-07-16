---
title: "Dev Spot VM"
description: "The cheapest way to get a full Linux machine on GCP: an `e2-medium` Spot VM booting the latest Debian 12, on the default network with an ephemeral external IP, deleted outright when GCP reclaims the..."
type: "preset"
rank: "01"
presetSlug: "01-dev-spot-vm"
componentSlug: "compute-instance"
componentTitle: "Compute Instance"
provider: "gcp"
icon: "package"
order: 1
---

# Dev Spot VM

The cheapest way to get a full Linux machine on GCP: an `e2-medium` Spot
VM booting the latest Debian 12, on the default network with an ephemeral
external IP, deleted outright when GCP reclaims the capacity.

## What this preset creates

A Spot-provisioned Debian VM. `provisioningModel: SPOT` is the single
switch — the engines derive the API's legacy preemptible flag and force
automatic restart off. `instanceTerminationAction: DELETE` makes the VM
fully disposable: preemption removes the instance and its auto-delete
boot disk.

## Prerequisites

None — the default network exists in every fresh project, and the
Compute Engine API is enabled automatically.

## The Spot trade

Spot capacity is discounted 60–91% against on-demand, and GCP may
reclaim it at any time with 30 seconds' notice. Perfect for dev boxes,
CI runners, and batch workers; wrong for anything that must stay up.

## Remix ideas

- Use `instanceTerminationAction: STOP` to keep the stopped VM and its
  disks across preemptions — restart it manually when capacity returns.
- Add `metadata: {enable-oslogin: "TRUE"}` to log in with IAM identities
  instead of managing SSH keys.
- Drop the `accessConfigs` entry to make the VM private, then reach it
  via IAP tunneling.
- For an on-demand VM, remove the `scheduling` block entirely.
