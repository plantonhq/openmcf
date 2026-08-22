locals {
  # Optional pool user: null when unset creates DigitalOcean's "inbound
  # user" pool (clients authenticate with their own database credentials).
  # Reads echo the empty value back, so omission is drift-stable.
  user = try(var.spec.user, "") != "" ? var.spec.user : null
}
