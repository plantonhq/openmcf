package resources_test

import (
	"testing"

	"github.com/plantonhq/planton/operator/internal/resources"
)

func TestRenderHelmChart_OpenFGA(t *testing.T) {
	chartData := resources.LoadOpenFGAChart()
	if len(chartData) == 0 {
		t.Fatal("Embedded OpenFGA chart is empty")
	}

	values := resources.OpenFGAHelmValues("test", "default")
	objs, err := resources.RenderHelmChart(chartData, "test-openfga", "default", values)
	if err != nil {
		t.Fatalf("RenderHelmChart failed: %v", err)
	}

	if len(objs) == 0 {
		t.Fatal("Expected rendered objects, got none")
	}

	kindCounts := map[string]int{}
	for _, obj := range objs {
		kindCounts[obj.GetKind()]++
	}

	if kindCounts["Deployment"] < 1 {
		t.Errorf("Expected at least 1 Deployment, got %d", kindCounts["Deployment"])
	}
	if kindCounts["Service"] < 1 {
		t.Errorf("Expected at least 1 Service, got %d", kindCounts["Service"])
	}

	for _, obj := range objs {
		if obj.GetKind() == "Deployment" || obj.GetKind() == "Service" {
			if obj.GetNamespace() != "default" {
				t.Errorf("Expected namespace 'default' on %s %s, got %q",
					obj.GetKind(), obj.GetName(), obj.GetNamespace())
			}
		}
	}
}

func TestRenderHelmChart_Temporal(t *testing.T) {
	chartData := resources.LoadTemporalChart()
	if len(chartData) == 0 {
		t.Fatal("Embedded Temporal chart is empty")
	}

	values := resources.TemporalHelmValues("test", "default")
	objs, err := resources.RenderHelmChart(chartData, "test-temporal", "default", values)
	if err != nil {
		t.Fatalf("RenderHelmChart failed: %v", err)
	}

	if len(objs) == 0 {
		t.Fatal("Expected rendered objects, got none")
	}

	kindCounts := map[string]int{}
	for _, obj := range objs {
		kindCounts[obj.GetKind()]++
	}

	if kindCounts["Deployment"] < 4 {
		t.Errorf("Expected at least 4 Deployments (frontend, history, matching, worker), got %d", kindCounts["Deployment"])
	}
	if kindCounts["Service"] < 1 {
		t.Errorf("Expected at least 1 Service, got %d", kindCounts["Service"])
	}
	if kindCounts["Job"] < 1 {
		t.Errorf("Expected at least 1 Job (schema setup), got %d", kindCounts["Job"])
	}

	t.Logf("Temporal chart rendered %d objects: %v", len(objs), kindCounts)
}

func TestRenderHelmChart_EmptyChart(t *testing.T) {
	_, err := resources.RenderHelmChart([]byte{}, "test", "default", nil)
	if err == nil {
		t.Error("Expected error for empty chart data, got nil")
	}
}
