# Amazon Bedrock AgentCore agent runtime: a serverless, session-isolated
# execution environment hosting YOUR agent code (container image or
# AWS-managed code bundle) with folded satellites -- named serving
# endpoints and the runtime's resource-based policy.
#
# Lifecycle facts the renders below depend on:
#   - every configuration change creates a new runtime VERSION in place;
#     switching the artifact arm (code <-> container) replaces the runtime
#     (provider-enforced);
#   - endpoints pin or float across versions -- an endpoint without an
#     explicit version tracks the latest;
#   - the resource policy attaches to the runtime's own ARN (the provider
#     resource accepts any AgentCore ARN; this module scopes it to the
#     runtime it deploys).

# The shared JWT-authorizer render: both the private_endpoint and each
# private_endpoint_overrides entry carry the same two-arm endpoint shape,
# rendered by the nested dynamic blocks below.
resource "aws_bedrockagentcore_agent_runtime" "this" {
  # AWS's runtime-name charset (letter first, then letters/digits/_) is
  # stricter than metadata.name conventions, so the name is an explicit
  # spec field. Changing it replaces the runtime.
  agent_runtime_name = var.spec.runtime_name

  # Required by AWS: the role the AgentCore service assumes to pull the
  # image / read the code bundle and run the agent.
  role_arn = var.spec.role_arn

  description = var.spec.description != "" ? var.spec.description : null

  environment_variables = length(var.spec.environment_variables) > 0 ? var.spec.environment_variables : null

  # Optional+Computed single-entry list ATTRIBUTE at the provider (not a
  # block): sent only when set so the module never fights AWS's defaults.
  lifecycle_configuration = local.has_lifecycle ? [{
    idle_runtime_session_timeout = var.spec.lifecycle.idle_runtime_session_timeout_seconds != 0 ? var.spec.lifecycle.idle_runtime_session_timeout_seconds : null
    max_lifetime                 = var.spec.lifecycle.max_lifetime_seconds != 0 ? var.spec.lifecycle.max_lifetime_seconds : null
  }] : null

  # Exactly one artifact arm (spec-validated). Switching arms replaces
  # the runtime.
  agent_runtime_artifact {
    dynamic "container_configuration" {
      for_each = var.spec.artifact.container != null ? [var.spec.artifact.container] : []
      content {
        container_uri = container_configuration.value.image_uri
      }
    }
    dynamic "code_configuration" {
      for_each = local.has_code_artifact ? [var.spec.artifact.code] : []
      content {
        entry_point = code_configuration.value.entry_point
        runtime     = code_configuration.value.runtime
        code {
          s3 {
            bucket     = code_configuration.value.s3.bucket
            prefix     = code_configuration.value.s3.prefix
            version_id = code_configuration.value.s3.version_id != "" ? code_configuration.value.s3.version_id : null
          }
        }
      }
    }
  }

  # Required by AWS. VPC mode carries the placement (spec-validated
  # pairing).
  network_configuration {
    network_mode = var.spec.network.mode
    dynamic "network_mode_config" {
      for_each = var.spec.network.vpc_config != null ? [var.spec.network.vpc_config] : []
      content {
        subnets         = network_mode_config.value.subnets
        security_groups = network_mode_config.value.security_groups
      }
    }
  }

  # HTTP is AWS's default protocol; the block renders only on an explicit
  # choice so the module never fights the default.
  dynamic "protocol_configuration" {
    for_each = var.spec.server_protocol != "" ? [var.spec.server_protocol] : []
    content {
      server_protocol = protocol_configuration.value
    }
  }

  dynamic "request_header_configuration" {
    for_each = length(var.spec.request_header_allowlist) > 0 ? [var.spec.request_header_allowlist] : []
    content {
      request_header_allowlist = request_header_configuration.value
    }
  }

  # Inbound JWT authorization. Omitted = AWS IAM (SigV4) callers only.
  dynamic "authorizer_configuration" {
    for_each = local.has_jwt ? [var.spec.custom_jwt_authorizer] : []
    content {
      custom_jwt_authorizer {
        discovery_url    = authorizer_configuration.value.discovery_url
        allowed_audience = length(authorizer_configuration.value.allowed_audience) > 0 ? authorizer_configuration.value.allowed_audience : null
        allowed_clients  = length(authorizer_configuration.value.allowed_clients) > 0 ? authorizer_configuration.value.allowed_clients : null
        allowed_scopes   = length(authorizer_configuration.value.allowed_scopes) > 0 ? authorizer_configuration.value.allowed_scopes : null

        dynamic "allowed_workload_configuration" {
          for_each = authorizer_configuration.value.allowed_workloads != null ? [authorizer_configuration.value.allowed_workloads] : []
          content {
            workload_identities = length(allowed_workload_configuration.value.workload_identities) > 0 ? allowed_workload_configuration.value.workload_identities : null
            dynamic "hosting_environment" {
              for_each = allowed_workload_configuration.value.hosting_environment_arns
              content {
                arn = hosting_environment.value
              }
            }
          }
        }

        dynamic "custom_claim" {
          for_each = authorizer_configuration.value.custom_claims
          content {
            inbound_token_claim_name       = custom_claim.value.claim_name
            inbound_token_claim_value_type = custom_claim.value.value_type
            authorizing_claim_match_value {
              claim_match_operator = custom_claim.value.match_operator
              claim_match_value {
                match_value_string      = custom_claim.value.match_value != "" ? custom_claim.value.match_value : null
                match_value_string_list = length(custom_claim.value.match_values) > 0 ? custom_claim.value.match_values : null
              }
            }
          }
        }

        dynamic "private_endpoint" {
          for_each = authorizer_configuration.value.private_endpoint != null ? [authorizer_configuration.value.private_endpoint] : []
          content {
            dynamic "managed_vpc_resource" {
              for_each = private_endpoint.value.managed_vpc != null ? [private_endpoint.value.managed_vpc] : []
              content {
                vpc_identifier           = managed_vpc_resource.value.vpc_id
                subnet_ids               = managed_vpc_resource.value.subnet_ids
                security_group_ids       = length(managed_vpc_resource.value.security_group_ids) > 0 ? managed_vpc_resource.value.security_group_ids : null
                endpoint_ip_address_type = managed_vpc_resource.value.endpoint_ip_address_type
                routing_domain           = managed_vpc_resource.value.routing_domain != "" ? managed_vpc_resource.value.routing_domain : null
                tags                     = length(managed_vpc_resource.value.tags) > 0 ? managed_vpc_resource.value.tags : null
              }
            }
            dynamic "self_managed_lattice_resource" {
              for_each = private_endpoint.value.self_managed_lattice != null ? [private_endpoint.value.self_managed_lattice] : []
              content {
                resource_configuration_identifier = self_managed_lattice_resource.value.resource_configuration_id
              }
            }
          }
        }

        dynamic "private_endpoint_overrides" {
          for_each = authorizer_configuration.value.private_endpoint_overrides
          content {
            domain = private_endpoint_overrides.value.domain
            private_endpoint {
              dynamic "managed_vpc_resource" {
                for_each = private_endpoint_overrides.value.private_endpoint.managed_vpc != null ? [private_endpoint_overrides.value.private_endpoint.managed_vpc] : []
                content {
                  vpc_identifier           = managed_vpc_resource.value.vpc_id
                  subnet_ids               = managed_vpc_resource.value.subnet_ids
                  security_group_ids       = length(managed_vpc_resource.value.security_group_ids) > 0 ? managed_vpc_resource.value.security_group_ids : null
                  endpoint_ip_address_type = managed_vpc_resource.value.endpoint_ip_address_type
                  routing_domain           = managed_vpc_resource.value.routing_domain != "" ? managed_vpc_resource.value.routing_domain : null
                  tags                     = length(managed_vpc_resource.value.tags) > 0 ? managed_vpc_resource.value.tags : null
                }
              }
              dynamic "self_managed_lattice_resource" {
                for_each = private_endpoint_overrides.value.private_endpoint.self_managed_lattice != null ? [private_endpoint_overrides.value.private_endpoint.self_managed_lattice] : []
                content {
                  resource_configuration_identifier = self_managed_lattice_resource.value.resource_configuration_id
                }
              }
            }
          }
        }
      }
    }
  }

  # Filesystem mounts (max 5; exactly one source arm each,
  # spec-validated). session_storage's only argument is its mount path.
  dynamic "filesystem_configuration" {
    for_each = var.spec.filesystems
    content {
      dynamic "efs_access_point" {
        for_each = filesystem_configuration.value.efs_access_point_arn != "" ? [filesystem_configuration.value] : []
        content {
          access_point_arn = efs_access_point.value.efs_access_point_arn
          mount_path       = efs_access_point.value.mount_path
        }
      }
      dynamic "s3_files_access_point" {
        for_each = filesystem_configuration.value.s3_files_access_point_arn != "" ? [filesystem_configuration.value] : []
        content {
          access_point_arn = s3_files_access_point.value.s3_files_access_point_arn
          mount_path       = s3_files_access_point.value.mount_path
        }
      }
      dynamic "session_storage" {
        for_each = filesystem_configuration.value.session_storage ? [filesystem_configuration.value] : []
        content {
          mount_path = session_storage.value.mount_path
        }
      }
    }
  }

  tags = local.aws_tags
}

