# OpenStackKeypair

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1`

OpenStackKeypairSpec defines the configuration for an OpenStack compute keypair.

A keypair is an SSH key pair used for authenticating access to compute instances.
You can either import an existing public key or let OpenStack generate a new keypair.

When public_key is provided, the given key is imported into OpenStack.
When public_key is omitted, OpenStack generates a new keypair and the private key
is available as a one-time sensitive output from the IaC module (not stored in
platform outputs). The generated private key is stored in IaC state and should be
retrieved immediately after creation.

The keypair name is derived from metadata.name.

Terraform resource: openstack_compute_keypair_v2
Pulumi resource: openstack.compute.Keypair

## Example

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackKeypair
metadata:
  name: test-keypair
spec:
  public_key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDxyz test@local"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.publicKey` | `string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.publicKey

`string`

public_key is the SSH public key to import, in OpenSSH authorized_keys format.
Example: "ssh-rsa AAAAB3NzaC1yc2EAAAA... user@host"

If omitted, OpenStack will generate a new keypair. In that case, the private key
is available as a sensitive IaC-level output only (not in platform stack outputs).
Security note: generated private keys are stored unencrypted in IaC state.
For production use, generating keys locally with ssh-keygen and importing the
public key is the recommended approach.

### spec.region

`string`

region overrides the region from the provider config for this keypair.
If omitted, the region from the OpenStack provider config is used.
Example: "RegionOne"

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackKeypair, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.name` | `string` | name is the name of the keypair (derived from metadata.name). |
| `status.outputs.fingerprint` | `string` | fingerprint is the MD5 fingerprint of the SSH public key. Example: "d7:62:43:93:10:a8:7e:7c:01:c8:c5:67:ba:99:5c:25" |
| `status.outputs.public_key` | `string` | public_key is the SSH public key in OpenSSH authorized_keys format. This is the imported key, or the generated public key when no public_key was provided. |
| `status.outputs.region` | `string` | region is the OpenStack region where the keypair was created. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackInstance | `spec.keyPair` | `status.outputs.name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
