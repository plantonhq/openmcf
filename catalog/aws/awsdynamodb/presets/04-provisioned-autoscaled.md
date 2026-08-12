# Provisioned Table with Auto Scaling

This preset creates a provisioned DynamoDB table whose read and write capacity is managed by Application Auto Scaling: target tracking holds utilization near 70%, and two scheduled adjustments raise the read floor for business hours and lower it after. Point-in-time recovery, AWS-managed encryption, and deletion protection round out the production posture.

## When to Use

- Sustained daily traffic with a predictable shape (business-hours peaks, quiet nights) where provisioned + auto scaling undercuts on-demand pricing
- Workloads planning reserved capacity purchases (which apply only to provisioned tables) that still want elasticity above the reserved baseline
- Tables migrating from on-demand once traffic became forecastable

## Key Configuration Choices

- **Auto scaling bounds** (`autoscaling.read` / `autoscaling.write`) -- The floor is the guaranteed baseline, the ceiling is the cost guardrail; the scaler moves capacity inside those bounds to hold **70% utilization** (the production sweet spot: headroom for spikes without paying for idle)
- **Scheduled adjustments** (`scheduledAdjustments`) -- Named, timezone-aware floor changes; each targets one dimension (READ here). The morning entry pre-warms capacity before traffic arrives instead of reacting to it
- **Initial capacity** (`provisionedThroughput`) -- Seeds the table at create; from then on the scaler owns live capacity, and capacity edits in this manifest land through the scaling target
- **Production safety** -- Point-in-time recovery, the AWS-managed `aws/dynamodb` encryption key (reference an `AwsKmsKey` via `serverSideEncryption.kmsKeyArn` to hold your own key), and deletion protection

## Placeholders to Replace

- `<aws-region>` -- The AWS region for the table (for example `us-west-2`)
- `metadata.name` -- Your table name
- The schedule expressions and timezone -- match your traffic's actual shape
