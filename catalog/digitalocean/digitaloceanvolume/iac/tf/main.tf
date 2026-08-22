# DigitalOcean block storage volume.
#
# Attachment to Droplets is a property of the Droplet (its volume_ids list),
# never of the volume. Size can only be EXPANDED after creation -- the
# provider rejects a shrink at plan time. Every other argument is create-only
# and replaces the volume when changed.
resource "digitalocean_volume" "this" {
  name = var.spec.volume_name

  # Create-only at the current provider pin: a description change REPLACES
  # the volume. Empty stays unset (the provider rejects empty strings).
  description = var.spec.description != "" ? var.spec.description : null

  # Volumes attach only to Droplets in the same region.
  region = var.spec.region

  size = var.spec.size_gib

  # Formatting happens once at creation; the API never reports these
  # arguments back (the resulting filesystem is observable through separate
  # computed attributes).
  initial_filesystem_type  = local.filesystem_type
  initial_filesystem_label = local.filesystem_label

  # When set, the volume is created from the snapshot, inheriting its region
  # and minimum size. Create-only, never reported back by the API.
  snapshot_id = var.spec.snapshot_id != "" ? var.spec.snapshot_id : null

  tags = local.tags
}
