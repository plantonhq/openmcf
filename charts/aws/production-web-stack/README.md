# AWS Production Web Stack

Startup-in-a-box: the complete production architecture for a web
application in one deploy. Its own two-tier VPC network, a load-balanced
and request-count-autoscaled Fargate service with an ECR registry and
least-privilege roles, Aurora PostgreSQL Serverless v2 with an AWS-managed
master password, HA Redis, HTTPS on your own domain, WAF, and CloudWatch
alarms wired to email. Twenty-six resources that would take days to
assemble correctly, composed as one dependency graph — and every piece a
first-class resource you can evolve independently afterwards.

This is the deliberate kitchen-sink of the AWS catalog. Every other chart
stays focused on one architectural intent; this one exists so a team can
go from an empty AWS account to a defensible production posture in a
single afternoon.

## Architecture

```
 internet ──▶ Route 53 A/AAAA aliases (custom_domain_enabled)
                │
        AwsAlb (public tier, delete-protected) ◀── AwsWafWebAcl (waf_enabled)
                ├─ :443 HTTPS + :80 redirect ◀── AwsCertManagerCert (custom domain)
                └─ :80 HTTP forward (before a domain exists)
                │
        AwsLbTargetGroup (ip mode, health-checked)
                │
        AwsEcsService (circuit breaker, requests-per-target autoscaling)
                ├── AwsEcsCluster (Fargate, Container Insights)
                ├── AwsEcsTaskDefinition ── AwsEcrRepo
                └── IAM: execution role + task role (reads the DB secret)
                │
   private tier │ (no public IPs; egress via one shared NAT)
        AwsRdsCluster (Aurora PG Serverless v2, writer + reader,
                │      managed master password → Secrets Manager)
        AwsRedisElasticache (2-node HA, TLS; redis_enabled)
                │
        AwsVpc + public/private AwsSubnet per AZ + IGW + EIP/NAT
        Security groups: world→ALB→tasks→{db, redis}; deny-all DB/cache egress

        AwsSnsTopic + email + 4 AwsCloudwatchAlarms (alarms_enabled):
        service CPU/memory, database CPU, Aurora capacity
```

## Included Cloud Resources

| Resource | Kind | Purpose |
|----------|------|---------|
| Network | `AwsVpc`, `AwsSubnet` ×4, `AwsInternetGateway`, `AwsElasticIp`, `AwsNatGateway` | Two-tier network striped across the AZs |
| Security groups | `AwsSecurityGroup` ×3-4 | The reachability contract: world→ALB→tasks→data |
| IAM | `AwsIamRole` ×2 | Execution role; task role with DB-secret read |
| Registry | `AwsEcrRepo` | Immutable-tag images, scan-on-push |
| Front door | `AwsAlb`, `AwsLbTargetGroup`, `AwsLbListener` ×1-2 | Load balancing, health gating, TLS/redirect |
| WAF | `AwsWafWebAcl` | Core + Known-Bad-Inputs managed packs (conditional) |
| Compute | `AwsEcsCluster`, `AwsEcsTaskDefinition`, `AwsEcsService` | The autoscaled application |
| Database | `AwsRdsCluster` | Aurora PG Serverless v2, writer + reader, managed password |
| Cache | `AwsRedisElasticache` | 2-node HA Redis, encrypted (conditional) |
| Domain | `AwsRoute53Zone`, `AwsCertManagerCert`, `AwsRoute53DnsRecord` ×2 | Zone (optional), cert, dual-stack aliases (conditional) |
| Alerting | `AwsSnsTopic`, `AwsSnsSubscription`, `AwsCloudwatchAlarm` ×4 | The pager wire (conditional) |

## Parameters

