# One Cloud Map namespace (HTTP XOR private DNS XOR public DNS) with
# its services and statically registered instances managed in-line.
#
# Lifecycle facts the render below depends on:
#   - the three namespace types are three provider resources; exactly
#     one exists per the spec's type field, all exposing the same
#     downstream surface (id/arn);
#   - the HTTP namespace has NO update path at the provider - changing
#     its description replaces it; the DNS namespaces update
#     description in place;
#   - the private namespace's vpc is never read back by the provider
#     (imports carry it as "{namespace_id}:{vpc_id}");
#   - a service binds its namespace through dns_config.namespace_id
#     when it publishes DNS records, and through the top-level
#     namespace_id otherwise (the provider's own documented split for
#     its legacy duplicated pointer);
#   - instance registration is an AWS upsert (create and update are
#     the same RegisterInstance call, keyed by instance_id); the
#     provider derives AWS_INSTANCE_IPV4 from AWS_EC2_INSTANCE_ID and
#     drops the derived echo on read;
#   - deregistering an already-gone instance errors at the provider
#     (no NotFound tolerance) - destroy instances through this module,
#     never out-of-band;
#   - a service's force_destroy deregisters EVERY instance in the
#     service first, including runtime-registered ones this manifest
#     never declared.

resource "aws_service_discovery_http_namespace" "this" {
  count = var.spec.type == "HTTP" ? 1 : 0

  name        = var.metadata.name
  description = var.spec.description != "" ? var.spec.description : null

  tags = local.aws_tags
}

resource "aws_service_discovery_private_dns_namespace" "this" {
  count = var.spec.type == "PRIVATE_DNS" ? 1 : 0

  name        = var.metadata.name
  vpc         = var.spec.vpc_id
  description = var.spec.description != "" ? var.spec.description : null

  tags = local.aws_tags
}

resource "aws_service_discovery_public_dns_namespace" "this" {
  count = var.spec.type == "PUBLIC_DNS" ? 1 : 0

  name        = var.metadata.name
  description = var.spec.description != "" ? var.spec.description : null

  tags = local.aws_tags
}

locals {
  # The one namespace whichever arm produced - the downstream surface
  # is identical.
  namespace_id = var.spec.type == "HTTP" ? aws_service_discovery_http_namespace.this[0].id : (var.spec.type == "PRIVATE_DNS" ? aws_service_discovery_private_dns_namespace.this[0].id : aws_service_discovery_public_dns_namespace.this[0].id)

  namespace_arn = var.spec.type == "HTTP" ? aws_service_discovery_http_namespace.this[0].arn : (var.spec.type == "PRIVATE_DNS" ? aws_service_discovery_private_dns_namespace.this[0].arn : aws_service_discovery_public_dns_namespace.this[0].arn)

  hosted_zone_id = var.spec.type == "PRIVATE_DNS" ? aws_service_discovery_private_dns_namespace.this[0].hosted_zone : (var.spec.type == "PUBLIC_DNS" ? aws_service_discovery_public_dns_namespace.this[0].hosted_zone : "")

  http_name = var.spec.type == "HTTP" ? aws_service_discovery_http_namespace.this[0].http_name : ""
}

# Services, keyed by service name.
resource "aws_service_discovery_service" "this" {
  for_each = { for service in var.spec.services : service.name => service }

  name        = each.value.name
  description = each.value.description != "" ? each.value.description : null

  # DNS-publishing services bind the namespace inside dns_config;
  # API-only services bind it at the top level.
  namespace_id = each.value.dns_config == null ? local.namespace_id : null

  dynamic "dns_config" {
    for_each = each.value.dns_config != null ? [each.value.dns_config] : []
    content {
      namespace_id   = local.namespace_id
      routing_policy = dns_config.value.routing_policy != "" ? dns_config.value.routing_policy : null

      dynamic "dns_records" {
        for_each = dns_config.value.records
        content {
          type = dns_records.value.type
          ttl  = dns_records.value.ttl
        }
      }
    }
  }

  dynamic "health_check_config" {
    for_each = each.value.health_check_config != null ? [each.value.health_check_config] : []
    content {
      type              = health_check_config.value.type != "" ? health_check_config.value.type : null
      resource_path     = health_check_config.value.resource_path != "" ? health_check_config.value.resource_path : null
      failure_threshold = health_check_config.value.failure_threshold > 0 ? health_check_config.value.failure_threshold : null
    }
  }

  dynamic "health_check_custom_config" {
    for_each = each.value.health_check_custom_config != null ? [each.value.health_check_custom_config] : []
    content {}
  }

  force_destroy = each.value.force_destroy ? true : null

  tags = local.aws_tags
}

locals {
  # One registration per (service, instance), keyed
  # "service//instance_id" - the "//" separator is the import bridge's
  # address-key-segment convention.
  service_instances = flatten([
    for service in var.spec.services : [
      for instance in service.instances : {
        key          = "${service.name}//${instance.instance_id}"
        service_name = service.name
        instance     = instance
      }
    ]
  ])

  service_instances_by_key = { for entry in local.service_instances : entry.key => entry }
}

resource "aws_service_discovery_instance" "this" {
  for_each = local.service_instances_by_key

  instance_id = each.value.instance.instance_id
  service_id  = aws_service_discovery_service.this[each.value.service_name].id

  # Typed fields compose AWS's documented attribute keys; anything
  # else rides custom_attributes verbatim.
  attributes = merge(
    each.value.instance.ip != "" ? { AWS_INSTANCE_IPV4 = each.value.instance.ip } : {},
    each.value.instance.ipv6 != "" ? { AWS_INSTANCE_IPV6 = each.value.instance.ipv6 } : {},
    each.value.instance.port > 0 ? { AWS_INSTANCE_PORT = tostring(each.value.instance.port) } : {},
    each.value.instance.cname != "" ? { AWS_INSTANCE_CNAME = each.value.instance.cname } : {},
    each.value.instance.alias_dns_name != "" ? { AWS_ALIAS_DNS_NAME = each.value.instance.alias_dns_name } : {},
    each.value.instance.ec2_instance_id != "" ? { AWS_EC2_INSTANCE_ID = each.value.instance.ec2_instance_id } : {},
    each.value.instance.custom_attributes
  )
}
