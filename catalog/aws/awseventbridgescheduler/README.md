# AwsEventBridgeScheduler

One EventBridge Scheduler schedule: cron for the cloud — a cron, rate, or one-time expression invoking one target under an execution role, with flexible time windows, bounded retries, and a dead-letter queue. The schedule group folds own-XOR-existing (a name-and-tags container whose provider update path is literally tags-only).

## Highlights

- **Full target depth**: the universal ARN+role target plus all five service parameter blocks (ECS RunTask at full depth, EventBridge PutEvents, Kinesis, SageMaker pipelines, SQS FIFO), guarded by an at-most-one CEL mirroring the provider.
- **The group is own-XOR-existing**: create it here (it carries the identity tags — the schedule itself is untaggable at AWS) or join any group by name; unset means AWS's `default` group. Fixed for life.
- **Contracts taught in place**: the two-minute IAM-propagation retry on first deploys, `action_after_completion: DELETE` deleting completed one-time schedules out from under IaC state, FLEXIBLE windows requiring a size (and OFF forbidding one) as CELs.
- **The target ARN is bare polymorphic** — no single kind dominates, so a `valueFrom` on it states its `kind:` explicitly.

## Both Engines

Both modules render the schedule (and optional group) identically and export the same outputs: `schedule_arn`, `group_name` (with metadata.name, the `{group}/{name}` import ID), `group_arn`.

## Chart Wiring

`target.role_arn` → AwsIamRole `role_arn`; `target.dead_letter_queue_arn` → AwsSqsQueue `queue_arn`; ECS parameters → AwsEcsTaskDefinition `task_definition_arn`, AwsSubnet `subnet_id`, AwsSecurityGroup `security_group_id`; `kms_key_arn` → AwsKmsKey `key_arn`. The target ARN wires to any invocable kind with an explicit `kind:`.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
