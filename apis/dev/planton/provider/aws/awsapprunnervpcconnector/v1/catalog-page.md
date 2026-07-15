# AWS App Runner VPC Connector

The managed network attachment that gives AWS App Runner services outbound access into a VPC -- databases, caches, and internal APIs become reachable without making anything public.

## Why a First-Class Resource

AWS designed the connector to be shared: one connector serves any number of App Runner services, and its managed ENIs persist independently of any single service's lifecycle. Planton models it the same way -- a referenceable node in the resource graph, composed from first-class subnets and security groups, adopted by services through one ARN reference.

## Key Capabilities

- **Private egress** -- route all service outbound traffic into the VPC through AWS-managed ENIs.
- **Multi-AZ reach** -- spread ENIs across availability zones; App Runner routes egress only through AZs the connector covers.
- **Security-group governance** -- the connector's groups define exactly what connected services can reach.
- **Honest immutability** -- every attribute is create-time immutable (AWS has no update API); changes register a new revision under the same name.

## Composes With

- `AwsSubnet` -- the subnets the ENIs land in.
- `AwsSecurityGroup` -- the egress governance.
- `AwsAppRunnerService` -- adopts the connector via `vpcConnectorArn`.
