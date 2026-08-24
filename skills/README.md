# Planton Skills

This tree is the source of truth for the skill definitions that power
Planton's AI assistant surfaces. A skill packages the judgment, workflow,
and working craft an agent loads at runtime — what to do, in what order,
and how to decide — as a `SKILL.md` instruction file plus on-demand
reference documents.

Agent identity definitions (the instructions an agent carries independent
of any one skill) live in the sibling `agents/` tree, one directory per
agent slug.

## Structure contract

```
skills/
  <skill-slug>/
    SKILL.md            # frontmatter: name (must equal the directory slug),
                        # description; then the instruction body
    references/         # on-demand reference documents cited by SKILL.md
      <topic>.md
  compat.yaml           # compatibility floor carried by every release manifest
  README.md
agents/
  <agent-slug>/
    instructions.md     # the agent's system instructions
```

Rules the lint gate enforces on every pull request:

- `SKILL.md` frontmatter `name:` equals the directory slug exactly.
- Reference integrity is bidirectional: every `references/<file>.md` that
  `SKILL.md` cites must exist and be non-empty, and every file present in
  `references/` must be cited by `SKILL.md`. Orphaned or empty reference
  files fail the gate.

## How releases carry this content

Every semver release of this repository packages each skill directory into
a deterministic zip and publishes them to the downloads CDN together with a
`definitions-manifest.json` carrying the release version, per-archive
SHA-256 checksums, and the compatibility floor from `compat.yaml`
(the minimum daemon and CLI versions the content assumes). Consumers verify
checksums before installing; the Planton desktop's local daemon carries no
embedded copy — it resolves a stable release-channel pointer at boot and
seeds the pointed-at release into its engine, verifying each artifact's
checksum, the compatibility floor, and content shape before adopting it.

## The authoring bar

Every skill in this tree is crafted individually, grounded in thorough
research of the platform's actual behavior — real command output, real
schemas, real failure modes. There is no bulk-import path: a skill earns
its place here by being verifiably correct, and behavioral claims about an
agent's judgment are proven against a live engine before they ship.

Skills carry judgment and workflow. Component facts (what a deployment
component is, its fields, its examples) live in the generated reference
pages beside each component's protos (`catalog/...`), and the verified
data layer (cost fact-sheets with generated estimates, control posture,
provisioning-permission manifests) lives in the catalog's committed
sidecars and central `_pricing/`/`_compliance/` trees — both are composed
into agents by reference (the multi-cloud-catalog skill's reference pack
teaches reading them as files), never duplicated into a skill's prose,
where they would rot.
