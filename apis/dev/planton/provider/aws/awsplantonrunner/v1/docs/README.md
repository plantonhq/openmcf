# AwsPlantonRunner: Running a Planton Runner Appliance on AWS

## Introduction

AwsPlantonRunner deploys a standing Planton runner inside an AWS VPC. The
runner is the execution arm of the Planton control plane: it receives
deploy operations (infrastructure-as-code applies and destroys) and cloud
operations (live resource browsing) and executes them against targets it
can reach. Most of the time, runners live in hosted fleets or on operator
workstations. This component exists for the cases where neither can work:
**targets reachable only from inside the customer's network**.

The canonical case is a Kubernetes cluster with a private API endpoint.
Private endpoints are the production security posture — the cluster's
control plane is simply not addressable from the internet — but that same
property means no external system can deploy into the cluster. The
industry answer is an in-network agent that dials OUT to its coordinator
and executes work locally. GitHub's self-hosted Actions runners, GitLab
runners, Buildkite agents, and Atlantis all follow this shape. The Planton
runner follows it too, and this component packages it as a first-class,
declaratively-managed appliance rather than a hand-run install script.

Two properties define the design. First, **outbound-only networking**: the
runner initiates every connection it uses — to the control plane, to its
work queue, to the registry it pulls its image from — so its security
group needs no inbound rule at all, and running it adds zero attack
surface to the VPC. Second, **standing, not ephemeral**: the appliance
outlives the clusters it deploys to. That is what makes lifecycle
ordering tractable — in-cluster workloads are destroyed through the
runner, the cluster is destroyed over the AWS path, and the runner is
destroyed last. An ephemeral "bootstrap runner" that deleted itself after
the first install would leave day-2 operations (upgrades, destroys,
incident response) with no execution path back into the network.

## The Landscape: How Teams Run In-Network Agents Today

### Level 0: Manual (a VM and a systemd unit)

The traditional approach: launch an EC2 instance in the VPC, SSH in,
download the agent binary, write a systemd unit, and paste credentials
into an environment file.

```shell
ssh ec2-user@10.0.1.23
curl -Lo /usr/local/bin/agent https://example.com/agent
sudo systemctl enable --now agent
```

**Pros:** conceptually simple; full control over the host.

**Cons:** a pet server — you patch its OS forever; credentials sit in a
plaintext file on disk; nothing restarts or replaces it when the
instance degrades; the setup is invisible to code review and
irreproducible six months later.

**Verdict:** how most in-network agents actually die: forgotten,
unpatched, and holding credentials nobody remembers issuing.

### Level 1: CLI provisioning (scripted AWS calls)

A script drives the AWS CLI or SDK: create a cluster, register a task
definition, create a service.

```shell
aws ecs create-cluster --cluster-name runner
aws ecs register-task-definition --family runner --cli-input-json file://taskdef.json
aws ecs create-service --cluster runner --service-name runner \
  --task-definition runner --desired-count 1 --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-...],securityGroups=[sg-...]}"
```

**Pros:** repeatable-ish; no console clicking; serverless compute removes
the OS-patching burden.

**Cons:** imperative — the script creates but does not reconcile; the IAM
roles, log group, and secret handling are each another dozen lines that
every team writes slightly differently; teardown is a second script that
drifts from the first.

**Verdict:** better than a pet VM, but the operational knowledge lives in
a script, not in a declarative model anything can reason about.

### Level 2: Terraform

The infrastructure-as-code answer: `aws_ecs_cluster`,
`aws_ecs_task_definition`, `aws_ecs_service`, two `aws_iam_role`s, an
`aws_cloudwatch_log_group`, an `aws_security_group`, and an
`aws_secretsmanager_secret`, wired together by hand.

```hcl
resource "aws_ecs_task_definition" "runner" {
  family                   = "runner"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "512"
  memory                   = "1024"
  execution_role_arn       = aws_iam_role.execution.arn
  container_definitions    = jsonencode([{
    name      = "runner"
    image     = "ghcr.io/plantonhq/planton/runner:latest"
    essential = true
    secrets   = [{ name = "PLANTON_RUNNER_CREDENTIALS", valueFrom = aws_secretsmanager_secret.creds.arn }]
    # ... log configuration, ports, environment ...
  }])
}
```

**Pros:** declarative, reviewable, reconciled; state tracks every piece;
destroy is symmetric.

**Cons:** ~200 lines of HCL to get right, and the sharp edges are all
invisible until runtime — the execution role missing
`secretsmanager:GetSecretValue` on exactly the right ARN, the awslogs
group not existing before the first task starts, the Fargate cpu/memory
pairing rejected at apply time, the secret's 30-day recovery window
blocking a re-create with the same name.

**Verdict:** the right foundation, which is why this component's own
modules are built on it — but as a per-team exercise it re-derives the
same 200 lines and re-discovers the same pitfalls every time.

### Level 3: Pulumi

The same resource set in a general-purpose language:

