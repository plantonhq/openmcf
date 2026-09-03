// Package helmcrds derives the CustomResourceDefinition set a Helm-installed
// kind must own, from the exact version the user pinned, and speaks every
// failure in a form a person or an agent can act on.
//
// It is the engine-neutral half of the catalog's answer to Helm's CRD hole
// (Helm installs a chart's crds/ directory once and never upgrades or removes
// it; templated CRDs are release-owned and die with the release). Charts that
// neither own their CRD lifecycle nor publish a separate CRD chart land on the
// DERIVE branch of the decision tree: the module renders the pinned chart with
// its CRD switch on (or fetches the pinned upstream bundle), applies each CRD
// as a kept resource ahead of the release, and installs the release with CRDs
// skipped. Both engines consume this package's contract: the Pulumi twin calls
// it in-process; the Terraform twin's generated helm_crds.tf mirrors its
// semantics with the helm_template data source and quotes its stamp keys.
//
// Two invariants the package exists to keep:
//
//   - A module derives, never copies, what a pinned artifact carries. A
//     vendored CRD is frozen at one version while chart_version moves with the
//     user, so the schema silently stops matching the operator.
//   - The render sees the release's FULL merged values with only the CRD
//     switch flipped on. Templated CRDs routinely depend on other values (the
//     release's fullname, whether cert-manager injects the conversion CA, the
//     webhook service port); a render from a minimal value set produces CRDs
//     that point at the wrong webhook or freeze a self-signed CA.
//
// Every applied CRD is stamped with the source chart and version
// (AnnotationSourceChart, AnnotationSourceVersion) and a selection label
// (LabelSource). The stamps are how both engines re-adopt kept CRDs, how the
// never-downgrade check reads what a cluster already carries, and how an agent
// with kubectl can see where a CRD came from.
//
// Failures are Failure values: what was observed (with the value), what it
// most likely means (one root cause), and the exact next step. Text that names
// only a mechanism is a defect here.
package helmcrds
