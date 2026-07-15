# AWS Route53 Zone

Deploys an Amazon Route 53 hosted zone — the DNS container for a domain. The zone's domain name comes from `metadata.name`. Public zones resolve globally (name servers exported for registrar delegation) and support reusable delegation sets, DNSSEC signing, query logging, and accelerated recovery; private zones resolve only inside their associated VPCs (split-horizon DNS). Individual DNS records are separate [AwsRoute53DnsRecord](/docs/catalog/aws/awsroute53dnsrecord) resources composing onto the zone's `zone_id` output.

## What Gets Created

When you deploy an AwsRoute53Zone resource, Planton provisions:

- **Route 53 Hosted Zone** — public or private (with its VPC association set), with comment, force-destroy behavior, optional delegation set, optional accelerated recovery, and Planton resource tags
- **DNSSEC Key-Signing Key + Signing Toggle** (optional) — a KSK backed by the referenced KMS key, with zone signing switched on after the key exists
- **Query Logging Configuration** (optional) — DNS query delivery to the referenced CloudWatch Logs log group

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **For private zones**: the VPCs to associate (DNS support + DNS hostnames enabled)
- **For DNSSEC**: a us-east-1 KMS key with key spec `ECC_NIST_P256` and a key policy allowing `dnssec-route53.amazonaws.com`
- **For query logging**: a us-east-1 CloudWatch Logs log group AND an account-level CloudWatch Logs resource policy allowing `route53.amazonaws.com` (max 10 resource policies per region — typically one shared policy covering `/aws/route53/*`)

## Quick Start

Create a file `zone.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53Zone
metadata:
  name: example.com
spec:
  region: us-east-1
```

Deploy:

```shell
planton apply -f zone.yaml
```

Then delegate the domain at your registrar using the `nameservers` output.

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `region` | `string` | Region used for provider API calls (Route 53 itself is global). |

The zone's **domain name** comes from `metadata.name` and is create-time immutable.

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `comment` | `string` | — | Note stored on the zone (max 256 chars), visible in the console. |
| `isPrivate` | `bool` | `false` | Private zone resolving only inside `vpcAssociations`. |
| `vpcAssociations` | `list(object)` | — | The private zone's VPCs (required, ≥1, when private; forbidden when public). Each entry: `vpcId` (ref to AwsVpc) + optional `vpcRegion` (defaults to the zone's region). |
| `delegationSetId` | `string` | — | Reusable delegation set ID (public zones only; create-time immutable). |
| `forceDestroy` | `bool` | `false` | Purge all records (and disable DNSSEC) before deletion. |
| `enableAcceleratedRecovery` | `bool` | `false` | Pre-stage zone data for faster recovery-event propagation (public only). |
| `queryLogging.cloudwatchLogGroupArn` | `string \| ref` | — | Destination log group ARN (us-east-1; see Prerequisites). Public only. |
| `dnssec.kmsKeyArn` | `string \| ref` | — | The KSK's KMS key (see Prerequisites). Public only. |
| `dnssec.keySigningKeyName` | `string` | derived | KSK name (3–128 chars, `0-9A-Za-z._-`). |

## Examples

### Public zone with safe teardown

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53Zone
metadata:
  name: example.com
spec:
  region: us-east-1
  comment: production apex zone
  forceDestroy: true
```

### Private split-horizon zone

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53Zone
metadata:
  name: internal.example.com
spec:
  region: us-west-2
  isPrivate: true
  vpcAssociations:
    - vpcId:
        valueFrom:
          kind: AwsVpc
          name: platform-vpc
          fieldPath: status.outputs.vpc_id
```

### DNSSEC-signed public zone

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53Zone
metadata:
  name: example.com
spec:
  region: us-east-1
  dnssec:
    kmsKeyArn:
      valueFrom:
        kind: AwsKmsKey
        name: dnssec-signing-key
        fieldPath: status.outputs.key_arn
```

### Query-logged public zone

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53Zone
metadata:
  name: example.com
spec:
  region: us-east-1
  queryLogging:
    cloudwatchLogGroupArn:
      valueFrom:
        kind: AwsCloudwatchLogGroup
        name: route53-query-logs
        fieldPath: status.outputs.log_group_arn
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `zone_id` | `string` | The hosted zone ID — the join key records, certificates, and load balancers reference. |
| `zone_name` | `string` | The domain name as normalized by Route 53. |
| `nameservers` | `list(string)` | The four authoritative name servers (registrar delegation values). |
| `primary_name_server` | `string` | The first name server (the SOA MNAME). |
| `zone_arn` | `string` | The zone's ARN, for IAM policies scoping record changes. |

## Related Components

- [AwsRoute53DnsRecord](/docs/catalog/aws/awsroute53dnsrecord) — the records that live in this zone
- [AwsRoute53HealthCheck](/docs/catalog/aws/awsroute53healthcheck) — gates records' DNS answers on endpoint health
- [AwsCertManagerCert](/docs/catalog/aws/awscertmanagercert) — DNS-validates certificates against this zone
- [AwsVpc](/docs/catalog/aws/awsvpc) — the VPCs a private zone resolves in
- [AwsKmsKey](/docs/catalog/aws/awskmskey) — backs the DNSSEC key-signing key
