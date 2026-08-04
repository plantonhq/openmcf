# AwsAlb

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsAlbSpec defines an Application Load Balancer: the Layer-7 entry point
that terminates HTTP/HTTPS and hands requests to the routing graph.

The load balancer itself carries no routing configuration -- that is
deliberate. Listeners (AwsLbListener) attach to it and own ports, TLS
material, and default actions; listener rules (AwsLbListenerRule) attach to
listeners and own per-service routing; target groups (AwsLbTargetGroup)
receive the traffic. This spec owns only what is truly load-balancer-wide:
placement (subnets, scheme), security groups, and the HTTP behavior knobs
AWS models as load balancer attributes.

The ALB name comes from metadata.name. AWS limits the name to 32
characters; both IaC modules truncate longer names deterministically.
Name, scheme (internal), and subnets-vs-AZ coverage aside, most attributes
update in place; changing "internal" replaces the load balancer.

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAlb
metadata:
  name: awsalb-demo
spec:
  region: us-west-2
  subnets:
    - value: subnet-12345678
    - value: subnet-12345679
  # Behavior attributes deliberately exercise the tri-state bool
  # (http2Enabled), the enum-validated strings, and the log-delivery object
  # so the fixture proves the full variable contract, not just the scalars.
  internal: false
  ipAddressType: ipv4
  idleTimeoutSeconds: 120
  http2Enabled: true
  desyncMitigationMode: defensive
  xffHeaderProcessingMode: append
  accessLogs:
    bucket:
      value: demo-alb-logs-bucket
    prefix: alb/demo
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.subnets` | `[]string \| valueFrom` | yes |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.securityGroups` | `[]string \| valueFrom` |  |  | AwsSecurityGroup (`status.outputs.security_group_id`) |
| `spec.internal` | `bool` |  |  |  |
| `spec.ipAddressType` | `string` |  |  |  |
| `spec.deleteProtectionEnabled` | `bool` |  |  |  |
| `spec.idleTimeoutSeconds` | `int32` |  | `60` |  |
| `spec.clientKeepAliveSeconds` | `int32` |  |  |  |
| `spec.http2Enabled` | `bool` |  |  |  |
| `spec.wafFailOpenEnabled` | `bool` |  |  |  |
| `spec.webAclArn` | `string \| valueFrom` |  |  | AwsWafWebAcl (`status.outputs.web_acl_arn`) |
| `spec.zonalShiftEnabled` | `bool` |  |  |  |
| `spec.dropInvalidHeaderFields` | `bool` |  |  |  |
| `spec.preserveHostHeader` | `bool` |  |  |  |
| `spec.xffClientPortEnabled` | `bool` |  |  |  |
| `spec.xffHeaderProcessingMode` | `string` |  |  |  |
| `spec.desyncMitigationMode` | `string` |  |  |  |
| `spec.tlsVersionAndCipherSuiteHeadersEnabled` | `bool` |  |  |  |
| `spec.accessLogs` | `AwsAlbLogDelivery` |  |  |  |
| `spec.accessLogs.bucket` | `string \| valueFrom` | yes |  | AwsS3Bucket (`status.outputs.bucket_id`) |
| `spec.accessLogs.prefix` | `string` |  |  |  |
| `spec.connectionLogs` | `AwsAlbLogDelivery` |  |  |  |
| `spec.connectionLogs.bucket` | `string \| valueFrom` | yes |  | AwsS3Bucket (`status.outputs.bucket_id`) |
| `spec.connectionLogs.prefix` | `string` |  |  |  |
| `spec.healthCheckLogs` | `AwsAlbLogDelivery` |  |  |  |
| `spec.healthCheckLogs.bucket` | `string \| valueFrom` | yes |  | AwsS3Bucket (`status.outputs.bucket_id`) |
| `spec.healthCheckLogs.prefix` | `string` |  |  |  |
| `spec.dns` | `AwsAlbDns` |  |  |  |
| `spec.dns.enabled` | `bool` |  |  |  |
| `spec.dns.route53ZoneId` | `string \| valueFrom` |  |  | AwsRoute53Zone (`status.outputs.zone_id`) |
| `spec.dns.hostnames` | `[]string` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the ALB is created.
Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.subnets

`[]string | valueFrom` · required

