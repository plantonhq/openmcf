package resources

import (
	"testing"
)

func TestValkeyHelmValues_FullnameOverride(t *testing.T) {
	vals := ValkeyHelmValues("my-planton", "1Gi", "")
	if vals["fullnameOverride"] != "my-planton-redis" {
		t.Errorf("expected fullnameOverride my-planton-redis, got %v", vals["fullnameOverride"])
	}
}

func TestValkeyHelmValues_Architecture(t *testing.T) {
	vals := ValkeyHelmValues("my-planton", "1Gi", "")
	if vals["architecture"] != "standalone" {
		t.Errorf("expected architecture standalone, got %v", vals["architecture"])
	}
}

func TestValkeyHelmValues_Auth(t *testing.T) {
	vals := ValkeyHelmValues("my-planton", "1Gi", "")
	auth, ok := vals["auth"].(map[string]any)
	if !ok {
		t.Fatal("expected auth to be a map")
	}
	if auth["existingSecret"] != "my-planton-redis-credentials" {
		t.Errorf("expected existingSecret my-planton-redis-credentials, got %v", auth["existingSecret"])
	}
	if auth["existingSecretPasswordKey"] != RedisSecretKey {
		t.Errorf("expected existingSecretPasswordKey %s, got %v", RedisSecretKey, auth["existingSecretPasswordKey"])
	}
}

func TestValkeyHelmValues_Persistence(t *testing.T) {
	vals := ValkeyHelmValues("my-planton", "5Gi", "")
	primary, ok := vals["primary"].(map[string]any)
	if !ok {
		t.Fatal("expected primary to be a map")
	}
	persistence, ok := primary["persistence"].(map[string]any)
	if !ok {
		t.Fatal("expected persistence to be a map")
	}
	if persistence["size"] != "5Gi" {
		t.Errorf("expected persistence size 5Gi, got %v", persistence["size"])
	}
	if persistence["enabled"] != true {
		t.Error("expected persistence enabled")
	}
}

func TestValkeyHelmValues_Image(t *testing.T) {
	// The image must stay pinned to the BSD-licensed Valkey engine; a drift
	// back to a Redis 8+ image would reintroduce the RSALv2/SSPLv1/AGPLv3
	// licensing family the platform deliberately avoids distributing or
	// directing self-hosted customers to run.
	vals := ValkeyHelmValues("my-planton", "1Gi", "")
	image, ok := vals["image"].(map[string]any)
	if !ok {
		t.Fatal("expected image to be a map")
	}
	if image["repository"] != "bitnamilegacy/valkey" {
		t.Errorf("expected repository bitnamilegacy/valkey, got %v", image["repository"])
	}
}

func TestRedisSecretName(t *testing.T) {
	name := RedisSecretName("my-planton")
	if name != "my-planton-redis-credentials" {
		t.Errorf("expected my-planton-redis-credentials, got %s", name)
	}
}

func TestRedisServiceHost(t *testing.T) {
	host := RedisServiceHost("my-planton", "planton-system")
	expected := "my-planton-redis-primary.planton-system.svc.cluster.local"
	if host != expected {
		t.Errorf("expected %s, got %s", expected, host)
	}
}

func TestValkeyHelmValues_ChartRendering(t *testing.T) {
	chartData := LoadValkeyChart()
	if len(chartData) == 0 {
		t.Fatal("Valkey chart data is empty")
	}

	values := ValkeyHelmValues("test", "1Gi", "")
	objs, err := RenderHelmChart(chartData, "test-redis", "default", values)
	if err != nil {
		t.Fatalf("failed to render Valkey chart: %v", err)
	}
	if len(objs) == 0 {
		t.Fatal("expected rendered objects, got none")
	}

	kinds := make(map[string]bool)
	var stsName string
	for _, obj := range objs {
		kinds[obj.GetKind()] = true
		if obj.GetKind() == "StatefulSet" {
			stsName = obj.GetName()
		}
	}

	if !kinds["StatefulSet"] {
		t.Error("expected StatefulSet in rendered objects")
	}
	if !kinds["Service"] {
		t.Error("expected Service in rendered objects")
	}
	// The readiness check and RedisServiceHost depend on the chart's
	// "primary" naming; a chart bump that renames the workload must fail here
	// rather than in a live cluster.
	if stsName != "test-redis-primary" {
		t.Errorf("expected StatefulSet test-redis-primary, got %s", stsName)
	}
}
