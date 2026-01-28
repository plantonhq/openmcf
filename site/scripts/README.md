# Documentation Build System

This directory contains the build scripts that power OpenMCF's documentation site. If you're wondering how deployment component documentation magically appears on the website, you're in the right place.

## The Problem We're Solving

Imagine you're maintaining a cookbook with hundreds of recipes. The actual recipes live in your kitchen (the `apis/` directory), organized by cuisine (provider) and dish type (component). But you also want to publish a beautiful website where people can browse and search these recipes.

You could manually copy each recipe to your website folder every time you update one. But that's tedious, error-prone, and creates two versions of the truth. What if you could just write the recipe once in your kitchen, and have it automatically appear in the published cookbook whenever you go to print?

That's exactly what `copy-component-docs.ts` does.

## The Mental Model

Think of this script as a **build-time librarian** with a simple job:

1. **Scan** the `apis/` directory for deployment component documentation
2. **Transform** each README into a web-ready page with metadata
3. **Organize** components into a browsable catalog by provider
4. **Generate** index pages that list all available components
5. **Place** everything where Next.js expects to find it

All of this happens automatically before every build, ensuring your documentation site always reflects the latest component docs without manual intervention.

## Architecture Overview

### The Single Source of Truth

```
apis/org/openmcf/provider/
├── aws/
│   ├── awsalb/v1/docs/README.md           ← Source of truth
│   └── awsroute53zone/v1/docs/README.md   ← Source of truth
├── gcp/
│   └── gcpcloudrun/v1/docs/README.md      ← Source of truth
└── kubernetes/
    └── argocdkubernetes/v1/docs/README.md ← Source of truth
```

These README files live alongside the protobuf definitions for each deployment component. They're created using OpenMCF's research-driven documentation workflow and contain comprehensive deployment guides.

### The Build-Time Transformation

When you run `yarn build`, the script springs into action:

```
Build Process Flow:
─────────────────

1. prebuild: Copy component docs
   ↓
   copy-component-docs.ts runs
   ├─ Scans 10+ provider directories
   ├─ Finds ~118 component READMEs
   ├─ Generates titles (awsalb → "ALB")
   ├─ Creates frontmatter metadata
   ├─ Writes to site/public/docs/catalog/
   └─ Generates provider index pages

2. build: Next.js builds the site
   ↓
   Reads site/public/docs/catalog/
   Generates static HTML pages

3. postbuild: Pagefind indexes content
   ↓
   Enables full-text search across all docs
```

### The Output Structure

```
site/public/docs/catalog/
├── index.md                    ← Auto-generated catalog landing
├── aws/
│   ├── index.md               ← Auto-generated AWS index
│   ├── awsalb.md              ← Component doc with frontmatter
│   └── awsroute53zone.md      ← Component doc with frontmatter
├── gcp/
│   ├── index.md               ← Auto-generated GCP index
│   └── gcpcloudrun.md         ← Component doc with frontmatter
└── kubernetes/
    ├── index.md               ← Auto-generated Kubernetes index
    └── argocdkubernetes.md    ← Component doc with frontmatter
```

**Key Point**: The entire `site/public/docs/catalog/` directory is auto-generated and git-ignored. It's rebuilt fresh on every build, ensuring zero drift between source and site.

## How It Works Under the Hood

### 1. Provider Scanning

The script walks through each provider directory and looks for the characteristic structure of a deployment component:

```typescript
// Looking for this pattern:
provider/{component}/v1/docs/README.md
```

If a directory doesn't have docs at the immediate component level, it scans one level deeper to handle any potential nested structures.

**Why this matters**: You don't need to register new components anywhere. Just create a README in the right place, and the build system finds it automatically.

### 2. Title Generation

Raw component names aren't user-friendly. The script transforms them:

```
awsalb              → "ALB"
gcpgkecluster       → "GKE Cluster"
argocdkubernetes    → "ArgoCD"
clickhousekubernetes → "ClickHouse"
```