The subnets the ALB places its nodes in -- at least two, in different
Availability Zones (an AWS requirement that also buys zonal redundancy).
Public subnets for internet-facing ALBs, private for internal ones.

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: {"required":true,"repeated":{"minItems":"2"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.securityGroups

`[]string | valueFrom`

Security groups controlling traffic to and from the ALB. When omitted,
AWS attaches the VPC's default security group -- fine for a first boot,
wrong for production; attach explicit groups that open exactly the
listener ports.

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.internal

`bool`

When true, the ALB is internal (reachable only inside the VPC); when
false (default), it is internet-facing with public DNS. Immutable:
changing the scheme replaces the load balancer.

### spec.ipAddressType

`string`

The address family of the ALB's nodes. Valid values: "ipv4" (default),
"dualstack" (IPv4 + IPv6), "dualstack-without-public-ipv4" (public IPv6
with private IPv4 -- avoids public-IPv4 charges for IPv6-capable
clients).

### spec.deleteProtectionEnabled

`bool`

Prevents deletion of the ALB while enabled. Recommended for production:
deleting an ALB silently orphans every listener and rule attached to it.

### spec.idleTimeoutSeconds

`int32`

Seconds an idle connection stays open. Range 1-4000. AWS default: 60.
Raise it above the application's slowest response time to avoid 504s on
long-running requests; keep it above any upstream keep-alive interval.

- default: `60`

### spec.clientKeepAliveSeconds

`int32`

Seconds an HTTP client connection may stay alive across requests.
Range 60-604800. AWS default: 3600.

### spec.http2Enabled

`bool` · optional (explicit presence)

Whether HTTP/2 is offered to clients. AWS default: true. Optional rather
than plain bool so that false ("disable HTTP/2") is distinguishable from
unset ("keep the AWS default").

### spec.wafFailOpenEnabled

`bool`

What happens to requests when an attached WAF is unreachable: when true,
requests pass through ("fail open"); when false (AWS default), they are
rejected ("fail closed"). A deliberate availability-versus-security call.
Only meaningful when web_acl_arn attaches a WAF to this ALB.

### spec.webAclArn

`string | valueFrom`

The REGIONAL-scope WAFv2 web ACL protecting this ALB, by ARN — the
modules create the web-ACL association alongside the load balancer.
An ALB has at most one web ACL; leave unset for no WAF. The web ACL
must live in the same region as the ALB. Can reference an
AwsWafWebAcl resource.

- references: AwsWafWebAcl (`status.outputs.web_acl_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsWafWebAcl, name: <that resource's name>, fieldPath: status.outputs.web_acl_arn}} -- a bare string does not parse

### spec.zonalShiftEnabled

`bool`

Allows Amazon Application Recovery Controller to shift this ALB's
traffic away from an impaired Availability Zone.

### spec.dropInvalidHeaderFields

`bool`

Drop request headers with names that are not valid HTTP header fields
instead of forwarding them. Hardens against header-smuggling tricks;
AWS default: false.

### spec.preserveHostHeader

`bool`

Forward the client's original Host header to the target unchanged
instead of rewriting it to the target address. AWS default: false.

### spec.xffClientPortEnabled

`bool`

Append the client's source port to the X-Forwarded-For header. AWS
default: false.

### spec.xffHeaderProcessingMode

`string`

How the ALB handles the X-Forwarded-For header. Valid values: "append"
(AWS default -- add the client IP), "preserve" (pass it through
untouched), "remove" (strip it). "preserve"/"remove" matter when the ALB
sits behind another proxy layer whose XFF chain must win.

### spec.desyncMitigationMode

`string`

Protection level against HTTP desync (request-smuggling) attacks.
Valid values: "monitor" (classify only), "defensive" (AWS default --
block ambiguous requests likely to poison caches), "strictest" (block
everything not RFC 7230 compliant).

### spec.tlsVersionAndCipherSuiteHeadersEnabled

`bool`

Inject the negotiated TLS version and cipher suite as request headers
(x-amzn-tls-version, x-amzn-tls-cipher-suite) toward targets, for
applications that audit their clients' TLS posture. AWS default: false.

### spec.accessLogs

`AwsAlbLogDelivery`

Access logs: one entry per request, delivered to S3. The bucket must
carry the ELB log-delivery bucket policy. When omitted, access logging
is off.

### spec.accessLogs.bucket

`string | valueFrom` · required

The S3 bucket receiving the logs.

- references: AwsS3Bucket (`status.outputs.bucket_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsS3Bucket, name: <that resource's name>, fieldPath: status.outputs.bucket_id}} -- a bare string does not parse

### spec.accessLogs.prefix

`string`

Key prefix inside the bucket (e.g. "alb/production"), for sharing one
bucket across several load balancers or log types.

### spec.connectionLogs

`AwsAlbLogDelivery`

Connection logs: one entry per client connection (TLS handshake
details, client address) -- the place TLS negotiation failures that
never become requests show up. Delivered to S3; when omitted,
connection logging is off.

### spec.connectionLogs.bucket

`string | valueFrom` · required

The S3 bucket receiving the logs.

- references: AwsS3Bucket (`status.outputs.bucket_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsS3Bucket, name: <that resource's name>, fieldPath: status.outputs.bucket_id}} -- a bare string does not parse

### spec.connectionLogs.prefix

`string`

Key prefix inside the bucket (e.g. "alb/production"), for sharing one
bucket across several load balancers or log types.

### spec.healthCheckLogs

`AwsAlbLogDelivery`

Health-check logs: one entry per health-check result, for debugging
flapping targets without packet captures. Delivered to S3; when
omitted, health-check logging is off.

### spec.healthCheckLogs.bucket

`string | valueFrom` · required

The S3 bucket receiving the logs.

- references: AwsS3Bucket (`status.outputs.bucket_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsS3Bucket, name: <that resource's name>, fieldPath: status.outputs.bucket_id}} -- a bare string does not parse

### spec.healthCheckLogs.prefix

`string`

Key prefix inside the bucket (e.g. "alb/production"), for sharing one
bucket across several load balancers or log types.

### spec.dns

`AwsAlbDns`

Optional Route53 DNS: alias A records pointing the given hostnames at
the ALB. Alias records are preferred over CNAMEs because they work at
the zone apex, cost nothing per query, and inherit the ALB's health.

### spec.dns.enabled

`bool`

When true, creates Route53 alias records for the ALB.

### spec.dns.route53ZoneId

`string | valueFrom`

Route53 hosted zone ID where alias records are created.
Required when enabled is true.

- references: AwsRoute53Zone (`status.outputs.zone_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsRoute53Zone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.dns.hostnames

`[]string`

Domain names that will point to the ALB via Route53 alias records.
Each hostname gets its own A record aliased to the ALB's DNS name.

- rule: {"repeated":{"unique":true}}

## Validation Rules

- `ip_address_type_valid`: ip_address_type must be 'ipv4', 'dualstack', or 'dualstack-without-public-ipv4' when set
- `idle_timeout_range`: idle_timeout_seconds must be between 1 and 4000 when set
- `client_keep_alive_range`: client_keep_alive_seconds must be between 60 and 604800 when set
- `xff_header_processing_mode_valid`: xff_header_processing_mode must be 'append', 'preserve', or 'remove' when set
- `desync_mitigation_mode_valid`: desync_mitigation_mode must be 'monitor', 'defensive', or 'strictest' when set

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsAlb, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.loadBalancerArn` | `string` | load_balancer_arn is the ARN of the created Application Load Balancer. |
| `status.outputs.loadBalancerName` | `string` | load_balancer_name is the final name assigned to the ALB (may differ from metadata.name). |
| `status.outputs.loadBalancerDnsName` | `string` | load_balancer_dns_name is the DNS name automatically assigned to the ALB. |
| `status.outputs.loadBalancerHostedZoneId` | `string` | load_balancer_hosted_zone_id is the Route53 hosted zone ID for the ALB's DNS entry. |
| `status.outputs.arnSuffix` | `string` | arn_suffix is the ARN suffix (e.g. "app/my-alb/50dc6c495c0c9188") used as the LoadBalancer dimension in CloudWatch metrics -- the handle alarms, dashboards, and request-count autoscaling policies need. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.subnets` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.securityGroups` | AwsSecurityGroup | `status.outputs.security_group_id` |
| `spec.webAclArn` | AwsWafWebAcl | `status.outputs.web_acl_arn` |
| `spec.accessLogs.bucket` | AwsS3Bucket | `status.outputs.bucket_id` |
| `spec.connectionLogs.bucket` | AwsS3Bucket | `status.outputs.bucket_id` |
| `spec.healthCheckLogs.bucket` | AwsS3Bucket | `status.outputs.bucket_id` |
| `spec.dns.route53ZoneId` | AwsRoute53Zone | `status.outputs.zone_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
