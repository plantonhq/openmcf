variable "digitalocean_token" {
  description = "DigitalOcean API token for authentication"
  type        = string
  sensitive   = true
}

# Spaces is an S3-compatible credential plane the API token cannot reach.
# When null, the provider falls back to its own env defaults
# (SPACES_ACCESS_KEY_ID / SPACES_SECRET_ACCESS_KEY).
variable "spaces_access_id" {
  description = "DigitalOcean Spaces access key id"
  type        = string
  default     = null
  sensitive   = true
}

variable "spaces_secret_key" {
  description = "DigitalOcean Spaces secret access key"
  type        = string
  default     = null
  sensitive   = true
}
