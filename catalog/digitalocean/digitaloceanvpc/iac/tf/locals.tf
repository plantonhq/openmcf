locals {
  # The VPC's name is its Planton identity -- it comes from metadata.name,
  # never from separate spec surface.
  vpc_name = var.metadata.name
}
