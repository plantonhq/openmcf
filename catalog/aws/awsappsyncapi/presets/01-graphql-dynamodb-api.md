# GraphQL API over DynamoDB

The serverless CRUD backbone: a GraphQL API whose resolvers read and write a DynamoDB table directly — no Lambda in the middle, no servers anywhere. Cognito users sign the requests.

## What this shape gives you

- **Zero-hop resolvers**: APPSYNC_JS resolvers call DynamoDB's GetItem/PutItem straight from the API layer. Latency is the table's, and there is no function cold start.
- **Cognito authorization**: the primary auth provider validates the caller's user pool token; `defaultAction: ALLOW` passes matched requests through (schema directives can still restrict fields).
- **Chart LEGO**: the user pool, the table, and the data-source role are all references — deploy this beside AwsCognitoUserPool, AwsDynamodb, and AwsIamRole resources and the wiring resolves at deploy.

## Adapt it

- The role needs dynamodb:GetItem/PutItem on the table and a trust policy for appsync.amazonaws.com.
- Add API_KEY as an additional provider for public read-only fields.
- Turn on `cache` (SMALL) when read traffic concentrates on hot items — remember it bills per hour.
