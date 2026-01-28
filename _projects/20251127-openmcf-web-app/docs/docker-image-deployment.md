# Docker Image Deployment Guide

**Last Updated:** December 23, 2025
**Purpose:** Build and deploy the unified OpenMCF Docker image to GitHub Container Registry (GHCR)

---

## Overview

This guide covers building the unified Docker image (MongoDB + Backend + Frontend) and publishing it to GitHub Container Registry (GHCR) so users can install it via the CLI with `planton webapp init`.

**GHCR Repository:** `ghcr.io/plantonhq/openmcf`
**Image Tag Strategy:** `latest` for stable releases, versioned tags for specific releases
**Access:** Publicly accessible - no authentication required for pulling images

---

## Automated Builds via GitHub Actions

The project uses GitHub Actions to automatically build and push Docker images to GHCR. Images are built with multi-architecture support (AMD64 + ARM64).

### Triggering a Build (Testing Phase)

1. **Navigate to GitHub Actions** in your repository: `https://github.com/plantonhq/openmcf/actions`
2. **Select "Build and Push Docker Image to GHCR"** workflow
3. **Click "Run workflow"** button (top right)
4. **Optionally specify a version tag** (e.g., `test-v1`, `v0.9.0-rc1`) or leave empty for `latest` only
5. **Click "Run workflow"** to start the build

The workflow will:
- Build for both `linux/amd64` and `linux/arm64` architectures
- Push to `ghcr.io/plantonhq/openmcf:latest`
- Also push to the specified version tag if provided
- Make the images publicly accessible (no credentials needed to pull)

### Future: Tag-Based Releases

Once testing is complete, the workflow will be updated to trigger automatically on git tags:
- Push a tag: `git tag v1.0.0 && git push origin v1.0.0`
- GitHub Actions automatically builds and publishes the image

## Prerequisites for Manual Builds

### 1. Docker Engine

**Verify Docker is running:**
```bash
docker --version
# Output: Docker version 24.0.x or higher

docker info
# Should show server information
```

### 3. Sufficient Disk Space

**Check available space:**
```bash
df -h

# You need at least 10GB free for building
# - Build cache: ~3GB
# - Intermediate layers: ~4GB
# - Final image: ~500MB
```

---

## Manual Building (Optional)

> **Note:** Manual builds are not required as GitHub Actions handles automated builds. This section is for reference only.

### Step 1: Navigate to Project Root

```bash
cd /Volumes/Others/Work/crafts/leftbin/planton/openmcf
```

### Step 2: Build the Image

**Build for current architecture:**

```bash
# Build the unified image
docker build -f app/Dockerfile.unified -t ghcr.io/plantonhq/openmcf:latest .

# This will take 5-10 minutes depending on your machine
# You'll see output from all three build stages:
#   [1/3] Building backend...
#   [2/3] Building frontend...
#   [3/3] Creating final image...
```

**Build output you'll see:**
```
[+] Building 450.5s (45/45) FINISHED
 => [backend-builder 1/8] FROM golang:1.24.7-alpine
 => [backend-builder 2/8] RUN apk add --no-cache git make bash curl wget
 => [backend-builder 3/8] RUN curl -fsSL https://github.com/bufbuild/buf/...
 => [backend-builder 4/8] WORKDIR /build
 => [backend-builder 5/8] COPY go.work go.work.sum ./
 => [backend-builder 6/8] COPY go.mod go.sum ./
 => [backend-builder 7/8] COPY pkg/ ./pkg/
 => [backend-builder 8/8] RUN CGO_ENABLED=0 GOOS=linux go build...
 => [frontend-builder 1/6] FROM node:20-alpine
 => [frontend-builder 2/6] RUN corepack enable
 => [frontend-builder 3/6] COPY package.json yarn.lock ./
 => [frontend-builder 4/6] RUN yarn install --frozen-lockfile
 => [frontend-builder 5/6] COPY . .
 => [frontend-builder 6/6] RUN yarn build
 => [stage-2 1/15] FROM ubuntu:22.04
 => [stage-2 2/15] RUN apt-get update && apt-get install -y...
 => [stage-2 3/15] RUN curl -fsSL https://deb.nodesource.com/setup_20.x...
 => [stage-2 4/15] COPY --from=backend-builder /usr/local/go /usr/local/go
 => [stage-2 5/15] RUN wget -q "https://get.pulumi.com/releases/sdk/pulumi...
 => [stage-2 6/15] RUN useradd -u 1001 -r -g 0 -d /home/appuser...
 => [stage-2 7/15] COPY --from=backend-builder /build/app/backend/bin/server...
 => [stage-2 8/15] COPY --from=frontend-builder /app/public /app/frontend/public
 => exporting to image
 => => naming to docker.io/satishlleftbin/openmcf:latest
```

### Step 3: Verify the Build

