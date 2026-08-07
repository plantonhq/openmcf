# HetznerCloudSshKey

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `hetzner-cloud.planton.dev/v1alpha1`

HetznerCloudSshKeySpec defines the specification for a Hetzner Cloud SSH key.

An SSH key registered in Hetzner Cloud can be injected into servers at creation
time, providing secure, password-less access. The key is identified by its
metadata.name and referenced by other components (e.g., HetznerCloudServer)
via its ssh_key_id output through StringValueOrRef.

Hetzner Cloud supports RSA (>= 1024 bits), ED25519, and ECDSA key types.
The public key must be in OpenSSH authorized_keys format.

Changing the public_key after creation forces replacement of the resource
because Hetzner Cloud does not allow in-place updates to the key material.

## Example

```yaml
apiVersion: hetzner-cloud.planton.dev/v1alpha1
kind: HetznerCloudSshKey
metadata:
  name: hetznercloudsshkey-demo
spec:
  publicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAEXAMPLE demo@example.com"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.publicKey` | `string` | yes |  |  |

## Field Details

### spec.publicKey

`string` · required

SSH public key content in OpenSSH authorized_keys format.

Supports RSA (>= 1024 bits), ED25519, and ECDSA key types.
Example: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA... user@host"

Changing this value forces replacement of the SSH key resource.

- rule: {"string":{"minLen":"1"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: HetznerCloudSshKey, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.ssh_key_id` | `string` | The Hetzner Cloud numeric ID of the created SSH key (as a string). Referenced by HetznerCloudServer.ssh_key_ids via StringValueOrRef. |
| `status.outputs.fingerprint` | `string` | MD5 fingerprint of the SSH public key (e.g., "aa:bb:cc:dd:..."). Computed by Hetzner Cloud from the uploaded public key material. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| HetznerCloudServer | `spec.sshKeys` | `status.outputs.ssh_key_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
