locals {
  # The spec's two-value access enum renders the provider's free-string
  # canned ACL.
  acl = var.spec.access_control == "PUBLIC_READ" ? "public-read" : "private"
}
