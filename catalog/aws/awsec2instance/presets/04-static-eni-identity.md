# Static Network Identity (Pre-Provisioned ENI)

This preset attaches a pre-provisioned network interface (ENI) as the instance's primary interface instead of creating one. The ENI -- not the instance -- owns the network identity: the subnet, the security groups, the fixed private IP, and the MAC address. Replace the instance (a new AMI, a bigger type) and every one of them stays put. The shape for network appliances, license-bound software keyed to a MAC, and any endpoint whose IP is baked into firewall rules or DNS.

## When to Use

- Virtual appliances (firewalls, routers, VPN concentrators) whose identity must survive instance replacement
- License servers where the vendor binds the license to a MAC address
- Fixed-IP endpoints referenced by external firewall rules, allow-lists, or DNS you do not control

## Key Configuration Choices

- **Pre-provisioned primary ENI** (`primaryNetworkInterfaceId`) -- The ENI carries the whole network configuration; the instance simply inherits it at launch and releases it intact at termination
- **No inline networking** -- `subnetId`, `securityGroupIds`, addressing fields, and `sourceDestCheck` must stay unset: the ENI defines them, and both validation and the EC2 API reject the mix. For router/NAT appliances, disable source/dest checking on the ENI itself
- **Instance stays replaceable** -- Everything instance-side (AMI, type, root volume) can change freely; the identity lives on the interface

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region where the instance will be created | AWS region list |
| `eni-0123456789abcdef0` | Pre-provisioned network interface ID | EC2 console > Network Interfaces |
| `ami-0123456789abcdef0` | AMI matching the appliance/workload | AWS EC2 AMI catalog or vendor listing |
| `<instance-profile-name>` | NAME of the IAM instance profile | `AwsIamInstanceProfile` status outputs |

## Related Presets

- **01-ssm-managed** -- The baseline where the instance creates its own interface inline
- **03-spot-worker** -- Interruption-tolerant workers (a static identity and Spot rarely mix)
