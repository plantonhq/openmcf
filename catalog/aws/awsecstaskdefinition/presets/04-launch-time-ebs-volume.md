# Launch-Time EBS Volume Task

This preset declares a task whose volume is configured at launch time: the task definition names the volume slot, and the ECS service that runs it attaches a fresh, service-owned EBS volume per task at deployment time. Real block storage (gp3/io1 IOPS, snapshots, encryption) beyond the ephemeral scratch space -- without baking storage details into the task blueprint.

## When to Use

- Stateful containers needing real block storage per task (databases, queues, build caches restored from snapshots)
- Workloads whose storage sizing/type should be a DEPLOYMENT decision, not a blueprint decision
- Pairing with the `AwsEcsService` `04-stateful-managed-ebs` preset -- the two halves of one contract

## Key Configuration Choices

- **`configureAtLaunch: true` with no backing** -- The volume declares only its NAME here; declaring an EFS/host/Docker/S3 backing alongside it is rejected at validation
- **The name is the join key** -- The consuming service's `volumeConfiguration.name` must match `data` exactly; ECS pairs them at deployment time
- **Mount point as usual** (`mountPoints`) -- Containers mount launch-time volumes exactly like statically-backed ones

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region the task definition is registered in | AWS region list |
| `<execution-role-arn>` | ARN of the task execution role (needs `AmazonECSTaskExecutionRolePolicy`) | `AwsIamRole` status outputs |
| `<container-image>` | The application container image | Your registry |

## Related Presets

- **01-web-app** -- Stateless baseline without volumes
- **AwsEcsService / 04-stateful-managed-ebs** -- The consuming service that supplies this volume's EBS backing
