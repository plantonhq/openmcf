# AwsRoute53ResolverEndpoint — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## The endpoint is the cost; rules are free

Endpoint ENIs bill hourly (two minimum) whether or not queries flow; rules, associations, and per-query charges are negligible beside them. Share ONE outbound endpoint across many rules and VPCs — never create an endpoint per rule. Verified figures live in the generated estimate at `catalog/_pricing/estimates/awsroute53resolverendpoint.yaml`, computed from the pinned, source-dated price book.

## FORWARD needs reachability; the control plane does not check it

Creating a FORWARD rule with unreachable targets succeeds — AWS validates the shape, not the path. Queries then time out at runtime. The endpoint's security group must allow DNS egress to the targets, and the network path (VPN, Direct Connect, peering) must exist. "Rule is COMPLETE but resolution fails" is a path problem, not a rule problem.

## SYSTEM rules are subdomain surgery

A SYSTEM rule only means something as an override of a broader FORWARD rule: forward `corp.example.com` on-prem, but let `aws.corp.example.com` resolve recursively in AWS. The most specific matching domain wins. A SYSTEM rule without a covering FORWARD rule changes nothing.

## Detaching a rule replaces it

Swapping a rule to a different endpoint updates in place, but clearing its endpoint binding forces the provider to replace the rule — associations re-create with it. Plan rule-to-endpoint moves, not detaches.

## The delegation arms are typed but unproven

INBOUND_DELEGATION endpoints and DELEGATE rules (the 2025 subdomain-delegation feature) are modeled per the provider's vocabulary, but their server-side pairing contract is not yet live-verified here — expect AWS-side validation to be the guard until a delegation scenario runs.

## Metrics toggles never revert by omission

Once `rni_enhanced_metrics_enabled` or `target_name_server_metrics_enabled` is set, removing the field leaves the last value at AWS — revert with an explicit false (the fields are tri-state on purpose).
