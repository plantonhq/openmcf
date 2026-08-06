# Terraform Module to Deploy AwsAlb

This module provisions an AWS Application Load Balancer and, when DNS is
enabled, Route53 alias A records for each hostname. It creates no listeners,
rules, or target groups -- routing attaches to the ALB by ARN through the
`AwsLbListener`, `AwsLbListenerRule`, and `AwsLbTargetGroup` components.

The module owns what is load-balancer-wide: placement (subnets, scheme),
security groups, and the ALB attributes (timeouts, HTTP/2, desync
mitigation, XFF handling, WAF fail-open, zonal shift, and the three S3 log
streams -- access, connection, and health-check logs). Only explicitly set
attributes are sent to AWS, so everything left unset keeps its AWS default.
Names longer than AWS's 32-character limit are truncated deterministically.

Generated `variables.tf` reflects the proto schema for `AwsAlb`.

## Usage

Use the Planton CLI (tofu) with the default local backend:

```shell
planton tofu init --manifest hack/manifest.yaml
planton tofu plan --manifest hack/manifest.yaml
planton tofu apply --manifest hack/manifest.yaml --auto-approve
planton tofu destroy --manifest hack/manifest.yaml --auto-approve
```

**Note**: Credentials are provided via stack input (CLI), not in the manifest `spec`.

For a working example, see [`hack/manifest.yaml`](../hack/manifest.yaml).
