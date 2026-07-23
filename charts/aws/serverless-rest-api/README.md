# AWS Serverless REST API

A production serverless REST API: your handler code on Lambda behind an
HTTP API Gateway catch-all, a DynamoDB table it owns end to end, optional
Cognito JWT authentication in front of every route, and error alarms
wired to email. Nothing runs (or bills) between requests — the whole
stack scales from zero to a traffic spike and back without a capacity
decision.

The path from zero to "my API is serving requests" is: upload your
deployment zip to S3, point the chart at it, deploy, and call the
invoke URL from the gateway's outputs.

## Architecture

```
 clients ──▶ AwsHttpApiGateway ($default stage, auto-deploy)
               │        ▲
               │   JWT authorizer ─── AwsCognitoUserPool + Client (auth_enabled)
               │   CORS answered at the gateway            (cors_enabled)
               │
               └─ $default route ── AWS_PROXY (payload 2.0)
                                       │
                              AwsLambda (zip XOR container image)
                                       │  TABLE_NAME env var
                                       │  invoke permission: apigateway, this account
                                       ├── AwsIamRole (logs + item-level table access)
                                       └── AwsDynamodb (on-demand, PITR, keyed pk[/sk])

 AwsCloudwatchAlarm (Errors, Throttles) ──▶ AwsSnsTopic ──▶ email   (alarms_enabled)
```

Deployment order derives from the references: the role and table first,
then the function, then the gateway (and the Cognito pair before it when
auth is on); alarms join once the topic exists.

## Included Cloud Resources

| Resource | Kind | Purpose |
|----------|------|---------|
| Execution role | `AwsIamRole` | Logs write + item-level access to this chart's table, nothing else |
| Table | `AwsDynamodb` | On-demand storage with PITR and deletion protection |
| Function | `AwsLambda` | Your handler — S3 zip or container image, one artifact for all routes |
| HTTP API | `AwsHttpApiGateway` | The front door: catch-all route, auto-deployed `$default` stage |
| User pool | `AwsCognitoUserPool` | Email sign-in, strong passwords, self-service recovery (conditional) |
| App client | `AwsCognitoUserPoolClient` | Public SRP client for your frontend (conditional) |
| Alerts topic | `AwsSnsTopic` | Alarm fan-out target (conditional) |
| Email subscription | `AwsSnsSubscription` | Delivers alerts to your inbox (conditional) |
| Error/throttle alarms | `AwsCloudwatchAlarm` ×2 | Pages when callers see 5xx or 429 (conditional) |

## Parameters

| Name | Description | Default | Required |
|------|-------------|---------|----------|
| `aws_region` | Region for every resource (and the code bucket) | `us-east-1` | yes |
| `aws_account_id` | Account id scoping the invoke permission and table policy | `123456789012` | yes |
| `api_name` | Name prefix for every resource | `rest-api` | yes |
| `lambda_code_bucket` | S3 bucket holding the deployment zip | `my-lambda-artifacts` | zip arm |
| `lambda_code_key` | Object key of the zip | `api/function.zip` | zip arm |
| `lambda_runtime` | Runtime for the zip (e.g. `python3.12`) | `python3.12` | zip arm |
| `lambda_handler` | Entry point inside the zip | `app.handler` | zip arm |
| `container_image_enabled` | Package as an ECR image instead of a zip | `false` | no |
| `container_image_uri` | Private ECR image URI | placeholder | image arm |
| `lambda_memory_mb` | Memory (and proportional CPU) per invocation | `256` | no |
| `lambda_timeout_seconds` | Per-invocation ceiling (gateway caps requests at 30 s) | `29` | no |
| `table_partition_key` | Partition (HASH) key attribute — immutable | `pk` | yes |
| `sort_key_enabled` | Add a RANGE key (single-table design shape) — immutable | `true` | no |
| `table_sort_key` | Sort key attribute name | `sk` | when sort key on |
| `table_deletion_protection` | Refuse table deletion while enabled | `true` | no |
| `auth_enabled` | Cognito JWT auth on every route | `false` | no |
| `cors_enabled` | Return CORS headers for browser callers | `false` | no |
| `cors_allowed_origins` | Origins allowed cross-origin calls | `["https://app.example.com"]` | when CORS on |
| `alarms_enabled` | Error/throttle alarms wired to email | `true` | no |
| `alert_email` | Alert destination (confirm AWS's first email) | `ops@example.com` | when alarms on |

## First deploy

1. Build your deployment zip and upload it:

   ```bash
   zip -r function.zip app.py
   aws s3 cp function.zip s3://my-lambda-artifacts/api/function.zip
   ```

   The minimal Python handler this chart's defaults expect (`app.handler`,
   payload format 2.0):

   ```python
   import json, os

   def handler(event, context):
       return {
           "statusCode": 200,
           "body": json.dumps({"table": os.environ["TABLE_NAME"],
                               "path": event["rawPath"]}),
       }
   ```

