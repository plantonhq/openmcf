package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Tree is the loaded definitions content: every skill under skills/, every
// agent under agents/, and the compatibility floor from skills/compat.yaml.
type Tree struct {
	Root   string
	Skills []Skill
	Agents []Agent
	Floor  CompatFloor
}

// Skill is one skills/<slug>/ directory: SKILL.md plus its reference files.
type Skill struct {
	Slug            string
	Dir             string
	FrontmatterName string
	Description     string
	SkillMD         []byte
	// ReferenceFiles maps a references/-relative name (e.g. "chart-format.md")
	// to its content, for every file present on disk.
	ReferenceFiles map[string][]byte
	// CitedReferences lists the references/<name>.md tokens SKILL.md cites,
	// deduplicated and sorted.
	CitedReferences []string
	// PackFiles carries content assembled from OUTSIDE the skill directory
	// at package time, keyed by archive path. Today only the catalog skill
	// uses it (the component reference pack under components/ -- see
	// catalogpack.go); the git tree never duplicates these files.
	PackFiles map[string][]byte
}

// Agent is one agents/<slug>/ directory carrying the agent's instructions.
type Agent struct {
	Slug         string
	Instructions []byte
}

// CompatFloor is the minimum consumer version set the current content
// assumes (see skills/compat.yaml for the raising discipline).
type CompatFloor struct {
	MinimumDaemonVersion string `yaml:"minimum_daemon_version" json:"minimumDaemonVersion"`
	MinimumCliVersion    string `yaml:"minimum_cli_version" json:"minimumCliVersion"`
}

// citationPattern matches the reference citations SKILL.md uses in prose
// (e.g. `references/chart-format.md`). The path shape is the contract: a
// reference an agent can load must be citable this way, so the same token
// is what validation keys on.
var citationPattern = regexp.MustCompile(`references/([A-Za-z0-9._-]+\.md)`)

var semverPattern = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

// LoadTree reads skills/, agents/, and skills/compat.yaml under root.
// Loading is shape-tolerant (missing pieces surface in Validate, not here)
// but real I/O errors fail immediately.
func LoadTree(root string) (*Tree, error) {
	tree := &Tree{Root: root}

	skillsDir := filepath.Join(root, "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("reading skills tree %s: %w", skillsDir, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		skill, err := loadSkill(filepath.Join(skillsDir, entry.Name()), entry.Name())
		if err != nil {
			return nil, err
		}
		if skill.Slug == catalogPackSkillSlug {
			pack, err := loadCatalogPack(root)
			if err != nil {
				return nil, err
			}
			skill.PackFiles = pack
		}
		tree.Skills = append(tree.Skills, *skill)
	}

	agentsDir := filepath.Join(root, "agents")
	agentEntries, err := os.ReadDir(agentsDir)
	if err != nil {
		return nil, fmt.Errorf("reading agents tree %s: %w", agentsDir, err)
	}
	for _, entry := range agentEntries {
		if !entry.IsDir() {
			continue
		}
		instructions, err := os.ReadFile(filepath.Join(agentsDir, entry.Name(), "instructions.md"))
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("reading agent instructions for %s: %w", entry.Name(), err)
		}
		tree.Agents = append(tree.Agents, Agent{Slug: entry.Name(), Instructions: instructions})
	}

	compatRaw, err := os.ReadFile(filepath.Join(skillsDir, "compat.yaml"))
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("reading compat floor: %w", err)
		}
	} else if err := yaml.Unmarshal(compatRaw, &tree.Floor); err != nil {
		return nil, fmt.Errorf("parsing skills/compat.yaml: %w", err)
	}

	return tree, nil
}

