# DigitalOcean Droplet Autoscale Pool
#
# Provisions a pool of identical droplets DigitalOcean keeps at a fixed
# size or scales on CPU/memory utilization -- the complete
# digitalocean_droplet_autoscale resource surface. The spec's scaling
# oneof (static XOR dynamic) is rendered into the provider's single config
# block, which the provider itself never validates for shape consistency.
#
# Create waits for the pool AND every member droplet to reach "active"
# (up to 15 minutes upstream). DESTROY DESTROYS THE MEMBERS: the API's
# only delete is the dangerous variant that terminates every droplet the
# pool owns.

resource "digitalocean_droplet_autoscale" "pool" {
  name = var.spec.pool_name

  config {
    # Static branch: exact member count, nothing else.
    target_number_instances = local.is_dynamic ? null : var.spec.static.target_instances

    # Dynamic branch: bounds, utilization targets, cooldown.
    min_instances             = local.is_dynamic ? var.spec.dynamic.min_instances : null
    max_instances             = local.is_dynamic ? var.spec.dynamic.max_instances : null
    target_cpu_utilization    = local.is_dynamic ? var.spec.dynamic.target_cpu_utilization : null
    target_memory_utilization = local.is_dynamic ? var.spec.dynamic.target_memory_utilization : null
    cooldown_minutes          = local.is_dynamic ? var.spec.dynamic.cooldown_minutes : null
  }

  droplet_template {
    size = var.spec.droplet_template.size

    # The spec's region enum value names ARE the provider's region slugs.
    region = var.spec.droplet_template.region

    # Slug or numeric image id; DigitalOcean reads the image back as a
    # numeric id, and the provider itself persists the configured value to
    # avoid the drift.
    image = var.spec.droplet_template.image

    # SSH key references resolve to literal numeric key ids before the
    # module runs.
    ssh_keys = var.spec.droplet_template.ssh_keys

    vpc_uuid   = local.vpc_uuid
    project_id = local.project_id

    tags = local.tags

    with_droplet_agent = var.spec.droplet_template.with_droplet_agent
    ipv6               = var.spec.droplet_template.ipv6

    user_data = try(var.spec.droplet_template.user_data, "") != "" ? var.spec.droplet_template.user_data : null

    # public_networking is deliberately never rendered: the provider
    # declares it but never copies it into any create/update request --
    # dead on write at the pinned version.
  }
}
