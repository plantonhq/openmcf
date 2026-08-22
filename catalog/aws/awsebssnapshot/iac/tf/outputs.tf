output "snapshot_id" {
  description = "The snapshot's id (snap-...) - what volume restores, copies, and permission grants reference, and the provider's import ID (volume arm only)"
  value       = local.snapshot_id
}

output "snapshot_arn" {
  description = "The snapshot's ARN"
  value       = local.snapshot_arn
}

output "owner_id" {
  description = "The AWS account that owns the snapshot"
  value       = local.snapshot_owner_id
}

output "volume_size_gb" {
  description = "The size (GiB) of the volume the snapshot captures - for imports, the size AWS derived from the disk image"
  value       = tostring(local.snapshot_volume_size)
}
