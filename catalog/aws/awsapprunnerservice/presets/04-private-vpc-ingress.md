# Private Internal Service (VPC Ingress)

This preset creates an App Runner service with NO public endpoint: incoming traffic is disabled on the internet side, and the service is published into a VPC through an interface VPC endpoint (AWS PrivateLink). Clients inside that VPC reach the service at a private domain name; nothing outside the VPC can reach it at all.

## When to Use

- Internal APIs and admin services that must never be internet-reachable
- Microservices called only by other workloads inside your VPC (ECS, EKS, EC2, Lambda in-VPC)
- Compliance postures that forbid public endpoints for a class of services

## Key Configuration Choices

- **Private incoming traffic** (`isPubliclyAccessible: false`) -- The service gets no public URL. Pair this with at least one VPC Ingress Connection or nothing can reach the service.
- **VPC Ingress Connection** (`vpcIngressConnections`) -- Publishes the service into the named VPC through the referenced interface endpoint. The endpoint must be an **interface VPC endpoint for `com.amazonaws.REGION.apprunner.requests`** in the target VPC (model it with `AwsVpcEndpoint`). The per-connection private domain name is exported in `status.outputs.vpc_ingress_connections` -- point internal DNS or service discovery at it.
- **Connection entries are create-time immutable** -- changing an entry's VPC or endpoint replaces that one connection (keyed by name; adding/removing entries updates in place).
- **Private ECR image** with an `accessRoleArn` pull role, same as the production preset.

## How Clients Connect

1. Deploy this service; read `status.outputs.vpc_ingress_connections[0].domain_name`.
2. Clients inside the connected VPC resolve that domain (App Runner manages the private hosted zone wiring through the endpoint) and call it over HTTPS as usual.
3. Egress is unchanged by this preset -- add a `vpcConnectorArn` if the service must also REACH private resources (ingress and egress are independent).

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `imageIdentifier` | Your private ECR image path (ACCOUNT_ID.dkr.ecr.REGION.amazonaws.com/REPO:TAG) | AWS ECR Console |
| `ecr-access-role` | Name of the `AwsIamRole` granting ECR pull access | Your resource graph |
| `main-vpc` | Name of the `AwsVpc` to publish the service into | Your resource graph |
| `apprunner-requests-endpoint` | Name of the `AwsVpcEndpoint` (interface endpoint for the App Runner requests service) in that VPC | Your resource graph |

## Related Presets

- **01-basic-public-image** -- The opposite posture: a public prototype endpoint.
- **02-production-vpc-encrypted** -- Public endpoint with private EGRESS (VPC connector); combine its egress story with this preset's ingress story for a fully VPC-internal service.
