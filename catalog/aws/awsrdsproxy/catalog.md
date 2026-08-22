# AWS RDS Proxy

Stop burning database connections: the proxy pools and multiplexes thousands of application connections onto the handful the database can actually hold. The standard cure for Lambda-to-RDS connection storms, and the graceful-failover front for Aurora.

## What Gets Managed

- The proxy: engine family, network placement, TLS enforcement, idle timeout, debug logging.
- Its sign-ins: Secrets Manager credentials with per-secret IAM auth posture.
- The default target group's connection-pool tuning (max connections, idle ceiling, borrow timeout, pinning filters).
- Additional endpoints (read-only endpoints for Aurora reader farms).
- The database target: one RDS instance or one Aurora cluster.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with RDS, IAM, and Secrets Manager permissions.

### AWS Prerequisites

- A Secrets Manager secret holding the database credentials (the standard username/password JSON).
- An IAM role trusting rds.amazonaws.com with GetSecretValue on that secret (plus kms:Decrypt for customer-managed keys).
- At least two subnets in different availability zones.

## After You Deploy

- Point applications at the `endpoint` output instead of the database endpoint — connection behavior is otherwise identical.
- With `iam_auth: REQUIRED`, clients sign in with IAM auth tokens over TLS instead of passwords.

## Common Changes

- Tune the pool (in-place), add sign-ins or endpoints (in-place), swap the target (re-registration).
- Changing engine family or subnets replaces the proxy — its endpoint DNS names change; roll application config with it.
