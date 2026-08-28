output "webhook_id" {
  description = "The Cloudflare-assigned UUID of the webhook destination (what notification policies reference)"
  value       = cloudflare_notification_policy_webhooks.main.id
}

output "type" {
  description = "The destination type Cloudflare inferred from the URL (datadog, discord, feishu, gchat, generic, opsgenie, slack, splunk)"
  # From the read-after-create data source, never the resource: the create
  # response omits type, so the resource attribute is empty until the first
  # refresh (see main.tf).
  value       = data.cloudflare_notification_policy_webhooks.main.type
}
