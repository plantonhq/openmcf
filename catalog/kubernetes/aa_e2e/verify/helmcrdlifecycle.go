package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// Stamp keys the modules write on every derived CRD; the source of truth is
// pkg/kubernetes/helmcrds, spelled out here so the verifiers read exactly
// what a user would read with kubectl.
const (
	crdSourceVersionAnnotation = "planton.ai/crd-source-version"
	crdSourceLabel             = "planton.ai/crd-source"
)

// The failure classes the CRD-lifecycle lanes pin. Each names one of the
// primitive's three-part refusals; the verifier asserts the engine's error
// carries all three parts, so a mechanism-only message can never pass.
const (
	// The pinned chart version is not published (the index read, or Helm's
	// own locate error, whichever the engine surfaces first).
	failureChartVersionNotPublished = "chart-version-not-published"
	// The manifest lowers the chart version below what the cluster's CRDs
	// were derived from.
	failureCrdSchemaDowngrade = "crd-schema-downgrade"
	// The chart templates CRDs as release resources without Helm's keep
	// mark and the spec did not accept Helm-managed CRDs.
	failureHelmManagedCrds = "helm-managed-crds"
)

// helmCRDLifecycle is the part of a Helm-based kind's verification that the
// catalog's CRD primitive promises identically for every kind on it: the
// module-owned CRDs are Established and stamped with the pinned chart
// version; destroy keeps or deletes them as the manifest declared; a refused
// deploy or upgrade explains itself in three parts. A kind's verifier embeds
// it and adds the workload checks that are its own.
type helmCRDLifecycle struct {
	// CRDs are the module-owned CRDs the lane expects, by metadata.name.
	CRDs []string
	// ChartVersion is the version the manifest pins (or the kind's default),
	// the value every derived CRD must carry in its source-version stamp.
	ChartVersion string
	// SourceLabel is the value of the planton.ai/crd-source label the module
	// stamps: the chart name.
	SourceLabel string
	// KeepOnUninstall is the manifest's crds.keepOnUninstall (default true):
	// whether destroy must leave the CRDs or take them.
	KeepOnUninstall bool
	// InstallCrds is the manifest's crds.install (default true).
	InstallCrds bool
	// DeployRefused marks a scenario whose deploy is DESIGNED to be refused
	// (the expect-deploy-failure lane): nothing was created, so the destroy
	// assertion has no keep or delete to check.
	DeployRefused bool
}

// readHelmCRDLifecycle reads the lifecycle-relevant fields from a scenario
// manifest so one verifier serves the plain, upgrade, cleanup, reinstall and
// refusal lanes. Scenario files in the catalog are written in both
// proto-JSON key cases; both are read.
func readHelmCRDLifecycle(manifestPath, defaultChartVersion, sourceLabel string, crds []string) helmCRDLifecycle {
	l := helmCRDLifecycle{
		CRDs:            crds,
		ChartVersion:    defaultChartVersion,
		SourceLabel:     sourceLabel,
		KeepOnUninstall: true,
		InstallCrds:     true,
	}
	spec := manifestSpecMap(manifestPath)
	if version, _ := specField(spec, "chartVersion").(string); version != "" {
		l.ChartVersion = version
	}
	if crdsBlock, ok := specField(spec, "crds").(map[string]interface{}); ok {
		if keep, ok := specField(crdsBlock, "keepOnUninstall").(bool); ok {
			l.KeepOnUninstall = keep
		}
		if install, ok := specField(crdsBlock, "install").(bool); ok {
			l.InstallCrds = install
		}
	}
	l.DeployRefused = manifestAnnotation(manifestPath, "planton.dev/e2e-expect-deploy-failure") != ""
	return l
}

