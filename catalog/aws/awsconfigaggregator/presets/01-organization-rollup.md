# Organization Rollup

This preset creates the org-wide compliance view: one aggregator in
the management (or delegated Config administrator) account collecting
every member account across every region — no per-account grants
needed.

## When to Use

- The organization's single pane of glass for Config data and rule
  compliance
- Security teams that need new member accounts to join the view
  automatically

## What You Get

- Every account and region in one queryable rollup
- Self-discovering membership through AWS Organizations
- Zero cost (aggregators are free; recorders in source accounts bill
  as themselves)

## Customize

- Point `roleArn` at a role trusting `config.amazonaws.com` with the
  `AWSConfigRoleForOrganizations` managed policy
- Narrow `allRegions: true` to a `regions` list for a region-scoped
  view
- Data appears only from accounts/regions running a Config recorder —
  pair with AwsConfigRecorder deployments

## Composing

```yaml
# The organization-reader role this preset expects:
apiVersion: aws.planton.dev/v1alpha1
kind: AwsIamRole
metadata:
  name: config-org-role
spec:
  region: <aws-region>
  managedPolicyArns:
    - value: arn:aws:iam::aws:policy/service-role/AWSConfigRoleForOrganizations
  trustPolicy:
    Version: "2012-10-17"
    Statement:
      - Effect: Allow
        Principal: { Service: config.amazonaws.com }
        Action: sts:AssumeRole
```