It handles:
- Removing redundant provider prefixes (`gcpgke` becomes `GKE`, not `GCP GKE`)
- Recognizing 100+ special cases (MongoDB, PostgreSQL, CloudFront)
- Uppercasing acronyms (ALB, EKS, GKE, VPC, DNS)
- Inserting spaces appropriately (certmanager → "Cert Manager")

**The result**: Human-readable titles without manual configuration.

### 3. Frontmatter Generation

Every component doc gets consistent metadata added to the top:

```yaml
---
title: "ALB"
description: "ALB deployment documentation"
icon: "package"
order: 100
componentName: "awsalb"
---
```

This frontmatter:
- Powers the documentation site's navigation
- Provides SEO metadata
- Controls sidebar display
- Enables search indexing

The original README content follows immediately after, completely unchanged.

### 4. Provider Index Pages

For each provider with components, the script generates an index page:

```markdown
---
title: "AWS"
description: "Deploy AWS resources using OpenMCF"
icon: "cloud"
order: 10
---

# AWS

The following AWS resources can be deployed using OpenMCF:

- [ALB](/docs/catalog/aws/awsalb)
- [Route53 Zone](/docs/catalog/aws/awsroute53zone)
- [EKS Cluster](/docs/catalog/aws/awsekscluster)
...
```

These index pages are **completely auto-generated**. The list stays in sync with available components automatically.

### 5. Catalog Index

The top-level catalog index provides a visual grid of provider cards with icons and component counts:

```markdown
---
title: "Catalog"
description: "Browse deployment components organized by cloud provider"
icon: "package"
order: 50
---

# Catalog

Browse deployment components by cloud provider:

<div class="grid grid-cols-1 md:grid-cols-2 gap-4 my-6">
  <a href="/docs/catalog/aws">
    <img src="/images/providers/aws.svg" />
    <div>AWS</div>
    <div>22 components</div>
  </a>
  <!-- More providers... -->
</div>
```

## Developer Workflows

### Adding a New Component

The beauty of this system is how little you need to do:

```bash
# 1. Create your component docs
mkdir -p apis/org/openmcf/provider/aws/awsnewservice/v1/docs/
cat > apis/org/openmcf/provider/aws/awsnewservice/v1/docs/README.md << 'EOF'
# AWS New Service

Comprehensive deployment guide for AWS New Service...
EOF

# 2. That's it! Build and see it live:
cd site && yarn build

# Output shows it was discovered:
# 📦 AWS: Found 23 components
#    ✓ awsnewservice
```

The documentation is now available at `/docs/catalog/aws/awsnewservice`.

### Updating Existing Documentation

Just edit the source README:

```bash
vim apis/org/openmcf/provider/gcp/gcpcloudrun/v1/docs/README.md
cd site && yarn build
```

Changes appear on the next build. No additional steps.

### Testing Locally

```bash
# Development server (hot reload)
cd site && yarn dev
# Visit http://localhost:3000/docs/catalog/aws/awsalb

# Production preview (includes search)
make preview-site
# Visit http://localhost:3000/docs/catalog/aws/awsalb
```

Both modes show your documentation with full styling and navigation.

### Manual Re-generation

Normally the build handles everything, but you can run the script directly:

```bash
cd site
yarn copy-docs

# Output:
# 🚀 Starting component documentation copy process...
# 📁 Scanning 10 providers...
# 📦 AWS: Found 22 components
# 📦 GCP: Found 5 components
# ...
# 📊 Summary:
#    Providers: 8
#    Components copied: 118
```

## Technical Deep Dive

### Why This Architecture?

**Single Source of Truth**  
Documentation lives with the code it describes. The protobuf definitions and deployment docs are versioned together in `apis/{provider}/{component}/v1/`.

**Build-Time Generation**  
Next.js requires files in `public/` at build time for static export. GitHub Pages serves static HTML with no server-side processing. This means we **must** copy docs during the build, not at runtime.

