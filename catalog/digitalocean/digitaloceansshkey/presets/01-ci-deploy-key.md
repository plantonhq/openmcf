# CI Deploy Key

This preset registers a dedicated ed25519 public key for a CI pipeline, so automation reaches droplets with its own credential instead of borrowing a person's. Replace the example key with your pipeline's real public key -- the private half stays in the CI system's secret store.

## When to Use

- Giving CI/CD pipelines their own droplet SSH access, revocable without touching human keys
- Any automation that provisions or configures droplets over SSH

## Key Configuration Choices

- **ed25519** -- small, fast, modern; use RSA only for clients that cannot speak it.
- **One key per system** -- a dedicated key deletes cleanly when the pipeline is retired (deleting never touches droplets created with it).

## What You Get

An account-registered key whose `ssh_key_id` and `fingerprint` outputs wire into droplet manifests' `sshKeys` lists -- new droplets are born trusting the pipeline.
