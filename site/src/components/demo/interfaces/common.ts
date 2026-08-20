/**
 * Common interfaces used across demo components
 */

// Base timestamp interface
export interface Timestamp {
  seconds: number;
  nanos: number;
}

// Version information for API resource metadata
export interface ApiResourceMetadataVersion {
  id: string;
  message: string;
}

// Base metadata interface
export interface BaseMetadata {
  id?: string;
  name: string;
  org?: string;
  env?: string;
  version?: ApiResourceMetadataVersion;
  labels?: { [key: string]: string };
  annotations?: { [key: string]: string };
}

// Extended metadata for API resources (exact match from protobuf)
export interface ApiResourceMetadata {
  name: string;
  slug: string;
  id: string;
  org: string;
  env: string;
  labels: { [key: string]: string };
  annotations: { [key: string]: string };
  tags: string[];
  version?: ApiResourceMetadataVersion;
}


// Lifecycle and audit interfaces
export interface ApiResourceLifecycle {
  isActive: boolean;
}

export interface ApiResourceAuditActor {
  id: string;
  avatar: string;
}

export interface ApiResourceAuditInfo {
  createdBy?: ApiResourceAuditActor;
  createdAt?: Timestamp;
  updatedBy?: ApiResourceAuditActor;
  updatedAt?: Timestamp;
  event: string;
}

export interface ApiResourceAudit {
  specAudit?: ApiResourceAuditInfo;
  statusAudit?: ApiResourceAuditInfo;
}

// Resource reference interface
export interface ApiResourceKindAndNameAndId {
  kind: string;
  name: string;
  id: string;
}