**Git-Ignored Output**  
The `site/public/docs/catalog/` directory is git-ignored because it's a build artifact. This prevents:
- Merge conflicts on generated files
- Git history bloat from duplicate content
- Drift between source and output

**Automatic Title Normalization**  
Component names follow protobuf conventions (`awsalb`, `gcpgkecluster`), which aren't human-readable. The title generation logic encodes domain knowledge about cloud provider terminology, ensuring consistent display names without manual configuration.

### Performance Characteristics

**Script Execution Time**  
Scanning 10 providers with 118 components takes ~2-3 seconds on modern hardware. The bottleneck is file I/O (reading READMEs and writing output files), not computation.

**Build Impact**  
Total build time is dominated by Next.js static generation (~30-60 seconds for 135+ pages). The documentation copy adds ~5% overhead.

**Index Size**  
The generated `catalog/` directory contains ~118 component docs plus index pages. At ~10-50KB per doc, total output is ~2-3MB of markdown files that Next.js transforms into HTML.

### Special Cases Handled

**Flat Provider Structures**  
All providers organize components in a flat structure:

```
kubernetes/
├── certmanager/v1/docs/README.md
├── kubernetesargocd/v1/docs/README.md
└── kubernetesredis/v1/docs/README.md

aws/
├── awsalb/v1/docs/README.md
└── awsekscluster/v1/docs/README.md
```

The scanner handles this flat structure efficiently, checking for `v1/docs/README.md` at the component level first.

**Provider Icon Mapping**  
The catalog index includes visual provider cards. Icons are mapped by provider name:

```typescript
const iconMap: Record<string, string> = {
  'aws': '/images/providers/aws.svg',
  'gcp': '/images/providers/gcp.svg',
  'kubernetes': '/images/providers/kubernetes.svg',
  // ... more providers
};
```

If a provider has no icon, it falls back to a generic cloud icon.

**Title Special Cases**  
The title generator maintains a dictionary of 100+ special cases for proper capitalization:

- Product names: MongoDB, PostgreSQL, MySQL, ClickHouse
- Acronyms: ALB, EKS, GKE, VPC, DNS, IAM
- Compound terms: Cert Manager, External DNS, Ingress Nginx

This ensures professional display names without manual intervention.

## Configuration

The script has minimal configuration, all embedded in the code:

### Provider List

```typescript
const providerDirs = [
  'aws', 'gcp', 'azure', 'kubernetes',
  'cloudflare', 'civo', 'digitalocean',
  'atlas', 'confluent', 'snowflake'
];
```

When you add a new provider with components, update this list so the script clears the generated directory properly on rebuild.

### Path Definitions

```typescript
const scriptDir = __dirname;
const projectRoot = path.join(scriptDir, '../..');
const apisRoot = path.join(projectRoot, 'apis/org/openmcf/provider');
const siteDocsRoot = path.join(scriptDir, '../public/docs/catalog');
```

These assume the standard project structure. If you move directories, update these paths.

## Integration Points

### package.json Scripts

```json
{
  "scripts": {
    "copy-docs": "tsx scripts/copy-component-docs.ts",
    "prebuild": "yarn copy-docs",
    "build": "next build --turbopack",
    "postbuild": "pagefind --site out ..."
  }
}
```

The `prebuild` hook ensures docs are copied before Next.js builds. You can also run `yarn copy-docs` manually during development.

### .gitignore

```gitignore
# generated docs (copied from apis/ during build)
public/docs/catalog/
```

**Important**: The `public/docs/` directory itself is NOT ignored. Manual docs like `public/docs/index.md` (the welcome page) and `public/docs/getting-started.md` are committed to git. Only the `catalog/` subdirectory is generated.

### Next.js File-Based Routing

The catch-all route `site/src/app/docs/[[...slug]]/page.tsx` handles all documentation URLs:

```
/docs/catalog/aws/awsalb
   ↓
reads: public/docs/catalog/aws/awsalb.md
   ↓
renders: Markdown with frontmatter-based metadata
```

## Maintenance

### Adding New Provider Support

