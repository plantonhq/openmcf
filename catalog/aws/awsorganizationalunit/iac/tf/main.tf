# An organizational unit in the organization's OU tree.
#
# Lifecycle facts the render below depends on:
#   - the parent reference is immutable (AWS moves accounts between
#     OUs, never OUs themselves) - a parent change forces replacement;
#   - the display name renames in place;
#   - creation retries through the organization's finalization window
#     (the provider handles FinalizingOrganizationException for up to
#     four minutes after CreateOrganization);
#   - AWS identifies the OU as "ou-..." (the import ID).

resource "aws_organizations_organizational_unit" "this" {
  name      = var.spec.ou_name
  parent_id = var.spec.parent_id

  tags = local.aws_tags
}
