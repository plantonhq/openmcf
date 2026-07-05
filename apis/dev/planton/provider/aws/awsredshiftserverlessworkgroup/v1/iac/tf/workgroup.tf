# The workgroup is the COMPUTE plane of the serverless warehouse: RPU
# capacity, VPC placement, reachability, and query-level configuration.
# The data it serves lives on the AwsRedshiftServerlessNamespace it
# attaches to by name -- this module never creates or mutates a resource
# that deserves to be its own node (subnets and security groups attach
# by reference; ingress rules live on the referenced AwsSecurityGroup
# nodes).
#
# Create-only in AWS: the workgroup name and the namespace it serves.
# Capacity edits are in place, though AWS serializes them (the provider
# issues one capacity change per update call for exactly that reason).
resource "aws_redshiftserverless_workgroup" "this" {
  workgroup_name = local.workgroup_name
  namespace_name = var.spec.namespace_name

  # The capacity contract (CEL enforces the shape): a fixed RPU baseline
  # OR an enabled price_performance_target where AWS owns the baseline.
  # 0 keeps the AWS default (128 RPU); AWS validates the exact RPU
  # increments at deploy, since they have changed over time. A max of 0
  # leaves scaling uncapped (the AWS default).
  base_capacity = var.spec.base_capacity != 0 ? var.spec.base_capacity : null
  max_capacity  = var.spec.max_capacity != 0 ? var.spec.max_capacity : null

  # The price-performance dial is rendered even when disabled IF the
  # message is present, so a later enable/disable edits in place; an
  # absent message keeps the block off entirely. level 0 keeps the AWS
  # default (50, balanced).
  dynamic "price_performance_target" {
    for_each = var.spec.price_performance_target != null ? [var.spec.price_performance_target] : []
    content {
      enabled = price_performance_target.value.enabled
      level   = price_performance_target.value.level != 0 ? price_performance_target.value.level : null
    }
  }

  # VPC placement: at least three subnets in three AZs (CEL-enforced
  # when set); empty falls back to the account's default VPC. The VPC
  # default security group applies when none are given (AWS's own
  # default).
  subnet_ids         = length(var.spec.subnet_ids) > 0 ? var.spec.subnet_ids : null
  security_group_ids = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null

  # Enhanced VPC routing forces COPY/UNLOAD data movement through the
  # VPC where flow logs and endpoints can see and govern it.
  # Reachability is opt-in; a private workgroup is the norm.
  enhanced_vpc_routing = var.spec.enhanced_vpc_routing
  publicly_accessible  = var.spec.publicly_accessible

  # 0 keeps the AWS default (5439); CEL constrains set values to the
  # only ranges Redshift Serverless accepts (5431-5455, 8191-8215).
  port = var.spec.port != 0 ? var.spec.port : null

  # Query-level parameters apply directly to the workgroup --
  # serverless has no parameter groups, so there is nothing to fold or
  # reference.
  dynamic "config_parameter" {
    for_each = var.spec.config_parameters
    content {
      parameter_key   = config_parameter.value.name
      parameter_value = config_parameter.value.value
    }
  }

  # Empty keeps the AWS default release track ("current").
  track_name = var.spec.track_name != "" ? var.spec.track_name : null

  tags = local.aws_tags
}
