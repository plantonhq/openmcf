# Demo Interfaces

This directory contains TypeScript interface definitions for demo components in the Planton.ai application.

## Structure

### Common Interfaces (`common.ts`)

- **BaseMetadata**: Base metadata structure for all resources
- **ApiResourceMetadata**: Extended metadata for API resources (exact match from protobuf)
- **ApiResourceMetadataVersion**: Version information for API resource metadata
- **Timestamp**: Base timestamp interface
- **ApiResourceLifecycle**: Resource lifecycle management
- **ApiResourceAuditActor**: Actor information for audit trails
- **ApiResourceAuditInfo**: Audit information with created/updated details
- **ApiResourceAudit**: Complete audit structure with spec and status audit
- **ApiResourceKindAndNameAndId**: Resource reference interface

### AWS Credential Interfaces (`aws-credential.ts`)

- **AwsCredential**: Main AWS credential interface (exact match from original form)
- **AwsCredentialSpec**: AWS credential specification with account details

### StackJob Interfaces (`stackjob.ts`)

- **StackJob**: Main StackJob interface for infrastructure operations
- **StackJobSpec**: StackJob specification and configuration
- **StackJobStatus**: StackJob status and progress tracking
- **StackJobEssentials**: Essential StackJob configuration
- **StackJobFormData**: Form data structure for StackJob creation
- **StackJobStatusDisplay**: UI display helpers for StackJob status
- **IacOperationStatusDisplay**: UI display helpers for IaC operations
- **StackJobProgressEvent**: Progress event tracking
- **StackJobList**: List response structure
- **NewStackJobInput**: Input for creating new StackJobs
- **CloudResource**: Cloud resource interface
- **CloudResourceMutator**: Cloud resource mutation interface
- **IacResource**: Infrastructure as Code resource interface
- **StackJobProgressStatus**: Progress status tracking
- **StackJobStatusIacOperationsState**: IaC operations state
- **StackJobProgressIacOperationState**: Individual operation state
- **StackJobProgressIacOperationSnapshot**: Operation snapshot
- **IacOperationResourceChangesSnapshotRow**: Resource change tracking
- **IacOperationSummarySnapshot**: Operation summary
- **IacOperationDiagnosticMessage**: Diagnostic messages
- **IacProvisioner**: Infrastructure provisioner interface

#### StackJob Enums

- **StackJobOperationType**: Types of StackJob operations
- **StackJobProgressEventType**: Types of progress events
- **IacResourceProgressStatusLabel**: IaC resource status labels
- **IacDiagnosticEventSeverityType**: Diagnostic event severity levels
- **WorkflowExecutionStatus**: Workflow execution status
- **WorkflowExecutionResult**: Workflow execution results

## Usage

```typescript
import {
  AwsCredential,
  StackJob,
  StackJobFormData,
  StackJobOperationType,
  ApiResourceAuditActor,
  ApiResourceAuditInfo,
  ApiResourceMetadataVersion,
} from '../interfaces';

// Use in components - exact match from original form
const credential: AwsCredential = {
  metadata: {
    name: 'my-aws-credential',
    slug: 'my-aws-credential',
    id: 'cred-123',
    org: 'my-org',
    env: 'dev',
    labels: {},
    annotations: {},
    tags: [],
    version: {
      id: 'v1.0.0',
      message: 'Initial version',
    },
  },
  spec: {
    accountId: '123456789012',
    region: 'us-east-1',
    accessKeyId: 'AKIA...',
    secretAccessKey: '...',
  },
};

// Using ApiResourceAuditActor for audit information
const auditActor: ApiResourceAuditActor = {
  id: 'user-123',
  avatar: 'https://example.com/avatar.jpg',
};

const auditInfo: ApiResourceAuditInfo = {
  createdBy: auditActor,
  createdAt: { seconds: 1234567890, nanos: 0 },
  event: 'created',
};

// Using ApiResourceMetadataVersion
const version: ApiResourceMetadataVersion = {
  id: 'v1.0.0',
  message: 'Initial version with AWS credentials',
};

// StackJob usage
const stackJob: StackJobFormData = {
  operationType: StackJobOperationType.UPDATE,
  essentials: {
    stackJobRunnerInfo: {
      kind: 'stack-job-runner',
      name: 'my-runner',
      id: 'runner-123',
    },
  },
};
```

## Integration with Demo Components

The demo components now use these interfaces for:

- **CloudConnections**: Type-safe modal form data handling and AWS credential form submission
- **StackJob Components**: Infrastructure operation tracking and status display
- **Form Components**: Consistent form data structures and validation
- **Future Components**: Ready for additional cloud provider integrations

## Future Extensions

This structure can be extended for other cloud providers:

- `gcp-credential.ts`
- `azure-credential.ts`
- `cloudflare-credential.ts`
- `digitalocean-credential.ts`

Each provider would follow the same pattern with provider-specific account details interfaces.
