# Azure Compute Gallery

Deploys an Azure Compute Gallery -- the shared library an organization keeps its approved VM images in. Image definitions live inside it; VMs and scale sets deploy from their published, region-replicated versions. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Compute gallery** -- the gallery itself: description, optional sharing profile (Groups or Community), tags

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Resource Group** -- reference an AzureResourceGroup or provide an existing group's name.

### Azure Subscription

- **Gallery names forbid dashes** -- up to 80 letters, numbers, dots, and underscores (dots separate logical segments by convention, e.g. `platform.images`).
- **The sharing choice is permanent** -- the entire sharing configuration is fixed at creation; changing it later replaces the gallery. Start Private unless you know you need Community.
- **Community sharing is public** -- it publishes the gallery under a generated public name built from your prefix, with your publisher email and EULA attached.

## Deploy

### Console

Open the deployment store, find **Azure Compute Gallery**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Private Gallery** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f compute-gallery.yaml
```

## After Deploy

Create image definitions inside the gallery with **AzureComputeGalleryImage** (they reference the gallery by its `gallery_name` output) and publish versions from your image build pipeline. VMs deploy from a version's ARM ID or from the definition's ID to get the latest version.
