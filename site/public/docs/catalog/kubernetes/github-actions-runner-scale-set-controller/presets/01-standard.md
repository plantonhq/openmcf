---
title: "Standard controller preset"
description: "One cluster-wide controller in its own namespace — the shape almost every cluster wants. It installs the runner CRDs and the manager; runner fleets are declared separately (one..."
type: "preset"
rank: "01"
presetSlug: "01-standard"
componentSlug: "github-actions-runner-scale-set-controller"
componentTitle: "GitHub Actions Runner Scale Set Controller"
provider: "kubernetes"
icon: "package"
order: 1
---

# Standard controller preset

One cluster-wide controller in its own namespace — the shape almost
every cluster wants. It installs the runner CRDs and the manager;
runner fleets are declared separately (one KubernetesGhaRunnerScaleSet
per repository/organization registration), and this controller serves
all of them.

Know the destroy contract: the CRDs delete with the controller, which
cascade-deletes every runner scale set on the cluster — destroy the
scale sets first.

Change first: `flags.log_level: info` once the install is settled (the
chart ships debug — noisy), and metrics when Prometheus scrapes the
cluster.

See [01-standard.yaml](./01-standard.yaml) for the manifest.
