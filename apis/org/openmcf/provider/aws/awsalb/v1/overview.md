The **AwsAlb** component in OpenMCF offers a unified, streamlined way to provision and manage Application Load
Balancers (ALBs) on AWS. It encapsulates the core AWS ALB settings—like subnets, security groups, and optional DNS/SSL
configuration—into a single resource that fits naturally into the OpenMCF multi-cloud framework.

## Purpose and Functionality

- **Automated ALB Provisioning**: Quickly create and configure internet-facing or internal ALBs, handling subnets,
  security groups, and IP address types in one cohesive manifest.
- **Built-In DNS Management**: Optionally manage Route53 DNS records (e.g., `hostnames`, `route53ZoneId`) to map your
  ALB to custom domains without separate steps.
- **SSL Made Simple**: Enable SSL and specify a certificate ARN in one place, eliminating tedious manual setup.
- **Consistent Multi-Cloud Support**: Use the same CLI and YAML-based workflow as other OpenMCF components,
  whether you deploy to AWS or any other supported cloud.

## Key Benefits

- **Simplified Configuration**: Define all critical ALB properties under a single `spec`, with validations and defaults
  guided by OpenMCF’s Protobuf schemas.
- **Safe by Default**: Built-in fields such as `enableDeletionProtection` and `idleTimeoutSeconds` reinforce best
  practices and reduce unintentional downtime.
- **Seamless Integration**: Combine **AwsAlb** with other AWS and non-AWS components for holistic multi-cloud
  architectures, all driven by the same CLI commands.
- **Extended Observability**: Easily point domain names to your ALB, enable SSL, and integrate with your existing AWS
  monitoring setup as needed.

Below is a minimal YAML example showing how you might define an ALB (notice the **camel-case** keys):

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsAlb
metadata:
  name: exampleAlb
spec:
  subnets:
    - subnet-12345
    - subnet-67890
  securityGroups:
    - sg-1111
    - sg-2222
  scheme: internet-facing
  ipAddressType: ipv4
  enableDeletionProtection: false
  idleTimeoutSeconds: 60
  dns:
    enabled: true
    route53ZoneId: Z123456789
    hostnames:
      - alb.mydomain.com
  ssl:
    enabled: true
    certificateArn: arn:aws:acm:us-east-1:123456789012:certificate/abc-123
```

Use the **AwsAlb** resource to seamlessly incorporate AWS load balancing into your multi-cloud environment. By following
OpenMCF’s schema validations and standardized CLI workflows, teams can reduce the complexity of configuring ALBs
and ensure consistent, reliable deployments across all their services.

---

Thanks for reading! By leveraging **AwsAlb**, users can bring consistent load balancing into the same familiar process
they already use for other OpenMCF components—allowing simpler operations, better governance, and an easier path
to multi-cloud adoption.
