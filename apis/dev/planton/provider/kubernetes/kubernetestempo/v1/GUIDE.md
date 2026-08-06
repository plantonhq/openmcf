# KubernetesTempo Guide

The judgment this guide carries: Tempo's scaling floor is a storage
decision — one replica runs honestly on a local volume, but the moment
you need two, object storage stops being optional. Decide by trace
volume up front, not mid-incident.

## The replica/storage floor

`local` storage (the default) keeps trace blocks on a PersistentVolume —
honest for a single replica. MORE than one replica REQUIRES an object
storage backend; on-cluster that means composing a
[KubernetesSeaweedFs](../../kubernetesseaweedfs/v1/reference.md) and
pointing `storage.s3` at it (the validated shape is in this kind's own
example). Proposing `replicas: 2` with local storage is a manifest that
will not deploy. This kind deliberately models single-binary Tempo — by
the time per-component microservices are needed, that is its own design
conversation (the reference page says so).

## How traces arrive, how they are read

Applications send OTLP straight to the exported `otlp_grpc_endpoint` /
`otlp_http_endpoint`, or a
[KubernetesOtelCollector](../../kubernetesotelcollector/v1/GUIDE.md)
gateway sits in between (sampling, enrichment, fan-out). Grafana reads
traces back through a `tempo` datasource on the exported `http_endpoint`.
The full wired composition:
[observability-stack pattern](../../../patterns/observability-stack.md).

## Namespace ownership

Shares the observability namespace with its siblings — wire
`spec.namespace` through that KubernetesNamespace
([namespace-ownership pattern](../../../patterns/namespace-ownership.md)).

## Pairs well with

- KubernetesOtelCollector — the optional gateway in front of ingestion.
- KubernetesGrafana — the reader (`tempo` datasource).
- KubernetesSeaweedFs — the storage floor above one replica.
