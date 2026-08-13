# Terraform Module: AWS Secrets Manager Secret

Provisions an AWS Secrets Manager secret using Terraform.

## Resources Created

- `aws_secretsmanager_secret` — The named, KMS-encrypted container:
  description, key choice, cross-region replicas, deletion recovery
  window, and the managed-external-secret partner type.
- `aws_secretsmanager_secret_policy` — The resource policy (only when
  `policy` is declared), rendered through the standalone resource so
  `block_public_policy` (default on) rejects policies granting anonymous
  access.
- `aws_secretsmanager_secret_version` — The managed version (only when a
  value arm is set). Custom staging labels always ride WITH `AWSCURRENT`
  — providing `version_stages` replaces AWS's automatic assignment, so
  the module concats them.
- `aws_secretsmanager_secret_rotation` — Rotation (only when `rotation`
  is declared): exactly one of the self-managed rotation Lambda or the
  partner-managed external rotation role, ordered after the version so
  the immediate first rotation has a value to read.

## Usage

```hcl
module "secret" {
  source = "./path/to/module"

  metadata = {
    name = "prod/payments/db"
    org  = "my-org"
    env  = "production"
    id   = "awssm-abc123"
  }

  spec = {
    region       = "us-west-2"
    description  = "Payments database credentials"
    string_value = "resolved-just-in-time"

    policy = {
      Version = "2012-10-17"
      Statement = [{
        Sid       = "AllowReaderRole"
        Effect    = "Allow"
        Principal = { AWS = "arn:aws:iam::123456789012:role/reader" }
        Action    = "secretsmanager:GetSecretValue"
        Resource  = "*"
      }]
    }

    replica_regions = [{ region = "us-east-1" }]

    recovery_window_in_days = 7
  }
}
```

The value arms are mutually exclusive (spec validation): `string_value`
for text/JSON, `binary_value` for base64 binary. Both arrive resolved
just-in-time from the platform's managed-secret store — plaintext never
lives in the control plane. A shell secret (no value) is legal; an
application or rotation function writes the first version then.

Deletion is soft by default (`recovery_window_in_days` 30); 0 deletes
immediately and permanently — the right choice for ephemeral secrets,
since a soft-deleted secret reserves its name for the window.
