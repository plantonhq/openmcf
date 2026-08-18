/**
 * Animation configuration constants for auto-fill forms
 * Adjust these values to change the speed and timing of form animations
 */

// Character typing speed (milliseconds between each character)
export const ANIMATION_DELAY = 50; // 50ms = moderate typing speed

// Pause between fields (milliseconds)
export const ANIMATION_PAUSE_BETWEEN_FIELDS = 300; // 300ms = natural pause

// Scroll delay (milliseconds to wait before scrolling to field)
export const ANIMATION_SCROLL_DELAY = 100; // 100ms = ensures field is rendered

// Auto-start delay (milliseconds to wait before starting animation)
export const ANIMATION_AUTO_START_DELAY = 1000; // 1000ms = 1 second delay

// Header height for scroll calculations (pixels)
export const ANIMATION_HEADER_HEIGHT = 80; // 80px = approximate header height

/**
 * Preset configurations for different animation speeds
 */
export const ANIMATION_PRESETS = {
  // Very fast animation for quick demos
  FAST: {
    delay: 25,
    pauseBetweenFields: 150,
    scrollDelay: 50,
    autoStartDelay: 500,
    headerHeight: 80,
  },
  
  // Normal speed (current default)
  NORMAL: {
    delay: 50,
    pauseBetweenFields: 300,
    scrollDelay: 100,
    autoStartDelay: 1000,
    headerHeight: 80,
  },
  
  // Slower animation for detailed demos
  SLOW: {
    delay: 100,
    pauseBetweenFields: 500,
    scrollDelay: 150,
    autoStartDelay: 1500,
    headerHeight: 80,
  },
  
  // Very slow for step-by-step demonstrations
  ULTRA_SLOW: {
    delay: 200,
    pauseBetweenFields: 1000,
    scrollDelay: 200,
    autoStartDelay: 2000,
    headerHeight: 80,
  },
} as const;

// Current preset being used
export const CURRENT_PRESET = ANIMATION_PRESETS.NORMAL;
