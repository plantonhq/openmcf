# Operator Break-Glass Key

This preset registers an operations team's emergency-access key, kept separate from day-to-day automation keys so incident access survives a CI credential rotation (and vice versa). Replace the example key with the team's real public key -- ideally one whose private half lives in a hardware token or a vault.

## When to Use

- Emergency SSH access to droplets during incidents, independent of automation credentials
- Fleets where every droplet should carry both an automation key and a human-escape hatch

## Key Configuration Choices

- **Separate lifecycle** -- rotating CI keys never locks operators out, and retiring an operator's key never breaks pipelines.
- **Key changes replace** -- rotating the material mints a new id and fingerprint; droplets created with the old key keep it until rebuilt or re-keyed.

## What You Get

An account-registered key whose `ssh_key_id`/`fingerprint` outputs join droplet manifests alongside the automation key.
