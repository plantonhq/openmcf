'use client';

import React from 'react';
import { FormField } from './FormField';

interface ReferenceFieldGroupProps {
  cloudResourceKind: string;
  environment: string;
  resourceName: string;
  fieldPath: string;
  onCloudResourceKindChange: (value: string) => void;
  onEnvironmentChange: (value: string) => void;
  onResourceNameChange: (value: string) => void;
  onFieldPathChange: (value: string) => void;
  onRemove?: () => void;
  showLabels?: boolean;
  helpTextMapping?: Record<string, string>;
}

export const ReferenceFieldGroup: React.FC<ReferenceFieldGroupProps> = ({
  cloudResourceKind,
  environment,
  resourceName,
  fieldPath,
  onCloudResourceKindChange,
  onEnvironmentChange,
  onResourceNameChange,
  onFieldPathChange,
  onRemove,
  showLabels = true,
  helpTextMapping,
}) => {
  return (
    <div className="flex items-end gap-2">
      <div className="flex-1 grid grid-cols-4 gap-2">
        <FormField
          label={showLabels ? "Cloud Resource Kind" : ""}
          value={cloudResourceKind}
          onChange={onCloudResourceKindChange}
          placeholder="Select kind"
          helpText={showLabels ? helpTextMapping?.['spec.securityGroups.kind'] : undefined}
        />
        <FormField
          label={showLabels ? "Environment" : ""}
          value={environment}
          onChange={onEnvironmentChange}
          placeholder="Select environment"
          helpText={showLabels ? helpTextMapping?.['spec.securityGroups.env'] : undefined}
        />
        <FormField
          label={showLabels ? "Resource Name" : ""}
          value={resourceName}
          onChange={onResourceNameChange}
          placeholder="Select resource"
          helpText={showLabels ? helpTextMapping?.['spec.securityGroups.name'] : undefined}
        />
        <FormField
          label={showLabels ? "Field Path" : ""}
          value={fieldPath}
          onChange={onFieldPathChange}
          placeholder="Enter field path"
          helpText={showLabels ? helpTextMapping?.['spec.securityGroups.fieldPath'] : undefined}
        />
      </div>
      {onRemove && (
        <button
          type="button"
          onClick={onRemove}
          className="p-1 text-red-600 hover:text-red-800 flex-shrink-0 mb-2"
        >
          <svg
            className="w-4 h-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
            />
          </svg>
        </button>
      )}
    </div>
  );
};
