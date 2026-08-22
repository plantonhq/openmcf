# CloudflareZeroTrustDevicePostureRule guide

The judgment this guide protects you from: the type/input pairing is API-owned, and a posture rule's lifecycle is entangled with the policies that consume it.

## Which input fields a type reads is Cloudflare's contract, not the schema's

The spec enum-walls every value list the provider walls (operators, statuses, states), but it deliberately does NOT wall which input fields belong to which type -- 23 types times ~36 fields is Cloudflare's pairing to evolve. A wrong pairing (say, `version` on a `file` rule) fails at the API with a clear message, not silently. Start from the check family named in each field's comment.

## Delete consumers before rules

Access and Gateway policies reference rules by UUID. Deleting a rule a policy still requires makes that requirement unevaluable -- depending on the policy's shape, that can mean nobody passes (a soft outage) or the check silently stops gating (a security hole). Retire the policy reference first, then the rule.

## Integration-backed checks depend on an integration this catalog does not manage

The `*_s2s`, `intune`, and `workspace_one` types read a posture INTEGRATION (CrowdStrike, Intune, and kin) -- a separate Cloudflare object holding vendor credentials, managed outside this catalog today. `input.connection_id` takes its UUID as a literal. If the integration is deleted or its credentials lapse, every rule pointing at it starts failing devices; treat integration health as part of these rules' operational surface.

## schedule and expiration trade freshness against noise

`schedule` is how often the client re-checks (min 1m); `expiration` is how long a result may be trusted. An expiration shorter than the schedule means devices oscillate to "unknown" between polls -- keep expiration comfortably above schedule, or leave it empty to trust the latest result.

## Pairs well with

- [CloudflareZeroTrustAccessPolicy](../cloudflarezerotrustaccesspolicy/README.md) -- require checks in front of applications.
- [CloudflareZeroTrustGatewayPolicy](../cloudflarezerotrustgatewaypolicy/README.md) -- require checks for network egress.
- [CloudflareZeroTrustDeviceDefaultProfile](../cloudflarezerotrustdevicedefaultprofile/README.md) -- the WARP client the checks run under.
