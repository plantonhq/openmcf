# Mesh Spoke

The default spoke shape: two-AZ subnet coverage on a full-mesh hub. Default route table membership is inherited from the gateway, so the spoke lands in the mesh with zero routing configuration.

## When to Use

- Joining a VPC to a full-mesh hub where every attachment should reach every other

## What It Configures

- **Two-AZ subnet spread** — the gateway only routes traffic to/from AZs it has an ENI in; cover every AZ your workloads run in
- **Inherited default-table membership** — association and propagation follow the gateway's dials

## What to Customize

- Replace `<aws-region>` and the referenced hub/VPC/subnet names with your resources
- Dedicated /28 subnets per AZ for attachments are the AWS-recommended pattern — they keep routing independent of workload subnets
- Remember the return path: spoke subnets need a route with `targetType: transit_gateway` toward the gateway
