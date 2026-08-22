output "event_id" {
  description = "The ID of the created waiting room event"
  value       = cloudflare_waiting_room_event.main.id
}

output "waiting_room_id" {
  description = "The waiting room the event runs on"
  value       = var.spec.waiting_room_id
}

output "zone_id" {
  description = "The zone the waiting room belongs to"
  value       = var.spec.zone_id
}
