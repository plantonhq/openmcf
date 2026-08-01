# Traces-gateway-to-Tempo preset

The gateway shape for traces: a two-replica deployment applications
push OTLP spans to (gRPC on 4317, HTTP on 4318 — the exported
`otlp_grpc_endpoint`/`otlp_http_endpoint` outputs are what they point
at), with the `memory_limiter` and `batch` processors in front of an
`otlp` gRPC exporter shipping to Tempo. The exporter's
`tls: {insecure: true}` is the in-cluster plaintext posture — both
ends live on the same cluster network.

PREREQUISITE: a `KubernetesOtelOperator` on the cluster. No custom
ServiceAccount is needed — this pipeline reads no cluster state, so
the operator's default account suffices.

The sizing is deliberate: the container's 512 Mi memory limit and the
`memory_limiter`'s 400 MiB limit (plus 100 MiB spike allowance) agree,
so under pressure the collector refuses data visibly — the sender
retries — instead of being OOM-killed silently.

Change first: the `otlp` exporter's endpoint host — compose it from
your Tempo install's OTLP gRPC endpoint output. Then `replicas` as
push volume grows, keeping the memory limit and `limit_mib` in step.

See [02-traces-gateway-to-tempo.yaml](./02-traces-gateway-to-tempo.yaml)
for the manifest.
