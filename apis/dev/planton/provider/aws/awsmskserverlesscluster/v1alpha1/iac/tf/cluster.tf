# A serverless MSK cluster is a single, essentially immutable resource: AWS
# manages brokers, storage, and Kafka version, so the whole declaration is
# WHERE it lives (subnets + security groups). Everything except tags is
# create-time (ForceNew) -- changing networking or auth replaces the cluster.
resource "aws_msk_serverless_cluster" "this" {
  cluster_name = local.cluster_name

  vpc_config {
    subnet_ids = var.spec.subnet_ids
    # Optional: AWS attaches the VPC's default security group when omitted.
    # The ingress rule for the SASL/IAM listener port (9098) lives on the
    # referenced first-class security group nodes.
    security_group_ids = length(var.spec.security_group_ids) > 0 ? var.spec.security_group_ids : null
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
