---
title: "Client VPN"
description: "Client VPN deployment documentation"
icon: "package"
order: 100
componentName: "awsclientvpn"
---

# AWS Client VPN

Creates an AWS Client VPN endpoint — the managed OpenVPN server for secure
remote access into AWS networks — with its folded target network
associations, authorization rules, and routes, so one manifest declares
who may connect, what they may reach, and how traffic flows.

## What Gets Created

When you deploy an AwsClientVpn resource, Planton provisions:

- **Client VPN endpoint** — an `aws_ec2_client_vpn_endpoint` /
  `ec2clientvpn.Endpoint` carrying authentication, tunnel shape, session,
  logging, and client-experience settings
- **Target network associations** — one per entry in `subnetIds`, each
  attaching the endpoint to that subnet's Availability Zone (associations
  bill hourly and take several minutes to attach/detach)
- **Authorization rules** — one per entry in `authorizationRules`; nothing
  is reachable until a rule authorizes it
- **Routes** — one per entry in `routes`, beyond the per-VPC route AWS
  auto-creates for each associated subnet

## Prerequisites

- **AWS credentials** configured via the Planton provider config (keyless
  SSO/OIDC).
- **An ACM server certificate** in the same region (reference an
  AwsCertManagerCert; an imported self-signed certificate works for
  private setups).
- **For certificate authentication** — a client CA chain certificate in
  ACM (may be the same imported certificate in a self-signed setup).
- **For federated authentication** — an IAM SAML identity provider.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsClientVpn
metadata:
  name: corp-access
spec:
  region: us-west-2
  authenticationOptions:
    - type: certificate-authentication
      rootCertificateChainArn:
        valueFrom:
          kind: AwsCertManagerCert
          name: vpn-ca
          fieldPath: status.outputs.cert_arn
  serverCertificateArn:
    valueFrom:
      kind: AwsCertManagerCert
      name: vpn-server
      fieldPath: status.outputs.cert_arn
  clientCidrBlock: 10.100.0.0/22
  splitTunnel: true
  vpcId:
    valueFrom:
      kind: AwsVpc
      name: main-vpc
      fieldPath: status.outputs.vpc_id
  subnetIds:
    - valueFrom:
        kind: AwsSubnet
        name: private-a
        fieldPath: status.outputs.subnet_id
  authorizationRules:
    - targetNetworkCidr: 10.0.0.0/16
      authorizeAllGroups: true
      description: reach the whole VPC
  connectionLog:
    cloudwatchLogGroup:
      valueFrom:
        kind: AwsCloudwatchLogGroup
        name: vpn-logs
        fieldPath: status.outputs.log_group_name
