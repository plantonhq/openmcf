output "policy_id" {
  description = "The Cloudflare-assigned UUID of the notification policy"
  value       = cloudflare_notification_policy.main.id
}
