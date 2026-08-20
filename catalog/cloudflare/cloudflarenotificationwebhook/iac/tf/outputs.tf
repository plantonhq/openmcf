output "webhook_id" {
  description = "The Cloudflare-assigned UUID of the webhook destination (what notification policies reference)"
  value       = cloudflare_notification_policy_webhooks.main.id
}

output "type" {
  description = "The destination type Cloudflare inferred from the URL (datadog, discord, feishu, gchat, generic, opsgenie, slack, splunk)"
  value       = cloudflare_notification_policy_webhooks.main.type
}
