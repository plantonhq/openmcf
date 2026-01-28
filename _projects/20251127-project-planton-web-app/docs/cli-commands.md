# OpenMCF CLI Command Reference

**Last Updated:** December 12, 2025

This comprehensive guide covers all OpenMCF CLI commands including web app management, configuration, deployment components, cloud resources, credentials, and stack-updates.

---

## Table of Contents

- [Web App Management](#web-app-management) (Installation & Setup)
- [Configuration Management](#configuration-management)
- [Deployment Components](#deployment-components)
- [Cloud Resources](#cloud-resources)
- [Credentials](#credentials)
- [Stack Jobs](#stack-update)
- [Common Workflows](#common-workflows)
- [Troubleshooting](#troubleshooting)
- [Advanced Usage](#advanced-usage)

---

## Web App Management

The OpenMCF web app provides a unified web interface for managing cloud resources and deployments in a single Docker container. This section covers installation and lifecycle management.

### Prerequisites

- **Docker Engine** must be installed and running
- **OpenMCF CLI** installed via Homebrew

### Installation Quick Start

```bash
# Install CLI
brew install plantonhq/tap/openmcf

# Initialize web app (pulls Docker image, creates volumes, configures CLI)
openmcf webapp init

# Start the web app
openmcf webapp start

# Access the web interface
open http://localhost:3000
```

### Commands

#### `openmcf webapp init`

Initialize and set up the OpenMCF web app.

**Description:**
Downloads the unified Docker image, creates volumes for data persistence, sets up the container, and automatically configures the CLI to use the local backend.

**Usage:**
```bash
openmcf webapp init
```

**What it does:**
1. Checks Docker availability and validates Docker Engine is running
2. Verifies no existing installation to prevent conflicts
3. Pulls Docker image (`ghcr.io/plantonhq/openmcf:latest`)
4. Creates data volumes (MongoDB, Pulumi state, Go cache)
5. Creates container with port mappings (3000, 50051)
6. Automatically configures CLI backend URL to `http://localhost:50051`

**Output:**
```
========================================
🚀 OpenMCF Web App Initialization
========================================

📋 Step 1/5: Checking Docker availability...
✅ Docker is available and running

📋 Step 2/5: Checking for existing installation...
✅ No existing installation found

📋 Step 3/5: Pulling Docker image...
   Pulling ghcr.io/plantonhq/openmcf:latest...
✅ Docker image pulled successfully

📋 Step 4/5: Creating Docker volumes and container...
   ✓ Created MongoDB data volume
   ✓ Created Pulumi state volume
   ✓ Created Go cache volume
   ✓ Created container
✅ Container created successfully

📋 Step 5/5: Configuring CLI...
✅ CLI configured to use local backend

========================================
✨ Initialization Complete!
========================================

Next steps:
  1. Start the web app:     planton webapp start
  2. Check status:          planton webapp status
  3. View logs:             planton webapp logs

Once started, access the web interface at:
  Frontend:  http://localhost:3000
  Backend:   http://localhost:50051
```

**Errors:**

| Error | Cause | Solution |
|-------|-------|----------|
| Docker not found | Docker not installed | Install Docker (see installation guide) |
| Docker not running | Docker daemon stopped | Start Docker: `docker info` |
| Container already exists | Previous installation | Run `planton webapp uninstall` first |

**Time:** 2-5 minutes (depending on internet speed)

---

#### `openmcf webapp start`

Start the OpenMCF web app.

**Description:**
Starts the container and waits for all services (MongoDB, backend, frontend) to be healthy.

**Usage:**
```bash
openmcf webapp start
```

**What it does:**
1. Checks if container exists
2. Checks if already running
3. Starts the container
4. Waits for services to be healthy (up to 60 seconds)
5. Displays access URLs

**Output:**
```
========================================
🚀 Starting OpenMCF Web App
========================================

🔄 Starting container...
⏳ Waiting for services to start (this may take 30-60 seconds)...
✅ All services are healthy

========================================
✨ Web App Started Successfully!
========================================

Access the web interface at:
  🌐 Frontend:  http://localhost:3000
  🔌 Backend:   http://localhost:50051

Useful commands:
  planton webapp status    # Check service status
  planton webapp logs      # View service logs
  planton webapp stop      # Stop the web app
```

**Errors:**

| Error | Cause | Solution |
|-------|-------|----------|
| Container not found | Not initialized | Run `planton webapp init` |
| Already running | Container is running | No action needed |
| Timeout waiting for health | Services slow to start | Check logs: `planton webapp logs` |

**Time:** 30-60 seconds

---

#### `openmcf webapp stop`

Stop the OpenMCF web app.

**Description:**
Gracefully stops the container. All data is preserved in volumes.

**Usage:**
```bash
openmcf webapp stop
```

**What it does:**
1. Checks if container exists and is running
2. Stops the container gracefully
3. Confirms data is preserved

**Output:**
```
========================================
🛑 Stopping OpenMCF Web App
========================================

🔄 Stopping container...

✅ Web app stopped successfully

Data is preserved. To start again, run:
  planton webapp start
```

**Time:** 5-10 seconds

---

#### `openmcf webapp status`

Check the status of the OpenMCF web app.

**Description:**
Displays container status, service health, and access URLs.

**Usage:**
```bash
openmcf webapp status
```

**Output (Running):**
```
========================================
📊 OpenMCF Web App Status
========================================

Container Information:
  Name:       openmcf-webapp
  Status:     🟢 running
  Image:      ghcr.io/plantonhq/openmcf:latest

Service Status:
  MongoDB:     🟢 running (port 27017)
  Backend:     🟢 running (port 50051)
  Frontend:    🟢 running (port 3000)

Access URLs:
  🌐 Frontend:  http://localhost:3000
  🔌 Backend:   http://localhost:50051

Data Volumes:
  MongoDB:     openmcf-mongodb-data
  Pulumi:      openmcf-pulumi-state
  Go Cache:    openmcf-go-cache
```

**Output (Stopped):**
```
========================================
📊 OpenMCF Web App Status
========================================

Container Information:
  Name:       openmcf-webapp
  Status:     🔴 stopped
  Image:      ghcr.io/plantonhq/openmcf:latest

The web app is not running.

To start the web app, run:
  planton webapp start
```

---

#### `openmcf webapp logs`

View logs from the OpenMCF web app.

**Description:**
Displays logs from all services in the container.

**Usage:**
```bash
# View last 100 lines
openmcf webapp logs

# Follow logs in real-time
openmcf webapp logs -f

# Show last 500 lines
openmcf webapp logs -n 500

# Follow with custom tail
openmcf webapp logs -f -n 200
```

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--follow` | `-f` | Stream logs in real-time | `false` |
| `--tail` | `-n` | Number of lines to show | `100` |
| `--service` | | Filter by service (future) | All services |

**Output:**
```
[mongodb] 2025-12-11T10:30:00.000Z I NETWORK  [initandlisten] waiting for connections on port 27017
[backend] 2025-12-11T10:30:05.000Z INFO Successfully connected to MongoDB
[backend] 2025-12-11T10:30:05.000Z INFO Starting backend server on port 50051
[frontend] 2025-12-11T10:30:08.000Z INFO Server started on http://0.0.0.0:3000
```

**Controls:**
- Press `Ctrl+C` to stop following logs

---

#### `openmcf webapp restart`

Restart the OpenMCF web app.

**Description:**
Restarts the container and all services. Useful after configuration changes or when services are unresponsive.

**Usage:**
```bash
openmcf webapp restart
```

**What it does:**
1. Restarts the Docker container
2. Waits for services to be healthy
3. Displays access URLs

**Output:**
```
========================================
🔄 Restarting OpenMCF Web App
========================================

🔄 Restarting container...
⏳ Waiting for services to start (this may take 30-60 seconds)...
✅ All services are healthy

========================================
✨ Web App Restarted Successfully!
========================================

Access the web interface at:
  🌐 Frontend:  http://localhost:3000
  🔌 Backend:   http://localhost:50051
```

**Time:** 30-60 seconds

---

#### `openmcf webapp uninstall`

Uninstall the OpenMCF web app.

**Description:**
Removes the container. Data volumes are preserved by default unless `--purge-data` is specified.

**Usage:**
```bash
# Remove container, keep data
openmcf webapp uninstall

# Remove everything (including data)
openmcf webapp uninstall --purge-data

# Skip confirmation prompt
openmcf webapp uninstall -f
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--purge-data` | Delete all data volumes | `false` |
| `--force` / `-f` | Skip confirmation | `false` |

**Output (Default):**
```
========================================
🗑️  Uninstalling OpenMCF Web App
========================================

This will:
  - Stop the web app container
  - Remove the container
  - Keep data volumes (MongoDB, Pulumi state, Go cache)

Are you sure you want to continue? (yes/no): yes

🔄 Stopping container...
✅ Container stopped
🔄 Removing container...
✅ Container removed
ℹ️  Data volumes preserved:
   - openmcf-mongodb-data
   - openmcf-pulumi-state
   - openmcf-go-cache

   To remove them manually, run:
     docker volume rm openmcf-mongodb-data
     docker volume rm openmcf-pulumi-state
     docker volume rm openmcf-go-cache
🔄 Cleaning up CLI configuration...
✅ CLI configuration cleaned up

========================================
✨ Uninstall Complete!
========================================

To reinstall with existing data:
  planton webapp init
```

**Output (With --purge-data):**
```
========================================
🗑️  Uninstalling OpenMCF Web App
========================================

This will:
  - Stop the web app container
  - Remove the container
  - ⚠️  DELETE ALL DATA (MongoDB, Pulumi state, Go cache)

Are you sure you want to continue? (yes/no): yes

🔄 Stopping container...
✅ Container stopped
🔄 Removing container...
✅ Container removed
🔄 Removing data volumes...
   ✓ Removed openmcf-mongodb-data
   ✓ Removed openmcf-pulumi-state
   ✓ Removed openmcf-go-cache
✅ Data volumes removed
🔄 Cleaning up CLI configuration...
✅ CLI configuration cleaned up

========================================
✨ Uninstall Complete!
========================================
```

**Time:** 10-20 seconds

---

### Web App Workflows

#### First Time Setup

```bash
# 1. Initialize
openmcf webapp init

# 2. Start
openmcf webapp start

# 3. Access http://localhost:3000
```

#### Daily Usage

```bash
# Morning: Start the web app
openmcf webapp start

# Check if running
openmcf webapp status

# Evening: Stop to save resources
openmcf webapp stop
```

#### Troubleshooting

```bash
# Check status
openmcf webapp status

# View recent logs
openmcf webapp logs -n 500

# Follow logs in real-time
openmcf webapp logs -f

# Restart if unresponsive
openmcf webapp restart
```

#### Clean Reinstall

```bash
# Remove everything
openmcf webapp uninstall --purge-data -f

# Initialize fresh
openmcf webapp init

# Start
openmcf webapp start
```

### Web App Tips

- **Keep it running:** The web app starts quickly after initial setup
- **Check logs first:** Most issues can be diagnosed from logs
- **Data is safe:** Stopping the container doesn't delete data
- **Resource usage:** ~500MB RAM, ~2GB disk when running

---

## Configuration Management

The OpenMCF CLI uses a configuration system similar to Git, allowing you to set and manage settings that persist across commands.

### Commands

#### `openmcf config set <key> <value>`

Set a configuration value.

**Available Keys:**

- `backend-url` - URL of the OpenMCF backend service

**Example:**

```bash
openmcf config set backend-url http://localhost:50051
openmcf config set backend-url https://api.openmcf.com
```

**Validation:**

- `backend-url` must start with `http://` or `https://`

#### `openmcf config get <key>`

Get a configuration value.

**Example:**

```bash
openmcf config get backend-url
# Output: http://localhost:50051
```

**Error Handling:**

- Returns exit code 1 if the key is not set
- Prints error message for unknown keys

#### `openmcf config list`

List all configuration values.

**Example:**

```bash
openmcf config list
# Output: backend-url=http://localhost:50051

# If no configuration is set:
# Output: No configuration values set
```

### Configuration Storage

- Configuration is stored in `~/.openmcf/config.yaml`
- The configuration directory is created automatically with permissions `0755`
- The configuration file has permissions `0600` (user read/write only)

### Commands

#### `openmcf config set <key> <value>`

Set a configuration value.

**Available Keys:**

- `backend-url` - URL of the OpenMCF backend service

**Example:**

```bash
openmcf config set backend-url http://localhost:50051
openmcf config set backend-url https://api.openmcf.com
```

**Validation:**

- `backend-url` must start with `http://` or `https://`

#### `openmcf config get <key>`

Get a configuration value.

**Example:**

```bash
openmcf config get backend-url
# Output: http://localhost:50051
```

**Error Handling:**

- Returns exit code 1 if the key is not set
- Prints error message for unknown keys

#### `openmcf config list`

List all configuration values.

**Example:**

```bash
openmcf config list
# Output: backend-url=http://localhost:50051

# If no configuration is set:
# Output: No configuration values set
```

### Configuration Storage

- Configuration is stored in `~/.openmcf/config.yaml`
- The configuration directory is created automatically with permissions `0755`
- The configuration file has permissions `0600` (user read/write only)

---

## Deployment Components

### `openmcf list-deployment-components`

List deployment components from the backend service with optional filtering.

#### Basic Usage

```bash
# List all deployment components
openmcf list-deployment-components
```

**Sample Output:**

```
NAME                KIND                PROVIDER    VERSION  ID PREFIX  SERVICE KIND  CREATED
PostgresKubernetes  PostgresKubernetes  kubernetes  v1       k8spg      Yes           2025-11-25
AwsRdsInstance      AwsRdsInstance      aws         v1       rdsins     Yes           2025-11-25
GcpCloudSql         GcpCloudSql         gcp         v1       gcpsql     Yes           2025-11-25

Total: 3 deployment component(s)
```

#### Filtering by Kind

Use the `--kind` flag to filter deployment components by their kind.

```bash
# Filter by specific kind
openmcf list-deployment-components --kind PostgresKubernetes
openmcf list-deployment-components -k AwsRdsInstance
```

**Sample Output:**

```
NAME                KIND                PROVIDER    VERSION  ID PREFIX  SERVICE KIND  CREATED
PostgresKubernetes  PostgresKubernetes  kubernetes  v1       k8spg      Yes           2025-11-25

Total: 1 deployment component(s) (filtered by kind: PostgresKubernetes)
```

#### Flags

- `--kind, -k` - Filter deployment components by kind (optional)
- `--help, -h` - Show help information

#### Output Format

The command displays results in a table with the following columns:

- **NAME** - Display name of the deployment component
- **KIND** - Type/kind of the deployment component
- **PROVIDER** - Cloud provider (aws, gcp, kubernetes, etc.)
- **VERSION** - API version of the component
- **ID PREFIX** - Prefix used for resource IDs
- **SERVICE KIND** - Whether this component can launch services (Yes/No)
- **CREATED** - Creation date in YYYY-MM-DD format

### Prerequisites

**Option 1: Using Local Web App (Recommended)**

If you're using the OpenMCF web app, the backend URL is automatically configured:

```bash
openmcf webapp init
openmcf webapp start
```

**Option 2: Using Remote/External Backend**

If connecting to a remote backend, configure the URL manually:

```bash
openmcf config set backend-url <your-backend-url>
```

---

## Cloud Resources

Cloud resources represent infrastructure resources that can be created and managed through the OpenMCF backend service. You can perform complete lifecycle management including create, read, update, delete, and list operations.

### `openmcf cloud-resource:create`

Create a new cloud resource from a YAML manifest file. This command automatically triggers deployment using credentials resolved from the database.

#### Basic Usage

```bash
openmcf cloud-resource:create --arg=path/to/manifest.yaml
```

**Example:**

```bash
# Create a cloud resource from a YAML file
openmcf cloud-resource:create --arg=my-vpc.yaml
```

**Sample Output:**

```
✅ Cloud resource created successfully!

ID: 507f1f77bcf86cd799439011
Name: my-vpc
Kind: CivoVpc
Created At: 2025-11-28 13:14:12
```

#### Automatic Deployment

When you create a cloud resource, the system automatically:

1. Saves the resource to the database
2. Determines the cloud provider from the resource kind (e.g., `GcpCloudSql` → `gcp`)
3. Resolves credentials from the database based on the provider
4. Triggers a Pulumi deployment with the resolved credentials

**Note:** Credentials are automatically resolved from the database based on the cloud provider. Ensure credentials are configured in the database before creating cloud resources (typically done through the backend API or web console).

#### Flags

- `--arg, -a` - Path to the YAML manifest file (required)
- `--help, -h` - Show help information

#### Manifest Requirements

The YAML manifest must contain:

- `kind` - The type of cloud resource (e.g., `CivoVpc`, `AwsRdsInstance`)
- `metadata.name` - A unique name for the resource

**Example Manifest:**

```yaml
kind: CivoVpc
metadata:
  name: my-vpc
spec:
  region: NYC1
  cidr: 10.0.0.0/16
```

#### Error Handling

**Missing manifest:**

```
Error: --arg flag is required. Provide path to YAML manifest file
Usage: openmcf cloud-resource:create --arg=<yaml-file>
```

**Invalid YAML:**

```
Error: Invalid manifest - invalid YAML format: yaml: line 2: found character that cannot start any token
```

**Duplicate resource:**

```
Error: Invalid manifest - cloud resource with name 'my-vpc' already exists
```

**Connection issues:**

```
Error: Cannot connect to backend service at http://localhost:50051. Please check:
  1. The backend service is running
  2. The backend URL is correct
  3. Network connectivity
```

### `openmcf cloud-resource:list`

List all cloud resources from the backend service with optional filtering by kind.

#### Basic Usage

```bash
# List all cloud resources
openmcf cloud-resource:list
```

**Sample Output:**

```
ID                     NAME      KIND            CREATED
507f1f77bcf86cd799439011  my-vpc   CivoVpc        2025-11-28 13:14:12
507f1f77bcf86cd799439012  my-db    AwsRdsInstance  2025-11-28 13:15:00

Total: 2 cloud resource(s)
```

#### Filtering by Kind

Use the `--kind` flag to filter cloud resources by their kind.

```bash
# Filter by specific kind
openmcf cloud-resource:list --kind CivoVpc
openmcf cloud-resource:list -k AwsRdsInstance
```

**Sample Output:**

```
ID                     NAME      KIND     CREATED
507f1f77bcf86cd799439011  my-vpc   CivoVpc  2025-11-28 13:14:12

Total: 1 cloud resource(s) (filtered by kind: CivoVpc)
```

#### Flags

- `--kind, -k` - Filter cloud resources by kind (optional)
- `--help, -h` - Show help information

#### Output Format

The command displays results in a table with the following columns:

- **ID** - Unique identifier (MongoDB ObjectID)
- **NAME** - Resource name (from metadata.name)
- **KIND** - Resource type/kind (e.g., CivoVpc, AwsRdsInstance)
- **CREATED** - Creation timestamp

### `openmcf cloud-resource:get`

Retrieve detailed information about a specific cloud resource by its ID.

#### Basic Usage

```bash
openmcf cloud-resource:get --id=<resource-id>
```

**Example:**

```bash
# Get a cloud resource by ID
openmcf cloud-resource:get --id=507f1f77bcf86cd799439011
```

**Sample Output:**

```
Cloud Resource Details:
======================
ID:         507f1f77bcf86cd799439011
Name:       my-vpc
Kind:       CivoVpc
Created At: 2025-11-28 13:14:12
Updated At: 2025-11-28 14:05:23

Manifest:
----------
kind: CivoVpc
metadata:
  name: my-vpc
spec:
  region: NYC1
  cidr: 10.0.0.0/16
  description: Production VPC
```

#### Flags

- `--id, -i` - Unique identifier of the cloud resource (required)
- `--help, -h` - Show help information

#### Error Handling

**Missing ID:**

```
Error: --id flag is required. Provide the cloud resource ID
Usage: openmcf cloud-resource:get --id=<resource-id>
```

**Resource not found:**

```
Error: Cloud resource with ID '507f1f77bcf86cd799439011' not found
```

**Invalid ID format:**

```
Error: Invalid manifest - invalid ID format
```

### `openmcf cloud-resource:update`

Update an existing cloud resource by providing a new YAML manifest. The manifest's `name` and `kind` must match the existing resource. This command automatically triggers deployment using credentials resolved from the database.

#### Basic Usage

```bash
openmcf cloud-resource:update --id=<resource-id> --arg=<yaml-file>
```

**Example:**

```bash
# Update a cloud resource
openmcf cloud-resource:update --id=507f1f77bcf86cd799439011 --arg=my-vpc-updated.yaml
```

**Sample Output:**

```
✅ Cloud resource updated successfully!

ID: 507f1f77bcf86cd799439011
Name: my-vpc
Kind: CivoVpc
Updated At: 2025-11-28 14:05:23
```

#### Automatic Deployment

When you update a cloud resource, the system automatically:

1. Updates the resource in the database
2. Determines the cloud provider from the resource kind
3. Resolves credentials from the database based on the provider
4. Triggers a Pulumi deployment with the resolved credentials

**Note:** Credentials are automatically resolved from the database based on the cloud provider. Ensure credentials are configured in the database before updating cloud resources (typically done through the backend API or web console).

#### Flags

- `--id, -i` - Unique identifier of the cloud resource (required)
- `--arg, -a` - Path to the YAML manifest file (required)
- `--help, -h` - Show help information

#### Update Validation

**CRITICAL**: The update operation validates that the manifest's `name` and `kind` match the existing resource to prevent accidental data corruption.

**Validation Rules:**

- Manifest `metadata.name` must match existing resource name
- Manifest `kind` must match existing resource kind
- Resource ID and creation timestamp are preserved

**Example Valid Update:**

```yaml
# Existing resource: name=my-vpc, kind=CivoVpc
# This update will succeed
kind: CivoVpc
metadata:
  name: my-vpc
spec:
  region: NYC1
  cidr: 10.0.0.0/16
  description: Updated description
  tags:
    - production
```

#### Error Handling

**Missing arguments:**

```
Error: --id flag is required. Provide the cloud resource ID
Error: --arg flag is required. Provide path to YAML manifest file
Usage: openmcf cloud-resource:update --id=<resource-id> --arg=<yaml-file>
```

**Resource not found:**

```
Error: Cloud resource with ID '507f1f77bcf86cd799439011' not found
```

**Name mismatch:**

```
Error: Invalid manifest - manifest name 'different-name' does not match existing resource name 'my-vpc'
```

**Kind mismatch:**

```
Error: Invalid manifest - manifest kind 'AwsVpc' does not match existing resource kind 'CivoVpc'
```

**Invalid YAML:**

```
Error: Invalid manifest - invalid YAML format: yaml: line 2: found character that cannot start any token
```

### `openmcf cloud-resource:delete`

Delete a cloud resource by its ID. This operation is irreversible.

#### Basic Usage

```bash
openmcf cloud-resource:delete --id=<resource-id>
```

**Example:**

```bash
# Delete a cloud resource
openmcf cloud-resource:delete --id=507f1f77bcf86cd799439011
```

**Sample Output:**

```
✅ Cloud resource 'my-vpc' deleted successfully
```

#### Flags

- `--id, -i` - Unique identifier of the cloud resource (required)
- `--help, -h` - Show help information

#### Error Handling

**Missing ID:**

```
Error: --id flag is required. Provide the cloud resource ID
Usage: openmcf cloud-resource:delete --id=<resource-id>
```

**Resource not found:**

```
Error: Cloud resource with ID '507f1f77bcf86cd799439011' not found
```

**Connection issues:**

```
Error: Cannot connect to backend service at http://localhost:50051. Please check:
  1. The backend service is running
  2. The backend URL is correct
  3. Network connectivity
```

### `openmcf cloud-resource:apply`

Apply a cloud resource from a YAML manifest file. This command performs an **upsert operation**: if a resource with the same `name` and `kind` already exists, it will be updated; otherwise, a new resource will be created. This command automatically triggers deployment using credentials resolved from the database.

#### Basic Usage

```bash
openmcf cloud-resource:apply --arg=path/to/manifest.yaml
```

**Example:**

```bash
# Apply a cloud resource (create or update)
openmcf cloud-resource:apply --arg=my-vpc.yaml
```

#### Automatic Deployment

When you apply a cloud resource, the system automatically:

1. Creates or updates the resource in the database
2. Determines the cloud provider from the resource kind
3. Resolves credentials from the database based on the provider
4. Triggers a Pulumi deployment with the resolved credentials

**Note:** Credentials are automatically resolved from the database based on the cloud provider. Ensure credentials are configured in the database before applying cloud resources (typically done through the backend API or web console).

#### Key Features

- **Idempotent**: Can be run multiple times safely with the same manifest
- **Declarative**: Declare the desired state, let the system figure out create vs update
- **Name + Kind uniqueness**: Resources are identified by the combination of `metadata.name` and `kind`
- **Kubernetes-style**: Follows the familiar `kubectl apply` pattern

#### Sample Output (Create)

When the resource doesn't exist, it will be created:

```
Applying cloud resource: kind=GcpCloudSql, name=gcp-postgres-example
Checking if resource exists...
✅ Cloud resource created successfully!

Action: Created
ID: 507f1f77bcf86cd799439011
Name: gcp-postgres-example
Kind: GcpCloudSql
Created At: 2025-11-28 13:14:12
Updated At: 2025-11-28 13:14:12

🚀 Pulumi deployment has been triggered automatically.
   Deployment is running in the background.
   Use 'openmcf stack-update:list' to check deployment status.
```

#### Sample Output (Update)

When the resource already exists (same name and kind), it will be updated:

```
Applying cloud resource: kind=GcpCloudSql, name=gcp-postgres-example
Checking if resource exists...
✅ Cloud resource updated successfully!

Action: Updated
ID: 507f1f77bcf86cd799439011
Name: gcp-postgres-example
Kind: GcpCloudSql
Created At: 2025-11-28 13:14:12
Updated At: 2025-11-28 15:30:45

🚀 Pulumi deployment has been triggered automatically.
   Deployment is running in the background.
   Use 'openmcf stack-update:list' to check deployment status.
```

#### Flags

- `--arg, -a` - Path to the YAML manifest file (required)
- `--help, -h` - Show help information

#### Manifest Requirements

The YAML manifest must contain:

- `kind` - The type of cloud resource (e.g., `CivoVpc`, `AwsRdsInstance`)
- `metadata.name` - A unique name for the resource

**Example Manifest:**

```yaml
kind: CivoVpc
metadata:
  name: my-vpc
spec:
  region: NYC1
  cidr: 10.0.0.0/16
  description: Production VPC
```

#### How It Works

1. **Extracts** `metadata.name` and `kind` from the manifest
2. **Calls** the `ApplyCloudResource` API which:
   - Checks if resource exists by `name` + `kind`
   - **Creates** the resource if it doesn't exist
   - **Updates** the resource if it already exists (preserves ID and creation timestamp)
3. **Automatically triggers** Pulumi deployment:
   - Resolves credentials from database based on provider
   - Creates a stack-update with status "in_progress"
   - Executes `pulumi up` asynchronously in the background
4. **Returns** the resource with a flag indicating whether it was created or updated
5. **Displays** deployment status message with instructions to check stack-updates

#### Idempotency

The `apply` command is fully idempotent - you can run it multiple times with the same manifest:

```bash
# First run - creates the resource
$ openmcf cloud-resource:apply --arg=my-vpc.yaml
Action: Created

# Second run - updates the resource (even if nothing changed)
$ openmcf cloud-resource:apply --arg=my-vpc.yaml
Action: Updated

# Third run - still works
$ openmcf cloud-resource:apply --arg=my-vpc.yaml
Action: Updated
```

#### Name + Kind Uniqueness

The combination of `metadata.name` and `kind` uniquely identifies a resource. This means you can have resources with the same name but different kinds:

```bash
# Create a CivoVpc named "my-vpc"
$ cat > civo-vpc.yaml <<EOF
kind: CivoVpc
metadata:
  name: my-vpc
spec:
  region: NYC1
  cidr: 10.0.0.0/16
EOF
$ openmcf cloud-resource:apply --arg=civo-vpc.yaml
Action: Created

# Create an AwsVpc with the same name - this is allowed!
$ cat > aws-vpc.yaml <<EOF
kind: AwsVpc
metadata:
  name: my-vpc
spec:
  region: us-east-1
  cidr: 10.1.0.0/16
EOF
$ openmcf cloud-resource:apply --arg=aws-vpc.yaml
Action: Created

# Now you have TWO resources named "my-vpc" with different kinds
```

#### Use Cases

**1. Initial Resource Creation**

```bash
# Create infrastructure from scratch
openmcf cloud-resource:apply --arg=vpc.yaml
openmcf cloud-resource:apply --arg=database.yaml
openmcf cloud-resource:apply --arg=cache.yaml
```

**2. Configuration Updates**

```bash
# Modify vpc.yaml to change CIDR or add tags
# Then apply the changes
openmcf cloud-resource:apply --arg=vpc.yaml
# The resource is updated automatically
```

**3. GitOps Workflows**

```bash
# In CI/CD pipeline - always apply the latest manifest
git pull origin main
openmcf cloud-resource:apply --arg=manifests/production-vpc.yaml
```

**4. Disaster Recovery**

```bash
# Resources deleted accidentally? Just reapply
openmcf cloud-resource:apply --arg=all-resources/*.yaml
# Creates any missing resources, updates existing ones
```

#### Comparison with Other Commands

**Apply vs Create:**

- `create`: Fails if resource already exists (by name only, regardless of kind)
- `apply`: Creates if not exists, updates if exists (by name AND kind)

**Apply vs Update:**

- `update`: Requires resource ID, fails if resource doesn't exist
- `apply`: No ID needed, works whether resource exists or not

**When to use each:**

- Use `apply` for declarative, idempotent workflows (recommended for most cases)
- Use `create` when you want to ensure a resource doesn't already exist
- Use `update` when you have the resource ID and want explicit update semantics

#### Error Handling

**Missing manifest:**

```
Error: --arg flag is required. Provide path to YAML manifest file
Usage: openmcf cloud-resource:apply --arg=<yaml-file>
```

**Invalid YAML:**

```
Error: Invalid manifest - invalid YAML format: yaml: line 2: found character that cannot start any token
```

**Missing required fields:**

```
Error: Invalid manifest - manifest must contain 'kind' field
Error: Invalid manifest - manifest must contain 'metadata.name' field
```

**Connection issues:**

```
Error: Cannot connect to backend service at http://localhost:50051. Please check:
  1. The backend service is running
  2. The backend URL is correct
  3. Network connectivity
```

### Prerequisites

**Option 1: Using Local Web App (Recommended)**

If you're using the OpenMCF web app, the backend URL is automatically configured:

```bash
openmcf webapp init
openmcf webapp start
```

**Option 2: Using Remote/External Backend**

If connecting to a remote backend, configure the URL manually:

```bash
openmcf config set backend-url <your-backend-url>
```

---

## Credentials

Credentials are cloud provider authentication configurations stored in the backend database. They are automatically used when deploying cloud resources to the corresponding provider.

### `openmcf credential:create`

Create a new cloud provider credential for GCP, AWS, or Azure. The command uses a unified interface with a `--provider` flag to specify the cloud provider type.

#### Basic Usage

```bash
openmcf credential:create --name=<credential-name> --provider=<gcp|aws|azure> [provider-specific-flags]
```

**Examples:**

```bash
# Create a GCP credential
openmcf credential:create \
  --name=my-gcp-production-credential \
  --provider=gcp \
  --service-account-key=~/Downloads/my-project-12345-abcdef.json

# Create an AWS credential
openmcf credential:create \
  --name=my-aws-production-credential \
  --provider=aws \
  --account-id=123456789012 \
  --access-key-id=AKIA... \
  --secret-access-key=... \
  --region=us-east-1

# Create an Azure credential
openmcf credential:create \
  --name=my-azure-production-credential \
  --provider=azure \
  --client-id=... \
  --client-secret=... \
  --tenant-id=... \
  --subscription-id=...
```

**Sample Output:**

```
✅ Credential created successfully!

ID: 507f1f77bcf86cd799439011
Name: my-gcp-production-credential
Provider: GCP
Created At: 2025-12-08 15:30:45

💡 This credential can now be automatically used when deploying GCP resources.
```

#### How It Works

The command handles different providers seamlessly:

**For GCP:**

1. Reads the service account key JSON file
2. Base64-encodes the key content automatically
3. Stores in database with `provider=gcp`

**For AWS:**

1. Collects AWS credentials from command flags
2. Stores in database with `provider=aws`

**For Azure:**

1. Collects Azure credentials from command flags
2. Stores in database with `provider=azure`

#### Automatic Credential Resolution

Once a credential is created, it will be automatically used when deploying resources for that provider:

- **GCP resources**: `GcpCloudSql`, `GcpGke`, `GcpComputeInstance`, etc.
- **AWS resources**: `AwsRdsInstance`, `AwsEksCluster`, `AwsEc2Instance`, etc.
- **Azure resources**: `AzureSqlDatabase`, `AzureAksCluster`, `AzureVm`, etc.

The system automatically:

1. Detects the cloud provider from the resource `kind` field
2. Queries the unified `credentials` database collection filtered by provider
3. Uses the first available credential for that provider

**Example Workflow:**

```bash
# Step 1: Create credentials for different providers
openmcf credential:create \
  --name=my-gcp-prod \
  --provider=gcp \
  --service-account-key=~/gcp-key.json

openmcf credential:create \
  --name=my-aws-prod \
  --provider=aws \
  --account-id=123456789012 \
  --access-key-id=AKIA... \
  --secret-access-key=...

# Step 2: Deploy resources (credentials auto-used based on resource kind)
openmcf cloud-resource:create --arg=gcp-cloudsql.yaml  # Uses GCP credential
openmcf cloud-resource:create --arg=aws-rds.yaml       # Uses AWS credential

# The system automatically:
# - Detects provider from resource kind (GcpCloudSql → gcp, AwsRdsInstance → aws)
# - Fetches matching credential from database
# - Deploys using that credential
```

#### Flags

**Common Flags** (required for all providers):

- `--name, -n` - Name of the credential (required, must be unique)
- `--provider, -p` - Cloud provider: `gcp`, `aws`, or `azure` (required)
- `--help, -h` - Show help information

**GCP-specific Flags** (required when `--provider=gcp`):

- `--service-account-key` - Path to GCP service account key JSON file

**AWS-specific Flags** (required when `--provider=aws`):

- `--account-id` - AWS account ID
- `--access-key-id` - AWS access key ID
- `--secret-access-key` - AWS secret access key
- `--region` - AWS region (optional)
- `--session-token` - AWS session token (optional)

**Azure-specific Flags** (required when `--provider=azure`):

- `--client-id` - Azure client ID
- `--client-secret` - Azure client secret
- `--tenant-id` - Azure tenant ID
- `--subscription-id` - Azure subscription ID

#### Provider-Specific Requirements

**GCP Service Account Key:**

The key file must be valid JSON from Google Cloud Console with appropriate IAM permissions:

- `roles/compute.admin`, `roles/container.admin`, `roles/cloudsql.admin`, etc.

**AWS Credentials:**

Use IAM access keys with appropriate permissions:

- Create via AWS IAM Console > Users > Security Credentials
- Grant policies like `AdministratorAccess` or specific resource policies

**Azure Service Principal:**

Create via Azure CLI or Portal with required RBAC roles:

```bash
az ad sp create-for-rbac --name="myServicePrincipal" --role="Contributor"
```

#### Error Handling

**Missing name:**

```
Error: --name flag is required
Usage: openmcf credential:create --name=<credential-name> --provider=<gcp|aws|azure> [provider-specific-flags]
```

**Missing provider:**

```
Error: --provider flag is required
Usage: openmcf credential:create --name=<name> --provider=<gcp|aws|azure> [provider-specific-flags]
```

**Unsupported provider:**

```
Error: Unsupported provider 'unknown'. Supported providers: gcp, aws, azure
```

**Missing provider-specific flags:**

```
# For GCP
Error: --service-account-key is required for GCP provider

# For AWS
Error: --account-id is required for AWS provider
Error: --access-key-id is required for AWS provider

# For Azure
Error: --client-id is required for Azure provider
```

**File not found (GCP):**

```
Error: failed to read service account key file '/path/to/key.json': no such file or directory
```

**Duplicate credential name:**

```
Error: A credential with name 'my-prod-cred' already exists
```

**Connection issues:**

```
Error: Cannot connect to backend service at http://localhost:50051. Please check:
  1. The backend service is running
  2. The backend URL is correct
  3. Network connectivity
```

#### Security Best Practices

1. **Never commit service account keys to version control**

   ```bash
   # Add to .gitignore
   echo "*.json" >> .gitignore
   echo "*-key.json" >> .gitignore
   ```

2. **Use separate credentials for different environments**

```bash
# GCP Development
openmcf credential:create \
  --name=gcp-dev \
  --provider=gcp \
  --service-account-key=~/keys/gcp-dev.json

# AWS Production
openmcf credential:create \
  --name=aws-prod \
  --provider=aws \
  --account-id=123456789012 \
  --access-key-id=AKIA... \
  --secret-access-key=...
```

3. **Rotate credentials regularly**

   - Create new service account keys periodically
   - Delete old credentials from the database
   - Revoke old keys in Google Cloud Console

4. **Use principle of least privilege**

   - Grant only the minimum required permissions
   - Use separate service accounts for different resource types

5. **Store key files securely**

   ```bash
   # Set restrictive permissions
   chmod 600 ~/keys/*.json

   # Store in a secure location
   mkdir -p ~/.config/gcp-keys
   chmod 700 ~/.config/gcp-keys
   mv ~/Downloads/*-key.json ~/.config/gcp-keys/
   ```

#### Use Cases

**1. Initial Setup**

```bash
# Configure backend
openmcf config set backend-url http://localhost:50051

# Add GCP credential
openmcf credential:create \
  --name=default-gcp \
  --provider=gcp \
  --service-account-key=~/gcp-service-account.json
```

**2. Multi-Environment Setup**

```bash
# Development environment
openmcf credential:create \
  --name=gcp-development \
  --provider=gcp \
  --service-account-key=~/keys/dev-sa-key.json

# Staging environment
openmcf credential:create \
  --name=gcp-staging \
  --provider=gcp \
  --service-account-key=~/keys/staging-sa-key.json

# Production environment
openmcf credential:create \
  --name=gcp-production \
  --provider=gcp \
  --service-account-key=~/keys/prod-sa-key.json
```

**3. Multi-Project Setup**

```bash
# Project A
openmcf credential:create \
  --name=gcp-project-a \
  --provider=gcp \
  --service-account-key=~/keys/project-a-key.json

# Project B
openmcf credential:create \
  --name=gcp-project-b \
  --provider=gcp \
  --service-account-key=~/keys/project-b-key.json
```

**4. Team Collaboration**

```bash
# Each team member creates their own credential
openmcf credential:create \
  --name=gcp-alice-dev \
  --provider=gcp \
  --service-account-key=~/alice-dev-key.json

openmcf credential:create \
  --name=gcp-bob-dev \
  --provider=gcp \
  --service-account-key=~/bob-dev-key.json
```

#### Comparison with CLI Flags

**Before (using --gcp-credential flag):**

```bash
# Had to provide credential file every time
openmcf pulumi up \
  --manifest resource.yaml \
  --gcp-credential ~/gcp-key.json
```

**Now (using database-stored credentials):**

```bash
# Create credential once
openmcf credential:create \
  --name=my-gcp-cred \
  --provider=gcp \
  --service-account-key=~/gcp-key.json

# Use cloud-resource commands without credential flags
openmcf cloud-resource:create --arg=resource.yaml
openmcf cloud-resource:apply --arg=resource.yaml
# Credentials are automatically resolved!
```

**Benefits:**

- ✅ No need to specify credentials for every deployment
- ✅ Centralized credential management
- ✅ Credentials are securely stored in the database
- ✅ Easier team collaboration
- ✅ Supports multiple credentials per provider

#### Prerequisites

**Option 1: Using Local Web App (Recommended)**

If you're using the OpenMCF web app:

```bash
# Backend URL is automatically configured during init
openmcf webapp init
openmcf webapp start
```

**Option 2: Using Remote/External Backend**

If connecting to a remote backend:

```bash
openmcf config set backend-url <your-backend-url>
```

### `openmcf credential:list`

List all stored cloud provider credentials with optional filtering by provider type.

#### Basic Usage

```bash
# List all credentials
openmcf credential:list

# List credentials for a specific provider
openmcf credential:list --provider=gcp
openmcf credential:list --provider=aws
openmcf credential:list --provider=azure
```

**Sample Output (all credentials):**

```
ID                       NAME                           PROVIDER   CREATED
------------------------------------------------------------------------------------
507f1f77bcf86cd799439011 my-gcp-production-credential   GCP        2025-12-08 15:30:45
507f1f77bcf86cd799439012 my-aws-production-credential   AWS        2025-12-08 15:32:10
507f1f77bcf86cd799439013 my-azure-production-credenti   AZURE      2025-12-08 15:33:22

Total: 3 credential(s)
```

**Sample Output (filtered by provider):**

```bash
$ openmcf credential:list --provider=gcp
```

```
ID                       NAME                           PROVIDER   CREATED
------------------------------------------------------------------------------------
507f1f77bcf86cd799439011 my-gcp-production-credential   GCP        2025-12-08 15:30:45
507f1f77bcf86cd799439014 my-gcp-dev-credential          GCP        2025-12-08 15:34:00

Total: 2 credential(s) (filtered by provider: GCP)
```

#### Flags

- `--provider, -p` - Filter by cloud provider: `gcp`, `aws`, or `azure` (optional)
- `--help, -h` - Show help information

#### Output Format

The command displays results in a table with the following columns:

- **ID** - Unique credential identifier (MongoDB ObjectID)
- **NAME** - Credential name
- **PROVIDER** - Cloud provider type (GCP, AWS, AZURE)
- **CREATED** - Creation timestamp

**Note:** Sensitive credential data (keys, secrets, passwords) is not displayed in the list output for security reasons.

#### Use Cases

**1. View All Credentials**

```bash
openmcf credential:list
```

**2. Check GCP Credentials**

```bash
openmcf credential:list --provider=gcp
```

**3. Verify Credentials Before Deployment**

```bash
# Check if AWS credentials exist before deploying AWS resources
openmcf credential:list --provider=aws

# If credentials exist, proceed with deployment
openmcf cloud-resource:create --arg=aws-resource.yaml
```

**4. Audit Credentials**

```bash
# List all credentials to audit what's configured
openmcf credential:list

# Check when credentials were created
openmcf credential:list --provider=azure
```

#### Error Handling

**Backend URL Not Configured:**

```
Error: backend URL not configured. Run: openmcf config set backend-url <url>
```

**Unsupported Provider:**

```
Error: Unsupported provider 'unknown'. Supported providers: gcp, aws, azure
```

**Connection Issues:**

```
Error: Cannot connect to backend service at http://localhost:50051. Please check:
  1. The backend service is running
  2. The backend URL is correct
  3. Network connectivity
```

**No Credentials Found:**

```
# When no credentials exist
No credentials found

# When no credentials exist for specific provider
No credentials found for provider: GCP
```

#### Prerequisites

Before using the credential:list command, ensure the backend is configured (automatically done if using `planton webapp init`):

```bash
# If using remote backend, configure manually:
openmcf config set backend-url <your-backend-url>
```

---

### `openmcf credential:get`

Retrieve detailed information about a credential by providing its unique ID. This command displays all credential metadata and masked sensitive data for security.

#### Basic Usage

```bash
openmcf credential:get --id=<credential-id>
```

**Example:**

```bash
openmcf credential:get --id=507f1f77bcf86cd799439011
```

**Sample Output:**

```
Credential Details:
===================
ID:         507f1f77bcf86cd799439011
Name:       my-gcp-production-credential
Provider:   GCP
Created At: 2025-12-08 15:30:45
Updated At: 2025-12-08 16:20:10

Credential Data:
----------------
Service Account Key (Base64): ewog...IA==
```

#### Flags

- `--id, -i` - Unique identifier of the credential (required)
- `--help, -h` - Show help information

#### Output Format

The command displays:

- **ID** - Unique credential identifier
- **Name** - Credential name
- **Provider** - Cloud provider type (GCP, AWS, Azure)
- **Created At** - Creation timestamp
- **Updated At** - Last update timestamp (if available)
- **Credential Data** - Provider-specific credential information with sensitive data masked

**Security Note:** Sensitive data (keys, secrets, passwords) is automatically masked in the output. Only partial data is shown (first 4 and last 4 characters) for security purposes.

#### Use Cases

**1. Verify Credential Details**

```bash
# Get full details of a credential
openmcf credential:get --id=507f1f77bcf86cd799439011
```

**2. Check Credential Before Deployment**

```bash
# List credentials to find ID
openmcf credential:list

# Get details of specific credential
openmcf credential:get --id=<credential-id>
```

**3. Audit Credential Configuration**

```bash
# Verify credential settings after creation or update
openmcf credential:get --id=507f1f77bcf86cd799439011
```

#### Error Handling

**Missing ID:**

```
Error: required flag(s) "id" not set
Usage: openmcf credential:get --id=<credential-id>
```

**Invalid ID Format:**

```
Error getting credential: internal: failed to get credential: invalid ID format: the provided hex string is not a valid ObjectID
```

**Credential Not Found:**

```
Error: Credential with ID '507f1f77bcf86cd799439011' not found
```

**Connection Issues:**

```
Error: Cannot connect to backend service at http://localhost:50051. Please check:
  1. The backend service is running
  2. The backend URL is correct
  3. Network connectivity
```

#### Prerequisites

Before using the credential:get command, ensure the backend is configured (automatically done if using `planton webapp init`):

```bash
# If using remote backend, configure manually:
openmcf config set backend-url <your-backend-url>
```

---

### `openmcf credential:update`

Update an existing cloud provider credential. The provider type must match the existing credential. You can update the credential name and all provider-specific credential data.

#### Basic Usage

```bash
openmcf credential:update --id=<credential-id> --name=<new-name> --provider=<gcp|aws|azure> [provider-specific-flags]
```

**Examples:**

```bash
# Update a GCP credential
openmcf credential:update \
  --id=507f1f77bcf86cd799439011 \
  --name=updated-gcp-credential \
  --provider=gcp \
  --service-account-key=~/Downloads/new-gcp-key.json

# Update an AWS credential
openmcf credential:update \
  --id=507f1f77bcf86cd799439012 \
  --name=updated-aws-credential \
  --provider=aws \
  --account-id=123456789012 \
  --access-key-id=AKIA... \
  --secret-access-key=... \
  --region=us-west-2

# Update an Azure credential
openmcf credential:update \
  --id=507f1f77bcf86cd799439013 \
  --name=updated-azure-credential \
  --provider=azure \
  --client-id=12345678-1234-1234-1234-123456789012 \
  --client-secret=new-client-secret-here \
  --tenant-id=87654321-4321-4321-4321-210987654321 \
  --subscription-id=11111111-2222-3333-4444-555555555555
```

**Sample Output:**

```
✅ Credential updated successfully!

ID:       507f1f77bcf86cd799439011
Name:     updated-gcp-credential
Provider: GCP
Updated At: 2025-12-08 16:20:10
```

#### Flags

**Common Flags** (required for all providers):

- `--id, -i` - Unique identifier of the credential (required)
- `--name, -n` - New name for the credential (required)
- `--provider, -p` - Cloud provider: `gcp`, `aws`, or `azure` (required, must match existing credential)
- `--help, -h` - Show help information

**GCP-specific Flags** (required when `--provider=gcp`):

- `--service-account-key` - Path to new GCP service account key JSON file

**AWS-specific Flags** (required when `--provider=aws`):

- `--account-id` - AWS account ID
- `--access-key-id` - AWS access key ID
- `--secret-access-key` - AWS secret access key
- `--region` - AWS region (optional)
- `--session-token` - AWS session token (optional)

**Azure-specific Flags** (required when `--provider=azure`):

- `--client-id` - Azure client ID
- `--client-secret` - Azure client secret
- `--tenant-id` - Azure tenant ID
- `--subscription-id` - Azure subscription ID

#### Use Cases

**1. Rotate Credentials**

```bash
# Update with new service account key
openmcf credential:update \
  --id=507f1f77bcf86cd799439011 \
  --name=my-gcp-prod \
  --provider=gcp \
  --service-account-key=~/new-key.json
```

**2. Update Credential Name**

```bash
# Rename credential while keeping same credentials
openmcf credential:update \
  --id=507f1f77bcf86cd799439011 \
  --name=new-credential-name \
  --provider=gcp \
  --service-account-key=~/existing-key.json
```

**3. Change AWS Region**

```bash
# Update AWS credential with new region
openmcf credential:update \
  --id=507f1f77bcf86cd799439012 \
  --name=my-aws-prod \
  --provider=aws \
  --account-id=123456789012 \
  --access-key-id=AKIA... \
  --secret-access-key=... \
  --region=us-west-2
```

#### Error Handling

**Missing Required Flags:**

```
Error: --id, --name, and --provider flags are required
Usage: openmcf credential:update --id=<id> --name=<name> --provider=<provider> [provider-specific-flags]
```

**Missing Provider-Specific Flags:**

```
# For GCP
Error: --service-account-key is required for GCP provider

# For AWS
Error: --account-id, --access-key-id, and --secret-access-key are required for AWS provider

# For Azure
Error: --client-id, --client-secret, --tenant-id, and --subscription-id are required for Azure provider
```

**Invalid Provider:**

```
Error: Invalid provider 'unknown'. Valid values: gcp, aws, azure
```

**Credential Not Found:**

```
Error: Credential with ID '507f1f77bcf86cd799439011' not found
```

**File Not Found (GCP):**

```
Error: Failed to read service account key file '/path/to/key.json': open /path/to/key.json: no such file or directory
```

**Connection Issues:**

```
Error: Cannot connect to backend service at http://localhost:50051. Please check:
  1. The backend service is running
  2. The backend URL is correct
  3. Network connectivity
```

#### Prerequisites

Before using the credential:update command, ensure the backend is configured (automatically done if using `planton webapp init`):

```bash
# If using remote backend, configure manually:
openmcf config set backend-url <your-backend-url>
```

---

### `openmcf credential:delete`

Delete a credential by providing its unique ID. This action is irreversible and will permanently remove the credential from the database.

#### Basic Usage

```bash
openmcf credential:delete --id=<credential-id>
```

**Example:**

```bash
openmcf credential:delete --id=507f1f77bcf86cd799439011
```

**Sample Output:**

```
✅ Credential deleted successfully
```

#### Flags

- `--id, -i` - Unique identifier of the credential (required)
- `--help, -h` - Show help information

#### Use Cases

**1. Remove Unused Credentials**

```bash
# List credentials to find ID
openmcf credential:list

# Delete specific credential
openmcf credential:delete --id=507f1f77bcf86cd799439011
```

**2. Clean Up Old Credentials**

```bash
# Delete credentials that are no longer needed
openmcf credential:delete --id=507f1f77bcf86cd799439011
```

**3. Rotate Credentials**

```bash
# After creating new credential, delete old one
openmcf credential:create --name=new-cred --provider=gcp --service-account-key=~/new-key.json
openmcf credential:delete --id=<old-credential-id>
```

#### Error Handling

**Missing ID:**

```
Error: required flag(s) "id" not set
Usage: openmcf credential:delete --id=<credential-id>
```

**Invalid ID Format:**

```
Error deleting credential: internal: failed to delete credential: invalid ID format: the provided hex string is not a valid ObjectID
```

**Credential Not Found:**

```
Error: Credential with ID '507f1f77bcf86cd799439011' not found
```

**Connection Issues:**

```
Error: Cannot connect to backend service at http://localhost:50051. Please check:
  1. The backend service is running
  2. The backend URL is correct
  3. Network connectivity
```

#### Prerequisites

Before using the credential:delete command, ensure the backend is configured (automatically done if using `planton webapp init`):

```bash
# If using remote backend, configure manually:
openmcf config set backend-url <your-backend-url>
```

#### Important Notes

- ⚠️ **This operation is irreversible** - Once deleted, the credential cannot be recovered
- ⚠️ **Impact on Deployments** - If this credential is being used by active deployments, those deployments may fail
- ✅ **Safe to Delete** - You can safely delete credentials that are no longer in use

---

## Future Enhancements

Additional credential providers coming soon:

- Support for additional providers: Cloudflare, Atlas, Confluent, Snowflake

---

## Stack Jobs

Stack jobs represent deployment operations for cloud resources. You can stream real-time output from stack-updates to monitor deployment progress.

### `openmcf stack-update:stream-output`

Stream real-time deployment logs from a stack-update. Shows stdout and stderr output as it's generated during deployment.

#### Basic Usage

```bash
# Stream output from a stack-update
openmcf stack-update:stream-output --id=<stack-update-id>
```

**Sample Output:**

```
🚀 Streaming output for stack-update: 69369e4ec78ad326a6e5aa8b

[15:04:05.123] [stdout] [Seq: 1] Updating (example-env.GcpCloudSql.gcp-postgres-example):
[15:04:05.234] [stdout] [Seq: 2]     pulumi:pulumi:Stack openmcf-examples-example-env.GcpCloudSql.gcp-postgres-example  Compiling the program ...
[15:04:06.456] [stdout] [Seq: 3]     pulumi:pulumi:Stack openmcf-examples-example-env.GcpCloudSql.gcp-postgres-example  Finished compiling
[15:04:07.789] [stdout] [Seq: 4] +  gcp:sql:DatabaseInstance gcp-postgres-example creating (0s)
[15:04:10.123] [stdout] [Seq: 5] +  gcp:sql:DatabaseInstance gcp-postgres-example created (3s)

✅ Stream completed successfully

📊 Total messages received: 5 (last sequence: 5)
```

#### Resuming from a Specific Sequence

If you need to resume streaming from a specific point (e.g., after disconnection), use the `--last-sequence` flag:

```bash
# Resume from sequence 100
openmcf stack-update:stream-output --id=<stack-update-id> --last-sequence=100
```

**Sample Output:**

```
🚀 Streaming output for stack-update: 69369e4ec78ad326a6e5aa8b
   Resuming from sequence: 100

[15:05:01.234] [stdout] [Seq: 101] Continuing deployment...
[15:05:02.456] [stdout] [Seq: 102] Finalizing resources...
```

#### Flags

- `--id, -i` - Unique identifier of the stack-update (required)
- `--last-sequence, -s` - Last sequence number received (for resuming stream from a specific point, default: 0)
- `--help, -h` - Show help information

#### Output Format

Each stream message displays:

- **Timestamp** - Time when the message was generated (HH:MM:SS.mmm format)
- **Stream Type** - `[stdout]` for standard output or `[stderr]` for error output
- **Sequence** - Sequence number of the message (for ordering and resuming)
- **Content** - The actual log line content

#### Graceful Shutdown

The stream command supports graceful shutdown via interrupt signals:

- Press `Ctrl+C` to cancel the stream
- The command will finish processing the current message and exit cleanly
- A summary showing total messages received and last sequence number will be displayed

**Example:**

```bash
$ openmcf stack-update:stream-output --id=69369e4ec78ad326a6e5aa8b
🚀 Streaming output for stack-update: 69369e4ec78ad326a6e5aa8b

[15:04:05.123] [stdout] [Seq: 1] Starting deployment...
[15:04:06.456] [stdout] [Seq: 2] Compiling program...
^C

⚠️  Interrupt received, stopping stream...

⚠️  Stream cancelled

📊 Total messages received: 2 (last sequence: 2)
```

#### Error Handling

**Backend URL Not Configured:**

```
Error: backend URL not configured. Run: openmcf config set backend-url <url>
```

**Solution:**

```bash
openmcf config set backend-url http://localhost:50051
```

**Connection Issues:**

```
❌ Error: Cannot connect to backend service at http://localhost:50051. Please check:
  1. The backend service is running
  2. The backend URL is correct
  3. Network connectivity
```

**Solutions:**

1. **Check if backend service is running:**

   ```bash
   # Check if port is accessible
   curl http://localhost:50051
   ```

2. **Verify backend URL configuration:**

   ```bash
   openmcf config get backend-url
   ```

3. **Update backend URL if needed:**
   ```bash
   openmcf config set backend-url <correct-url>
   ```

**Stack Job Not Found:**

```
❌ Error: Stack job with ID 'invalid-id' not found
```

**Solution:**

- Verify the stack-update ID is correct
- Use `openmcf cloud-resource:get` to find the associated cloud resource and its stack-updates

**Stream Error:**

```
❌ Stream error: <error details>
```

**Possible Causes:**

- Backend service disconnected during streaming
- Network interruption
- Backend service error

**Solutions:**

1. Check backend service logs
2. Verify network connectivity
3. Retry the stream command
4. Use `--last-sequence` to resume from the last received sequence number

#### Use Cases

**1. Monitoring Active Deployments**

```bash
# Stream output from an in-progress deployment
openmcf stack-update:stream-output --id=<stack-update-id>
```

**2. Reviewing Completed Deployments**

```bash
# Stream all logs from a completed deployment
openmcf stack-update:stream-output --id=<stack-update-id>
```

**3. Resuming After Disconnection**

```bash
# If disconnected, resume from the last sequence number you saw
openmcf stack-update:stream-output --id=<stack-update-id> --last-sequence=150
```

#### Prerequisites

Before using the stack-update:stream-output command, ensure the backend is configured (automatically done if using `planton webapp init`):

```bash
# If using remote backend, configure manually:
openmcf config set backend-url <your-backend-url>
```

---

## Common Workflows

### Initial Setup

1. **Configure the backend URL:**

   ```bash
   openmcf config set backend-url http://localhost:50051
   ```

2. **Verify configuration:**

   ```bash
   openmcf config get backend-url
   ```

3. **Test connectivity:**
   ```bash
   openmcf list-deployment-components
   ```

### Daily Usage

1. **List all available deployment components:**

   ```bash
   openmcf list-deployment-components
   ```

2. **Find specific component types:**

   ```bash
   # List all Kubernetes components
   openmcf list-deployment-components --kind PostgresKubernetes

   # List all AWS components
   openmcf list-deployment-components --kind AwsRdsInstance
   ```

3. **Manage cloud resources:**

   ```bash
   # Apply a cloud resource (recommended - works for create and update)
   openmcf cloud-resource:apply --arg=my-vpc.yaml

   # Or use explicit create/update commands
   openmcf cloud-resource:create --arg=my-vpc.yaml

   # Get resource details by ID
   openmcf cloud-resource:get --id=507f1f77bcf86cd799439011

   # Update a resource
   openmcf cloud-resource:update --id=507f1f77bcf86cd799439011 --arg=updated.yaml

   # Delete a resource
   openmcf cloud-resource:delete --id=507f1f77bcf86cd799439011
   ```

4. **List cloud resources:**

   ```bash
   # List all cloud resources
   openmcf cloud-resource:list

   # Filter by kind
   openmcf cloud-resource:list --kind CivoVpc
   ```

5. **Check configuration:**

   ```bash
   openmcf config list
   ```

### Environment-Specific Setup

#### Local Development

```bash
openmcf config set backend-url http://localhost:50051
```

#### Staging Environment

```bash
openmcf config set backend-url https://staging-api.openmcf.com
```

#### Production Environment

```bash
openmcf config set backend-url https://api.openmcf.com
```

---

## Troubleshooting

### Backend URL Not Configured

**Error:**

```
Error: backend URL not configured. Run: openmcf config set backend-url <url>
```

**Solution:**

```bash
openmcf config set backend-url http://localhost:50051
```

### Connection Refused

**Error:**

```
Error: Cannot connect to backend service at http://localhost:50051. Please check:
  1. The backend service is running
  2. The backend URL is correct
  3. Network connectivity
```

**Solutions:**

1. **Check if backend service is running:**

   ```bash
   # If using Docker
   docker ps | grep backend

   # Check if port is accessible
   curl http://localhost:50051
   ```

2. **Verify backend URL configuration:**

   ```bash
   openmcf config get backend-url
   ```

3. **Update backend URL if needed:**
   ```bash
   openmcf config set backend-url <correct-url>
   ```

### Invalid Backend URL

**Error:**

```
Error: backend-url must start with http:// or https://
```

**Solution:**

```bash
# Correct format
openmcf config set backend-url http://localhost:50051
openmcf config set backend-url https://api.example.com
```

### No Results Found

**Output:**

```
No deployment components found
# or
No deployment components found with kind 'YourKind'
```

**Possible Causes:**

1. Backend database is empty
2. Wrong kind filter applied
3. Backend service not properly initialized

**Solutions:**

1. **Check without filters:**

   ```bash
   openmcf list-deployment-components
   ```

2. **Verify backend service logs**

3. **Check available kinds by listing all components first**

### Cloud Resource Creation Errors

**Invalid Manifest:**

```
Error: Invalid manifest - invalid YAML format: yaml: line 2: found character that cannot start any token
```

**Solution:**

- Verify the YAML file is valid
- Ensure `kind` and `metadata.name` fields are present
- Check YAML syntax (indentation, quotes, etc.)

**Duplicate Resource Name:**

```
Error: Invalid manifest - cloud resource with name 'my-vpc' already exists
```

**Solution:**

- Use a different name for the resource
- Check existing resources: `openmcf cloud-resource:list`
- Delete the existing resource if needed: `openmcf cloud-resource:delete --id=<id>`

### Cloud Resource Update Errors

**Name Mismatch:**

```
Error: Invalid manifest - manifest name 'different-name' does not match existing resource name 'my-vpc'
```

**Solution:**

- Ensure the manifest `metadata.name` matches the existing resource name
- Get current resource details: `openmcf cloud-resource:get --id=<id>`
- Update the manifest to use the correct name

**Kind Mismatch:**

```
Error: Invalid manifest - manifest kind 'AwsVpc' does not match existing resource kind 'CivoVpc'
```

**Solution:**

- Ensure the manifest `kind` matches the existing resource kind
- If you need to change the kind, delete and recreate the resource
- Get current resource details: `openmcf cloud-resource:get --id=<id>`

**Resource Not Found:**

```
Error: Cloud resource with ID '507f1f77bcf86cd799439011' not found
```

**Solution:**

- Verify the resource ID is correct
- List all resources: `openmcf cloud-resource:list`
- The resource may have been deleted

### Cloud Resource Deletion Errors

**Resource Not Found:**

```
Error: Cloud resource with ID '507f1f77bcf86cd799439011' not found
```

**Solution:**

- Verify the resource ID is correct
- List all resources: `openmcf cloud-resource:list`
- The resource may have already been deleted

**Empty Results:**

```
No cloud resources found
# or
No cloud resources found with kind 'YourKind'
```

**Possible Causes:**

1. No cloud resources have been created yet
2. Wrong kind filter applied
3. Backend database is empty

**Solutions:**

1. **Check without filters:**

   ```bash
   openmcf cloud-resource:list
   ```

2. **Create a test resource:**

   ```bash
   openmcf cloud-resource:create --arg=test-resource.yaml
   ```

3. **Verify backend service is running and initialized**

### Configuration File Issues

**Error:** Permission denied or file access issues

**Solutions:**

1. **Check file permissions:**

   ```bash
   ls -la ~/.openmcf/
   ```

2. **Reset configuration directory:**
   ```bash
   rm -rf ~/.openmcf/
   openmcf config set backend-url <your-url>
   ```

### Network Connectivity

**Error:** Timeout or DNS resolution issues

**Solutions:**

1. **Test basic connectivity:**

   ```bash
   ping <your-backend-host>
   curl <your-backend-url>
   ```

2. **Check firewall/proxy settings**

3. **Try different backend URL (HTTP vs HTTPS)**

---

## Advanced Usage

### Scripting

The CLI commands are designed to be script-friendly:

```bash
#!/bin/bash

# Check if backend is configured
if ! openmcf config get backend-url > /dev/null 2>&1; then
    echo "Backend not configured"
    exit 1
fi

# Get component count
COMPONENT_COUNT=$(openmcf list-deployment-components | grep "Total:" | grep -o '[0-9]\+')
echo "Found $COMPONENT_COUNT deployment components"

# List specific kinds
for kind in PostgresKubernetes AwsRdsInstance GcpCloudSql; do
    echo "=== $kind ==="
    openmcf list-deployment-components --kind "$kind"
done

# Apply cloud resources from directory (recommended - idempotent)
echo "=== Applying resources ==="
for manifest in resources/*.yaml; do
    echo "Applying resource from $manifest"
    openmcf cloud-resource:apply --arg="$manifest"
done

# Or use explicit create for new resources
for manifest in resources/*.yaml; do
    echo "Creating resource from $manifest"
    openmcf cloud-resource:create --arg="$manifest"
done

# List all cloud resources and get details
openmcf cloud-resource:list | grep -v "^Total:" | tail -n +2 | while read -r id name kind created; do
    echo "=== Resource: $name ($kind) ==="
    openmcf cloud-resource:get --id="$id"
    echo ""
done

# Apply updates to resources (no ID needed!)
echo "=== Applying updates ==="
for manifest in updates/*.yaml; do
    name=$(grep "name:" "$manifest" | awk '{print $2}')
    kind=$(grep "kind:" "$manifest" | awk '{print $2}')
    echo "Applying $kind/$name from $manifest"
    openmcf cloud-resource:apply --arg="$manifest"
done

# Cleanup old resources
echo "=== Cleaning up old resources ==="
openmcf cloud-resource:list --kind TestResource | grep -v "^Total:" | tail -n +2 | while read -r id rest; do
    echo "Deleting test resource: $id"
    openmcf cloud-resource:delete --id="$id"
done
```

### Complete Cloud Resource Lifecycle

#### Using Apply (Recommended - Simpler)

```bash
#!/bin/bash

# Complete workflow example using apply command
set -e

# 1. Apply a resource (creates it)
echo "Applying VPC resource..."
cat > temp-vpc.yaml <<EOF
kind: CivoVpc
metadata:
  name: automation-vpc
spec:
  region: NYC1
  cidr: 10.0.0.0/16
EOF

# First apply creates the resource
OUTPUT=$(openmcf cloud-resource:apply --arg=temp-vpc.yaml)
echo "$OUTPUT"
RESOURCE_ID=$(echo "$OUTPUT" | grep "^ID:" | awk '{print $2}')
echo "Resource ID: $RESOURCE_ID"

# 2. Get resource details
echo "Fetching resource details..."
openmcf cloud-resource:get --id="$RESOURCE_ID"

# 3. Modify and apply again (updates it)
echo "Updating resource..."
cat > temp-vpc.yaml <<EOF
kind: CivoVpc
metadata:
  name: automation-vpc
spec:
  region: NYC1
  cidr: 10.0.0.0/16
  description: Updated via automation
  tags:
    - automated
    - production
EOF

# Apply again - automatically updates the resource
openmcf cloud-resource:apply --arg=temp-vpc.yaml

# 4. Verify update
echo "Verifying update..."
openmcf cloud-resource:get --id="$RESOURCE_ID" | grep "description"

# 5. Apply is idempotent - run it again, still works
echo "Applying again (idempotency test)..."
openmcf cloud-resource:apply --arg=temp-vpc.yaml

# 6. Delete resource
echo "Cleaning up..."
openmcf cloud-resource:delete --id="$RESOURCE_ID"

# Cleanup temp file
rm temp-vpc.yaml

echo "Workflow complete!"
```

#### Using Create/Update (Explicit)

```bash
#!/bin/bash

# Complete workflow example using explicit create/update
set -e

# 1. Create a resource
echo "Creating VPC resource..."
cat > temp-vpc.yaml <<EOF
kind: CivoVpc
metadata:
  name: automation-vpc
spec:
  region: NYC1
  cidr: 10.0.0.0/16
EOF

RESOURCE_ID=$(openmcf cloud-resource:create --arg=temp-vpc.yaml | grep "^ID:" | awk '{print $2}')
echo "Created resource with ID: $RESOURCE_ID"

# 2. Get resource details
echo "Fetching resource details..."
openmcf cloud-resource:get --id="$RESOURCE_ID"

# 3. Update the resource
echo "Updating resource..."
cat > temp-vpc.yaml <<EOF
kind: CivoVpc
metadata:
  name: automation-vpc
spec:
  region: NYC1
  cidr: 10.0.0.0/16
  description: Updated via automation
EOF

openmcf cloud-resource:update --id="$RESOURCE_ID" --arg=temp-vpc.yaml

# 4. Verify update
echo "Verifying update..."
openmcf cloud-resource:get --id="$RESOURCE_ID" | grep "description"

# 5. Delete resource
echo "Cleaning up..."
openmcf cloud-resource:delete --id="$RESOURCE_ID"

# Cleanup temp file
rm temp-vpc.yaml

echo "Workflow complete!"
```

### GitOps-Style Infrastructure Management

```bash
#!/bin/bash

# GitOps workflow - sync infrastructure from Git repository
set -e

MANIFEST_DIR="infrastructure/manifests"

echo "=== Syncing infrastructure from Git ==="

# Pull latest changes
git pull origin main

# Apply all resources (creates new, updates existing)
for manifest in "$MANIFEST_DIR"/*.yaml; do
    echo "Applying $(basename $manifest)..."

    # Parse manifest for name and kind
    name=$(grep "name:" "$manifest" | head -1 | awk '{print $2}')
    kind=$(grep "kind:" "$manifest" | head -1 | awk '{print $2}')

    # Apply the resource
    OUTPUT=$(openmcf cloud-resource:apply --arg="$manifest")

    # Check if created or updated
    if echo "$OUTPUT" | grep -q "Action: Created"; then
        echo "✅ Created $kind/$name"
    else
        echo "✅ Updated $kind/$name"
    fi
done

echo "=== Infrastructure sync complete ==="

# List all resources to verify
echo "Current infrastructure state:"
openmcf cloud-resource:list
```

### JSON Output (Future Enhancement)

Currently, the CLI outputs human-readable tables. JSON output support may be added in future versions:

```bash
# Future enhancement
openmcf list-deployment-components --output json
```

---

## Support

For additional help:

- Check the main CLI help: `openmcf --help`
- Command-specific help: `openmcf <command> --help`
- Project documentation: [OpenMCF Documentation](https://openmcf.org)
- GitHub Issues: [Report Issues](https://github.com/plantonhq/openmcf/issues)
