# The addressing choice selects which provider resource materializes:
# a scope-addressed subscription attaches to any ARM resource by id,
# while a system-topic subscription is a CHILD of the system topic and
# the provider addresses it by (resource group, system topic name) --
# the spec takes the topic's ARM id (the composable reference shape),
# split into its segments here.
#
# The system topic id's shape is /subscriptions/{sub}/resourceGroups/
# {rg}/providers/Microsoft.EventGrid/systemTopics/{name}. Segment
# names are matched case-insensitively (ARM treats them that way), so
# ids composed by hand survive the split.
#
# No tags: the provider carries no tags argument on event
# subscriptions (Event Grid stores free-form labels instead -- the
# spec's labels field).
locals {
  system_topic_id_segments = var.spec.system_topic_id != null ? split("/", var.spec.system_topic_id) : []

  system_topic_resource_group = var.spec.system_topic_id != null ? [
    for i, segment in local.system_topic_id_segments :
    local.system_topic_id_segments[i + 1]
    if lower(segment) == "resourcegroups"
  ][0] : ""

  system_topic_name = var.spec.system_topic_id != null ? [
    for i, segment in local.system_topic_id_segments :
    local.system_topic_id_segments[i + 1]
    if lower(segment) == "systemtopics"
  ][0] : ""

  # The spec enum's value names mapped to the provider's identity
  # tokens (delivery and dead-letter identities allow exactly these
  # two -- there is no combined mode on a subscription).
  identity_type_map = {
    "SYSTEM_ASSIGNED" = "SystemAssigned"
    "USER_ASSIGNED"   = "UserAssigned"
  }
}