# Named serving endpoints. An endpoint without an explicit version tracks
# the runtime's latest version; a pinned endpoint serves that version
# until re-pointed.
resource "aws_bedrockagentcore_agent_runtime_endpoint" "this" {
  for_each = local.endpoints

  agent_runtime_id = aws_bedrockagentcore_agent_runtime.this.agent_runtime_id
  name             = each.value.name

  description           = each.value.description != "" ? each.value.description : null
  agent_runtime_version = each.value.agent_runtime_version != "" ? each.value.agent_runtime_version : null

  tags = local.aws_tags
}

# Resource-based policy on the runtime's own ARN -- grant other accounts
# or services permission to invoke this runtime. AWS requires every
# statement's Resource to be exactly the attached runtime's ARN
# (PutResourcePolicy 400 "must contain exactly one resource ARN that
# matches the provided resource ARN", live-caught 2026-08-14) -- an ARN
# no author can know before create, so the module owns the Resource
# member on every statement (standard IAM JSON casing assumed).
resource "aws_bedrockagentcore_resource_policy" "this" {
  count = var.spec.resource_policy != null ? 1 : 0

  resource_arn = aws_bedrockagentcore_agent_runtime.this.agent_runtime_arn
  policy = jsonencode(merge(var.spec.resource_policy, {
    Statement = [
      for s in try(var.spec.resource_policy.Statement, []) :
      merge(s, { Resource = aws_bedrockagentcore_agent_runtime.this.agent_runtime_arn })
    ]
  }))
}
