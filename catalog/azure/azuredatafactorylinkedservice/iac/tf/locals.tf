# A linked service carries no tags (ARM sub-resources of a factory
# expose none), so there is no tag map to derive; the metadata
# variable stays part of the module contract for uniformity.
locals {
  linked_service_name = var.spec.name
}
