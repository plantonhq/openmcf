# Azure Planton Runner Terraform Module
#
# Provisions a standing Planton runner appliance on Azure Container Apps:
# an always-on, outbound-only worker that executes deploy and cloud
# operations from inside your network perimeter -- the piece that makes
# private endpoints (most notably private AKS API servers) deployable and
# operable. The resource group and the Container App Environment are
# referenced resources -- the module never creates or mutates them (the
# environment decides the network boundary; a VNet-integrated one gives
# the runner reach into private endpoints).
#
# ENROLLMENT IS TOKEN-FIRST: the app ships the runner TOKEN (in the app's
# own secret store), never an identity. The runner joins the control plane
# on first boot, registers itself, and receives its own individually
# revocable identity; replica replacement re-joins with the same token
# (its lineage re-admits the runner it originally admitted).

resource "azurerm_container_app" "runner" {
  name                         = local.runner_name
  resource_group_name          = var.spec.resource_group
  container_app_environment_id = var.spec.container_app_environment_id

  # Single revision mode: new revisions replace the old one -- the only
  # sane model for a workload whose identity contract forbids two live
  # copies. (Revision rollover briefly overlaps two replicas; the draining
  # one's revoked key dies with it, and a rollback self-heals by
  # re-joining.)
  revision_mode = "Single"

  tags = local.azure_tags

  # The token lives in the app's own secret store; the container's env
  # references it by name -- never a plain env value, so reading the app
  # definition (a common, low-sensitivity permission) reveals nothing.
  secret {
    name  = local.token_secret_name
    value = var.spec.token
  }

  # NO ingress block at all: the runner accepts no inbound traffic -- it
  # initiates every connection it uses (control plane, work queue, image
  # pulls).

  template {
    # Exactly one runner per enrollment: a runner's identity is minted
    # for ONE live replica -- a second replica joining under the same
    # name would revoke the first's key (token lineage: re-admission
    # re-mints and revokes). Never enable scaling here without
    # redesigning enrollment for fleets.
    min_replicas = 1
    max_replicas = 1

    container {
      name   = "planton-runner"
      image  = "${var.spec.image_repository}:${var.spec.runner_version}"
      cpu    = local.cpu
      memory = local.memory

      # The runner binary's own start command -- the image's entrypoint
      # takes the subcommand as args.
      args = ["start"]

      # The runner's environment contract: the token via the app's secret
      # store, the name it registers itself under, and the control-plane
      # endpoint when one is declared (omitted, the runner's built-in
      # hosted default applies). No EXECUTION_MODE: only the control
      # plane knows whether its work queue is reachable, so the runner
      # derives its mode from the identity the join returns -- a mode
      # knob here would silently strip capability.
      env {
        name        = "PLANTON_RUNNER_TOKEN"
        secret_name = local.token_secret_name
      }
      env {
        name  = "PLANTON_RUNNER_NAME"
        value = local.registration_name
      }
      dynamic "env" {
        for_each = try(var.spec.control_plane_endpoint, "") != "" ? [var.spec.control_plane_endpoint] : []
        content {
          name  = "PLANTON_RUNNER_ENDPOINT"
          value = env.value
        }
      }
      env {
        name  = "PORT"
        value = "50051"
      }
      env {
        name  = "LOG_LEVEL"
        value = "info"
      }

      # The health server answers independently of control-plane
      # reachability, so a runner whose control plane is momentarily
      # unreachable still starts -- its readiness contract is the work
      # queue, not the probe.
      startup_probe {
        transport               = "HTTP"
        port                    = 8093
        path                    = "/healthz"
        initial_delay           = 5
        interval_seconds        = 10
        timeout                 = 5
        failure_count_threshold = 3
      }
    }
  }
}
