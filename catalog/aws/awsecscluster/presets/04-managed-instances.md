# Managed Instances Cluster

This preset adds ECS Managed Instances capacity: ECS launches, patches, and retires the EC2 instances itself from attribute-based requirements -- no auto-scaling group, no AMI selection, no user data. You describe the compute (memory, vCPUs, CPU manufacturer) and the network to launch into; ECS owns the fleet end to end. EC2 pricing with Fargate-grade operational hands-off.

## When to Use

- EC2 unit economics (or instance features Fargate lacks) without owning an auto-scaling group and AMI pipeline
- Fleets best described by requirements ("2 GiB+, Graviton") instead of a named instance type
- Teams consolidating on capacity providers who want AWS to handle instance patching and lifecycle

## Key Configuration Choices

- **Infrastructure role** (`infrastructureRoleArn`) -- The role ECS assumes to launch and retire instances; it must trust `ecs.amazonaws.com` and carry `AmazonECSInfrastructureRolePolicyForManagedInstances`. The identity applying the manifest needs `iam:PassRole` on it
- **Instance profile** (`ec2InstanceProfileArn`) -- The instance-side identity; the ECS agent's permissions come from here
- **Attribute-based requirements** (`instanceRequirements`) -- `memoryMib.min` and `vcpuCount.min` are the two required dimensions; every other field narrows the candidate set. `cpuManufacturers: [amazon-web-services]` expresses a Graviton-only fleet
- **Idle scale-in** (`scaleInAfterSeconds: 300`) -- Retires instances idle for 5 minutes; `-1` disables scale-in entirely
- **Creating the provider launches nothing** -- Instances appear only when a service's `capacityProviderStrategy` schedules tasks onto `mi-general`

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region where the cluster will be created (e.g., `us-west-2`) | AWS region list |
| `<infrastructure-role-arn>` | ARN of the ECS Managed Instances infrastructure role | `AwsIamRole` status outputs |
| `<instance-profile-arn>` | ARN of the instance profile for the launched instances | `AwsIamInstanceProfile` status outputs |
| `<private-subnet-1-id>` / `<private-subnet-2-id>` | Subnets the instances launch into (two AZs for availability) | `AwsSubnet` status outputs |

## Related Presets

- **01-fargate-standard** -- Serverless-only baseline with no instances at all
- **03-ec2-capacity** -- EC2 capacity you own: wrap your auto-scaling group when you need full control of AMI and user data
