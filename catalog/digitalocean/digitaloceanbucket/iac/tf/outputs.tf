output "bucket_id" {
  description = "The provider's resource id for the bucket, which IS the bucket name"
  value       = digitalocean_spaces_bucket.main.id
}

output "endpoint" {
  description = "The region-level Spaces endpoint host (<region>.digitaloceanspaces.com)"
  value       = digitalocean_spaces_bucket.main.endpoint
}

output "region" {
  description = "The region slug the bucket lives in, read back from the API"
  value       = digitalocean_spaces_bucket.main.region
}

output "bucket_domain_name" {
  description = "The bucket's virtual-host-style FQDN (<bucket>.<region>.digitaloceanspaces.com)"
  value       = digitalocean_spaces_bucket.main.bucket_domain_name
}

output "urn" {
  description = "The uniform resource name of the bucket (do:space:<name>)"
  value       = digitalocean_spaces_bucket.main.urn
}
