package resources

import (
	"testing"
)

func TestLoadCloudNativePGManifests(t *testing.T) {
	objs, err := LoadCloudNativePGManifests()
	if err != nil {
		t.Fatalf("failed to load manifests: %v", err)
	}

	if len(objs) == 0 {
		t.Fatal("expected at least one manifest object, got 0")
	}

	kindCounts := make(map[string]int)
	for _, obj := range objs {
		kindCounts[obj.GetKind()]++
	}

	if kindCounts["CustomResourceDefinition"] == 0 {
		t.Error("expected at least one CRD in the manifests")
	}
	if kindCounts["Deployment"] == 0 {
		t.Error("expected at least one Deployment in the manifests")
	}
	if kindCounts["ServiceAccount"] == 0 {
		t.Error("expected at least one ServiceAccount in the manifests")
	}
	if kindCounts["ClusterRole"] == 0 {
		t.Error("expected at least one ClusterRole in the manifests")
	}
	// The release ships admission webhooks -- whoever applies it needs the
	// webhookconfiguration RBAC, so a version bump that adds/removes them
	// should be a visible event.
	if kindCounts["MutatingWebhookConfiguration"] == 0 || kindCounts["ValidatingWebhookConfiguration"] == 0 {
		t.Error("expected both webhook configurations in the manifests")
	}

	t.Logf("Loaded %d manifest objects:", len(objs))
	for kind, count := range kindCounts {
		t.Logf("  %s: %d", kind, count)
	}
}

func TestLoadCloudNativePGManifests_AllHaveKindAndName(t *testing.T) {
	objs, err := LoadCloudNativePGManifests()
	if err != nil {
		t.Fatalf("failed to load manifests: %v", err)
	}

	for i, obj := range objs {
		if obj.GetKind() == "" {
			t.Errorf("object %d has empty kind", i)
		}
		if obj.GetName() == "" {
			t.Errorf("object %d (%s) has empty name", i, obj.GetKind())
		}
	}
}

// The Cluster CRD doubles as the detect-or-install probe, and the controller
// Deployment must exist under the name the component's gate waits on -- a
// CloudNativePG version bump that renames either must fail here, not in a
// forever-Pending component.
func TestLoadCloudNativePGManifests_ContainsDetectAndReadinessTargets(t *testing.T) {
	objs, err := LoadCloudNativePGManifests()
	if err != nil {
		t.Fatalf("failed to load manifests: %v", err)
	}

	foundCRD := false
	foundDeployment := false
	for _, obj := range objs {
		if obj.GetKind() == "CustomResourceDefinition" && obj.GetName() == "clusters.postgresql.cnpg.io" {
			foundCRD = true
		}
		if obj.GetKind() == "Deployment" && obj.GetName() == "cnpg-controller-manager" &&
			obj.GetNamespace() == "cnpg-system" {
			foundDeployment = true
		}
	}

	if !foundCRD {
		t.Error("expected to find the Cluster CRD (clusters.postgresql.cnpg.io) in manifests")
	}
	if !foundDeployment {
		t.Error("expected Deployment cnpg-controller-manager in cnpg-system -- the component's gate waits on it")
	}
}

// ---------------------------------------------------------------------------
// Tekton Pipelines manifests
// ---------------------------------------------------------------------------

func TestLoadTektonPipelinesManifests(t *testing.T) {
	objs, err := LoadTektonPipelinesManifests()
	if err != nil {
		t.Fatalf("failed to load manifests: %v", err)
	}

	if len(objs) == 0 {
		t.Fatal("expected at least one manifest object, got 0")
	}

	kindCounts := make(map[string]int)
	for _, obj := range objs {
		kindCounts[obj.GetKind()]++
	}

	if kindCounts["CustomResourceDefinition"] == 0 {
		t.Error("expected at least one CRD in the manifests")
	}
	if kindCounts["Deployment"] == 0 {
		t.Error("expected at least one Deployment in the manifests")
	}
	if kindCounts["Namespace"] == 0 {
		t.Error("expected the tekton-pipelines Namespace in the manifests")
	}
	if kindCounts["ClusterRole"] == 0 {
		t.Error("expected at least one ClusterRole in the manifests")
	}

	t.Logf("Loaded %d manifest objects:", len(objs))
	for kind, count := range kindCounts {
		t.Logf("  %s: %d", kind, count)
	}
}

func TestLoadTektonPipelinesManifests_AllHaveKindAndName(t *testing.T) {
	objs, err := LoadTektonPipelinesManifests()
	if err != nil {
		t.Fatalf("failed to load manifests: %v", err)
	}

	for i, obj := range objs {
		if obj.GetKind() == "" {
			t.Errorf("object %d has empty kind", i)
		}
		if obj.GetName() == "" {
			t.Errorf("object %d (%s) has empty name", i, obj.GetKind())
		}
	}
}

// The PipelineRun CRD doubles as the detect-or-install probe, and the two
// readiness-gate Deployments must exist under the names the component waits
// on -- a Tekton version bump that renames either must fail here, not in a
// forever-Pending component.
func TestLoadTektonPipelinesManifests_ContainsDetectAndReadinessTargets(t *testing.T) {
	objs, err := LoadTektonPipelinesManifests()
	if err != nil {
		t.Fatalf("failed to load manifests: %v", err)
	}

	foundCRD := false
	foundDeployments := map[string]bool{}
	for _, obj := range objs {
		if obj.GetKind() == "CustomResourceDefinition" && obj.GetName() == TektonPipelineRunCRDName {
			foundCRD = true
		}
		if obj.GetKind() == "Deployment" && obj.GetNamespace() == TektonPipelinesNamespace {
			foundDeployments[obj.GetName()] = true
		}
	}

	if !foundCRD {
		t.Errorf("expected to find the PipelineRun CRD (%s) in manifests", TektonPipelineRunCRDName)
	}
	for _, want := range []string{TektonControllerDeploymentName, TektonWebhookDeploymentName} {
		if !foundDeployments[want] {
			t.Errorf("expected Deployment %s in %s -- the component's readiness gate waits on it",
				want, TektonPipelinesNamespace)
		}
	}
}
