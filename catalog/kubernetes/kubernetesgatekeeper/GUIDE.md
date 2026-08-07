# KubernetesGatekeeper Guide

The judgment this guide carries: Gatekeeper or Kyverno is the admission
policy decision — choose by policy LANGUAGE and lifecycle model, and run
ONE of them, not both.

## Gatekeeper vs Kyverno (the comparison home)

- **Gatekeeper (this kind)** — policies are ConstraintTemplates (Rego or
  Kubernetes CEL) instantiated by Constraints: the OPA ecosystem's
  power, at the cost of a policy language to learn. The chart OWNS its
  webhook configurations as release objects (clean uninstall), and the
  policy webhook defaults FAIL-OPEN — an engine outage never blocks
  admission (the namespace-label guard webhook is the deliberate
  exception; the reference page has both).
- **[KubernetesKyverno](../kuberneteskyverno/GUIDE.md)** — policies
  are plain Kubernetes resources (no new language), and the engine also
  MUTATES, GENERATES, and CLEANS UP resources, not just validates. Its
  webhooks are registered at RUNTIME per installed policies — more
  dynamic, with the messier uninstall story its guide covers.

Choose Gatekeeper when the organization already speaks Rego/OPA or wants
constraint-framework portability; choose Kyverno when
policies-as-plain-YAML and mutation/generation matter more. Running both
means two admission webhooks judging every request — never propose that
without a stated split.

## Engine here, policies beside it

This installs the ENGINE (webhook + audit controllers + CRDs).
ConstraintTemplates and Constraints are authored separately — via
KubernetesManifest or GitOps — the same engine-plus-declarations shape as
KEDA and the Spark operator. An installed Gatekeeper with no Constraints
enforces nothing.

## Once per cluster; the infra exception

Admission engines are cluster singletons in practice; shared-cluster chart, own
namespace ([namespace-ownership pattern](../../_patterns/namespace-ownership.md)
sole-tenant case;
[operator-prerequisite pattern](../../_patterns/operator-prerequisite.md)
for the CRD seam).

## Pairs well with

- KubernetesManifest — carries the ConstraintTemplates/Constraints today.
- KubernetesKyverno — the alternative lane (comparison above).
