# Stateful Service with Managed EBS

This preset runs a stateful service whose every task gets a fresh, service-owned EBS volume attached at deployment time. The task definition declares the volume slot (`configureAtLaunch`); this service supplies the actual storage -- size, type, IOPS, encryption, tags -- so storage is a deployment decision that can change per environment without re-registering the task blueprint.

## When to Use

- Stateful containers needing real block storage per task (embedded databases, queues, caches hydrated from snapshots)
- Different storage shapes per environment (50 GiB gp3 in staging, io2 with provisioned IOPS in production) over one task definition
- Pairing with the `AwsEcsTaskDefinition` `04-launch-time-ebs-volume` preset -- the two halves of one contract

## Key Configuration Choices

- **The name is the join key** (`volumeConfiguration.name: data`) -- Must match the task definition's `configureAtLaunch` volume name exactly; ECS pairs them at deployment time
- **EBS infrastructure role** (`managedEbsVolume.roleArn`) -- ECS assumes it to create and attach the per-task volumes; it must trust `ecs.amazonaws.com` and carry `AmazonECSInfrastructureRolePolicyForVolumes`
- **Volume tags** (`tagSpecifications`) -- Without these, the per-task volumes carry no cost-allocation tags at all; `propagateTags: SERVICE` inherits the service's tags onto every created volume
- **Circuit breaker on** -- A deployment whose tasks cannot mount their volumes rolls back instead of hanging

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region where the service runs | AWS region list |
| `<cluster-name>` | The `AwsEcsCluster` resource name | Your cluster resource |
| `<task-definition-name>` | The `AwsEcsTaskDefinition` declaring the `data` launch-time volume | Your task definition resource |
| `<private-subnet-1>` / `<private-subnet-2>` | `AwsSubnet` resource names for the task ENIs | Your network resources |
| `<ebs-infrastructure-role-arn>` | ARN of the ECS volumes infrastructure role | `AwsIamRole` status outputs |

## Related Presets

- **01-web-service** -- Stateless baseline without volumes
- **AwsEcsTaskDefinition / 04-launch-time-ebs-volume** -- The task definition declaring the volume slot this service fills
