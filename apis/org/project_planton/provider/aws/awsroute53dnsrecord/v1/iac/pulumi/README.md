# AWS Route53 DNS Record - Pulumi Module

This Pulumi module creates and manages AWS Route53 DNS records.

## Overview

The module provisions a single DNS record in an existing Route53 hosted zone with support for:

- Standard DNS record types (A, AAAA, CNAME, MX, TXT, etc.)
- Alias records (Route53's killer feature for AWS resources)
- Advanced routing policies (weighted, latency, failover, geolocation)
- Health check integration

## Usage

### With Project Planton CLI

```bash
# Deploy a DNS record
project-planton pulumi up --manifest dns-record.yaml

# Preview changes
project-planton pulumi preview --manifest dns-record.yaml

# Destroy the record
project-planton pulumi destroy --manifest dns-record.yaml
```

### Standalone Usage

```bash
# Set the stack input as environment variable
export STACK_INPUT=$(cat manifest.yaml | base64)

# Initialize and run
pulumi login --local
pulumi stack init dev
pulumi up
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `STACK_INPUT` | Base64-encoded manifest YAML | Yes |
| `AWS_ACCESS_KEY_ID` | AWS access key (if not in manifest) | No |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key (if not in manifest) | No |
| `AWS_REGION` | AWS region (if not in manifest) | No |

## Outputs

| Output | Description |
|--------|-------------|
| `fqdn` | Fully qualified domain name of the record |
| `record_type` | DNS record type (A, AAAA, CNAME, etc.) |
| `hosted_zone_id` | Route53 hosted zone ID |
| `is_alias` | Whether this is an alias record |
| `set_identifier` | Set identifier for routing policies |

## Example Manifests

### Basic A Record

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: www-example
spec:
  hosted_zone_id: Z1234567890ABC
  name: www.example.com
  type: A
  ttl: 300
  values:
    - 192.0.2.1
```

### Alias Record to CloudFront

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRoute53DnsRecord
metadata:
  name: apex-cloudfront
spec:
  hosted_zone_id: Z1234567890ABC
  name: example.com
  type: A
  alias_target:
    dns_name: d1234abcd.cloudfront.net
    hosted_zone_id: Z2FDTNDATAQYW2
```

## Troubleshooting

### Record already exists

If you get an error that the record already exists, the record may have been created outside of Pulumi. You can import it:

```bash
pulumi import aws:route53/record:Record www Z1234567890ABC_www.example.com_A
```

### Invalid hosted zone ID

Ensure you're using the correct hosted zone ID:
- For your records: Use your Route53 hosted zone ID
- For alias targets: Use the AWS service's hosted zone ID (e.g., CloudFront: Z2FDTNDATAQYW2)

### TTL ignored for alias records

This is expected behavior. Alias records use the TTL of the target resource automatically.
