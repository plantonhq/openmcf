# Demo Modals

This directory contains modal components for the demo application, providing a comprehensive and reusable modal system for different cloud provider configurations.

## Architecture

The modal system is built with a layered architecture for maximum reusability:

```
┌─────────────────────────────────────────────────────────┐
│ FormModalRegistry (Entry Point)                  │
├─────────────────────────────────────────────────────────┤
│ FormModalFactory (Factory Pattern)               │
├─────────────────────────────────────────────────────────┤
│ CredentialModal (Base Modal)                           │
├─────────────────────────────────────────────────────────┤
│ Specific Modals (AWS, GCP, Azure, etc.)                │
└─────────────────────────────────────────────────────────┘
```

## Core Components

### CredentialModal

The base modal component that provides the core modal functionality.

**Features:**

- **Full-screen modal** with proper animations
- **Click outside to close** functionality
- **Keyboard navigation** support
- **Consistent styling** across all modals

**Props:**

```typescript
interface CredentialModalProps {
  isOpen: boolean; // Controls modal visibility
  onClose: () => void; // Called when modal should close
  onSubmit: (data: any) => void; // Called with form data on submit
  providerName: string; // Provider name
  providerDescription: string; // Provider description
  providerIcon?: React.ReactNode; // Provider icon component
  children: React.ReactNode; // Modal content
  title?: string; // Optional custom title
  submitLabel?: string; // Custom submit button text
  cancelLabel?: string; // Custom cancel button text
}
```

### FormModalFactory

A factory component that creates modals for any form.

**Features:**

- **Generic form integration** - works with any form component
- **Automatic form handling** - manages form submission and cancellation
- **Props forwarding** - passes through form-specific props
- **Consistent behavior** - all modals behave the same way

**Props:**

```typescript
interface FormModalFactoryProps {
  isOpen: boolean; // Controls modal visibility
  onClose: () => void; // Called when modal should close
  onSubmit: (data: any) => void; // Called with form data on submit
  providerName: string; // Provider name
  providerDescription: string; // Provider description
  providerIcon?: React.ReactNode; // Provider icon component
  FormComponent: React.ComponentType<FormProps>; // Form component
  formProps?: any; // Props to pass to form component
  title?: string; // Optional custom title
  submitLabel?: string; // Custom submit button text
  cancelLabel?: string; // Custom cancel button text
}
```

### FormModalRegistry

The main entry point that automatically selects the appropriate modal based on provider ID.

**Features:**

- **Automatic provider detection** - selects the right modal based on provider ID
- **Fallback support** - provides generic modal for unsupported providers
- **Easy extension** - adding new providers is simple
- **Type safety** - TypeScript support for supported providers

**Props:**

```typescript
interface FormModalRegistryProps {
  providerId: string; // Provider identifier (aws, gcp, azure, etc.)
  isOpen: boolean; // Controls modal visibility
  onClose: () => void; // Called when modal should close
  onSubmit: (data: any) => void; // Called with form data on submit
  providerName: string; // Provider name
  providerDescription: string; // Provider description
  providerIcon?: React.ReactNode; // Provider icon component
  initialData?: any; // Initial data for the form
}
```

## Specific Credential Modals

### AwsCredentialModal

Specialized modal for AWS credential configuration.

**Features:**

- **Integrated AWS form** with all features
- **Dynamic data loading** from JSON
- **Full-screen layout** with proper scrolling
- **Summary sidebar** with progress indicators

### GcpCredentialModal

Modal for Google Cloud Platform credentials (placeholder implementation).

**Features:**

- **Placeholder form** ready for implementation
- **Consistent styling** with other modals
- **Easy to extend** with actual GCP form

### AzureCredentialModal

Modal for Microsoft Azure credentials (placeholder implementation).

**Features:**

- **Placeholder form** ready for implementation
- **Consistent styling** with other modals
- **Easy to extend** with actual Azure form

## Usage Examples

### Using the Registry (Recommended)

```typescript
import { FormModalRegistry } from '../modals';

<FormModalRegistry
  providerId="aws"
  isOpen={isModalOpen}
  onClose={() => setIsModalOpen(false)}
  onSubmit={(data) => handleCredentials(data)}
  providerName="AWS"
  providerDescription="Amazon Web Services"
  providerIcon={<AwsIcon />}
/>;
```

