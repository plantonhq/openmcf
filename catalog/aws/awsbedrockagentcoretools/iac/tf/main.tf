# Amazon Bedrock AgentCore built-in tools: managed cloud browsers (with
# optional session recording, Chrome enterprise policies, and mTLS
# certificates), reusable browser profiles, and code-execution sandboxes.
#
# Lifecycle facts the renders below depend on:
#   - AWS exposes NO update for any of these resources -- every field
#     change recreates the tool. That is cheap: the tools are
#     configuration shells; AWS bills only per session at runtime;
#   - the browser and code interpreter share the certificate shape (a
#     Secrets Manager location); the enterprise policy and recording
#     shapes are browser-only.

resource "aws_bedrockagentcore_browser" "this" {
  for_each = local.browsers

  name = each.value.name

  description        = each.value.description != "" ? each.value.description : null
  execution_role_arn = each.value.execution_role_arn != "" ? each.value.execution_role_arn : null

  # Required by AWS. VPC mode carries the placement (spec-validated
  # pairing).
  network_configuration {
    network_mode = each.value.network.mode
    dynamic "vpc_config" {
      for_each = each.value.network.vpc_config != null ? [each.value.network.vpc_config] : []
      content {
        subnets         = vpc_config.value.subnets
        security_groups = vpc_config.value.security_groups
      }
    }
  }

  # Cryptographic traffic signing -- rendered only on an explicit choice
  # so the module never fights AWS's default.
  dynamic "browser_signing" {
    for_each = each.value.signing_enabled != null ? [each.value.signing_enabled] : []
    content {
      enabled = browser_signing.value
    }
  }

  # Session recording to S3.
  dynamic "recording" {
    for_each = each.value.recording != null ? [each.value.recording] : []
    content {
      enabled = recording.value.enabled
      dynamic "s3_location" {
        for_each = recording.value.s3_location != null ? [recording.value.s3_location] : []
        content {
          bucket = s3_location.value.bucket
          prefix = s3_location.value.prefix
        }
      }
    }
  }

  # Chrome enterprise policy files (max 100).
  dynamic "enterprise_policy" {
    for_each = each.value.enterprise_policies
    content {
      type = enterprise_policy.value.type != "" ? enterprise_policy.value.type : null
      location {
        s3 {
          bucket     = enterprise_policy.value.s3.bucket
          prefix     = enterprise_policy.value.s3.prefix
          version_id = enterprise_policy.value.s3.version_id != "" ? enterprise_policy.value.s3.version_id : null
        }
      }
    }
  }

  # Client certificates for mTLS-protected sites (max 200).
  dynamic "certificate" {
    for_each = each.value.certificates
    content {
      location {
        secrets_manager {
          secret_arn = certificate.value.secret_arn
        }
      }
    }
  }

  tags = local.aws_tags
}

# Reusable saved browser state (cookies, logins).
resource "aws_bedrockagentcore_browser_profile" "this" {
  for_each = local.browser_profiles

  name = each.value.name

  description = each.value.description != "" ? each.value.description : null

  tags = local.aws_tags
}

# Managed code-execution sandboxes.
resource "aws_bedrockagentcore_code_interpreter" "this" {
  for_each = local.code_interpreters

  name = each.value.name

  description        = each.value.description != "" ? each.value.description : null
  execution_role_arn = each.value.execution_role_arn != "" ? each.value.execution_role_arn : null

  # Required by AWS. SANDBOX blocks all network access (the safest for
  # untrusted code); VPC mode carries the placement (spec-validated
  # pairing).
  network_configuration {
    network_mode = each.value.network.mode
    dynamic "vpc_config" {
      for_each = each.value.network.vpc_config != null ? [each.value.network.vpc_config] : []
      content {
        subnets         = vpc_config.value.subnets
        security_groups = vpc_config.value.security_groups
      }
    }
  }

  # Client certificates for mTLS-protected endpoints the code calls
  # (max 200).
  dynamic "certificate" {
    for_each = each.value.certificates
    content {
      location {
        secrets_manager {
          secret_arn = certificate.value.secret_arn
        }
      }
    }
  }

  tags = local.aws_tags
}
