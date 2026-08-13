# AzureComputeGalleryImage Terraform Module

## Overview

Creates a gallery image: one image definition inside an Azure Compute Gallery (marketplace-style identity, OS type, security posture, recommended sizing) plus its published versions, each replicated to its own target regions.

## Resources Created

- `azurerm_shared_image` -- the image definition
- `azurerm_shared_image_version` -- one per entry in `spec.versions`, keyed by version name (the composed release artifacts VMs actually deploy)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureComputeGalleryImageSpec fields; the resource-group, gallery, snapshot, storage-account, and disk-encryption-set references arrive as resolved literals

## Outputs

- `image_id` -- the definition's ARM resource ID (deploying from it gets the latest version)
- `image_name` -- the definition's name within its gallery
- `version_ids` -- map of version name to ARM ID (pin VMs to an exact release)

## Behavior Notes

- **The security flags fire on argument PRESENCE**: the four trusted-launch/confidential-VM flags are a provider ConflictsWith clique, so each is sent ONLY when true -- an explicit false alongside another flag is provider-rejected. The spec's CEL enforces at most one.
- **Each version has exactly one source** (blob + storage account, OS disk snapshot, or managed image/VM) -- the spec's CEL mirrors the provider's ExactlyOneOf.
- **Almost the whole definition is create-only**; updatable in place: description, disk-type exclusions, recommended sizing, release notes, end-of-life, tags. On versions: target regions, exclude-from-latest, end-of-life, tags.
- **Clearing end_of_life_date forces replacement** on both the definition and a version (the provider's CustomizeDiff).
- **A target region's storage_account_type cannot be updated**: the API rejects it and the provider cannot force replacement for it (region-list membership changes in place) -- remove and re-add the region instead.
- **Shallow replication cannot use disk encryption sets** (provider expand-time check, front-loaded as a spec CEL).
- **The definition is free**; each version bills for the storage its regional replicas consume.
