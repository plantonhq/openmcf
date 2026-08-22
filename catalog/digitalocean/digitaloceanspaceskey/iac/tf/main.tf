# DigitalOcean Spaces Key
#
# Provisions an access-key pair for Spaces object storage -- the complete
# digitalocean_spaces_key resource surface. Keys are managed through
# DigitalOcean's REST API (the account token authenticates this module; no
# Spaces credentials are needed to CREATE credentials).
#
# The SECRET KEY EXISTS ONLY IN THE CREATE RESPONSE: DigitalOcean never
# returns it again, so the secret_key output is the only place it ever
# lives. Name and grants update in place (the grant list is replaced
# wholesale); the key material itself never changes.

resource "digitalocean_spaces_key" "key" {
  name = var.spec.key_name

  # The provider's grant grammar: read/readwrite grants name their bucket;
  # a fullaccess grant carries an EMPTY bucket string (spec validation
  # guarantees the pairing). Bucket references resolve to literal bucket
  # names before the module runs.
  dynamic "grant" {
    for_each = coalesce(var.spec.grants, [])
    content {
      bucket     = try(grant.value.bucket, "") != "" ? grant.value.bucket : ""
      permission = grant.value.permission
    }
  }
}
