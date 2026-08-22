# DigitalOcean SSH Key
#
# Registers an SSH public key on the DigitalOcean account, modeling the
# complete digitalocean_ssh_key resource surface. Only the name updates in
# place; the key material is create-only upstream (the provider compares it
# after trimming only leading/trailing whitespace, so any in-line change --
# a different comment, a re-encoded body -- REPLACES the key and rotates
# the numeric id and fingerprint droplets reference).

resource "digitalocean_ssh_key" "ssh_key" {
  name       = var.spec.key_name
  public_key = var.spec.public_key
}
