output "volume_id" {
  description = "The volume's id (vol-...) - what attachments, snapshots, and copies reference, and the provider's import ID"
  value       = local.volume_id
}

output "volume_arn" {
  description = "The volume's ARN"
  value       = local.volume_arn
}

output "availability_zone" {
  description = "The availability zone the volume actually lives in - notably useful for copies, which inherit the source's zone"
  value       = local.volume_availability_zone
}

output "size_gb" {
  description = "The volume's actual size in GiB - the snapshot's size when size_gb was left unset"
  value       = tostring(local.volume_size)
}

output "create_time" {
  description = "When AWS created the volume (RFC3339); empty on the copy arm (the provider does not expose it there)"
  value       = local.volume_create_time
}
