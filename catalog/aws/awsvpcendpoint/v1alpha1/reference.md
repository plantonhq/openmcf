# AwsVpcEndpoint

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsVpcEndpointSpec creates a private connection from a VPC to an AWS
service, a PrivateLink service, or a VPC Lattice resource -- traffic
stays on the AWS network instead of crossing the internet through a
NAT or internet gateway.

Two endpoint types carry nearly all real-world use:

  - Gateway (S3 and DynamoDB only): free; works by injecting a
    prefix-list route into the route tables you attach. The default
    private path for S3/DynamoDB traffic -- it also removes that
    traffic from NAT gateway data-processing charges.
  - Interface (everything else, plus S3/DynamoDB when on-premises or
    cross-VPC access is needed): billed per AZ-hour + per GB; places
    an ENI in each subnet you attach and (optionally) overrides the
    service's public DNS name inside the VPC via private DNS.

The remaining types are specialized: GatewayLoadBalancer front a
GWLB appliance fleet, Resource and ServiceNetwork attach to VPC
Lattice. All five are modeled; the type gates which attachment
fields apply, and CEL enforces the gating at validate time instead
of letting AWS reject the create.

The endpoint composes onto its neighbors instead of embedding them:
the VPC attaches by reference, gateway endpoints reference route
tables (an AwsSubnet's route_table_id output when the subnet owns its
table, or the AwsVpc's main/default route-table outputs), and
interface endpoints reference AwsSubnet and AwsSecurityGroup nodes.
This component never modifies a resource it merely references.

Create-time immutable in AWS: the endpoint type, the service target
(service_name / resource_configuration_arn / service_network_arn),
and service_region. Route tables, subnets, security groups, policy,
DNS options, and IP address type all update in place.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsVpcEndpoint
metadata:
  name: awsvpcendpoint-demo
spec:
  region: us-west-2
  vpcId:
    value: vpc-0a1b2c3d4e5f67890
  endpointType: Interface
  serviceName: com.amazonaws.us-west-2.ecr.dkr
  subnetIds:
    - value: subnet-0a1b2c3d4e5f67890
    - value: subnet-0f9e8d7c6b5a43210
  securityGroupIds:
    - value: sg-0a1b2c3d4e5f67890
  privateDnsEnabled: true
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.vpcId` | `string \| valueFrom` | yes |  | AwsVpc (`status.outputs.vpc_id`) |
| `spec.endpointType` | `string` |  |  |  |
| `spec.serviceName` | `string` |  |  |  |
| `spec.resourceConfigurationArn` | `string` |  |  |  |
| `spec.serviceNetworkArn` | `string` |  |  |  |
| `spec.routeTableIds` | `[]string \| valueFrom` |  |  | AwsSubnet (`status.outputs.route_table_id`) |
| `spec.subnetIds` | `[]string \| valueFrom` |  |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.securityGroupIds` | `[]string \| valueFrom` |  |  | AwsSecurityGroup (`status.outputs.security_group_id`) |
| `spec.privateDnsEnabled` | `bool` |  |  |  |
| `spec.dnsOptions` | `AwsVpcEndpointDnsOptions` |  |  |  |
| `spec.dnsOptions.dnsRecordIpType` | `string` |  |  |  |
| `spec.dnsOptions.privateDnsOnlyForInboundResolverEndpoint` | `bool` |  |  |  |
| `spec.dnsOptions.privateDnsPreference` | `string` |  |  |  |
| `spec.dnsOptions.privateDnsSpecifiedDomains` | `[]string` |  |  |  |
| `spec.ipAddressType` | `string` |  |  |  |
| `spec.policy` | `string` |  |  |  |
| `spec.subnetConfigurations` | `[]AwsVpcEndpointSubnetConfiguration` |  |  |  |
| `spec.subnetConfigurations[].subnetId` | `string \| valueFrom` | yes |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.subnetConfigurations[].ipv4` | `string` |  |  |  |
| `spec.subnetConfigurations[].ipv6` | `string` |  |  |  |
| `spec.serviceRegion` | `string` |  |  |  |
| `spec.autoAccept` | `bool` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region the endpoint is created in. Must match the VPC's
region. Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.vpcId

`string | valueFrom` · required

The VPC the endpoint lives in. Reference an AwsVpc's vpc_id output
or pass a literal VPC id. Create-only in AWS.

- references: AwsVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.endpointType

`string`

The endpoint type. "Gateway" (S3/DynamoDB via route tables),
"Interface" (ENI-based PrivateLink -- most AWS services and every
third-party PrivateLink service), "GatewayLoadBalancer" (fronts a
GWLB appliance fleet), "Resource" and "ServiceNetwork" (VPC
Lattice). Empty defaults to "Gateway" -- AWS's own default.
Create-only in AWS: changing the type replaces the endpoint.

### spec.serviceName

`string`

The AWS service to connect to, e.g. "com.amazonaws.us-west-2.s3"
for S3 in us-west-2 or a PrivateLink provider's
"com.amazonaws.vpce.<region>.vpce-svc-..." name. Exactly one of
service_name, resource_configuration_arn, or service_network_arn
must be set. Create-only in AWS.

### spec.resourceConfigurationArn

`string`

The ARN of a VPC Lattice resource configuration to connect to --
the "Resource" endpoint type's target. Exactly one of the three
service-target fields must be set. Create-only in AWS.

### spec.serviceNetworkArn

`string`

The ARN of a VPC Lattice service network to connect to -- the
"ServiceNetwork" endpoint type's target. Exactly one of the three
service-target fields must be set. Create-only in AWS.

### spec.routeTableIds

`[]string | valueFrom`

Route tables a GATEWAY endpoint injects its prefix-list route into.
Traffic to the service from any subnet on these tables flows
through the endpoint. Reference an AwsSubnet's route_table_id
output when the subnet owns its table (inline routes), or the
AwsVpc's main_route_table_id / default_route_table_id outputs when
subnets ride the VPC main table; literals also work. Gateway
endpoints only. Updates in place.

- references: AwsSubnet (`status.outputs.route_table_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.route_table_id}} -- a bare string does not parse

