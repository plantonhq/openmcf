package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// ManifestFileName is the manifest's filename inside a definitions release.
const ManifestFileName = "definitions-manifest.json"

// Manifest describes one definitions release: which artifacts it carries,
// their integrity checksums, and the compatibility floor consumers must
// honor. Consumers verify an artifact's SHA-256 after download and refuse
// the release entirely when they sit below the floor.
type Manifest struct {
	Version string      `json:"version"`
	Floor   CompatFloor `json:"compatibilityFloor"`
	Skills  []Artifact  `json:"skills"`
	Agents  []Artifact  `json:"agents"`
	// omitempty keeps older-shaped manifests byte-identical for releases
	// carrying no automations; consumers treat the absent field as zero
	// automations rather than an error.
	Automations []Artifact `json:"automations,omitempty"`
}

// Artifact is one downloadable file in a definitions release. Skills ship
// as zip archives (SKILL.md at the archive root, references/ beside it);
// agents ship their instructions as a bare markdown file; automations ship
// their definition as a bare YAML file.
type Artifact struct {
	Slug      string `json:"slug"`
	File      string `json:"file"`
	Sha256    string `json:"sha256"`
	SizeBytes int64  `json:"sizeBytes"`
}

// zipEpoch is the fixed modification time stamped on every archive entry.
// Archive bytes must be a pure function of the tree's content so a rebuild
// of unchanged content reproduces the recorded checksum exactly; wall-clock
// timestamps would break that. 1980-01-01 is the earliest time the ZIP
// format's MS-DOS timestamp can represent.
var zipEpoch = time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)

// skillEntries is the ONE definition of what a skill's published content
// is, keyed by archive-relative path. Both carriers of that content -- the
// zip archive (BuildSkillArchive) and the browse manifest (see browse.go)
// -- compose from this seam, which is what makes "an exploded file is the
// archive entry, byte for byte" a structural fact rather than a promise:
// the release workflow derives the browsable exploded tree by unzipping
// the archives and verifies it against the browse manifest before upload.
func skillEntries(skill Skill) map[string][]byte {
	entries := map[string][]byte{"SKILL.md": skill.SkillMD}
	for name, content := range skill.ReferenceFiles {
		entries["references/"+name] = content
	}
	for path, content := range skill.PackFiles {
		entries[path] = content
	}
	return entries
}

// BuildSkillArchive packs one skill into a deterministic zip: lexically
// ordered entries, the fixed timestamp, stored (uncompressed) content, and
// CreateRaw with pre-computed sizes -- descriptor-free entries are the
// maximally compatible ZIP shape, and storing keeps the bytes immune even
// to compressor changes across Go releases. Storing (not deflating) also
// keeps the archive bytes a pure function of file content, which is what
// lets every consumer -- the release manifest, the daemon's re-pack
// verifier, and the serving engine's content-addressed storage -- agree on
// one checksum for one content state.
func BuildSkillArchive(skill Skill) ([]byte, error) {
	entries := skillEntries(skill)
	var paths []string
	for path := range entries {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for _, path := range paths {
		content := entries[path]
		entry, err := w.CreateRaw(&zip.FileHeader{
			Name:               path,
			Method:             zip.Store,
			Modified:           zipEpoch,
			CRC32:              crc32.ChecksumIEEE(content),
			CompressedSize64:   uint64(len(content)),
			UncompressedSize64: uint64(len(content)),
		})
		if err != nil {
			return nil, fmt.Errorf("creating archive entry %s: %w", path, err)
		}
		if _, err := entry.Write(content); err != nil {
			return nil, fmt.Errorf("writing archive entry %s: %w", path, err)
		}
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("finalizing archive for %s: %w", skill.Slug, err)
	}
	return buf.Bytes(), nil
}

// PackageRelease writes every release artifact plus the manifest into
// outDir and returns the manifest. The caller is expected to have run
// Validate first; packaging an invalid tree is a programming error.
func PackageRelease(tree *Tree, version, outDir string) (*Manifest, error) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating output directory: %w", err)
	}

	manifest := &Manifest{Version: version, Floor: tree.Floor}

	for _, skill := range tree.Skills {
		archive, err := BuildSkillArchive(skill)
		if err != nil {
			return nil, err
		}
		fileName := fmt.Sprintf("skill-%s.zip", skill.Slug)
		if err := os.WriteFile(filepath.Join(outDir, fileName), archive, 0o644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", fileName, err)
		}
		manifest.Skills = append(manifest.Skills, describe(skill.Slug, fileName, archive))
	}

	for _, agent := range tree.Agents {
		fileName := fmt.Sprintf("agent-%s-instructions.md", agent.Slug)
		if err := os.WriteFile(filepath.Join(outDir, fileName), agent.Instructions, 0o644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", fileName, err)
		}
		manifest.Agents = append(manifest.Agents, describe(agent.Slug, fileName, agent.Instructions))
	}

	for _, automation := range tree.Automations {
		fileName := fmt.Sprintf("automation-%s.yaml", automation.Slug)
		if err := os.WriteFile(filepath.Join(outDir, fileName), automation.Content, 0o644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", fileName, err)
		}
		manifest.Automations = append(manifest.Automations, describe(automation.Slug, fileName, automation.Content))
	}

	manifestJSON, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encoding manifest: %w", err)
	}
	manifestJSON = append(manifestJSON, '\n')
	if err := os.WriteFile(filepath.Join(outDir, ManifestFileName), manifestJSON, 0o644); err != nil {
		return nil, fmt.Errorf("writing %s: %w", ManifestFileName, err)
	}

	browseJSON, err := EncodeBrowseManifest(BuildBrowseManifest(tree, version))
	if err != nil {
		return nil, fmt.Errorf("encoding browse manifest: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outDir, BrowseManifestFileName), browseJSON, 0o644); err != nil {
		return nil, fmt.Errorf("writing %s: %w", BrowseManifestFileName, err)
	}

	return manifest, nil
}

func describe(slug, fileName string, content []byte) Artifact {
	sum := sha256.Sum256(content)
	return Artifact{
		Slug:      slug,
		File:      fileName,
		Sha256:    hex.EncodeToString(sum[:]),
		SizeBytes: int64(len(content)),
	}
}
