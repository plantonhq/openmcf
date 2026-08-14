package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	stigmer "github.com/stigmer/stigmer/sdk/go/v3"
	skillv1 "github.com/stigmer/stigmer/sdk/go/v3/proto/ai/stigmer/agentic/skill/v1"
)

// skillPushTag mirrors the tag every other publish lane uses, so CI-pushed
// and operator-pushed versions are indistinguishable on the engine.
const skillPushTag = "latest"

// skillEngine is the slice of the SDK's skill client this tool uses; tests
// serve fakes behind the genuine SDK client, and this seam only narrows the
// surface, never reimplements it.
type skillEngine interface {
	Push(ctx context.Context, input *skillv1.PushSkillRequest) (*skillv1.Skill, error)
	GetByReference(ctx context.Context, ref stigmer.ResourceRef) (*skillv1.Skill, error)
	ListVersions(ctx context.Context, input *skillv1.ListSkillVersionsInput) (*skillv1.ListSkillVersionsResponse, error)
}

// manifest mirrors the defspack release manifest's skill-relevant shape.
// The checksum grammar is shared by construction: defspack records the
// SHA-256 of each archive's bytes, and the engine keys a version on the
// SHA-256 of the uploaded bytes, so equality IS the delivery proof.
type manifest struct {
	Version string `json:"version"`
	Skills  []struct {
		Slug   string `json:"slug"`
		File   string `json:"file"`
		Sha256 string `json:"sha256"`
	} `json:"skills"`
}

// run publishes (or, with verifyOnly, audits) every selected skill from a
// defspack output directory against the target org. Any divergence between
// the manifest, the local archives, and the engine's registered state is an
// error -- this tool's one job is keeping those three identical.
func run(ctx context.Context, engine skillEngine, dir, org, onlySkill string, verifyOnly bool, out io.Writer) error {
	raw, err := os.ReadFile(filepath.Join(dir, "definitions-manifest.json"))
	if err != nil {
		return fmt.Errorf("reading manifest (run defspack first): %w", err)
	}
	var m manifest
	if err := json.Unmarshal(raw, &m); err != nil {
		return fmt.Errorf("parsing manifest: %w", err)
	}

	selected := 0
	for _, entry := range m.Skills {
		if onlySkill != "" && entry.Slug != onlySkill {
			continue
		}
		selected++

		archive, err := os.ReadFile(filepath.Join(dir, entry.File))
		if err != nil {
			return fmt.Errorf("reading %s: %w", entry.File, err)
		}
		// A stale or hand-touched build directory must never publish: the
		// manifest's checksum is the authority the archive must reproduce.
		sum := sha256.Sum256(archive)
		localHash := hex.EncodeToString(sum[:])
		if localHash != entry.Sha256 {
			return fmt.Errorf("%s does not match the manifest checksum -- stale build directory? re-run defspack", entry.File)
		}

		if verifyOnly {
			if err := verifySkill(ctx, engine, org, entry.Slug, entry.Sha256, m.Version, out); err != nil {
				return err
			}
			continue
		}
		if err := publishSkill(ctx, engine, org, entry.Slug, archive, entry.Sha256, out); err != nil {
			return err
		}
	}

	if selected == 0 {
		if onlySkill != "" {
			return fmt.Errorf("skill %q is not in the manifest", onlySkill)
		}
		return fmt.Errorf("manifest lists no skills")
	}
	return nil
}

// verifySkill reports whether the engine's LATEST version of the skill is
// exactly the manifest's content. Absence and divergence are both errors:
// verify-only exists to answer "does the engine serve what this release
// ships?" with an exit code.
func verifySkill(ctx context.Context, engine skillEngine, org, slug, wantHash, version string, out io.Writer) error {
	skill, err := engine.GetByReference(ctx, stigmer.ResourceRef{Org: org, Slug: slug})
	switch {
	case stigmer.IsNotFound(err):
		return fmt.Errorf("skill %s/%s: absent on the engine (manifest %s)", org, slug, version)
	case err != nil:
		return fmt.Errorf("reading skill %s/%s: %w", org, slug, err)
	}
	got := skill.GetStatus().GetVersionHash()
	if got != wantHash {
		return fmt.Errorf("skill %s/%s: engine serves %.12s, manifest %s ships %.12s -- the engine is NOT on this release's content", org, slug, got, version, wantHash)
	}
	fmt.Fprintf(out, "skill %s/%s: engine matches the manifest (%.12s)\n", org, slug, wantHash)
	return nil
}

// publishSkill pushes the archive and proves delivery: the engine's
// registered version hash must equal the manifest checksum, and the
// version count reveals whether content actually changed (the engine
// registers no new version for identical bytes).
func publishSkill(ctx context.Context, engine skillEngine, org, slug string, archive []byte, wantHash string, out io.Writer) error {
	before, err := versionCount(ctx, engine, org, slug)
	if err != nil {
		return err
	}

	pushed, err := engine.Push(ctx, &skillv1.PushSkillRequest{
		Org:      org,
		Artifact: archive,
		Tag:      skillPushTag,
	})
	if err != nil {
		return fmt.Errorf("pushing skill %s/%s: %w", org, slug, err)
	}
	if got := pushed.GetStatus().GetVersionHash(); got != wantHash {
		return fmt.Errorf("skill %s/%s: engine registered %.12s, manifest says %.12s -- pushed bytes are not the release bytes", org, slug, got, wantHash)
	}

	after, err := versionCount(ctx, engine, org, slug)
	if err != nil {
		return err
	}
	if after > before {
		fmt.Fprintf(out, "skill %s/%s: new version registered (%d -> %d, %.12s)\n", org, slug, before, after, wantHash)
	} else {
		fmt.Fprintf(out, "skill %s/%s: content unchanged (no new version registered)\n", org, slug)
	}
	return nil
}

func versionCount(ctx context.Context, engine skillEngine, org, slug string) (int, error) {
	resp, err := engine.ListVersions(ctx, &skillv1.ListSkillVersionsInput{Org: org, Slug: slug})
	switch {
	case stigmer.IsNotFound(err):
		return 0, nil
	case err != nil:
		return 0, fmt.Errorf("listing skill %s/%s versions: %w", org, slug, err)
	}
	return int(resp.GetTotalCount()), nil
}
