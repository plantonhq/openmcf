'use client';

import { useState, useEffect, useCallback, useRef } from 'react';

interface AutoFillAnimationOptions {
  delay?: number; // Delay between characters (ms)
  pauseBetweenFields?: number; // Pause between fields (ms)
  onComplete?: () => void;
  onFieldComplete?: (fieldName: string) => void;
  onFieldStart?: (fieldName: string) => void; // Called when a field starts animating
}

export interface FieldConfig {
  name: string;
  value: string | boolean | string[];
  order: number;
}

export const useAutoFillAnimation = (options: AutoFillAnimationOptions = {}) => {
  const {
    delay = 50,
    pauseBetweenFields = 300,
    onComplete,
    onFieldComplete,
    onFieldStart,
  } = options;

  const [isAnimating, setIsAnimating] = useState(false);
  const [currentField, setCurrentField] = useState<string | null>(null);
  const [fieldValues, setFieldValues] = useState<Record<string, string | boolean | string[]>>({});
  const [completedFields, setCompletedFields] = useState<Set<string>>(new Set());
  
  const animationRef = useRef<NodeJS.Timeout | null>(null);
  const fieldOrderRef = useRef<FieldConfig[]>([]);

  const startAnimation = useCallback((fields: FieldConfig[]) => {
    console.log('🍎 iOS DEBUG - useAutoFillAnimation: startAnimation called');
    console.log('🍎 iOS DEBUG - Fields received:', fields.length, 'fields');
    console.log('🍎 iOS DEBUG - Current isAnimating state:', isAnimating);
    console.log('🍎 iOS DEBUG - User Agent:', navigator.userAgent);
    console.log('🍎 iOS DEBUG - isIOS:', /iPad|iPhone|iPod/.test(navigator.userAgent));
    
    if (isAnimating) {
      console.log('🍎 iOS DEBUG - Already animating, skipping');
      return;
    }

    // Sort fields by order
    const sortedFields = [...fields].sort((a, b) => a.order - b.order);
    fieldOrderRef.current = sortedFields;

    console.log('🍎 iOS DEBUG - Sorted fields:', sortedFields);
    console.log('🍎 iOS DEBUG - Setting isAnimating to true');
    
    setIsAnimating(true);
    setCompletedFields(new Set());
    setFieldValues({});

    let currentFieldIndex = 0;
    let currentCharIndex = 0;

    const animateNextCharacter = () => {
      console.log('🍎 iOS DEBUG - animateNextCharacter called');
      console.log('🍎 iOS DEBUG - currentFieldIndex:', currentFieldIndex, 'total fields:', sortedFields.length);
      
      if (currentFieldIndex >= sortedFields.length) {
        // Animation complete
        console.log('🍎 iOS DEBUG - Animation complete! All fields processed');
        setIsAnimating(false);
        setCurrentField(null);
        onComplete?.();
        return;
      }

      const field = sortedFields[currentFieldIndex];
      const fieldValue = field.value;
      
      console.log('🍎 iOS DEBUG - Processing field:', field.name);
      console.log('🍎 iOS DEBUG - Field value:', fieldValue, 'type:', typeof fieldValue);
      console.log('🍎 iOS DEBUG - currentCharIndex:', currentCharIndex);

      if (currentCharIndex === 0) {
        // Starting new field
        setCurrentField(field.name);
        onFieldStart?.(field.name); // Notify that field is starting
        
        // Handle different value types
        if (typeof fieldValue === 'string') {
          setFieldValues(prev => ({ ...prev, [field.name]: '' }));
        } else {
          // For non-string values (boolean, array), set immediately
          setFieldValues(prev => ({ ...prev, [field.name]: fieldValue }));
          setCompletedFields(prev => new Set([...prev, field.name]));
          onFieldComplete?.(field.name);
          
          // Move to next field immediately
          currentFieldIndex++;
          currentCharIndex = 0;
          animationRef.current = setTimeout(animateNextCharacter, pauseBetweenFields);
          return;
        }
      }

      // Only animate string values character by character
      if (typeof fieldValue === 'string' && currentCharIndex < fieldValue.length) {
        // Add next character
        const newValue = fieldValue.substring(0, currentCharIndex + 1);
        console.log('🍎 iOS DEBUG - Setting field value:', field.name, '=', newValue);
        setFieldValues(prev => ({ ...prev, [field.name]: newValue }));
        currentCharIndex++;

        console.log('🍎 iOS DEBUG - Scheduling next character in', delay, 'ms');
        animationRef.current = setTimeout(animateNextCharacter, delay);
      } else {
        // Field complete
        console.log('🍎 iOS DEBUG - Field complete:', field.name);
        setCompletedFields(prev => new Set([...prev, field.name]));
        onFieldComplete?.(field.name);
        
        // Move to next field
        currentFieldIndex++;
        currentCharIndex = 0;
        
        animationRef.current = setTimeout(animateNextCharacter, pauseBetweenFields);
      }
    };

    animateNextCharacter();
  }, [delay, pauseBetweenFields, onComplete, onFieldComplete, isAnimating, onFieldStart]);

  const stopAnimation = useCallback(() => {
    if (animationRef.current) {
      clearTimeout(animationRef.current);
      animationRef.current = null;
    }
    setIsAnimating(false);
    setCurrentField(null);
  }, []);

  const resetAnimation = useCallback(() => {
    stopAnimation();
    setFieldValues({});
    setCompletedFields(new Set());
  }, [stopAnimation]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (animationRef.current) {
        clearTimeout(animationRef.current);
      }
    };
  }, []);

  return {
    isAnimating,
    currentField,
    fieldValues,
    completedFields,
    startAnimation,
    stopAnimation,
    resetAnimation,
  };
};
