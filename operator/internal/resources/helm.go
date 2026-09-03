package resources

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
)

// RenderHelmChart loads a Helm chart from an in-memory .tgz archive, renders
// it with the given values, and returns the resulting Kubernetes objects as
// unstructured resources ready for Server-Side Apply.
//
// Rendering uses the Helm action package (equivalent to `helm template`) rather
// than the raw engine, so subchart dependency conditions defined in Chart.yaml
// (e.g., cassandra.enabled, postgresql.enabled) are correctly evaluated. Charts
// with disabled subcharts will not produce resources for those subcharts.
//
// The releaseName and namespace are used in chart template rendering (they
// populate .Release.Name and .Release.Namespace respectively). Only manifests
// that decode to valid Kubernetes objects are returned; empty documents, NOTES
// files, and Helm test resources are filtered out.
func RenderHelmChart(chartData []byte, releaseName, namespace string, values map[string]any) ([]*unstructured.Unstructured, error) {
	chart, err := loader.LoadArchive(bytes.NewReader(chartData))
	if err != nil {
		return nil, fmt.Errorf("loading helm chart archive: %w", err)
	}

	client := action.NewInstall(&action.Configuration{})
	client.DryRun = true
	client.ClientOnly = true
	client.ReleaseName = releaseName
	client.Namespace = namespace
	client.Replace = true
	client.IncludeCRDs = false
	client.KubeVersion = &chartutil.KubeVersion{
		Version: "v1.31.0",
		Major:   "1",
		Minor:   "31",
	}

	rel, err := client.Run(chart, values)
	if err != nil {
		return nil, fmt.Errorf("rendering helm chart templates: %w", err)
	}

	var allManifests []string
	if m := strings.TrimSpace(rel.Manifest); m != "" {
		allManifests = append(allManifests, m)
	}
	for _, hook := range rel.Hooks {
		if isTestHook(hook.Events) {
			continue
		}
		if m := strings.TrimSpace(hook.Manifest); m != "" {
			allManifests = append(allManifests, m)
		}
	}

	if len(allManifests) == 0 {
		return nil, nil
	}

	var result []*unstructured.Unstructured
	for _, manifest := range allManifests {
		objs, err := parseMultiDocYAML([]byte(manifest))
		if err != nil {
			return nil, fmt.Errorf("parsing rendered manifests: %w", err)
		}
		for _, obj := range objs {
			if obj.GetNamespace() == "" && isNamespacedKind(obj.GetKind()) {
				obj.SetNamespace(namespace)
			}
			result = append(result, obj)
		}
	}

	return result, nil
}

// shouldSkipRenderedFile returns true for Helm output files that should not
// be applied as Kubernetes resources (NOTES.txt, test manifests).
func shouldSkipRenderedFile(name string) bool {
	if strings.HasSuffix(name, "NOTES.txt") {
		return true
	}
	if strings.Contains(name, "/tests/") {
		return true
	}
	return false
}

// isTestHook returns true if the hook events include only "test" events.
// Test hooks should not be applied as regular resources.
func isTestHook(events []release.HookEvent) bool {
	return slices.Contains(events, release.HookTest)
}

// isNamespacedKind returns true for resource kinds that require a namespace.
// Cluster-scoped resources should not have a namespace set.
func isNamespacedKind(kind string) bool {
	switch kind {
	case "Namespace", "ClusterRole", "ClusterRoleBinding",
		"CustomResourceDefinition", "PriorityClass",
		"PersistentVolume", "StorageClass", "IngressClass":
		return false
	default:
		return true
	}
}
