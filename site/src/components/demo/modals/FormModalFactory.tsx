'use client';

import React from 'react';
import { CredentialModal } from './CredentialModal';

interface FormProps {
  onSubmit: (data: any) => void;
  onCancel: () => void;
  initialData?: any;
}

interface FormModalFactoryProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: any) => void;
  FormComponent: React.ComponentType<FormProps>;
  formProps?: any;
}

export const FormModalFactory: React.FC<FormModalFactoryProps> = ({
  isOpen,
  onClose,
  onSubmit,
  FormComponent,
  formProps = {},
}) => {
  const handleFormSubmit = (data: any) => {
    onSubmit(data);
    onClose();
  };

  const handleFormCancel = () => {
    onClose();
  };

  return (
    <CredentialModal
      isOpen={isOpen}
      onClose={onClose}
    >
      <FormComponent
        {...formProps}
        onSubmit={handleFormSubmit}
        onCancel={handleFormCancel}
      />
    </CredentialModal>
  );
};

// Factory function that creates a modal component
export function createFormModal(
  FormComponent: React.ComponentType<FormProps>
): React.FC<Omit<FormModalFactoryProps, 'FormComponent'>> {
  const ModalComponent = (props: Omit<FormModalFactoryProps, 'FormComponent'>) => (
    <FormModalFactory
      {...props}
      FormComponent={FormComponent}
    />
  );
  
  ModalComponent.displayName = `FormModal(${FormComponent.displayName || FormComponent.name || 'Component'})`;
  
  return ModalComponent;
}