```go
service, err := ecs.NewService(ctx, "runner", &ecs.ServiceArgs{
    Cluster:        cluster.Arn,
    TaskDefinition: taskDef.Arn,
    DesiredCount:   pulumi.Int(1),
    LaunchType:     pulumi.String("FARGATE"),
    NetworkConfiguration: &ecs.ServiceNetworkConfigurationArgs{
        Subnets:        pulumi.StringArray{...},
        SecurityGroups: pulumi.StringArray{sg.ID()},
    },
})
```

**Pros:** everything Terraform offers, plus real-language abstraction and
testing.

**Cons:** the same domain knowledge burden; loops and helper functions
make each team's variant more, not less, unique.

**Verdict:** equivalent power, same re-derivation cost.

### Other approaches

- **Kubernetes-hosted agents (Helm):** run the runner in a cluster. The
  right pattern for customer-managed self-hosted runners on EXISTING
  clusters — but structurally wrong as the appliance that manages the
  cluster's own lifecycle: an in-cluster runner uninstalling its own
  deployment kills itself mid-operation, and it dies with the very
  cluster it would be needed to rebuild.
- **EC2 with an auto-scaling group of one:** heavier than Fargate for a
  single always-on container and reintroduces AMI/OS lifecycle.
- **AWS App Runner / Lambda:** built for request-serving and event
  handling respectively; a long-polling worker executing hour-scale IaC
  operations fits neither runtime model.

## Comparative Analysis

| Method | Declarative | OS burden | Secret handling | Day-2 story | Effort |
| --- | --- | --- | --- | --- | --- |
| Manual VM | No | Yours, forever | Plaintext on disk | None | Low once, high forever |
| CLI scripts | No | None (Fargate) | Varies | Re-run scripts | Medium |
| Terraform (hand-rolled) | Yes | None | Yours to design | Reconciled | High (once per team) |
| Pulumi (hand-rolled) | Yes | None | Yours to design | Reconciled | High (once per team) |
| **AwsPlantonRunner** | **Yes** | **None** | **Secrets Manager, native injection** | **Reconciled + composable** | **One manifest** |

## The Planton Approach

The component models the runner's **intent** and derives everything else:

- **Placement** — `subnets` (references or literals) put the runner in
  the network whose private endpoints it must reach. The VPC is derived
  from the subnets rather than asked for separately: a second field could
  only ever agree with the subnets or contradict them.
- **Sizing** — Fargate task `cpu`/`memory`, defaulting to 512/1024, with
  the valid pairing validated at manifest time. AWS itself rejects an
  invalid combination only at deploy time, deep into an apply; the spec
  turns that into an instant, explainable validation error.
- **Version** — `runner_version` is the image tag of the official runner
  image; `image_repository` exists only for air-gapped mirrors. The
  component runs Planton's own container, so "which build" is the real
  intent — not an arbitrary image string.
- **Execution mode** — `temporal` (pull-based deploy worker, the
  private-endpoint default), `dual` (adds the real-time CloudOps channel
  through an outbound-initiated tunnel), or `grpc` (CloudOps only).
- **Identity** — `credentials` is the registration's credentials
  document, handled as a secret end to end; `task_role` is the runner's
  own AWS identity at runtime, composed as a first-class `AwsIamRole`
  reference when the runner needs permissions of its own.

Just as deliberately, the spec does **not** model:

- **The compute substrate.** No ECS cluster reference, no launch-type
  knob, no capacity-provider blend. The appliance is the product; ECS
  Fargate is the implementation. Keeping the substrate out of the API
  means the contract survives implementation evolution untouched.
- **Connectivity endpoints.** The credentials document is the single
  source of the runner's identity AND connectivity; asking a manifest
  author for coordinates the platform already knows would be setup
  ceremony with a failure mode (a wrong address) attached.
- **A replica count.** A runner registration corresponds to exactly one
  running agent; its work queue serializes operations. Scaling execution
  capacity is a matter of registering more runners, not scaling one
  service — so a `desired_count` field could only introduce a
  misconfiguration class.
- **Spot capacity.** An appliance that must be present to receive work is
  the textbook wrong fit for interruptible capacity; the cost difference
  on a 0.5-vCPU service is noise.

### Design decisions worth recording

- **Credentials ride Secrets Manager, never the task definition.** The
  container definition carries only the secret's ARN; the ECS agent
  fetches the value at task start using the execution role, whose inline
  policy grants `secretsmanager:GetSecretValue` on exactly that one ARN.
  Anyone who can read the task definition (a common, low-sensitivity
  permission) learns nothing.
- **The secret uses a zero-day recovery window.** It holds re-mintable
  credential material, not data. Secrets Manager's default 30-day
  soft-delete would block re-creating a same-named runner for a month
  after a destroy — an operational trap with no compensating benefit
  here.
- **Two IAM roles, strictly separated.** The execution role is the
  *setup* identity (pull image, write logs, read the one secret); the
  task role is the *runtime* identity the runner holds while executing
  work. Collapsing them would hand the runner's workloads the
  infrastructure permissions and vice versa.
