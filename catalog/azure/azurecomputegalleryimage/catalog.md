# Azure Compute Gallery Image

Deploys a gallery image: one image definition inside an Azure Compute Gallery plus its published, region-replicated versions -- the release artifacts VMs actually deploy. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Image definition** -- the image's identity (publisher/offer/SKU), OS type, security posture, and recommended sizing
- **Image versions** -- one per entry in the versions list, each built from its source (disk snapshot, VHD blob, or managed image/VM) and replicated to its target regions

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Compute Gallery** -- reference an AzureComputeGallery or provide an existing gallery's name.
- **A version source** -- publishing a version needs an OS disk snapshot (reference an AzureDiskSnapshot), a VHD blob + storage account, or a managed image/VM ID.

### Azure Subscription

- **The definition's identity is permanent** -- publisher/offer/SKU, OS type, generation, and the security posture are all fixed at creation.
- **Pick ONE security posture** -- trusted-launch supported/enabled and confidential-VM supported/enabled are mutually exclusive; trusted launch requires Hyper-V generation V2.
- **Versions bill for replica storage** -- each version stores replicas in every target region; the definition itself is free.

## Deploy

### Console

Open the deployment store, find **Azure Compute Gallery Image**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Linux Gen2 Trusted Launch** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f gallery-image.yaml
```

## After Deploy

Publish new versions by adding entries to the versions list (your image pipeline snapshots a prepared disk, then the manifest adds the version). Point production VMs at exact version IDs from the `version_ids` output; let development float on the definition's `image_id` (latest). Set `excludeFromLatest` on a version to pull it out of "latest" resolution without unpublishing it.
