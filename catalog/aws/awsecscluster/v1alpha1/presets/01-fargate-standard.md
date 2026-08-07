# Standard Fargate Cluster

This preset creates an ECS cluster on AWS Fargate with enhanced Container Insights and audited ECS Exec. Fargate eliminates instance management entirely -- AWS runs the compute -- and the cluster itself is free: only the tasks it schedules cost money. The standard starting point for any ECS deployment.

## When to Use

- Running containerized workloads without managing EC2 instances
- Standard production ECS deployments where cost optimization via Spot is not yet needed
- Any ECS cluster that will run Fargate tasks and services

## Key Configuration Choices

- **Fargate only** (`capacityProviders: [FARGATE]`) -- All tasks run On-Demand Fargate; predictable pricing with no Spot interruptions
- **Enhanced Container Insights** (`containerInsights: enhanced`) -- Container-level observability with automatic CloudWatch dashboards; the production posture (use `enabled` for the lighter task/service-level tier)
- **Exec auditing** (`executeCommandConfiguration.logging: DEFAULT`) -- Interactive `aws ecs execute-command` sessions log through each task's own log configuration

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region where the cluster will be created (e.g., `us-west-2`) | AWS region list |

## Related Presets

- **02-fargate-cost-optimized** -- Run a portion of tasks on Fargate Spot for cost savings
- **03-ec2-capacity** -- Add EC2 capacity backed by an auto-scaling group
