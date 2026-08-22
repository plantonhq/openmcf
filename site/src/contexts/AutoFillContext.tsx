'use client';

import React, { createContext, useContext, ReactNode } from 'react';
import { useAutoFillAnimation, FieldConfig } from '../hooks/useAutoFillAnimation';
import { CURRENT_PRESET } from '../constants/animationConfig';

interface AutoFillContextType {
  isAnimating: boolean;
  currentField: string | null;
  fieldValues: Record<string, string | boolean | string[]>;
  completedFields: Set<string>;
  startAnimation: (fields: FieldConfig[]) => void;
  stopAnimation: () => void;
  resetAnimation: () => void;
  getFieldValue: (fieldName: string) => string | boolean | string[];
  isFieldAnimating: (fieldName: string) => boolean;
  isFieldCompleted: (fieldName: string) => boolean;
  scrollToField: (fieldName: string) => void;
  registerFieldRef: (fieldName: string, element: HTMLDivElement | null) => void;
}

export const AutoFillContext = createContext<AutoFillContextType | undefined>(undefined);

interface AutoFillProviderProps {
  children: ReactNode;
  delay?: number;
  pauseBetweenFields?: number;
  onComplete?: () => void;
  onFieldComplete?: (fieldName: string) => void;
}

export const AutoFillProvider: React.FC<AutoFillProviderProps> = ({
  children,
  delay = 50,
  pauseBetweenFields = 300,
  onComplete,
  onFieldComplete,
}) => {
  const fieldRefs = React.useRef<Record<string, HTMLDivElement | null>>({});

  const scrollToField = React.useCallback((fieldName: string) => {
    const fieldElement = fieldRefs.current[fieldName];
    if (fieldElement) {
      // Get the form container to ensure we don't scroll past the header
      const formContainer = fieldElement.closest('.overflow-y-auto');
      
      if (formContainer) {
        // Calculate the position to scroll to, ensuring header stays visible
        const containerRect = formContainer.getBoundingClientRect();
        const fieldRect = fieldElement.getBoundingClientRect();
        const headerHeight = CURRENT_PRESET.headerHeight;
        
        // Calculate scroll position to keep field visible but header in view
        const scrollTop = formContainer.scrollTop + (fieldRect.top - containerRect.top) - headerHeight;
        
        formContainer.scrollTo({
          top: Math.max(0, scrollTop), // Don't scroll above the top
          behavior: 'smooth'
        });
      } else {
        // Fallback to original behavior if container not found
        fieldElement.scrollIntoView({
          behavior: 'smooth',
          block: 'start', // Changed from 'center' to 'start' to avoid cutting header
          inline: 'nearest'
        });
      }
    }
  }, []);

  const handleFieldStart = React.useCallback((fieldName: string) => {
    // Small delay to ensure the field is rendered before scrolling
    setTimeout(() => {
      scrollToField(fieldName);
    }, CURRENT_PRESET.scrollDelay);
  }, [scrollToField]);

  const handleComplete = React.useCallback(() => {
    // Scroll to top when animation completes to show header
    const formContainer = document.querySelector('.overflow-y-auto');
    if (formContainer) {
      formContainer.scrollTo({
        top: 0,
        behavior: 'smooth'
      });
    }
    onComplete?.();
  }, [onComplete]);

  const {
    isAnimating,
    currentField,
    fieldValues,
    completedFields,
    startAnimation,
    stopAnimation,
    resetAnimation,
  } = useAutoFillAnimation({
    delay,
    pauseBetweenFields,
    onComplete: handleComplete,
    onFieldComplete,
    onFieldStart: handleFieldStart,
  });

  const getFieldValue = (fieldName: string): string | boolean | string[] => {
    return fieldValues[fieldName] || '';
  };

  const isFieldAnimating = (fieldName: string): boolean => {
    return isAnimating && currentField === fieldName;
  };

  const isFieldCompleted = (fieldName: string): boolean => {
    return completedFields.has(fieldName);
  };

  const registerFieldRef = React.useCallback((fieldName: string, element: HTMLDivElement | null) => {
    fieldRefs.current[fieldName] = element;
  }, []);

  const contextValue: AutoFillContextType = {
    isAnimating,
    currentField,
    fieldValues,
    completedFields,
    startAnimation,
    stopAnimation,
    resetAnimation,
    getFieldValue,
    isFieldAnimating,
    isFieldCompleted,
    scrollToField,
    registerFieldRef,
  };

  return (
    <AutoFillContext.Provider value={contextValue}>
      {children}
    </AutoFillContext.Provider>
  );
};

export const useAutoFill = (): AutoFillContextType => {
  const context = useContext(AutoFillContext);
  if (context === undefined) {
    throw new Error('useAutoFill must be used within an AutoFillProvider. Make sure the component using this hook is wrapped with <AutoFillProvider>.');
  }
  return context;
};
