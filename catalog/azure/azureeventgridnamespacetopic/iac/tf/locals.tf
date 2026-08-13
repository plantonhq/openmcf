# A namespace topic carries no tags (the provider exposes none -- it
# is a pure naming-and-retention entry), so there is no tag map to
# derive; the metadata variable stays part of the module contract for
# uniformity.
locals {
  namespace_topic_name = var.spec.name
}
