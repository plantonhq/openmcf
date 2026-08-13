# A user grant carries no tags (ARM data-plane user entries are
# untagged -- the provider exposes none), so there is no tag map to
# derive; the metadata variable stays part of the module contract for
# uniformity.
locals {
  object_id = var.spec.object_id
}
