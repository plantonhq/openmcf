package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestCommittedTreeValidates is the lint gate's core: the real committed
// skills/ and agents/ trees must satisfy the structure contract at all
// times. A truncated reference, an orphaned file, or a frontmatter/slug
// mismatch in a pull request fails here with the exact repair message.
func TestCommittedTreeValidates(t *testing.T) {
	root := repoRoot(t)
	if _, err := os.Stat(filepath.Join(root, "skills")); err != nil {
		// Under Bazel the sandbox carries only declared inputs, not the
		// repository tree; the enforcing lane for this gate is `go test`
		// from the repo root (the skills lint workflow), same as the
		// reference drift gate.
		t.Skipf("skills tree not present at %s (sandboxed run)", root)
	}
	tree, err := LoadTree(root)
	if err != nil {
		t.Fatalf("loading committed tree: %v", err)
	}
	for _, err := range Validate(tree) {
		t.Errorf("committed tree invalid: %v", err)
	}
	if len(tree.Skills) == 0 {
		t.Fatal("no skills loaded from the committed tree")
	}
}

// TestValidateCatchesBrokenTrees proves every guarded failure class is
// actually caught (a validator that passes everything would also pass the
// committed tree). Each case plants exactly one defect in an otherwise
// valid tree and asserts the violation is reported.
func TestValidateCatchesBrokenTrees(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(files map[string]string)
		wantErr string
	}{
		{
			name: "cited reference file is empty",
			mutate: func(files map[string]string) {
				files["skills/demo/references/topic.md"] = ""
			},
			wantErr: "is empty",
		},
		{
			name: "cited reference file is missing",
			mutate: func(files map[string]string) {
				delete(files, "skills/demo/references/topic.md")
			},
			wantErr: "does not exist",
		},
		{
			name: "reference file exists but is never cited",
			mutate: func(files map[string]string) {
				files["skills/demo/references/orphan.md"] = "never cited"
			},
			wantErr: "never cites it",
		},
		{
			name: "frontmatter name does not match directory slug",
			mutate: func(files map[string]string) {
				files["skills/demo/SKILL.md"] = strings.Replace(
					files["skills/demo/SKILL.md"], "name: demo", "name: other", 1)
			},
			wantErr: "must equal the directory slug",
		},
		{
			name: "agent instructions are empty",
			mutate: func(files map[string]string) {
				files["agents/helper/instructions.md"] = ""
			},
			wantErr: "instructions.md is missing or empty",
		},
		{
			name: "compat floor is not a version",
			mutate: func(files map[string]string) {
				files["skills/compat.yaml"] = "minimum_daemon_version: soon\nminimum_cli_version: v0.0.0\n"
			},
			wantErr: "is not a vX.Y.Z version",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := validFixtureTree()
			tc.mutate(files)
			root := writeTree(t, files)

			tree, err := LoadTree(root)
			if err != nil {
				t.Fatalf("loading fixture tree: %v", err)
			}
			errs := Validate(tree)
			if len(errs) == 0 {
				t.Fatalf("planted defect was not caught")
			}
			for _, e := range errs {
				if strings.Contains(e.Error(), tc.wantErr) {
					return
				}
			}
			t.Fatalf("no violation mentions %q; got: %v", tc.wantErr, errs)
		})
	}

	// The unmutated fixture must be clean, or the cases above prove nothing.
	t.Run("unmutated fixture is valid", func(t *testing.T) {
		tree, err := LoadTree(writeTree(t, validFixtureTree()))
		if err != nil {
			t.Fatalf("loading fixture tree: %v", err)
		}
		if errs := Validate(tree); len(errs) > 0 {
			t.Fatalf("fixture tree should be valid, got: %v", errs)
		}
	})
}

