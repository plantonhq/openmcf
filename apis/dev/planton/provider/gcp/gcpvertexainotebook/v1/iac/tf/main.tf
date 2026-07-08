# Enable the Notebooks API — the control plane that owns Workbench
# instances. disable_on_destroy is false: tearing down one notebook must
# never disable the API for everything else in the project.
resource "google_project_service" "notebooks_api" {
  project = local.project_id
  service = "notebooks.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Workbench VM itself is a Compute Engine instance provisioned by the
# Notebooks service agent, so a fresh project also needs the Compute API.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Workbench instance — a managed JupyterLab VM. The resource's `name`
# is the instance ID GCP uses in the create call, so it is the single
# identity field (the provider's separate instance_id attribute is
# vestigial at create time and deliberately not sent).
resource "google_workbench_instance" "this" {
  name     = local.instance_name
  location = local.location
  project  = local.project_id
  labels   = local.final_labels

  disable_proxy_access = var.spec.disable_proxy_access ? true : null

  # desired_state drives declarative stop/start: ACTIVE boots the VM,
  # STOPPED suspends compute billing while keeping disks.
  desired_state   = var.spec.desired_state != "" ? var.spec.desired_state : null
  instance_owners = length(var.spec.instance_owners) > 0 ? var.spec.instance_owners : null

  # Managed end-user credentials: JupyterLab runs as the accessing user's
  # own identity instead of the VM service account.
  enable_managed_euc          = var.spec.enable_managed_euc ? true : null
  enable_third_party_identity = var.spec.enable_third_party_identity ? true : null

  gce_setup {
    machine_type = var.spec.machine_type

    # Boot disk. disk_encryption is derived, never spec-set: presence of a
    # KMS key means CMEK, absence means Google-managed encryption.
    dynamic "boot_disk" {
      for_each = var.spec.boot_disk != null ? [var.spec.boot_disk] : []
      content {
        disk_type       = boot_disk.value.disk_type != "" ? boot_disk.value.disk_type : null
        disk_size_gb    = boot_disk.value.disk_size_gb != 0 ? tostring(boot_disk.value.disk_size_gb) : null
        disk_encryption = boot_disk.value.kms_key != "" ? "CMEK" : null
        kms_key         = boot_disk.value.kms_key != "" ? boot_disk.value.kms_key : null
      }
    }

    # Data disk (the API supports exactly one).
    dynamic "data_disks" {
      for_each = var.spec.data_disk != null ? [var.spec.data_disk] : []
      content {
        disk_type       = data_disks.value.disk_type != "" ? data_disks.value.disk_type : null
        disk_size_gb    = data_disks.value.disk_size_gb != 0 ? tostring(data_disks.value.disk_size_gb) : null
        disk_encryption = data_disks.value.kms_key != "" ? "CMEK" : null
        kms_key         = data_disks.value.kms_key != "" ? data_disks.value.kms_key : null
      }
    }

    # GPU accelerator (the API supports exactly one configuration).
    dynamic "accelerator_configs" {
      for_each = var.spec.accelerator_config != null && var.spec.accelerator_config.type != "" ? [var.spec.accelerator_config] : []
      content {
        type       = accelerator_configs.value.type
        core_count = accelerator_configs.value.core_count != 0 ? tostring(accelerator_configs.value.core_count) : null
      }
    }

    # Network interface (the API supports exactly one). An external_ip
    # pins a static address (ONE_TO_ONE_NAT); without it — and with public
    # IP not disabled — GCP assigns an ephemeral external address that
    # changes across stop/start cycles.
    dynamic "network_interfaces" {
      for_each = var.spec.network_interface != null ? [var.spec.network_interface] : []
      content {
        network  = network_interfaces.value.network != "" ? network_interfaces.value.network : null
        subnet   = network_interfaces.value.subnet != "" ? network_interfaces.value.subnet : null
        nic_type = network_interfaces.value.nic_type != "" ? network_interfaces.value.nic_type : null

        dynamic "access_configs" {
          for_each = network_interfaces.value.external_ip != "" ? [network_interfaces.value.external_ip] : []
          content {
            external_ip = access_configs.value
          }
        }
      }
    }

    disable_public_ip    = var.spec.disable_public_ip ? true : null
    enable_ip_forwarding = var.spec.enable_ip_forwarding ? true : null

    # VM identity (the API supports exactly one service account); scopes
    # are fixed to cloud-platform by the Workbench service.
    dynamic "service_accounts" {
      for_each = var.spec.service_account != "" ? [var.spec.service_account] : []
      content {
        email = service_accounts.value
      }
    }

    tags     = length(var.spec.tags) > 0 ? var.spec.tags : null
    metadata = length(var.spec.metadata) > 0 ? var.spec.metadata : null

    # VM image (mutually exclusive with container_image — enforced
    # pre-deploy by the spec's CEL rule).
    dynamic "vm_image" {
      for_each = var.spec.vm_image != null ? [var.spec.vm_image] : []
      content {
        project = vm_image.value.project != "" ? vm_image.value.project : null
        family  = vm_image.value.family != "" ? vm_image.value.family : null
        name    = vm_image.value.name != "" ? vm_image.value.name : null
      }
    }

    # Container image (mutually exclusive with vm_image).
    dynamic "container_image" {
      for_each = var.spec.container_image != null ? [var.spec.container_image] : []
      content {
        repository = container_image.value.repository
        tag        = container_image.value.tag != "" ? container_image.value.tag : null
      }
    }

    # Shielded VM posture (rootkit/bootkit protection). False values are
    # sent as null (omitted) rather than explicit false: the API enables
    # vTPM and integrity monitoring by default, and an explicit false would
    # actively disable them — the Pulumi module omits false the same way,
    # so both engines leave server defaults intact for unset flags.
    dynamic "shielded_instance_config" {
      for_each = var.spec.shielded_instance_config != null ? [var.spec.shielded_instance_config] : []
      content {
        enable_secure_boot          = shielded_instance_config.value.enable_secure_boot ? true : null
        enable_vtpm                 = shielded_instance_config.value.enable_vtpm ? true : null
        enable_integrity_monitoring = shielded_instance_config.value.enable_integrity_monitoring ? true : null
      }
    }

    # Confidential Computing (AMD SEV): guest memory encrypted in use.
    # Requires an SEV-capable machine type (n2d family).
    dynamic "confidential_instance_config" {
      for_each = var.spec.confidential_instance_config != null ? [var.spec.confidential_instance_config] : []
      content {
        confidential_instance_type = confidential_instance_config.value.confidential_instance_type != "" ? confidential_instance_config.value.confidential_instance_type : null
      }
    }

    # Reservation affinity: consume pre-purchased Compute capacity — how
    # organizations guarantee GPU availability for ML workloads.
    dynamic "reservation_affinity" {
      for_each = var.spec.reservation_affinity != null ? [var.spec.reservation_affinity] : []
      content {
        consume_reservation_type = reservation_affinity.value.consume_reservation_type != "" ? reservation_affinity.value.consume_reservation_type : null
        key                      = reservation_affinity.value.key != "" ? reservation_affinity.value.key : null
        values                   = length(reservation_affinity.value.values) > 0 ? reservation_affinity.value.values : null
      }
    }
  }

  depends_on = [
    google_project_service.notebooks_api,
    google_project_service.compute_api,
  ]
}
