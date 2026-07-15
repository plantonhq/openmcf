# AWS Fargate Web Service

The production container service: your application image running on ECS
Fargate behind an Application Load Balancer, health-checked, autoscaled on
real request traffic, with its own ECR repository and least-privilege IAM
roles. No instances to patch, no capacity to plan — the service scales
between the task bounds you set and bills for exactly the vCPU-seconds it
runs.

It deploys green out of the box on a public sample image, so the path from
zero to "my code is serving production traffic" is: deploy the chart, push
your image to the repository it created, point `container_image` at it,
redeploy.

## Architecture

```
 internet ──▶ AwsAlb (public subnets, AwsSecurityGroup: 80/443 from world)
                │        ▲                        ▲
                │   AwsWafWebAcl (toggle)   AwsRoute53DnsRecord A/AAAA (toggle)
                │
                ├─ AwsLbListener :443 HTTPS ──▶ forward   (https_enabled)
                ├─ AwsLbListener :80  HTTP ───▶ redirect  (https_enabled)
                └─ AwsLbListener :80  HTTP ───▶ forward   (otherwise)
                                       │
                              AwsLbTargetGroup (ip mode, health check)
                                       │ registers task IPs
              AwsEcsService ───────────┘
                │  circuit breaker + rollback
                │  autoscaling: requests-per-target (ALB + TG arn_suffix)
                ├── AwsEcsCluster (FARGATE, Container Insights toggle)
                ├── AwsEcsTaskDefinition ── AwsEcrRepo (toggle)
                │        │ awslogs → /ecs/<family> (automatic)
                │        ├── AwsIamRole (execution: pull image, write logs)
                │        └── AwsIamRole (task: the app's own AWS permissions)
                └── AwsSecurityGroup (tasks: app port from the ALB SG only)
```

Deployment order derives from the references: roles, repository, and
security groups first; then the load-balancer chain; then cluster and task
definition; the service last, joining all of them.

## Included Cloud Resources

| Resource | Kind | Purpose |
|----------|------|---------|
| Execution role | `AwsIamRole` | What the ECS agent uses before your code runs: pull image, write logs |
| Task role | `AwsIamRole` | What your code gets when it calls AWS — ships empty, grant per API |
| Image repository | `AwsEcrRepo` | Immutable-tag private registry with scan-on-push (conditional) |
| ALB security group | `AwsSecurityGroup` | 80 (and 443 with HTTPS) from the internet |
| Task security group | `AwsSecurityGroup` | App port from the load balancer only |
| Load balancer | `AwsAlb` | The internet-facing front door, spread across two-plus AZs |
| Target group | `AwsLbTargetGroup` | IP-mode target set with the service's health probe |
| Listeners | `AwsLbListener` | HTTPS + HTTP-redirect pair, or plain HTTP (toggle-shaped) |
| Web ACL | `AwsWafWebAcl` | Core + Known-Bad-Inputs managed packs on the ALB (conditional) |
| Cluster | `AwsEcsCluster` | Fargate-only scheduling home with Container Insights toggle |
| Task definition | `AwsEcsTaskDefinition` | The container blueprint: image, port, sizing, roles, logging |
| Service | `AwsEcsService` | Keeps tasks running, load-balanced, and autoscaled |
| Alias records | `AwsRoute53DnsRecord` | A + AAAA aliases to the ALB in your zone (conditional) |

## Parameters