When you add documentation for a new cloud provider:

1. Create the provider directory structure:
   ```bash
   mkdir -p apis/org/openmcf/provider/newprovider/newcomponent/v1/docs/
   ```

2. Add the provider to the clearing list in `copy-component-docs.ts`:
   ```typescript
   const providerDirs = [
     'aws', 'gcp', 'azure', 'kubernetes',
     'newprovider', // ← Add here
   ];
   ```

3. Optionally add a provider icon to `site/public/images/providers/newprovider.svg` and register it:
   ```typescript
   const iconMap: Record<string, string> = {
     'newprovider': '/images/providers/newprovider.svg',
   };
   ```

4. Build and verify:
   ```bash
   yarn copy-docs
   # Should show: 📦 NEWPROVIDER: Found X components
   ```

### Updating Title Generation Logic

If you encounter a component name that doesn't format well, add it to the special cases:

```typescript
const specialCases: Record<string, string> = {
  'mynewservice': 'My New Service',
  // ... existing cases
};
```

Or for compound names:

```typescript
const compoundSpecialCases: Record<string, string> = {
  'complexservicename': 'Complex Service Name',
};
```

Run `yarn copy-docs` to verify the output.

## Troubleshooting

### Documentation Not Appearing

**Symptom**: Component README exists but doesn't show up in catalog.

**Checklist**:
1. Verify path structure: `apis/org/openmcf/provider/{provider}/{component}/v1/docs/README.md`
2. Check file isn't empty
3. Run `yarn copy-docs` and look for errors
4. Verify provider is in `providerDirs` list if it's a new provider

### Incorrect Title Formatting

**Symptom**: Component title shows as "Gcpgkecluster" instead of "GKE Cluster".

**Solution**: Add to special cases or verify acronym list includes the terms.

### Manual Docs Deleted on Build

**Symptom**: Your handwritten docs in `public/docs/` disappeared after building.

**Cause**: Script only clears `public/docs/catalog/` and provider-specific subdirectories. If you created a file that matches a provider name (like `public/docs/aws.md`), it gets deleted.

**Solution**: Keep manual docs outside the `catalog/` directory and outside provider-named directories.

### Stale Documentation

**Symptom**: Updated the README but changes don't appear on the site.

**Checklist**:
1. Did you rebuild? Changes only appear after `yarn build` or `yarn copy-docs`.
2. Check browser cache—hard refresh (Cmd+Shift+R / Ctrl+Shift+R).
3. Verify you edited the source in `apis/`, not the generated file in `site/public/docs/catalog/`.

## Future Enhancements

This script could evolve to support:

**Metadata Extraction from Protos**  
Parse the component's `.proto` files to extract descriptions, examples, and field documentation, enriching the generated pages beyond just the README.

**Cross-Reference Generation**  
Detect relationships between components (e.g., VPC referenced by Subnet) and automatically generate "Related Components" sections.

**Version Support**  
Handle multiple API versions (`v1`, `v2`) and generate version-specific documentation with migration guides.

**Category Tags**  
Auto-detect component categories (compute, storage, networking) based on proto annotations or directory structure, enabling filtered browsing.

**Documentation Quality Checks**  
Validate that each README includes required sections (Overview, Prerequisites, Examples) and warn on missing content.

## Summary

The `copy-component-docs.ts` script is the bridge between your API definitions and your documentation site. It embodies the principle of **single source of truth**: write documentation once, alongside the code it describes, and let the build system make it web-ready.

For developers, this means:
- ✅ **No manual copying** between directories
- ✅ **No frontmatter boilerplate** to write
- ✅ **No index pages** to maintain
- ✅ **No git conflicts** on generated files
- ✅ **No drift** between source and site

Just write great documentation in `apis/`, and the build system handles the rest.

---

**Questions or Issues?**  
If something isn't working as expected or you have ideas for improvement, the script is thoroughly commented and designed to be modified. Start by reading through the main `copyComponentDocs()` function to understand the flow, then dive into specific helpers as needed.

