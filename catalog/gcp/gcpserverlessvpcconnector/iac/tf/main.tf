# Enable the Serverless VPC Access API before creating the connector so a
# fresh project works first try. disable_on_destroy=false: turning an API
# off on teardown is a project-wide blast radius no single resource should
# own.
resource "google_project_service" "vpcaccess_api" {
  project = local.project_id
  service = "vpcaccess.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Serverless VPC Access connector: a managed fleet of forwarding
# instances inside the VPC that serverless workloads (Cloud Functions,
# Cloud Run, App Engine) route egress through to reach private IPs.
#
# Placement is an exactly-one contract enforced pre-deploy by the spec's
# CEL rules: either network + ip_cidr_range (the connector carves a new /28
# out of the network) or an existing dedicated /28 subnet (the
# Shared-VPC-capable mode).
resource "google_vpc_access_connector" "main" {
  name    = local.connector_name
  project = local.project_id
  region  = var.spec.region

  network       = local.network
  ip_cidr_range = local.ip_cidr_range

  dynamic "subnet" {
    for_each = var.spec.subnet != null ? [var.spec.subnet] : []
    content {
      name       = subnet.value.name
      project_id = subnet.value.project_id != "" ? subnet.value.project_id : null
    }
  }

  machine_type = local.machine_type

  # Scaling: only the instance-based contract is modeled. The legacy
  # min/max_throughput fields are deliberately not set — the provider
  # discourages them in favor of instances, they conflict with the
  # instance fields, and they force replacement on change.
  #
  # Shrink asymmetry (worth knowing before an update): the provider applies
  # INCREASES to min/max_instances in place but forces the connector to be
  # REPLACED when either value is DECREASED — a brief egress outage for
  # every workload using the connector.
  min_instances = var.spec.min_instances
  max_instances = var.spec.max_instances

  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.vpcaccess_api]
}
