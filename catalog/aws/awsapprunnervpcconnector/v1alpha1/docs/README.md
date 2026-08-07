# AwsAppRunnerVpcConnector -- Research Notes

## Provider Surface

Modeled 1:1 on `aws_apprunner_vpc_connector`:

- `vpc_connector_name` (Required, ForceNew, 4-40 chars) -- derived from `metadata.name`, never a spec field.
- `subnets` (Required, ForceNew, set) -- spec `subnet_ids` as `repeated StringValueOrRef` -> `AwsSubnet` (min 1).
- `security_groups` (Required, ForceNew, set) -- spec `security_group_ids` as `repeated StringValueOrRef` -> `AwsSecurityGroup` (min 1). The provider marks this Required (requiredness-honesty: the spec does too).
- Computed: `arn`, `vpc_connector_revision`, `status`.

Every settable attribute is ForceNew -- AWS has no update API for connectors. Recreating under the same name yields a new revision.

## Design Decisions

- **Registry prerequisites `[AwsSubnet, AwsSecurityGroup]`** -- both are hard deploy requirements (required refs), so they belong in the registry rather than scenario annotations; the chains expand transitively (both reach AwsVpc).
- **Egress only, stated loudly** -- the connector is often confused with inbound VPC access; the spec and docs state the ingress plane is a separate resource.
- **`status` exported** -- unlike the versioned configurations, the connector transitions through PENDING_CREATION; the resting state is the useful signal.

## Deferral Ledger

- `aws_apprunner_vpc_ingress_connection` -- DEFER: the INBOUND private-access plane (a PrivateLink ingress into a private service, requiring a VPC endpoint for `apprunner.requests`). A separate product surface that composes against the service's exported `service_arn` with zero rework; revisit on concrete pull.

## Verification

- Spec tests cover required-ness of region/subnets/security groups.
- E2E: the registry prerequisite chains (VPC -> two-AZ subnets; VPC -> SG) deploy first, then create -> DescribeVpcConnector (ACTIVE) -> destroy -> verify INACTIVE (deletion flips status rather than removing the resource).
