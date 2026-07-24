---
title: "Private mirror preset"
description: "The standard cluster-wide operator, but every image byte comes from your own registry: an `image` override pointing at the mirrored manager image, `image_pull_secrets` naming the registry credential,..."
type: "preset"
rank: "03"
presetSlug: "03-private-mirror"
componentSlug: "opensearch-operator"
componentTitle: "OpenSearch Operator"
provider: "kubernetes"
icon: "package"
order: 3
---

# Private mirror preset

The standard cluster-wide operator, but every image byte comes from
your own registry: an `image` override pointing at the mirrored
manager image, `image_pull_secrets` naming the registry credential,
and explicit `resources` for quota-governed namespaces. This is the
shape for air-gapped clusters and organizations that forbid direct
Docker Hub pulls.

Do not use this when the cluster can reach Docker Hub — the override
is one more thing to keep in sync with the pinned chart version for
no benefit. And note the mirror covers only the OPERATOR image here:
the OpenSearch node images that KubernetesOpenSearch clusters run are
declared on those resources (their own `image` field), not inherited
from this install.

Change the `image.repository` to your mirror path, create the pull
secret in the install namespace first, and keep `image.tag` aligned
with the chart's appVersion whenever you bump `chart_version` — a
mismatched pair is the classic silent-drift failure.

See [03-private-mirror.yaml](./03-private-mirror.yaml) for the
manifest.
