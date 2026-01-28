```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsClientVpn
metadata:
  name: example
spec: {}
```

CLI:

```bash
openmcf pulumi preview \
  --manifest ../hack/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir .

openmcf pulumi update \
  --manifest ../hack/manifest.yaml \
  --stack organization/<project>/<stack> \
  --module-dir . \
  --yes
```

# Examples

## Minimal manifest (YAML)
```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsClientVpn
metadata:
  name: my-client-vpn
spec:
  vpcId:
    value: vpc-12345678
  subnets:
    - value: subnet-abc123
  clientCidrBlock: 10.0.0.0/22
  authenticationType: certificate
  serverCertificateArn:
    value: arn:aws:acm:us-east-1:123456789012:certificate/abc
  vpnPort: 443
  transportProtocol: tcp
  cidrAuthorizationRules:
    - 10.0.0.0/16
```

## CLI flows
- Validate:
```bash
openmcf validate --manifest ./manifest.yaml
```

- Pulumi preview/update via CLI:
```bash
openmcf pulumi preview --manifest ./manifest.yaml --stack <org>/<project>/<stack> --module-dir .
openmcf pulumi update --manifest ./manifest.yaml --stack <org>/<project>/<stack> --module-dir . --yes
```

Note: Provider credentials are supplied via stack input, not in the spec.


