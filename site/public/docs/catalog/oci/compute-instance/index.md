---
title: "Compute Instance"
description: "Compute Instance deployment documentation"
icon: "package"
order: 100
componentName: "ocicomputeinstance"
---

# Compute Instance on OCI

Deploys an OCI Compute Instance -- a virtual machine or bare metal host with configurable flex shapes (OCPU/memory), boot source, VNIC-based networking, and optional platform security hardening (Secure Boot, TPM, memory encryption). Flex shapes let you pick the exact OCPU and memory allocation instead of choosing from fixed T-shirt sizes. Integrates with Planton's Provider Connections for OCI credential management and ValueFromRef for compartment, subnet, and security group wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Compute Instance** -- a VM or bare metal host in the specified compartment, availability domain, and subnet, with the chosen shape, boot source, and VNIC configuration
- **Primary VNIC** -- configured inline via `createVnicDetails`, determining the instance's subnet, security groups, public/private IP, and DNS hostname
- **Boot Volume** -- created implicitly from the specified image or cloned from an existing boot volume, with configurable size and performance tier
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the instance

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the instance in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A subnet for the primary VNIC. Public subnets with an Internet Gateway allow public IP assignment; private subnets require a NAT Gateway or Bastion for outbound/SSH access. Provide the subnet OCID directly or reference an OciSubnet Cloud Resource via ValueFromRef.
- An OS image OCID for the boot source. Find available images via the OCI Console or `oci compute image list`.
- An SSH public key for the `metadata.ssh_authorized_keys` field (the standard mechanism for SSH access to OCI instances).

## Deploy

### Console

Open the deployment store, find **Compute Instance on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **General-Purpose Flex Instance** preset in the [Presets](#presets) tab to pre-populate a working VM.Standard.E4.Flex configuration.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciComputeInstance
metadata:
  name: app-server
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  availabilityDomain: "Ixxj:US-ASHBURN-AD-1"
  shape: VM.Standard.E4.Flex
  shapeConfig:
    ocpus: 1
    memoryInGbs: 16
  sourceDetails:
    sourceType: image
    sourceId: "ocid1.image.oc1..example"
  createVnicDetails:
    subnetId:
      value: "ocid1.subnet.oc1..example"
  metadata:
    ssh_authorized_keys: "ssh-rsa AAAA..."
```

```shell
planton apply -f compute-instance.yaml
```

This creates a 1-OCPU flex instance with 16 GiB memory, booting from the specified image, in a subnet with default IP assignment. No public IP override, no security hardening, and no preemptible configuration.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the instance to a compartment, subnet, and security group deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: platform-compartment
      fieldPath: status.outputs.compartmentId
  createVnicDetails:
    subnetId:
      valueFrom:
        kind: OciSubnet
        name: private-app-subnet
        fieldPath: status.outputs.subnetId
    nsgIds:
      - valueFrom:
          kind: OciSecurityGroup
          name: app-nsg
          fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, subnet, and security group first, then provisions the compute instance with the resolved values.

## Key Configuration

These are the most important decisions when configuring a compute instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Shape and sizing** -- The `shape` field selects the hardware profile (e.g., `VM.Standard.E4.Flex` for AMD EPYC, `VM.Standard.A1.Flex` for Arm). Flex shapes require `shapeConfig` with `ocpus` and `memoryInGbs`. Each OCPU maps to a physical core with SMT, so 1 OCPU provides 2 vCPUs. Standard shapes use fixed configurations.

**Boot source** -- The `sourceDetails.sourceType` selects between `image` (platform or custom image) and `boot_volume` (clone an existing boot volume). Set `bootVolumeSizeInGbs` to override the default size and `bootVolumeVpusPerGb` to control performance (10=Balanced, 20=Higher, 30-120=Ultra High).

**Network placement** -- The `createVnicDetails` block configures the primary VNIC: subnet, NSGs, public/private IP, and DNS hostname. Set `assignPublicIp: false` for production instances behind a load balancer or bastion. Associate `nsgIds` for fine-grained traffic control.

**Preemptible instances** -- Set `preemptibleInstanceConfig` to create a spot-like instance at significantly lower cost. OCI can reclaim it when capacity is needed. Set `preserveBootVolume: true` to keep the boot disk intact on preemption.

**Security hardening** -- The `platformConfig` block enables Secure Boot, Measured Boot, TPM, and memory encryption. The `instanceOptions.areLegacyImdsEndpointsDisabled` flag forces IMDSv2 (token-based), preventing SSRF attacks on the metadata endpoint.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciSubnet** | `createVnicDetails.subnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `createVnicDetails.nsgIds` | `status.outputs.networkSecurityGroupId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `instance_id` | OCID of the compute instance | Monitoring dashboards, OCI CLI operations |
| `private_ip` | Private IP address of the primary VNIC | Load balancer backends, DNS records, application configuration |
| `public_ip` | Public IP address of the primary VNIC (empty when not assigned) | DNS records, direct SSH access |
| `boot_volume_id` | OCID of the boot volume attached to the instance | Backup policies, volume group membership |
| `availability_domain` | Availability domain where the instance was placed | Placement constraints for related resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**General-purpose flex** -- A 1-OCPU AMD EPYC instance with 16 GiB memory, public IP, and SSH access. The standard starting point for web servers, API backends, and general workloads. Start from the **General-Purpose Flex Instance** preset.

**Private backend** -- A production-hardened instance in a private subnet with NSG association, in-transit encryption, legacy IMDS disabled, live migration preferred, and automatic recovery. The standard pattern for application servers behind a load balancer. Start from the **Private Backend Instance** preset.

**Preemptible dev** -- A cost-optimized spot-like instance for development, testing, and CI workloads. Boot volume preserved on preemption so work-in-progress is not lost. Start from the **Preemptible Dev Instance** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this instance
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the subnet for the primary VNIC's network placement
- [**Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides network security groups for fine-grained traffic control on the VNIC