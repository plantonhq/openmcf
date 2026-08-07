---
title: "File System"
description: "File System deployment documentation"
icon: "package"
order: 100
componentName: "ocifilesystem"
---

# File System on OCI

Deploys an Oracle Cloud Infrastructure File Storage file system -- an NFS-compatible, fully managed network file system bundled with a dedicated mount target and one or more NFS exports. The mount target provides the NFS endpoint (private IP) in a subnet, and exports define the paths at which the file system is accessible. Each export supports per-client NFS access control rules covering read/write permissions, identity squashing, and privileged port requirements. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments, subnets, security groups, and encryption keys.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **File System** -- the NFS file system in the specified compartment and availability domain, with optional KMS encryption and snapshot policy attachment
- **Mount Target** -- the NFS endpoint in the specified subnet, providing a private IP address that clients use to mount the file system. Created in the same availability domain as the file system.
- **Export Set Configuration** -- created only when `maxFsStatBytes` or `maxFsStatFiles` is set on the mount target; controls the NFS capacity values reported to clients via `statfs`
- **Exports** -- one export per entry in `exports`, each binding the file system to the mount target at a specific NFS path with optional per-source access control rules
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the file system and mount target

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the file system and mount target in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- An availability domain for the file system and mount target (both must reside in the same AD). OCI's default service limit is 2 mount targets per availability domain -- request a limit increase if deploying more.
- A private subnet for the mount target. NFS clients must be able to reach port 2049/TCP and 111/TCP (portmapper) on the mount target IP.
- Optionally, network security groups to control NFS traffic to the mount target.
- For customer-managed encryption: an OCI KMS key. When omitted, Oracle-managed encryption is used.

## Deploy

### Console

Open the deployment store, find **File System on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Shared Application Storage** preset in the [Presets](#presets) tab to pre-populate a file system with a single export, root squashing, and privileged port requirements.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciFileSystem
metadata:
  name: shared-data
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  availabilityDomain: "Uocm:US-ASHBURN-AD-1"
  displayName: shared-data
  mountTarget:
    subnetId:
      value: "ocid1.subnet.oc1..example"
    displayName: shared-data-mt
    hostnameLabel: shareddata
  exports:
    - path: /shared
      exportOptions:
        - source: "10.0.1.0/24"
          access: read_write
          identitySquash: root_squash
          requirePrivilegedSourcePort: true
          anonymousUid: 65534
          anonymousGid: 65534
```

```shell
planton apply -f file-system.yaml
```

This creates a file system with a mount target in the specified subnet, a single `/shared` export with read-write access for the 10.0.1.0/24 CIDR, root squashing, and privileged port enforcement. Oracle-managed encryption is used.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the file system to a compartment, subnet, and security group deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: storage
      fieldPath: status.outputs.compartmentId
  mountTarget:
    subnetId:
      valueFrom:
        kind: OciSubnet
        name: nfs-subnet
        fieldPath: status.outputs.subnetId
    nsgIds:
      - valueFrom:
          kind: OciSecurityGroup
          name: nfs-nsg
          fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, subnet, and security group first, then provisions the file system and mount target with the resolved values.

## Key Configuration

These are the most important decisions when configuring a file system. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Availability domain** -- Set `availabilityDomain` to the AD where compute instances will mount this file system. The file system and mount target are both created in this AD. AD selection is immutable -- changing it forces recreation of all resources.

**Mount target networking** -- Configure `mountTarget.subnetId` for the subnet where the NFS endpoint will be placed. Optionally set `hostnameLabel` for VCN-internal DNS (produces an FQDN like `<hostname>.<subnet>.<vcn>.oraclevcn.com`). Optionally set `ipAddress` to assign a specific private IP. Subnet, hostname, and IP are all immutable after creation.

**Export paths and access control** -- Define at least one export in `exports`. Each export specifies a `path` (must start with `/`, unique within the mount target). Add `exportOptions` entries to control per-source CIDR access: `access` (read_write or read_only), `identitySquash` (root_squash for production, no_squash for dev), `requirePrivilegedSourcePort`, and anonymous UID/GID mappings.

**Capacity reporting** -- Set `mountTarget.maxFsStatBytes` and `mountTarget.maxFsStatFiles` to control the capacity values reported to NFS clients via `statfs`. Useful for applications that make provisioning decisions based on reported free space. When omitted, the actual metered file system size is reported.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciSubnet** | `mountTarget.subnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `mountTarget.nsgIds` | `status.outputs.networkSecurityGroupId` |
| **OciKmsKey** (optional) | `kmsKeyId` | `status.outputs.keyId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `file_system_id` | OCID of the file system | Snapshot policies, replication targets, IAM policy scoping |
| `mount_target_id` | OCID of the mount target | Monitoring alarms, resource management |
| `mount_target_ip_address` | Private IP address of the mount target | NFS mount commands (`mount -t nfs <ip>:<path> /mnt`), DNS records |
| `export_set_id` | OCID of the export set on the mount target | Advanced export set management |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Shared application storage** -- A file system with a single `/shared` export, root squashing, and privileged port enforcement. Covers the standard use case of shared NFS storage for compute instances or container workloads in a single subnet. Start from the **Shared Application Storage** preset.

**Restricted multi-export** -- A file system with KMS encryption, NSG-restricted mount target, and multiple exports with differentiated access. Application subnets get read-write access while monitoring subnets get read-only access with full identity squashing. Start from the **Restricted Multi-Export** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes the file system and mount target
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the subnet where the mount target NFS endpoint is placed
- [**Network Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides network security rules for NFS traffic to the mount target
- [**KMS Key on OCI**](/cloud-catalog/oci-kms-key) -- provides the customer-managed encryption key for the file system