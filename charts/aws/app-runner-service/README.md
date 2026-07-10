# AWS App Runner Service

The simplest production deploy on AWS: point App Runner at a container
image and get a managed HTTPS endpoint, request-based autoscaling with a
warm floor, health-checked instances, and per-request tracing — no load
balancer to assemble, no cluster to run, no TLS to terminate yourself.

The chart deploys green out of the box on AWS's public sample image. The
path to your own code: push to a private ECR repository, flip
`private_ecr_enabled`, set `image_identifier` — and with auto-deployments
on, every subsequent push of that tag rolls the service by itself.

## Architecture

```
 internet ──▶ App Runner managed endpoint (HTTPS, TLS terminated by AWS)
                │
        AwsAppRunnerService
                ├── AwsAppRunnerAutoScalingConfiguration   (always — the scaling policy)
                ├── AwsIamRole  instance role              (runtime: the app's AWS identity)
                ├── AwsIamRole  ECR access role            (deploy-time image pull; private_ecr_enabled)
                ├── AwsAppRunnerObservabilityConfiguration (X-Ray tracing; observability_enabled)
                └── AwsAppRunnerVpcConnector ── AwsSecurityGroup
                        (outbound-only path into your VPC; vpc_connector_enabled)
```

The companions are first-class, versioned AWS resources shared by
reference — a scaling or tracing change registers a new revision whose ARN
rolls every referencing service deliberately through the graph.

## Included Cloud Resources

| Resource | Kind | Purpose |
|----------|------|---------|
| Service | `AwsAppRunnerService` | The running application and its managed HTTPS endpoint |
| Scaling configuration | `AwsAppRunnerAutoScalingConfiguration` | Concurrency trigger + min/max instance bounds |
| Instance role | `AwsIamRole` | The app's runtime AWS identity — ships empty, grant per API |
| ECR access role | `AwsIamRole` | Deploy-time image pull from private ECR (conditional) |
| Observability configuration | `AwsAppRunnerObservabilityConfiguration` | X-Ray request tracing, enabled by reference (conditional) |
| VPC connector | `AwsAppRunnerVpcConnector` | Outbound path into your private subnets (conditional) |
| Connector security group | `AwsSecurityGroup` | The app's egress identity inside the VPC (conditional) |

## Parameters

| Parameter | Description | Default | Type |
|-----------|-------------|---------|------|
| `aws_region` | Region for the service and companions | `us-east-1` | string |
| `service_name` | Name prefix for everything the chart creates | `api` | string |
| `image_identifier` | The image to run (public sample by default) | hello sample | string |
| `container_port` | The app's listening port behind AWS-terminated TLS | `8080` | string |
| `health_check_path` | HTTP readiness path — unhealthy instances are replaced | `/` | string |
| `cpu` / `memory` | Per-instance sizing on App Runner's pairing grid | `1024` / `2048` | string |
| `max_concurrency` | Concurrent requests per instance before scale-out | `100` | number |
| `min_instances` | Warm floor (instant response, provisioned-rate billing) | `1` | number |
| `max_instances` | Scale-out ceiling — the worst bill you accept | `25` | number |
| `private_ecr_enabled` | Image lives in private ECR (creates the access role) | `false` | bool |
| `auto_deployments_enabled` | Push-to-deploy on the private-ECR arm | `true` | bool |
| `observability_enabled` | X-Ray per-request tracing | `true` | bool |
| `vpc_connector_enabled` | Outbound access to private VPC resources | `false` | bool |
| `vpc_id` / `connector_subnet_ids` | Where the connector attaches (used with the connector) | placeholders | string / list |
| `custom_domain_enabled` | Associate your own subdomain with the service | `false` | bool |
| `domain_name` | The subdomain to serve on (used with custom domain) | `api.example.com` | string |

## First deploy to your own image

1. Deploy with the defaults — the sample serves immediately at the
   service's URL (`AwsAppRunnerService` → `status.outputs.service_url`).
2. Create a private ECR repository (one resource, or reuse another chart's)
   and push:

```bash
aws ecr create-repository --repository-name api --region <region>
aws ecr get-login-password --region <region> | docker login --username AWS --password-stdin <account>.dkr.ecr.<region>.amazonaws.com
docker build -t <account>.dkr.ecr.<region>.amazonaws.com/api:latest .
docker push <account>.dkr.ecr.<region>.amazonaws.com/api:latest
```

3. Flip `private_ecr_enabled: true`, set `image_identifier` to the pushed
   URL, redeploy. From here, with `auto_deployments_enabled`, every push of
   that tag rolls the service — deploys become `docker push`.

## Custom domain (day-2, one-time per domain)

With `custom_domain_enabled` the service associates `domain_name`, and its
outputs then carry everything DNS needs — the values only exist once AWS
issues them, which is why the records cannot be pre-created by the chart:

- `status.outputs.custom_domains[0].dns_target` — point the domain here
- `status.outputs.custom_domains[0].certificate_validation_records[]` —
  prove domain ownership so App Runner can issue and renew the certificate

Create one CNAME per validation record plus the domain CNAME (an
`AwsRoute53DnsRecord` each), reading the values from the service's outputs:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRoute53DnsRecord
metadata:
  name: api-domain-cname
spec:
  region: us-east-1
  zoneId:
    value: <your-zone-id>
  name: api.example.com
  type: CNAME
  ttl: 300
  values:
    - <status.outputs.custom_domains[0].dns_target>
```

Leave the validation CNAMEs in place — App Runner re-uses them at every
certificate renewal. For a zone APEX (example.com), create an A record with
`aliasTarget` pointing at the DNS target instead of a CNAME (DNS forbids
apex CNAMEs).

## After deploying

- `AwsAppRunnerService` → `status.outputs.service_url` (the HTTPS
  endpoint), `status.outputs.service_arn`
- `AwsIamRole` (instance role) → grant the app's AWS permissions here, by
  resource ARN — including Secrets Manager read for anything you wire into
  `environment_secrets`
- `AwsSecurityGroup` (connector) → reference this group as the ingress
  source on database/cache security groups: "the API may reach the
  database" as one auditable rule

## Day-2 guidance

- **Private dependencies**: flip `vpc_connector_enabled` with your VPC and
  private subnets; then admit the connector's security group on each
  dependency's ingress. Egress-heavy apps keep NAT costs in mind — the
  connector routes through the subnets' route tables.
- **Secrets**: add `environmentSecrets` on the service (values are Secrets
  Manager/SSM ARNs, never material) and grant the instance role read on
  exactly those ARNs.
- **Scaling tuning**: change the scaling configuration's values — a new
  revision rolls the service. Lower `max_concurrency` for CPU-bound
  request profiles; raise `min_instances` ahead of anticipated spikes.
- **Deploy gating**: turn `auto_deployments_enabled` off to make releases
  explicit manifest changes (new tag in `image_identifier`) instead of
  registry-driven — the reviewable-deploys posture.

---

© Planton. Licensed under [Apache-2.0](../../../LICENSE).
