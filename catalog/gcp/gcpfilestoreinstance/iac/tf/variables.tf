variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Specification for the Filestore instance"
  type = object({
    # StringValueOrRef fields arrive from the proto→tfvars converter as
    # plain strings (already resolved), never as object({value}).
    # Empty falls back to the provider's default project.
    project_id = optional(string, "")

    # Instance name; empty falls back to metadata.name. Immutable.
    instance_name = optional(string, "")

    # Zone for zonal tiers, region for ENTERPRISE/REGIONAL. Immutable.
    location = string

    # Service tier. Immutable.
    tier = string

    description = optional(string, "")

    # NFS_V3 (default) or NFS_V4_1. Immutable.
    protocol = optional(string, "")

    # Resolved CMEK key id (projects/../cryptoKeys/..). Immutable.
    kms_key_name = optional(string, "")

    deletion_protection_enabled = optional(bool, false)
    deletion_protection_reason  = optional(string, "")

    # The single file share on the instance.
    file_share = object({
      name        = string
      capacity_gb = number
      nfs_export_options = optional(list(object({
        ip_ranges   = optional(list(string), [])
        access_mode = optional(string, "")
        squash_mode = optional(string, "")
        anon_uid    = optional(number)
        anon_gid    = optional(number)
        # Source VPC network (name) for ip_ranges; required by GCP for
        # PSC instances, optional otherwise.
        network = optional(string, "")
      })), [])
      # Restore from an existing Filestore backup. Create-time only.
      source_backup = optional(string, "")
      # Restore from a Backup and DR Service backup. Create-time only;
      # mutually exclusive with source_backup (CEL-enforced pre-deploy).
      source_backupdr_backup = optional(string, "")
    })

    # The single VPC attachment. Immutable.
    network_config = object({
      # Resolved VPC network (name or self link).
      network           = string
      connect_mode      = optional(string, "")
      reserved_ip_range = optional(string, "")
      # IP versions; empty means ["MODE_IPV4"].
      modes = optional(list(string), [])
      # Consumer project for the PSC endpoint; only meaningful with
      # connect_mode PRIVATE_SERVICE_CONNECT (CEL-enforced pre-deploy).
      psc_endpoint_project = optional(string, "")
    })

    # IOPS tuning (ZONAL/REGIONAL/ENTERPRISE tiers).
    performance_config = optional(object({
      fixed_iops = optional(object({
        max_iops = number
      }), null)
      iops_per_tb = optional(object({
        max_iops_per_tb = number
      }), null)
    }), null)

    # Create-time cross-instance replication (this instance's role +
    # resolved peer instance paths).
    initial_replication = optional(object({
      role           = optional(string, "")
      peer_instances = list(string)
    }), null)

    # User labels merged beneath the platform attribution labels.
    labels = optional(map(string), {})

    # Resource Manager tags (tagKeys/{id} => tagValues/{id}). Create-time.
    tags = optional(map(string), {})

    # LDAP directory services for NFSv4.1 identity mapping. Requires
    # protocol NFS_V4_1 (CEL-enforced pre-deploy).
    ldap = optional(object({
      domain    = string
      servers   = list(string)
      groups_ou = optional(string, "")
      users_ou  = optional(string, "")
    }), null)

    # Replica-relationship state: READY (replicating, the default) or
    # PAUSED. A virtual lever the provider drives via pause/resume
    # replica calls; no effect on instances without a replica pair.
    desired_replica_state = optional(string, "READY")

    # Client-side destroy behavior: DELETE (default), PREVENT, ABANDON.
    deletion_policy = optional(string, "")
  })
}