// specField reads one key from a spec map in either proto-JSON case: the
// camelCase name given, or its snake_case form. Both appear in the catalog's
// scenario files, and a verifier that read one case alone would silently
// take its defaults on the other.
func specField(m map[string]interface{}, camel string) interface{} {
	if m == nil {
		return nil
	}
	if v, ok := m[camel]; ok {
		return v
	}
	return m[snakeCase(camel)]
}

func snakeCase(camel string) string {
	var b strings.Builder
	for i, r := range camel {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// verifyEstablishedAndStamped waits for every expected CRD to be Established
// and, when the module owns them, asserts each one's source stamp: the
// version it was derived at and the chart it came from. After an upgrade
// the stamp must have moved to the new version; after a reinstall it must
// be present on the re-adopted CRD.
func (l helmCRDLifecycle) verifyEstablishedAndStamped(ctx context.Context, kubeconfig string) error {
	for _, crd := range l.CRDs {
		if out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"wait", "--for=condition=Established", "crd/"+crd, "--timeout=120s").CombinedOutput(); err != nil {
			return errors.Wrapf(err, "CRD %q never became Established: %s", crd, firstLines(string(out), 3))
		}
		if !l.InstallCrds {
			continue
		}
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "crd", crd, "-o",
			fmt.Sprintf("jsonpath={.metadata.annotations.%s}|{.metadata.labels.%s}",
				strings.ReplaceAll(crdSourceVersionAnnotation, ".", "\\."),
				strings.ReplaceAll(crdSourceLabel, ".", "\\."))).CombinedOutput()
		if err != nil {
			return errors.Wrapf(err, "reading the source stamp of CRD %q: %s", crd, firstLines(string(out), 2))
		}
		parts := strings.SplitN(strings.TrimSpace(string(out)), "|", 2)
		if len(parts) != 2 || parts[0] != l.ChartVersion {
			return errors.Errorf("CRD %q carries source version %q, the manifest pins %q -- the derived CRDs did not follow the chart version", crd, parts[0], l.ChartVersion)
		}
		if parts[1] != l.SourceLabel {
			return errors.Errorf("CRD %q carries source label %q, expected %q (the chart name)", crd, parts[1], l.SourceLabel)
		}
	}
	if len(l.CRDs) > 0 {
		fmt.Printf("  [verify] all %d module-owned CRDs Established and stamped with chart %s\n", len(l.CRDs), l.ChartVersion)
	}
	return nil
}

// verifyDestroyPosture asserts what a destroy must have done to the CRDs:
// kept them (the default), deleted them (keepOnUninstall: false), or nothing
// to assert (bring-your-own CRDs, or a deploy that was refused before
// anything was created).
func (l helmCRDLifecycle) verifyDestroyPosture(ctx context.Context, kubeconfig string) error {
	switch {
	case len(l.CRDs) == 0:
		return nil
	case !l.InstallCrds:
		fmt.Printf("  [verify] DESTROY: bring-your-own CRDs, nothing to assert about them\n")
		return nil
	case l.DeployRefused:
		fmt.Printf("  [verify] DESTROY: nothing to tear down -- the deploy was refused before anything was created\n")
		return nil
	case l.KeepOnUninstall:
		// The designed keep: the module-owned CRDs survive the destroy so an
		// uninstall never cascade-deletes the custom resources built on them.
		for _, crd := range l.CRDs {
			if err := KubectlResourceExists(ctx, kubeconfig, "crd", crd, ""); err != nil {
				return errors.Wrapf(err, "CRD %q was DELETED on destroy -- the module-owned keep posture broke", crd)
			}
		}
		fmt.Printf("  [verify] DESTROY: all %d CRDs RETAINED by design\n", len(l.CRDs))
		return nil
	default:
		// crds.keepOnUninstall: false -- the destroy must take the CRDs with it.
		for _, crd := range l.CRDs {
			if err := KubectlResourceAbsent(ctx, kubeconfig, "crd", crd, ""); err != nil {
				return errors.Wrapf(err, "CRD %q survived a destroy that declared keepOnUninstall: false", crd)
			}
		}
		fmt.Printf("  [verify] DESTROY: all %d CRDs DELETED as the manifest asked\n", len(l.CRDs))
		return nil
	}
}

