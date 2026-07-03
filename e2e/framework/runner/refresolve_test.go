package runner

import (
	"os"
	"path/filepath"
	"testing"

	awssubnetv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awssubnet/v1"
	gcpbackendservicev1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpbackendservice/v1"
	gcpurlmapv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpurlmap/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/internal/manifest"
)

const subnetManifestWithRef = `apiVersion: aws.planton.dev/v1
kind: AwsSubnet
metadata:
  name: ref-subnet
spec:
  region: us-west-2
  vpcId:
    valueFrom:
      kind: AwsVpc
      name: my-vpc
      fieldPath: status.outputs.vpc_id
  availabilityZone: us-west-2a
  cidrBlock: 10.0.1.0/24
`

func writeTempManifest(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "manifest.yaml")
	if err := os.WriteFile(path, []byte(body), 0600); err != nil {
		t.Fatalf("failed to write temp manifest: %v", err)
	}
	return path
}

func TestResolveManifestRefs_ResolvesVpcIdFromPrerequisite(t *testing.T) {
	manifestPath := writeTempManifest(t, subnetManifestWithRef)

	depOutputs := map[cloudresourcekind.CloudResourceKind]map[string]interface{}{
		cloudresourcekind.CloudResourceKind_AwsVpc: {
			"vpc_id":   "vpc-resolved123",
			"vpc_cidr": "10.0.0.0/16",
		},
	}

	resolvedPath, err := ResolveManifestRefs(manifestPath, depOutputs)
	if err != nil {
		t.Fatalf("ResolveManifestRefs failed: %v", err)
	}
	if resolvedPath == manifestPath {
		t.Fatal("expected a new resolved manifest path, got the original")
	}

	obj, err := manifest.LoadManifest(resolvedPath)
	if err != nil {
		t.Fatalf("failed to load resolved manifest: %v", err)
	}
	subnet, ok := obj.(*awssubnetv1.AwsSubnet)
	if !ok {
		t.Fatalf("resolved manifest is not an AwsSubnet: %T", obj)
	}
	if got := subnet.GetSpec().GetVpcId().GetValue(); got != "vpc-resolved123" {
		t.Errorf("vpc_id value = %q, want %q", got, "vpc-resolved123")
	}
	if subnet.GetSpec().GetVpcId().GetValueFrom() != nil {
		t.Error("vpc_id should be a literal after resolution, but value_from is still set")
	}
}

func TestResolveManifestRefs_NoDependenciesReturnsOriginal(t *testing.T) {
	manifestPath := writeTempManifest(t, subnetManifestWithRef)

	resolvedPath, err := ResolveManifestRefs(manifestPath, nil)
	if err != nil {
		t.Fatalf("ResolveManifestRefs failed: %v", err)
	}
	if resolvedPath != manifestPath {
		t.Errorf("expected original path when there are no dependencies, got %q", resolvedPath)
	}
}

func TestResolveManifestRefs_MissingOutputErrors(t *testing.T) {
	manifestPath := writeTempManifest(t, subnetManifestWithRef)

	depOutputs := map[cloudresourcekind.CloudResourceKind]map[string]interface{}{
		cloudresourcekind.CloudResourceKind_AwsVpc: {
			"vpc_cidr": "10.0.0.0/16", // vpc_id intentionally absent
		},
	}

	if _, err := ResolveManifestRefs(manifestPath, depOutputs); err == nil {
		t.Fatal("expected an error when the prerequisite output is missing, got nil")
	}
}

const backendServiceNegManifest = `apiVersion: gcp.planton.dev/v1
kind: GcpBackendService
metadata:
  name: neg-backend
spec:
  protocol: HTTP
  loadBalancingScheme: EXTERNAL_MANAGED
  backends:
    - group:
        valueFrom:
          kind: GcpRegionNetworkEndpointGroup
          name: my-neg
          fieldPath: status.outputs.self_link
      balancingMode: RATE
      maxRate: 100
`

const urlMapExplicitKindManifest = `apiVersion: gcp.planton.dev/v1
kind: GcpUrlMap
metadata:
  name: my-url-map
spec:
  defaultService:
    valueFrom:
      kind: GcpBackendService
      name: my-backend
      fieldPath: status.outputs.self_link
`

func TestResolveManifestRefs_ResolvesNestedBackendGroupByExplicitKind(t *testing.T) {
	manifestPath := writeTempManifest(t, backendServiceNegManifest)

	depOutputs := map[cloudresourcekind.CloudResourceKind]map[string]interface{}{
		cloudresourcekind.CloudResourceKind_GcpRegionNetworkEndpointGroup: {
			"self_link": "https://www.googleapis.com/compute/v1/projects/p/regions/r/networkEndpointGroups/n",
		},
	}

	resolvedPath, err := ResolveManifestRefs(manifestPath, depOutputs)
	if err != nil {
		t.Fatalf("ResolveManifestRefs failed: %v", err)
	}

	obj, err := manifest.LoadManifest(resolvedPath)
	if err != nil {
		t.Fatalf("failed to load resolved manifest: %v", err)
	}
	bs, ok := obj.(*gcpbackendservicev1.GcpBackendService)
	if !ok {
		t.Fatalf("resolved manifest is not a GcpBackendService: %T", obj)
	}
	if len(bs.GetSpec().GetBackends()) != 1 {
		t.Fatalf("expected 1 backend, got %d", len(bs.GetSpec().GetBackends()))
	}
	got := bs.GetSpec().GetBackends()[0].GetGroup().GetValue()
	want := "https://www.googleapis.com/compute/v1/projects/p/regions/r/networkEndpointGroups/n"
	if got != want {
		t.Errorf("backends[0].group = %q, want %q", got, want)
	}
}

func TestResolveManifestRefs_ResolvesExplicitValueFromKindWithoutDefaultKind(t *testing.T) {
	manifestPath := writeTempManifest(t, urlMapExplicitKindManifest)

	depOutputs := map[cloudresourcekind.CloudResourceKind]map[string]interface{}{
		cloudresourcekind.CloudResourceKind_GcpBackendService: {
			"self_link": "https://www.googleapis.com/compute/v1/projects/p/global/backendServices/bs",
		},
	}

	resolvedPath, err := ResolveManifestRefs(manifestPath, depOutputs)
	if err != nil {
		t.Fatalf("ResolveManifestRefs failed: %v", err)
	}

	obj, err := manifest.LoadManifest(resolvedPath)
	if err != nil {
		t.Fatalf("failed to load resolved manifest: %v", err)
	}
	um, ok := obj.(*gcpurlmapv1.GcpUrlMap)
	if !ok {
		t.Fatalf("resolved manifest is not a GcpUrlMap: %T", obj)
	}
	got := um.GetSpec().GetDefaultService().GetValue()
	want := "https://www.googleapis.com/compute/v1/projects/p/global/backendServices/bs"
	if got != want {
		t.Errorf("default_service = %q, want %q", got, want)
	}
}
