---
title: "Smoother Connection Setup Experience"
date: 2026-03-04
category: improvement
tags:
  - console
  - cli
excerpt: "Connection wizards now surface errors clearly, the CLI supports all connection types, and connection URLs are cleaner and more consistent."
author:
  - name: Swarup Donepudi
    title: Founder
---

Connections are how you link Planton to your cloud providers, source control, container registries, and IaC backends. Over the past week, we've refined the connection setup experience across both the console and the CLI to eliminate silent failures, clean up the URL structure, and make the overall experience more polished.

## Better Error Handling in Connection Wizards

Previously, if something went wrong when you clicked "Create" in a connection wizard, the button would reset with no feedback — you had no idea what happened. Now, the wizard surfaces errors inline with clear context:

- **Validation errors** appear in an amber panel listing exactly which fields need attention, with guidance to go back and fix them
- **Server errors** are handled by the platform's standard error dialog so you're never left guessing
- **The error clears** automatically when you retry or navigate back to fix the issue

This applies to all connection types across the platform.

We also fixed an issue where creating platform-managed Pulumi or Terraform backends would silently fail because a required field wasn't being set behind the scenes. The wizard now handles default values correctly for all backend methods, whether they have configuration steps or not.

## CLI Support for All Connection Types

The `planton apply` command now works with YAML manifests for every connection type on the platform — cloud providers, source control, container registries, IaC backends, and everything in between. Previously, running `planton apply -f github-connection.yaml` would fail with a "backend config not found" error because the CLI didn't know how to route these resources. That's fixed — any connection you can create in the console can now be managed through the CLI as well.

If you manage connections as YAML files in your ops repository, you can now apply them directly without switching to the console.

## Cleaner Connection URLs

Connection pages in the console now follow a consistent, org-scoped URL structure:

| Before | After |
|--------|-------|
| `/resource/connect/aws-provider-connection/create` | `/orgs/acme/connections/aws/create` |
| `/resource/connect/github-connection/create` | `/orgs/acme/connections/github/create` |

The new URLs are shorter, show which organization you're working in, and use friendly slugs (`aws`, `gcp`, `github`) instead of internal resource type names. Old URLs redirect gracefully, so existing bookmarks continue to work.

OAuth callback URLs for GCP, Azure, and GitHub were also consolidated under a clean `/oauth/callbacks/` path.

## Visual Polish

Several visual inconsistencies across the 19 connection screens have been addressed:

- **Provider icons** are now consistent between the connection list view and detail pages — GitHub and Kubernetes icons that were previously missing from the list view now display correctly
- **URL slugs** are simplified across all connection detail pages, removing redundant `-provider` and `-connection` suffixes
- **GCP service account key uploads** now show a confirmation step with an editable secret name before creating the secret, giving you control over naming
- **Path selection** in the GCP and Kubernetes wizards uses radio buttons instead of ambiguous icon cards, making the active selection unmistakable
- **Cloudflare R2 Access Key ID** is now stored as a secret reference rather than a plain value, consistent with the platform's security model for credential material
