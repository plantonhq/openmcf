/**
 * Demo modals index
 * Exports all modal components for demo
 */

// Core modal components
export { CredentialModal } from './CredentialModal';
export { FormModalFactory, createFormModal } from './FormModalFactory';
export { FormModalRegistry } from './FormModalRegistry';


// Specific credential modals
export { AwsCredentialModal } from './AwsCredentialModal';
export { AwsAlbModal } from './AwsAlbModal';

// Legacy component (kept for backward compatibility)
export { CloudProviderModal } from './CloudProviderModal';
