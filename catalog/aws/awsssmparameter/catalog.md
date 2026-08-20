# AWS SSM Parameter

One named configuration value your applications read at runtime —
plain config as a String, lists as a StringList, secrets as a
KMS-encrypted SecureString.

## What Gets Managed

- The parameter's name (hierarchical paths like `/prod/db/url`
  organize config and enable by-path reads), type, and value.
- Secrets stay secrets: SecureString values are supplied as
  managed-secret references and resolved just-in-time at deploy —
  plaintext never lives in the control plane or plan output.
- Guardrails: a validation pattern AWS enforces on every write, the
  storage tier, AMI-ID validation via `dataType`, and the encryption
  key.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SSM permissions (plus KMS when using
  a customer-managed key).

### AWS Account

- Nothing — parameters stand alone. For SecureString with your own
  key, a KMS key ([AWS KMS Key](/cloud-catalog/aws-kms-key)).

## Deploy

### Console

Create the resource from the AWS catalog, pick the type, set the value
(or a managed-secret reference for SecureString), and deploy.

### CLI

```bash
planton apply -f ssm-parameter.yaml
```

## After Deploy

- Applications read the value via `GetParameter`/`GetParametersByPath`
  (SecureString decrypts for callers with KMS access).
- Every value write increments the version — the previous value stays
  retrievable by version.
- Standard tier is free; Advanced bills per parameter and cannot be
  downgraded in place.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
