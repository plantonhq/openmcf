# AWS Client VPN Examples

Below are several examples demonstrating how to define an AWS Client VPN component in
OpenMCF. After creating one of these YAML manifests, apply it with Terraform using the OpenMCF CLI:

```shell
openmcf tofu apply --manifest <yaml-path> --stack <stack-name>
```

---

## Basic Client VPN Endpoint

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsClientVpn
metadata:
  name: basic-client-vpn
spec:
  vpcId:
    value: vpc-12345678
  subnets:
    - value: subnet-abc123
  clientCidrBlock: 10.0.0.0/22
  serverCertificateArn:
    value: arn:aws:acm:us-east-1:123456789012:certificate/abc
  cidrAuthorizationRules:
    - 10.0.0.0/16
```

This example creates a basic Client VPN endpoint:
• Uses certificate-based authentication (default).
• Connects to a single subnet in the VPC.
• Allows access to the entire VPC CIDR (10.0.0.0/16).
• Uses default port 443 and TCP protocol.

---

## Client VPN with Multiple Subnets

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsClientVpn
metadata:
  name: multi-subnet-vpn
spec:
  description: "Client VPN for development team access"
  vpcId:
    value: vpc-12345678
  subnets:
    - value: subnet-private-1a
    - value: subnet-private-1b
    - value: subnet-private-1c
  clientCidrBlock: 10.1.0.0/22
  serverCertificateArn:
    value: arn:aws:acm:us-east-1:123456789012:certificate/vpn-cert
  cidrAuthorizationRules:
    - 10.0.0.0/16
    - 172.16.0.0/12
  vpnPort: 443
  transportProtocol: tcp
```

This example provides access to multiple subnets:
• Connects to three private subnets across different AZs.
• Allows access to VPC and additional private networks.
• Uses custom client CIDR block for IP allocation.
• Includes description for better resource management.

---

## Client VPN with Custom Security Groups

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsClientVpn
metadata:
  name: secure-client-vpn
spec:
  vpcId:
    value: vpc-12345678
  subnets:
    - value: subnet-private-1a
  clientCidrBlock: 10.2.0.0/22
  serverCertificateArn:
    value: arn:aws:acm:us-east-1:123456789012:certificate/secure-vpn
  securityGroups:
    - value: sg-vpn-access
  cidrAuthorizationRules:
    - 10.0.0.0/16
  disableSplitTunnel: false
  logGroupName: /aws/client-vpn/secure-vpn
```

This example includes security and logging:
• Uses custom security group for fine-grained access control.
• Enables connection logging to CloudWatch.
• Maintains split-tunnel routing (only internal traffic through VPN).
• Restricts access to specific VPC resources.

---

## Client VPN with Custom DNS

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsClientVpn
metadata:
  name: dns-client-vpn
spec:
  vpcId:
    value: vpc-12345678
  subnets:
    - value: subnet-private-1a
  clientCidrBlock: 10.3.0.0/22
  serverCertificateArn:
    value: arn:aws:acm:us-east-1:123456789012:certificate/dns-vpn
  dnsServers:
    - 8.8.8.8
    - 8.8.4.4
  cidrAuthorizationRules:
    - 10.0.0.0/16
  vpnPort: 1194
  transportProtocol: udp
```

This example uses custom DNS configuration:
• Configures Google DNS servers for VPN clients.
• Uses UDP protocol on port 1194 for potentially better performance.
• Maintains standard authorization rules.

---

## Full-Tunnel Client VPN

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsClientVpn
metadata:
  name: full-tunnel-vpn
spec:
  description: "Full-tunnel VPN for complete traffic routing"
  vpcId:
    value: vpc-12345678
  subnets:
    - value: subnet-private-1a
    - value: subnet-private-1b
  clientCidrBlock: 10.4.0.0/22
  serverCertificateArn:
    value: arn:aws:acm:us-east-1:123456789012:certificate/full-tunnel
  disableSplitTunnel: true
  cidrAuthorizationRules:
    - 0.0.0.0/0
  logGroupName: /aws/client-vpn/full-tunnel
```

This example routes all traffic through the VPN:
• Disables split-tunnel to route all client traffic through VPN.
• Authorizes access to all networks (0.0.0.0/0).
• Includes logging for monitoring all traffic.
• Connects to multiple subnets for high availability.

---

## Minimal Client VPN

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsClientVpn
metadata:
  name: minimal-vpn
spec:
  vpcId:
    value: vpc-12345678
  subnets:
    - value: subnet-abc123
  clientCidrBlock: 10.5.0.0/22
  serverCertificateArn:
    value: arn:aws:acm:us-east-1:123456789012:certificate/minimal
```

A minimal configuration with:
• Only required fields specified.
• Uses default authentication (certificate-based).
• Default port 443 and TCP protocol.
• No authorization rules (clients won't have access until rules are added).

---

## After Deploying

Once you've applied your manifest with OpenMCF tofu, you can confirm that the Client VPN endpoint is active via the AWS console or by
using the AWS CLI:

```shell
aws ec2 describe-client-vpn-endpoints
```

You should see your new Client VPN endpoint with its endpoint ID and DNS name. You'll need to download the client configuration file
and distribute client certificates for users to connect.


