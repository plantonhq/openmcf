package main

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	stigmer "github.com/stigmer/stigmer/sdk/go/v3"
	skillv1 "github.com/stigmer/stigmer/sdk/go/v3/proto/ai/stigmer/agentic/skill/v1"
	apiresource "github.com/stigmer/stigmer/sdk/go/v3/proto/ai/stigmer/commons/apiresource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The fake serves the real gRPC controller interfaces behind the genuine
// SDK client, so the tests exercise the exact marshaling path the live
// engine sees. Its content-addressing mirrors the engine's: a version is
// keyed on the SHA-256 of the uploaded archive bytes, and pushing the hash
// the engine already holds registers no new version.
type fakeEngine struct {
	mu sync.Mutex
	// latestHash and versionCount per slug; a missing slug is NotFound.
	latestHash   map[string]string
	versionCount map[string]int
	// When set, Push reports THIS hash instead of the artifact's real one
	// (the delivery-corruption arm the parity assertion must catch).
	lieOnPushHash string
}

func newFakeEngine() *fakeEngine {
	return &fakeEngine{latestHash: map[string]string{}, versionCount: map[string]int{}}
}

// archiveSlug derives the pushed archive's skill slug the way the engine
// does: from the SKILL.md frontmatter name at the archive root.
func archiveSlug(data []byte) string {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "unreadable-archive"
	}
	rc, err := reader.Open("SKILL.md")
	if err != nil {
		return "missing-skill-md"
	}
	defer func() { _ = rc.Close() }()
	content, err := io.ReadAll(rc)
	if err != nil {
		return "unreadable-skill-md"
	}
	for _, line := range strings.Split(string(content), "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "name:"); ok {
			return strings.TrimSpace(rest)
		}
	}
	return "nameless-skill"
}

type fakeSkillCommand struct {
	skillv1.UnimplementedSkillCommandControllerServer
	rec *fakeEngine
}

func (s *fakeSkillCommand) Push(_ context.Context, in *skillv1.PushSkillRequest) (*skillv1.Skill, error) {
	slug := archiveSlug(in.GetArtifact())
	sum := sha256.Sum256(in.GetArtifact())
	hash := hex.EncodeToString(sum[:])

	s.rec.mu.Lock()
	defer s.rec.mu.Unlock()
	if s.rec.latestHash[slug] != hash {
		s.rec.latestHash[slug] = hash
		s.rec.versionCount[slug]++
	}
	reported := hash
	if s.rec.lieOnPushHash != "" {
		reported = s.rec.lieOnPushHash
	}
	return &skillv1.Skill{Status: &skillv1.SkillStatus{VersionHash: reported}}, nil
}

type fakeSkillQuery struct {
	skillv1.UnimplementedSkillQueryControllerServer
	rec *fakeEngine
}

func (s *fakeSkillQuery) GetByReference(_ context.Context, ref *apiresource.ApiResourceReference) (*skillv1.Skill, error) {
	s.rec.mu.Lock()
	defer s.rec.mu.Unlock()
	hash, ok := s.rec.latestHash[ref.GetSlug()]
	if !ok {
		return nil, status.Error(codes.NotFound, "skill not found")
	}
	return &skillv1.Skill{Status: &skillv1.SkillStatus{VersionHash: hash}}, nil
}

func (s *fakeSkillQuery) ListVersions(_ context.Context, in *skillv1.ListSkillVersionsInput) (*skillv1.ListSkillVersionsResponse, error) {
	s.rec.mu.Lock()
	defer s.rec.mu.Unlock()
	count, ok := s.rec.versionCount[in.GetSlug()]
	if !ok {
		return nil, status.Error(codes.NotFound, "skill not found")
	}
	return &skillv1.ListSkillVersionsResponse{TotalCount: int32(count)}, nil
}

