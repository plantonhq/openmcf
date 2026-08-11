# Azure Traffic Manager Endpoint -- Operational Guide

Judgment calls that matter when you run Traffic Manager endpoints in production.

## Set priorities explicitly, or creation order owns your failover plan

Unset priority means Azure assigns the next free value in creation order -- fine until someone recreates an endpoint and it silently moves to the back of the failover line. On Priority-routed profiles, give every endpoint an explicit value with gaps (10, 20, 30) so inserting a tier later never renumbers the plan.

## Drain with enabled: false; reserve always-serve for broken probes

`enabled: false` is the maintenance drain: the endpoint leaves DNS answers while its configuration stays. `always_serve_enabled` is the opposite tool -- it DISABLES health checking and keeps the endpoint in answers even when probes fail, for targets probes cannot reach (locked-down networks, probe-hostile firewalls). Never leave always-serve on a target you expect health-based failover to protect; it opts that endpoint out of failover entirely.

## Azure endpoints follow the resource; external endpoints follow the string

An azure endpoint tracks its target resource -- if the Public IP's address changes, Traffic Manager follows automatically, and its region feeds Performance routing with no location field at all. An external endpoint is a frozen string: retargeting is an edit, and Performance routing needs the explicit `endpointLocation` you gave it. Prefer azure endpoints for Azure resources; use external only for what genuinely lives elsewhere.

## Geographic claims are exclusive and validated live

Every geographic code (WORLD, GEO-EU, country codes) must be claimed by exactly ONE endpoint in the profile -- Azure rejects overlaps at apply time against its live hierarchy. Claim WORLD on a catch-all endpoint so unmatched callers get an answer; a Geographic profile with unclaimed regions returns NXDOMAIN for those callers.

## Nested trees: the child floor is your blast-radius dial

`minimumChildEndpoints` decides when a whole child profile counts as down: with a floor of 1, the parent keeps sending traffic to a region running on its last healthy instance; with a floor near the child's endpoint count, one instance failure fails the region over entirely. Set it to the child's genuine minimum serving capacity. The child profile must not use MultiValue routing (Azure enforces the composition rules).

## Endpoint names are per-type namespaces

Uniqueness is per (profile, endpoint type): an azure endpoint and an external endpoint may share a name. Do not lean on that -- give endpoints names that say what they front ("eastus-pip", "onprem-dc2"), because the portal's monitor view and the ARM ids read much better for it.