### spec.subnetIds

`[]string | valueFrom`

Subnets an INTERFACE / GatewayLoadBalancer / Resource /
ServiceNetwork endpoint places its network interfaces in -- one ENI
per subnet, so spread across AZs for availability (each AZ is
billed separately for interface endpoints). Reference AwsSubnet
subnet_id outputs or pass literal subnet ids. Not applicable to
gateway endpoints. Updates in place.

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.securityGroupIds

`[]string | valueFrom`

Security groups attached to an INTERFACE endpoint's network
interfaces -- they must allow inbound traffic from the clients on
the service's port (443 for AWS APIs). Empty means AWS attaches
the VPC's DEFAULT security group. Reference AwsSecurityGroup
security_group_id outputs or pass literal ids. Interface endpoints
only. Updates in place.

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.privateDnsEnabled

`bool`

Resolve the service's PUBLIC DNS name (e.g. sts.us-west-2.
amazonaws.com) to the endpoint's private IPs inside the VPC, via an
AWS-managed private hosted zone. Most interface-endpoint users want
this on -- clients keep their default SDK endpoints and privately
reach the service with zero code changes. Requires the VPC to have
BOTH DNS support and DNS hostnames enabled. Interface endpoints
only (AWS replaces the endpoint if this changes on other types --
on Interface it updates in place). Off by default, matching AWS.

### spec.dnsOptions

`AwsVpcEndpointDnsOptions`

Fine-grained DNS behavior for the endpoint. Only meaningful when
the endpoint creates DNS records (interface endpoints, and the
Lattice types for the preference/domain fields).

- rule: dns_record_ip_type must be 'ipv4', 'ipv6', 'dualstack', or 'service-defined' when set
- rule: private_dns_preference must be 'ALL_DOMAINS', 'VERIFIED_DOMAINS_ONLY', 'VERIFIED_DOMAINS_AND_SPECIFIED_DOMAINS', or 'SPECIFIED_DOMAINS_ONLY' when set
- rule: private_dns_specified_domains is required when private_dns_preference includes specified domains, and must be empty otherwise

### spec.dnsOptions.dnsRecordIpType

`string`

Which record types the endpoint's DNS names resolve to: "ipv4",
"ipv6", "dualstack", or "service-defined" (the service picks).
Empty lets AWS choose based on the endpoint's ip_address_type.

### spec.dnsOptions.privateDnsOnlyForInboundResolverEndpoint

`bool`

Route only INBOUND (on-premises / cross-VPC) resolver traffic
through this endpoint's private DNS, keeping in-VPC traffic on the
gateway endpoint. Applies to the S3 dual-stack pattern -- a service
with BOTH a gateway and an interface endpoint in the same VPC:
in-VPC S3 traffic rides the free gateway while on-premises clients
resolve to the interface endpoint. Requires private_dns_enabled.

### spec.dnsOptions.privateDnsPreference

`string`

Which private domains get a private hosted zone, for Resource /
ServiceNetwork endpoints: "ALL_DOMAINS", "VERIFIED_DOMAINS_ONLY",
"VERIFIED_DOMAINS_AND_SPECIFIED_DOMAINS", or
"SPECIFIED_DOMAINS_ONLY". Empty keeps the AWS default. Create-only
in AWS.

### spec.dnsOptions.privateDnsSpecifiedDomains

`[]string`

The private domains to create hosted zones for -- required when
private_dns_preference includes specified domains
("VERIFIED_DOMAINS_AND_SPECIFIED_DOMAINS" or
"SPECIFIED_DOMAINS_ONLY"), and must be empty otherwise. AWS allows
1-10 domains. Create-only in AWS.

- rule: {"repeated":{"maxItems":"10"}}

### spec.ipAddressType

`string`

The IP address type of the endpoint: "ipv4", "dualstack", or
"ipv6". Empty lets AWS pick based on the service and subnets --
effectively ipv4 for nearly every service today. The service must
support the chosen type. Updates in place.

