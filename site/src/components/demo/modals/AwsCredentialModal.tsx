'use client';

import React from 'react';
import { AwsCredential } from '../interfaces';
import { AwsCredentialForm } from '../forms';
import { FormModalFactory } from './FormModalFactory';

interface AwsCredentialModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: AwsCredential) => void;
  initialData?: Partial<AwsCredential>;
}

export const AwsCredentialModal: React.FC<AwsCredentialModalProps> = ({
  isOpen,
  onClose,
  onSubmit,
  initialData,
}) => {
  return (
    <FormModalFactory
      isOpen={isOpen}
      onClose={onClose}
      onSubmit={onSubmit}
      FormComponent={AwsCredentialForm}
      formProps={{ initialData }}
    />
  );
};