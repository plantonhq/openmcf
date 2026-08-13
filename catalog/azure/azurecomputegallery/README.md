# Overview

The **AzureComputeGallery** component deploys an Azure Compute Gallery -- the shared library an organization keeps its approved VM images in. Image definitions (AzureComputeGalleryImage) live inside the gallery, each with published, region-replicated versions; VMs and scale sets deploy from those versions.

## Purpose

- **One home for golden images**: the gallery is the organizational unit -- per team, per platform, or per environment -- under which image definitions and their versions are published, shared, and consumed.
- **Sharing as a deliberate posture**: private (RBAC-only) by default; Groups sharing extends to chosen subscriptions/tenants, Community sharing publishes publicly under a generated community name.
- **Free at rest**: the gallery bills nothing itself -- image versions bill for the storage their replicas consume.

## Key Features

- Full azurerm v5 surface: description, the whole sharing tree (permission + community-gallery publishing details), and tags.
- Chart-ready: `resource_group` defaults its reference to AzureResourceGroup; the `gallery_name` output is what AzureComputeGalleryImage definitions reference; `unique_name` and `community_gallery_name` carry Azure's generated identities.
- The provider's expand-time contract (Community sharing requires the community block) is front-loaded as a spec validation rule, so a bad manifest fails before it reaches Azure.

## Use Cases

- **The platform team's image library**: one gallery per platform, definitions per OS flavor, versions published by the image build pipeline.
- **Cross-subscription image distribution**: share one gallery to every workload subscription instead of copying images around.
- **Public image publishing**: a Community-shared gallery puts your images in the public catalog under your publisher identity.

## Future Enhancements

- Gallery applications (VM app packages) are deliberately deferred -- a niche surface (`gallery_application` ×2 + the VM assignment) that enters the catalog by its own decision when demand shows.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
