locals {
  # metadata.name is the rule's cloud name (create-time immutable in AWS).
  resource_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsEventBridgeRule"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Event bus name — null targets the account's default bus. Note the AWS
  # constraint documented on the spec: schedule rules only run on the
  # default bus.
  event_bus_name = var.spec.event_bus_name != "" ? var.spec.event_bus_name : null

  # The event pattern arrives from the tfvars layer as a nested object (the
  # spec models it as a Struct); the provider wants a JSON string.
  event_pattern = var.spec.event_pattern != null ? jsonencode(var.spec.event_pattern) : null

  schedule_expression = var.spec.schedule_expression != "" ? var.spec.schedule_expression : null

  # State — null lets AWS default to ENABLED.
  state = var.spec.state != "" ? var.spec.state : null

  # Targets keyed by name: each materializes as its own provider resource,
  # so list edits add/remove targets in place instead of churning neighbors.
  targets = {
    for target in var.spec.targets : target.name => target
  }
}
