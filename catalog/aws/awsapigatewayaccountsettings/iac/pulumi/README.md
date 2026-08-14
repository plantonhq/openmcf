# AwsApiGatewayAccountSettings — Pulumi module (Go)

Manages the region's API Gateway account object
(`apigateway.Account`) — a settings singleton whose one lever is the
CloudWatch Logs role REST API stages log through.

Module facts worth knowing before editing:

- **One object per account+region.** The module always renders exactly
  one resource; deploying two instances against the same region makes
  them fight over the same settings object.
- **Empty string IS the reset.** The spec value passes through
  unconditionally: `""` and unset both patch the role to none
  upstream, so the clear posture needs no conditional.
- **Destroy resets, never deletes.** The account object survives
  destroy with its role cleared.
- **No tags.** The upstream resource carries no tags argument — the
  catalog's identity-tag map does not apply here.

Outputs mirror the Terraform module key-for-key: `account_id`,
`api_key_version`, `features`, `throttle_burst_limit`,
`throttle_rate_limit`.
