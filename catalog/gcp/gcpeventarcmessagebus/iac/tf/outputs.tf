# The full bus resource name — the value cross-bus pipelines (another
# bus's destination.message_bus) and external publishers consume.
output "message_bus_name" {
  description = "Full message bus resource name"
  value       = google_eventarc_message_bus.this.name
}
