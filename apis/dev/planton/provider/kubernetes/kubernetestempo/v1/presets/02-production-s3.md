# Production Tempo on object storage

Trace blocks in an S3-compatible object store (an in-cluster
KubernetesSeaweedFs here; AWS S3, GCS or Azure by swapping the storage
block), with the metrics generator deriving service-graph and span metrics
from the trace stream and remote-writing them to the cluster's Prometheus
— the seam that lights up Grafana's service map.

**When to use:** production tracing at scale, and whenever you want the
service map / RED metrics in Grafana.

**Prerequisite:** the target Prometheus must accept remote-write
(`enable_remote_write_receiver: true` on the KubernetesKubePrometheusStack).

**Cross-cloud:** on EKS drop the endpoint, set a region, and leave
credentials empty for keyless IRSA.
