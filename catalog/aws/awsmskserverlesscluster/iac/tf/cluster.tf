# A serverless MSK cluster is a single, essentially immutable resource: AWS
# manages brokers, storage, and Kafka version, so the whole declaration is
# WHERE it lives (one or more VPC placements). Everything except tags is
# create-time (ForceNew) -- changing networking or auth replaces the cluster.
resource "aws_msk_serverless_cluster" "this" {
  cluster_name = local.cluster_name

  # One block per VPC placement: AWS provisions client-facing network
  # interfaces in EACH declared VPC, so applications in every listed VPC
  # connect privately without peering or PrivateLink setup.
  dynamic "vpc_config" {
    for_each = var.spec.vpc_configs
    content {
      subnet_ids = vpc_config.value.subnet_ids
      # Optional: AWS attaches the VPC's default security group when omitted.
      # The ingress rule for the SASL/IAM listener port (9098) lives on the
      # referenced first-class security group nodes.
      security_group_ids = length(vpc_config.value.security_group_ids) > 0 ? vpc_config.value.security_group_ids : null
    }
  }

  # SASL/IAM is the ONLY client-authentication scheme serverless MSK supports,
  # and it is mandatory -- so it is hardcoded rather than exposed as a spec
  # field that could only ever hold one value. Clients authenticate with AWS
  # IAM credentials on port 9098.
  client_authentication {
    sasl {
      iam {
        enabled = true
      }
    }
  }

  tags = local.aws_tags
}