```

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `region` | `string` | AWS region for the endpoint |
| `authenticationOptions` | `AwsClientVpnAuthenticationOption[]` | 1–2 options; a client passes on ANY. All ForceNew. |
| `serverCertificateArn` | `StringValueOrRef` | The ACM certificate the VPN server presents — required for every auth type. Updates in place. |

### Authentication Options

| Field | Type | Description |
|-------|------|-------------|
| `type` | `string` | `certificate-authentication`, `directory-service-authentication`, or `federated-authentication` |
| `rootCertificateChainArn` | `StringValueOrRef` | Certificate type: the client CA chain (AwsCertManagerCert) |
| `activeDirectoryId` | `string` | Directory type: the AWS Directory Service directory ID |
| `samlProviderArn` | `string` | Federated type: the IAM SAML provider ARN |
| `selfServiceSamlProviderArn` | `string` | Federated type, optional: a separate SAML app for the self-service portal |

### Tunnel Shape

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `clientCidrBlock` | `string` | — | /22–/12 client address pool. Required unless `trafficIpAddressType: ipv6` (where it must be empty). ForceNew. |
| `splitTunnel` | `bool` | `false` | `true`: only routed traffic enters the VPN. `false` (full tunnel): ALL client traffic — pair with a 0.0.0.0/0 route + rule through a NAT-ed subnet. |
| `transportProtocol` | `string` | `udp` | `udp` or `tcp`. ForceNew. |
| `vpnPort` | `int32` | `443` | `443` or `1194` — either works with either protocol. |
| `endpointIpAddressType` | `string` | `ipv4` | Which stacks clients reach the endpoint over: `ipv4`, `ipv6`, `dual-stack`. ForceNew. |
| `trafficIpAddressType` | `string` | `ipv4` | Addressing inside the tunnel: `ipv4`, `ipv6`, `dual-stack`. ForceNew. |

### Network Attachment

| Field | Type | Description |
|-------|------|-------------|
| `vpcId` | `StringValueOrRef` | The VPC whose security groups apply. Optional — inferred from the first associated subnet. Mutually exclusive with the TGW arm. |
| `securityGroupIds` | `StringValueOrRef[]` | Max 5. When omitted, the VPC default security group applies. |
| `subnetIds` | `StringValueOrRef[]` | Folded target network associations (may be empty — a pre-staged endpoint routes no traffic). |
| `transitGatewayConfiguration` | `object` | Attach to an AwsTransitGateway instead of a VPC (AZ names XOR AZ IDs). ForceNew; excludes `vpcId`/`securityGroupIds`/`subnetIds`. |

### Access Control

| Field | Type | Description |
|-------|------|-------------|
| `authorizationRules[]` | `object` | `targetNetworkCidr` + exactly one of `accessGroupId` (one IdP group) or `authorizeAllGroups`. |
| `routes[]` | `object` | `destinationCidrBlock` + `targetSubnetId` (must be one of `subnetIds`). |

### Sessions and Client Experience

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `sessionTimeoutHours` | `int32` | `24` | 8, 10, 12, or 24 |
| `disconnectOnSessionTimeout` | `bool` | `false` | Hard-disconnect at timeout instead of auto-reconnect |
| `selfServicePortalEnabled` | `bool` | `false` | Federated auth only |
| `clientConnectOptions` | `object` | — | Posture-check Lambda (name must start with `AWSClientVPN-`) run on every connection |
| `clientLoginBanner` | `object` | — | Banner text (≤1400 chars) shown on session establishment |
| `clientRouteEnforcementEnabled` | `bool` | `false` | Block client-side route manipulation |
| `dnsServers` | `string[]` | device DNS | Up to 2 IPs pushed to clients (set the VPC resolver for private zone resolution) |
| `connectionLog` | `object` | — | CloudWatch log group (ref) + optional stream; presence enables logging |

## Full-Tunnel Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsClientVpn
metadata:
  name: secure-egress
spec:
  region: us-west-2
  authenticationOptions:
    - type: certificate-authentication
      rootCertificateChainArn:
        value: arn:aws:acm:us-west-2:123456789012:certificate/client-ca
  serverCertificateArn:
    value: arn:aws:acm:us-west-2:123456789012:certificate/server
  clientCidrBlock: 10.100.0.0/22
  splitTunnel: false
  subnetIds:
    - value: subnet-nat-routed
  authorizationRules:
    - targetNetworkCidr: 0.0.0.0/0
      authorizeAllGroups: true
  routes:
    - destinationCidrBlock: 0.0.0.0/0
      targetSubnetId:
        value: subnet-nat-routed
      description: internet egress via the NAT-ed subnet
  clientRouteEnforcementEnabled: true
```

## Stack Outputs

| Output | Description |
|--------|-------------|
| `client_vpn_endpoint_id` | The `cvpn-endpoint-...` identifier (used to export the client configuration) |
| `client_vpn_endpoint_arn` | The endpoint ARN for IAM policies |
| `endpoint_dns_name` | The DNS name clients connect to |
| `self_service_portal_url` | Portal URL (empty when disabled) |
| `subnet_association_ids` | Map of subnet ID → association ID |
| `transit_gateway_attachment_id` | The TGW attachment (empty for VPC-attached endpoints) |

## Related Components

- [AwsCertManagerCert](/docs/catalog/aws/cert-manager-cert) — server certificate and client CA chain
- [AwsVpc](/docs/catalog/aws/vpc) / [AwsSubnet](/docs/catalog/aws/subnet) — the networks clients reach
- [AwsSecurityGroup](/docs/catalog/aws/security-group) — traffic control between clients and VPC resources
- [AwsCloudwatchLogGroup](/docs/catalog/aws/cloudwatch-log-group) — connection logging
- [AwsLambda](/docs/catalog/aws/lambda) — the connection posture-check hook
- [AwsTransitGateway](/docs/catalog/aws/transit-gateway) — hub attachment for multi-network reach
