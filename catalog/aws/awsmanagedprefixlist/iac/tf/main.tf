# One customer-managed prefix list with its entries managed in-line.
#
# Lifecycle facts the render below depends on:
#   - address_family is fixed for life (replace-on-change) - every
#     referencing rule breaks with the old pl- id;
#   - the in-line entry set is the single declarative owner; the
#     standalone aws_ec2_managed_prefix_list_entry resource is the
#     identical payload and fights this form, so this module never
#     uses it;
#   - AWS versions the list on every entry change (the version output);
#     the provider orders max_entries increases BEFORE entry changes
#     and decreases AFTER, so a resize never transiently strands
#     entries;
#   - a description-only edit costs two API round trips (the provider
#     removes and re-adds the entry) - expected, not drift.

resource "aws_ec2_managed_prefix_list" "this" {
  name           = var.metadata.name
  address_family = var.spec.address_family
  max_entries    = var.spec.max_entries

  dynamic "entry" {
    for_each = var.spec.entries != null ? var.spec.entries : []
    content {
      cidr        = entry.value.cidr
      description = entry.value.description != "" ? entry.value.description : null
    }
  }

  tags = local.aws_tags
}
