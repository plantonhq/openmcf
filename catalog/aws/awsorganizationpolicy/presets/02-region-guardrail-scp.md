# Region Guardrail SCP

This preset creates the classic data-residency guardrail: workload
accounts can only act in the approved regions, with AWS's global
services exempted so nothing fundamental breaks.

## When to Use

- Estates with data-residency or cost-control reasons to fence regions
- Preventing shadow infrastructure in regions nobody monitors

## What You Get

- A deny on everything outside the approved region list, scoped to one
  OU (start narrow — widen to the root once proven)
- A `NotAction` exemption for the global services (IAM, Organizations,
  Route53, CloudFront, STS, Support) that would otherwise break under
  a region fence

## Customize

- Edit the `aws:RequestedRegion` list to your approved regions
- Extend the `NotAction` list per your stack (e.g. `waf:*`,
  `globalaccelerator:*` for other global services you use)
- Attach to the root only after a soak on the workloads OU — a wrong
  region fence at the root is the widest possible blast radius