### Using Specific Modals

```typescript
import { AwsCredentialModal } from '../modals';

<AwsCredentialModal
  isOpen={isModalOpen}
  onClose={() => setIsModalOpen(false)}
  onSubmit={(data) => handleAwsCredentials(data)}
  providerName="AWS"
  providerDescription="Amazon Web Services"
/>;
```

### Using the Factory

```typescript
import { FormModalFactory } from '../modals';

<FormModalFactory
  isOpen={isModalOpen}
  onClose={() => setIsModalOpen(false)}
  onSubmit={(data) => handleCredentials(data)}
  providerName="Custom Provider"
  providerDescription="Custom provider description"
  FormComponent={CustomCredentialForm}
  formProps={{ customProp: 'value' }}
/>;
```

### Using Higher-Order Component

```typescript
import { withCredentialModal } from '../modals';

const ModalizedForm = withCredentialModal(MyCredentialForm);

<ModalizedForm
  isOpen={isModalOpen}
  onClose={() => setIsModalOpen(false)}
  onSubmit={(data) => handleCredentials(data)}
  providerName="Custom Provider"
  providerDescription="Custom provider description"
  // ... other form props
/>;
```

## Adding New Providers

### 1. Create the Form Component

```typescript
// NewProviderCredentialForm.tsx
const NewProviderCredentialForm: React.FC<{
  onSubmit: (data: any) => void;
  onCancel: () => void;
  initialData?: any;
}> = ({ onSubmit, onCancel, initialData }) => {
  // Implement your form here
  return <div className="h-full flex flex-col">{/* Your form content */}</div>;
};
```

### 2. Create the Modal Component

```typescript
// NewProviderCredentialModal.tsx
import { FormModalFactory } from './FormModalFactory';
import { NewProviderCredentialForm } from './NewProviderCredentialForm';

export const NewProviderCredentialModal: React.FC<{
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: any) => void;
  providerName?: string;
  providerDescription?: string;
  providerIcon?: React.ReactNode;
  initialData?: any;
}> = ({
  isOpen,
  onClose,
  onSubmit,
  providerName,
  providerDescription,
  providerIcon,
  initialData,
}) => {
  return (
    <FormModalFactory
      isOpen={isOpen}
      onClose={onClose}
      onSubmit={onSubmit}
      providerName={providerName}
      providerDescription={providerDescription}
      providerIcon={providerIcon}
      FormComponent={NewProviderCredentialForm}
      formProps={{ initialData }}
      title={`Create ${providerName} Credentials`}
    />
  );
};
```

### 3. Register the Provider

```typescript
// In FormModalRegistry.tsx
const CREDENTIAL_MODALS = {
  aws: AwsCredentialModal,
  gcp: GcpCredentialModal,
  azure: AzureCredentialModal,
  newprovider: NewProviderCredentialModal, // Add your new provider
} as const;
```

### 4. Export the Component

```typescript
// In index.ts
export { NewProviderCredentialModal } from './NewProviderCredentialModal';
```

## File Structure

```
modals/
├── CredentialModal.tsx              # Base modal component
├── FormModalFactory.tsx       # Factory for creating modals
├── FormModalRegistry.tsx      # Registry for provider selection
├── withCredentialModal.tsx          # Higher-order component
├── AwsCredentialModal.tsx           # AWS-specific modal
├── GcpCredentialModal.tsx           # GCP-specific modal
├── AzureCredentialModal.tsx         # Azure-specific modal
├── CloudProviderModal.tsx           # Legacy generic modal
├── index.ts                         # Export index
└── README.md                        # This documentation
```

## Benefits

### Reusability

- **Single modal system** for all credential types
- **Consistent behavior** across all providers
- **Easy to extend** with new providers

### Maintainability

- **Centralized modal logic** in base components
- **Type safety** with TypeScript
- **Clear separation** of concerns

### Developer Experience

- **Simple API** for common use cases
- **Flexible options** for custom implementations
- **Comprehensive documentation** and examples

## Dependencies

- **React**: For component logic
- **Framer Motion**: For animations
- **TypeScript**: For type safety
- **Tailwind CSS**: For styling
- **Form Components**: For modal content