- **A dedicated ECS cluster per runner.** Clusters are free scheduling
  boundaries (no cost until tasks run). A dedicated one keeps the
  appliance self-contained: its teardown removes everything, and no
  shared-cluster coupling can ever block or complicate it.
- **The service never gates on task health.** ECS reports a service
  ACTIVE independently of its tasks, and the modules deliberately do not
  wait for steady state or enable the deployment circuit breaker: a
  runner whose control plane is momentarily unreachable must still
  deploy (and destroy) cleanly. The runner's real readiness contract is
  its work queue — operations wait there until the worker polls, so
  nothing downstream depends on ECS-level timing.
- **Log retention is explicit.** The runner's logs are the audit trail of
  every change it executed; the log group is created with a configured
  retention (default 30 days) instead of CloudWatch's never-expire
  default.

## Implementation Landscape

Both engines provision the same eight resources in the same dependency
order:

1. **Secrets Manager secret** (`<name>-credentials`) — the credentials
   document, zero-day recovery window.
2. **Execution role** — trust `ecs-tasks.amazonaws.com`; managed
   `AmazonECSTaskExecutionRolePolicy` plus an inline read grant scoped to
   the one secret.
3. **Runtime (task) role** — created permissionless, or skipped entirely
   when `task_role` references an existing `AwsIamRole`.
4. **CloudWatch log group** (`/ecs/<name>`) — explicit retention.
5. **Security group** — outbound-only (the VPC default egress rule is
   exactly the posture needed; no inbound rules), in the VPC looked up
   from the first subnet.
6. **ECS cluster** — dedicated to this runner.
7. **Task definition** — Fargate, awsvpc networking, the runner container
   with its command, ports, environment (execution mode, tunnel toggle,
   log level), the secret injection, and awslogs wiring. Task
   definitions are immutable; every configuration change registers a new
   revision, and the service rolls to it.
8. **Fargate service** — desired count 1, no public IP unless asked for,
   the created security group plus any referenced extras.

The Pulumi module lives at `iac/pulumi/` (Go, one file per resource
concern) and the OpenTofu module at `iac/tf/`; both flatten their outputs
onto the same `AwsPlantonRunnerStackOutputs` contract and are held to
behavioral parity by the outputs-conformance guard.

## Production Best Practices

- **Prefer private subnets with a NAT route.** The runner needs outbound
  internet (image pulls, control plane). A NAT gateway keeps it
  unaddressable from outside; `assign_public_ip` exists for NAT-less
  public-subnet setups but puts a public address on the appliance.
- **Spread across two AZs.** Two subnets in different availability zones
  let the service reschedule the runner across an AZ event.
- **Let targets trust the security group, not IPs.** The
  `security_group_id` output is the stable handle: a private cluster API
  endpoint or database admits the runner by referencing it, and the rule
  survives every task replacement (Fargate task IPs churn).
- **Pin `runner_version` in production.** `latest` tracks releases and is
  right for evaluation; pinned tags plus an explicit bump give change
  control and a clean rollback.
- **Grant the runtime role deliberately.** Start permissionless; add
  exactly what keyless operations need, on your own `AwsIamRole`, and
  reference it. The role ARN in `task_role_arn` is the single place to
  audit what the runner can do in the account.
- **Rotate credentials by re-issuing the registration document** and
  updating the managed secret; the next task restart picks it up. The
  document is re-mintable — treat rotation as routine, not exceptional.
- **Watch the logs, not the task count.** The log group is the
  operational signal: operation starts, failures, and control-plane
  connectivity all land there. `aws logs tail <log_group_name> --follow`
  during any incident.
- **Size memory before CPU.** IaC operations are memory-hungry (provider
  plugins, plans held in memory); a failed apply from memory pressure
  looks like a runner fault but is a sizing fault.

## Cost

The default 0.5 vCPU / 1 GiB Fargate task runs continuously:
approximately $18/month in us-east-1 (on-demand, x86). The ECS cluster,
security group, IAM roles, and task definition are free; Secrets Manager
adds $0.40/month per secret; CloudWatch logs bill by ingestion and
retention (the runner's steady-state log volume is small). A NAT gateway,
when one is added for private subnets, dominates this bill — but it is
network infrastructure the VPC usually has anyway.

## Conclusion

AwsPlantonRunner turns "we need something inside the network that can
deploy" from a hand-run install script or 200 lines of bespoke IaC into
one declarative resource that composes with the rest of the graph: subnets
and security groups and IAM roles in, a trusted security-group handle and
an auditable log trail out. Use it whenever a deployment target is
reachable only from inside an AWS network — private-endpoint Kubernetes
clusters above all — and keep it standing: it is the network's operability,
not a bootstrap step.

## References

- ECS Fargate task sizing: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html
- ECS secret injection: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/specifying-sensitive-data-secrets.html
- ECS task execution role: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_execution_IAM_role.html
- Fargate networking: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/fargate-task-networking.html
- CloudWatch Logs retention: https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/Working-with-log-groups-and-streams.html
