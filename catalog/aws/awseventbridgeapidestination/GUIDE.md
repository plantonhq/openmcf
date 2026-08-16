# AwsEventBridgeApiDestination — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## Your credentials live in AWS's secret, not yours

CreateConnection stores the credential values in a Secrets Manager secret AWS creates and owns (`events!connection/...` — the `connection_secret_arn` output). No read API ever returns them: the provider reads secrets back from prior state, drift in a credential is invisible, and a fresh import re-sends the configured values as a one-time server-side no-op. Treat the manifest's secret references as the source of truth.

## The auth state machine gates everything

Connection creates and credential updates walk CREATING/AUTHORIZING → AUTHORIZED (up to 20 minutes budgeted; usually seconds-to-minutes). A connection stuck DEAUTHORIZED after an update means the downstream rejected the new credentials — the StateReason on the connection says why.

## Rate limits queue, not drop

Invocations beyond `invocation_rate_limit_per_second` are queued by EventBridge for up to 24 hours, then dead-lettered (if the invoking rule has a DLQ) or dropped. A too-low limit shows up as delay first, loss second — size it to the downstream API's real ceiling.

## The endpoint is not validated at create

AWS accepts any HTTPS URL — a typo'd endpoint deploys green and fails at first invocation. The reserved-domain examples in presets are deliberate; replace them with the real endpoint before wiring a rule.

## Grant-shaped KMS policies

With a customer-managed key, the key policy must allow Secrets Manager decryption scoped to the AWS-created secret — condition on `kms:EncryptionContext:SecretARN` matching `arn:aws:secretsmanager:*:*:secret:events!connection/*`.
