# The module-derived name, not the provider's `name` echo: the provider
# stores the SHORT name at create but reads back the FULL resource name
# (projects/{p}/locations/{l}/triggers/{name}) on refresh, so echoing it
# flips the output between plans (live-caught by the idempotency gate).
# The name is immutable and the create succeeded with exactly this value.
output "trigger_name" {
  description = "The trigger name in GCP"
  value       = local.trigger_name
}

# The full trigger resource name (projects/{p}/locations/{l}/triggers/
# {name}) with the ambient project resolved — the trigger's canonical
# API handle (the resource ID carries this form).
output "trigger_id" {
  description = "Full trigger resource name"
  value       = google_eventarc_trigger.this.id
}

# Partner triggers only: the one-time token the SaaS partner needs to
# complete the channel handshake. Empty for non-partner triggers.
output "partner_channel_activation_token" {
  description = "One-time partner channel activation token (sensitive)"
  value       = local.create_partner_channel ? google_eventarc_channel.this[0].activation_token : ""
  sensitive   = true
}
