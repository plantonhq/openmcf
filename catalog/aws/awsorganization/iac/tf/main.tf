# THE AWS Organization of the deploying account - creating it makes
# the caller the management account.
#
# Lifecycle facts the render below depends on:
#   - feature_set upgrades (CONSOLIDATED_BILLING -> ALL) apply in place
#     via EnableAllFeatures; the DOWNGRADE forces replacement, which is
#     delete-and-recreate of the ENTIRE organization (AWS only permits
#     it once every member account, OU, and policy is gone);
#   - service-access principals and policy types are applied as
#     enable/disable calls diffed on update (disables first); trusted
#     access, policy types, delegated administrators, and the resource
#     policy all require ALL features (spec CELs front-load this);
#   - destroy calls DeleteOrganization - the whole organization ends;
#   - the resource policy is a per-organization SINGLETON
#     (PutResourcePolicy upserts it), so the arm renders at most one;
#   - the organization resource is untaggable; the resource policy is
#     the kind's one taggable surface.

resource "aws_organizations_organization" "this" {
  # Empty keeps the provider default (ALL - the level every advanced
  # arm requires).
  feature_set = var.spec.feature_set != "" ? var.spec.feature_set : null

  aws_service_access_principals = length(var.spec.aws_service_access_principals) > 0 ? var.spec.aws_service_access_principals : null

  enabled_policy_types = length(var.spec.enabled_policy_types) > 0 ? var.spec.enabled_policy_types : null
}

# Each registration names one member account as the administrator for
# one AWS service. Both leaves are immutable - a change re-registers.
resource "aws_organizations_delegated_administrator" "this" {
  for_each = local.delegated_administrators

  account_id        = each.value.account_id
  service_principal = each.value.service_principal

  depends_on = [aws_organizations_organization.this]
}

# The organization's single resource-based delegation policy. The spec
# carries the document structured; it is serialized to JSON here.
resource "aws_organizations_resource_policy" "this" {
  count = var.spec.resource_policy != null ? 1 : 0

  content = jsonencode(var.spec.resource_policy)

  tags = local.aws_tags

  depends_on = [aws_organizations_organization.this]
}
