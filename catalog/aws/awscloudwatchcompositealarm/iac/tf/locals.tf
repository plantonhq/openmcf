locals {
  # The composite alarm's cloud name is the resource's metadata.name — the
  # same basis the Pulumi module uses, and the identity other composite alarms
  # address in their own rule expressions.
  resource_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key (the canonical
  # six-key identity map -- user labels never merge into cloud tags).
  aws_tags = {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsCloudwatchCompositeAlarm"
    "planton.ai/resource-id"   = var.metadata.id
  }

  alarm_description = var.spec.alarm_description != "" ? var.spec.alarm_description : null

  # actions_enabled is a genuine tri-state: null lets AWS default to true; an
  # explicit false is a real user choice. On this resource the provider marks
  # the flag ForceNew, so flipping it replaces the alarm.
  actions_enabled = var.spec.actions_enabled

  # Action ARNs arrive pre-resolved to plain strings by the orchestrator.
  alarm_actions             = var.spec.alarm_actions
  ok_actions                = var.spec.ok_actions
  insufficient_data_actions = var.spec.insufficient_data_actions
}
