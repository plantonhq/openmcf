# DigitalOcean Droplet
#
# Provisions a virtual machine on DigitalOcean, modeling the complete
# digitalocean_droplet resource surface: base image and sizing, region and
# VPC placement, SSH key injection, automated backups with a policy window,
# IPv6 and public-network toggles, the monitoring and web-console agents,
# block volume attachments, cloud-init user data, tags, GPU partitioning,
# and the resize/shutdown behavior flags.

resource "digitalocean_droplet" "this" {
  name  = var.spec.droplet_name
  image = var.spec.image
  size  = var.spec.size

  # Placement (both create-only). Unset region lets DigitalOcean choose;
  # unset VPC means the region's default VPC.
  region   = local.region
  vpc_uuid = local.vpc_uuid

  # SSH keys are create-only: the standard access path to the droplet.
  ssh_keys = local.ssh_keys

  # Networking toggles. public_networking and ipv6-disable force
  # recreation; null defers to DigitalOcean's defaults (public on).
  ipv6              = var.spec.enable_ipv6
  public_networking = var.spec.public_networking

  # Automated backups and their optional policy window (policy requires
  # backups — enforced by the spec's validation, mirrored by the provider).
  backups = var.spec.enable_backups
  dynamic "backup_policy" {
    for_each = var.spec.backup_policy != null ? [var.spec.backup_policy] : []
    content {
      plan    = try(backup_policy.value.plan, "") != "" ? backup_policy.value.plan : null
      weekday = try(backup_policy.value.weekday, "") != "" ? backup_policy.value.weekday : null
      # Hour 0 is a real window start (midnight), so it is always sent.
      hour = backup_policy.value.hour
    }
  }

  # Agents (both create-only). droplet_agent is tri-state: null lets
  # DigitalOcean install where the image supports it, true makes install
  # failures fatal, false prevents installation.
  monitoring    = var.spec.monitoring
  droplet_agent = var.spec.droplet_agent

  # Resize behavior: null defers to DigitalOcean's default (grow the disk
  # permanently). Never coalesced to false — that would silently flip the
  # provider default.
  resize_disk = var.spec.resize_disk

  # Destroy-time behavior: ACPI power-off before delete when true.
  graceful_shutdown = var.spec.graceful_shutdown

  # GPU partitioning, create-only, GPU sizes only.
  gpu_partition_mode = local.gpu_partition_mode

  # Block volume attachments; empty is sent as null so the computed set
  # never diffs.
  volume_ids = length(local.volume_ids) > 0 ? local.volume_ids : null

  # User tags plus Planton labels — how firewalls and load balancers
  # target this droplet.
  tags = local.tags

  # Cloud-init user data (create-only, hash-stored).
  user_data = local.user_data
}
