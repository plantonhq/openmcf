# The resource ID is the fully qualified key ring path
# (projects/{project}/locations/{location}/keyRings/{name}) — the exact
# string a GcpKmsKey's key_ring_id reference consumes.
output "key_ring_id" {
  description = "Fully qualified key ring resource path"
  value       = google_kms_key_ring.this.id
}

output "key_ring_name" {
  description = "The short name of the key ring"
  value       = google_kms_key_ring.this.name
}

# Consumers that take a bare ring name plus a separately supplied location
# (rather than the fully qualified path) compose from key_ring_name + this.
output "location" {
  description = "Location the key ring resides in"
  value       = google_kms_key_ring.this.location
}
