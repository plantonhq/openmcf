# SigNoz for production

The destination posture: SigNoz runs the observability product; a
`KubernetesClickHouse` (with its `KubernetesAltinityOperator`) runs the
database — each with its own lifecycle, sizing and operations. Every
connection detail is a reference to that resource's outputs, the
database connection runs over verified TLS, alert email flows through
secret-safe SMTP, and the external URL makes every alert link resolve
to your real hostname. The ingestion collector autoscales with ingest
volume; the server is sized for rule evaluation over real traffic.

The ClickHouse side's contract (its own production preset carries the
sizing): SigNoz's tested version (`25.12.5` at chart 0.133.0 — older
servers fail the schema migrations), a keeper-backed coordination
service, TLS enabled on the
client port this preset points at, and the `signoz` user — no grants
(unrestricted) or a grant set including `GRANT CLUSTER ON *.*`, with
explicit `networks` (a networks-less user is fenced to the ClickHouse
pods and localhost; SigNoz's pods get what reads as a password
failure). Deploy SigNoz into the SAME namespace as the ClickHouse: the
password travels by secretKeyRef, which cannot cross namespaces (or
replicate the Secret when they must live apart).

**When to use:** real workloads reporting into a platform your team
operates — one ClickHouse, one SigNoz, clean separations.

**When to move on:** this IS the destination; scale by sizing the
ClickHouse resource and the collector bounds, not by changing shape.
