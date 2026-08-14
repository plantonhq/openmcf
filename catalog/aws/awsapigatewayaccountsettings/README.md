<p align="center">
  <img src="logo.svg" alt="AWS API Gateway Account Settings" width="80"/>
</p>

# AWS API Gateway Account Settings

Manage the region's [API Gateway account object](https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-logging.html)
— the CloudWatch Logs role every REST API stage in the region logs
through.

This is a **settings singleton**: AWS keeps exactly one API Gateway
account object per account+region, so deploy at most one instance per
region. `metadata.name` never reaches AWS. Stage-level logging on
[AwsRestApiGateway](../awsrestapigateway) silently does nothing until
this role is set — that cross-component contract is why this kind
exists.

## What Gets Managed

- **The region's CloudWatch logging role** for API Gateway. The role
  must trust `apigateway.amazonaws.com` and carry CloudWatch Logs
  write permissions (AWS's managed
  `AmazonAPIGatewayPushToCloudWatchLogs` policy is the canonical
  grant).

Destroying this component **resets the role to none** — the account
object itself always exists and cannot be deleted. Managing the
settings is free.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