2. Set `aws_region`, `aws_account_id`, `api_name`, and the three code
   params; deploy the chart.

3. Call the API — the invoke URL is in the gateway's outputs
   (`status.outputs.stage_invoke_url`):

   ```bash
   curl https://<api-id>.execute-api.<region>.amazonaws.com/
   ```

4. If `alarms_enabled` is on, click the confirmation link AWS emailed to
   `alert_email` — no alert is delivered until it is confirmed.

## Using authentication

With `auth_enabled: true`, every route rejects requests without a valid
Cognito JWT before your code runs. To exercise it end to end:

```bash
# Create a user (admin-side) and set a permanent password
aws cognito-idp admin-create-user --user-pool-id <pool-id> \
  --username you@example.com --message-action SUPPRESS
aws cognito-idp admin-set-user-password --user-pool-id <pool-id> \
  --username you@example.com --password '<StrongPassword1!>' --permanent

# Sign in through the app client and call the API with the token
aws cognito-idp initiate-auth --auth-flow USER_PASSWORD_AUTH \
  --client-id <client-id> \
  --auth-parameters USERNAME=you@example.com,PASSWORD='<StrongPassword1!>'
curl -H "Authorization: Bearer <IdToken>" https://<invoke-url>/
```

The pool id and client id are in the Cognito resources' outputs. Note:
`USER_PASSWORD_AUTH` for the CLI test above requires temporarily adding
`ALLOW_USER_PASSWORD_AUTH` to the client's `explicitAuthFlows`; real
frontends use the SRP flow (Amplify, amazon-cognito-identity-js) that the
client ships with, where the password never crosses the wire.

## Day-2 guidance

- **Gateway-level alarms.** 5xx-rate and integration-latency alarms need
  the `ApiId` CloudWatch dimension, which only exists after the first
  deploy. Read `status.outputs.api_id` from the gateway, then create
  `AwsCloudwatchAlarm` resources with `namespace: AWS/ApiGateway`,
  `metricName: 5xx` (and `IntegrationLatency`), and
  `dimensions: {ApiId: <api-id>, Stage: $default}`, pointing
  `alarmActions` at this chart's alerts topic.
- **Tighten the invoke permission.** The chart scopes the function's
  invoke grant to API Gateway calls from your account. After the first
  deploy, read `status.outputs.execution_arn` from the gateway and
  replace `sourceAccount` in the function's `invokePermissions` with
  `sourceArn: <execution-arn>/*/*` — the grant then names exactly this
  API.
- **Custom domain.** Map `api.example.com` to this API with an
  `AwsHttpApiDomain` resource (ACM certificate + API mapping), then set
  `disableExecuteApiEndpoint: true` on the gateway so the default URL
  stops bypassing your domain's TLS policy.
- **Rolling code.** Upload a new zip under a versioned key and update
  `lambda_code_key` (or push a new image tag and update
  `container_image_uri`) — rollback is re-pointing at the previous
  artifact.
- **Spend guardrails.** On-demand tables scale with whatever traffic
  arrives, including abusive traffic. `onDemandThroughput` ceilings on
  the table and a `defaultThrottle` on the gateway stage cap the worst
  case; both are in-place updates on the deployed resources.
- **Teardown.** `table_deletion_protection` is on by default — flip it
  to `false` and apply before destroying, or the teardown stops at the
  table (by design).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
