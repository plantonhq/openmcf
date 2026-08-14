# The provider addresses domain topics by (resource group, domain
# name) while the spec takes the domain's ARM id (the composable
# reference shape) -- the same ARM object either way, so the id is
# split into its segments here.
#
# The domain id's shape is /subscriptions/{sub}/resourceGroups/{rg}
# /providers/Microsoft.EventGrid/domains/{domain}. Segment names are
# matched case-insensitively (ARM treats them that way), so ids
# composed by hand survive the split.
#
# No tags: the provider carries no tags argument on domain topics
# (they are addressing entries under the domain).
locals {
  domain_id_segments = split("/", var.spec.domain_id)

  resource_group_name = [
    for i, segment in local.domain_id_segments :
    local.domain_id_segments[i + 1]
    if lower(segment) == "resourcegroups"
  ][0]

  domain_name = [
    for i, segment in local.domain_id_segments :
    local.domain_id_segments[i + 1]
    if lower(segment) == "domains"
  ][0]
}