// TestBuildSkillArchiveDeterministic holds the archive to its contract:
// byte-identical across builds, stored entries only, descriptor-free
// (CreateRaw), SKILL.md at the archive root.
func TestBuildSkillArchiveDeterministic(t *testing.T) {
	tree, err := LoadTree(writeTree(t, validFixtureTree()))
	if err != nil {
		t.Fatalf("loading fixture tree: %v", err)
	}
	skill := tree.Skills[0]

	first, err := BuildSkillArchive(skill)
	if err != nil {
		t.Fatalf("first build: %v", err)
	}
	second, err := BuildSkillArchive(skill)
	if err != nil {
		t.Fatalf("second build: %v", err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("two builds of identical content produced different bytes")
	}

	reader, err := zip.NewReader(bytes.NewReader(first), int64(len(first)))
	if err != nil {
		t.Fatalf("reading archive back: %v", err)
	}
	var names []string
	for _, f := range reader.File {
		names = append(names, f.Name)
		if f.Method != zip.Store {
			t.Errorf("%s: entry is not stored", f.Name)
		}
		if f.Flags&0x8 != 0 {
			t.Errorf("%s: entry carries a data descriptor", f.Name)
		}
	}
	if names[0] != "SKILL.md" {
		t.Errorf("first entry is %s, want SKILL.md at the archive root", names[0])
	}
}

// TestPackageReleaseManifest packages the fixture tree and verifies the
// manifest's integrity claims against the files actually written: every
// artifact exists, its size matches, and its SHA-256 matches.
func TestPackageReleaseManifest(t *testing.T) {
	tree, err := LoadTree(writeTree(t, validFixtureTree()))
	if err != nil {
		t.Fatalf("loading fixture tree: %v", err)
	}
	outDir := t.TempDir()

	manifest, err := PackageRelease(tree, "v9.9.9", outDir)
	if err != nil {
		t.Fatalf("packaging: %v", err)
	}
	if manifest.Version != "v9.9.9" {
		t.Errorf("manifest version = %q", manifest.Version)
	}
	if manifest.Floor.MinimumDaemonVersion != "v0.0.0" {
		t.Errorf("manifest floor = %+v, want the fixture's compat.yaml values", manifest.Floor)
	}

	var reread Manifest
	raw, err := os.ReadFile(filepath.Join(outDir, ManifestFileName))
	if err != nil {
		t.Fatalf("manifest not written: %v", err)
	}
	if err := json.Unmarshal(raw, &reread); err != nil {
		t.Fatalf("manifest does not parse: %v", err)
	}

	for _, artifact := range append(reread.Skills, reread.Agents...) {
		content, err := os.ReadFile(filepath.Join(outDir, artifact.File))
		if err != nil {
			t.Errorf("%s: manifest lists a file that was not written: %v", artifact.File, err)
			continue
		}
		if int64(len(content)) != artifact.SizeBytes {
			t.Errorf("%s: size %d does not match manifest's %d", artifact.File, len(content), artifact.SizeBytes)
		}
		sum := sha256.Sum256(content)
		if hex.EncodeToString(sum[:]) != artifact.Sha256 {
			t.Errorf("%s: checksum does not match manifest", artifact.File)
		}
	}
	if len(reread.Skills) != 1 || len(reread.Agents) != 1 {
		t.Errorf("manifest lists %d skill(s) and %d agent(s), want 1 and 1", len(reread.Skills), len(reread.Agents))
	}
}

// validFixtureTree is a minimal tree satisfying the whole structure
// contract; broken-tree cases mutate exactly one thing at a time.
func validFixtureTree() map[string]string {
	return map[string]string{
		"skills/demo/SKILL.md":            "---\nname: demo\ndescription: A demo skill.\n---\n\n# Demo\n\nRead `references/topic.md` first.\n",
		"skills/demo/references/topic.md": "The topic, explained.\n",
		"skills/compat.yaml":              "minimum_daemon_version: v0.0.0\nminimum_cli_version: v0.0.0\n",
		"agents/helper/instructions.md":   "You are the helper.\n",
	}
}

func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for path, content := range files {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// Validate requires the agents/ dir to exist even in skill-only cases.
	if err := os.MkdirAll(filepath.Join(root, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate test file")
	}
	return filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(thisFile))))
}
