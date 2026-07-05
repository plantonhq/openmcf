# PrivateLink-Hardened NLB

This preset creates the NLB posture for exposing a service over AWS
PrivateLink: an internal load balancer with security groups attached,
security-group rules enforced on traffic arriving through consumer VPC
endpoints, deletion protection, and access logs to S3. An internal NLB is
the required target of a VPC endpoint service, and this configuration keeps
your own inbound rules in the path even for endpoint traffic. Attach
`AwsLbListener` resources (TLS listeners if you want the traffic to appear
in access logs) against its `load_balancer_arn` output.

## When to Use

- Backing a VPC endpoint service that other AWS accounts or VPCs consume via
  PrivateLink
- Internal NLBs in regulated environments where "no security groups" is not
  an acceptable posture
- Any Layer-4 service where you want defense in depth: endpoint policies
  *and* your own security-group rules

## Key Configuration Choices

- **Security groups attached** — a deliberate one-way door: once attached,
  the last group can never be removed, only replaced. This preset commits
  because PrivateLink exposure is exactly the case where inbound filtering
  earns its keep
- **`enforceSecurityGroupInboundRulesOnPrivateLinkTraffic: "on"`** — traffic
  arriving through PrivateLink endpoints is evaluated against the NLB's
  inbound rules instead of bypassing them; the consumer's endpoint policy
  alone is not the last line of defense
- **Access logs to S3** — NLB access logs capture TLS-listener traffic only
  (an AWS limitation); the bucket must carry the ELB log-delivery bucket
  policy or delivery fails silently
- **Internal scheme + deletion protection** — the endpoint-service target
  should be unreachable from the internet and should refuse casual deletion
  (which would orphan listeners and break every consumer endpoint)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<nlb-name>` | Unique name for the NLB (AWS caps it at 32 characters) | Choose a descriptive name (e.g., `endpoint-svc-nlb`) |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `<private-subnet-id-az1>` | Private subnet in the first Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `<private-subnet-id-az2>` | Private subnet in the second Availability Zone | AWS VPC console or `AwsSubnet` status outputs |
| `<nlb-security-group-resource-name>` | `AwsSecurityGroup` resource admitting exactly the listener ports | Your AwsSecurityGroup manifest's `metadata.name` |
| `<access-logs-bucket-name>` | S3 bucket with the ELB log-delivery bucket policy | AWS S3 console or `AwsS3Bucket` status outputs |
| `<log-prefix>` | Key prefix inside the bucket (e.g., `nlb/endpoint-svc`) | Your logging layout |

## Common Additions

- Add `zonalShiftEnabled: true` to let Amazon Application Recovery Controller
  drain an impaired AZ
- Add `crossZoneLoadBalancingEnabled: true` if targets are unevenly
  distributed across AZs
- Set `privateIpv4Address` on the mappings if consumers reference the NLB
  nodes by fixed address

## Related Presets

- **01-internal** — the same scheme without the hardening commitments
- **02-static-ip-internet-facing** — the public variant with Elastic-IP
  static addresses
