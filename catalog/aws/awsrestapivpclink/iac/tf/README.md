# AwsRestApiVpcLink — Terraform/OpenTofu module

Deploys an API Gateway v1 VPC link (`aws_api_gateway_vpc_link`)
fronting exactly one Network Load Balancer.

Module facts worth knowing before editing:

- **`target_arns` is a one-element list upstream.** The spec's single
  `target_arn` renders as that list.
- **The balancer is create-time immutable.** Changing it replaces the
  link; integrations must re-home onto the new ID.
- **Provisioning is slow** (several minutes) while AWS builds the
  network attachment. Create is free; NLB charges apply to the
  balancer.

Outputs mirror the Pulumi module key-for-key: `vpc_link_id`,
`vpc_link_arn`.
