# SSH Key on DigitalOcean

Registers an SSH public key on the DigitalOcean account, ready to be injected into droplets and droplet autoscale pools at create time. Only the public half is stored -- the private key never leaves your machine. Integrates with Planton's Provider Connections for DigitalOcean API token management; droplets reference the key by its exported numeric id or fingerprint.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **SSH Key** -- the named public key on the account, with a DigitalOcean-computed fingerprint

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.

### DigitalOcean Account

- Nothing: SSH keys are free account metadata.

## After You Deploy

Reference the `ssh_key_id` (or `fingerprint`) output from droplet manifests' `sshKeys` lists so new droplets are born with the key installed -- DigitalOcean injects keys only at droplet create time.