### spec.policy

`string`

An IAM policy document (JSON) controlling which principals may use
the endpoint to reach which resources -- e.g. "only this account's
buckets" on an S3 gateway endpoint, turning the endpoint into a
data-exfiltration control. Empty means full access (AWS's
default). All gateway and most interface endpoints support
policies. Updates in place.

### spec.subnetConfigurations

`[]AwsVpcEndpointSubnetConfiguration`

Pin specific IPv4/IPv6 addresses for the endpoint's ENI in chosen
subnets -- for appliances or firewall rules that need stable
endpoint IPs. Every subnet_id listed here must also appear in
subnet_ids. Rarely needed: AWS assigns addresses automatically when
omitted. Updates in place.

### spec.subnetConfigurations[].subnetId

`string | valueFrom` · required

The subnet whose ENI addresses are being pinned. Reference an
AwsSubnet's subnet_id output or pass a literal subnet id -- and
list the same subnet in spec.subnet_ids.

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.subnetConfigurations[].ipv4

`string`

The private IPv4 address to assign to the ENI in this subnet. Must
be a free address inside the subnet's CIDR.

### spec.subnetConfigurations[].ipv6

`string`

The IPv6 address to assign to the ENI in this subnet. Requires a
dualstack or ipv6 endpoint and an IPv6-enabled subnet.

### spec.serviceRegion

`string`

Connect to a service in ANOTHER region (interface endpoints only),
e.g. reach us-east-1-only services from a us-west-2 VPC without
cross-region networking of your own. Empty means the endpoint's own
region. Create-only in AWS.

### spec.autoAccept

`bool`

Accept the endpoint connection automatically when the PrivateLink
service requires acceptance and lives in the SAME AWS account.
Cross-account services must accept on their side regardless.

## Validation Rules

- `service_target_exactly_one`: exactly one of service_name, resource_configuration_arn, or service_network_arn must be set
- `endpoint_type_valid`: endpoint_type must be one of 'Gateway', 'Interface', 'GatewayLoadBalancer', 'Resource', or 'ServiceNetwork' when set (empty defaults to Gateway)
- `lattice_targets_match_type`: resource_configuration_arn requires endpoint_type 'Resource', and service_network_arn requires endpoint_type 'ServiceNetwork'
- `route_tables_gateway_only`: route_table_ids apply only to Gateway endpoints (the default type when endpoint_type is empty)
- `eni_fields_not_gateway`: subnet_ids, security_group_ids, private_dns_enabled, subnet_configurations, and service_region do not apply to Gateway endpoints -- set endpoint_type to 'Interface' (or another ENI-based type)
- `sg_and_private_dns_interface_only`: security_group_ids and private_dns_enabled apply only to Interface endpoints
- `service_region_interface_only`: service_region (cross-region endpoints) applies only to Interface endpoints
- `ip_address_type_valid`: ip_address_type must be 'ipv4', 'dualstack', or 'ipv6' when set

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsVpcEndpoint, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.vpc_endpoint_id` | `string` | vpc_endpoint_id is the endpoint's id (e.g. "vpce-0abc..."), the handle every AWS API and route inspection uses. |
| `status.outputs.arn` | `string` | arn is the endpoint's Amazon Resource Name -- arn:aws:ec2:<region>:<account>:vpc-endpoint/<vpce-id>. |
| `status.outputs.state` | `string` | state is the endpoint's lifecycle state after provisioning -- "available" on a successful create ("pendingAcceptance" when the PrivateLink service requires manual acceptance). |
| `status.outputs.prefix_list_id` | `string` | prefix_list_id is the service's prefix list (e.g. "pl-68a54001" for S3) -- gateway endpoints only. Reference it in security-group or route rules that must scope traffic to the service's address ranges. Empty for ENI-based endpoint types. |
| `status.outputs.dns_name` | `string` | dns_name is the endpoint's primary private DNS name -- the regional endpoint-specific name clients use when private DNS is off, and the name Route53 aliases target. Interface endpoints only; empty for gateway endpoints (which have no DNS presence). |
| `status.outputs.hosted_zone_id` | `string` | hosted_zone_id is the Route53 hosted zone of dns_name, needed alongside it when creating a Route53 alias record to the endpoint. Interface endpoints only. |
| `status.outputs.network_interface_ids` | `[]string` | network_interface_ids are the endpoint's ENIs, one per attached subnet -- the objects flow logs, firewall rules, and IP lookups operate on. Empty for gateway endpoints. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | AwsVpc | `status.outputs.vpc_id` |
| `spec.routeTableIds` | AwsSubnet | `status.outputs.route_table_id` |
| `spec.subnetIds` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.securityGroupIds` | AwsSecurityGroup | `status.outputs.security_group_id` |
| `spec.subnetConfigurations[].subnetId` | AwsSubnet | `status.outputs.subnet_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsAppRunnerService | `spec.vpcIngressConnections[].vpcEndpointId` | `status.outputs.vpc_endpoint_id` |

## See Also

- [Overview](../README.md)
