---
title: "Air-gapped mirror preset"
description: "The air-gap posture: every image this install — and the clusters it will create — pulls, re-pointed at a private registry. Three image surfaces travel together:"
type: "preset"
rank: "03"
presetSlug: "03-air-gapped-mirror"
componentSlug: "rabbitmq-cluster-operator"
componentTitle: "RabbitMQ Cluster Operator"
provider: "kubernetes"
icon: "package"
order: 3
---

# Air-gapped mirror preset

The air-gap posture: every image this install — and the clusters it
will create — pulls, re-pointed at a private registry. Three image
surfaces travel together:

- **`operator_image`** is the operator itself; empty, the release
  manifest's pinned `ghcr.io/rabbitmq/cluster-operator:2.22.3` is
  used, so keep the mirror's tag at the pinned release.
- **`default_rabbitmq_image`** is fleet-wide: every RabbitmqCluster
  that does not pin its own image runs this one (the operator's
  `DEFAULT_RABBITMQ_IMAGE`). The `-management` variant is REQUIRED —
  the operator's generated configuration expects the management
  plugin — so mirror `rabbitmq:4.2.6-management`, not the bare server
  image.
- **`default_user_updater_image`** is the credential-updater sidecar,
  consulted only by KubernetesRabbitMq resources using the Vault
  secret backend (upstream default
  `ghcr.io/rabbitmq/default-user-credential-updater:1.0.14`). Mirror
  it anyway — a Vault-backed cluster declared later should not be the
  moment the air gap is discovered.

`image_pull_secrets` names Secrets in the fixed `rabbitmq-system`
namespace; the `pull_secret_name` on the operator image points at the
same Secret here.

The first thing to change is the mirror hostname and the pull-secret
name. The tags should stay put: they are the pinned release's own.

See [03-air-gapped-mirror.yaml](./03-air-gapped-mirror.yaml) for the
manifest.
