# Pulumi Module to Deploy AwsAlb

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

It exports the same outputs as the Terraform module: `load_balancer_arn`,
`load_balancer_name`, `load_balancer_dns_name`, and
`load_balancer_hosted_zone_id`.

## CLI usage (Planton pulumi)

```bash
# Preview
planton pulumi preview \
  --manifest ../../e2e/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir .

# Update (apply)
planton pulumi update \
  --manifest ../../e2e/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir . \
  --yes

# Refresh
planton pulumi refresh \
  --manifest ../../e2e/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir .

# Destroy
planton pulumi destroy \
  --manifest ../../e2e/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir .
```

## Debugging

This module includes a `debug.sh` helper. To enable debugging, edit `Pulumi.yaml` and uncomment the `runtime.options.binary` line so Pulumi runs the program via the script:

```yaml
name: aws-module-test-pulumi-project
runtime:
  name: go
#  options:
#    binary: ./debug.sh
```

Then make the script executable and run your command (e.g., `preview` or `update`). See `docs/pages/docs/guide/debug-pulumi-modules.mdx` for full instructions.

```bash
chmod +x debug.sh
planton pulumi preview \
  --manifest ../../e2e/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir .
```
