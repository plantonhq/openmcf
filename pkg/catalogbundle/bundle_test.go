package catalogbundle

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// The full bundle lifecycle, hermetically: a descriptor set built from the
// compiled-in registry + the real provider tree -> Build -> Load (checksum
// verification) -> CheckConformance across every registered kind. The CI
// lane runs the same verify against a buf-built descriptor set -- proving
// the proto toolchain and the compiled registry agree.
func TestBundleRoundTrip(t *testing.T) {
	dir := t.TempDir()
	descriptorsPath := filepath.Join(dir, "descriptors.binpb")
	writeLinkedDescriptorSet(t, descriptorsPath)

	bundlePath := filepath.Join(dir, "catalog-bundle.zip")
	manifest, err := Build(BuildInput{
		DescriptorSetPath: descriptorsPath,
		ProviderBaseDir:   providerBaseDir(t),
		ReleaseTag:        "vTEST",
		OutputPath:        bundlePath,
	})
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}
	if len(manifest.Checksums) < 3 {
		t.Fatalf("suspiciously small bundle: %d entries", len(manifest.Checksums))
	}

	bundle, err := Load(bundlePath)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if err := CheckConformance(bundle); err != nil {
		t.Fatalf("conformance failed: %v", err)
	}

	// The torture kind's cargo must be aboard: its conversion spec and its
	// served version's preset.
	if _, ok := bundle.ConversionSpecs()["conversions/_test/testcloudresourcegeneric/v1alpha1_to_v1alpha2.yaml"]; !ok {
		t.Error("the torture kind's conversion spec is missing from the bundle")
	}
	foundPreset := false
	for name := range bundle.Presets() {
		if strings.HasPrefix(name, "presets/_test/testcloudresourcegeneric/") {
			foundPreset = true
			break
		}
	}
	if !foundPreset {
		t.Error("the torture kind's presets are missing from the bundle")
	}
}

// A corrupted entry must fail the load whole -- consumers never run on
// partially verified data.
func TestLoadRefusesCorruption(t *testing.T) {
	dir := t.TempDir()
	descriptorsPath := filepath.Join(dir, "descriptors.binpb")
	writeLinkedDescriptorSet(t, descriptorsPath)

	bundlePath := filepath.Join(dir, "catalog-bundle.zip")
	if _, err := Build(BuildInput{
		DescriptorSetPath: descriptorsPath,
		ProviderBaseDir:   providerBaseDir(t),
		OutputPath:        bundlePath,
	}); err != nil {
		t.Fatal(err)
	}

	// Re-zip the bundle with one entry's CONTENT altered while the manifest
	// keeps the original checksum -- exactly the tamper/mis-build class the
	// self-verification exists for.
	corruptedPath := filepath.Join(dir, "corrupted.zip")
	rezipWithAlteredEntry(t, bundlePath, corruptedPath, "descriptors.binpb")

	if _, err := Load(corruptedPath); err == nil {
		t.Fatal("a bundle whose entry does not match its recorded checksum must be refused")
	}
}

func rezipWithAlteredEntry(t *testing.T, src, dst, alterEntry string) {
	t.Helper()
	reader, err := zip.OpenReader(src)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	out, err := os.Create(dst)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	zw := zip.NewWriter(out)
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			t.Fatal(err)
		}
		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatal(err)
		}
		if file.Name == alterEntry {
			content[0] ^= 0xFF
		}
		w, err := zw.Create(file.Name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(content); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
}

func writeLinkedDescriptorSet(t *testing.T, path string) {
	t.Helper()
	var fds descriptorpb.FileDescriptorSet
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		fds.File = append(fds.File, protodesc.ToFileDescriptorProto(fd))
		return true
	})
	raw, err := proto.Marshal(&fds)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}
}

func providerBaseDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve caller location")
	}
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "apis", "dev", "planton", "provider")
}
