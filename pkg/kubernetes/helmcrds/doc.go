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
// One ownership rule decides what the module applies. Helm has two CRD
// surfaces: the crds/ directory (installed once, never upgraded, skipped by
// skip_crds) and the templates (ordinary release resources). A typed kind
// supplies a CRDOverride that turns the chart's switch on for the render and
// pins it off for the release, so every CRD the render produces is the
// module's. A kind that supplies no override renders exactly what the release
// renders: the crds/ surface is the module's, and templated CRDs stay Helm's.
// Templated CRDs the chart marks with helm.sh/resource-policy: keep are the
// chart owning its lifecycle correctly and need nothing; templated CRDs
// without that mark would be deleted with the release along with every custom
// resource built on them, so they are refused unless the kind's Policy accepts
// them. Only the generic Helm kind, whose chart is arbitrary, ever meets the
// second case.
//
// Every applied CRD is stamped with the source chart and version
// (AnnotationSourceChart, AnnotationSourceVersion) and a source label
// (LabelSource). The stamps are the ownership mark: before either engine
// writes a CRD it reads the cluster's copy by name, and one read answers two
// questions. Is it ours (the label)? If not, someone else owns it and the
// module refuses with that owner named rather than take it over with
// server-side apply, because a takeover could lower a schema the module never
// stamped and cannot order. Is it newer (the version stamp)? If so, the module
// refuses the downgrade. A CRD it stamped earlier is re-adopted, and the
// deploy log says so. An agent with kubectl reads the same stamps.
//
// Failures are Failure values (pkg/failure, the repository-wide shape): what
// was observed (with the value), what it most likely means (one root cause),
// and the exact next step. Text that names only a mechanism is a defect here.
// The Helm texts are classified through pkg/failure's constructors so the
// Terraform twin, whose raw provider output the CLI and the harness explain
// with pkg/failure.Explain, ends in the same words.
package helmcrds
