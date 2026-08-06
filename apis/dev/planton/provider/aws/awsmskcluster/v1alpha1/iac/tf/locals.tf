locals {
  # The cluster name is metadata.name -- create-only in AWS (max 64 chars), and
  # the basis both engines share so a manifest deploys identically on either.
  cluster_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsMskCluster"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # An MSK Configuration is managed here only for inline server_properties;
  # bringing an existing configuration (configuration_arn + revision) and inline
  # properties are mutually exclusive (CEL-enforced). The configuration itself is
  # pure glue -- a named server.properties document consumed by exactly one
  # cluster -- which is why it stays inside this module instead of being its own
  # node.
  manage_configuration = length(var.spec.server_properties) > 0

  # Kafka's server.properties is a flat "key = value" document; the map keeps the
  # manifest structured while this serialization feeds the AWS API.
  server_properties = local.manage_configuration ? join("\n", [
    for k, v in var.spec.server_properties : "${k} = ${v}"
  ]) : ""

  # connectivity_info is emitted only when some connectivity surface is
  # configured; an empty block would still trigger AWS's create-then-update
  # connectivity flow for no reason.
  vpc_connectivity_enabled = var.spec.vpc_connectivity != null ? (
    var.spec.vpc_connectivity.sasl_iam_enabled ||
    var.spec.vpc_connectivity.sasl_scram_enabled ||
    var.spec.vpc_connectivity.tls_enabled
  ) : false
  manage_connectivity_info = var.spec.public_access_type != "" || local.vpc_connectivity_enabled || var.spec.network_type != ""
}