| Parameter | Description | Default | Type |
|-----------|-------------|---------|------|
| `aws_region` | Region for every resource (cert must match it when HTTPS is on) | `us-east-1` | string |
| `service_name` | Name prefix for everything the chart creates | `web` | string |
| `vpc_id` | The hosting VPC (network-foundation's output in the common case) | placeholder | string |
| `alb_subnet_ids` | Public subnets for the ALB, two-plus AZs | placeholders | list |
| `service_subnet_ids` | Subnets tasks run in — private-with-NAT in production | placeholders | list |
| `assign_public_ip` | Public IPs on tasks (only for public subnets without NAT) | `false` | bool |
| `container_image` | The image to run; swap to your ECR image after first push | httpd sample | string |
| `container_port` | The app's listening port — drives TG, SG, and port mapping together | `80` | number |
| `health_check_path` | The readiness endpoint the ALB probes | `/` | string |
| `task_cpu` / `task_memory` | Fargate sizing pair (per-task, the main cost dial) | `512` / `1024` | number |
| `min_tasks` / `max_tasks` | Autoscaling floor (also initial count) and ceiling | `2` / `10` | number |
| `requests_per_target` | Per-task requests/minute the autoscaler holds | `1000` | number |
| `ecr_repo_enabled` | Create the private image repository | `true` | bool |
| `container_insights_enabled` | Per-service/task CloudWatch metrics (billed) | `true` | bool |
| `https_enabled` | HTTPS on 443 + HTTP→HTTPS redirect (needs `certificate_arn`) | `false` | bool |
| `certificate_arn` | ISSUED ACM cert in this region (used when HTTPS is on) | placeholder | string |
| `waf_enabled` | Managed-rules WAF in front of the ALB (~$10+/mo) | `false` | bool |
| `dns_enabled` | Alias records to the ALB in an existing zone | `false` | bool |
| `zone_id` / `domain_name` | The Route 53 zone and FQDN (used when DNS is on) | placeholders | string |

## First deploy to your own image

1. Deploy the chart with the defaults (plus your VPC/subnet ids). The
   sample image serves immediately at the ALB's DNS name
   (`AwsAlb` → `status.outputs.load_balancer_dns_name`).
2. Push your application image to the chart's repository
   (`AwsEcrRepo` → `status.outputs.repository_url`):

```bash
aws ecr get-login-password --region <region> | docker login --username AWS --password-stdin <repository_url>
docker build -t <repository_url>:v1.0.0 .
docker push <repository_url>:v1.0.0
```

3. Set `container_image` to `<repository_url>:v1.0.0` (and `container_port`
   to your app's port) and redeploy. The task definition's new revision
   rolls the service with the circuit breaker guarding the rollout — a bad
   image rolls back by itself instead of taking the service down.

Ship every release the same way: push a new immutable tag, update
`container_image`, redeploy.

## After deploying

The useful join points:

- `AwsAlb` → `status.outputs.load_balancer_dns_name` (point DNS here, or
  flip `dns_enabled`) and `status.outputs.arn_suffix` (CloudWatch
  `LoadBalancer` dimension for dashboards and alarms)
- `AwsLbTargetGroup` → `status.outputs.arn_suffix` (CloudWatch
  `TargetGroup` dimension)
- `AwsEcrRepo` → `status.outputs.repository_url` (CI pushes here)
- `AwsEcsCluster` / `AwsEcsService` → cluster and service names are the
  `ClusterName`/`ServiceName` dimensions for ECS alarms
- `AwsIamRole` (task role) → grant your app's AWS permissions on this role,
  by resource ARN, as the app grows

Swapping the literal network parameters for references: if the network came
from the network-foundation chart, any `value: subnet-...` in your manifests
can become

```yaml
- valueFrom:
    kind: AwsSubnet
    name: core-private-us-east-1a
    fieldPath: status.outputs.subnet_id
```

so the network and the service share one dependency graph.

## Day-2 guidance

- **HTTPS**: request a certificate for your domain (the
  `AwsCertManagerCert` resource with DNS validation makes issuance
  hands-free), then flip `https_enabled` with its ARN. Port 80 becomes a
  permanent redirect — no cleartext body ever crosses after that.
- **Alarms**: the service's names are stable, so ECS alarms
  (`ClusterName`/`ServiceName` dimensions) can be created any time. ALB
  alarms (5xx rate, target response time) use the `arn_suffix` outputs as
  dimensions — create them after the first deploy when those values exist.
- **Environment variables and secrets**: add them to the task definition's
  container block — plain values for configuration, Secrets Manager ARNs
  under `secrets` for credentials (never plaintext env values).
- **Blue/green and canary rollouts**: the service's
  `loadBalancers[].advancedConfiguration` composes an alternate target
  group and listener rules when rolling deployments stop being enough.
- **Spot for tolerant workloads**: add `FARGATE_SPOT` to the cluster's
  capacity providers and a `capacityProviderStrategy` on the service
  (replacing `launchType`) to blend spot capacity into non-critical
  services.

---

© Planton. Licensed under [Apache-2.0](../../../LICENSE).
