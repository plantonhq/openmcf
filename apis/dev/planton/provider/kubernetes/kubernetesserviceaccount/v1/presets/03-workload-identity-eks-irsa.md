# EKS IRSA (IAM Roles for Service Accounts)

This preset creates a ServiceAccount bound to an AWS IAM role via IRSA. Pods running as this identity call AWS APIs keylessly — the AWS SDK inside the pod exchanges the projected ServiceAccount token for role credentials automatically; no access keys are stored in the cluster.

## When to Use

- Pods on EKS that need AWS APIs (S3, SQS, DynamoDB, Secrets Manager, Route 53, ...)
- Replacing static AWS access keys in secrets or node-instance-profile over-granting — the two patterns IRSA exists to eliminate

## Key Configuration Choices

- **`workloadIdentity.eks.roleArn`** — the IAM role to assume; the module emits it as the `eks.amazonaws.com/role-arn` annotation, the exact key the EKS pod-identity webhook expects
- **Both halves are required** — the annotation alone grants nothing. The IAM role's trust policy must trust the cluster's OIDC provider with a condition on `system:serviceaccount:<namespace>:<ksa-name>`, and the cluster must have an OIDC provider associated
- **Name and namespace are part of the trust** — the trust policy condition matches the subject string exactly; renaming or moving this ServiceAccount silently breaks the federation (pods get the node role or no credentials, with no Kubernetes-side error)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Namespace for the ServiceAccount (embedded in the IAM trust policy condition) | Your namespace management |
| `<your-aws-iam-role-arn>` | IAM role ARN, e.g. `arn:aws:iam::123456789012:role/app-role` | AWS Console → IAM → Roles, or your AwsIamRole resource's outputs |

Also rename `aws-app-identity` to match your workload — the name is embedded in the trust policy condition too.

## Related Presets

- **01-basic** — identity with no cloud federation
- **02-workload-identity-gke** — the GCP equivalent
- **04-image-pull-secrets** — private registry credentials and automount hardening
