# Isolated Branches WAN

This preset creates the hub-and-spoke-only shape: branches (VPN sites) can reach Azure through their hubs but can NOT reach each other through the WAN, and latency-sensitive Office 365 traffic (the Optimize category: Exchange, SharePoint, Teams media) exits at local branch internet breakouts instead of hauling through the WAN.

## When to Use

- Estates that treat branches as untrusted islands (retail, OT/industrial networks) and force inter-branch traffic through hub-hosted inspection
- Branch fleets with decent local internet where O365 performance matters

## Key Configuration Choices

- **Branch isolation is a security posture, not a tuning knob** -- flipping `allowBranchToBranchTraffic` later changes routing behavior across the whole estate; decide deliberately at creation
- **OPTIMIZE is the conservative breakout** -- it covers only the endpoints Microsoft marks latency-critical; step up to OPTIMIZE_AND_ALLOW or ALL as your branch security stack allows

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The ARM metadata region | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
