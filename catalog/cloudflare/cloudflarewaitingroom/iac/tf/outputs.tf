output "waiting_room_id" {
  description = "The ID of the created waiting room (what events reference)"
  value       = cloudflare_waiting_room.main.id
}

output "zone_id" {
  description = "The zone the waiting room belongs to"
  value       = var.spec.zone_id
}
