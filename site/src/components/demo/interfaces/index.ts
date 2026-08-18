/**
 * Demo interfaces index
 * Exports all interface definitions for demo components
 */

// Common interfaces
export * from './common';

// AWS-specific interfaces
export * from './aws-credential';
export * from './aws-alb';

// StackJob interfaces
export * from './stackjob';

// Re-export commonly used types for convenience
export type {
  BaseMetadata,
  ApiResourceMetadata,
  ApiResourceMetadataVersion,
  Timestamp,
  ApiResourceLifecycle,
  ApiResourceAudit,
  ApiResourceAuditInfo,
  ApiResourceAuditActor,
  ApiResourceKindAndNameAndId,
} from './common';

export type {
  AwsCredential,
  AwsCredentialSpec,
} from './aws-credential';

export type {
  AwsAlb,
  AwsAlbSpec,
  AwsAlbDns,
  AwsAlbSsl,
  AwsAlbStatus,
  AwsAlbStackOutputs,
  StringValueOrRef,
} from './aws-alb';

export type {
  StackJob,
  StackJobSpec,
  StackJobStatus,
  StackJobEssentials,
  StackJobFormData,
  StackJobStatusDisplay,
  IacOperationStatusDisplay,
  StackJobProgressEvent,
  StackJobList,
  NewStackJobInput,
  CloudResource,
  CloudResourceMutator,
  IacResource,
  StackJobProgressStatus,
  StackJobStatusIacOperationsState,
  StackJobProgressIacOperationState,
  StackJobProgressIacOperationSnapshot,
  IacOperationResourceChangesSnapshotRow,
  IacOperationSummarySnapshot,
  IacOperationDiagnosticMessage,
  IacProvisioner,
} from './stackjob';

// Re-export enums
export {
  StackJobOperationType,
  StackJobProgressEventType,
  IacResourceProgressStatusLabel,
  IacDiagnosticEventSeverityType,
  WorkflowExecutionStatus,
  WorkflowExecutionResult,
} from './stackjob';
