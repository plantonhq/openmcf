import { ApiResourceMetadata } from './common';

/**
 * AWS Credential interface
 * Exact match from the original AWS credential form
 */

export interface AwsCredentialSpec {
  accountId: string;
  region: string;
  accessKeyId: string;
  secretAccessKey: string;
}

export interface AwsCredential {
  metadata: ApiResourceMetadata;
  spec: AwsCredentialSpec;
}

