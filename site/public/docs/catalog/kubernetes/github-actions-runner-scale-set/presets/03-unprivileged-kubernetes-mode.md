---
title: "Unprivileged Kubernetes-mode preset"
description: "Container jobs WITHOUT privileged pods: the Kubernetes container hook runs each job container as its own pod, so security-restricted clusters (Pod Security `restricted`, regulated environments) still..."
type: "preset"
rank: "03"
presetSlug: "03-unprivileged-kubernetes-mode"
componentSlug: "github-actions-runner-scale-set"
componentTitle: "GitHub Actions Runner Scale Set"
provider: "kubernetes"
icon: "package"
order: 3
---

# Unprivileged Kubernetes-mode preset

Container jobs WITHOUT privileged pods: the Kubernetes container hook
runs each job container as its own pod, so security-restricted
clusters (Pod Security `restricted`, regulated environments) still get
container-based workflows. Each runner claims an ephemeral work volume
from the declared StorageClass (a reference field — a
KubernetesStorageClass composes naturally) and releases it with the
runner.

Know the semantics trade against dind: jobs MUST run in a container
(`ACTIONS_RUNNER_REQUIRE_JOB_CONTAINER` is the hook's contract), and a
few Docker-specific actions that talk to a daemon socket do not apply.
For teams standardized on container jobs, that constraint is usually
already true.

Change first: `kubernetes_work_volume.size` when jobs check out large
repositories or produce big artifacts — the work volume is where the
whole job workspace lives.

See
[03-unprivileged-kubernetes-mode.yaml](./03-unprivileged-kubernetes-mode.yaml)
for the manifest.
