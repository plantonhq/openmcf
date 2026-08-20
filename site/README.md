# planton.ai — Website, docs, and blog

This is the source for the planton.ai website, living in the `site/` folder of the Planton open-source repository. It contains:

- All written content and visuals for the public website: product pages, docs, blog, case studies, tutorials, and marketing pages
- File system–based content for docs and blog (Markdown/MDX) alongside the components, templates, and logic used to render them
- All the components, styling, and configurations used to run the site as an application

## Table of contents
- [Quick start](#quick-start)
- [Advanced setup](#advanced-setup)
- [Contributing](#contributing)
 - [Rules for docs](#rules-for-docs)
- [Media management](#media-management)

## Quick start

1. **Pre-installation**

    Install Node.js (≥ 18). If you use `nvm`, run `nvm use` in the repo to switch to the recommended version. Enable Yarn (Berry) via Corepack.

    ```bash
    node -v   # should be >= 18
    corepack enable  # enables Yarn 4 from packageManager field
    ```

2.  **Start developing**

    Clone the repo and navigate into the site folder:

    ```bash
    git clone git@github.com:plantonhq/planton.git && cd planton/site
    ```

    Install dependencies and start the dev server (Next.js):

    ```bash
    yarn
    yarn dev
    ```

    Or, from the repository root: `make run-site`. Your site is now running at `http://localhost:3000`.

3.  **Open the source code and start editing!**

    Content lives under `content/` and app code under `src/`. Changes in most files hot-reload automatically.

## Advanced setup

### Scripts

```bash
# Lint
yarn lint

# Build for production
yarn build

# Start production server (after build)
yarn start

# Audit the static export: crawl every page in a headless browser, record
# what each page actually fetches, and capture full-page screenshots.
# Run against a served out/ (e.g. `npx serve out -l 4175`).
node scripts/audit-static-export.mjs --base http://localhost:4175 --out /tmp/site-audit
```

### Content model

- Docs: Markdown files under `public/docs/**`
- Blog: MD/MDX under `public/blog/**`
- Static assets: `public/`
- Media (images, videos): See [Media management](#media-management) below

The app uses file system–based routing and content rendering. To add a new page:
1) Create a Markdown/MDX file under `public/docs` or `public/blog`
2) Include frontmatter (e.g., `title`, `description` for docs; `title`, `date`, `author`, `tags` for blog)
3) Run `yarn dev` and navigate to the route

### Project structure

- `src/app/docs/**`: Docs routes and components
- `src/app/blog/**`: Blog routes and components
- `src/app/components/**`: Shared UI components
- `src/lib/**`: Utilities for parsing/reading content

## Contributing

We <3 contributions big and small. In priority order (although everything is appreciated):

- Open a PR
  - Follow the [Quick start](#quick-start) to run the site locally
  - For basic content edits, navigate to the file in GitHub and click the edit icon (pencil)
- Open an issue for bugs or content ideas: `https://github.com/plantonhq/planton/issues/new`

Guidelines:
- Keep content changes focused and include a brief summary in the PR description
- For docs/blog posts, prefer copy-pastable examples and reference exact file paths or commands where possible
- Commit and PR conventions follow the repository-wide rules under `_rules/git/` and `_rules/pull-requests/`

## Rules for docs

We use coding agents as a native part of our workflow for writing docs and blog posts. To make this transparent and repeatable, we maintain rule files under:

```
site/_rules/docs/
  00-planton-ai-docs-style.mdc     # Style and tone guide for docs/blogs
  01-planton-ai-docs-sources.mdc   # Source-of-truth and referencing guide
  02-planton-ai-docs-elements.mdc  # Rendering capabilities (Mermaid, code languages)
  03-planton-ai-docs-media.mdc     # Media and screenshot URL requirements
```

Why we do this: we want our documentation to be accurate, humble, and useful. That means grounding every page in real code (protobufs, IaC, backend/CLI), and writing in a reader-first tone. The rules above encode that process so contributors can open PRs with confidence.

## Media management

We store screenshots and image assets in Cloudflare R2, mirrored from the `content/` folder. This keeps the repo lean, provides CDN delivery, and ensures predictable URLs.

### Architecture

```
content/
├── assets/
│   ├── _inbox/           # Drop zone for raw images (gitignored)
│   ├── _rules/           # Cursor rules for asset management
│   └── images/           # Organized, processed images
└── ...

↓ syncs to ↓

assets.planton.ai/site/
├── assets/
│   └── images/           # Same structure as local
└── ...
```

**Mental model**: Everything in `content/assets/` is available at `https://assets.planton.ai/site/`.

### Adding images (recommended workflow)

1. **Drop images** in `content/assets/_inbox/`

2. **Invoke the cursor rule** with context:
   ```
   @process-planton-ai-images
   
   I've added screenshots of the deployment pipeline feature for Service Hub docs.
   ```

3. **The tool processes images**:
   - Compresses (jpegoptim for JPEG, pngquant for PNG)
   - Renames following convention: `YYYY-MM-DD-HHMMSS-context-description.ext`
   - Moves to appropriate folder
   - Uploads to R2

4. **Use the R2 URL** in docs/blog:
   ```md
   ![Pipeline overview](https://assets.planton.ai/site/images/service-hub/2025-12-31-103045-service-hub-pipeline-overview.png)
   ```

### One-time setup

```bash
# Install CLI tools
brew install jpegoptim pngquant awscli

# Create virtual environment and install Python dependencies
python3 -m venv .venv
source .venv/bin/activate
pip install -r tools/image_processor/requirements.txt

# Configure R2 access (get credentials from Cloudflare dashboard)
aws configure --profile r2
```

### Manual commands

```bash
# Activate virtual environment first
source .venv/bin/activate

# Check prerequisites
make check-images

# List inbox contents
make inbox

# Sync assets to R2
make sync-assets

# Process images with specific context
python -m tools.image_processor process -c "service-hub" -d "feature-name"
```

### Naming convention

```
YYYY-MM-DD-HHMMSS-context-description.ext
```

Examples:
- `2025-12-31-103045-service-hub-deployment-pipeline.png`
- `2025-12-31-143500-kubernetes-dashboard-pod-logs.png`

### Folder structure

| Local Path | R2 URL |
|------------|--------|
| `content/assets/images/service-hub/x.png` | `https://assets.planton.ai/site/images/service-hub/x.png` |

### Troubleshooting

- **Preflight check fails**: Run `make check-images` to see what's missing
- **R2 access denied**: Verify `aws configure list --profile r2` shows credentials
- **No images found**: Ensure images are in `content/assets/_inbox/` (not a subdirectory)

For detailed documentation, see `content/assets/_rules/README.md`.
