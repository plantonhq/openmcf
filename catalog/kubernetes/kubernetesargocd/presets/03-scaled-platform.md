# Scaled platform preset

Argo CD sized for a real platform team: the API server and repo server
on autoscalers (the repo server is what saturates first — every sync
renders manifests there), the application controller given real
memory, the disposable cache upgraded to the three-node Sentinel arm,
and every component scraped by the monitoring stack. This is the shape
for hundreds of Applications and the humans watching them.

Two scheduling truths to respect: the HA Redis arm carries a required
anti-affinity, so this preset needs at least three schedulable worker
nodes before its pods leave Pending; and the HPAs meter the declared
CPU requests, so removing those requests silently freezes scaling.
The controller stays single-replica deliberately — its replicas shard
across target CLUSTERS, and sharding one cluster buys nothing.

Change first: pair with the team SSO preset's identity block (this
one focuses on capacity), and when the monitoring stack is not
deployed yet, drop `service_monitors_enabled` until it is — the
ServiceMonitor objects need the Prometheus Operator CRDs to admit.

See [03-scaled-platform.yaml](./03-scaled-platform.yaml) for the
manifest.
