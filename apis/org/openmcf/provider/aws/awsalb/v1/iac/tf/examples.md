# AWS ALB Examples

Below are several examples demonstrating how to define an AWS ALB (Application Load Balancer) component in
OpenMCF. After creating one of these YAML manifests, apply it with Terraform using the OpenMCF CLI:

```shell
openmcf tofu apply --manifest <yaml-path> --stack <stack-name>
```

---

## Basic Internet-Facing ALB

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsAlb
metadata:
  name: basic-alb
spec:
  subnets:
    - subnet-12345abc
    - subnet-67890def
  securityGroups:
    - sg-abcdef1234567890
```

This example:
• Creates an ALB in two subnets (`subnet-12345abc` and `subnet-67890def`).
• Attaches one security group (`sg-abcdef1234567890`).
• Uses default HTTP listener on port 80 with auto-redirect to HTTPS when SSL is enabled.

---

## ALB with HTTPS

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsAlb
metadata:
  name: alb-with-https
spec:
  subnets:
    - subnet-1111aaaa
    - subnet-2222bbbb
  securityGroups:
    - sg-zzzzxxxxccccdddd
  ssl:
    enabled: true
    certificateArn: "arn:aws:acm:us-east-1:123456789012:certificate/abcd1234-5678-efgh-ijkl-123456abcdef"
```

Here:
• SSL is enabled with a certificate ARN for TLS termination.
• HTTP (port 80) automatically redirects to HTTPS (port 443).
• Uses the default SSL policy `ELBSecurityPolicy-2016-08`.

---

## Internal ALB

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsAlb
metadata:
  name: internal-alb
spec:
  subnets:
    - subnet-internal1
    - subnet-internal2
  securityGroups:
    - sg-internal1234
  internal: true
```

This ALB is internal-facing:
• `internal: true` makes it private/internal traffic only.
• It can serve traffic only accessible within the specified subnets or VPC.

---

## ALB with DNS Management

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsAlb
metadata:
  name: alb-with-dns
spec:
  subnets:
    - subnet-abc12345
    - subnet-def67890
  securityGroups:
    - sg-1234abcd5678efgh
  dns:
    enabled: true
    hostnames:
      - "app.example.com"
      - "api.example.com"
    route53ZoneId: "Z1234567890ABCDEF"
```

In this example:
• DNS management is enabled for Route53 .
• Two hostnames point to the ALB.
• Route53 zone ID is specified for DNS record creation.

---

## ALB with Custom Idle Timeout and Deletion Protection

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsAlb
metadata:
  name: alb-custom-config
spec:
  subnets:
    - subnet-custom1
    - subnet-custom2
  securityGroups:
    - sg-custom
  deleteProtectionEnabled: true
  idleTimeoutSeconds: 120
  ssl:
    enabled: true
    certificateArn: "arn:aws:acm:us-east-1:987654321098:certificate/fedcba98-7654-3210-fedc-ba9876543210"
```

This configuration:
• Sets a custom idle timeout of 120 seconds.
• Protects the ALB from accidental deletion by enabling deletion protection.
• Enables SSL with a certificate ARN.

---

## Complete ALB with All Features

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsAlb
metadata:
  name: production-alb
spec:
  subnets:
    - subnet-prod-public-1a
    - subnet-prod-public-1b
  securityGroups:
    - sg-alb-production
  internal: false
  deleteProtectionEnabled: true
  idleTimeoutSeconds: 60
  dns:
    enabled: true
    hostnames:
      - "app.production.example.com"
    route53ZoneId: "Z1234567890ABCDEF"
  ssl:
    enabled: true
    certificateArn: "arn:aws:acm:us-east-1:123456789012:certificate/prod-cert-1234"
```

A production-ready ALB with:
• Internet-facing configuration in public subnets.
• DNS management for custom domain.
• SSL termination with certificate.
• Deletion protection enabled.
• Standard idle timeout.

---

## After Deploying

Once you've applied your manifest with OpenMCF tofu, you can confirm that the ALB is active via the AWS console or by
using the AWS CLI:

```shell
aws elbv2 describe-load-balancers --names <your-load-balancer-name>
```

You should see details about your Application Load Balancer, such as its DNS name and ARN. Feel free to adjust subnets,
security groups, or other settings to fit your needs.
