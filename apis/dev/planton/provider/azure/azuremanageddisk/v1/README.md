# Overview

The **Azure Managed Disk API Resource** provides a consistent and standardized interface for deploying and managing Azure Managed Disks within our infrastructure. This resource treats block storage the way Azure itself does: as a first-class resource with its own lifecycle, SKU, encryption, and network posture — storage whose data outlives any one virtual machine.

## Purpose

We developed this API resource to make persistent block storage a deliberate architectural decision rather than a side effect of creating a VM. A standalone disk is what makes stateful infrastructure survivable, enabling:

- **Data That Outlives the Machine**: Detach a disk from a failed or replaced VM and re-attach it to the next one — the data never moves
- **The Full Origin Matrix**: Create disks empty, cloned from snapshots or disks, from platform or gallery images, imported from VHD blobs, restored from backup recovery points, or as direct-upload targets
- **Performance as a Dial**: Premium SSD v2 and Ultra Disk decouple capacity from performance — size, IOPS, and throughput are provisioned independently and updated in place
- **Shared-Disk Clustering**: `maxShares` lets several VMs attach one disk simultaneously, the seam failover databases and scale-out file systems build on

## Key Features

- **Consistent Interface**: Aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Attachment Lives on the VM**: The disk spec deliberately knows nothing about its consumers — an `AzureVirtualMachine`'s `dataDiskAttachments` references this disk's `disk_id` output with a LUN and caching mode, so replacing the VM never touches the disk
- **Validation That Mirrors ARM**: The create option's required source fields, the SKU-gated performance dials, the encryption pairings, and the network-posture rules are all enforced at spec level, exactly as Azure Resource Manager enforces them — misconfigurations fail before a deploy ever starts
- **Seven Storage SKUs**: From Standard HDD for cold data through Premium SSD for production databases to Premium SSD v2 and Ultra Disk with independently dialed IOPS and throughput
- **Complete Security Surface**: Customer-managed-key encryption via disk encryption sets, confidential-VM security profiles, trusted launch, and a network export posture that can be locked down entirely

## Use Cases

- **Database Volumes**: Premium SSD data and log disks whose data survives VM replacement, resizing, and rebuild
- **Performance-Sensitive Workloads**: Premium SSD v2 disks with IOPS and throughput dialed independently of capacity
- **Failover Clusters**: Shared disks (`maxShares`) attached to several VMs at once for clustered databases and scale-out file systems
- **Snapshot Restore and Environment Duplication**: Clone a production snapshot into staging, or restore a backup recovery point as a fresh disk
- **VHD Migration**: Import existing VHD blobs — including the secure-import path for confidential-VM scenarios
- **Regulated Workloads**: Customer-managed keys, confidential-VM encryption, and fully private (or disabled) network export

## Future Enhancements

Future updates will include:

- **Disk Encryption Set Modeling**: A first-class DiskEncryptionSet resource so customer-managed-key encryption composes by reference instead of raw ARM IDs
- **Snapshot Lifecycle**: First-class snapshot resources and scheduled snapshot policies
- **Monitoring Integration**: Built-in Azure Monitor metrics and alerts for IOPS, throughput, and burst-credit consumption
- **Comprehensive Documentation**: Expanded capacity-planning and performance-tuning guides
