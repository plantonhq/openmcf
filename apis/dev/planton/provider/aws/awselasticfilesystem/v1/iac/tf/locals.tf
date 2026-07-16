locals {
  # The file system's only human-controlled physical identity is the creation
  # token (EFS has no "name" argument; the console name is the Name tag).
  # metadata.name is the same basis the Pulumi module pins, keeping the two
  # engines' physical identity converged.
  resource_name = var.metadata.name

  # Resource-identity tags follow the catalog convention; user labels merge in
  # without being able to override the identity keys.
  aws_tags = merge(try(var.metadata.labels, {}), {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsElasticFileSystem"
    "planton.ai/resource-id"   = var.metadata.id
  })

  # Mount targets keyed by subnet ID (references arrive pre-resolved as plain
  # strings). AWS allows at most one mount target per AZ, so subnet IDs are
  # unique keys by construction; a duplicate-AZ mistake fails at the AWS API.
  mount_targets = { for mt in var.spec.mount_targets : mt.subnet_id => mt }

  # Empty security_group_ids means "let AWS attach the VPC default SG" --
  # pass null so the provider omits the argument entirely.
  security_groups = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null

  # The provider models each lifecycle transition as its own lifecycle_policy
  # element (up to 3: IA, Archive, back-to-primary).
  lifecycle_policies = concat(
    var.spec.transition_to_ia != "" ? [{
      transition_to_ia = var.spec.transition_to_ia
    }] : [],
    var.spec.transition_to_archive != "" ? [{
      transition_to_archive = var.spec.transition_to_archive
    }] : [],
    var.spec.transition_to_primary_storage_class != "" ? [{
      transition_to_primary_storage_class = var.spec.transition_to_primary_storage_class
    }] : []
  )

  # Null when not specified: EFS then encrypts with the AWS-managed key
  # aws/elasticfilesystem (when encryption is enabled at all).
  kms_key_id = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null

  # Null for regional (multi-AZ) file systems; a value selects One Zone
  # storage in that AZ. ForceNew.
  availability_zone_name = var.spec.availability_zone_name != "" ? var.spec.availability_zone_name : null

  # Only meaningful in provisioned mode (the coupling is CEL-enforced at
  # validation time); null otherwise so the provider omits the argument.
  provisioned_throughput_in_mibps = var.spec.throughput_mode == "provisioned" && var.spec.provisioned_throughput_in_mibps > 0 ? var.spec.provisioned_throughput_in_mibps : null

  # Empty strings become null so unset stays indistinguishable from the AWS
  # defaults (generalPurpose / bursting / protection ENABLED).
  performance_mode                 = var.spec.performance_mode != "" ? var.spec.performance_mode : null
  throughput_mode                  = var.spec.throughput_mode != "" ? var.spec.throughput_mode : null
  replication_overwrite_protection = var.spec.replication_overwrite_protection != "" ? var.spec.replication_overwrite_protection : null
}
