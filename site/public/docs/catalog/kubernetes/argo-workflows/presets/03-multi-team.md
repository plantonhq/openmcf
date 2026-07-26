---
title: "Multi-team preset"
description: "One cluster, several workflow engines, zero crosstalk. The instance ID is the mechanism: this controller reconciles only Workflows carrying its `controller-instanceid` label and ignores the rest, so..."
type: "preset"
rank: "03"
presetSlug: "03-multi-team"
componentSlug: "argo-workflows"
componentTitle: "Argo Workflows"
provider: "kubernetes"
icon: "package"
order: 3
---

# Multi-team preset

One cluster, several workflow engines, zero crosstalk. The instance ID
is the mechanism: this controller reconciles only Workflows carrying
its `controller-instanceid` label and ignores the rest, so team A's
engine, team B's engine and a platform engine coexist without stealing
each other's runs. The workflow namespaces list then places the runner
identity — ServiceAccount plus RBAC — into the team's app namespaces,
so workflows run where the team's services live without a second
install.

The submission contract that comes with instancing: every Workflow
this engine should run MUST carry the instance-ID label (the UI and
CLI add it automatically when pointed at this server; raw CR authors
add it themselves). An unlabeled Workflow in these namespaces sits
Pending forever — by design, because some other engine may own it.

Change first: pair each team's runner ServiceAccount with its own
cloud identity (IRSA/workload-identity annotations via helm_values on
the SA, or KubernetesServiceAccount composition) so team A's pipelines
can never sign requests as team B — the parallelism caps are fairness,
the identities are the actual wall.

See [03-multi-team.yaml](./03-multi-team.yaml) for the manifest.
