# Cloudflare D1 Database Examples

This guide provides concrete, copy-and-paste examples for common Cloudflare D1 database deployment scenarios using OpenMCF.

## Table of Contents

- [Minimal Configuration](#minimal-configuration)
- [Development Database](#development-database)
- [Preview/Staging Database](#previewstaging-database)
- [Production Database with Region Specification](#production-database-with-region-specification)
- [Global Database with Read Replication](#global-database-with-read-replication)
- [Multi-Environment Setup](#multi-environment-setup)
- [Pulumi Go Example](#pulumi-go-example)
- [Terraform HCL Example](#terraform-hcl-example)
- [Complete CI/CD Workflow](#complete-cicd-workflow)

---

## Minimal Configuration

The simplest possible D1 database with only required fields.

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareD1Database
metadata:
  name: minimal-db
spec:
  accountId: "abc123def456..."
  databaseName: "my-minimal-db"
```

**Deploy:**

```bash
planton apply -f minimal-db.yaml
```

**Use Case:** Quick experimentation or proof-of-concept. Not recommended for production.

---

## Development Database

Database optimized for a developer's local workflow.

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareD1Database
metadata:
  name: my-app-dev-db
  labels:
    environment: development
    app: my-application
spec:
  accountId: "abc123def456..."
  databaseName: "my-app-dev-db"
  region: "weur"  # Western Europe - close to EU developers
```

**Deploy:**

```bash
planton apply -f dev-db.yaml
```

**Schema Setup:**

After provisioning, apply migrations:

```bash
# Create migration
npx wrangler d1 migrations create my-app-dev-db create_users_table

# Edit migrations/0001_create_users_table.sql
# CREATE TABLE users (
#   id INTEGER PRIMARY KEY AUTOINCREMENT,
#   email TEXT NOT NULL UNIQUE,
#   name TEXT,
#   created_at DATETIME DEFAULT CURRENT_TIMESTAMP
# );

# Apply migration
npx wrangler d1 migrations apply my-app-dev-db --remote
```

**Use Case:** Developer sandbox for local testing and experimentation.

---

## Preview/Staging Database

Separate database for preview deployments and staging environments.

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareD1Database
metadata:
  name: my-app-preview-db
  labels:
    environment: preview
    app: my-application
spec:
  accountId: "abc123def456..."
  databaseName: "my-app-preview-db"
  region: "wnam"  # Western North America
```

**Worker Binding (wrangler.toml):**

```toml
name = "my-app-worker"

# Production binding
[[d1_databases]]
binding = "DB"
database_name = "my-app-production-db"
database_id = "aaaa-bbbb-cccc-dddd"

# Preview binding (points to preview database)
preview_database_id = "eeee-ffff-gggg-hhhh"  # ID of my-app-preview-db
```

**Deploy:**

```bash
# Provision preview database
planton apply -f preview-db.yaml

# Apply schema migrations
npx wrangler d1 migrations apply my-app-preview-db --remote

# Deploy Worker in preview mode
npx wrangler deploy --preview
```

**Use Case:** Testing changes in a staging environment before promoting to production.

---

## Production Database with Region Specification

Production database with explicit region selection for optimal latency.

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareD1Database
metadata:
  name: my-app-prod-db
  labels:
    environment: production
    app: my-application
    tier: critical
spec:
  accountId: "abc123def456..."
  databaseName: "my-app-production-db"
  region: "enam"  # Eastern North America - close to primary user base
```

**Deploy:**

```bash
planton apply -f prod-db.yaml
```

**Use Case:** Production database for a US-based application with users primarily in the Eastern United States.

---

## Global Database with Read Replication

Database with read replication enabled for global applications with users across multiple continents.

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareD1Database
metadata:
  name: my-app-global-db
  labels:
    environment: production
    app: my-application
    tier: critical
    replication: enabled
spec:
  accountId: "abc123def456..."
  databaseName: "my-app-global-db"
  region: "enam"  # Primary region: Eastern North America
  readReplication:
    mode: "auto"  # Enable automatic read replication
```

**Worker Code (Sessions API Required):**

```javascript
export default {
  async fetch(request, env, ctx) {
    // Create a D1 session to maintain consistency across replicas
    const session = env.DB.withSession();

    // Write operation - goes to primary
    await session.prepare(
      "INSERT INTO users (email, name) VALUES (?, ?)"
    ).bind("user@example.com", "Alice").run();

    // Read operation - may use a replica, but guaranteed to see previous write
    const user = await session.prepare(
      "SELECT * FROM users WHERE email = ?"
    ).bind("user@example.com").first();

    return new Response(JSON.stringify(user), {
      headers: { "Content-Type": "application/json" },
    });
  },
};
```

**Deploy:**

```bash
# Provision database with replication
planton apply -f global-db.yaml

# Apply schema migrations
npx wrangler d1 migrations apply my-app-global-db --remote

# Deploy Worker with Sessions API code
npx wrangler deploy
```

**⚠️ Critical:** Failing to use the D1 Sessions API (`env.DB.withSession()`) with read replication enabled **will cause data consistency errors**.

**Use Case:** Global SaaS application with users in North America, Europe, and Asia Pacific.

---

## Multi-Environment Setup

Complete multi-environment setup with development, preview, and production databases.

### Directory Structure

```
my-project/
├── databases/
│   ├── dev-db.yaml
│   ├── preview-db.yaml
│   └── prod-db.yaml
├── migrations/
│   ├── 0001_create_users_table.sql
│   └── 0002_create_posts_table.sql
└── wrangler.toml
```

### Development Database

**File:** `databases/dev-db.yaml`

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareD1Database
metadata:
  name: my-app-dev-db
  labels:
    environment: development
spec:
  accountId: "abc123def456..."
  databaseName: "my-app-dev-db"
  region: "weur"
```

### Preview Database

**File:** `databases/preview-db.yaml`

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareD1Database
metadata:
  name: my-app-preview-db
  labels:
    environment: preview
spec:
  accountId: "abc123def456..."
  databaseName: "my-app-preview-db"
  region: "wnam"
```

### Production Database

**File:** `databases/prod-db.yaml`

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareD1Database
metadata:
  name: my-app-prod-db
  labels:
    environment: production
spec:
  accountId: "abc123def456..."
  databaseName: "my-app-production-db"
  region: "enam"
  readReplication:
    mode: "auto"
```

### Deployment Script

**File:** `deploy.sh`

```bash
#!/bin/bash
set -e

ENVIRONMENT=$1

if [ "$ENVIRONMENT" = "dev" ]; then
  planton apply -f databases/dev-db.yaml
  npx wrangler d1 migrations apply my-app-dev-db --remote
  npx wrangler deploy
elif [ "$ENVIRONMENT" = "preview" ]; then
  planton apply -f databases/preview-db.yaml
  npx wrangler d1 migrations apply my-app-preview-db --remote
  npx wrangler deploy --preview
elif [ "$ENVIRONMENT" = "prod" ]; then
  planton apply -f databases/prod-db.yaml
  npx wrangler d1 migrations apply my-app-production-db --remote
  npx wrangler deploy
else
  echo "Usage: ./deploy.sh [dev|preview|prod]"
  exit 1
fi
```

**Usage:**

```bash
chmod +x deploy.sh
./deploy.sh dev      # Deploy to development
./deploy.sh preview  # Deploy to preview
./deploy.sh prod     # Deploy to production
```

---

## Pulumi Go Example

Direct Pulumi Go code for provisioning a D1 database (without OpenMCF CLI).

```go
package main

import (
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create D1 database
		db, err := cloudflare.NewD1Database(ctx, "my-app-db", &cloudflare.D1DatabaseArgs{
			AccountId:           pulumi.String("abc123def456..."),
			Name:                pulumi.String("my-app-production-db"),
			PrimaryLocationHint: pulumi.String("enam"),
			ReadReplication: &cloudflare.D1DatabaseReadReplicationArgs{
				Mode: pulumi.String("auto"),
			},
		})
		if err != nil {
			return err
		}

		// Export outputs
		ctx.Export("databaseId", db.ID())
		ctx.Export("databaseName", db.Name)

		return nil
	})
}
```

**Deploy:**

```bash
pulumi up
```

---

## Terraform HCL Example

Direct Terraform HCL code for provisioning a D1 database (without OpenMCF CLI).

```hcl
terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
  }
}

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

variable "cloudflare_api_token" {
  description = "Cloudflare API token"
  type        = string
  sensitive   = true
}

variable "account_id" {
  description = "Cloudflare account ID"
  type        = string
}

resource "cloudflare_d1_database" "main" {
  account_id            = var.account_id
  name                  = "my-app-production-db"
  primary_location_hint = "enam"

  read_replication {
    mode = "auto"
  }
}

output "database_id" {
  description = "The ID of the created D1 database"
  value       = cloudflare_d1_database.main.id
}

output "database_name" {
  description = "The name of the database"
  value       = cloudflare_d1_database.main.name
}
```

**Deploy:**

```bash
terraform init
terraform apply \
  -var="cloudflare_api_token=$CLOUDFLARE_API_TOKEN" \
  -var="account_id=abc123def456..."
```

---

## Complete CI/CD Workflow

GitHub Actions workflow for automated deployment of D1 databases across environments.

**File:** `.github/workflows/deploy-d1.yml`

```yaml
name: Deploy Cloudflare D1 Database

on:
  push:
    branches:
      - main
      - develop
  pull_request:
    branches:
      - main

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install Wrangler
        run: npm install -g wrangler

      - name: Install OpenMCF CLI
        run: |
          curl -fsSL https://get.openmcf.org | bash
          planton version

      - name: Determine environment
        id: env
        run: |
          if [[ "${{ github.ref }}" == "refs/heads/main" ]]; then
            echo "environment=prod" >> $GITHUB_OUTPUT
            echo "db_name=my-app-production-db" >> $GITHUB_OUTPUT
            echo "manifest=databases/prod-db.yaml" >> $GITHUB_OUTPUT
          elif [[ "${{ github.ref }}" == "refs/heads/develop" ]]; then
            echo "environment=preview" >> $GITHUB_OUTPUT
            echo "db_name=my-app-preview-db" >> $GITHUB_OUTPUT
            echo "manifest=databases/preview-db.yaml" >> $GITHUB_OUTPUT
          else
            echo "environment=dev" >> $GITHUB_OUTPUT
            echo "db_name=my-app-dev-db" >> $GITHUB_OUTPUT
            echo "manifest=databases/dev-db.yaml" >> $GITHUB_OUTPUT
          fi

      - name: Provision D1 Database
        env:
          CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
        run: |
          planton apply -f ${{ steps.env.outputs.manifest }}

      - name: Apply Database Migrations
        env:
          CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
        run: |
          npx wrangler d1 migrations apply ${{ steps.env.outputs.db_name }} --remote

      - name: Deploy Worker
        env:
          CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
        run: |
          if [[ "${{ steps.env.outputs.environment }}" == "preview" ]]; then
            npx wrangler deploy --preview
          else
            npx wrangler deploy
          fi
```

**Required GitHub Secrets:**

- `CLOUDFLARE_API_TOKEN`: Cloudflare API token with D1 and Workers permissions

---

## Region Selection Guide

Choose the region closest to your primary user base:

| Region Code | Location | Best For |
|-------------|----------|----------|
| `weur` | Western Europe | EU users (Ireland, UK, France, Germany) |
| `eeur` | Eastern Europe | Eastern EU users (Poland, Romania) |
| `apac` | Asia Pacific | Asian users (Singapore, Japan, Australia) |
| `oc` | Oceania | Australia/New Zealand users |
| `wnam` | Western North America | West Coast US, Western Canada |
| `enam` | Eastern North America | East Coast US, Eastern Canada |

**Example:** If your users are primarily in California, choose `wnam`. If they're in New York, choose `enam`.

---

## Validation

After deployment, verify your database:

```bash
# List databases in your account
npx wrangler d1 list

# Get database info
npx wrangler d1 info my-app-production-db

# Query the database
npx wrangler d1 execute my-app-production-db --remote --command="SELECT name FROM sqlite_master WHERE type='table'"
```

---

## Common Patterns

### Per-Tenant Databases

For multi-tenant SaaS applications requiring data isolation:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareD1Database
metadata:
  name: tenant-acme-corp-db
  labels:
    tenant: acme-corp
    app: my-saas
spec:
  accountId: "abc123def456..."
  databaseName: "tenant-acme-corp-db"
  region: "wnam"
```

**Note:** Create one database per tenant. Configure Worker bindings dynamically based on tenant context.

### Feature Branch Databases

Temporary databases for feature branch testing:

```yaml
apiVersion: cloudflare.openmcf.org/v1
kind: CloudflareD1Database
metadata:
  name: feature-new-auth-db
  labels:
    branch: feature/new-auth
    ephemeral: "true"
spec:
  accountId: "abc123def456..."
  databaseName: "feature-new-auth-db"
  region: "weur"
```

**Cleanup:** Delete after feature branch is merged:

```bash
planton delete -f feature-branch-db.yaml
npx wrangler d1 delete feature-new-auth-db
```

---

## Next Steps

- **Explore Architecture**: Read [docs/README.md](docs/README.md) for in-depth architectural guidance
- **User Guide**: See [README.md](README.md) for general usage instructions
- **Pulumi Docs**: Check [iac/pulumi/README.md](iac/pulumi/README.md) for Pulumi-specific details
- **Terraform Docs**: Check [iac/tf/README.md](iac/tf/README.md) for Terraform-specific details

---

## Support

For questions or issues:
- **OpenMCF**: [openmcf.org](https://openmcf.org)
- **Cloudflare D1 Docs**: [developers.cloudflare.com/d1](https://developers.cloudflare.com/d1)
- **Wrangler CLI**: [developers.cloudflare.com/workers/wrangler](https://developers.cloudflare.com/workers/wrangler)