| Parameter | Description | Default | Type |
|-----------|-------------|---------|------|
| `aws_region` / `stack_name` | Region and the name prefix for everything | `us-east-1` / `prod` | string |
| `vpc_cidr` | Immutable VPC range — never overlap future peers | `10.0.0.0/16` | string |
| `availability_zones` | AZs to stripe (CIDR lists are index-matched) | 2 AZs | list |
| `public_subnet_cidrs` / `private_subnet_cidrs` | One entry per AZ | /20 slices | list |
| `container_image` | The app image (public sample by default) | httpd | string |
| `container_port` / `health_check_path` | The app's port and readiness probe | `80` / `/` | number / string |
| `task_cpu` / `task_memory` | Fargate sizing pair | `512` / `1024` | number |
| `min_tasks` / `max_tasks` / `requests_per_target` | Autoscaling bounds + per-task target | `2` / `10` / `1000` | number |
| `db_name` / `db_master_username` | Initial database and master user (password is AWS-managed) | `appdb` / `dbadmin` | string |
| `db_min_capacity` / `db_max_capacity` | Serverless v2 ACU bounds (0 floor = auto-pause for staging) | `0.5` / `16` | number |
| `deletion_protection` | Refuse to delete the database while on | `true` | bool |
| `redis_enabled` / `redis_node_type` | HA cache tier and its node class | `true` / `cache.t4g.small` | bool / string |
| `custom_domain_enabled` / `domain_name` | HTTPS on your own domain | `false` / placeholder | bool / string |
| `dns_zone_enabled` / `dns_zone_name` / `existing_zone_id` | Create the zone, or bring your own | `true` / placeholders | bool / string |
| `waf_enabled` | Managed-rules WAF on the ALB | `true` | bool |
| `alarms_enabled` / `alert_email` | The alarm set and where it pages | `true` / placeholder | bool / string |

## First deploy

1. Deploy with the defaults (adjust `stack_name`, region, and
   `alert_email`). The sample app serves at the ALB's DNS name
   (`AwsAlb` → `status.outputs.load_balancer_dns_name`); the database and
   cache come up on the private tier.
2. **Click the SNS confirmation link** AWS emails to `alert_email` — no
   alert is delivered until it is confirmed.
3. Push your image to the stack's registry
   (`AwsEcrRepo` → `status.outputs.repository_url`), set
   `container_image` (and `container_port`), redeploy. The circuit breaker
   guards the rollout; a bad image rolls back by itself.
4. When you have a domain: flip `custom_domain_enabled` (with
   `dns_zone_enabled: false` + `existing_zone_id` if the zone already
   exists). Certificate issuance is hands-free through DNS validation, 80
   becomes a permanent redirect, and the aliases go live. If the chart
   created the zone, delegate your registrar's NS records to it
   (`AwsRoute53Zone` → `status.outputs.nameservers`).

## Wiring the app to its data tier

The endpoints and credentials the application needs, all outputs:

- `AwsRdsCluster` → `status.outputs.endpoint` (writer),
  `status.outputs.reader_endpoint` (reads),
  `status.outputs.master_user_secret_arn` (the credentials secret — the
  task role can already read it; fetch at startup, or wire it into the
  container's `secrets` by ARN)
- `AwsRedisElasticache` → `status.outputs.primary_endpoint_address`
  (connect with `rediss://` — transit encryption is on)

Add these as environment variables/secrets on the task definition's
container block as your app expects them.

## Day-2 guidance

- **ALB alarms**: 5xx rate and target latency use the deploy-time
  `arn_suffix` outputs as CloudWatch dimensions, so they cannot be
  pre-rendered honestly. After the first deploy, add two
  `AwsCloudwatchAlarm` resources (namespace `AWS/ApplicationELB`, metrics
  `HTTPCode_Target_5XX_Count` and `TargetResponseTime`, dimensions
  `LoadBalancer`/`TargetGroup` from the two `arn_suffix` outputs) pointed
  at the existing alerts topic.
- **Tighten the DB-secret grant**: the task role's Secrets Manager grant
  is name-pattern-scoped (`rds!*`) because the secret's suffixed ARN only
  exists after creation; replace it with the exact
  `master_user_secret_arn` value.
- **Per-AZ NAT**: this stack runs one shared NAT gateway (the cost-sane
  default). For per-AZ egress survivability, stand up the
  network-foundation chart (its `nat_gateway_per_az` dial) and compose the
  application pieces onto it.
- **Staging twin**: a second instance with `stack_name: staging`,
  `db_min_capacity: 0` (auto-pause), `redis_enabled: false`,
  `waf_enabled: false`, `min_tasks: 1` is the same architecture at a
  fraction of the cost.
- **CI/CD**: the ci-cd-pipeline chart pairs directly — point its ECS
  deploy stage at `<stack_name>-cluster` / `<stack_name>-web` and its
  buildspec at this stack's repository.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
