# Compliance Log Scrubbing

This preset creates a Standard profile with access-log scrubbing turned
all the way up -- client IP addresses, request URIs, and query-string
arguments are masked before Front Door writes its logs.

## When to Use

- Workloads under GDPR/CCPA-style rules where client IPs are personal
  data
- APIs that carry tokens or identifiers in paths or query strings
  (which otherwise land verbatim in access logs)
- Any environment whose log pipeline is broader-access than its request
  path

## Key Configuration Choices

- **All three scrubbing variables** -- Azure scrubs ALL values of each
  selected part (the service supports only the match-everything
  operator on profiles); scrub selectively by listing fewer variables
- **Scrubbing is enabled by presence** -- an empty list disables it;
  there is no separate on/off switch
- **Tags carry the compliance context** -- `data-classification` on the
  profile makes the posture visible to Azure Policy and cost tooling

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The AzureResourceGroup's Planton resource name | Your Azure composition |
| `profileName` (example value) | 2-90 chars, letters/digits/hyphens -- rename to your convention | Your naming convention |
| `<cost-center>` | Your org's cost-center tag value | Your governance conventions |