**Check image size:**
```bash
docker images satishlleftbin/openmcf

# Expected output:
# REPOSITORY                         TAG       IMAGE ID       CREATED          SIZE
# satishlleftbin/openmcf    latest    abc123def456   2 minutes ago    520MB
```

**Expected image size:** ~500-550MB

**Test the image locally:**
```bash
# Create a test container
docker run -d \
  --name test-planton \
  -p 3000:3000 \
  -p 50051:50051 \
  satishlleftbin/openmcf:latest

# Wait 60 seconds for services to start
sleep 60

# Check if services are running
curl http://localhost:3000
curl http://localhost:50051

# Check logs
docker logs test-planton

# Clean up
docker stop test-planton
docker rm test-planton
```

---

## Pushing to GHCR (Manual - Not Recommended)

> **Note:** GitHub Actions handles this automatically. Manual pushes are not recommended.

If you need to manually push:

### Step 1: Authenticate with GHCR

```bash
# Create a GitHub Personal Access Token with write:packages scope
# Visit: https://github.com/settings/tokens

# Login to GHCR
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

### Step 2: Tag and Push

**Tag with version number:**
```bash
# Get current version (or manually set)
VERSION="1.0.0"

# Tag with version
docker tag ghcr.io/plantonhq/openmcf:latest ghcr.io/plantonhq/openmcf:${VERSION}

# Verify tags
docker images ghcr.io/plantonhq/openmcf
```

**Push the tags:**
```bash
docker push ghcr.io/plantonhq/openmcf:latest
docker push ghcr.io/plantonhq/openmcf:${VERSION}
```

**Expected push time:** 5-10 minutes (depending on internet speed for ~500MB)

### Step 3: Verify on GitHub Container Registry

**Check the package:**
1. Visit https://github.com/orgs/openmcf/packages/container/package/openmcf
2. Or check your repository's Packages tab
3. Verify the `latest` tag is present
4. Check the size (~520MB)
5. Verify the last push timestamp

**Pull and test:**
```bash
# Pull to verify (no authentication required for public images)
docker pull ghcr.io/plantonhq/openmcf:latest

# Test the image
docker run -d -p 3000:3000 -p 50051:50051 ghcr.io/plantonhq/openmcf:latest
```

---

## Multi-Architecture Support

GitHub Actions automatically builds for multiple architectures:
- `linux/amd64` - Intel/AMD processors
- `linux/arm64` - Apple Silicon, ARM servers

**Verification:**
```bash
# Check manifest (no authentication required)
docker manifest inspect ghcr.io/plantonhq/openmcf:latest

# Should show both platforms:
# - linux/amd64
# - linux/arm64
```

**Manual multi-arch build** (if needed):
```bash
# Create a new builder
docker buildx create --name multiarch --use

# Build for both amd64 and arm64
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -f app/Dockerfile.unified \
  -t ghcr.io/plantonhq/openmcf:latest \
  --push \
  .
```

---

## Version Management Strategy

### Semantic Versioning

Use semantic versioning: `MAJOR.MINOR.PATCH`

**Example tags:**
```bash
# Latest stable (always updated)
ghcr.io/plantonhq/openmcf:latest

# Specific version
ghcr.io/plantonhq/openmcf:v1.0.0
ghcr.io/plantonhq/openmcf:v1.0.1
ghcr.io/plantonhq/openmcf:v1.1.0

# Development/preview (testing phase)
ghcr.io/plantonhq/openmcf:test-v1
ghcr.io/plantonhq/openmcf:v0.9.0-rc1
```

### Tagging Workflow (Automated)

**Testing Phase (Current):**
1. Go to GitHub Actions
2. Run "Build and Push Docker Image to GHCR" workflow
3. Specify version tag (e.g., `test-v1`)
4. Workflow builds and pushes automatically

**Production Phase (Future):**
```bash
# Create and push a tag
VERSION="v1.0.1"
git tag ${VERSION}
git push origin ${VERSION}

# GitHub Actions automatically:
# - Builds multi-arch image
# - Pushes as ghcr.io/plantonhq/openmcf:v1.0.1
# - Updates ghcr.io/plantonhq/openmcf:latest
```

### CLI Configuration

The CLI is configured to pull from GHCR:

```go
// cmd/openmcf/root/webapp/init.go
const (
    DockerImageName  = "ghcr.io/plantonhq/openmcf"
    DockerImageTag   = "latest"
)
```

---

## GitHub Actions Workflow

The project includes a GitHub Actions workflow at `.github/workflows/docker-build-push.yml`.

**Current Configuration (Testing Phase):**
- Trigger: Manual `workflow_dispatch`
- Input: Optional version tag
- Builds: Multi-architecture (amd64 + arm64)
- Pushes to: `ghcr.io/plantonhq/openmcf`
- Authentication: Uses built-in `GITHUB_TOKEN`

**Future Configuration (Production):**
Will be updated to trigger automatically on git tags matching `v*` pattern.

**Key Features:**
- No secrets configuration needed (uses `GITHUB_TOKEN`)
- Automatic multi-arch builds
- Build cache using GitHub Actions cache
- Public package visibility
- Build summary with pull commands

---

## Troubleshooting

### Build Fails

**Issue: "no space left on device"**

```bash
# Clean up Docker
docker system prune -a --volumes

