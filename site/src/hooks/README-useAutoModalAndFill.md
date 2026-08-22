# useAutoModalAndFill Hook

A reusable hook for auto modal open functionality that handles the complex timing and state management needed for reliable auto-clicking on page load.

## Features

- ✅ **Auto-click on page load** - Automatically clicks the first item after a configurable delay
- ✅ **React Strict Mode safe** - Uses global state to survive component remounting
- ✅ **Manual testing** - Provides trigger function for testing
- ✅ **Comprehensive logging** - Detailed console logs for debugging
- ✅ **Configurable delays** - Customizable auto-click timing

## Usage

### Basic Usage

```tsx
import { useAutoModalAndFill } from '../../../hooks/useAutoModalAndFill';

const MyComponent = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleOpenModal = () => {
    setIsModalOpen(true);
  };

  const { triggerAutoClick } = useAutoModalAndFill({
    enabled: true, // Set to false to disable auto-click
    autoClickDelay: 1500, // 1.5 seconds delay
    onAutoClick: handleOpenModal,
    debugPrefix: 'My Component',
  });

  return (
    <div>
      {/* Your component content */}

      {/* Manual test button */}
      <button onClick={triggerAutoClick}>🧪 Test Auto-Click</button>
    </div>
  );
};
```

### Advanced Usage

```tsx
const { triggerAutoClick } = useAutoModalAndFill({
  enabled: AUTO_CLICK_FIRST_ITEM, // Use a constant
  autoClickDelay: 2000, // 2 seconds delay
  onAutoClick: () => {
    // Complex logic for opening modal
    const firstItem = items[0];
    handleItemClick(firstItem);
  },
  debugPrefix: 'Advanced Component',
});
```

## Parameters

| Parameter        | Type         | Default               | Description                                |
| ---------------- | ------------ | --------------------- | ------------------------------------------ |
| `enabled`        | `boolean`    | `true`                | Whether to enable auto-click functionality |
| `autoClickDelay` | `number`     | `1500`                | Delay in milliseconds before auto-clicking |
| `onAutoClick`    | `() => void` | **Required**          | Function to call when auto-clicking        |
| `debugPrefix`    | `string`     | `'Auto Modal & Fill'` | Prefix for console logs                    |

## Returns

| Return Value       | Type         | Description                         |
| ------------------ | ------------ | ----------------------------------- |
| `triggerAutoClick` | `() => void` | Manual trigger function for testing |

## Console Logs

The hook provides comprehensive logging with the following format:

```
🚀 [debugPrefix]: useEffect running
🚀 [debugPrefix]: enabled: true
🚀 [debugPrefix]: Auto-click scheduled in 1500 ms
🚀 [debugPrefix]: Countdown: 2
🚀 [debugPrefix]: Countdown: 1
🚀 [debugPrefix]: GLOBAL TIMER FIRED! Auto-clicking
🚀 [debugPrefix]: onAutoClick called successfully
```

## Migration from Existing Code

### Before (CloudConnections.tsx)

```tsx
// Global timer to survive component remounting
let globalAutoClickTimer: NodeJS.Timeout | null = null;
let globalAutoClickExecuted = false;

// Auto-click the first card on page load
useEffect(() => {
  if (AUTO_CLICK_FIRST_CARD && !globalAutoClickExecuted && !globalAutoClickTimer) {
    globalAutoClickExecuted = true;
    const autoClickDelay = 1500;

    globalAutoClickTimer = setTimeout(() => {
      const firstProvider = providers[0];
      handleCardClick(firstProvider);
      globalAutoClickTimer = null;
    }, autoClickDelay);
  }
}, []);

// Reset global state when component unmounts
useEffect(() => {
  return () => {
    globalAutoClickExecuted = false;
    if (globalAutoClickTimer) {
      clearTimeout(globalAutoClickTimer);
      globalAutoClickTimer = null;
    }
  };
}, []);
```

### After (Using the Hook)

```tsx
import { useAutoModalAndFill } from '../../../hooks/useAutoModalAndFill';

const { triggerAutoClick } = useAutoModalAndFill({
  enabled: AUTO_CLICK_FIRST_CARD,
  autoClickDelay: 1500,
  onAutoClick: () => {
    const firstProvider = providers[0];
    handleCardClick(firstProvider);
  },
  debugPrefix: 'Cloud Connections',
});
```

## Benefits

1. **Reusability** - Use the same logic across multiple components
2. **Maintainability** - Centralized logic in one place
3. **Testing** - Built-in manual trigger for testing
4. **Debugging** - Comprehensive logging
5. **Safety** - Handles React Strict Mode and component remounting
6. **Flexibility** - Configurable delays and conditions

## Future Enhancements

This hook can be extended to include:

- Auto-fill functionality
- More complex timing controls
- Animation coordination
- State persistence across page navigations
