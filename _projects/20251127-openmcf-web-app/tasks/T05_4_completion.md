# T05: Pulumi CLI Stack Job API Implementation

**Status:** ✅ COMPLETED
**Date:** December 3, 2025
**Type:** Feature
**Changelog:** `2025-12-03-171627-pulumi-cli-stack-update-api-implementation.md`

---

## Overview

Implemented a comprehensive Stack Job API service that enables asynchronous Pulumi CLI deployments for cloud resources. The system provides gRPC APIs to deploy cloud resources using Pulumi, track deployment jobs, and retrieve deployment status and output.

## What Was Accomplished

### 1. Stack Job Service API

Three main RPC methods:

**DeployCloudResource:**
- Initiates deployment for a cloud resource
- Validates cloud resource exists
- Creates stack-update record
- Returns immediately with job ID
- Executes Pulumi deployment asynchronously

**GetStackUpdate:**
- Retrieves a specific stack-update by ID
- Returns job status, output, and timestamps
- Used for polling deployment status

**ListStackUpdates:**
- Lists stack-updates with optional filters
- Filter by cloud resource ID
- Filter by status (success, failed, in_progress)
- Sorted by creation date (newest first)

### 2. Pulumi CLI Integration

Executes Pulumi CLI commands directly:
```bash
pulumi up --stack <stack_fqdn> --yes --skip-preview
```

**Key capabilities:**
- Extracts stack FQDN from manifest labels
- Builds stack input YAML from manifest
- Sets working directory to Pulumi module path
- Sets `STACK_INPUT_YAML` environment variable
- Captures stdout, stderr, and exit codes
- 10-minute timeout for deployments
- Handles errors gracefully with detailed messages

### 3. Deployment Output Format

Results stored as JSON in the `output` field:
```json
{
  "status": "success" | "failed",
  "timestamp": "2025-12-03T17:16:27Z",
  "stack_fqdn": "org/project/stack",
  "stdout": "...",
  "stderr": "...",
  "exit_code": 0,
  "error": "..." // only on failure
}
```

### 4. Asynchronous Execution

- Deployments run in background goroutines
- API returns immediately with job ID
- Clients can poll for status using GetStackUpdate
- No request timeouts for long-running deployments

### 5. Docker Integration

- Pulumi CLI v3.206.0 installed in backend Docker image
- Environment variables configured (PULUMI_HOME, PULUMI_STATE_DIR)
- Volume mounts for Pulumi state persistence

## Technical Implementation

### Data Model

```go
type StackUpdate struct {
    ID              primitive.ObjectID `bson:"_id,omitempty"`
    CloudResourceID string             `bson:"cloud_resource_id"`
    Status          string             `bson:"status"` // success, failed, in_progress
    Output          string             `bson:"output,omitempty"` // JSON string
    CreatedAt       time.Time          `bson:"created_at"`
    UpdatedAt       time.Time          `bson:"updated_at"`
}
```

**MongoDB Collection:** `stackupdates`

### Deployment Flow

```
1. User calls DeployCloudResource API
2. Backend creates StackUpdate with "in_progress" status
3. Returns job immediately
4. Background goroutine starts
5. Load manifest, extract stack FQDN
6. Build stack input YAML
7. Execute: pulumi up --stack <fqdn> --yes --skip-preview
8. Capture stdout/stderr/exit_code
9. Update StackUpdate with status and output
```

## Files Created

### Backend API
- `app/backend/apis/proto/stack_job_service.proto` - gRPC service definition
- `app/backend/internal/service/stack_job_service.go` - Service implementation
- `app/backend/internal/database/stack_job_repo.go` - Repository layer
- `app/backend/pkg/models/stack_job.go` - Data model

## Files Modified

### Backend
- `app/backend/go.mod` - Updated dependencies
- `app/backend/apis/buf.gen.yaml` - Fixed TypeScript proto output paths
- `app/backend/apis/Makefile` - Fixed TypeScript stubs cleanup paths

### Infrastructure
- `Dockerfile.backend` - Added Pulumi CLI installation
- `docker-compose.yml` - Added Pulumi environment variables

### Frontend
- `app/frontend/src/components/shared/cloud-resources-list/cloud-resources-list.tsx` - Minor updates
- `app/frontend/src/app/dashboard/page.tsx` - Fixed TypeScript error

## Key Features Delivered

✅ **Asynchronous Pulumi deployments** with background execution
✅ **Stack job tracking** with status and output storage
✅ **Deployment history** with list and filter capabilities
✅ **Pulumi output capture** (stdout, stderr, exit codes)
✅ **Error handling** with detailed error messages
✅ **10-minute timeout** for deployments
✅ **Docker integration** with Pulumi CLI

## Technical Metrics

- **3 gRPC RPC methods** for stack-update management
- **1 MongoDB collection** (stackupdates)
- **10-minute timeout** for Pulumi deployments
- **Asynchronous execution** via goroutines
- **JSON output format** for flexible result storage
- **Pulumi CLI v3.206.0** integrated

## Benefits

### For End Users
- Cloud resources can be deployed to actual infrastructure
- Asynchronous execution prevents request timeouts
- Deployment history provides visibility
- Full Pulumi output for debugging
- Exit codes indicate success/failure clearly

### For Developers
- Simple gRPC APIs for deployment operations
- Polling pattern for status checking
- Filtering capabilities for deployment history
- Comprehensive error capture from Pulumi
- Scalable asynchronous architecture

## Post-Implementation Fixes

### Frontend TypeScript Proto Generation Path Fix

**Issue:** TypeScript proto files generated in wrong location (`app/backend/frontend/` instead of `app/frontend/`)

**Fix:**
1. Updated `app/backend/apis/buf.gen.yaml`: Changed path from `../frontend/src/gen` to `../../frontend/src/gen`
2. Updated `app/backend/apis/Makefile`: Changed cleanup paths from `../frontend/` to `../../frontend/`

**Result:** Proto generation now correctly targets `app/frontend/src/gen/proto/`

## Related Work

**Built on:**
- Cloud Resource APIs
- Pulumi Integration infrastructure
- Manifest Processing

**Enables:**
- Deployment UI in frontend
- Webhooks for deployment notifications
- Retry logic for failed deployments
- Deployment scheduling
- Multi-environment deployments
- Rollback capability using history

## Known Limitations

- **No cancellation**: Once started, deployments cannot be cancelled
- **No progress updates**: Status only updates on completion
- **Single deployment at a time**: No explicit concurrency control
- **No deployment preview**: Always executes with `--skip-preview`
- **Fixed timeout**: 10-minute timeout may be insufficient for large deployments
- **No retry logic**: Failed deployments must be manually retried

---

**Completion Date:** December 3, 2025
**Status:** ✅ Production Ready
**Location:** `app/backend/internal/service/`, `app/backend/internal/database/`, `app/backend/apis/proto/`