# Remove unused images
docker image prune -a

# Check space
df -h
```

**Issue: "failed to solve: process ... did not complete successfully"**

```bash
# Check which stage failed
# Common issues:
# - Network issues (npm/go dependencies)
# - Missing files in COPY commands
# - Build errors in backend/frontend

# Build with verbose output
docker build -f app/Dockerfile.unified --progress=plain -t test-image .
```

### Push Fails

**Issue: "denied: requested access to the resource is denied"**

For automated builds, this should not occur as GitHub Actions uses `GITHUB_TOKEN`.

For manual pushes:
```bash
# Re-authenticate with GHCR
docker logout ghcr.io
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

**Issue: "connection timeout"**

```bash
# Check internet connection
ping ghcr.io

# For GitHub Actions, check workflow logs
# For manual push, try again:
docker push ghcr.io/plantonhq/openmcf:latest
```

### Image Too Large

**Current size: ~520MB**

If you need to reduce size:

```dockerfile
# In Dockerfile.unified, consider:
# 1. Use Alpine instead of Ubuntu (complex due to MongoDB)
# 2. Remove unnecessary packages
# 3. Combine RUN commands to reduce layers
# 4. Use .dockerignore to exclude files
```

Create `.dockerignore`:
```
node_modules
.git
.next
*.log
*.md
.DS_Store
```

---

## Testing After Push

### Test Fresh Install on New Machine

**Simulate fresh install:**
```bash
# Remove local images
docker rmi ghcr.io/plantonhq/openmcf:latest

# Test CLI installation flow (no authentication required)
planton webapp init
planton webapp start

# Verify it pulled from GHCR
docker images | grep ghcr.io

# Test the web app
open http://localhost:3000
```

### Test on Different Architectures

**Intel Mac / Linux:**
```bash
docker pull ghcr.io/plantonhq/openmcf:latest
# Should pull linux/amd64 variant automatically
```

**Apple Silicon Mac:**
```bash
docker pull ghcr.io/plantonhq/openmcf:latest
# Should pull linux/arm64 variant automatically
```

---

## Quick Reference

### Automated Build Workflow (Recommended)

```bash
# 1. Navigate to GitHub Actions
# https://github.com/plantonhq/openmcf/actions

# 2. Select "Build and Push Docker Image to GHCR" workflow

# 3. Click "Run workflow"

# 4. Optionally specify version tag (e.g., test-v1, v0.9.0-rc1)

# 5. Click "Run workflow" to start

# 6. Wait for build to complete (~15-20 minutes)

# 7. Test the image
docker pull ghcr.io/plantonhq/openmcf:latest
docker run -d --name test-planton -p 3000:3000 -p 50051:50051 ghcr.io/plantonhq/openmcf:latest
sleep 60
curl http://localhost:3000
docker stop test-planton && docker rm test-planton

# 8. Verify on GitHub Packages
# Check your repository's Packages tab
```

### Manual Local Build (Not Recommended)

```bash
# Build locally for testing
docker build -f app/Dockerfile.unified -t ghcr.io/plantonhq/openmcf:test .

# Test locally
docker run -d --name test-planton -p 3000:3000 -p 50051:50051 ghcr.io/plantonhq/openmcf:test
```

---

## Next Steps

After building and pushing images to GHCR:

1. ✅ **Test the CLI installation flow** on a clean machine
2. ✅ **Verify multi-arch support** on both Intel and ARM machines
3. ✅ **Document the version** in release notes
4. ✅ **Update changelog** with image information
5. ✅ **Monitor GitHub Packages** for download statistics
6. ✅ **Transition to tag-based releases** after testing phase

---

## Support

**GHCR Repository:** `ghcr.io/plantonhq/openmcf`
**Package URL:** https://github.com/plantonhq/openmcf/pkgs/container/openmcf
**Default Tag:** `latest`
**Access:** Public (no authentication required)

For issues:
- Check GitHub Packages page for your repository
- Review GitHub Actions workflow logs
- Verify image digest matches
- Test pull on different machines and architectures
- Check workflow run history for build logs

---

**Registry:** GitHub Container Registry (GHCR)
**Build Method:** Automated via GitHub Actions
**Multi-Arch Support:** Yes (linux/amd64, linux/arm64)
**Public Access:** Yes (no credentials required for pulling)

