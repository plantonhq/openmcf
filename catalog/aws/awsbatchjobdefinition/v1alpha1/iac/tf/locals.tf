locals {
  # Resource-identity tags follow the catalog convention. With
  # spec.propagate_tags they also reach the ECS tasks jobs run as.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBatchJobDefinition"
    "planton.ai/resource-id"   = var.metadata.id
  }

  container = var.spec.container

  # ---------------------------------------------------------------------
  # containerProperties document.
  #
  # The provider takes containerProperties as an opaque JSON string and
  # compares it semantically, so the one requirement beyond correctness is
  # that this document and the Pulumi module's serializer produce
  # semantically identical JSON for the same spec: Batch API camelCase
  # names, absent optionals ABSENT (never null -- merge() of conditional
  # maps, since jsonencode renders null attributes), and map-derived lists
  # sorted by name for determinism.
  # ---------------------------------------------------------------------

  # Sizing goes through resourceRequirements (the modern shape; the
  # top-level vcpus/memory API fields are deprecated). Values are API
  # strings.
  resource_requirements = concat(
    [
      { type = "VCPU", value = tostring(local.container.vcpus) },
      { type = "MEMORY", value = tostring(local.container.memory_mib) },
    ],
    local.container.gpus > 0 ? [{ type = "GPU", value = tostring(local.container.gpus) }] : [],
  )

  volumes = [
    for volume in local.container.volumes : merge(
      { name = volume.name },
      volume.efs != null ? {
        efsVolumeConfiguration = merge(
          {
            fileSystemId = volume.efs.file_system_id
            # Transit encryption is always on: AWS requires it for access
            # points and IAM auth, and there is no good reason to mount EFS
            # unencrypted in transit.
            transitEncryption = "ENABLED"
          },
          volume.efs.root_directory != "" ? { rootDirectory = volume.efs.root_directory } : {},
          volume.efs.access_point_id != "" || volume.efs.iam_authorization ? {
            authorizationConfig = merge(
              volume.efs.access_point_id != "" ? { accessPointId = volume.efs.access_point_id } : {},
              volume.efs.iam_authorization ? { iam = "ENABLED" } : {},
            )
          } : {},
        )
      } : {},
      volume.host_path != "" ? { host = { sourcePath = volume.host_path } } : {},
    )
  ]

  # A ternary whose false arm is `{}` and whose true arm is a merge() result
  # cannot type-unify in HCL ("Inconsistent conditional result types"), so
  # each linuxParameters sub-field is guarded individually with try() -- when
  # the whole block is null every guard is false and the merge collapses to
  # an empty object.
  linux_parameters_inner = merge(
    try(local.container.linux_parameters.init_process_enabled, false) ? { initProcessEnabled = true } : {},
    try(local.container.linux_parameters.shared_memory_size_mib, 0) > 0 ? { sharedMemorySize = local.container.linux_parameters.shared_memory_size_mib } : {},
    try(local.container.linux_parameters.max_swap_mib, 0) > 0 ? { maxSwap = local.container.linux_parameters.max_swap_mib } : {},
    try(local.container.linux_parameters.swappiness, 0) > 0 ? { swappiness = local.container.linux_parameters.swappiness } : {},
    length(try(local.container.linux_parameters.tmpfs, [])) > 0 ? {
      tmpfs = [
        for mount in local.container.linux_parameters.tmpfs : merge(
          { containerPath = mount.container_path, size = mount.size_mib },
          length(mount.mount_options) > 0 ? { mountOptions = mount.mount_options } : {},
        )
      ]
    } : {},
    length(try(local.container.linux_parameters.devices, [])) > 0 ? {
      devices = [
        for device in local.container.linux_parameters.devices : merge(
          { hostPath = device.host_path },
          device.container_path != "" ? { containerPath = device.container_path } : {},
          length(device.permissions) > 0 ? { permissions = device.permissions } : {},
        )
      ]
    } : {},
  )

  container_properties = merge(
    {
      image                = local.container.image
      resourceRequirements = local.resource_requirements
    },
    length(local.container.command) > 0 ? { command = local.container.command } : {},
    local.container.job_role != "" ? { jobRoleArn = local.container.job_role } : {},
    local.container.execution_role != "" ? { executionRoleArn = local.container.execution_role } : {},
    length(local.container.environment) > 0 ? {
      environment = [
        for name in sort(keys(local.container.environment)) : {
          name  = name
          value = local.container.environment[name]
        }
      ]
    } : {},
    length(local.container.secrets) > 0 ? {
      secrets = [
        for name in sort(keys(local.container.secrets)) : {
          name      = name
          valueFrom = local.container.secrets[name]
        }
      ]
    } : {},
    local.container.log_configuration != null ? {
      logConfiguration = merge(
        { logDriver = try(local.container.log_configuration.log_driver, "") },
        length(try(local.container.log_configuration.options, {})) > 0 ? {
          options = local.container.log_configuration.options
        } : {},
        length(try(local.container.log_configuration.secret_options, {})) > 0 ? {
          secretOptions = [
            for name in sort(keys(local.container.log_configuration.secret_options)) : {
              name      = name
              valueFrom = local.container.log_configuration.secret_options[name]
            }
          ]
        } : {},
      )
    } : {},
    length(local.container.mount_points) > 0 ? {
      mountPoints = [
        for mount in local.container.mount_points : {
          sourceVolume  = mount.source_volume
          containerPath = mount.container_path
          readOnly      = mount.read_only
        }
      ]
    } : {},
    length(local.container.volumes) > 0 ? { volumes = local.volumes } : {},
    length(local.container.ulimits) > 0 ? {
      ulimits = [
        for ulimit in local.container.ulimits : {
          name      = ulimit.name
          softLimit = ulimit.soft_limit
          hardLimit = ulimit.hard_limit
        }
      ]
    } : {},
    length(local.linux_parameters_inner) > 0 ? { linuxParameters = local.linux_parameters_inner } : {},
    local.container.privileged ? { privileged = true } : {},
    local.container.user != "" ? { user = local.container.user } : {},
    local.container.readonly_root_filesystem ? { readonlyRootFilesystem = true } : {},
    local.container.repository_credentials_secret_arn != "" ? {
      repositoryCredentials = { credentialsParameter = local.container.repository_credentials_secret_arn }
    } : {},
    local.container.runtime_platform != null ? {
      runtimePlatform = merge(
        try(local.container.runtime_platform.cpu_architecture, "") != "" ? { cpuArchitecture = local.container.runtime_platform.cpu_architecture } : {},
        try(local.container.runtime_platform.operating_system_family, "") != "" ? { operatingSystemFamily = local.container.runtime_platform.operating_system_family } : {},
      )
    } : {},
    local.container.fargate_platform_version != "" ? {
      fargatePlatformConfiguration = { platformVersion = local.container.fargate_platform_version }
    } : {},
    local.container.assign_public_ip ? {
      networkConfiguration = { assignPublicIp = "ENABLED" }
    } : {},
    local.container.ephemeral_storage_gib > 0 ? {
      ephemeralStorage = { sizeInGiB = local.container.ephemeral_storage_gib }
    } : {},
  )
}
