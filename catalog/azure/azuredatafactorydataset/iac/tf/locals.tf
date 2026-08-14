# A dataset carries no tags (ARM sub-resources of a factory expose
# none), so there is no tag map to derive; the metadata variable
# stays part of the module contract for uniformity.
locals {
  dataset_name = var.spec.name
}