func loadSkill(dir, slug string) (*Skill, error) {
	skill := &Skill{Slug: slug, Dir: dir, ReferenceFiles: map[string][]byte{}}

	skillMD, err := os.ReadFile(filepath.Join(dir, "SKILL.md"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading SKILL.md for %s: %w", slug, err)
	}
	skill.SkillMD = skillMD
	skill.FrontmatterName, skill.Description = parseFrontmatter(skillMD)

	cited := map[string]bool{}
	for _, match := range citationPattern.FindAllStringSubmatch(string(skillMD), -1) {
		cited[match[1]] = true
	}
	for name := range cited {
		skill.CitedReferences = append(skill.CitedReferences, name)
	}
	sort.Strings(skill.CitedReferences)

	refsDir := filepath.Join(dir, "references")
	refEntries, err := os.ReadDir(refsDir)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading references for %s: %w", slug, err)
	}
	for _, ref := range refEntries {
		if ref.IsDir() {
			continue
		}
		content, err := os.ReadFile(filepath.Join(refsDir, ref.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading reference %s/%s: %w", slug, ref.Name(), err)
		}
		skill.ReferenceFiles[ref.Name()] = content
	}

	return skill, nil
}

// validatePack holds the catalog skill's assembled pack to its contract:
// present and recognizably a pack (the commons root marker plus at least
// one component reference page -- an empty or misrooted assembly would
// ship a research skill with nothing to research), every file non-empty,
// and the whole archive comfortably inside the serving engine's push
// ceilings so catalog growth fails HERE with a clear message instead of
// failing the engine push at release time.
func validatePack(skill Skill, report func(format string, args ...any)) {
	if skill.Slug != catalogPackSkillSlug {
		return
	}
	if len(skill.PackFiles) == 0 {
		report("skills/%s: catalog pack is empty -- no catalog/ tree was found to assemble", skill.Slug)
		return
	}
	if _, ok := skill.PackFiles[packDirName+"/_docs/reference-commons.md"]; !ok {
		report("skills/%s: assembled pack is missing its root marker (%s/_docs/reference-commons.md)", skill.Slug, packDirName)
	}
	referencePages := 0
	totalBytes := len(skill.SkillMD)
	for _, content := range skill.ReferenceFiles {
		totalBytes += len(content)
	}
	for path, content := range skill.PackFiles {
		if len(content) == 0 {
			report("skills/%s: pack file %s is empty", skill.Slug, path)
		}
		if strings.HasSuffix(path, "/reference.md") {
			referencePages++
		}
		totalBytes += len(content)
	}
	if referencePages == 0 {
		report("skills/%s: assembled pack carries no component reference pages", skill.Slug)
	}
	totalFiles := 1 + len(skill.ReferenceFiles) + len(skill.PackFiles)
	if totalFiles > maxSkillArchiveFiles {
		report("skills/%s: archive would carry %d files (limit %d, engine ceiling 10000) -- the pack has outgrown its margins", skill.Slug, totalFiles, maxSkillArchiveFiles)
	}
	if totalBytes > maxSkillArchiveBytes {
		report("skills/%s: archive would carry %d bytes (limit %d, engine ceiling 100MB compressed) -- the pack has outgrown its margins", skill.Slug, totalBytes, maxSkillArchiveBytes)
	}
}

// parseFrontmatter extracts name and description from SKILL.md's leading
// YAML frontmatter block. Absent or malformed frontmatter yields empty
// values, which Validate reports with a precise message.
func parseFrontmatter(skillMD []byte) (name, description string) {
	text := string(skillMD)
	if !strings.HasPrefix(text, "---\n") {
		return "", ""
	}
	rest := text[len("---\n"):]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return "", ""
	}
	var fm struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	}
	if err := yaml.Unmarshal([]byte(rest[:end]), &fm); err != nil {
		return "", ""
	}
	return fm.Name, fm.Description
}

// Validate enforces the structure contract documented in skills/README.md.
// It returns every violation rather than stopping at the first, so a PR
// author sees the complete repair list in one run.
func Validate(tree *Tree) []error {
	var errs []error
	report := func(format string, args ...any) {
		errs = append(errs, fmt.Errorf(format, args...))
	}

	if len(tree.Skills) == 0 {
		report("skills/: no skill directories found")
	}
	for _, skill := range tree.Skills {
		if len(skill.SkillMD) == 0 {
			report("skills/%s: SKILL.md is missing or empty", skill.Slug)
			continue
		}
		if skill.FrontmatterName == "" {
			report("skills/%s: SKILL.md frontmatter is missing a name", skill.Slug)
		} else if skill.FrontmatterName != skill.Slug {
			// The serving engine derives the skill's identity from the
			// frontmatter name; a mismatch would publish content under a
			// different slug than its directory advertises.
			report("skills/%s: frontmatter name %q must equal the directory slug", skill.Slug, skill.FrontmatterName)
		}
		if skill.Description == "" {
			report("skills/%s: SKILL.md frontmatter is missing a description", skill.Slug)
		}

		for _, cited := range skill.CitedReferences {
			content, exists := skill.ReferenceFiles[cited]
			if !exists {
				report("skills/%s: SKILL.md cites references/%s which does not exist", skill.Slug, cited)
			} else if len(content) == 0 {
				report("skills/%s: references/%s is empty -- an agent would load a blank document", skill.Slug, cited)
			}
		}
		citedSet := map[string]bool{}
		for _, cited := range skill.CitedReferences {
			citedSet[cited] = true
		}
		for name := range skill.ReferenceFiles {
			if !citedSet[name] {
				report("skills/%s: references/%s exists but SKILL.md never cites it -- orphaned content rots", skill.Slug, name)
			}
			if len(skill.ReferenceFiles[name]) == 0 && citedSet[name] {
				continue // already reported above
			}
			if len(skill.ReferenceFiles[name]) == 0 {
				report("skills/%s: references/%s is empty", skill.Slug, name)
			}
		}

		validatePack(skill, report)
	}

	if len(tree.Agents) == 0 {
		report("agents/: no agent directories found")
	}
	for _, agent := range tree.Agents {
		if len(agent.Instructions) == 0 {
			report("agents/%s: instructions.md is missing or empty", agent.Slug)
		}
	}

	if !semverPattern.MatchString(tree.Floor.MinimumDaemonVersion) {
		report("skills/compat.yaml: minimum_daemon_version %q is not a vX.Y.Z version", tree.Floor.MinimumDaemonVersion)
	}
	if !semverPattern.MatchString(tree.Floor.MinimumCliVersion) {
		report("skills/compat.yaml: minimum_cli_version %q is not a vX.Y.Z version", tree.Floor.MinimumCliVersion)
	}

	return errs
}
