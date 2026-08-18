'use client';

import React, { useContext } from 'react';
import { AwsCredential } from '../interfaces';
import { AnimatedFormField } from './AnimatedFormField';
import { AutoFillContext } from '../../../contexts/AutoFillContext';

interface AwsCredentialFormElementsProps {
  formData: AwsCredential;
  onInputChange: (field: string, value: string) => void;
  helpTextMapping?: { [key: string]: string };
}

export const AwsCredentialFormElements: React.FC<AwsCredentialFormElementsProps> = ({
  formData,
  onInputChange,
  helpTextMapping,
}) => {
  // Safely get the auto-fill context - it might not be available in all forms
  const autoFillContext = useContext(AutoFillContext);
  const isAnimating = autoFillContext?.isAnimating || false;
  const getFieldValue = autoFillContext?.getFieldValue || (() => '');
  const isFieldAnimating = autoFillContext?.isFieldAnimating || (() => false);
  const isFieldCompleted = autoFillContext?.isFieldCompleted || (() => false);
  
  const accountIdRef = React.useRef<HTMLDivElement>(null);
  const regionRef = React.useRef<HTMLDivElement>(null);
  const accessKeyIdRef = React.useRef<HTMLDivElement>(null);
  const secretAccessKeyRef = React.useRef<HTMLDivElement>(null);

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
      {/* Account ID */}
      <AnimatedFormField
        ref={accountIdRef}
        label="Account ID"
        value={formData.spec.accountId}
        onChange={(value) => onInputChange('spec.accountId', value as string)}
        helpText={helpTextMapping?.['spec.accountId']}
        fieldName="spec.accountId"
        isAnimating={isAnimating}
        isCurrentField={isFieldAnimating('spec.accountId')}
        isCompleted={isFieldCompleted('spec.accountId')}
        animatedValue={getFieldValue('spec.accountId')}
      />

      {/* Default Region */}
      <AnimatedFormField
        ref={regionRef}
        label="Default Region"
        value={formData.spec.region}
        onChange={(value) => onInputChange('spec.region', value as string)}
        helpText={helpTextMapping?.['spec.region']}
        fieldName="spec.region"
        isAnimating={isAnimating}
        isCurrentField={isFieldAnimating('spec.region')}
        isCompleted={isFieldCompleted('spec.region')}
        animatedValue={getFieldValue('spec.region')}
      />

      {/* Access Key ID */}
      <AnimatedFormField
        ref={accessKeyIdRef}
        label="Access Key ID"
        value={formData.spec.accessKeyId}
        onChange={(value) => onInputChange('spec.accessKeyId', value as string)}
        helpText={helpTextMapping?.['spec.accessKeyId']}
        fieldName="spec.accessKeyId"
        isAnimating={isAnimating}
        isCurrentField={isFieldAnimating('spec.accessKeyId')}
        isCompleted={isFieldCompleted('spec.accessKeyId')}
        animatedValue={getFieldValue('spec.accessKeyId')}
      />

      {/* Secret Access Key */}
      <AnimatedFormField
        ref={secretAccessKeyRef}
        label="Secret Access Key"
        value={formData.spec.secretAccessKey}
        onChange={(value) => onInputChange('spec.secretAccessKey', value as string)}
        type="password"
        helpText={helpTextMapping?.['spec.secretAccessKey']}
        fieldName="spec.secretAccessKey"
        isAnimating={isAnimating}
        isCurrentField={isFieldAnimating('spec.secretAccessKey')}
        isCompleted={isFieldCompleted('spec.secretAccessKey')}
        animatedValue={getFieldValue('spec.secretAccessKey')}
      />
    </div>
  );
};
