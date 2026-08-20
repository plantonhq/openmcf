'use client';

import React, { useState } from 'react';
import { TabNavigation } from './TabNavigation';
import { AnimatedFormField } from './AnimatedFormField';
import { ReferenceFieldGroup } from './ReferenceFieldGroup';
import { useAutoFill } from '../../../contexts/AutoFillContext';

interface FieldArrayWithTabsProps {
  items: Array<{ value?: string; valueFrom?: { kind?: string; env?: string; name?: string; fieldPath?: string } }>;
  onItemsChange: (items: Array<{ value?: string; valueFrom?: { kind?: string; env?: string; name?: string; fieldPath?: string } }>) => void;
  fieldName: string;
  helpText?: string;
}

export const FieldArrayWithTabs: React.FC<FieldArrayWithTabsProps> = ({
  items,
  onItemsChange,
  fieldName,
  helpText,
}) => {
  const [activeTab, setActiveTab] = useState<'direct' | 'reference'>('direct');
  const { isAnimating, getFieldValue, isFieldAnimating, isFieldCompleted } = useAutoFill();

  const handleDirectValueChange = (index: number, value: string) => {
    const newItems = [...items];
    newItems[index] = { value };
    onItemsChange(newItems);
  };

  const handleReferenceChange = (index: number, field: string, value: string) => {
    const newItems = [...items];
    if (!newItems[index]) {
      newItems[index] = { valueFrom: {} };
    }
    if (!newItems[index].valueFrom) {
      newItems[index].valueFrom = {};
    }
    newItems[index].valueFrom = {
      ...newItems[index].valueFrom,
      [field]: value,
    };
    onItemsChange(newItems);
  };

  const addItem = () => {
    const newItems = [...items];
    if (activeTab === 'direct') {
      newItems.push({ value: '' });
    } else {
      newItems.push({ valueFrom: {} });
    }
    onItemsChange(newItems);
  };

  const removeItem = (index: number) => {
    const newItems = items.filter((_, i) => i !== index);
    onItemsChange(newItems);
  };

  return (
    <TabNavigation
      tabs={[
        {
          id: 'direct',
          label: 'Direct Value',
          icon: (
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              />
            </svg>
          ),
          content: (
            <div className="space-y-3">
              {items.map((item, index) => (
                <div key={index} className="flex items-center gap-2">
                  <AnimatedFormField
                    label=""
                    value={item.value || ''}
                    onChange={(value) => handleDirectValueChange(index, value as string)}
                    type="text"
                    fieldName={`${fieldName}.${index}`}
                    isAnimating={isAnimating}
                    isCurrentField={isFieldAnimating(`${fieldName}.${index}`)}
                    isCompleted={isFieldCompleted(`${fieldName}.${index}`)}
                    animatedValue={getFieldValue(`${fieldName}.${index}`)}
                    helpText={index === 0 ? helpText : undefined}
                  />
                  <button
                    type="button"
                    onClick={() => removeItem(index)}
                    className="p-1 text-red-600 hover:text-red-800"
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
                </div>
              ))}
              <button
                type="button"
                onClick={addItem}
                className="text-white hover:text-white text-sm font-medium flex items-center gap-1"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 4v16m8-8H4"
                  />
                </svg>
                Add Another
              </button>
            </div>
          ),
        },
        {
          id: 'reference',
          label: 'Reference',
          icon: (
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
              />
            </svg>
          ),
          content: (
            <div className="space-y-3">
              <div className="space-y-2">
                {items.map((item, index) => (
                  <ReferenceFieldGroup
                    key={index}
                    cloudResourceKind={item.valueFrom?.kind || ''}
                    environment={item.valueFrom?.env || ''}
                    resourceName={item.valueFrom?.name || ''}
                    fieldPath={item.valueFrom?.fieldPath || ''}
                    onCloudResourceKindChange={(value) => handleReferenceChange(index, 'kind', value)}
                    onEnvironmentChange={(value) => handleReferenceChange(index, 'env', value)}
                    onResourceNameChange={(value) => handleReferenceChange(index, 'name', value)}
                    onFieldPathChange={(value) => handleReferenceChange(index, 'fieldPath', value)}
                    onRemove={() => removeItem(index)}
                    showLabels={index === 0}
                    helpTextMapping={index === 0 && helpText ? { [fieldName]: helpText } : undefined}
                  />
                ))}
                <button
                  type="button"
                  onClick={addItem}
                  className="text-white hover:text-white text-sm font-medium flex items-center gap-1"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 4v16m8-8H4"
                    />
                  </svg>
                  Add Another
                </button>
              </div>
            </div>
          ),
        },
      ]}
      activeTab={activeTab}
      onTabChange={(tabId) => setActiveTab(tabId as 'direct' | 'reference')}
    />
  );
};