// verifyRefusal pins a refused deploy or upgrade to the primitive's
// three-part refusal. The texts are the primitive's own (see
// pkg/kubernetes/helmcrds for the Go side and the generated helm_crds.tf for
// the Terraform side); asserting all three parts is what keeps "count
// mismatch"-class messages from ever passing again.
func (l helmCRDLifecycle) verifyRefusal(ctx context.Context, kubeconfig, expectation string, deployErr error) error {
	// The engines wrap long diagnostics at terminal width and the Terraform
	// runner joins lines with " | ", so a sentence can be split anywhere;
	// collapse all whitespace and separators before matching phrases.
	text := strings.Join(strings.Fields(strings.ReplaceAll(deployErr.Error(), "|", " ")), " ")
	for _, part := range []string{"observed:", "meaning:", "next step:"} {
		if !strings.Contains(text, part) {
			return errors.Errorf("the refusal lacks its %q part -- every CRD-lifecycle failure explains itself in three parts; got: %s", part, firstLines(text, 12))
		}
	}
	switch expectation {
	case failureChartVersionNotPublished:
		// The render or install could not locate the pinned chart; the
		// engines surface Helm's own text inside the observation.
		if !strings.Contains(text, l.ChartVersion) || !(strings.Contains(text, "not found") || strings.Contains(text, "not in the index")) {
			return errors.Errorf("expected the version-not-published refusal naming %s; got: %s", l.ChartVersion, firstLines(text, 12))
		}
		fmt.Printf("  [verify] REFUSED as expected: chart version %s is not published, and the message says what to do\n", l.ChartVersion)
	case failureCrdSchemaDowngrade:
		if !strings.Contains(text, "derived from chart version") || !strings.Contains(text, "asks for chart version "+l.ChartVersion) {
			return errors.Errorf("expected the schema-downgrade refusal naming the cluster's version and %s; got: %s", l.ChartVersion, firstLines(text, 12))
		}
		if !strings.Contains(text, "kubectl delete crd") {
			return errors.Errorf("the downgrade refusal must name the deliberate remedy; got: %s", firstLines(text, 12))
		}
		// Nothing was touched: the CRDs still carry the higher version.
		for _, crd := range l.CRDs {
			out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig, "get", "crd", crd, "-o",
				fmt.Sprintf("jsonpath={.metadata.annotations.%s}", strings.ReplaceAll(crdSourceVersionAnnotation, ".", "\\."))).CombinedOutput()
			if err != nil {
				return errors.Wrapf(err, "reading CRD %q after the refused downgrade", crd)
			}
			if strings.TrimSpace(string(out)) == l.ChartVersion {
				return errors.Errorf("CRD %q was downgraded to %s despite the refusal", crd, l.ChartVersion)
			}
		}
		fmt.Printf("  [verify] REFUSED as expected: the schema downgrade to %s was stopped before any CRD changed\n", l.ChartVersion)
	case failureHelmManagedCrds:
		if !strings.Contains(text, "templates") || !strings.Contains(text, "helm.sh/resource-policy") || !strings.Contains(text, "allow_helm_managed") {
			return errors.Errorf("expected the Helm-managed-CRDs refusal naming the keep mark and the dial; got: %s", firstLines(text, 12))
		}
		for _, crd := range l.CRDs {
			if !strings.Contains(text, crd) {
				return errors.Errorf("the Helm-managed-CRDs refusal must name CRD %q; got: %s", crd, firstLines(text, 12))
			}
		}
		fmt.Printf("  [verify] REFUSED as expected: the chart's Helm-managed CRDs were named, with the remedies\n")
	default:
		return errors.Errorf("unknown expected failure class %q", expectation)
	}
	return nil
}
