# T07: Stack Jobs UI Integration and Backend Pagination

**Status:** ✅ COMPLETED
**Date:** December 4, 2025
**Type:** Feature, Enhancement
**Changelog:** `2025-12-04-121750-stack-update-ui-integration-and-pagination.md`

---

## Overview

Integrated stack-updates functionality into the cloud resources web interface, enabling users to view and navigate stack-updates directly from the cloud resources list. Added server-side pagination to ListStackUpdates API, enhanced DeployCloudResource API to accept user-provided credentials, and fixed module directory path resolution.

## What Was Accomplished

### 1. Stack Jobs Menu Integration

Added "Stack Jobs" action menu that:

- Opens drawer showing paginated stack-updates for selected cloud resource
- Displays stack-updates in sortable table
- Provides clickable rows that navigate to detailed pages
- Shows job status with color-coded chips

### 2. Stack Jobs Detail Page

Created dedicated page (`/stack-update/[id]`) with:

- Breadcrumb navigation with clickable "Stack Jobs" link
- Job ID with copy-to-clipboard functionality
- Status chip for visual status indication
- Last updated timestamp
- Full JSON output with syntax highlighting
- Loading states with skeleton placeholders

### 3. Backend Stack Jobs Pagination

Implemented server-side pagination for ListStackUpdates:

- Added `PageInfo` support to request message
- Added `total_pages` field to response
- Repository supports pagination with MongoDB skip/limit
- Total pages calculated using ceiling division
- Default pagination: page 0, size 20 if not provided
- Backward compatible (pagination is optional)

### 4. User-Provided Credentials Support

Enhanced DeployCloudResource API to accept provider credentials:

- Added `ProviderConfig` message supporting all providers (AWS, GCP, Azure, Atlas, Cloudflare, Confluent, Snowflake, Kubernetes)
- Credentials priority: User-provided > Environment variables
- Automatic credential validation based on resource provider
- Temporary credential files created from API request
- Automatic cleanup of temporary files after deployment
- Clear error messages for missing credentials

### 5. Module Path Fixes

Fixed module directory path resolution:

- **Pulumi**: Corrected from `apis/org/openmcf/provider` to `apis/project/planton/provider`
- **OpenTofu**: Corrected from `apis/org/openmcf/provider` to `apis/project/planton/provider`
- Fixed version check logic in Pulumi module directory
- Improved version handling in module directory logic

### 6. Code Quality Improvements

Removed redundant logrus logging from service layer:

- Removed from `stack_job_service.go`
- Removed from `cloud_resource_service.go`
- Removed from `deployment_component_service.go`
- Errors still properly returned via Connect RPC

## Technical Implementation

### User Flow

```
Cloud Resources List Page
    ↓ User clicks "Stack Jobs" menu item
Stack Jobs Drawer (opens)
    ↓ Shows paginated list of stack-updates
    ↓ User clicks on a stack-update row
Stack Job Detail Page (/stack-update/[id])
    ↓ Shows full stack-update details
    ↓ Can navigate back via breadcrumb
```

### Backend API Flow

```
Frontend Request (ListStackUpdates)
    ↓ With PageInfo (page, size)
Backend Service (ListStackUpdates)
    ↓ Applies filters and pagination
Repository Layer
    ↓ MongoDB query with skip/limit
    ↓ Count query for total pages
Response (jobs + totalPages)
```

### Credential Flow

```
DeployCloudResource Request
    ↓ With ProviderConfig (optional)
Backend Service
    ↓ Priority: User credentials > Env vars
    ↓ Create temporary credential files
    ↓ Build provider config from credentials
    ↓ Validate based on resource provider
Pulumi Deployment
    ↓ Uses provided credentials
Cleanup temporary files
```

## Files Created

### Frontend Pages

- `app/frontend/src/app/stack-update/[id]/page.tsx` - Stack job detail page
- `app/frontend/src/app/stack-update/_services/index.ts` - Service exports
- `app/frontend/src/app/stack-update/_services/query.ts` - Stack jobs query service
- `app/frontend/src/app/stack-update/styled.ts` - Styled components

### UI Components

- `app/frontend/src/components/shared/stackupdate/index.ts` - Stack job component exports
- `app/frontend/src/components/shared/stackupdate/stack-update-header.tsx` - Job header component
- `app/frontend/src/components/shared/stackupdate/stack-update-drawer.tsx` - Drawer component
- `app/frontend/src/components/shared/stackupdate/stack-update-list.tsx` - List component with pagination
- `app/frontend/src/components/shared/breadcrumb/index.tsx` - Breadcrumb navigation
- `app/frontend/src/components/shared/breadcrumb/styled.ts` - Breadcrumb styling
- `app/frontend/src/components/shared/status-chip/index.ts` - Status chip exports
- `app/frontend/src/components/shared/status-chip/status-chip.tsx` - Status chip component
- `app/frontend/src/components/shared/syntax-highlighter/index.ts` - Syntax highlighter exports
- `app/frontend/src/components/shared/syntax-highlighter/json-code.tsx` - JSON syntax highlighter

