---
title: "Pipelines engine preset"
description: "Tekton as an embedded engine: the `lite` profile runs Pipelines alone — no Triggers, no Dashboard, no Chains — for platforms that create PipelineRuns programmatically and present their own UI. Every..."
type: "preset"
rank: "02"
presetSlug: "02-pipelines-engine"
componentSlug: "tekton"
componentTitle: "Tekton"
provider: "kubernetes"
icon: "package"
order: 2
---

# Pipelines engine preset

Tekton as an embedded engine: the `lite` profile runs Pipelines alone
— no Triggers, no Dashboard, no Chains — for platforms that create
PipelineRuns programmatically and present their own UI. Every run's
lifecycle event streams to one receiver through the CloudEvents sink
(one cluster-global URL by Tekton's design — fan out downstream if
tenants need separate streams), the internet-reaching resolvers are
off, and the controller runs two replicas sharded into two buckets —
Tekton's actual HA mechanism (replicas without buckets are idle
standbys, not capacity).

Aggressive retention fits the embedded shape: runs older than two days
are pruned every six hours, because the platform consuming the events
owns the durable history.

Change first: raise `performance.kube_api_qps`/`kube_api_burst`
together with run volume — the controller's API budget, not pod
resources, is the usual throughput ceiling.

See [02-pipelines-engine.yaml](./02-pipelines-engine.yaml) for the
manifest.
