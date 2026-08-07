package certification

import (
	"os"
	"path/filepath"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/plantonhq/planton/pkg/catalogbundle"

	// This test assembles its descriptor set from the LINKED registry, and
	// only the served version reaches a binary through the kind registry --
	// the conversion-source version must be linked explicitly. (The real
	// bundle is buf-built from the proto tree and carries every version
	// regardless of linkage.)
	_ "github.com/plantonhq/planton/catalog/_test/testcloudresourcegeneric/v1alpha1"
)

// The catalog-as-data certification case: a built bundle must carry BOTH of
// the torture kind's versions -- the old one is the conversion source, and a
// runtime serving documents stored under it needs its schema from the same
// bundle. This is the groundwork multi-version serving stands on.
func TestCertify_BundleCarriesBothTortureVersions(t *testing.T) {
	dir := t.TempDir()

	var fds descriptorpb.FileDescriptorSet
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		fds.File = append(fds.File, protodesc.ToFileDescriptorProto(fd))
		return true
	})
	raw, err := proto.Marshal(&fds)
	if err != nil {
		t.Fatal(err)
	}
	descriptorsPath := filepath.Join(dir, "descriptors.binpb")
	if err := os.WriteFile(descriptorsPath, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	bundlePath := filepath.Join(dir, "catalog-bundle.zip")
	if _, err := catalogbundle.Build(catalogbundle.BuildInput{
		DescriptorSetPath: descriptorsPath,
		CatalogDir:        filepath.Join(TortureKindRoot(t), "..", ".."),
		OutputPath:        bundlePath,
	}); err != nil {
		t.Fatal(err)
	}
	bundle, err := catalogbundle.Load(bundlePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := catalogbundle.CheckConformance(bundle); err != nil {
		t.Fatal(err)
	}

	for _, version := range []string{"v1alpha1", "v1alpha2"} {
		name := protoreflect.FullName(
			"dev.planton._test.testcloudresourcegeneric." + version + ".TestCloudResourceGeneric")
		if _, err := bundle.Files.FindDescriptorByName(name); err != nil {
			t.Errorf("the bundle must carry %s (the %s schema) -- stored documents at that version cannot be served without it", name, version)
		}
	}

	if _, ok := bundle.ConversionSpecs()["conversions/_test/testcloudresourcegeneric/v1alpha1_to_v1alpha2.yaml"]; !ok {
		t.Error("the bundle must carry the bridge between the versions it serves")
	}
}
