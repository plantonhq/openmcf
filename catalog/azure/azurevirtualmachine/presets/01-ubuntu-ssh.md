# Ubuntu Server with SSH Keys

This preset creates a zonal Ubuntu 24.04 LTS VM authenticated by SSH keys only, attached to a referenced `AzureNetworkInterface`, with managed boot diagnostics. It is the canonical Linux production shape: password authentication stays disabled (the default), and every composable piece -- the NIC, and through it the subnet, public IP, and NSG -- is a first-class referenced resource.

## When to Use

- General-purpose Linux application and web servers
- Bastion-less fleets managed over SSH from inside the network
- The starting point for any Linux workload before layering data disks, identity, or spot

## Key Configuration Choices

- **`networkInterfaceIds` reference** -- resolves to the NIC's `network_interface_id` output. The VM's entire network posture (subnet, public exposure, NSG filtering) lives on the NIC, so changing it never touches the VM
- **SSH keys, no password** -- `disablePasswordAuthentication` defaults to true; the key is public material and safe in a manifest
- **`osDisk` Premium** -- the production default; the OS disk is the one disk born and dying with the VM. Data belongs on referenced `AzureManagedDisk` resources (`dataDiskAttachments`) that outlive it
- **`version: latest` resolves at creation only** -- the VM does not follow new image releases; pin a version for reproducible fleets
- **Zonal placement** -- zone "1" here; NICs' public IPs and zonal disks must match the zone
- **`bootDiagnostics: {}`** -- managed storage, no account to operate; the serial console is the first tool when a VM will not boot

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the NIC's region) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-network-interface-resource-name>` | Planton metadata name of the `AzureNetworkInterface` | Your NIC resource |
| `<your-ssh-public-key>` | An OpenSSH public key, e.g. `ssh-ed25519 AAAA...` | `~/.ssh/id_ed25519.pub` or your key management |
