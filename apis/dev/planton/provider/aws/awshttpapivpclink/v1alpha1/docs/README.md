# AWS HTTP API VPC Link: Architecture and Design

## What a VPC Link Is

An API Gateway v2 VPC link is AWS's managed answer to "how does a public API reach a private backend?" When created, AWS provisions cross-account elastic network interfaces into the subnets you choose. HTTP API integrations with `connection_type: VPC_LINK` route requests through those ENIs to targets inside the VPC -- an internal Application Load Balancer listener, a Network Load Balancer listener, or an AWS Cloud Map service -- so the backend never needs an internet-facing endpoint.

## Design Decisions

### 1. First-Class Resource, Not a Field on the API

One VPC link serves any number of APIs and integrations. Modeling it inside the API spec would force every API to own (and churn) its own link, multiplying ENIs and losing the shared-infrastructure semantics AWS designed. As a first-class resource, the link is created once per VPC and referenced by ID from every private integration -- the same composition pattern as target groups and listeners in the load-balancing family.

### 2. Honest Immutability

AWS exposes no update API for the link's subnets or security groups -- both are create-time attributes, and the provider marks them ForceNew. The spec documents this directly on the fields: changing either replaces the link (new ENIs, new link ID), and every integration referencing the old ID follows the replacement through the reference graph. Only the name mutates in place.

### 3. Security Groups Optional

The AWS API requires subnets but not security groups. A link created without security groups applies no filtering on its side -- reachability is then governed solely by the target's security groups (which must admit ingress from the link's ENI addresses or their group). The spec models this honestly: `subnet_ids` requires at least one entry, `security_group_ids` may be empty. The recommended posture is still an explicit egress-scoped group, because it gives the target's security group a stable source to reference.

### 4. Multi-AZ Guidance, Not Enforcement

A VPC link can only reach targets in availability zones where it has an ENI. AWS accepts a single-subnet link, so the spec does too -- but the field comment and presets steer toward two or more AZs, matching the load balancer family's posture.

## Lifecycle Notes

- **Creation** takes a few minutes: AWS waits for the link to reach `AVAILABLE` before integrations can use it.
- **Deletion** waits for the ENIs to be reclaimed; a link still referenced by an integration cannot be deleted.
- **Status values**: `PENDING` → `AVAILABLE`; `DELETING` on the way out; `FAILED` if provisioning fails (usually IAM or subnet capacity).

## Relationship to the REST API VPC Link

REST APIs (API Gateway v1) have their own `aws_api_gateway_vpc_link` resource that fronts NLBs only and uses a different API. The two are not interchangeable; this component models the v2 link used by HTTP APIs. A REST API surface, if built, brings its own link kind.

## Dependencies

### Upstream

| Resource | Field | Relationship |
|----------|-------|-------------|
| AwsSubnet | `subnet_ids` | Required -- hosts the link's ENIs |
| AwsSecurityGroup | `security_group_ids` | Optional -- filters the link's reach |

### Downstream

| Resource | Output Used | Use Case |
|----------|-------------|----------|
| AwsHttpApiGateway | `vpc_link_id` | Private integrations set it as `connection_id` |
