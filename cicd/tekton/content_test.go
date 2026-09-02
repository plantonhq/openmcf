package tekton

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// The recorded digest holds Version and the embedded tree together: editing
// any embedded YAML without bumping Version (and re-recording the digest)
// fails the gate. Two different content states must never share a pin --
// the pin is what makes stamped provenance and byte-identical reruns honest.
func TestContent_versionNamesExactlyOneContentState(t *testing.T) {
	computed := computeDigest(t)
	digestPath := filepath.Join("testdata", "content-digest.txt")
	recorded, readErr := os.ReadFile(digestPath)
	want := strings.TrimSpace(string(recorded))
	if readErr != nil || computed != want {
		if os.Getenv("UPDATE_CONTENT_DIGEST") != "" {
			if err := os.WriteFile(digestPath, []byte(computed+"\n"), 0o644); err != nil {
				t.Fatalf("recording the new digest: %v", err)
			}
			t.Logf("content digest recorded as %s -- verify Version was bumped for this content change", computed)
			return
		}
		t.Fatalf("the embedded content changed but the recorded digest did not: computed %s, recorded %q.\nBump Version in content.go (two content states must never share a pin), then re-record with UPDATE_CONTENT_DIGEST=1 go test ./cicd/tekton/", computed, want)
	}
}

func TestContent_tracksAndTasksArePresent(t *testing.T) {
	wantTracks := []string{"buildpacks", "cloudflare-worker", "dockerfile"}
	if got := Tracks(); !reflect.DeepEqual(got, wantTracks) {
		t.Fatalf("expected exactly the service build tracks %v, got %v", wantTracks, got)
	}
	tasks, err := TaskFiles()
	if err != nil {
		t.Fatal(err)
	}
	for _, stem := range []string{"git-clone", "kustomize-build", "buildkit", "buildpacks"} {
		if _, ok := tasks[stem]; !ok {
			t.Fatalf("embedded task %s.yaml is missing", stem)
		}
	}
	if _, ok := Track("no-such-track"); ok {
		t.Fatal("an unknown track must not resolve")
	}
}

// Every image the content can run is digest-pinned: the tag documents
// intent, the digest is what the cluster pulls, so a build is reproducible
// and immune to upstream tag mutation -- and the derived allowlist an
// air-gapped cluster mirrors is immutable.
func TestImages_everyImageIsDigestPinned(t *testing.T) {
	images, err := Images()
	if err != nil {
		t.Fatal(err)
	}
	if len(images) == 0 {
		t.Fatal("the content runs no images at all -- the derivation is broken")
	}
	for _, image := range images {
		if !strings.Contains(image, "@sha256:") {
			t.Errorf("image %q is not digest-pinned", image)
		}
		if strings.Contains(image, "$(") {
			t.Errorf("image %q is an unresolved param reference", image)
		}
	}
}

// The derived image set is recorded so an image can never enter or leave
// silently: a content change that alters what the cluster pulls shows up
// here as a reviewed diff, beside the Version bump the digest test forces.
func TestImages_matchTheRecordedSet(t *testing.T) {
	images, err := Images()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"amazon/aws-cli:latest@sha256:310813a7eae8fd88da1cc9c37970e3500b0ff3984479e1012f0a6fd44e453f63",
		"bitnami/kubectl:latest@sha256:c62a62db80e777acdee87f76bc6f06a95239ad2ff210bf78f585e39e33da98e2",
		"bitnamilegacy/kubectl:latest@sha256:cd354d5b25562b195b277125439c23e4046902d7f1abc0dc3c75aad04d298c17",
		"docker.io/library/bash:5.1.4@sha256:b208215a4655538be652b2769d82e576bc4d0a2bb132144c060efc5be8c3f5d6",
		"ghcr.io/tektoncd-catalog/git-clone:v1.1.0@sha256:b7fe6c370322586feb555c807f3fae7ca5d62c20ebbcca987114e69366151957",
		"ghcr.io/tektoncd/github.com/tektoncd/pipeline/cmd/git-init:v0.45.0@sha256:8ab0f58d8381b0b71f5b2bae1f63522989d739e3154d8cab1bacfa0ef5317214",
		"moby/buildkit:v0.23.2-rootless@sha256:cab936745de5d673465948f1e93ff4d6e372bbe33f218afd3314eba45a6f85a9",
		"node:20-alpine@sha256:fb4cd12c85ee03686f6af5362a0b0d56d50c58a04632e6c0fb8363f609372293",
		"paketobuildpacks/builder-jammy-base:latest@sha256:f93da4e8abc73ab3555793d3a992b724ba7d0baafffeddb1219cf5c433fcf3b2",
	}
	if !reflect.DeepEqual(images, want) {
		t.Fatalf("the derived image set changed.\n got: %s\nwant: %s\nIf the content change is intended, update this list (and Version).", strings.Join(images, "\n      "), strings.Join(want, "\n      "))
	}
}

func computeDigest(t *testing.T) string {
	t.Helper()
	files, err := allFiles()
	if err != nil {
		t.Fatalf("walking embedded content: %v", err)
	}
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	h := sha256.New()
	for _, path := range paths {
		h.Write([]byte(path))
		h.Write([]byte{0})
		h.Write(files[path])
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}
