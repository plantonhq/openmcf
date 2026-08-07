# Hetzner Cloud SSH Key

Registers an SSH public key in Hetzner Cloud for injection into servers at creation time. Servers reference imported SSH keys via their `ssh_key_ids` field, enabling secure password-less access without storing private key material in IaC state. Supports RSA (>= 1024-bit), ED25519, and ECDSA key types in OpenSSH authorized_keys format.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **SSH Key** -- a single `hcloud_ssh_key` resource registered in Hetzner Cloud with the specified public key material and metadata name

## Before You Deploy

### Hetzner Cloud Account

- **A Hetzner Cloud account** with an active project and API token.
- **An SSH key pair** generated locally (e.g., `ssh-keygen -t ed25519`). Only the public key is imported -- the private key never enters IaC state or the Hetzner Cloud API.

## Deploy

### Console

Open the deployment store, find **Hetzner Cloud SSH Key**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Import Public Key** preset in the [Presets](#presets) tab to register an existing key.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudSshKey
metadata:
  name: deploy-key
  org: acme-corp
  env: prod
spec:
  publicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA... user@host"
```

```shell
planton apply -f hetznercloud-ssh-key.yaml
```

This registers the public key in Hetzner Cloud. A Stack Job tracks the provisioning in real time. Reference the key in HetznerCloudServer manifests via `ssh_key_ids`.

## Key Configuration

These are the most important decisions when configuring an SSH key. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Public key** -- The `publicKey` field accepts an SSH public key in OpenSSH authorized_keys format. Hetzner Cloud supports ED25519, RSA (>= 1024-bit), and ECDSA key types. Changing the public key after creation forces replacement of the resource because Hetzner Cloud does not allow in-place key material updates.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `ssh_key_id` | Hetzner Cloud numeric ID of the SSH key | HetznerCloudServer `ssh_key_ids` field |
| `fingerprint` | MD5 fingerprint of the public key | Key identification and verification |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Import existing key** -- Register a locally generated SSH public key for server injection. The most common pattern. Start from the **Import Public Key** preset.

## Works With

- [**Hetzner Cloud Server**](/cloud-catalog/hetznercloud-server) -- servers reference this SSH key via `ssh_key_ids` for password-less access
