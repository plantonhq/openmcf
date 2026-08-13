# The server-assigned resource name — the value alert policies reference in
# their notification_channels list.
output "channel_name" {
  description = "Resource name of the channel (projects/{project}/notificationChannels/{id})"
  value       = google_monitoring_notification_channel.this.name
}

# Verification state: SMS and email channels require verification before
# they deliver; types that need none report UNSPECIFIED (normal, not an
# error).
output "verification_status" {
  description = "Whether the channel has been verified"
  value       = google_monitoring_notification_channel.this.verification_status
}
