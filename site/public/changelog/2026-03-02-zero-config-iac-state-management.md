---
title: "Zero-Config IaC State Management for New Organizations"
date: 2026-03-02
category: improvement
tags:
  - infra-hub
  - console
excerpt: "New organizations now get Pulumi and Terraform state backends automatically, and self-hosted runners can use cloud-native identity instead of pasting credentials."
author:
  - name: Swarup Donepudi
    title: Founder
---

When you create a new organization, the platform now automatically provisions Pulumi and Terraform state backends alongside the existing secret backend. Previously, you had to manually create these connections before you could deploy any infrastructure — a friction point that tripped up every new user during onboarding. Now, your first `planton apply` just works.

## What Changed

Every new organization receives three platform-managed backends out of the box:

- **Secret backend** — for encrypting and storing sensitive configuration values (this was already auto-provisioned; now it also gets proper access controls and event tracking)
- **Pulumi backend** — for storing Pulumi state files
- **Terraform backend** — for storing Terraform state files

All three are ready to use immediately. You don't need to configure storage buckets, access credentials, or backend URLs. The platform handles all of it.

## Runner-Based Authentication for State Backends

If you're running a self-hosted runner in your own cloud environment, you no longer need to provide state backend credentials through the platform. A new **auth mode** setting on Pulumi and Terraform state backend connections lets you choose between two approaches:

**Provide Credentials** (the default) — You supply credentials through the platform's secret management system. The control plane resolves them and passes them to the runner at deployment time. This is the right choice when your runner doesn't have direct access to the state storage.

**Runner Environment** — The runner uses its own environment credentials (AWS IAM roles via IRSA, GCP Workload Identity, Azure Managed Identity, or environment variables) to authenticate with the state backend. No credentials are stored in Planton at all. This is ideal when your runner already has access to the storage account through your cloud provider's native identity system.

When you choose runner-based authentication, the connection wizard streamlines by hiding all credential fields — you only configure the storage location (bucket name, region, endpoint), and the runner handles the rest.

## Guard Rails

The platform prevents misconfiguration: if you try to use a platform-managed state backend with a runner that uses its own credentials, the system rejects the request with a clear error. Platform-managed backends are only accessible through the platform's credential resolution pipeline, not through runner-local identity.

## Why This Matters

- **Faster onboarding** — new organizations can deploy infrastructure immediately without manual backend setup
- **Better security posture** — self-hosted runners can use cloud-native identity, keeping credentials entirely out of the platform
- **Less configuration** — the wizard adapts based on your auth mode choice, showing only the fields you need