### Infrastructure

- `pkg/iac/stackinput/stackinputproviderconfig/user_provider.go` - User credentials handling

## Files Modified

### Backend API

- `app/backend/apis/proto/stack_job_service.proto` - Added PageInfo and ProviderConfig support
- `app/backend/internal/service/stack_job_service.go` - Pagination logic, user credentials support, removed logrus
- `app/backend/internal/service/cloud_resource_service.go` - Removed logrus logging
- `app/backend/internal/service/deployment_component_service.go` - Removed logrus logging
- `app/backend/internal/database/stack_job_repo.go` - Added pagination support

### Infrastructure Code

- `pkg/iac/pulumi/pulumimodule/module_directory.go` - Fixed API path and version check
- `pkg/iac/tofu/tofumodule/module_directory.go` - Fixed API path

### Frontend Components

- `app/frontend/src/components/shared/cloud-resources-list/cloud-resources-list.tsx` - Added Stack Jobs menu
- `app/frontend/src/components/shared/cloud-resources-list/index.ts` - Updated exports
- `app/frontend/src/components/layout/styled.ts` - Updated layout styling
- `app/frontend/src/components/shared/drawer/styled.ts` - Updated drawer styling

## Files Deleted

### Infrastructure

- `pkg/iac/stackinput/stackinputproviderconfig/env_provider.go` - Functionality consolidated into user_provider.go

### Frontend

- `app/frontend/src/components/shared/cloud-resources-list/styled.ts` - No longer needed

## Key Features Delivered

✅ **Stack Jobs UI integration** with cloud resources list
✅ **Stack Jobs detail page** with full job information
✅ **Server-side pagination** for stack-updates
✅ **User-provided credentials** support in deploy API
✅ **Module path fixes** for Pulumi and OpenTofu
✅ **Breadcrumb navigation** for better UX
✅ **JSON syntax highlighting** for outputs
✅ **Status chips** for visual feedback

## Technical Metrics

- **1 new detail page** with dynamic routing
- **6 new reusable components** for stack-updates UI
- **Server-side pagination** in backend and frontend
- **8 cloud providers** supported for user credentials
- **Default page size**: 10 items per page (frontend), 20 (backend)
- **Full TypeScript coverage** for all components
- **Module path fixes** for both IaC engines

## Benefits

### For End Users

- Stack jobs accessible directly from cloud resources
- Intuitive navigation flow from resources to jobs to details
- Paginated list prevents performance issues
- Detailed view shows complete deployment output
- Status chips provide quick visual feedback
- Can provide credentials per deployment via API
- No need to pre-configure credentials on server

### For Developers

- Reusable stack-updates components
- Consistent pagination pattern
- Type-safe implementation with TypeScript
- Scalable server-side pagination
- Fixed module path resolution for reliable deployments
- Provider credential support for all major clouds

## Post-Implementation Cleanup

### Logrus Removal from Service APIs

Removed all logrus logging from backend service layer:

- Most error logs were redundant (errors returned via Connect RPC)
- Info/warning logs were nice-to-have but not essential
- Simplifies codebase and reduces dependencies
- Errors still properly handled and returned to clients

## Known Limitations

- **Fixed page size**: Frontend uses 10 items per page, not user-configurable
- **No status filtering in UI**: Backend supports it but UI doesn't expose it
- **No real-time updates**: Status changes require manual refresh
- **No deployment actions**: Can't trigger deployments from UI
- **No output filtering**: Full JSON output always displayed
- **No credential encryption**: Credentials passed in plain text (should use encryption in production)

## Related Work

**Built on:**

- Pulumi CLI Stack Job API (Dec 3, 2025)
- Cloud Resource UI Enhancements (Dec 3, 2025)
- Cloud Resource Web UI (Dec 1, 2025)

**Enables:**

- Deployment actions from UI
- Real-time status updates
- Stack job filtering in UI
- Export functionality
- Bulk operations
- Search functionality

---

**Completion Date:** December 4, 2025
**Status:** ✅ Production Ready
**Location:** `app/frontend/src/app/stack-update/`, `app/backend/internal/service/`, `pkg/iac/`