func newTestEngine(t *testing.T, fake *fakeEngine) skillEngine {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	srv := grpc.NewServer()
	skillv1.RegisterSkillCommandControllerServer(srv, &fakeSkillCommand{rec: fake})
	skillv1.RegisterSkillQueryControllerServer(srv, &fakeSkillQuery{rec: fake})
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	client, err := stigmer.NewClient(stigmer.WithBaseURL(lis.Addr().String()), stigmer.WithInsecure())
	if err != nil {
		t.Fatalf("client: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })
	return client.Skill
}

// skillZip builds a minimal skill archive whose frontmatter carries the
// slug (the engine's identity derivation) plus a body payload so distinct
// contents hash differently.
func skillZip(t *testing.T, slug, body string) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	f, err := w.Create("SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fmt.Fprintf(f, "---\nname: %s\n---\n\n%s\n", slug, body); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// writeReleaseDir lays out a defspack-shaped output directory: one zip per
// skill plus the manifest recording each archive's true checksum.
func writeReleaseDir(t *testing.T, skills map[string][]byte) string {
	t.Helper()
	dir := t.TempDir()
	m := manifest{Version: "v9.9.9"}
	for slug, archive := range skills {
		file := "skill-" + slug + ".zip"
		if err := os.WriteFile(filepath.Join(dir, file), archive, 0o644); err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(archive)
		m.Skills = append(m.Skills, struct {
			Slug   string `json:"slug"`
			File   string `json:"file"`
			Sha256 string `json:"sha256"`
		}{Slug: slug, File: file, Sha256: hex.EncodeToString(sum[:])})
	}
	raw, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "definitions-manifest.json"), raw, 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestPublishRegistersAndReportsNewVersion(t *testing.T) {
	fake := newFakeEngine()
	engine := newTestEngine(t, fake)
	dir := writeReleaseDir(t, map[string][]byte{
		"alpha": skillZip(t, "alpha", "content one"),
		"beta":  skillZip(t, "beta", "content two"),
	})

	var out bytes.Buffer
	if err := run(context.Background(), engine, dir, "acme", "", false, &out); err != nil {
		t.Fatalf("run: %v", err)
	}
	for _, want := range []string{
		"skill acme/alpha: new version registered (0 -> 1",
		"skill acme/beta: new version registered (0 -> 1",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output missing %q; got:\n%s", want, out.String())
		}
	}
}

func TestPublishUnchangedContentIsNoOp(t *testing.T) {
	fake := newFakeEngine()
	engine := newTestEngine(t, fake)
	dir := writeReleaseDir(t, map[string][]byte{"alpha": skillZip(t, "alpha", "same content")})

	if err := run(context.Background(), engine, dir, "acme", "", false, io.Discard); err != nil {
		t.Fatalf("first run: %v", err)
	}
	var out bytes.Buffer
	if err := run(context.Background(), engine, dir, "acme", "", false, &out); err != nil {
		t.Fatalf("second run: %v", err)
	}
	if !strings.Contains(out.String(), "content unchanged (no new version registered)") {
		t.Errorf("second push of identical bytes should be a no-op; got:\n%s", out.String())
	}
	if fake.versionCount["alpha"] != 1 {
		t.Errorf("version count = %d, want 1", fake.versionCount["alpha"])
	}
}

func TestPublishFailsWhenEngineHashDiverges(t *testing.T) {
	fake := newFakeEngine()
	fake.lieOnPushHash = "not-the-manifest-hash"
	engine := newTestEngine(t, fake)
	dir := writeReleaseDir(t, map[string][]byte{"alpha": skillZip(t, "alpha", "content")})

	err := run(context.Background(), engine, dir, "acme", "", false, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "not the release bytes") {
		t.Fatalf("hash divergence not caught, err = %v", err)
	}
}

func TestPublishRefusesStaleBuildDir(t *testing.T) {
	fake := newFakeEngine()
	engine := newTestEngine(t, fake)
	dir := writeReleaseDir(t, map[string][]byte{"alpha": skillZip(t, "alpha", "content")})
	// Corrupt the archive after the manifest recorded its checksum.
	if err := os.WriteFile(filepath.Join(dir, "skill-alpha.zip"), skillZip(t, "alpha", "tampered"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := run(context.Background(), engine, dir, "acme", "", false, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "stale build directory") {
		t.Fatalf("stale build dir not caught, err = %v", err)
	}
	if fake.versionCount["alpha"] != 0 {
		t.Error("a stale archive must never be pushed")
	}
}

func TestVerifyOnlyBothVerdicts(t *testing.T) {
	fake := newFakeEngine()
	engine := newTestEngine(t, fake)
	archive := skillZip(t, "alpha", "content")
	dir := writeReleaseDir(t, map[string][]byte{"alpha": archive})

	// Absent skill: verify must fail.
	err := run(context.Background(), engine, dir, "acme", "", true, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "absent on the engine") {
		t.Fatalf("absent skill not reported, err = %v", err)
	}

	// After a push the engine matches: verify must pass.
	if err := run(context.Background(), engine, dir, "acme", "", false, io.Discard); err != nil {
		t.Fatalf("push: %v", err)
	}
	var out bytes.Buffer
	if err := run(context.Background(), engine, dir, "acme", "", true, &out); err != nil {
		t.Fatalf("verify after push: %v", err)
	}
	if !strings.Contains(out.String(), "engine matches the manifest") {
		t.Errorf("match not reported; got:\n%s", out.String())
	}

	// Engine moves ahead (someone pushed different content): verify must fail.
	fake.mu.Lock()
	fake.latestHash["alpha"] = "somethingelse"
	fake.mu.Unlock()
	err = run(context.Background(), engine, dir, "acme", "", true, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "NOT on this release's content") {
		t.Fatalf("divergence not reported, err = %v", err)
	}
}

func TestSkillFilter(t *testing.T) {
	fake := newFakeEngine()
	engine := newTestEngine(t, fake)
	dir := writeReleaseDir(t, map[string][]byte{
		"alpha": skillZip(t, "alpha", "one"),
		"beta":  skillZip(t, "beta", "two"),
	})

	if err := run(context.Background(), engine, dir, "acme", "beta", false, io.Discard); err != nil {
		t.Fatalf("run: %v", err)
	}
	if fake.versionCount["beta"] != 1 || fake.versionCount["alpha"] != 0 {
		t.Errorf("filter leaked: alpha=%d beta=%d", fake.versionCount["alpha"], fake.versionCount["beta"])
	}

	err := run(context.Background(), engine, dir, "acme", "nosuch", false, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "not in the manifest") {
		t.Fatalf("unknown slug not caught, err = %v", err)
	}
}
