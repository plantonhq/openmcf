# AwsRdsProxy — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## The role is the usual first failure

A proxy that sits INCOMPATIBLE_AUTH or never reaches AVAILABLE almost always has a role problem: the trust policy must name rds.amazonaws.com, and the policy needs secretsmanager:GetSecretValue on every auth secret plus kms:Decrypt when the secrets use a customer-managed key. The E2E verifier demands AVAILABLE, never mere existence, for exactly this class.

## IAM auth is a TLS decision too

`iam_auth: REQUIRED` only works over TLS — pair it with `require_tls`. Clients then generate 15-minute auth tokens (`rds generate-db-auth-token`) instead of passwords; the proxy still signs in to the database with the secret's credentials either way.

## Session pinning is the silent multiplexing killer

The proxy multiplexes only while connections stay interchangeable. Session state (prepared statements, session variables, temp tables) PINS a connection to one client and quietly turns the pool into 1:1 passthrough. `EXCLUDE_VARIABLE_SETS` relaxes the variable-set trigger (MySQL family); beyond that, pinning is an application-behavior fix, not a proxy dial.

## Watch max_idle_connections_percent on small databases

The pool's defaults are tuned for big instances. On a db.t-class instance with low max_connections, a 50% idle ceiling can starve the database's own headroom — drop it (10-20%) so the proxy returns connections aggressively.

## Endpoints multiply networks, not capacity

Additional endpoints exist for network topology (a second VPC's subnets, IPv6) and Aurora read-only routing — they do not add throughput. READ_ONLY endpoints distribute across Aurora readers; against a plain instance target they have nothing to distribute to.
