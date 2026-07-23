# EKS with IRSA for AWS Scalers

This preset installs KEDA on an EKS cluster with IAM Roles for Service
Accounts (IRSA) wired for the scalers: KEDA's own service account is
annotated with an IAM role, so scalers that read AWS metric sources (SQS
queue depth, CloudWatch metrics, Kinesis shard counts, DynamoDB streams)
authenticate through the cluster's OIDC provider — no stored access keys
anywhere. It also tightens the admission webhooks to `Fail`, so broken
ScaledObjects are rejected at apply time instead of surfacing as runtime
scaling failures.

## When to Use

- EKS clusters scaling workloads on AWS signals (SQS, CloudWatch, Kinesis,
  DynamoDB) with keyless identity
- Teams standardizing on IRSA for every addon that touches AWS APIs
- Postures that prefer apply-time rejection of invalid scaling
  declarations over silent unvalidated applies

## Key Configuration Choices

- **`podIdentity.awsIrsa`** — annotates KEDA's service account with the
  IAM role; ScaledObjects then use `podIdentity` provider `aws` in their
  TriggerAuthentication instead of credential Secrets. Per-trigger
  authentication beyond the ambient identity still lives in
  TriggerAuthentication resources next to the workloads.
- **`roleArn`** — must satisfy the IAM role ARN shape (the spec validates
  it); the example account ID `111111111111` is a placeholder-by-value.
  The role's trust policy must allow the cluster's OIDC provider and the
  `keda-operator` service account in the `keda` namespace.
- **`webhooks.failurePolicy: Fail`** — stricter than the chart's `Ignore`
  default: a webhook outage blocks ScaledObject changes, but nothing
  unvalidated ever lands.
- **Everything else at spec/chart defaults** — single replicas,
  operator-managed certificates, CRDs installed and kept on uninstall.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `arn:aws:iam::111111111111:role/keda-scalers` | The real IAM role ARN for KEDA's scalers (replace account ID and role name) | AWS IAM console or your `AwsIamRole` resource outputs |

## Related Presets

- **01-cluster-standard** — the default posture when scalers need no cloud
  identity
- **03-ha-production** — layer HA replicas, resources, and telemetry on
  top of this identity wiring
