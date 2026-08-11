# External Endpoint

This preset adds a destination that lives outside Azure -- an on-premises datacenter, another cloud, or any service with a resolvable name -- so Traffic Manager can steer across your whole estate during migrations or in hybrid steady-state.

## When to Use

- Hybrid fronts: cloud primary with on-premises standby (or the reverse, mid-migration)
- Multi-cloud steering across providers behind one DNS name

## Key Configuration Choices

- **The target is a frozen string** -- unlike azure endpoints, nothing tracks it; retargeting is an edit. A reference (`valueFrom`) works when another component mints the hostname
- **`endpointLocation` anchors latency routing** -- external targets carry no discoverable region; the service REQUIRES it under Performance routing and ignores it otherwise (give the nearest Azure region)
- **Priority 20 behind the preset-01 primary** -- the two presets together form an active/passive pair on a Priority profile

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-azure-traffic-manager-profile-resource-name>` | The AzureTrafficManagerProfile component's resource name | Your Planton catalog |
| `<your-target-hostname-or-ip>` | The external destination | Your datacenter/other-cloud DNS records |
| `<nearest-azure-region>` | The Azure region latency-wise closest to the target | Azure region list (e.g. `eastus`) |
