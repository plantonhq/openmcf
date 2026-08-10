locals {
  # Derive a stable resource ID
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # The trigger name defaults to metadata.name when the spec leaves
  # trigger_name empty — the same naming basis every kind uses.
  trigger_name = (
    var.spec.trigger_name != null && var.spec.trigger_name != ""
    ? var.spec.trigger_name
    : var.metadata.name
  )

  # Which companions the spec arms (count guards key off these).
  create_partner_channel       = var.spec.partner_channel != null
  create_google_channel_config = var.spec.google_channel_crypto_key != ""

  # The partner channel name defaults to "{trigger name}-channel".
  channel_name = (
    local.create_partner_channel && var.spec.partner_channel.channel_name != ""
    ? var.spec.partner_channel.channel_name
    : "${local.trigger_name}-channel"
  )

  # The same planton-ai_* label set the Pulumi module applies, so a resource
  # is attributable to its Planton object regardless of the engine that
  # created it.
  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.metadata.name
    "planton-ai_kind"     = "gcpeventarctrigger"
  }

  org_label = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "planton-ai_organization" = var.metadata.org } : {}

  env_label = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "planton-ai_environment" = var.metadata.env } : {}

  id_label = (
    var.metadata.id != null && var.metadata.id != ""
  ) ? { "planton-ai_id" = var.metadata.id } : {}

  # User labels first: the platform labels win on key conflicts. Applied to
  # the trigger and to the partner channel.
  final_labels = merge(var.spec.labels, local.base_labels, local.org_label, local.env_label, local.id_label)

  # Empty deletion_policy defers to the provider default (DELETE). One spec
  # lever wired to the trigger and both companions.
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null
}
