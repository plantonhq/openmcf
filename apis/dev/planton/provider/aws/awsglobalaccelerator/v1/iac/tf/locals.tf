locals {
  # The accelerator's cloud name is the resource's metadata.name — Global
  # Accelerator has a real name attribute (1-255 chars, alphanumeric and
  # hyphens), so the name is carried on the resource itself rather than
  # through a Name tag. Same basis as the Pulumi module.
  name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key. Identity
  # tagging is the only tagging surface this module manages; user-defined
  # custom tags are a platform-wide concern, not per-kind spec surface.
  tags = {
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsGlobalAccelerator"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Flatten listeners into a map keyed by listener name for for_each iteration.
  # The listener name is a Planton-side key (AWS listeners are anonymous); it
  # keys the resource address and the listener_arns output map, so renaming a
  # listener replaces that listener resource.
  listeners_map = {
    for listener in var.spec.listeners : listener.name => listener
  }

  # Flatten all endpoint groups across all listeners into a map keyed by
  # "listener_name/group_name" for for_each iteration. The composite key keeps
  # group names unique per listener without forcing global uniqueness, and is
  # the key format of the endpoint_group_arns output map.
  endpoint_groups_map = {
    for pair in flatten([
      for listener in var.spec.listeners : [
        for group in listener.endpoint_groups : {
          key           = "${listener.name}/${group.name}"
          listener_name = listener.name
          group         = group
        }
      ]
    ]) : pair.key => pair
  }
}
