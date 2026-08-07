---
title: "Private mirror preset"
description: "The air-gap posture: every image the install pulls re-pointed at a private registry, with a pull secret for all of them. Beyond the usual mirror hygiene there is one image here that genuinely needs..."
type: "preset"
rank: "03"
presetSlug: "03-private-mirror"
componentSlug: "altinity-clickhouse-operator"
componentTitle: "Altinity ClickHouse Operator"
provider: "kubernetes"
icon: "package"
order: 3
---

# Private mirror preset

The air-gap posture: every image the install pulls re-pointed at a
private registry, with a pull secret for all of them. Beyond the
usual mirror hygiene there is one image here that genuinely needs
this treatment: the CRD hook's upstream default is
`bitnami/kubectl:latest`, which still pulls today but has been frozen
since Bitnami retired its public catalog — it will silently age, and
`latest` is exactly the tag you do not want resolved differently on
different nodes. Pinning your own kubectl build turns that rot risk
into a controlled dependency.

The operator image tag should track the chart version (they release
in lockstep — chart 0.27.2 runs operator 0.27.2); when you bump one,
bump both. The metrics-exporter sidecar image stays on chart defaults
here because most mirrored environments proxy Docker Hub; on a fully
air-gapped cluster re-point it through `helm_values`
(`metrics.image.repository`/`tag`) the same way.

The first thing to change is the mirror hostname and the pull secret
name — then, as always, the operator password.

See [03-private-mirror.yaml](./03-private-mirror.yaml) for the
manifest.
