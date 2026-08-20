# AwsRestApiVpcLink — Pulumi module (Go)

Deploys an API Gateway v1 VPC link (`apigateway.VpcLink`) fronting
exactly one Network Load Balancer.

Module facts worth knowing before editing:

- **`targetArns` is a one-element array on the bridged type.** The
  spec's single `target_arn` renders as that array.
- **The balancer is create-time immutable.** Changing it replaces the
  link; integrations must re-home onto the new ID.
- **Provisioning is slow** (several minutes) while AWS builds the
  network attachment. Create is free; NLB charges apply to the
  balancer.

Outputs mirror the Terraform module key-for-key: `vpc_link_id`,
`vpc_link_arn`.
