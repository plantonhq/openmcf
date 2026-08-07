# Pulumi Module to Deploy AliCloudEipAddress

This module provisions an Alibaba Cloud Elastic IP Address (EIP) using the
`ecs.EipAddress` Pulumi resource. The EIP is a standalone public IPv4 address
that can be associated with NAT gateways, load balancers, VPN gateways, and
ECS instances.

## CLI Usage (Planton Pulumi)

```shell
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

**Note**: Credentials are provided via stack input (CLI), not in the manifest `spec`.

For more examples, see `../../e2e/manifest.yaml` and [`e2e/manifest.yaml`](../../e2e/manifest.yaml).
