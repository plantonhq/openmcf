# Azure Compute Gallery -- Operational Guide

Judgment calls that matter when you run Compute Galleries in production.

## Decide the sharing posture before you create

The whole sharing tree -- Private, Groups, or Community, and every community publishing detail -- is create-only: changing any of it replaces the gallery, and replacing a gallery means recreating every definition and republishing every version inside it. Treat the sharing decision as a naming-level decision: almost every gallery should start (and stay) Private, with RBAC granting consuming subscriptions read access; reach for Community only when you genuinely publish images to the world.

## Name galleries like packages, not like resources

Gallery names forbid dashes and allow dots -- the platform convention that works is a package-style name (`platform.images`, `data.golden`) that reads well in the image ARM IDs every consumer sees. The name is also permanent (ForceNew), and image references embed it everywhere, so a rename is an estate-wide migration; choose once.

## One gallery per audience, not per image

The gallery is the ACL boundary and the publishing boundary -- definitions inside it share the gallery's visibility. Split galleries by who consumes them (platform-wide, team-private, public/community), never by OS or application; that is what image definitions are for. A single well-named gallery with many definitions is easier to govern than many small galleries with one definition each.

## Community publishing is a storefront, not a switch

Community sharing attaches your EULA, publisher email, and URI to a PUBLIC identity Azure generates from your prefix (the `community_gallery_name` output). Everything published in that gallery becomes publicly deployable. Keep community galleries SEPARATE from internal ones, review what gets a version published into them, and remember the prefix is create-only -- the public name cannot be rebranded in place.
