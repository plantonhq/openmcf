# KubernetesPlantonOperator Guide

The judgment this guide carries: one Planton operator serves every
Planton platform on the cluster — the manager is a cluster-owner install,
the platforms are per-team declarations, and the coupling that actually
bites is DESTROY ORDERING, not installation.

## Once per cluster, platforms by the many

Install the operator once; teams declare KubernetesPlantonPlatform
resources — each platform gets its own namespace, its own URL, its own
identity server and databases, all reconciled by this one manager (it
watches every namespace). The operator polices its own singleton-ness at
startup: a second install refuses to start and says why, so the failure
mode of a mistake is a crash-looping pod with a remedy in its log, never
two managers fighting.

Two honest shared facts ride the topology: the cluster has ONE
`plantonplatforms.planton.ai` CRD schema (so one operator version defines
the spec surface for every platform, while each platform still pins its
own application version), and Tekton's cluster-wide build-events sink
means builds can feed only one build-enabled platform per cluster.

## The CRD is the load-bearing artifact

The chart ships its CRD in Helm's install-once `crds/` directory — never
upgraded, never removed. The modules deliberately take that lifecycle
over: the CRD is applied from a staged copy extracted from the pinned
published chart, kept on uninstall (a destroyed operator strands no
platform declarations), adopted on reinstall, and upgraded exactly when
`chart_version` moves. When pinning a NEWER chart than the module
default, know that the staged CRD stays at the default pin — the catalog
release that bumps the default re-stages the CRD with it, which is the
supported way to move.

## Destroy is safe by construction

The operator's own destroy strands nothing: the CRD survives (kept,
module-owned), declarations survive, running platforms keep serving —
merely unmanaged until the next install adopts them. Even platform
DELETION completes without the operator, because platform teardown is
Kubernetes garbage collection of the declaration's owner-referenced
objects, not operator work. Prefer platforms-then-operator as hygiene
(the operator publishes status while platforms drain), never as a
requirement.

## On the diagram

The operator draws one explicit `depends_on` edge FROM each
KubernetesPlantonPlatform — no spec field consumes an operator output
(the coupling is the cluster-global CRD contract), so composed charts
declare the edge in metadata. The `planton-on-kubernetes` infra-chart
carries namespace + operator + platform as one deployable arm.
