'use client';

import React, { useState, useEffect, useContext } from 'react';
import { AnimatedFormField } from './AnimatedFormField';
import { AutoFillContext } from '../../../contexts/AutoFillContext';

interface MetadataSectionProps {
  name: string;
  slug: string;
  onNameChange: (value: string) => void;
  onSlugChange: (value: string) => void;
  nameHelpText?: string;
  slugHelpText?: string;
}

export const MetadataSection: React.FC<MetadataSectionProps> = ({
  name,
  slug,
  onNameChange,
  onSlugChange,
  nameHelpText,
  slugHelpText,
}) => {
  const [isSlugManuallyEdited, setIsSlugManuallyEdited] = useState(false);

  // Safely get the auto-fill context - it might not be available in all forms
  const autoFillContext = useContext(AutoFillContext);
  const isAnimating = autoFillContext?.isAnimating || false;
  const getFieldValue = autoFillContext?.getFieldValue || (() => '');
  const isFieldAnimating = autoFillContext?.isFieldAnimating || (() => false);
  const isFieldCompleted = autoFillContext?.isFieldCompleted || (() => false);

  const nameFieldRef = React.useRef<HTMLDivElement>(null);
  const slugFieldRef = React.useRef<HTMLDivElement>(null);

  // Auto-generate slug when name changes (only if slug hasn't been manually edited)
  useEffect(() => {
    if (!isSlugManuallyEdited && name) {
      const generatedSlug = name
        .toLowerCase()
        .replace(/[^a-z0-9\s-]/g, '-') // Replace special characters with hyphens
        .replace(/\s+/g, '-') // Replace spaces with hyphens
        .replace(/-+/g, '-') // Replace multiple hyphens with single hyphen
        .replace(/^-|-$/g, ''); // Remove leading/trailing hyphens

      onSlugChange(generatedSlug);
    }
  }, [name, isSlugManuallyEdited, onSlugChange]);

  const handleNameChange = (value: string | boolean | string[]) => {
    // For name field, we only expect string values
    if (typeof value === 'string') {
      onNameChange(value);
    }
  };

  const handleSlugChange = (value: string | boolean | string[]) => {
    // For slug field, we only expect string values
    if (typeof value === 'string') {
      onSlugChange(value);

      // If user clears the slug (makes it blank), resume auto-generation
      if (value.trim() === '') {
        setIsSlugManuallyEdited(false);
      } else {
        // Mark slug as manually edited when user changes it to non-empty value
        setIsSlugManuallyEdited(true);
      }
    }
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
      <AnimatedFormField
        ref={nameFieldRef}
        label="Name"
        value={name}
        onChange={handleNameChange}
        helpText={nameHelpText}
        fieldName="metadata.name"
        isAnimating={isAnimating}
        isCurrentField={isFieldAnimating('metadata.name')}
        isCompleted={isFieldCompleted('metadata.name')}
        animatedValue={getFieldValue('metadata.name')}
      />
      <AnimatedFormField
        ref={slugFieldRef}
        label="Slug"
        value={slug}
        onChange={handleSlugChange}
        helpText={slugHelpText}
        fieldName="metadata.slug"
        isAnimating={isAnimating}
        isCurrentField={isFieldAnimating('metadata.slug')}
        isCompleted={isFieldCompleted('metadata.slug')}
        animatedValue={getFieldValue('metadata.slug')}
      />
    </div>
  );
};
