# HTTP Redirect VIP (Shared IP)

The other half of the standard pair: port 80 on the SAME static IP as the 443 rule, pointing at a target HTTP proxy that serves an http→https redirect URL map. Two forwarding rules may share one `[IPAddress, IPProtocol]` exactly because their port ranges do not overlap.

## When to Use

- Alongside preset 01 for every production HTTPS frontend — browsers and old links still arrive on port 80

## Remix Notes

- Both rules MUST reference the same `GcpGlobalAddress` (same `valueFrom`) — a different IP would leave port-80 traffic pointing at a stranger.
- The scheme should match the 443 rule's scheme; mixing `EXTERNAL` and `EXTERNAL_MANAGED` on one IP works but complicates migration bookkeeping.
- The redirect chain is: this rule → `GcpTargetHttpProxy` → a redirect-only `GcpUrlMap` (no backends anywhere).
